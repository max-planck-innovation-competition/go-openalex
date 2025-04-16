package openalex

import (
	"testing"
)

func TestParseEntity(t *testing.T) {

	workSamplePath := "../../sample/openalex/works/W2741809807"

	lineData := map[string]string{"id": "https://openalex.org/W2741809807", "doi": "https://doi.org/10.7717/peerj.4375", "title": "The state of OA: a large-scale analysis of the prevalence and impact of Open Access articles", "display_name": "The state of OA: a large-scale analysis of the prevalence and impact of Open Access articles", "publication_year": "2018", "publication_date": "2018-02-13"}

	jsonBytes, err := json.Marshal(lineData)
	if err != nil {
		panic(err)
	}

	lineSample := string(jsonBytes)

	data, err := ParseEntity(workSamplePath, lineSample)
	if err != nil {
		print(data)
		t.Error(err)
	}
}
