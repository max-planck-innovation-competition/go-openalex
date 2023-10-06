package openalex

import (
	"log"
	"testing"
)

// Tests the complete directory
func TestReadFromDirectory(t *testing.T) {
	results, err := ReadFromDirectory("../../sample/openalex/")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Amount of results", len(results))
	log.Println("Example:", results[0])
}
