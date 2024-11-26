package openalex

import (
	"fmt"
	"os"
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

func TestSync(t *testing.T) {
	path := os.Getenv("TEST_OPENALEX_PATH")
	//Please change to your directory
	err := Sync(path, nil)
	if err != nil {
		t.Error(err)
	}
}
