package openalex

import (
	"fmt"
	"testing"
)

func TestParseMergedIDsFile(t *testing.T) {

	var tests = []struct {
		data string
		path string
	}{
		{"gz", "../../sample/openalex/merged_ids/authors/2022-08-03.csv.gz"},
		{"csv", "../../sample/openalex/merged_ids/authors/2023-04-13.csv"},
		{"gz2", "../../sample/openalex/merged_ids/authors/2023-04-13.csv.gz"},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			res, err := ParseMergedIDsFile(tt.path)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(res)
		})
	}

}
