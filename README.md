# Eventstore 

Eventstore to be used with [khatru](https://github.com/fiatjaf/khatru) relays.

Current available stores:
  - typesense30142: A relay specialized for [AMB](https://dini-ag-kim.github.io/amb/latest/#context) Metadata events. See https://github.com/edufeed-org/nips/blob/edufeed-amb/edufeed.md for the current AMB NIP

## Search Query Syntax

The eventstore supports [NIP-50](https://github.com/nostr-protocol/nips/blob/master/50.md) full-text search with advanced filtering capabilities.

### Basic Full-Text Search

Search across all text fields (name, description, keywords, etc.):

```bash
nak req --search "mathematik" -k 30142 ws://localhost:3334
```

### Field-Specific Filtering

Filter by specific nested fields using dot notation `field.path:value`:

```bash
# Filter by publisher name
nak req --search "publisher.name:e-teaching.org" -k 30142 ws://localhost:3334

# Filter by learning resource type (with language-specific label)
nak req --search "learningResourceType.prefLabel.de:Video" -k 30142 ws://localhost:3334

# Filter by subject/topic
nak req --search "about.prefLabel.de:Mathematik" -k 30142 ws://localhost:3334
```

### Combined Queries

Combine free-text search with field filters:

```bash
# Search for "forschung" in text fields AND filter by publisher
nak req --search "forschung publisher.name:e-teaching.org" -k 30142 ws://localhost:3334

# Multiple filters with text search
nak req --search "bildung learningResourceType.prefLabel.de:Video publisher.name:test" -k 30142 ws://localhost:3334
```

### Multiple Values for Same Field

Multiple values for the same base field are combined with OR logic:

```bash
# Find resources about Mathematics OR Physics
nak req --search "about.prefLabel.de:Mathematik about.prefLabel.de:Physik" -k 30142 ws://localhost:3334
```

### Supported Filter Fields

| Field Path | Description | Example |
|------------|-------------|---------|
| `publisher.name` | Publisher organization name | `publisher.name:e-teaching.org` |
| `creator.name` | Content creator name | `creator.name:John` |
| `about.prefLabel.de` | Subject/topic (German) | `about.prefLabel.de:Mathematik` |
| `about.prefLabel.en` | Subject/topic (English) | `about.prefLabel.en:Mathematics` |
| `learningResourceType.prefLabel.de` | Resource type (German) | `learningResourceType.prefLabel.de:Video` |
| `audience.prefLabel.de` | Target audience (German) | `audience.prefLabel.de:Lehrkräfte` |
| `educationalLevel.prefLabel.de` | Educational level (German) | `educationalLevel.prefLabel.de:Sekundarstufe` |

Any nested field path supported by the Typesense schema can be used for filtering.

## Usage Example

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr"
    "github.com/edufeed-org/eventstore/typesense30142"
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
