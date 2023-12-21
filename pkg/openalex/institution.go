package openalex

import jsoniter "github.com/json-iterator/go"

// Institution is a struct that represents the JSON response from the OpenAlex API.
type Institution struct {
	ID                     string `json:"id"`
	AssociatedInstitutions []struct {
		CountryCode  string `json:"country_code"`
		DisplayName  string `json:"display_name"`
		ID           string `json:"id"`
		Relationship string `json:"relationship"`
		Ror          string `json:"ror"`
		Type         string `json:"type"`
	} `json:"associated_institutions"`
	CitedByCount int    `json:"cited_by_count"`
	CountryCode  string `json:"country_code"`
	CountsByYear []struct {
		CitedByCount int `json:"cited_by_count"`
		WorksCount   int `json:"works_count"`
		Year         int `json:"year"`
	} `json:"counts_by_year"`
	CreatedDate             string   `json:"created_date"`
	DisplayName             string   `json:"display_name"`
	DisplayNameAcronyms     []string `json:"display_name_acronyms"`
	DisplayNameAlternatives []string `json:"display_name_alternatives"`
	Geo                     struct {
		City           string   `json:"city"`
		Country        string   `json:"country"`
		CountryCode    string   `json:"country_code"`
		GeonamesCityID string   `json:"geonames_city_id"`
		Latitude       *float64 `json:"latitude"`
		Longitude      *float64 `json:"longitude"`
		Region         *string  `json:"region"`
	} `json:"geo"`
	HomepageURL *string `json:"homepage_url"`
	Ids         struct {
		Grid      string          `json:"grid"`
		Mag       jsoniter.Number `json:"mag,omitempty"`
		Openalex  string          `json:"openalex"`
		Ror       string          `json:"ror"`
		Wikidata  string          `json:"wikidata,omitempty"`
		Wikipedia string          `json:"wikipedia,omitempty"`
	} `json:"ids"`
	ImageThumbnailURL *string `json:"image_thumbnail_url"`
	ImageURL          *string `json:"image_url"`
	International     struct {
		DisplayName map[string]string `json:"display_name"`
	} `json:"international"`
	Lineage      []string `json:"lineage"`
	Repositories []struct {
		DisplayName             string   `json:"display_name"`
		HostOrganization        string   `json:"host_organization"`
		HostOrganizationLineage []string `json:"host_organization_lineage"`
		HostOrganizationName    string   `json:"host_organization_name"`
		ID                      string   `json:"id"`
	} `json:"repositories"`
	Roles []struct {
		ID         string `json:"id"`
		Role       string `json:"role"`
		WorksCount int    `json:"works_count"`
	} `json:"roles"`
	Ror          string `json:"ror"`
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
	Type        string `json:"type"`
	UpdatedDate string `json:"updated_date"`
	WorksAPIURL string `json:"works_api_url"`
	WorksCount  int    `json:"works_count"`
	XConcepts   []struct {
		DisplayName string  `json:"display_name"`
		ID          string  `json:"id"`
		Level       int     `json:"level"`
		Score       float64 `json:"score"`
		Wikidata    string  `json:"wikidata"`
	} `json:"x_concepts"`
}

// GetID returns the ID of the institution
func (i *Institution) GetID() string {
	return i.ID
}

// GetType returns the entity type
func (i *Institution) GetType() string {
	return string(InstitutionsFileEntityType)
}
