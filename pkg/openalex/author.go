package openalex

// Author is a struct that represents the data of an author of OpenAlex
type Author struct {
	ID                      string                     `json:"id"`
	CitedByCount            int                        `json:"cited_by_count"`
	CountsByYear            []AuthorCountsByYear       `json:"counts_by_year"`
	CreatedDate             string                     `json:"created_date"`
	DisplayName             string                     `json:"display_name"`
	DisplayNameAlternatives []string                   `json:"display_name_alternatives"`
	Ids                     AuthorIDs                  `json:"ids"`
	LastKnownInstitution    AuthorLastKnownInstitution `json:"last_known_institution"`
	MostCitedWork           string                     `json:"most_cited_work"`
	Orcid                   string                     `json:"orcid"`
	SummaryStats            AuthorSummaryStats         `json:"summary_stats"`
	UpdatedDate             string                     `json:"updated_date"`
	WorksAPIURL             string                     `json:"works_api_url"`
	WorksCount              int                        `json:"works_count"`
	XConcepts               []AuthorXConcept           `json:"x_concepts"`
}

type AuthorCountsByYear struct {
	Year         int `json:"year"`
	WorksCount   int `json:"works_count"`
	CitedByCount int `json:"cited_by_count"`
}

type AuthorIDs struct {
	Openalex  string `json:"openalex"`
	Orcid     string `json:"orcid"`
	Scopus    string `json:"scopus"`
	Wikipedia string `json:"wikipedia"`
}

type AuthorLastKnownInstitution struct {
	CountryCode string `json:"country_code"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Ror         string `json:"ror"`
	Type        string `json:"type"`
}

type AuthorSummaryStats struct {
	CitedByCount2yr  int     `json:"2yr_cited_by_count" graphql:"two_year_cited_by_count"`
	HIndex2yr        int     `json:"2yr_h_index" graphql:"two_year_h_index"`
	I10Index2yr      int     `json:"2yr_i10_index" graphql:"two_year_i10_index"`
	MeanCitedness2yr float64 `json:"2yr_mean_citedness" graphql:"two_year_mean_citedness"`
	WorksCount2yr    int     `json:"2yr_works_count" graphql:"two_year_works_count"`
	CitedByCount     int     `json:"cited_by_count"`
	HIndex           int     `json:"h_index"`
	I10Index         int     `json:"i10_index"`
	OaPercent        float64 `json:"oa_percent"`
	WorksCount       int     `json:"works_count"`
}

type AuthorXConcept struct {
	DisplayName string  `json:"display_name"`
	ID          string  `json:"id"`
	Level       float64 `json:"level"`
	Score       float64 `json:"score"`
	Wikidata    string  `json:"wikidata"`
}

// GetID returns the ID of the author
func (a *Author) GetID() string {
	return a.ID
}

// GetType returns the entity type
func (a *Author) GetType() string {
	return string(WorksFileEntityType)
}
