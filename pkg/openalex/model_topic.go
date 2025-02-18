package openalex

type Topic struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Subfield    struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"subfield"`
	Field struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"field"`
	Domain struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"domain"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Ids         struct {
		Openalex  string `json:"openalex"`
		Wikipedia string `json:"wikipedia"`
	} `json:"ids"`
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

// GetID returns the ID of the topic
func (t *Topic) GetID() string {
	return t.ID
}

// GetType returns the entity type
func (t *Topic) GetType() string {
	return string(TopicsFileEntityType)
}
