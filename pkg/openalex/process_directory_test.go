package openalex

import (
	"testing"
)

// Tests the complete directory
func TestReadFromDirectory(t *testing.T) {
	err := ProcessDirectory("../../sample/openalex/", PrintEntityHandler, PrintMergedIdRecordHandler)
	if err != nil {
		t.Error(err)
	}
}
