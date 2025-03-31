package typesense30142

import "github.com/nbd-wtf/go-nostr"

// BaseEntity contains common fields used across many entity types
type BaseEntity struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ControlledVocabulary struct {
	Type       string `json:"type,omitempty"`
	ID         string `json:"id"`
	PrefLabel  string `json:"prefLabel"`
	InLanguage string `json:"inLanguage,omitempty"`
}

// LanguageEntity adds language support to entities
type LanguageEntity struct {
	InLanguage string `json:"inLanguage,omitempty"`
}

// LabeledEntity adds prefLabel to entities
type LabeledEntity struct {
	PrefLabel string `json:"prefLabel,omitempty"`
}

// About represents a topic or subject
type About struct {
	ControlledVocabulary
}

// Creator represents the creator of content
type Creator struct {
	BaseEntity
	Affiliation   *Affiliation `json:"affiliation,omitempty"`
	HonoricPrefix string       `json:"honoricPrefix,omitempty"`
}

// Contributor represents someone who contributed to the content
type Contributor struct {
	BaseEntity
	Affiliation   *Affiliation `json:"affiliation,omitempty"`
	HonoricPrefix string       `json:"honoricPrefix,omitempty"`
}

// Publisher represents the publisher of the content
type Publisher struct {
	BaseEntity
}

// Funder represents an entity that funded the content
type Funder struct {
	BaseEntity
}

// Affiliation represents an organization affiliation
type Affiliation struct {
	BaseEntity
}

// ConditionsOfAccess represents access conditions
type ConditionsOfAccess struct {
	ControlledVocabulary
}

// LearningResourceType categorizes the learning resource
type LearningResourceType struct {
	ControlledVocabulary
}

// Audience represents the target audience
type Audience struct {
	ControlledVocabulary
}

// Teaches represents what the content teaches
type Teaches struct {
  ControlledVocabulary
}

// Assesses represents what the content assesses
type Assesses struct {
  ControlledVocabulary
}

// CompetencyRequired represents required competencies
type CompetencyRequired struct {
  ControlledVocabulary
}

// EducationalLevel represents the educational level
type EducationalLevel struct {
	ControlledVocabulary
}

// InteractivityType represents the type of interactivity
type InteractivityType struct {
	ControlledVocabulary
}

// IsBasedOn represents a reference to source material
type IsBasedOn struct {
  ID string `json:"id"`
	Type    string   `json:"type,omitempty"`
	Name    string   `json:"name"`
	Creator *Creator `json:"creator,omitempty"`
	License *License `json:"license,omitempty"`
}

// IsPartOf represents a parent relationship
type IsPartOf struct {
	BaseEntity
}

type HasPart struct {
	BaseEntity
}

// Trailer represents a media trailer
type Trailer struct {
	Type           string `json:"type"`
	ContentUrl     string `json:"contentUrl"`
	EncodingFormat string `json:"encodingFormat"`
	ContentSize    string `json:"contentSize,omitempty"`
	Sha256         string `json:"sha256,omitempty"`
	EmbedUrl       string `json:"embedUrl,omitempty"`
	Bitrate        string `json:"bitrate,omitempty"`
}

// License represents the content license
type License struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// NostrMetadata contains Nostr-specific metadata
type NostrMetadata struct {
	EventID        string          `json:"eventID"`
	EventKind      int             `json:"eventKind"`
	EventPubKey    string          `json:"eventPubKey"`
	EventSig       string          `json:"eventSignature"`
	EventCreatedAt nostr.Timestamp `json:"eventCreatedAt"`
	EventContent   string          `json:"eventContent"`
	EventRaw       string          `json:"eventRaw"`
}

// AMBMetadata represents the full metadata structure
type AMBMetadata struct {
	// Event ID
	ID string `json:"id"`
	// Document ID
	D           string     `json:"d"`
	Type        []string   `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	About       []*About   `json:"about,omitempty"`
	Keywords    []string   `json:"keywords,omitempty"`
	InLanguage  []string   `json:"inLanguage,omitempty"`
	Image       string     `json:"image,omitempty"`
	Trailer     []*Trailer `json:"trailer,omitempty"`

	// Provenience
	Creator       []*Creator     `json:"creator,omitempty"`
	Contributor   []*Contributor `json:"contributor,omitempty"`
	DateCreated   string         `json:"dateCreated,omitempty"`
	DatePublished string         `json:"datePublished,omitempty"`
	DateModified  string         `json:"dateModified,omitempty"`
	Publisher     []*Publisher   `json:"publisher,omitempty"`
	Funder        []*Funder      `json:"funder,omitempty"`

	// Costs and Rights
	IsAccessibleForFree bool                `json:"isAccessibleForFree,omitempty"`
	License             *License            `json:"license,omitempty"`
	ConditionsOfAccess  *ConditionsOfAccess `json:"conditionsOfAccess,omitempty"`

	// Educational metadata
	LearningResourceType []*LearningResourceType `json:"learningResourceType,omitempty"`
	Audience             []*Audience             `json:"audience,omitempty"`
	Teaches              []*Teaches              `json:"teaches,omitempty"`
	Assesses             []*Assesses             `json:"assesses,omitempty"`
	CompetencyRequired   []*CompetencyRequired   `json:"competencyRequired,omitempty"`
	EducationalLevel     []*EducationalLevel     `json:"educationalLevel,omitempty"`
	InteractivityType    *InteractivityType      `json:"interactivityType,omitempty"`

	// Relation
	IsBasedOn []*IsBasedOn `json:"isBasedOn,omitempty"`
	IsPartOf  []*IsPartOf  `json:"isPartOf,omitempty"`
	HasPart   []*HasPart   `json:"hasPart,omitempty"`

	// Technical
	Duration string `json:"duration,omitempty"`
	// TODO Encoding  ``
	// TODO Caption

	// Nostr integration
	NostrMetadata `json:",inline"`
}
