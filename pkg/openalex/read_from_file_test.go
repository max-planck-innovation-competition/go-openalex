package openalex

import (
	"fmt"
	"testing"
)

// Test std and gz file
func TestParseFile(t *testing.T) {

	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	var tests = []struct {
		data       string
		path       string
		wantLength int
	}{
		{"stdWorks", "C:/DOCDB/openalex/data/works/updated_date=2024-06-30/part_043.gz", 27},
		/*	{"stdWorks", "../../sample/openalex/works/updated_date=2023-05-16/part_000", 27},
			{"gzWorks", "../../sample/openalex/works/updated_date=2023-05-16/part_000.gz", 27},
			{"gzAuthors", "../../sample/openalex/authors/updated_date=2023-04-21/part_000.gz", 50},
			{"gzAuthors", "../../sample/openalex/authors/updated_date=2023-04-21/part_000", 50}, */
	}

	p := Processor{}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			_, err := p.ParseFile(tt.path)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestParseWork(t *testing.T) {
	p := Processor{}

	workSamplePath := "../../sample/openalex/works/W2741809807"
	_, err := p.ParseFile(workSamplePath)
	if err != nil {
		t.Error(err)
	}
}

func TestParseAuthor(t *testing.T) {
	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	p := Processor{}
	authorSamplePath := "../../sample/openalex/authors/A5023888391"
	_, err := p.ParseFile(authorSamplePath)
	if err != nil {
		t.Error(err)
	}
}

func TestParseSource(t *testing.T) {
	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	p := Processor{}
	sourceSamplePath := "../../sample/openalex/sources/S137773608"
	_, err := p.ParseFile(sourceSamplePath)
	if err != nil {
		t.Error(err)
	}
}

func TestParseInstitution(t *testing.T) {
	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	p := Processor{}
	institutionSamplePath := "../../sample/openalex/institutions/I27837315"
	_, err := p.ParseFile(institutionSamplePath)
	if err != nil {
		t.Error(err)
	}
}

func TestParseConcept(t *testing.T) {
	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	p := Processor{}
	conceptSamplePath := "../../sample/openalex/concepts/C71924100"
	_, err := p.ParseFile(conceptSamplePath)
	if err != nil {
		t.Error(err)
	}
}

func TestParsePublisher(t *testing.T) {
	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	p := Processor{}
	publisherSamplePath := "../../sample/openalex/publishers/P4310319965"
	_, err := p.ParseFile(publisherSamplePath)
	if err != nil {
		t.Error(err)
	}
}

func TestParseFunder(t *testing.T) {
	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	p := Processor{}
	funderSamplePath := "../../sample/openalex/funders/F4320332161"
	_, err := p.ParseFile(funderSamplePath)
	if err != nil {
		t.Error(err)
	}
}
