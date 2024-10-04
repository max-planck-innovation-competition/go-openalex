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
	if stateHandler.IsSnapshotFinished() {
		return
	}

	entities := [7]string{
		"authors",
		"concepts",
		"funders",
		"institution",
		"publisher",
		"sources",
		"works",
	}

	for _, entity := range entities {

		entityFolderDone, _ := stateHandler.RegisterOrSkipEntityFolder(entity)
		if entityFolderDone {
			fmt.Println("skipped")
			continue
		}

		for dateFolderIndex := 1; dateFolderIndex < 5; dateFolderIndex++ {
			dateFolderName := "updated_date=2024-06-0" + strconv.Itoa(dateFolderIndex)
			dateFolderDone, _ := stateHandler.RegisterOrSkipDateFolder(dateFolderName)
			if dateFolderDone {
				continue
			}

			for entityZipIndex := 1; entityZipIndex < 10; entityZipIndex++ {
				entityZipName := "part_0" + strconv.Itoa(entityZipIndex) + ".gz"
				entityZipDone, _ := stateHandler.RegisterOrSkipEntityZip(entityZipName)
				if entityZipDone {
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

				stateHandler.MarkEntityZipAsFinished()
			}

			stateHandler.MarkDateFolderAsFinished()
		}

		stateHandler.MarkEntityFolderAsFinished()
	}

	stateHandler.MarkSnapshotAsFinished()
}
