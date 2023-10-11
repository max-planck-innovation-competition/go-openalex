package openalex

import (
	"compress/gzip"
	"encoding/csv"
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

// ParseMergedIDsFile parses a CSV file (either plain or Gzipped) into a slice of CsvData
func ParseMergedIDsFile(filePath string) (results []MergedID, err error) {
	logger := slog.With("filePath", filePath)
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
			return nil, err
		}
		if len(record) != 3 {
			logger.With("record", strings.Join(record, ";")).Warn("Invalid CSV record")
			continue
		}
		if rowCount == 0 {
			rowCount++
			continue
		}
		// append the parsed record to the results
		results = append(results, MergedID{
			MergeDate:   record[0],
			ID:          record[1],
			MergeIntoID: record[2],
		})
		rowCount++
	}

	return
}
