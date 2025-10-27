package typesense30142

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

// tagGroup represents a grouped tag with its full key and values
type tagGroup struct {
	fullKey string
	values  []string
}

// groupTagsByBase groups tags by their base key (before first ':')
// Example: "creator:name" and "creator:affiliation:name" both group under "creator"
func groupTagsByBase(tags nostr.Tags) map[string][]tagGroup {
	grouped := make(map[string][]tagGroup)

	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		key := tag[0]
		values := tag[1:]

		// Skip special Nostr-native tags (handled separately)
		if key == "d" || key == "t" || key == "p" || key == "a" {
			continue
		}

		// Get base key (before first ':')
		baseKey := key
		if idx := strings.Index(key, ":"); idx != -1 {
			baseKey = key[:idx]
		}

		grouped[baseKey] = append(grouped[baseKey], tagGroup{
			fullKey: key,
			values:  values,
		})
	}

	return grouped
}

// parseNestedValue reconstructs a nested value from a colon-delimited key
// Example: "creator:affiliation:name" -> sets value in nested structure
func parseNestedValue(result map[string]interface{}, key string, value string) {
	parts := strings.Split(key, ":")

	// Navigate/create nested structure
	current := result
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}
		current = current[part].(map[string]interface{})
	}

	// Set the final value
	lastKey := parts[len(parts)-1]
	current[lastKey] = value
}

// separateArrayInstances detects multiple instances in array fields
// When we see a repeated base property (e.g., second "name" in creator:name),
// it indicates a new array item
func separateArrayInstances(tags []tagGroup) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	instances := []map[string]interface{}{}
	currentInstance := make(map[string]interface{})
	seenKeys := make(map[string]bool)

	for _, tag := range tags {
		// Get the immediate child key after base
		// E.g., "creator:name" -> "name", "creator:affiliation:name" -> "affiliation"
		parts := strings.Split(tag.fullKey, ":")
		if len(parts) < 2 {
			continue
		}

		immediateKey := parts[1]

		// If we've seen this immediate key before and current instance is not empty,
		// we're starting a new array item
		if seenKeys[immediateKey] && len(currentInstance) > 0 {
			instances = append(instances, currentInstance)
			currentInstance = make(map[string]interface{})
			seenKeys = make(map[string]bool)
		}

		seenKeys[immediateKey] = true

		// Remove base key from the path
		subKey := strings.Join(parts[1:], ":")
		parseNestedValue(currentInstance, subKey, tag.values[0])
	}

	// Don't forget the last instance
	if len(currentInstance) > 0 {
		instances = append(instances, currentInstance)
	}

	return instances
}

// handleTTags extracts keywords from Nostr-native 't' tags
func handleTTags(tags nostr.Tags) []string {
	var keywords []string
	for _, tag := range tags {
		if len(tag) >= 2 && tag[0] == "t" {
			keywords = append(keywords, tag[1])
		}
	}
	return keywords
}

// handlePTags extracts creator and contributor pubkeys from 'p' tags
func handlePTags(tags nostr.Tags) (creatorPubkeys []string, contributorPubkeys []string) {
	for _, tag := range tags {
		if len(tag) >= 4 && tag[0] == "p" {
			pubkey := tag[1]
			marker := tag[3]

			switch marker {
			case "creator":
				creatorPubkeys = append(creatorPubkeys, pubkey)
			case "contributor":
				contributorPubkeys = append(contributorPubkeys, pubkey)
			}
		}
	}
	return
}

// handleATags extracts addressable event references with relation markers
func handleATags(tags nostr.Tags) (isBasedOn, isPartOf, hasPart []string) {
	for _, tag := range tags {
		if len(tag) >= 4 && tag[0] == "a" {
			eventAddr := tag[1]
			marker := tag[3]

			switch marker {
			case "isBasedOn":
				isBasedOn = append(isBasedOn, eventAddr)
			case "isPartOf":
				isPartOf = append(isPartOf, eventAddr)
			case "hasPart":
				hasPart = append(hasPart, eventAddr)
			}
		}
	}
	return
}

// Field-specific parsers for complex nested structures

// parseControlledVocabularyArray parses an array of ControlledVocabulary items from grouped tags
func parseControlledVocabularyArray(instances []map[string]interface{}) []ControlledVocabulary {
	var result []ControlledVocabulary
	for _, inst := range instances {
		cv := ControlledVocabulary{}
		if id, ok := inst["id"].(string); ok {
			cv.ID = id
		}
		if prefLabel, ok := inst["prefLabel"].(string); ok {
			cv.PrefLabel = prefLabel
		}
		if inLanguage, ok := inst["inLanguage"].(string); ok {
			cv.InLanguage = inLanguage
		}
		if typ, ok := inst["type"].(string); ok {
			cv.Type = typ
		}
		result = append(result, cv)
	}
	return result
}

// parseCreators parses creator array from grouped tags and p tag pubkeys
func parseCreators(instances []map[string]interface{}, pubkeys []string) []*Creator {
	var creators []*Creator

	// First, add creators with pubkeys from p tags
	for _, pubkey := range pubkeys {
		creators = append(creators, &Creator{
			BaseEntity: BaseEntity{
				ID: pubkey,
			},
		})
	}

	// Then parse creators from flattened tags
	for _, inst := range instances {
		creator := &Creator{}

		if id, ok := inst["id"].(string); ok {
			creator.ID = id
		}
		if name, ok := inst["name"].(string); ok {
			creator.Name = name
		}
		if typ, ok := inst["type"].(string); ok {
			creator.Type = typ
		}
		if honorificPrefix, ok := inst["honorificPrefix"].(string); ok {
			creator.HonoricPrefix = honorificPrefix
		}

		// Parse affiliation if present
		if affData, ok := inst["affiliation"].(map[string]interface{}); ok {
			aff := &Affiliation{}
			if id, ok := affData["id"].(string); ok {
				aff.ID = id
			}
			if name, ok := affData["name"].(string); ok {
				aff.Name = name
			}
			if typ, ok := affData["type"].(string); ok {
				aff.Type = typ
			}
			creator.Affiliation = aff
		}

		creators = append(creators, creator)
	}

	return creators
}

// parseContributors parses contributor array from grouped tags and p tag pubkeys
func parseContributors(instances []map[string]interface{}, pubkeys []string) []*Contributor {
	var contributors []*Contributor

	// First, add contributors with pubkeys from p tags
	for _, pubkey := range pubkeys {
		contributors = append(contributors, &Contributor{
			BaseEntity: BaseEntity{
				ID: pubkey,
			},
		})
	}

	// Then parse contributors from flattened tags
	for _, inst := range instances {
		contributor := &Contributor{}

		if id, ok := inst["id"].(string); ok {
			contributor.ID = id
		}
		if name, ok := inst["name"].(string); ok {
			contributor.Name = name
		}
		if typ, ok := inst["type"].(string); ok {
			contributor.Type = typ
		}
		if honorificPrefix, ok := inst["honorificPrefix"].(string); ok {
			contributor.HonoricPrefix = honorificPrefix
		}

		// Parse affiliation if present
		if affData, ok := inst["affiliation"].(map[string]interface{}); ok {
			aff := &Affiliation{}
			if id, ok := affData["id"].(string); ok {
				aff.ID = id
			}
			if name, ok := affData["name"].(string); ok {
				aff.Name = name
			}
			if typ, ok := affData["type"].(string); ok {
				aff.Type = typ
			}
			contributor.Affiliation = aff
		}

		contributors = append(contributors, contributor)
	}

	return contributors
}

// parseBaseEntityArray parses an array of BaseEntity items
func parseBaseEntityArray(instances []map[string]interface{}) []BaseEntity {
	var result []BaseEntity
	for _, inst := range instances {
		entity := BaseEntity{}
		if id, ok := inst["id"].(string); ok {
			entity.ID = id
		}
		if name, ok := inst["name"].(string); ok {
			entity.Name = name
		}
		if typ, ok := inst["type"].(string); ok {
			entity.Type = typ
		}
		result = append(result, entity)
	}
	return result
}

// parseTrailers parses trailer array from grouped tags
func parseTrailers(instances []map[string]interface{}) []*Trailer {
	var trailers []*Trailer
	for _, inst := range instances {
		trailer := &Trailer{}
		if contentUrl, ok := inst["contentUrl"].(string); ok {
			trailer.ContentUrl = contentUrl
		}
		if typ, ok := inst["type"].(string); ok {
			trailer.Type = typ
		}
		if encodingFormat, ok := inst["encodingFormat"].(string); ok {
			trailer.EncodingFormat = encodingFormat
		}
		if contentSize, ok := inst["contentSize"].(string); ok {
			trailer.ContentSize = contentSize
		}
		if sha256, ok := inst["sha256"].(string); ok {
			trailer.Sha256 = sha256
		}
		if embedUrl, ok := inst["embedUrl"].(string); ok {
			trailer.EmbedUrl = embedUrl
		}
		if bitrate, ok := inst["bitrate"].(string); ok {
			trailer.Bitrate = bitrate
		}
		trailers = append(trailers, trailer)
	}
	return trailers
}

// parseEncodings parses encoding array from grouped tags
func parseEncodings(instances []map[string]interface{}) []*Encoding {
	var encodings []*Encoding
	for _, inst := range instances {
		encoding := &Encoding{}
		if typ, ok := inst["type"].(string); ok {
			encoding.Type = typ
		}
		if contentUrl, ok := inst["contentUrl"].(string); ok {
			encoding.ContentUrl = contentUrl
		}
		if embedUrl, ok := inst["embedUrl"].(string); ok {
			encoding.EmbedUrl = embedUrl
		}
		if encodingFormat, ok := inst["encodingFormat"].(string); ok {
			encoding.EncodingFormat = encodingFormat
		}
		if contentSize, ok := inst["contentSize"].(string); ok {
			encoding.ContentSize = contentSize
		}
		if sha256, ok := inst["sha256"].(string); ok {
			encoding.Sha256 = sha256
		}
		if bitrate, ok := inst["bitrate"].(string); ok {
			encoding.Bitrate = bitrate
		}
		encodings = append(encodings, encoding)
	}
	return encodings
}

// parseCaptions parses caption array from grouped tags
func parseCaptions(instances []map[string]interface{}) []*Caption {
	var captions []*Caption
	for _, inst := range instances {
		caption := &Caption{}
		if id, ok := inst["id"].(string); ok {
			caption.ID = id
		}
		if typ, ok := inst["type"].(string); ok {
			caption.Type = typ
		}
		if encodingFormat, ok := inst["encodingFormat"].(string); ok {
			caption.EncodingFormat = encodingFormat
		}
		if inLanguage, ok := inst["inLanguage"].(string); ok {
			caption.InLanguage = inLanguage
		}
		captions = append(captions, caption)
	}
	return captions
}

// parseMainEntityOfPages parses mainEntityOfPage array from grouped tags
func parseMainEntityOfPages(instances []map[string]interface{}) []*MainEntityOfPage {
	var pages []*MainEntityOfPage
	for _, inst := range instances {
		page := &MainEntityOfPage{}
		if id, ok := inst["id"].(string); ok {
			page.ID = id
		}
		if typ, ok := inst["type"].(string); ok {
			page.Type = typ
		}
		if dateCreated, ok := inst["dateCreated"].(string); ok {
			page.DateCreated = dateCreated
		}
		if dateModified, ok := inst["dateModified"].(string); ok {
			page.DateModified = dateModified
		}

		// Parse provider if present
		if providerData, ok := inst["provider"].(map[string]interface{}); ok {
			provider := &MainEntityProvider{}
			if id, ok := providerData["id"].(string); ok {
				provider.ID = id
			}
			if name, ok := providerData["name"].(string); ok {
				provider.Name = name
			}
			if typ, ok := providerData["type"].(string); ok {
				provider.Type = typ
			}
			page.Provider = provider
		}

		pages = append(pages, page)
	}
	return pages
}

// parseLicense parses a single license from grouped tags
func parseLicense(instances []map[string]interface{}) *License {
	if len(instances) == 0 {
		return nil
	}
	inst := instances[0]
	license := &License{}
	if id, ok := inst["id"].(string); ok {
		license.ID = id
	}
	if name, ok := inst["name"].(string); ok {
		license.Name = name
	}
	return license
}

// parseConditionsOfAccess parses a single conditionsOfAccess from grouped tags
func parseConditionsOfAccess(instances []map[string]interface{}) *ConditionsOfAccess {
	if len(instances) == 0 {
		return nil
	}
	inst := instances[0]
	coa := &ConditionsOfAccess{}
	if id, ok := inst["id"].(string); ok {
		coa.ID = id
	}
	if prefLabel, ok := inst["prefLabel"].(string); ok {
		coa.PrefLabel = prefLabel
	}
	if inLanguage, ok := inst["inLanguage"].(string); ok {
		coa.InLanguage = inLanguage
	}
	if typ, ok := inst["type"].(string); ok {
		coa.Type = typ
	}
	return coa
}

// parseInteractivityType parses a single interactivityType from grouped tags
func parseInteractivityType(instances []map[string]interface{}) *InteractivityType {
	if len(instances) == 0 {
		return nil
	}
	inst := instances[0]
	it := &InteractivityType{}
	if id, ok := inst["id"].(string); ok {
		it.ID = id
	}
	if prefLabel, ok := inst["prefLabel"].(string); ok {
		it.PrefLabel = prefLabel
	}
	if inLanguage, ok := inst["inLanguage"].(string); ok {
		it.InLanguage = inLanguage
	}
	if typ, ok := inst["type"].(string); ok {
		it.Type = typ
	}
	return it
}

// converts a nostr event to stringified JSON
func eventToStringifiedJSON(event *nostr.Event) (string, error) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	jsonString := string(jsonData)
	return jsonString, err
}

// NostrToAMB converts a Nostr event of kind 30142 to AMB metadata using the new flattened tag format
func NostrToAMB(event *nostr.Event) (*AMBMetadata, error) {
	if event == nil {
		return nil, fmt.Errorf("cannot convert nil event")
	}

	eventRaw, err := eventToStringifiedJSON(event)
	if err != nil {
		return nil, fmt.Errorf("error converting event to JSON: %w", err)
	}

	// Initialize AMB metadata
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
	}

	// Phase 1: Extract Nostr-native tags
	// Handle 'd' tag
	for _, tag := range event.Tags {
		if len(tag) >= 2 && tag[0] == "d" {
			amb.D = tag[1]
			amb.ID = event.ID
			break
		}
	}

	// Handle 't' tags (keywords)
	amb.Keywords = handleTTags(event.Tags)

	// Handle 'p' tags (creator/contributor pubkeys)
	creatorPubkeys, contributorPubkeys := handlePTags(event.Tags)

	// Handle 'a' tags (addressable event relations)
	isBasedOnAddrs, isPartOfAddrs, hasPartAddrs := handleATags(event.Tags)

	// Phase 2: Group remaining flattened tags by base key
	grouped := groupTagsByBase(event.Tags)

	// Phase 3: Parse simple fields (non-nested)
	if tags, ok := grouped["name"]; ok && len(tags) > 0 {
		amb.Name = tags[0].values[0]
	}
	if tags, ok := grouped["description"]; ok && len(tags) > 0 {
		amb.Description = tags[0].values[0]
	}
	if tags, ok := grouped["image"]; ok && len(tags) > 0 {
		amb.Image = tags[0].values[0]
	}
	if tags, ok := grouped["duration"]; ok && len(tags) > 0 {
		amb.Duration = tags[0].values[0]
	}
	if tags, ok := grouped["dateCreated"]; ok && len(tags) > 0 {
		amb.DateCreated = tags[0].values[0]
	}
	if tags, ok := grouped["datePublished"]; ok && len(tags) > 0 {
		amb.DatePublished = tags[0].values[0]
	}
	if tags, ok := grouped["dateModified"]; ok && len(tags) > 0 {
		amb.DateModified = tags[0].values[0]
	}
	if tags, ok := grouped["isAccessibleForFree"]; ok && len(tags) > 0 {
		amb.IsAccessibleForFree = (tags[0].values[0] == "true" || tags[0].values[0] == "1")
	}

	// Parse type (can have multiple values)
	if tags, ok := grouped["type"]; ok {
		amb.Type = []string{}
		for _, tag := range tags {
			amb.Type = append(amb.Type, tag.values[0])
		}
	}

	// Parse inLanguage (can have multiple values)
	if tags, ok := grouped["inLanguage"]; ok {
		for _, tag := range tags {
			amb.InLanguage = append(amb.InLanguage, tag.values[0])
		}
	}

	// Phase 4: Parse complex nested fields using helper functions

	// About (array of ControlledVocabulary)
	if tags, ok := grouped["about"]; ok {
		instances := separateArrayInstances(tags)
		cvArray := parseControlledVocabularyArray(instances)
		for _, cv := range cvArray {
			amb.About = append(amb.About, &About{ControlledVocabulary: cv})
		}
	}

	// LearningResourceType (array of ControlledVocabulary)
	if tags, ok := grouped["learningResourceType"]; ok {
		instances := separateArrayInstances(tags)
		cvArray := parseControlledVocabularyArray(instances)
		for _, cv := range cvArray {
			amb.LearningResourceType = append(amb.LearningResourceType, &LearningResourceType{ControlledVocabulary: cv})
		}
	}

	// Audience (array of ControlledVocabulary)
	if tags, ok := grouped["audience"]; ok {
		instances := separateArrayInstances(tags)
		cvArray := parseControlledVocabularyArray(instances)
		for _, cv := range cvArray {
			amb.Audience = append(amb.Audience, &Audience{ControlledVocabulary: cv})
		}
	}

	// Teaches (array of ControlledVocabulary)
	if tags, ok := grouped["teaches"]; ok {
		instances := separateArrayInstances(tags)
		cvArray := parseControlledVocabularyArray(instances)
		for _, cv := range cvArray {
			amb.Teaches = append(amb.Teaches, &Teaches{ControlledVocabulary: cv})
		}
	}

	// Assesses (array of ControlledVocabulary)
	if tags, ok := grouped["assesses"]; ok {
		instances := separateArrayInstances(tags)
		cvArray := parseControlledVocabularyArray(instances)
		for _, cv := range cvArray {
			amb.Assesses = append(amb.Assesses, &Assesses{ControlledVocabulary: cv})
		}
	}

	// CompetencyRequired (array of ControlledVocabulary)
	if tags, ok := grouped["competencyRequired"]; ok {
		instances := separateArrayInstances(tags)
		cvArray := parseControlledVocabularyArray(instances)
		for _, cv := range cvArray {
			amb.CompetencyRequired = append(amb.CompetencyRequired, &CompetencyRequired{ControlledVocabulary: cv})
		}
	}

	// EducationalLevel (array of ControlledVocabulary)
	if tags, ok := grouped["educationalLevel"]; ok {
		instances := separateArrayInstances(tags)
		cvArray := parseControlledVocabularyArray(instances)
		for _, cv := range cvArray {
			amb.EducationalLevel = append(amb.EducationalLevel, &EducationalLevel{ControlledVocabulary: cv})
		}
	}

	// InteractivityType (single ControlledVocabulary)
	if tags, ok := grouped["interactivityType"]; ok {
		instances := separateArrayInstances(tags)
		amb.InteractivityType = parseInteractivityType(instances)
	}

	// ConditionsOfAccess (single ControlledVocabulary)
	if tags, ok := grouped["conditionsOfAccess"]; ok {
		instances := separateArrayInstances(tags)
		amb.ConditionsOfAccess = parseConditionsOfAccess(instances)
	}

	// License (single object)
	if tags, ok := grouped["license"]; ok {
		instances := separateArrayInstances(tags)
		amb.License = parseLicense(instances)
	}

	// Creator (array with potential p tag pubkeys)
	if tags, ok := grouped["creator"]; ok {
		instances := separateArrayInstances(tags)
		amb.Creator = parseCreators(instances, creatorPubkeys)
	} else if len(creatorPubkeys) > 0 {
		// Only p tags, no flattened creator tags
		amb.Creator = parseCreators(nil, creatorPubkeys)
	}

	// Contributor (array with potential p tag pubkeys)
	if tags, ok := grouped["contributor"]; ok {
		instances := separateArrayInstances(tags)
		amb.Contributor = parseContributors(instances, contributorPubkeys)
	} else if len(contributorPubkeys) > 0 {
		// Only p tags, no flattened contributor tags
		amb.Contributor = parseContributors(nil, contributorPubkeys)
	}

	// Publisher (array of BaseEntity)
	if tags, ok := grouped["publisher"]; ok {
		instances := separateArrayInstances(tags)
		entities := parseBaseEntityArray(instances)
		for _, entity := range entities {
			amb.Publisher = append(amb.Publisher, &Publisher{BaseEntity: entity})
		}
	}

	// Funder (array of BaseEntity)
	if tags, ok := grouped["funder"]; ok {
		instances := separateArrayInstances(tags)
		entities := parseBaseEntityArray(instances)
		for _, entity := range entities {
			amb.Funder = append(amb.Funder, &Funder{BaseEntity: entity})
		}
	}

	// IsBasedOn (array, mixing a tags and flattened)
	if tags, ok := grouped["isBasedOn"]; ok {
		instances := separateArrayInstances(tags)
		for _, inst := range instances {
			isBasedOn := &IsBasedOn{}
			if id, ok := inst["id"].(string); ok {
				isBasedOn.ID = id
			}
			if name, ok := inst["name"].(string); ok {
				isBasedOn.Name = name
			}
			if typ, ok := inst["type"].(string); ok {
				isBasedOn.Type = typ
			}
			amb.IsBasedOn = append(amb.IsBasedOn, isBasedOn)
		}
	}
	// Add addressable event references from a tags
	for _, addr := range isBasedOnAddrs {
		amb.IsBasedOn = append(amb.IsBasedOn, &IsBasedOn{
			ID: addr,
		})
	}

	// IsPartOf (array, mixing a tags and flattened)
	if tags, ok := grouped["isPartOf"]; ok {
		instances := separateArrayInstances(tags)
		entities := parseBaseEntityArray(instances)
		for _, entity := range entities {
			amb.IsPartOf = append(amb.IsPartOf, &IsPartOf{BaseEntity: entity})
		}
	}
	// Add addressable event references from a tags
	for _, addr := range isPartOfAddrs {
		amb.IsPartOf = append(amb.IsPartOf, &IsPartOf{
			BaseEntity: BaseEntity{ID: addr},
		})
	}

	// HasPart (array, mixing a tags and flattened)
	if tags, ok := grouped["hasPart"]; ok {
		instances := separateArrayInstances(tags)
		entities := parseBaseEntityArray(instances)
		for _, entity := range entities {
			amb.HasPart = append(amb.HasPart, &HasPart{BaseEntity: entity})
		}
	}
	// Add addressable event references from a tags
	for _, addr := range hasPartAddrs {
		amb.HasPart = append(amb.HasPart, &HasPart{
			BaseEntity: BaseEntity{ID: addr},
		})
	}

	// Trailer (array of MediaObject)
	if tags, ok := grouped["trailer"]; ok {
		instances := separateArrayInstances(tags)
		amb.Trailer = parseTrailers(instances)
	}

	// Encoding (array of MediaObject)
	if tags, ok := grouped["encoding"]; ok {
		instances := separateArrayInstances(tags)
		amb.Encoding = parseEncodings(instances)
	}

	// Caption (array of MediaObject)
	if tags, ok := grouped["caption"]; ok {
		instances := separateArrayInstances(tags)
		amb.Caption = parseCaptions(instances)
	}

	// MainEntityOfPage (array of WebPage)
	if tags, ok := grouped["mainEntityOfPage"]; ok {
		instances := separateArrayInstances(tags)
		amb.MainEntityOfPage = parseMainEntityOfPages(instances)
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
