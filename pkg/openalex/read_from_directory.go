package openalex

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// visit walks over files in a directory
func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("visit", err)
		}
		// do not include directories
		// only include files with .gz extension
		if !info.IsDir() && strings.Contains(path, ".gz") {
			// only files
			*files = append(*files, path)
		}
		return nil
	}
}

// ReadFromDirectory parses the directory of separated files provided by semantic scholar
func ReadFromDirectory(directoryPath string) (results []map[string]interface{}, err error) {
	log.Println("Start reading directory:", directoryPath)

	var filePaths []string // stores the file paths of all the files in the directory

	// walk over the files in the directory
	err = filepath.Walk(directoryPath, visit(&filePaths))
	if err != nil {
		return nil, err
	}

	// Read all files
	for _, file := range filePaths {
		errFile := ParseFile(file)
		if errFile != nil {
			log.Println("error while reading file: ", file, " : ", errFile)
			return nil, errFile
		}
		// add the parsed documents to the results
		results = append(results)
	}
	return
}
