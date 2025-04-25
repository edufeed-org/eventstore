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

	nostrs, err := ts.SearchResources(filter.Search)
	if err != nil {
		log.Printf("Search failed: %v", err)
		return ch, err
	}

	go func() {
		for _, evt := range nostrs {
			ch <- &evt
		}
		close(ch)
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
	queryBy := "name,description"

	// Start building the search URL
	searchURL := fmt.Sprintf("%s/collections/%s/documents/search?q=%s&query_by=%s",
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

	return parseSearchResponse(body)
}

// SearchQuery represents a parsed search query with raw terms and field filters
type SearchQuery struct {
	RawTerms     []string
	FieldFilters map[string]string
}

// ParseSearchQuery parses a search string with support for quoted terms and field:value pairs
func ParseSearchQuery(searchStr string) SearchQuery {
	var query SearchQuery
	query.RawTerms = []string{}
	query.FieldFilters = make(map[string]string)

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
			// This is a field:value pair
			fieldName := match[2]
			fieldValue := match[3]
			query.FieldFilters[fieldName] = fieldValue
		} else if match[4] != "" {
			// This is a regular word, check if it's a simple field:value
			parts := strings.SplitN(match[4], ":", 2)
			if len(parts) == 2 && !strings.Contains(parts[0], ".") {
				// Simple field:value without dot notation
				query.FieldFilters[parts[0]] = parts[1]
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

	// Build filter expressions for field filters
	var filterExpressions []string

	for field, value := range query.FieldFilters {
		// Handle special fields with dot notation
		if strings.Contains(field, ".") {
			parts := strings.SplitN(field, ".", 2)
			fieldName := parts[0]
			subField := parts[1]

			filterExpressions = append(filterExpressions, fmt.Sprintf("%s.%s:%s", fieldName, subField, value))
		} else {
			filterExpressions = append(filterExpressions, fmt.Sprintf("%s:%s", field, value))
		}
	}

	// Combine all filter expressions
	if len(filterExpressions) > 0 {
		params["filter_by"] = strings.Join(filterExpressions, " && ")
	}

	return mainQuery, params, nil
}

func parseSearchResponse(responseBody []byte) ([]nostr.Event, error) {
	var searchResponse SearchResponse
	if err := json.Unmarshal(responseBody, &searchResponse); err != nil {
		return nil, fmt.Errorf("error parsing search response: %v", err)
	}

	nostrResults := make([]nostr.Event, 0, len(searchResponse.Hits))

	for _, hit := range searchResponse.Hits {
		// Extract the document from the hit
		docMap, ok := hit["document"]
		if !ok {
			return nil, fmt.Errorf("invalid document format in search results")
		}

		// Convert the map to AMB metadata
		docJSON, err := json.Marshal(docMap)
		if err != nil {
			return nil, fmt.Errorf("error marshaling document: %v", err)
		}

		var ambData AMBMetadata
		if err := json.Unmarshal(docJSON, &ambData); err != nil {
			return nil, fmt.Errorf("error unmarshaling to AMBMetadata: %v", err)
		}

		// Convert the AMB metadata to a Nostr event
		nostrEvent, err := StringifiedJSONToNostrEvent(ambData.EventRaw)
		if err != nil {
			fmt.Printf("Warning: failed to convert AMB to Nostr: %v\n", err)
			continue
		}

		nostrResults = append(nostrResults, nostrEvent)
	}

	// Print the number of results for logging
	fmt.Printf("Found %d results\n",
		len(nostrResults))

	return nostrResults, nil
}
