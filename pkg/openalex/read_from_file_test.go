package openalex

import (
	"log"
	"testing"
)

// Test std and gz file
func TestParseFile(t *testing.T) {
	var tests = []struct {
		data        string
		path        string
		want_length int
	}{
		{"std", "../../sample/openalex/works/updated_date=2021-11-03/part_000", 50},
		{"gz", "../../sample/openalex/works/updated_date=2021-11-03/part_000.gz", 50},
	}

	for _, tt := range tests {
		t.Run(tt.data, func(t *testing.T) {
			results, err := ParseFile(tt.path)
			log.Println("Amount of results:", len(results))
			log.Println("Example:", results[0]["title"].(string))

			ans := len(results)
			if err != nil {
				log.Fatal(err)
			}
			if ans != len(results) {
				t.Errorf("got %d, want %d", ans, tt.want_length)
			}
		})
	}
}
