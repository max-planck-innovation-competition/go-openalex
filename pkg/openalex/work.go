package openalex

import (
	"sort"
	"strings"
)

// AbstractInvertedIndex is a struct that represents the inverted index of an abstract
type AbstractInvertedIndex struct {
	IndexLength   int              `json:"IndexLength"`
	InvertedIndex map[string][]int `json:"InvertedIndex,omitempty"`
}

func (abstractInvertedIndex *AbstractInvertedIndex) ToAbstract() string {
	if abstractInvertedIndex == nil {
		return ""
	}
	// Create a list of (word, index) pairs.
	var wordIndex [][]interface{}
	for word, indices := range abstractInvertedIndex.InvertedIndex {
		for _, idx := range indices {
			wordIndex = append(wordIndex, []interface{}{word, idx})
		}
	}

	// Sort the list by index.
	sort.Slice(wordIndex, func(i, j int) bool {
		return wordIndex[i][1].(int) < wordIndex[j][1].(int)
	})

	// Create a list of words in index order.
	var words []string
	for _, pair := range wordIndex {
		words = append(words, pair[0].(string))
	}

	// Join the words with spaces and return as a string.
	return strings.Join(words, " ")
}

type Work struct {
	ID                    string                `json:"id"`
	Abstract              string                `json:"abstract"`
	AbstractInvertedIndex AbstractInvertedIndex `json:"abstract_inverted_index"`
	Authorships           []struct {
		Author struct {
			DisplayName string  `json:"display_name"`
			ID          string  `json:"id"`
			Orcid       *string `json:"orcid"`
		} `json:"author"`
		AuthorPosition string `json:"author_position"`
		Institutions   []struct {
			CountryCode *string `json:"country_code"`
			DisplayName string  `json:"display_name"`
			ID          *string `json:"id"`
			Ror         *string `json:"ror"`
			Type        *string `json:"type"`
		} `json:"institutions"`
		RawAffiliationString *string `json:"raw_affiliation_string"`
	} `json:"authorships"`
	Biblio struct {
		FirstPage *string `json:"first_page"`
		Issue     *string `json:"issue"`
		LastPage  *string `json:"last_page"`
		Volume    *string `json:"volume"`
	} `json:"biblio"`
	CitedByAPIURL string `json:"cited_by_api_url"`
	CitedByCount  int    `json:"cited_by_count"`
	Concepts      []struct {
		DisplayName string  `json:"display_name"`
		ID          string  `json:"id"`
		Level       int     `json:"level"`
		Score       float64 `json:"score"`
		Wikidata    string  `json:"wikidata"`
	} `json:"concepts"`
	CountsByYear []struct {
		CitedByCount int `json:"cited_by_count"`
		Year         int `json:"year"`
	} `json:"counts_by_year"`
	CreatedDate string `json:"created_date"`
	DisplayName string `json:"display_name"`
	Doi         string `json:"doi"`
	HostVenue   struct {
		DisplayName *string  `json:"display_name"`
		ID          *string  `json:"id"`
		IsOa        *bool    `json:"is_oa"`
		Issn        []string `json:"issn"`
		IssnL       *string  `json:"issn_l"`
		License     *string  `json:"license"`
		Publisher   *string  `json:"publisher"`
		Type        string   `json:"type"`
		URL         string   `json:"url"`
		Version     *string  `json:"version"`
	} `json:"host_venue"`
	Ids struct {
		Doi      string `json:"doi"`
		Mag      int    `json:"mag"`
		Openalex string `json:"openalex"`
		Pmid     string `json:"pmid,omitempty"`
	} `json:"ids"`
	IsParatext  bool `json:"is_paratext"`
	IsRetracted bool `json:"is_retracted"`
	Mesh        []struct {
		DescriptorName string  `json:"descriptor_name"`
		DescriptorUi   string  `json:"descriptor_ui"`
		IsMajorTopic   bool    `json:"is_major_topic"`
		QualifierName  *string `json:"qualifier_name"`
		QualifierUi    *string `json:"qualifier_ui"`
	} `json:"mesh"`
	OpenAccess struct {
		IsOa     bool    `json:"is_oa"`
		OaStatus string  `json:"oa_status"`
		OaURL    *string `json:"oa_url"`
	} `json:"open_access"`
	PublicationDate string   `json:"publication_date"`
	PublicationYear int      `json:"publication_year"`
	ReferencedWorks []string `json:"referenced_works"`
	RelatedWorks    []string `json:"related_works"`
	Title           string   `json:"title"`
	Type            string   `json:"type"`
	Updated         string   `json:"updated"`
	UpdatedDate     string   `json:"updated_date"`
}

// GetID returns the ID of the work
func (w *Work) GetID() string {
	return w.ID
}

// GetType returns the entity type
func (w *Work) GetType() string {
	return string(WorksFileEntityType)
}
