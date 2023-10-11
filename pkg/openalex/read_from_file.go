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

// ParseFile takes a file name and reads the data from within the file and parses every line it into structs
func ParseFile(filePath string) (err error) {
	logger := slog.With("filePath", filePath)

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
		if strings.Contains(filePath, "work") {
			data = &Work{}
		} else if strings.Contains(filePath, "author") {
			data = &Author{}
		} else {
			// handle unsupported filePath or struct type
			logger.Error("unsupported filePath")
			return ErrUnsupportedFileType
		}

		// Unmarshal the JSON line into the determined struct using jsoniter
		err = json.UnmarshalFromString(line, data)
		if err != nil {
			logger.With("err", err).Error("error unmarshalling line")
			return err
		}

		// You can now work with the parsed struct data as needed.
		switch data.(type) {
		case *Work:
			work := data.(*Work)
			fmt.Println(work.Title)
			fmt.Printf("Parsed StructA: %+v\n", work.AbstractInvertedIndex.ToAbstract())
		case *Author:
			structBData := data.(*Author)
			fmt.Printf("Parsed StructB: %+v\n", structBData)
		default:
			logger.Error("unsupported struct type")
			return ErrUnsupportedFileType
		}
	}

	err = scanner.Err()
	if err != nil {
		logger.With("err", err).Error("error scanning file")
		return err
	}

	return nil
}
