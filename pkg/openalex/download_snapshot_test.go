package openalex

import (
	"testing"
)

// WARNING
// Note that the Snapshot has around 422GB and 1.6TB after uncompression
func TestDownloadSnapshot(t *testing.T) {
	//Please change to your directory
	err := Sync("C:\\docdb\\openalex")
	if err != nil {
		t.Error(err)
	}
}
