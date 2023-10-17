package openalex

// Author is a struct that represents the data of an author of OpenAlex
type Author struct {
	ID                      string `json:"id"`
	CitedByCount            int    `json:"cited_by_count"`
	CountsByYear            []any  `json:"counts_by_year"`
	CreatedDate             string `json:"created_date"`
	DisplayName             string `json:"display_name"`
	DisplayNameAlternatives []any  `json:"display_name_alternatives"`
	Ids                     struct {
		Openalex string `json:"openalex"`
	} `json:"ids"`
	LastKnownInstitution struct {
		CountryCode string `json:"country_code"`
		DisplayName string `json:"display_name"`
		ID          string `json:"id"`
		Ror         string `json:"ror"`
		Type        string `json:"type"`
	} `json:"last_known_institution"`
	MostCitedWork any `json:"most_cited_work"`
	Orcid         any `json:"orcid"`
	SummaryStats  struct {
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
	WorksAPIURL string `json:"works_api_url"`
	WorksCount  int    `json:"works_count"`
	XConcepts   []any  `json:"x_concepts"`
}

// GetID returns the ID of the author
func (a *Author) GetID() string {
	return a.ID
}

// GetType returns the entity type
func (a *Author) GetType() string {
	return string(WorksFileEntityType)
}
