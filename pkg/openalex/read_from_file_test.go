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
		{"stdWorks", "../../sample/openalex/works/updated_date=2021-11-03/part_000", 50},
		{"gzWorks", "../../sample/openalex/works/updated_date=2021-11-03/part_000.gz", 50},
		{"gzAuthors", "../../sample/openalex/authors/updated_date=2023-04-21/part_000.gz", 50},
		{"gzAuthors", "../../sample/openalex/authors/updated_date=2023-04-21/part_000", 50},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			err := ParseFile(tt.path, PrintEntityHandler)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
