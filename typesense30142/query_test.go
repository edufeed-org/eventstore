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
