package openalex

import (
	"fmt"
	"strconv"
	"testing"
)

// to test: go test -timeout 99999s -run TestSQLHandlerFull -v
func TestSQLHandlerFull(t *testing.T) {
	// TODO: windows directory as env variable
	stateHandler := NewStateHandler("log.db", "C:\\go-openalex\\openalex-snapshot\\data", "C:\\go-openalex\\openalex-snapshot\\data")
	//Comment to test safe delete
	stateHandler.SetSafeDelete(false) //deletes everything under the Date Folder after finishing
	fmt.Println(stateHandler)

	stateHandler.MarkSnapshotAsUpdated()

	if stateHandler.IsSnapshotFinished() {
		return
	}

	filepaths := [4]string{
		"openalex-snapshot/data/works/updated_date=2023-09-20/part_000.gz",
		"openalex-snapshot/data/works/updated_date=2023-09-20/part_000.gz",
		"openalex-snapshot/data/works/updated_date=2023-09-20/part_001.gz",
		"openalex-snapshot/data/works/updated_date=2023-09-20/part_002.gz",
	}

	for _, filepath := range filepaths {

		entityFileDone, _ := stateHandler.RegisterOrSkipEntityFile(filepath)
		if entityFileDone {
			fmt.Println("skipped")
			continue
		}

		for entityLineIndex := 1; entityLineIndex < 20; entityLineIndex++ {
			entityLineName := "entity_line_" + strconv.Itoa(entityLineIndex) + "_end"
			entityLineDone, _ := stateHandler.RegisterOrSkipEntityLine(entityLineName)
			if entityLineDone {
				continue
			}

			stateHandler.MarkEntityLineAsFinished()
		}

		stateHandler.MarkEntityFileAsFinished()
	}

	stateHandler.MarkSnapshotAsFinished()
}
