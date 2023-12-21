package openalex

import jsoniter "github.com/json-iterator/go"

// Concept is a struct that represents a concept in OpenAlex
type Concept struct {
	ID        string `json:"id"`
	Ancestors []struct {
		DisplayName string `json:"display_name"`
		ID          string `json:"id"`
		Level       int    `json:"level"`
		Wikidata    string `json:"wikidata"`
	} `json:"ancestors"`
	CitedByCount int `json:"cited_by_count"`
	CountsByYear []struct {
		CitedByCount int `json:"cited_by_count"`
		OaWorksCount int `json:"oa_works_count"`
		WorksCount   int `json:"works_count"`
		Year         int `json:"year"`
	} `json:"counts_by_year"`
	CreatedDate string `json:"created_date"`
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
	Ids         struct {
		Mag       jsoniter.Number `json:"mag"`
		Openalex  string          `json:"openalex"`
		UmlsCui   []string        `json:"umls_cui,omitempty"`
		Wikidata  string          `json:"wikidata"`
		Wikipedia string          `json:"wikipedia"`
	} `json:"ids"`
	ImageThumbnailURL *string `json:"image_thumbnail_url"`
	ImageURL          *string `json:"image_url"`
	International     struct {
		Description map[string]string `json:"description"`
		DisplayName map[string]string `json:"display_name"`
	} `json:"international"`
	Level           int `json:"level"`
	RelatedConcepts []struct {
		DisplayName string  `json:"display_name"`
		ID          string  `json:"id"`
		Level       int     `json:"level"`
		Score       float64 `json:"score"`
		Wikidata    any     `json:"wikidata"` // TODO: replace any with struct
	} `json:"related_concepts"`
	SummaryStats struct {
		CitedByCount2yr  int     `json:"2yr_cited_by_count"`
		HIndex2yr        int     `json:"2yr_h_index"`
		I10Index2yr      int     `json:"2yr_i10_index"`
		MeanCitedness2yr float64 `json:"2yr_mean_citedness"`
		WorksCount2yr    int     `json:"2yr_works_count"`
		CitedByCount     int     `json:"cited_by_count"`
		HIndex           int     `json:"h_index"`
		I10Index         int     `json:"i10_index"`
		OaPercent        float64 `json:"oa_percent"`
		WorksCount       int     `json:"works_count"`
	} `json:"summary_stats"`
	UpdatedDate string `json:"updated_date"`
	Wikidata    string `json:"wikidata"`
	WorksAPIURL string `json:"works_api_url"`
	WorksCount  int    `json:"works_count"`
}

// GetID returns the ID of the concept
func (c *Concept) GetID() string {
	return c.ID
}

// GetType returns the entity type
func (c *Concept) GetType() string {
	return string(WorksFileEntityType)
}
