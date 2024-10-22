package openalex

import (
	"fmt"
	"testing"
)

// WARNING
// Note that the Snapshot has around 422GB and 1.6TB after uncompression
func TestDownloadSnapshot(t *testing.T) {

	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	//Please change to your directory
	err := Sync("C:\\docdb\\openalex", stateHandler)
	if err != nil {
		t.Error(err)
	}
}
