package openalex

import (
	"testing"
)

// Test std and gz file
func TestParseFile(t *testing.T) {
	var tests = []struct {
		data       string
		path       string
		wantLength int
	}{
		{"stdWorks", "../../sample/openalex/works/updated_date=2023-05-16/part_000", 27},
		{"gzWorks", "../../sample/openalex/works/updated_date=2023-05-16/part_000.gz", 27},
		{"gzAuthors", "../../sample/openalex/authors/updated_date=2023-04-21/part_000.gz", 50},
		{"gzAuthors", "../../sample/openalex/authors/updated_date=2023-04-21/part_000", 50},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			_, err := ParseFile(tt.path, PrintEntityHandler)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestParseWork(t *testing.T) {
	workSamplePath := "../../sample/openalex/works/W2741809807"
	_, err := ParseFile(workSamplePath, PrintEntityHandler)
	if err != nil {
		t.Error(err)
	}
}
