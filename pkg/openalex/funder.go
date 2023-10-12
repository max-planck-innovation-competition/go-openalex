package openalex

// Funder is a struct that represents the JSON response from the OpenAlex API.
type Funder struct {
	ID              string   `json:"id"`
	AlternateTitles []string `json:"alternate_titles"`
	CitedByCount    int      `json:"cited_by_count"`
	CountryCode     string   `json:"country_code"`
	CountsByYear    []struct {
		CitedByCount int `json:"cited_by_count"`
		OaWorksCount int `json:"oa_works_count"`
		WorksCount   int `json:"works_count"`
		Year         int `json:"year"`
	} `json:"counts_by_year"`
	CreatedDate string  `json:"created_date"`
	Description *string `json:"description"`
	DisplayName string  `json:"display_name"`
	GrantsCount int     `json:"grants_count"`
	HomepageURL *string `json:"homepage_url"`
	Ids         struct {
		Crossref int    `json:"crossref"`
		Doi      string `json:"doi"`
		Openalex string `json:"openalex"`
		Wikidata string `json:"wikidata,omitempty"`
	} `json:"ids"`
	ImageThumbnailURL any `json:"image_thumbnail_url"`
	ImageURL          any `json:"image_url"`
	Roles             []struct {
		ID         string `json:"id"`
		Role       string `json:"role"`
		WorksCount int    `json:"works_count"`
	} `json:"roles"`
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
	Updated     string `json:"updated"`
	UpdatedDate string `json:"updated_date"`
	WorksCount  int    `json:"works_count"`
	XConcepts   []struct {
		DisplayName string  `json:"display_name"`
		ID          string  `json:"id"`
		Level       int     `json:"level"`
		Score       float64 `json:"score"`
		Wikidata    string  `json:"wikidata"`
	} `json:"x_concepts"`
}

// GetID returns the ID of the funder
func (f *Funder) GetID() string {
	return f.ID
}

// GetType returns the entity type
func (f *Funder) GetType() string {
	return string(FundersFileEntityType)
}
