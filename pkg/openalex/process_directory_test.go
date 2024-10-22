package openalex

import (
	"fmt"
	"strings"
	"testing"
)

// Tests the complete directory
func TestReadFromDirectory(t *testing.T) {

	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	//Change to according directory
	err := ProcessDirectory("C:\\docdb\\openalex", PrintEntityHandler, PrintMergedIdRecordHandler, stateHandler)
	if err != nil {
		t.Error(err)
	}
}

func TestOrderByMergedIDsLast(t *testing.T) {

	filePaths := []string{
		"works/file1.txt",
		"merged_ids/file1.txt",
		"works/file2.txt",
		"authors/updated_date=2023-04-21/part_000",
		"authors/manifest",
		"funders/file3.txt",
		"merged_ids/file2.txt",
		"file4.txt",
	}

	orderedFilePaths := OrderByMergedIDsLast(filePaths)
	for _, path := range orderedFilePaths {
		println(path)
	}

	if orderedFilePaths[0] != "authors/manifest" {
		t.Error("file1.txt should be first", orderedFilePaths[0])
	}

	if !strings.Contains(orderedFilePaths[len(orderedFilePaths)-1], "merged_ids") {
		t.Error("merged_ids should be last", orderedFilePaths[len(orderedFilePaths)-1])
	}

}
