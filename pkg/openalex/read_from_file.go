package openalex

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log/slog"
	"os"
	"path"
	"strings"
)

// use faster parser
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ErrUnsupportedFileType is returned when the file type is not supported
var ErrUnsupportedFileType = errors.New("unsupported file type")

type FileEntityType string

const (
	AuthorsFileEntityType      FileEntityType = "authors"
	ConceptsFileEntityType     FileEntityType = "concepts"
	FundersFileEntityType      FileEntityType = "funders"
	InstitutionsFileEntityType FileEntityType = "institution"
	PublishersFileEntityType   FileEntityType = "publisher"
	SourcesFileEntityType      FileEntityType = "sources"
	WorksFileEntityType        FileEntityType = "works"
)

func getEntityType(filePath string) (result FileEntityType, err error) {
	if strings.Contains(filePath, "author") {
		result = AuthorsFileEntityType
	} else if strings.Contains(filePath, "concepts") {
		result = ConceptsFileEntityType
	} else if strings.Contains(filePath, "funders") {
		result = FundersFileEntityType
	} else if strings.Contains(filePath, "institutions") {
		result = InstitutionsFileEntityType
	} else if strings.Contains(filePath, "publisher") {
		result = PublishersFileEntityType
	} else if strings.Contains(filePath, "works") {
		result = WorksFileEntityType
	} else if strings.Contains(filePath, "sources") {
		result = SourcesFileEntityType
	} else {
		// handle unsupported filePath or struct type
		slog.Error("unsupported filePath")
		err = ErrUnsupportedFileType
		return
	}
	return
}

// ParsedEntityLineHandler is a function that handles a parsed line of a file
type ParsedEntityLineHandler func(fileEntityType FileEntityType, entity any) error

// PrintEntityHandler is a function that prints a parsed line of a file
func PrintEntityHandler(fileEntityType FileEntityType, entity any) error {
	fmt.Println(fileEntityType, entity)
	return nil
}

// ParseFile takes a file name and reads the data from within the file and parses every line it into structs
func ParseFile(filePath string, fn ParsedEntityLineHandler) (err error) {
	logger := slog.With("filePath", filePath)

	// determine the struct type based on the filePath
	entityType, err := getEntityType(filePath)
	if err != nil {
		logger.With("err", err).Error("error getting entity type")
		return err
	}

	// init the read
	var scanner *bufio.Scanner

	// check if rawContent is compressed
	fileExtension := path.Ext(filePath)
	if fileExtension == ".gz" {
		// if file has a .gz ending
		compressedFile, errOpen := os.Open(filePath)
		if errOpen != nil {
			slog.With("err", errOpen).Error("error opening file")
			return errOpen
		}
		defer compressedFile.Close()
		// get the raw content of the file
		rawContent, errGzip := gzip.NewReader(compressedFile)
		if errGzip != nil {
			slog.With("err", errGzip).Error("error opening gz file")
			return errGzip
		}
		scanner = bufio.NewScanner(rawContent)
	} else {
		// if file has no file ending that indicates compression
		fileContent, errOpen := os.Open(filePath)
		if errOpen != nil {
			slog.With("err", errOpen).Error("error opening file")
			return errOpen
		}
		defer fileContent.Close()
		// init scanner
		scanner = bufio.NewScanner(fileContent)
	}

	// iterate over the lines
	for scanner.Scan() {
		line := scanner.Text()
		// determine the struct type based on the filePath
		var data interface{}

		switch entityType {
		case AuthorsFileEntityType:
			data = &Author{}
		case ConceptsFileEntityType:
			data = &Concept{}
		case FundersFileEntityType:
			data = &Funder{}
		case InstitutionsFileEntityType:
			data = &Institution{}
		case SourcesFileEntityType:
			data = &Source{}
		case PublishersFileEntityType:
			data = &Publisher{}
		case WorksFileEntityType:
			data = &Work{}

		}

		// Unmarshal the JSON line into the determined struct using jsoniter
		err = json.UnmarshalFromString(line, data)
		if err != nil {
			logger.With("err", err).Error("error unmarshalling line")
			return err
		}

		// handle the parsed line
		err = fn(entityType, data)
	}

	err = scanner.Err()
	if err != nil {
		logger.With("err", err).Error("error scanning file")
		return err
	}

	return nil
}
