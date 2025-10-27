package typesense30142

import (
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestNostrToAmbEvent(t *testing.T) {
	assert := assert.New(t)
	sk := nostr.GeneratePrivateKey()
	event := &nostr.Event{
		Content:   "",
		CreatedAt: nostr.Now(),
		Kind:      30142,
		Tags: nostr.Tags{
			{"d", "test-resource-id"},
			{"name", "Test Resource"},
			{"description", "This is a test resource"},
		},
	}

	event.Sign(sk)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(event.ID, amb.EventID)
	assert.Equal(event.PubKey, amb.EventPubKey)
	assert.Equal(event.Content, amb.EventContent)
	assert.Equal(event.Kind, amb.EventKind)
	assert.Equal(event.Sig, amb.EventSig)
	assert.Equal(event.CreatedAt, amb.EventCreatedAt)

}

// Helper function to create a test event with a specific tag
func createTestEvent(tags nostr.Tags) *nostr.Event {
	sk := nostr.GeneratePrivateKey()
	event := &nostr.Event{
		Content:   "",
		CreatedAt: nostr.Now(),
		Kind:      30142,
		Tags:      tags,
	}
	event.Sign(sk)
	return event
}

func TestNostrToAMB_BasicFields(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"name", "Test Resource"},
		{"description", "This is a test resource"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(event.ID, amb.ID)
	assert.Equal("test-resource-id", amb.D)
	assert.Equal("Test Resource", amb.Name)
	assert.Equal("This is a test resource", amb.Description)

	assert.Equal(event.ID, amb.EventID)
	assert.Equal(event.PubKey, amb.EventPubKey)
	assert.Equal(event.Content, amb.EventContent)
	assert.Equal(event.CreatedAt, amb.EventCreatedAt)
	assert.Equal(event.Kind, amb.EventKind)
	assert.Equal(event.Sig, amb.EventSig)
}

func TestNostrToAMB_About(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"about:id", "http://w3id.org/kim/schulfaecher/s1009"},
		{"about:prefLabel", "Französisch"},
		{"about:inLanguage", "de"},
		{"about:type", "Concept"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.NotEmpty(amb.About)
	assert.Equal(1, len(amb.About))
	assert.Equal("http://w3id.org/kim/schulfaecher/s1009", amb.About[0].ID)
	assert.Equal("Französisch", amb.About[0].PrefLabel)
	assert.Equal("de", amb.About[0].InLanguage)
	assert.Equal("Concept", amb.About[0].Type)
}

func TestNostrToAMB_Keywords(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"t", "Französisch"},
		{"t", "Niveau A2"},
		{"t", "Sprache"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(3, len(amb.Keywords))
	assert.Contains(amb.Keywords, "Französisch")
	assert.Contains(amb.Keywords, "Niveau A2")
	assert.Contains(amb.Keywords, "Sprache")
}

func TestNostrToAMB_InLanguage(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"inLanguage", "fr"},
		{"inLanguage", "de"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.InLanguage))
	assert.Contains(amb.InLanguage, "fr")
	assert.Contains(amb.InLanguage, "de")
}

func TestNostrToAMB_Image(t *testing.T) {
	assert := assert.New(t)

	imageURL := "https://www.tutory.de/worksheet/fbbadf1a-145a-463d-9a43-1ae9965c86b9.jpg?width=1000"
	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"image", imageURL},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(imageURL, amb.Image)
}

func TestNostrToAMB_Creator(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"creator:id", "http://author1.org"},
		{"creator:name", "Autorin 1"},
		{"creator:type", "Person"},
		{"creator:id", "http://author2.org"},
		{"creator:name", "Autorin 2"},
		{"creator:type", "Person"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.Creator))

	assert.Equal("http://author1.org", amb.Creator[0].ID)
	assert.Equal("Autorin 1", amb.Creator[0].Name)
	assert.Equal("Person", amb.Creator[0].Type)

	assert.Equal("http://author2.org", amb.Creator[1].ID)
	assert.Equal("Autorin 2", amb.Creator[1].Name)
	assert.Equal("Person", amb.Creator[1].Type)
}

func TestNostrToAMB_Contributor(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"contributor:id", "http://author1.org"},
		{"contributor:name", "Autorin 1"},
		{"contributor:type", "Person"},
		{"contributor:id", "http://author2.org"},
		{"contributor:name", "Autorin 2"},
		{"contributor:type", "Person"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.Contributor))

	assert.Equal("http://author1.org", amb.Contributor[0].ID)
	assert.Equal("Autorin 1", amb.Contributor[0].Name)
	assert.Equal("Person", amb.Contributor[0].Type)

	assert.Equal("http://author2.org", amb.Contributor[1].ID)
	assert.Equal("Autorin 2", amb.Contributor[1].Name)
	assert.Equal("Person", amb.Contributor[1].Type)
}

func TestNostrToAMB_Dates(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"dateCreated", "2019-07-02"},
		{"datePublished", "2019-07-03"},
		{"dateModified", "2019-07-04"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal("2019-07-02", amb.DateCreated)
	assert.Equal("2019-07-03", amb.DatePublished)
	assert.Equal("2019-07-04", amb.DateModified)
}

func TestNostrToAMB_Publisher(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"publisher:id", "http://publisher1.org"},
		{"publisher:name", "Publisher 1"},
		{"publisher:type", "Person"},
		{"publisher:id", "http://publisher2.org"},
		{"publisher:name", "Publisher 2"},
		{"publisher:type", "Organization"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.Publisher))

	assert.Equal("http://publisher1.org", amb.Publisher[0].ID)
	assert.Equal("Publisher 1", amb.Publisher[0].Name)
	assert.Equal("Person", amb.Publisher[0].Type)

	assert.Equal("http://publisher2.org", amb.Publisher[1].ID)
	assert.Equal("Publisher 2", amb.Publisher[1].Name)
	assert.Equal("Organization", amb.Publisher[1].Type)
}

func TestNostrToAMB_Funder(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"funder:id", "http://funder1.org"},
		{"funder:name", "Funder 1"},
		{"funder:type", "Person"},
		{"funder:id", "http://funder2.org"},
		{"funder:name", "Funder 2"},
		{"funder:type", "Organization"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.Funder))

	assert.Equal("http://funder1.org", amb.Funder[0].ID)
	assert.Equal("Funder 1", amb.Funder[0].Name)
	assert.Equal("Person", amb.Funder[0].Type)

	assert.Equal("http://funder2.org", amb.Funder[1].ID)
	assert.Equal("Funder 2", amb.Funder[1].Name)
	assert.Equal("Organization", amb.Funder[1].Type)
}

func TestNostrToAMB_IsAccessibleForFree(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"isAccessibleForFree", "true"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.True(amb.IsAccessibleForFree)

	// Test with "false" value
	tags = nostr.Tags{
		{"d", "test-resource-id"},
		{"isAccessibleForFree", "false"},
	}
	event = createTestEvent(tags)

	amb, err = NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.False(amb.IsAccessibleForFree)
}

func TestNostrToAMB_License(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"license:id", "https://creativecommons.org/publicdomain/zero/1.0/"},
		{"license:name", "CC-0"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.NotNil(amb.License)
	assert.Equal("https://creativecommons.org/publicdomain/zero/1.0/", amb.License.ID)
	assert.Equal("CC-0", amb.License.Name)
}

func TestNostrToAMB_ConditionsOfAccess(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"conditionsOfAccess:id", "http://w3id.org/kim/conditionsOfAccess/no_login"},
		{"conditionsOfAccess:prefLabel", "Kein Login"},
		{"conditionsOfAccess:inLanguage", "de"},
		{"conditionsOfAccess:type", "Concept"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.NotNil(amb.ConditionsOfAccess)
	assert.Equal("http://w3id.org/kim/conditionsOfAccess/no_login", amb.ConditionsOfAccess.ID)
	assert.Equal("Kein Login", amb.ConditionsOfAccess.PrefLabel)
	assert.Equal("de", amb.ConditionsOfAccess.InLanguage)
}

func TestNostrToAMB_LearningResourceType(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"learningResourceType:id", "http://w3id.org/openeduhub/vocabs/new_lrt/video"},
		{"learningResourceType:prefLabel", "Video"},
		{"learningResourceType:inLanguage", "de"},
		{"learningResourceType:id", "http://w3id.org/openeduhub/vocabs/new_lrt/tutorial"},
		{"learningResourceType:prefLabel", "Tutorial"},
		{"learningResourceType:inLanguage", "en"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.LearningResourceType))

	assert.Equal("http://w3id.org/openeduhub/vocabs/new_lrt/video", amb.LearningResourceType[0].ID)
	assert.Equal("Video", amb.LearningResourceType[0].PrefLabel)
	assert.Equal("de", amb.LearningResourceType[0].InLanguage)

	assert.Equal("http://w3id.org/openeduhub/vocabs/new_lrt/tutorial", amb.LearningResourceType[1].ID)
	assert.Equal("Tutorial", amb.LearningResourceType[1].PrefLabel)
	assert.Equal("en", amb.LearningResourceType[1].InLanguage)
}

func TestNostrToAMB_Audience(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"audience:id", "http://purl.org/dcx/lrmi-vocabs/educationalAudienceRole/student"},
		{"audience:prefLabel", "Schüler:in"},
		{"audience:inLanguage", "de"},
		{"audience:id", "http://purl.org/dcx/lrmi-vocabs/educationalAudienceRole/teacher"},
		{"audience:prefLabel", "Lehrer:in"},
		{"audience:inLanguage", "de"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.Audience))

	assert.Equal("http://purl.org/dcx/lrmi-vocabs/educationalAudienceRole/student", amb.Audience[0].ID)
	assert.Equal("Schüler:in", amb.Audience[0].PrefLabel)
	assert.Equal("de", amb.Audience[0].InLanguage)

	assert.Equal("http://purl.org/dcx/lrmi-vocabs/educationalAudienceRole/teacher", amb.Audience[1].ID)
	assert.Equal("Lehrer:in", amb.Audience[1].PrefLabel)
	assert.Equal("de", amb.Audience[1].InLanguage)
}

func TestNostrToAMB_Teaches(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"teaches:id", "http://awesome-skills.org/1"},
		{"teaches:prefLabel", "Zuhören"},
		{"teaches:inLanguage", "de"},
		{"teaches:id", "http://awesome-skills.org/2"},
		{"teaches:prefLabel", "Sprechen"},
		{"teaches:inLanguage", "de"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.Teaches))

	assert.Equal("http://awesome-skills.org/1", amb.Teaches[0].ID)
	assert.Equal("Zuhören", amb.Teaches[0].PrefLabel)
	assert.Equal("de", amb.Teaches[0].InLanguage)

	assert.Equal("http://awesome-skills.org/2", amb.Teaches[1].ID)
	assert.Equal("Sprechen", amb.Teaches[1].PrefLabel)
	assert.Equal("de", amb.Teaches[1].InLanguage)
}

func TestNostrToAMB_Assesses(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"assesses:id", "http://awesome-skills.org/1"},
		{"assesses:prefLabel", "Hörverständnis"},
		{"assesses:inLanguage", "de"},
		{"assesses:id", "http://awesome-skills.org/2"},
		{"assesses:prefLabel", "Grammatik"},
		{"assesses:inLanguage", "de"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.Assesses))

	assert.Equal("http://awesome-skills.org/1", amb.Assesses[0].ID)
	assert.Equal("Hörverständnis", amb.Assesses[0].PrefLabel)
	assert.Equal("de", amb.Assesses[0].InLanguage)

	assert.Equal("http://awesome-skills.org/2", amb.Assesses[1].ID)
	assert.Equal("Grammatik", amb.Assesses[1].PrefLabel)
	assert.Equal("de", amb.Assesses[1].InLanguage)
}

func TestNostrToAMB_CompetencyRequired(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"competencyRequired:id", "http://awesome-skills.org/1"},
		{"competencyRequired:prefLabel", "Basisvokabular"},
		{"competencyRequired:inLanguage", "de"},
		{"competencyRequired:id", "http://awesome-skills.org/2"},
		{"competencyRequired:prefLabel", "Grundkenntnisse"},
		{"competencyRequired:inLanguage", "de"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.CompetencyRequired))

	assert.Equal("http://awesome-skills.org/1", amb.CompetencyRequired[0].ID)
	assert.Equal("Basisvokabular", amb.CompetencyRequired[0].PrefLabel)
	assert.Equal("de", amb.CompetencyRequired[0].InLanguage)

	assert.Equal("http://awesome-skills.org/2", amb.CompetencyRequired[1].ID)
	assert.Equal("Grundkenntnisse", amb.CompetencyRequired[1].PrefLabel)
	assert.Equal("de", amb.CompetencyRequired[1].InLanguage)
}

func TestNostrToAMB_EducationalLevel(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"educationalLevel:id", "https://w3id.org/kim/educationalLevel/level_2"},
		{"educationalLevel:prefLabel", "Sekundarstufe 1"},
		{"educationalLevel:inLanguage", "de"},
		{"educationalLevel:id", "https://w3id.org/kim/educationalLevel/level_3"},
		{"educationalLevel:prefLabel", "Sekundarstufe 2"},
		{"educationalLevel:inLanguage", "de"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.EducationalLevel))

	assert.Equal("https://w3id.org/kim/educationalLevel/level_2", amb.EducationalLevel[0].ID)
	assert.Equal("Sekundarstufe 1", amb.EducationalLevel[0].PrefLabel)
	assert.Equal("de", amb.EducationalLevel[0].InLanguage)

	assert.Equal("https://w3id.org/kim/educationalLevel/level_3", amb.EducationalLevel[1].ID)
	assert.Equal("Sekundarstufe 2", amb.EducationalLevel[1].PrefLabel)
	assert.Equal("de", amb.EducationalLevel[1].InLanguage)
}

func TestNostrToAMB_InteractivityType(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"interactivityType:id", "http://purl.org/dcx/lrmi-vocabs/interactivityType/active"},
		{"interactivityType:prefLabel", "aktiv"},
		{"interactivityType:inLanguage", "de"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.NotNil(amb.InteractivityType)
	assert.Equal("http://purl.org/dcx/lrmi-vocabs/interactivityType/active", amb.InteractivityType.ID)
	assert.Equal("aktiv", amb.InteractivityType.PrefLabel)
	assert.Equal("de", amb.InteractivityType.InLanguage)
}

func TestNostrToAMB_IsBasedOn(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"isBasedOn:id", "http://an-awesome-resource.org"},
		{"isBasedOn:name", "Französisch I"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(1, len(amb.IsBasedOn))
	assert.Equal("http://an-awesome-resource.org", amb.IsBasedOn[0].ID)
	assert.Equal("Französisch I", amb.IsBasedOn[0].Name)
}

func TestNostrToAMB_IsPartOf(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"isPartOf:id", "http://whole.org"},
		{"isPartOf:name", "Whole"},
		{"isPartOf:type", "PresentationDigitalDocument"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(1, len(amb.IsPartOf))
	assert.Equal("http://whole.org", amb.IsPartOf[0].ID)
	assert.Equal("Whole", amb.IsPartOf[0].Name)
	assert.Equal("PresentationDigitalDocument", amb.IsPartOf[0].Type)
}

func TestNostrToAMB_HasPart(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"hasPart:id", "http://part1.org"},
		{"hasPart:name", "Part 1"},
		{"hasPart:type", "LearningResource"},
		{"hasPart:id", "http://part2.org"},
		{"hasPart:name", "Part 2"},
		{"hasPart:type", "LearningResource"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(2, len(amb.HasPart))
	assert.Equal("http://part1.org", amb.HasPart[0].ID)
	assert.Equal("Part 1", amb.HasPart[0].Name)
	assert.Equal("LearningResource", amb.HasPart[0].Type)

	assert.Equal("http://part2.org", amb.HasPart[1].ID)
	assert.Equal("Part 2", amb.HasPart[1].Name)
	assert.Equal("LearningResource", amb.HasPart[1].Type)
}

func TestNostrToAMB_Duration(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"duration", "PT30M"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal("PT30M", amb.Duration)
}

func TestNostrToAMB_Trailer(t *testing.T) {
	assert := assert.New(t)

	tags := nostr.Tags{
		{"d", "test-resource-id"},
		{"trailer:contentUrl", "https://example.com/video.mp4"},
		{"trailer:type", "VideoObject"},
		{"trailer:encodingFormat", "video/mp4"},
		{"trailer:contentSize", "10MB"},
		{"trailer:sha256", "abc123"},
		{"trailer:embedUrl", "https://example.com/embed"},
		{"trailer:bitrate", "1Mbps"},
	}
	event := createTestEvent(tags)

	amb, err := NostrToAMB(event)

	assert.NoError(err)
	assert.NotNil(amb)

	assert.Equal(1, len(amb.Trailer))
	assert.Equal("https://example.com/video.mp4", amb.Trailer[0].ContentUrl)
	assert.Equal("VideoObject", amb.Trailer[0].Type)
	assert.Equal("video/mp4", amb.Trailer[0].EncodingFormat)
	assert.Equal("10MB", amb.Trailer[0].ContentSize)
	assert.Equal("abc123", amb.Trailer[0].Sha256)
	assert.Equal("https://example.com/embed", amb.Trailer[0].EmbedUrl)
	assert.Equal("1Mbps", amb.Trailer[0].Bitrate)
}
