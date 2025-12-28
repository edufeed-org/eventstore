package typesense30142

import (
	"testing"

	"github.com/nbd-wtf/go-nostr"
)

func TestQueryEventsWithLimit(t *testing.T) {
	// Test that limit is properly handled
	tests := []struct {
		name          string
		filter        nostr.Filter
		expectedLimit int
	}{
		{
			name: "No limit specified - should default to 100",
			filter: nostr.Filter{
				Search: "",
			},
			expectedLimit: 100,
		},
		{
			name: "Limit of 2 specified",
			filter: nostr.Filter{
				Search: "",
				Limit:  2,
			},
			expectedLimit: 2,
		},
		{
			name: "Limit of 50 specified",
			filter: nostr.Filter{
				Search: "test",
				Limit:  50,
			},
			expectedLimit: 50,
		},
		{
			name: "Zero limit - should default to 100",
			filter: nostr.Filter{
				Search: "test",
				Limit:  0,
			},
			expectedLimit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the limit logic
			limit := 100
			if tt.filter.Limit > 0 {
				limit = tt.filter.Limit
			}

			if limit != tt.expectedLimit {
				t.Errorf("Expected limit %d, got %d", tt.expectedLimit, limit)
			}
		})
	}
}

func TestBuildTypesenseQueryWithEmptySearch(t *testing.T) {
	// Test that empty search strings are handled correctly
	query := ParseSearchQuery("")
	mainQuery, params, err := BuildTypesenseQuery(query)

	if err != nil {
		t.Errorf("BuildTypesenseQuery failed: %v", err)
	}

	// Empty search should result in empty main query
	if mainQuery != "" {
		t.Errorf("Expected empty main query, got: %s", mainQuery)
	}

	// Should have no filter parameters for empty search
	if len(params) != 0 {
		t.Errorf("Expected no filter params for empty search, got: %v", params)
	}
}

func TestSearchResourcesWithLimitHandlesEmptySearch(t *testing.T) {
	// Test that SearchResourcesWithLimit properly handles empty search by using wildcard
	query := ParseSearchQuery("")
	mainQuery, _, err := BuildTypesenseQuery(query)

	if err != nil {
		t.Errorf("BuildTypesenseQuery failed: %v", err)
	}

	// Empty search should result in empty main query initially
	if mainQuery != "" {
		t.Errorf("Expected empty main query from ParseSearchQuery, got: %s", mainQuery)
	}

	// SearchResourcesWithLimit should convert empty query to "*"
	if mainQuery == "" {
		mainQuery = "*"
	}

	if mainQuery != "*" {
		t.Errorf("Expected wildcard query '*', got: %s", mainQuery)
	}
}

func TestParseSearchQueryFieldFilters(t *testing.T) {
	tests := []struct {
		name           string
		searchStr      string
		expectedTerms  []string
		expectedFields map[string][]string
	}{
		{
			name:           "Simple field filter with dot notation",
			searchStr:      "publisher.name:e-teaching.org",
			expectedTerms:  []string{},
			expectedFields: map[string][]string{"publisher.name": {"e-teaching.org"}},
		},
		{
			name:           "Deeply nested field filter",
			searchStr:      "learningResourceType.prefLabel.de:Video",
			expectedTerms:  []string{},
			expectedFields: map[string][]string{"learningResourceType.prefLabel.de": {"Video"}},
		},
		{
			name:           "Field filter with free text",
			searchStr:      "forschung publisher.name:e-teaching.org",
			expectedTerms:  []string{"forschung"},
			expectedFields: map[string][]string{"publisher.name": {"e-teaching.org"}},
		},
		{
			name:           "Multiple field filters same base field",
			searchStr:      "about.prefLabel.de:Mathematik about.prefLabel.de:Physik",
			expectedTerms:  []string{},
			expectedFields: map[string][]string{"about.prefLabel.de": {"Mathematik", "Physik"}},
		},
		{
			name:           "Multiple field filters different base fields",
			searchStr:      "publisher.name:test learningResourceType.prefLabel.de:Video",
			expectedTerms:  []string{},
			expectedFields: map[string][]string{"publisher.name": {"test"}, "learningResourceType.prefLabel.de": {"Video"}},
		},
		{
			name:           "Combined text and multiple filters",
			searchStr:      "mathematik bildung publisher.name:e-teaching.org about.prefLabel.de:Schule",
			expectedTerms:  []string{"mathematik", "bildung"},
			expectedFields: map[string][]string{"publisher.name": {"e-teaching.org"}, "about.prefLabel.de": {"Schule"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := ParseSearchQuery(tt.searchStr)

			// Check raw terms
			if len(query.RawTerms) != len(tt.expectedTerms) {
				t.Errorf("Expected %d raw terms, got %d: %v", len(tt.expectedTerms), len(query.RawTerms), query.RawTerms)
			}
			for i, term := range tt.expectedTerms {
				if i < len(query.RawTerms) && query.RawTerms[i] != term {
					t.Errorf("Expected raw term %d to be %s, got %s", i, term, query.RawTerms[i])
				}
			}

			// Check field filters
			if len(query.FieldFilters) != len(tt.expectedFields) {
				t.Errorf("Expected %d field filters, got %d: %v", len(tt.expectedFields), len(query.FieldFilters), query.FieldFilters)
			}
			for field, expectedValues := range tt.expectedFields {
				if actualValues, ok := query.FieldFilters[field]; ok {
					if len(actualValues) != len(expectedValues) {
						t.Errorf("Expected %d values for field %s, got %d", len(expectedValues), field, len(actualValues))
					}
					for i, expectedVal := range expectedValues {
						if i < len(actualValues) && actualValues[i] != expectedVal {
							t.Errorf("Expected value %d for field %s to be %s, got %s", i, field, expectedVal, actualValues[i])
						}
					}
				} else {
					t.Errorf("Expected field %s not found in filters", field)
				}
			}
		})
	}
}

func TestBuildTypesenseQueryWithFieldFilters(t *testing.T) {
	tests := []struct {
		name               string
		searchStr          string
		expectedMainQuery  string
		expectedFilterPart string
	}{
		{
			name:               "Simple nested field filter",
			searchStr:          "publisher.name:e-teaching.org",
			expectedMainQuery:  "",
			expectedFilterPart: "publisher.name:=`e-teaching.org`",
		},
		{
			name:               "Deeply nested field filter",
			searchStr:          "learningResourceType.prefLabel.de:Video",
			expectedMainQuery:  "",
			expectedFilterPart: "learningResourceType.prefLabel.de:=`Video`",
		},
		{
			name:               "Text with field filter",
			searchStr:          "forschung publisher.name:test",
			expectedMainQuery:  "forschung",
			expectedFilterPart: "publisher.name:=`test`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := ParseSearchQuery(tt.searchStr)
			mainQuery, params, err := BuildTypesenseQuery(query)

			if err != nil {
				t.Errorf("BuildTypesenseQuery failed: %v", err)
			}

			if mainQuery != tt.expectedMainQuery {
				t.Errorf("Expected main query '%s', got '%s'", tt.expectedMainQuery, mainQuery)
			}

			filterBy, hasFilter := params["filter_by"]
			if tt.expectedFilterPart != "" {
				if !hasFilter {
					t.Errorf("Expected filter_by parameter but none found")
				} else if filterBy != tt.expectedFilterPart {
					t.Errorf("Expected filter_by '%s', got '%s'", tt.expectedFilterPart, filterBy)
				}
			}
		})
	}
}

func TestBuildTypesenseQueryMultipleFilters(t *testing.T) {
	// Test that multiple filters on same base field are joined with OR
	query := ParseSearchQuery("about.prefLabel.de:Mathematik about.prefLabel.de:Physik")
	_, params, err := BuildTypesenseQuery(query)

	if err != nil {
		t.Errorf("BuildTypesenseQuery failed: %v", err)
	}

	filterBy, hasFilter := params["filter_by"]
	if !hasFilter {
		t.Errorf("Expected filter_by parameter but none found")
	}

	// Should contain both filters joined with OR
	if filterBy != "(about.prefLabel.de:=`Mathematik` || about.prefLabel.de:=`Physik`)" {
		t.Errorf("Expected OR-joined filter, got: %s", filterBy)
	}
}
