package openalex

type Domain struct {
	ID                      string `json:"id"`
	DisplayName             string `json:"display_name"`
	Description             string `json:"description"`
	DisplayNameAlternatives []any  `json:"display_name_alternatives"`
	Ids                     struct {
		Wikidata  string `json:"wikidata"`
		Wikipedia string `json:"wikipedia"`
	} `json:"ids"`
	Fields []struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"fields"`
	Siblings []struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"siblings"`
	WorksCount   int    `json:"works_count"`
	CitedByCount int    `json:"cited_by_count"`
	WorksAPIURL  string `json:"works_api_url"`
	UpdatedDate  string `json:"updated_date"`
	CreatedDate  string `json:"created_date"`
	Updated      string `json:"updated"`
}

// GetID returns the ID of the domain
func (d *Domain) GetID() string {
	return d.ID
}

// GetType returns the entity type
func (d *Domain) GetType() string {
	return string(DomainsFileEntityType)
}
