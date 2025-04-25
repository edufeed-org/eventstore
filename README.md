# Eventstore 

Eventstore to be used with [khatru](https://github.com/fiatjaf/khatru) relays.

Current available stores:
  - typesense30142: A relay specialiced for [AMB](https://dini-ag-kim.github.io/amb/latest/#context) Metadata events. See https://github.com/edufeed-org/nips/blob/edufeed-amb/edufeed.md for the current AMB NIP

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
    "github.com/sroertgen/eventstore/typesense30142"
)

func main() {
	relay := khatru.NewRelay()
	relay.Info.Name = "A edufeed relay for AMB metadata"
	relay.Info.PubKey = "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	relay.Info.Description = "this is the typesense custom relay"
	relay.Info.Icon = "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fliquipedia.net%2Fcommons%2Fimages%2F3%2F35%2FSCProbe.jpg&f=1&nofb=1&ipt=0cbbfef25bce41da63d910e86c3c343e6c3b9d63194ca9755351bb7c2efa3359&ipo=images"
    db := typesense30142.TSBackend{
        ApiKey: "xyz",
        Host: "http://localhost:8108",
        CollectionName: "amb",
    }
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay.OnConnect = append(relay.OnConnect, func(ctx context.Context) {
		khatru.RequestAuth(ctx)
	})

	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)
	relay.ReplaceEvent = append(relay.ReplaceEvent, db.ReplaceEvent)
	relay.Negentropy = true

	relay.RejectEvent = append(relay.RejectEvent,
		func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
			if event.Kind != 30142 {
				return true, "we don't allow these kinds here. It's a 30142 only place."
			}
			return false, ""
		},
	)

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}

```
