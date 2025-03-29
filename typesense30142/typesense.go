package typesense30142

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nbd-wtf/go-nostr"
)

type CollectionSchema struct {
	Name                string  `json:"name"`
	Fields              []Field `json:"fields"`
	DefaultSortingField string  `json:"default_sorting_field"`
	EnableNestedFields  bool    `json:"enable_nested_fields"`
}

type Field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Facet    bool   `json:"facet,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

type SearchResponse struct {
	Found   int              `json:"found"`
	Hits    []map[string]any `json:"hits"`
	Page    int              `json:"page"`
	Request map[string]any   `json:"request"`
}

// CheckOrCreateCollection checks if a collection exists and creates it if it doesn't
func (ts *TSBackend) CheckOrCreateCollection() error {
	exists, err := ts.collectionExists()
	if err != nil {
		log.Fatalf("Error checking collection: %v", err)
	}

	if !exists {
		log.Printf("Collection %s does not exist. Creating...\n", ts.CollectionName)
		if err := ts.createCollection(ts.CollectionName); err != nil {
			log.Fatalf("Error creating collection: %v", err)
		}
		log.Printf("Collection %s created successfully\n", ts.CollectionName)
	} else {
		log.Printf("Collection %s already exists\n", ts.CollectionName)
	}

	return nil
}

func (ts *TSBackend) collectionExists() (bool, error) {
	url := fmt.Sprintf("%s/collections/%s", ts.Host, ts.CollectionName)

	resp, err := ts.makehttpRequest(url, http.MethodGet, nil)
	if err != nil {
		return false, err
	}
	// 404 means collection doesn't exist
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	// Any status code other than 200 is an error
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return true, nil
}

// create a typesense collection
func (ts *TSBackend) createCollection(name string) error {
	schema := CollectionSchema{
		Name: name,
		Fields: []Field{
			// Base information
			{Name: "id", Type: "string"},
			{Name: "d", Type: "string"},
			{Name: "type", Type: "string"},
			{Name: "name", Type: "string"},
			{Name: "description", Type: "string", Optional: true},
			{Name: "about", Type: "object[]", Optional: true},
			{Name: "keywords", Type: "string[]", Optional: true},
			{Name: "inLanguage", Type: "string[]", Optional: true},
			{Name: "image", Type: "string", Optional: true},
			{Name: "trailer", Type: "object[]", Optional: true},

			// Provenience
			{Name: "creator", Type: "object[]", Optional: true},
			{Name: "contributor", Type: "object[]", Optional: true},
			{Name: "dateCreated", Type: "string", Optional: true},
			{Name: "datePublished", Type: "string", Optional: true},
			{Name: "dateModified", Type: "string", Optional: true},
			{Name: "publisher", Type: "object[]", Optional: true},
			{Name: "funder", Type: "object[]", Optional: true},

			// Costs and Rights
			{Name: "isAccessibleForFree", Type: "bool", Optional: true},
			{Name: "license", Type: "object", Optional: true},
			{Name: "conditionsOfAccess", Type: "object", Optional: true},

			// Educational Metadata
			{Name: "learningResourceType", Type: "object[]", Optional: true},
			{Name: "audience", Type: "object[]", Optional: true},
			{Name: "teaches", Type: "object[]", Optional: true},
			{Name: "assesses", Type: "object[]", Optional: true},
			{Name: "competencyRequired", Type: "object[]", Optional: true},
			{Name: "educationalLevel", Type: "object[]", Optional: true},
			{Name: "interactivityType", Type: "object", Optional: true},

			// Relation
			{Name: "isBasedOn", Type: "object[]", Optional: true},
			{Name: "isPartOf", Type: "object[]", Optional: true},
			{Name: "hasPart", Type: "object[]", Optional: true},

			// Technical
			{Name: "duration", Type: "string", Optional: true},

			// Nostr Event
			{Name: "eventID", Type: "string"},
			{Name: "eventKind", Type: "int32"},
			{Name: "eventPubKey", Type: "string"},
			{Name: "eventSignature", Type: "string"},
			{Name: "eventCreatedAt", Type: "int64"},
			{Name: "eventContent", Type: "string"},
			{Name: "eventRaw", Type: "string"},
		},
		DefaultSortingField: "eventCreatedAt",
		EnableNestedFields:  true,
	}

	url := fmt.Sprintf("%s/collections", ts.Host)

	jsonData, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	resp, err := ts.makehttpRequest(url, http.MethodPost, jsonData)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create collection, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Makes an http request to typesense
func (ts *TSBackend) makehttpRequest(url string, method string, reqBody []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-TYPESENSE-API-KEY", ts.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// TODO Count events
func CountEvents(filter nostr.Filter) (int64, error) {
	fmt.Println("filter", filter)
	// search by author
	// search by d-tag
	return 0, nil
}


