package openalex

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// e.g. C:/openalex/data/
type SnapshotSQL struct {
	gorm.Model
	SnapshotId   string `gorm:"unique"`
	ZipPath      string //Identifier
	DatabasePath string
	Done         bool
	Info         string
	EntityFiles  []EntityFileSQL `gorm:"foreignkey:SnapshotId"`
}

type EntityFileSQL struct {
	gorm.Model
	SnapshotId     uint `gorm:"index"` //foreign key of snapshot
	EntityFileName string
	Identifier     string //ZipPath + "::" + EntityFolderName + "::" + DateFolderName + "::" + EntityZipName
	FullPath       string
	Done           bool
	Info           string
	EntityLines    []EntityLineSQL `gorm:"foreignkey:EntityFileId"`
}

type EntityLineSQL struct {
	gorm.Model
	EntityFileId uint `gorm:"index"` // Foreign key to EntityFile
	LineInfo     string
	Identifier   string //ZipPath + "::" + EntityFolderName + "::" + DateFolderName + "::" + EntityZipName + "::" + LineInfo
	Info         string
	FullPath     string
	Done         bool
}

// Initialize loads Last Known State
// Creates DB if there is none
// returns false if the processing is already finished
// returns true if there is some processing left to be done
func (sh *StateHandler) Initialize() {
	db, err := gorm.Open(sqlite.Open(sh.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed to open " + sh.DatabasePath)
	}

	sh.db = db
	// this will create the tables in the database, or migrate them if they already exist
	err = sh.db.AutoMigrate(&SnapshotSQL{}, &EntityFileSQL{}, &EntityLineSQL{})
	if err != nil {
		slog.With("err", err).Error("could not migrate")
		return
	}

	// Load the Process Dir Struct if it doesn't exist
	var currentSnapshotSQL SnapshotSQL

	dirResult := sh.db.Where("zip_path = ?", sh.SnapshotZipPath).First(&currentSnapshotSQL)
	if dirResult.Error != nil {
		if errors.Is(dirResult.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this processing path, creating one")
			snapshotSQL := SnapshotSQL{
				ZipPath:      sh.SnapshotZipPath,
				DatabasePath: sh.DatabasePath,
				Done:         false,
				Info:         "New Snapshot Process Started",
			}
			sh.db.Create(&snapshotSQL)
			sh.currentSnapshotSQL = snapshotSQL
			return
		} else {
			panic(dirResult.Error)
		}
	}

	// loaded successfully, so cache it in the SqlLogger struct
	sh.currentSnapshotSQL = currentSnapshotSQL
}

// GetDirectoryProcessStatus no if directory.done = false or no entry exists
func (sh *StateHandler) IsSnapshotFinished() bool {
	return sh.currentSnapshotSQL.Done
}

// RegisterOrEntityFile returns True if the entity file is already processed
// If the EntityFile entry does not exist,
// creates a new one (using the current processDir as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (sh *StateHandler) RegisterOrSkipEntityFile(filePath string) (bool, error) {

	identifier := filePath
	fileName := filepath.Base(filePath)

	var entityFile EntityFileSQL
	errEntityFile := sh.db.Where("identifier = ?", identifier).First(&entityFile).Error
	if errEntityFile != nil {
		if errors.Is(errEntityFile, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newEntityFile := EntityFileSQL{
				SnapshotId:     sh.currentSnapshotSQL.ID,
				Identifier:     identifier,
				EntityFileName: fileName,
				FullPath:       filePath,
				Done:           false,
				Info:           "new entity file process started",
			}
			errCreate := sh.db.Create(&newEntityFile).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create entity zip entry in db")
			}
			sh.currentEntityFileSQL = newEntityFile
			return false, nil
		} else {
			return false, errEntityFile
		}
	}

	// won't be registered if already done
	if !entityFile.Done {
		sh.currentEntityFileSQL = entityFile
	}

	return entityFile.Done, nil
}

func (sh *StateHandler) RegisterOrSkipEntityLine(line_info string) (bool, error) {
	identifier := sh.currentEntityFileSQL.Identifier + "::" + line_info

	var entityLine EntityLineSQL
	errEntityLine := sh.db.Where("identifier = ?", identifier).First(&entityLine).Error
	if errEntityLine != nil {
		if errors.Is(errEntityLine, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newEntityLine := EntityLineSQL{
				EntityFileId: sh.currentEntityFileSQL.ID,
				Identifier:   identifier,
				LineInfo:     line_info,
				Done:         false,
				FullPath:     sh.currentEntityFileSQL.FullPath + "::" + line_info,
			}
			errCreate := sh.db.Create(&newEntityLine).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create entity zip entry in db")
			}
			sh.currentEntityLineSQL = newEntityLine
			return false, nil //new processing project
		} else {
			return false, errEntityLine
		}
	}

	// won't be registered if already done
	if !entityLine.Done {
		sh.currentEntityLineSQL = entityLine
	}

	return sh.currentEntityLineSQL.Done, nil
}

// SetSafeDelete no if directory.done = false or no entry exists
func (sh *StateHandler) SetSafeDelete(status bool) {
	sh.SafeDeleteOnly = status
}

func (sh *StateHandler) MarkSnapshotAsFinished() {
	// set the status of the directory as finished
	err := sh.db.Model(&sh.currentSnapshotSQL).Update("done", true).Error
	if err != nil {
		panic(err)
	}
	// set the info of the directory as finished
	err = sh.db.Model(&sh.currentSnapshotSQL).Update("info", "Finished").Error
	if err != nil {
		panic(err)
	}
}

func (sh *StateHandler) MarkSnapshotAsUpdated() {
	// set the status of the directory as finished
	err := sh.db.Model(&sh.currentSnapshotSQL).Update("done", false).Error
	if err != nil {
		panic(err)
	}
	// set the info of the directory as finished
	err = sh.db.Model(&sh.currentSnapshotSQL).Update("info", "Updated").Error
	if err != nil {
		panic(err)
	}
}

func (sh *StateHandler) MarkEntityFileAsFinished() {
	// set the status of the directory as finished
	err := sh.db.Model(&sh.currentEntityFileSQL).Update("done", true).Error
	if err != nil {
		panic(err)
	}
	// set the info of the directory as finished
	err = sh.db.Model(&sh.currentEntityFileSQL).Update("info", "Finished").Error
	if err != nil {
		panic(err)
	}
}

func (sh *StateHandler) MarkEntityLineAsFinished() {
	// set the status of the directory as finished
	err := sh.db.Model(&sh.currentEntityLineSQL).Update("done", true).Error
	if err != nil {
		panic(err)
	}
	// set the info of the directory as finished
	err = sh.db.Model(&sh.currentEntityLineSQL).Update("info", "Finished").Error
	if err != nil {
		panic(err)
	}
}
