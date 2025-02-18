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
	TopicsFileEntityType       FileEntityType = "topics"
	DomainsFileEntityType      FileEntityType = "domains"
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
	} else if strings.Contains(filePath, "topics") {
		result = TopicsFileEntityType
	} else if strings.Contains(filePath, "domains") {
		result = DomainsFileEntityType
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

// LineHandler is a function that handles a line of a file
type LineHandler func(filePath string, line string) error

// PrintLineHandler is a function that prints a parsed line of a file
func PrintLineHandler(filePath string, line string) error {
	fmt.Println(filePath, line)
	return nil
}

// ParseFile takes a file name and reads the data from within the file and parses every line it into structs
func (p *Processor) ParseFile(filePath string) (count int, err error) {
	logger := slog.With("filePath", filePath)
	count = 0

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
	const maxCapacity = 500 * 1024 * 1024 // 500 MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// iterate over the lines
	entityLineIndex := 0
	for scanner.Scan() {
		entityLineIndex++
		if p.StateHandler != nil {
			entityLineName := "entity_line_" + strconv.Itoa(entityLineIndex) + "_end"
			entityLineDone, _ := p.StateHandler.RegisterOrSkipEntityLine(entityLineName)
			if entityLineDone {
				continue
			}
		}
		line := scanner.Text()
		// handle the parsed line
		err = p.LineHandler(filePath, line)
		if err != nil {
			logger.With("err", err).Error("error handling parsed entity")
			return count, err
		}
		if p.StateHandler != nil {
			p.StateHandler.MarkEntityLineAsFinished()
		}
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
