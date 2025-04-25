package typesense30142

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nbd-wtf/go-nostr"
)

// Delete a nostr event from the index
func (ts *TSBackend) DeleteEvent(ctx context.Context, event *nostr.Event) error {
	fmt.Println("deleting event")
	d := event.Tags.GetD()

	url := fmt.Sprintf(
		"%s/collections/%s/documents?filter_by=d:=%s&&eventPubKey:=%s",
		ts.Host, ts.CollectionName, d, event.PubKey)

	resp, body, err := ts.makehttpRequest(url, http.MethodDelete, nil)
	if err != nil {
		return err
	}

	// Any status code other than 200 is an error
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

