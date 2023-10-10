package openalex

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
)

// ManifestUrl is a type that represents a manifest URL
type ManifestUrl string

const (
	ManifestUrlAuthors      ManifestUrl = "https://openalex.s3.amazonaws.com/data/authors/manifest"
	ManifestUrlConcepts     ManifestUrl = "https://openalex.s3.amazonaws.com/data/concepts/manifest"
	ManifestUrlFunders      ManifestUrl = "https://openalex.s3.amazonaws.com/data/funders/manifest"
	ManifestUrlInstitutions ManifestUrl = "https://openalex.s3.amazonaws.com/data/institutions/manifest"
	ManifestUrlPublishers   ManifestUrl = "https://openalex.s3.amazonaws.com/data/publishers/manifest"
	ManifestUrlSources      ManifestUrl = "https://openalex.s3.amazonaws.com/data/sources/manifest"
	ManifestUrlWorks        ManifestUrl = "https://openalex.s3.amazonaws.com/data/works/manifest"
)

// AllManifestUrls is a list of all manifest URLs
var AllManifestUrls = []ManifestUrl{
	ManifestUrlAuthors,
	ManifestUrlConcepts,
	ManifestUrlFunders,
	ManifestUrlInstitutions,
	ManifestUrlPublishers,
	ManifestUrlSources,
	ManifestUrlWorks,
}

// Manifest is a struct that represents the manifest file
type Manifest struct {
	Entries []struct {
		URL  string `json:"url"`
		Meta struct {
			ContentLength int `json:"content_length"`
			RecordCount   int `json:"record_count"`
		} `json:"meta"`
	} `json:"entries"`
	Meta struct {
		ContentLength int64 `json:"content_length"`
		RecordCount   int   `json:"record_count"`
	} `json:"meta"`
}

// Hash returns the SHA256 hash of the manifest
func (m *Manifest) Hash() (result string, err error) {
	data, err := json.Marshal(m)
	if err != nil {
		slog.With("error", err).Error("Failed to marshal manifest")
		return
	}
	result = fmt.Sprintf("%x", sha256.Sum256(data))
	return
}