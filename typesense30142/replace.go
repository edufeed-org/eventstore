package typesense30142

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nbd-wtf/go-nostr"
)

// IndexNostrEvent converts a Nostr event to AMB metadata and indexes it in Typesense
func (ts *TSBackend) ReplaceEvent(ctx context.Context, event *nostr.Event) error {
	ambData, err := NostrToAMB(event)
	if err != nil {
		return fmt.Errorf("error converting Nostr event to AMB metadata: %v", err)
	}

	// check if event is already there, if so replace it, else index it
	alreadyIndexed, err := ts.eventAlreadyIndexed(ambData)
	return ts.indexDocument(ctx, ambData, alreadyIndexed)
}

func (ts *TSBackend) eventAlreadyIndexed(doc *AMBMetadata) (*nostr.Event, error) {
	url := fmt.Sprintf(
		"%s/collections/%s/documents/search?filter_by=d:=%s&&eventPubKey:=%s&q=&query_by=d,eventPubKey",
		ts.Host, ts.CollectionName, doc.D, doc.EventPubKey)

	resp, err := ts.makehttpRequest(url, http.MethodGet, nil)
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Search for event failed, status: %d, body: %s", resp.StatusCode, string(body))
	}

	events, err := parseSearchResponse(body)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing search response: %v", err)
	}

	// Check if we found any events
	if len(events) == 0 {
		return nil, nil
	}
	return &events[0], nil
}

// Index a document in Typesense
func (ts *TSBackend) indexDocument(ctx context.Context, doc *AMBMetadata, alreadyIndexedEvent *nostr.Event) error {
	if alreadyIndexedEvent != nil {
		fmt.Println("deleting old event for new one")
		ts.DeleteEvent(ctx, alreadyIndexedEvent)
	}

	url := fmt.Sprintf("%s/collections/%s/documents", ts.Host, ts.CollectionName)

	jsonData, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	resp, err := ts.makehttpRequest(url, http.MethodPost, jsonData)
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to index document, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}


