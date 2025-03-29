package typesense30142

import (
	"context"
	"fmt"

	"github.com/fiatjaf/eventstore"
	"github.com/nbd-wtf/go-nostr"
)

var _ eventstore.Store = (*TSBackend)(nil)

type TSBackend struct {
	ApiKey         string
	Host           string
	CollectionName string
}

func (ts *TSBackend) Init() error {
	err := ts.CheckOrCreateCollection()
	if err != nil {
		return fmt.Errorf("Failed to check/create collection: %v", err)
	}

	return nil
}

func (ts *TSBackend) Close() {}

func (ts *TSBackend) SaveEvent(ctx context.Context, event *nostr.Event) error {return nil}
