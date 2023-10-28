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

func TestParseAuthor(t *testing.T) {
	authorSamplePath := "../../sample/openalex/authors/A5023888391"
	_, err := ParseFile(authorSamplePath, PrintEntityHandler)
	if err != nil {
		t.Error(err)
	}
}

func TestParseSource(t *testing.T) {
	sourceSamplePath := "../../sample/openalex/sources/S137773608"
	_, err := ParseFile(sourceSamplePath, PrintEntityHandler)
	if err != nil {
		t.Error(err)
	}
}

func TestParseInstitution(t *testing.T) {
	institutionSamplePath := "../../sample/openalex/institutions/I27837315"
	_, err := ParseFile(institutionSamplePath, PrintEntityHandler)
	if err != nil {
		t.Error(err)
	}
}
