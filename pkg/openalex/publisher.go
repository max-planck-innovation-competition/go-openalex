package openalex

type Publisher struct {
	ID              string   `json:"id"`
	AlternateTitles []string `json:"alternate_titles"`
	CitedByCount    int      `json:"cited_by_count"`
	CountryCodes    []string `json:"country_codes"`
	CountsByYear    []struct {
		CitedByCount int `json:"cited_by_count"`
		OaWorksCount int `json:"oa_works_count"`
		WorksCount   int `json:"works_count"`
		Year         int `json:"year"`
	} `json:"counts_by_year"`
	CreatedDate    string  `json:"created_date"`
	DisplayName    string  `json:"display_name"`
	HierarchyLevel int     `json:"hierarchy_level"`
	HomepageURL    *string `json:"homepage_url"`
	Ids            struct {
		Openalex string `json:"openalex"`
		Wikidata string `json:"wikidata,omitempty"`
		Ror      string `json:"ror"`
	} `json:"ids"`
	ImageThumbnailURL *string  `json:"image_thumbnail_url"`
	ImageURL          *string  `json:"image_url"`
	Lineage           []string `json:"lineage"`
	ParentPublisher   *string  `json:"parent_publisher"`
	Roles             []struct {
		ID         string `json:"id"`
		Role       string `json:"role"`
		WorksCount int    `json:"works_count"`
	} `json:"roles"`
	SourcesAPIURL string `json:"sources_api_url"`
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
		SourcesCount     int     `json:"sources_count"`
		WorksCount       int     `json:"works_count"`
	} `json:"summary_stats"`
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

// GetID returns the ID of the publisher
func (p *Publisher) GetID() string {
	return p.ID
}

// GetType returns the entity type
func (p *Publisher) GetType() string {
	return string(PublishersFileEntityType)
}
