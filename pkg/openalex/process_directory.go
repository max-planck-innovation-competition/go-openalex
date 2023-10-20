package openalex

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
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

// OrderByMergedIDsLast sorts the file paths so that the merged ids file is last
func OrderByMergedIDsLast(filePaths []string) []string {
	// Custom sorting function to place file paths with "merged_ids" at the end
	sort.SliceStable(filePaths, func(i, j int) bool {
		if containsMergedIDs(filePaths[i]) && containsMergedIDs(filePaths[j]) {
			return false
		} else if containsMergedIDs(filePaths[i]) {
			return false
		} else if containsMergedIDs(filePaths[j]) {
			return true
		} else {
			return filePaths[i] < filePaths[j]
		}
	})
	return filePaths
}

func containsMergedIDs(filePath string) bool {
	return strings.Contains(filePath, "merged_ids")
}

// ProcessDirectory parses the directory of separated files and processes them
func ProcessDirectory(directoryPath string, fnEntityHandler ParsedEntityLineHandler, fnMergedIdHandler MergedIdRecordHandler) (err error) {
	logger := slog.With("directoryPath", directoryPath)
	logger.Info("Start reading directory")
	// get the files
	filePaths, err := GetFiles(directoryPath)
	if err != nil {
		logger.With("err", err).Error("error while reading the directory")
		return
	}
	// process the files
	err = ProcessFiles(filePaths, fnEntityHandler, fnMergedIdHandler)
	if err != nil {
		logger.With("err", err).Error("error while processing the files")
		return
	}
	logger.Info("Finished reading directory")
	return
}

// GetFiles returns a list of files in a directory
func GetFiles(directoryPath string) (filePaths []string, err error) {
	logger := slog.With("directoryPath", directoryPath)
	logger.Info("Start listing directory")
	// walk over the files in the directory
	err = filepath.Walk(directoryPath, visit(&filePaths))
	if err != nil {
		logger.With("err", err).Error("error while walking the directory")
		return
	}
	// order the files
	filePaths = OrderByMergedIDsLast(filePaths)
	logger.Info("Finished listing directory")
	return
}

// ProcessFiles parses the files and processes them
func ProcessFiles(filePaths []string, fnEntityHandler ParsedEntityLineHandler, fnMergedIdHandler MergedIdRecordHandler) (err error) {
	logger := slog.With("method", "ProcessFiles")
	// process all files
	for _, filePath := range filePaths {
		if strings.Contains(filePath, "merged_ids") {
			// handle merged ids file
			errFile := ParseMergedIDsFile(filePath, fnMergedIdHandler)
			if errFile != nil {
				logger.
					With("err", errFile).
					With("filePath", filePath).
					Error("error while parsing the merged id file")
				return errFile
			}
		} else {
			// handle other files
			_, errFile := ParseFile(filePath, fnEntityHandler)
			if errFile != nil {
				logger.
					With("err", errFile).
					With("filePath", filePath).
					Error("error while parsing the file")
				return errFile
			}
		}
	}
	return
}
