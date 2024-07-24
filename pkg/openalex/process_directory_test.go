package openalex

import (
	"strings"
	"testing"
)

// Tests the complete directory
func TestReadFromDirectory(t *testing.T) {
	err := ProcessDirectory("C:\\docdb\\openalex", PrintEntityHandler, PrintMergedIdRecordHandler)
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
