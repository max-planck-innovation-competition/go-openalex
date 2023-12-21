package openalex

type Source struct {
	ID               string   `json:"id"`
	AbbreviatedTitle *string  `json:"abbreviated_title"`
	AlternateTitles  []string `json:"alternate_titles"`
	ApcPrices        []struct {
		Currency string `json:"currency"`
		Price    int    `json:"price"`
	} `json:"apc_prices"`
	ApcUsd       *int    `json:"apc_usd"`
	CitedByCount int     `json:"cited_by_count"`
	CountryCode  *string `json:"country_code"`
	CountsByYear []struct {
		CitedByCount int `json:"cited_by_count"`
		OaWorksCount int `json:"oa_works_count"`
		WorksCount   int `json:"works_count"`
		Year         int `json:"year"`
	} `json:"counts_by_year"`
	CreatedDate                  string   `json:"created_date"`
	DisplayName                  string   `json:"display_name"`
	HomepageURL                  *string  `json:"homepage_url"`
	HostInstitutionLineage       []string `json:"host_institution_lineage"`
	HostInstitutionLineageNames  []string `json:"host_institution_lineage_names"`
	HostOrganization             *string  `json:"host_organization"`
	HostOrganizationLineage      []string `json:"host_organization_lineage"`
	HostOrganizationLineageNames []string `json:"host_organization_lineage_names"`
	HostOrganizationName         *string  `json:"host_organization_name"`
	Ids                          struct {
		Fatcat   string   `json:"fatcat,omitempty"`
		Issn     []string `json:"issn,omitempty"`
		IssnL    string   `json:"issn_l,omitempty"`
		Mag      int      `json:"mag,omitempty"`
		Openalex string   `json:"openalex"`
		Wikidata string   `json:"wikidata,omitempty"`
	} `json:"ids"`
	IsInDoaj              bool     `json:"is_in_doaj"`
	IsOa                  bool     `json:"is_oa"`
	Issn                  []string `json:"issn"`
	IssnL                 *string  `json:"issn_l"`
	Publisher             *string  `json:"publisher"`
	PublisherID           *string  `json:"publisher_id"`
	PublisherLineage      []string `json:"publisher_lineage"`
	PublisherLineageNames []string `json:"publisher_lineage_names"`
	Societies             []any    `json:"societies"` // TODO: replace any with struct
	SummaryStats          struct {
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
	Updated     string `json:"updated"`
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

// GetID returns the ID of the source
func (s *Source) GetID() string {
	return s.ID
}

// GetType returns the entity type
func (s *Source) GetType() string {
	return string(SourcesFileEntityType)
}
