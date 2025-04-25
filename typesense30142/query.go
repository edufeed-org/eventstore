package typesense30142

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

func (ts *TSBackend) QueryEvents(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error) {
	ch := make(chan *nostr.Event)

	log.Printf("Processing query with search: %s", filter.Search)
	
	// If we have no search parameter, return an empty channel
	if filter.Search == "" {
		log.Printf("No search parameter provided, returning empty result")
		close(ch)
		return ch, nil
	}

	nostrsearch, err := ts.SearchResources(filter.Search)
	if err != nil {
		log.Printf("Search failed: %v", err)
		// Return the channel anyway, but close it immediately
		close(ch)
		return ch, fmt.Errorf("search failed: %w", err)
	}

	log.Printf("Search succeeded, found %d events", len(nostrsearch))

	go func() {
		// Check if context is done before sending events
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled before sending results")
			close(ch)
			return
		default:
			for _, evt := range nostrsearch {
				select {
				case <-ctx.Done():
					// Context was cancelled during event sending
					log.Printf("Context cancelled during event sending")
					break
				default:
					// Send the event
					ch <- &evt
				}
			}
			close(ch)
		}
	}()
	
	return ch, nil
}

// searches for resources and returns both the AMB metadata and converted Nostr events
func (ts *TSBackend) SearchResources(searchStr string) ([]nostr.Event, error) {
	parsedQuery := ParseSearchQuery(searchStr)

	mainQuery, params, err := BuildTypesenseQuery(parsedQuery)
	if err != nil {
		return nil, fmt.Errorf("error building Typesense query: %v", err)
	}

	// URL encode the main query
	encodedQuery := url.QueryEscape(mainQuery)

	// Default fields to search in
	queryBy := "name,description,about,learningResourceType,keywords,creator,publisher"

	// Start building the search URL
	searchURL := fmt.Sprintf("%s/collections/%s/documents/search?validate_field_names=false&q=%s&query_by=%s",
		ts.Host, ts.CollectionName, encodedQuery, queryBy)

	// Add additional parameters
	for key, value := range params {
		searchURL += fmt.Sprintf("&%s=%s", key, url.QueryEscape(value))
	}

	// Debug information
	fmt.Printf("Search URL: %s\n", searchURL)

	resp, body, err := ts.makehttpRequest(searchURL, http.MethodGet, nil)

	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Try to parse the raw JSON to understand its structure
	var rawResponse interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		fmt.Printf("Warning: Could not parse raw response as JSON: %v\n", err)
	} else {
		// Check if we got a hits array
		responseMap, ok := rawResponse.(map[string]interface{})
		if ok {
			if hits, exists := responseMap["hits"]; exists {
				hitsArray, ok := hits.([]interface{})
				if ok {
					fmt.Printf("Response contains %d hits\n", len(hitsArray))
					if len(hitsArray) > 0 {
						// Look at the structure of the first hit
						firstHit, ok := hitsArray[0].(map[string]interface{})
						if ok {
							fmt.Printf("First hit keys: %v\n", getMapKeys(firstHit))
						}
					}
				}
			}
		}
	}

	return parseSearchResponse(body)
}

// SearchQuery represents a parsed search query with raw terms and field filters
type SearchQuery struct {
	RawTerms     []string
	FieldFilters map[string][]string // Changed from map[string]string to map[string][]string to support multiple values
}

// ParseSearchQuery parses a search string with support for quoted terms and field:value pairs
func ParseSearchQuery(searchStr string) SearchQuery {
	var query SearchQuery
	query.RawTerms = []string{}
	query.FieldFilters = make(map[string][]string) // Initialize as map to array of strings

	// Regular expression to match quoted strings and field:value pairs
	// This regex handles:
	// 1. Quoted strings (preserving spaces and everything inside)
	// 2. Field:value pairs
	// 3. Regular words
	re := regexp.MustCompile(`"([^"]+)"|(\S+\.\S+):(\S+)|(\S+)`)
	matches := re.FindAllStringSubmatch(searchStr, -1)

	for _, match := range matches {
		if match[1] != "" {
			// This is a quoted string, add it to raw terms
			query.RawTerms = append(query.RawTerms, match[1])
		} else if match[2] != "" && match[3] != "" {
			// This is a field:value pair with dot notation
			fieldName := match[2]
			fieldValue := match[3]
			
			// Add to the array of values for this field
			query.FieldFilters[fieldName] = append(query.FieldFilters[fieldName], fieldValue)
		} else if match[4] != "" {
			// This is a regular word, check if it's a simple field:value
			parts := strings.SplitN(match[4], ":", 2)
			if len(parts) == 2 && !strings.Contains(parts[0], ".") {
				// Simple field:value without dot notation
				fieldName := parts[0]
				fieldValue := parts[1]
				
				// Add to the array of values for this field
				query.FieldFilters[fieldName] = append(query.FieldFilters[fieldName], fieldValue)
			} else {
				// Regular search term
				query.RawTerms = append(query.RawTerms, match[4])
			}
		}
	}

	return query
}

// BuildTypesenseQuery builds a Typesense search query from a parsed SearchQuery
func BuildTypesenseQuery(query SearchQuery) (string, map[string]string, error) {
	// Join raw terms for the main query
	mainQuery := strings.Join(query.RawTerms, " ")

	// Parameters for filter_by and other Typesense parameters
	params := make(map[string]string)

	// Group filter expressions by base field name
	fieldGroups := make(map[string][]string)

	for field, values := range query.FieldFilters {
		// Extract the base field name (part before the first dot)
		baseName := field
		if dotIndex := strings.Index(field, "."); dotIndex != -1 {
			baseName = field[:dotIndex]
		}

		for _, value := range values {
			// Create the filter expression
			filterExpr := fmt.Sprintf("%s:%s", field, value)
			
			// Add to the corresponding field group
			fieldGroups[baseName] = append(fieldGroups[baseName], filterExpr)
		}
	}

	// Build the final filter expressions
	var finalFilterExpressions []string

	for _, expressions := range fieldGroups {
		if len(expressions) == 1 {
			// Single expression, add as is
			finalFilterExpressions = append(finalFilterExpressions, expressions[0])
		} else {
			// Multiple expressions for same base field, join with OR
			orExpression := fmt.Sprintf("(%s)", strings.Join(expressions, " || "))
			finalFilterExpressions = append(finalFilterExpressions, orExpression)
		}
	}

	// Combine all filter expressions with AND
	if len(finalFilterExpressions) > 0 {
		params["filter_by"] = strings.Join(finalFilterExpressions, " && ")
	}

	return mainQuery, params, nil
}

func parseSearchResponse(responseBody []byte) ([]nostr.Event, error) {
	var searchResponse SearchResponse
	if err := json.Unmarshal(responseBody, &searchResponse); err != nil {
		return nil, fmt.Errorf("error parsing search response: %v", err)
	}

	// Debug: Print the raw response structure
	fmt.Printf("Search response found %d hits\n", searchResponse.Found)
	
	nostrResults := make([]nostr.Event, 0, len(searchResponse.Hits))

	for i, hit := range searchResponse.Hits {
		// Debug: Print hit structure information
		fmt.Printf("Processing hit %d, keys: %v\n", i, getMapKeys(hit))
		
		// Check if document exists in the hit
		docRaw, exists := hit["document"]
		if !exists {
			fmt.Printf("Warning: hit %d has no 'document' field\n", i)
			continue // Skip this hit
		}
		
		// Extract document directly as a map[string]interface{}
		docMap, ok := docRaw.(map[string]interface{})
		if !ok {
			fmt.Printf("Warning: hit %d document is not a map, type: %T\n", i, docRaw)
			continue // Skip this hit
		}
		
		// Debug: Print document keys
		fmt.Printf("Document keys: %v\n", getMapKeys(docMap))
		
		// Check for EventRaw field directly
		eventRawVal, hasEventRaw := docMap["eventRaw"]
		if !hasEventRaw {
			fmt.Printf("Warning: document has no 'eventRaw' field\n")
			continue // Skip this document
		}
		
		// Try to extract EventRaw as string
		eventRawStr, ok := eventRawVal.(string)
		if !ok {
			fmt.Printf("Warning: eventRaw is not a string, type: %T\n", eventRawVal)
			continue // Skip this document
		}
		
		// Convert the EventRaw string to a Nostr event
		nostrEvent, err := StringifiedJSONToNostrEvent(eventRawStr)
		if err != nil {
			fmt.Printf("Warning: failed to convert EventRaw to Nostr event: %v\n", err)
			continue
		}

		nostrResults = append(nostrResults, nostrEvent)
	}

	// Print the number of results for logging
	fmt.Printf("Successfully processed %d results\n", len(nostrResults))

	return nostrResults, nil
}

// Helper function to get keys from a map for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
