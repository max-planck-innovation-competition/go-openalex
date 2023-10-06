package main

import "github.com/max-planck-innovation-competition/go-openalex/pkg/openalex"

func main() {

	// start download and update the data
	err := openalex.Sync("openalex-snapshot")
	if err != nil {
		return
	}

	results, _ := openalex.ReadFromDirectory("openalex-snapshot/data/works")
	_ = results
}
