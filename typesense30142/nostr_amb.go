package typesense30142

import (
	"encoding/json"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

// converts a nostr event to stringified JSON
func eventToStringifiedJSON(event *nostr.Event) (string, error) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	jsonString := string(jsonData)
	return jsonString, err
}

// NostrToAMB converts a Nostr event of kind 30142 to AMB metadata
func NostrToAMB(event *nostr.Event) (*AMBMetadata, error) {
	if event == nil {
		return nil, fmt.Errorf("cannot convert nil event")
	}

	eventRaw, err := eventToStringifiedJSON(event)
	if err != nil {
		return nil, fmt.Errorf("error converting event to JSON: %w", err)
	}
	amb := &AMBMetadata{
		Type: []string{"LearningResource"},
		NostrMetadata: NostrMetadata{
			EventID:        event.ID,
			EventPubKey:    event.PubKey,
			EventContent:   event.Content,
			EventCreatedAt: event.CreatedAt,
			EventKind:      event.Kind,
			EventSig:       event.Sig,
			EventRaw:       eventRaw,
		},
		About:                []*About{},
		Keywords:             []string{},
		InLanguage:           []string{},
		Creator:              []*Creator{},
		Contributor:          []*Contributor{},
		Publisher:            []*Publisher{},
		Funder:               []*Funder{},
		LearningResourceType: []*LearningResourceType{},
		Audience:             []*Audience{},
		Teaches:              []*Teaches{},
		Assesses:             []*Assesses{},
		CompetencyRequired:   []*CompetencyRequired{},
		EducationalLevel:     []*EducationalLevel{},
		IsBasedOn:            []*IsBasedOn{},
		IsPartOf:             []*IsPartOf{},
		HasPart:              []*HasPart{},
		Trailer:              []*Trailer{},
	}

	for _, tag := range event.Tags {
		if len(tag) < 2 {
			continue
		}

		switch tag[0] {
		case "d":
			if len(tag) >= 2 {
				amb.D = tag[1]
				amb.ID = event.ID
			}
		case "type":
			if len(tag) >= 2 {
				amb.Type = tag[1:]
			}

		case "name":
			if len(tag) >= 2 {
				amb.Name = tag[1]
			}
		case "description":
			if len(tag) >= 2 {
				amb.Description = tag[1]
			}
		case "creator":
			if len(tag) >= 3 {
				creator := &Creator{
					BaseEntity: BaseEntity{
						ID:   tag[1],
						Name: tag[2],
					},
					Affiliation: &Affiliation{},
				}
				if len(tag) >= 4 {
					creator.Type = tag[3]
				}
				if len(tag) >= 5 {
					creator.Affiliation.Name = tag[4]
				}
				if len(tag) >= 6 {
					creator.Affiliation.Type = tag[5]
				}
				if len(tag) >= 7 {
					creator.Affiliation.ID = tag[6]
				}
				amb.Creator = append(amb.Creator, creator)
			}
		case "image":
			if len(tag) >= 2 {
				amb.Image = tag[1]
			}
		case "about":
			if len(tag) >= 4 {
				subject := &About{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				if len(tag) >= 5 {
					subject.Type = tag[4]
				}
				amb.About = append(amb.About, subject)
			}
		case "learningResourceType":
			if len(tag) >= 4 {
				lrt := &LearningResourceType{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.LearningResourceType = append(amb.LearningResourceType, lrt)
			}
		case "inLanguage":
			if len(tag) >= 2 {
				amb.InLanguage = append(amb.InLanguage, tag[1])
			}
		case "keywords":
			if len(tag) >= 2 {
				amb.Keywords = tag[1:]
			}
		case "license":
			if len(tag) >= 3 {
				amb.License = &License{
					ID:   tag[1],
					Name: tag[2],
				}
			}
		case "datePublished":
			if len(tag) >= 2 {
				amb.DatePublished = tag[1]
			}
		case "dateCreated":
			if len(tag) >= 2 {
				amb.DateCreated = tag[1]
			}
		case "dateModified":
			if len(tag) >= 2 {
				amb.DateModified = tag[1]
			}
		case "publisher":
			if len(tag) >= 3 {
				publisher := &Publisher{
					BaseEntity: BaseEntity{
						ID:   tag[1],
						Name: tag[2],
					},
				}
				if len(tag) >= 4 {
					publisher.Type = tag[3]
				}
				amb.Publisher = append(amb.Publisher, publisher)
			}
		case "contributor":
			if len(tag) >= 4 {
				contributor := &Contributor{
					BaseEntity: BaseEntity{
						ID:   tag[1],
						Name: tag[2],
						Type: tag[3],
					},
					Affiliation: &Affiliation{},
				}
				if len(tag) >= 5 {
					contributor.Affiliation.Name = tag[4]
				}
				if len(tag) >= 6 {
					contributor.Affiliation.Type = tag[5]
				}
				if len(tag) >= 7 {
					contributor.Affiliation.ID = tag[6]
				}
				amb.Contributor = append(amb.Contributor, contributor)
			}
		case "funder":
			if len(tag) >= 3 {
				funder := &Funder{
					BaseEntity: BaseEntity{
						ID:   tag[1],
						Name: tag[2],
					},
				}
				if len(tag) >= 4 {
					funder.Type = tag[3]
				}
				amb.Funder = append(amb.Funder, funder)
			}
		case "isAccessibleForFree":
			if len(tag) >= 2 && (tag[1] == "true" || tag[1] == "1") {
				amb.IsAccessibleForFree = true
			}
		case "audience":
			if len(tag) >= 4 {
				audience := &Audience{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.Audience = append(amb.Audience, audience)
			}
		case "duration":
			if len(tag) >= 2 {
				amb.Duration = tag[1]
			}
		case "conditionsOfAccess":
			if len(tag) == 4 {
				conditionsOfAccess := &ConditionsOfAccess{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.ConditionsOfAccess = conditionsOfAccess
			}
		case "teaches":
			if len(tag) >= 3 {
				teaches := &Teaches{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.Teaches = append(amb.Teaches, teaches)
			}
		case "assesses":
			if len(tag) >= 3 {
				assesses := &Assesses{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.Assesses = append(amb.Assesses, assesses)
			}
		case "competencyRequired":
			if len(tag) >= 3 {
				competencyRequired := &CompetencyRequired{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.CompetencyRequired = append(amb.CompetencyRequired, competencyRequired)
			}
		case "educationalLevel":
			if len(tag) >= 3 {
				educationalLevel := &EducationalLevel{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.EducationalLevel = append(amb.EducationalLevel, educationalLevel)
			}
		case "interactivityType":
			if len(tag) >= 3 {
				interactivityType := &InteractivityType{
					ControlledVocabulary: ControlledVocabulary{
						ID:         tag[1],
						PrefLabel:  tag[2],
						InLanguage: tag[3],
					},
				}
				amb.InteractivityType = interactivityType
			}
		case "isBasedOn":
			if len(tag) >= 3 {
				isBasedOn := &IsBasedOn{
					ID:   tag[1],
					Name: tag[2],
				}
				amb.IsBasedOn = append(amb.IsBasedOn, isBasedOn)
			}
		case "isPartOf":
			if len(tag) >= 4 {
				isPartOf := &IsPartOf{
					BaseEntity: BaseEntity{
						ID:   tag[1],
						Name: tag[2],
						Type: tag[3],
					},
				}
				amb.IsPartOf = append(amb.IsPartOf, isPartOf)
			}
		case "hasPart":
			if len(tag) >= 4 {
				hasPart := &HasPart{
					BaseEntity: BaseEntity{
						ID:   tag[1],
						Name: tag[2],
						Type: tag[3],
					},
				}
				amb.HasPart = append(amb.HasPart, hasPart)
			}
		case "trailer":
			if len(tag) >= 8 {
				trailer := &Trailer{
					ContentUrl:     tag[1],
					Type:           tag[2],
					EncodingFormat: tag[3],
					ContentSize:    tag[4],
					Sha256:         tag[5],
					EmbedUrl:       tag[6],
					Bitrate:        tag[7],
				}
				amb.Trailer = append(amb.Trailer, trailer)
			}
		}
	}

	return amb, nil
}

// converts a stringified JSON event to a nostr.Event
func StringifiedJSONToNostrEvent(jsonString string) (nostr.Event, error) {
	var event nostr.Event
	err := json.Unmarshal([]byte(jsonString), &event)
	if err != nil {
		return nostr.Event{}, err
	}
	return event, nil
}

func Hello(name string) (string, error) {
	resp := "Hello " + name + " nice to meet you"
	return resp, nil
}
