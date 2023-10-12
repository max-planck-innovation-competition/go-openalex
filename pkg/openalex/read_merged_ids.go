package openalex

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// MergedID represents a row in the merged IDs file
type MergedID struct {
	MergeDate   string
	ID          string
	MergeIntoID string
}

// MergedIdRecordHandler is a function that handles a parsed line of a file
type MergedIdRecordHandler func(fileEntityType FileEntityType, mergedID MergedID) error

func PrintMergedIdRecordHandler(fileEntityType FileEntityType, mergedID MergedID) error {
	fmt.Println("fileEntityType", fileEntityType, "mergedID", mergedID)
	return nil
}

// ParseMergedIDsFile parses a CSV file (either plain or Gzipped) into a slice of CsvData
func ParseMergedIDsFile(filePath string, fn MergedIdRecordHandler) (err error) {
	logger := slog.With("filePath", filePath)

	// get the file entity type
	fileEntityType, err := getEntityType(filePath)
	if err != nil {
		logger.With("err", err).Error("error getting file entity type")
		return
	}

	var reader io.Reader

	// Check if the file is Gzipped
	if strings.HasSuffix(filePath, ".gz") {
		// If it's Gzipped, open a Gzip reader
		file, errOpen := os.Open(filePath)
		if errOpen != nil {
			err = errOpen
			logger.With("err", err).Error("error opening file")
			return
		}
		defer file.Close()

		reader, err = gzip.NewReader(file)
		if err != nil {
			logger.With("err", err).Error("error opening gz file")
			return
		}
	} else {
		// If it's not Gzipped, open a regular file reader
		file, errOpen := os.Open(filePath)
		if err != nil {
			err = errOpen
			logger.With("err", err).Error("error opening file")
			return
		}
		defer file.Close()

		reader = file
	}

	// Create a CSV reader to read the data
	csvReader := csv.NewReader(reader)

	// Read and parse the CSV data into a slice of CsvData structs
	rowCount := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.With("err", err).Error("error reading CSV record")
			return err
		}
		if len(record) != 3 {
			logger.With("record", strings.Join(record, ";")).Warn("Invalid CSV record")
			continue
		}
		if rowCount == 0 {
			rowCount++
			continue
		}
		// generate the merged id struct
		mergedID := MergedID{
			MergeDate:   record[0],
			ID:          record[1],
			MergeIntoID: record[2],
		}
		// process the merged id
		errProcess := fn(fileEntityType, mergedID)
		if errProcess != nil {
			logger.With("err", errProcess).Error("error processing CSV record")
			return errProcess
		}
		rowCount++
	}
	return
}
