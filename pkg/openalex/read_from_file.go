package openalex

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
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

func GetEntityType(filePath string) (result FileEntityType, err error) {
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

func getUpdatedDate(filePath string) string {
	// Define the regex pattern to match "updated_date=YYYY-MM-DD"
	pattern := `updated_date=\d{4}-\d{2}-\d{2}`
	re := regexp.MustCompile(pattern)

	// Find the first occurrence of the pattern in the path
	match := re.FindString(filePath)
	return match
}

// ParsedEntityLineHandler is a function that handles a parsed line of a file
type ParsedEntityLineHandler func(fileEntityType FileEntityType, entity any) error

// PrintEntityHandler is a function that prints a parsed line of a file
func PrintEntityHandler(fileEntityType FileEntityType, entity any) error {
	fmt.Println(fileEntityType, entity)
	return nil
}

// ParseFile takes a file name and reads the data from within the file and parses every line it into structs
func ParseFile(filePath string, fn ParsedEntityLineHandler, sh *StateHandler) (count int, err error) {
	logger := slog.With("filePath", filePath)
	count = 0

	// determine the struct type based on the filePath
	entityType, err := GetEntityType(filePath)
	if err != nil {
		logger.With("err", err).Error("error getting entity type")
		return count, err
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
			return count, errOpen
		}
		defer compressedFile.Close()
		// get the raw content of the file
		rawContent, errGzip := gzip.NewReader(compressedFile)
		if errGzip != nil {
			slog.With("err", errGzip).Error("error opening gz file")
			return count, errGzip
		}
		scanner = bufio.NewScanner(rawContent)
	} else {
		// if file has no file ending that indicates compression
		fileContent, errOpen := os.Open(filePath)
		if errOpen != nil {
			slog.With("err", errOpen).Error("error opening file")
			return count, errOpen
		}
		defer fileContent.Close()
		// init scanner
		scanner = bufio.NewScanner(fileContent)
	}
	// set the max capacity of the scanner
	const maxCapacity = 100 * 1024 * 1024 // 100 MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// iterate over the lines
	entityLineIndex := 0
	for scanner.Scan() {
		entityLineIndex++
		entityLineName := "entity_line_" + strconv.Itoa(entityLineIndex) + "_end"
		entityLineDone, _ := sh.RegisterOrSkipEntityLine(entityLineName)
		if entityLineDone {
			continue
		}

		line := scanner.Text()
		// replace all open alex prefixes
		line = strings.ReplaceAll(line, "https://openalex.org/", "")
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
			return count, err
		}

		// convert the inverted abstract
		if entityType == WorksFileEntityType {
			work := data.(*Work)
			work.Abstract = work.ToAbstract()
			data = work
		}

		// handle the parsed line
		err = fn(entityType, data)

		sh.MarkEntityLineAsFinished()

		// increment the count of the parsed record
		count++

	}

	err = scanner.Err()
	if err != nil {
		logger.With("err", err).Error("error scanning file")
		return count, err
	}

	return count, nil
}
