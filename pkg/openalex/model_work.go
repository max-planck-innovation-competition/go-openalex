package openalex

import (
	"sort"
	"strings"
)

func TransformInvertedIndexToAbstract(invertedIndex map[string][]int) string {
	if invertedIndex == nil {
		return ""
	}
	if len(invertedIndex) == 0 {
		return ""
	}
	// Create a list of (word, index) pairs.
	var wordIndex [][]interface{}
	for word, indices := range invertedIndex {
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

func (w *Work) ToAbstract() string {
	return TransformInvertedIndexToAbstract(w.AbstractInvertedIndex)
}

// GenerateAbstractFromInvertedIndex sets the abstract
func (w *Work) GenerateAbstractFromInvertedIndex() *Work {
	w.Abstract = w.ToAbstract()
	return w
}

// SetAbstractInvertedIndexToNil sets the abstract inverted index to nil
func (w *Work) SetAbstractInvertedIndexToNil() *Work {
	w.AbstractInvertedIndex = nil
	return w
}

// Work is the struct for a work in the open alex database
type Work struct {
	ID                    string           `json:"id"`
	Abstract              string           `json:"abstract"`
	AbstractInvertedIndex map[string][]int `json:"abstract_inverted_index,omitempty"`
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
		RawAffiliationString *string  `json:"raw_affiliation_string"`
		RawAuthorName        *string  `json:"raw_author_name"`
		IsCorresponding      *bool    `json:"is_corresponding"`
		Countries            []string `json:"countries"`
	} `json:"authorships"`
	ApcList struct {
		Value      int     `json:"value"`
		Currency   *string `json:"currency"`
		Provenance *string `json:"provenance"`
		ValueUsd   int     `json:"value_usd"`
	} `json:"apc_list"`
	BestOALocation struct {
		IsOA           *bool   `json:"is_oa"`
		LandingPageUrl *string `json:"landing_page_url"`
		PdfUrl         *string `json:"pdf_url"`
		License        *string `json:"license"`
		Version        *string `json:"version"`
		Source         struct {
			Id               *string  `json:"id"`
			DisplayName      *string  `json:"display_name"`
			IssnL            *string  `json:"issn_l"`
			Issn             []string `json:"issn"`
			HostOrganization *string  `json:"host_organization"`
			Type             *string  `json:"type"`
		} `json:"source"`
	} `json:"best_oa_location"`
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
	CorrespondingAuthorIds      []string `json:"corresponding_author_ids"`
	CorrespondingInstitutionIds []string `json:"corresponding_institution_ids"`
	CountriesDistinctCount      int      `json:"countries_distinct_count"`
	CountsByYear                []struct {
		CitedByCount int `json:"cited_by_count"`
		Year         int `json:"year"`
	} `json:"counts_by_year"`
	CreatedDate string `json:"created_date"`
	DisplayName string `json:"display_name"`
	Doi         string `json:"doi"`
	Grants      []struct {
		Funder            string `json:"funder"`
		FunderDisplayName string `json:"funder_display_name"`
		AwardId           string `json:"award_id"`
	}
	HasFulltext               *bool  `json:"has_fulltext"`
	InstitutionsDistinctCount int    `json:"institutions_distinct_count"`
	Language                  string `json:"language"`
	Locations                 []struct {
		IsOA           *bool   `json:"is_oa"`
		LandingPageUrl *string `json:"landing_page_url"`
		PdfUrl         *string `json:"pdf_url"`
		Source         struct {
			Id               *string  `json:"id"`
			DisplayName      *string  `json:"display_name"`
			IssnL            *string  `json:"issn_l"`
			Issn             []string `json:"issn"`
			HostOrganization *string  `json:"host_organization"`
			Type             *string  `json:"type"`
		} `json:"source"`
		License *string `json:"license"`
		Version *string `json:"version"`
	} `json:"locations"`
	PrimaryLocation struct {
		IsOA           *bool   `json:"is_oa"`
		LandingPageUrl *string `json:"landing_page_url"`
		PdfUrl         *string `json:"pdf_url"`
		Source         struct {
			Id               *string  `json:"id"`
			DisplayName      *string  `json:"display_name"`
			IssnL            *string  `json:"issn_l"`
			Issn             []string `json:"issn"`
			HostOrganization *string  `json:"host_organization"`
			Type             *string  `json:"type"`
		} `json:"source"`
		License *string `json:"license"`
		Version *string `json:"version"`
	} `json:"primary_location"`
	LocationCount int `json:"location_count"`
	Ids           struct {
		Doi      string `json:"doi"`
		Openalex string `json:"openalex"`
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
		IsOa                     bool    `json:"is_oa"`
		OaStatus                 string  `json:"oa_status"`
		OaURL                    *string `json:"oa_url"`
		AnyRepositoryHasFulltext bool    `json:"any_repository_has_fulltext"`
	} `json:"open_access"`
	PublicationDate string   `json:"publication_date"`
	PublicationYear int      `json:"publication_year"`
	ReferencedWorks []string `json:"referenced_works"`
	RelatedWorks    []string `json:"related_works"`
	Title           string   `json:"title"`
	Type            string   `json:"type"`
	UpdatedDate     string   `json:"updated_date"`
	TypeCrossref    string   `json:"type_crossref"`
}

// GetID returns the ID of the work
func (w *Work) GetID() string {
	return w.ID
}

// GetType returns the entity type
func (w *Work) GetType() string {
	return string(WorksFileEntityType)
}
