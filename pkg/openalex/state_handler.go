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
	SnapshotId    string `gorm:"unique"`
	ZipPath       string //Identifier
	DatabasePath  string
	Done          bool
	Info          string
	EntityFolders []EntityFolderSQL `gorm:"foreignkey:SnapshotId"`
}

// From here on everything is a subdirectory of the SnapshotDirectory
// /works, /authors, etc...
type EntityFolderSQL struct {
	gorm.Model
	SnapshotId       uint `gorm:"index"` //foreign key fo snapshot
	EntityFolderName string
	Identifier       string //ZipPath + "::" + EntityFolderName
	EntityType       FileEntityType
	Done             bool
	Info             string
	FullPath         string
	DateFolders      []DateFolderSQL `gorm:"foreignkey:EntityFolderId"`
}

// /works/updated_date=2024-06-30, /authors/updated_date=2024-06-30 etc.
type DateFolderSQL struct {
	gorm.Model
	EntityFolderId uint `gorm:"index"` //foreign key of entityfolder
	DateFolderName string
	Identifier     string //ZipPath + "::" + EntityFolderName + "::" + DateFolderName
	FullPath       string
	Info           string
	Done           bool
	EntityFiles    []EntityZipSQL `gorm:"foreignkey:DateFolderId"`
}

// /works/updated_date=2024-06-30/part_043.gz
type EntityZipSQL struct {
	gorm.Model
	DateFolderId  uint `gorm:"index"` // Foreign key to DateFolderSQL
	EntityZipName string
	Identifier    string //ZipPath + "::" + EntityFolderName + "::" + DateFolderName + "::" + EntityZipName
	FullPath      string
	Done          bool
	Info          string
	EntityLines   []EntityLineSQL `gorm:"foreignkey:EntityZipId"`
}

type EntityLineSQL struct {
	gorm.Model
	EntityZipId uint `gorm:"index"` // Foreign key to EntityZip
	LineInfo    string
	Identifier  string //ZipPath + "::" + EntityFolderName + "::" + DateFolderName + "::" + EntityZipName + "::" + LineInfo
	Info        string
	FullPath    string
	Done        bool
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
	err = sh.db.AutoMigrate(&SnapshotSQL{}, &EntityFolderSQL{}, &DateFolderSQL{}, &EntityZipSQL{}, &EntityLineSQL{})
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

// RegisterOrEntityFolder returns True if the entity folder is already processed
// If the EntityFolder entry does not exist,
// creates a new one (using the current processDir as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (sh *StateHandler) RegisterOrSkipEntityFolder(folderName string) (bool, error) {
	identifier := sh.currentSnapshotSQL.ZipPath + "::" + folderName

	var entityFolder EntityFolderSQL
	errEntityFolder := sh.db.Where("identifier = ?", identifier).First(&entityFolder).Error
	if errEntityFolder != nil { //if entry does not exist
		if errors.Is(errEntityFolder, gorm.ErrRecordNotFound) {
			slog.With("fileName", folderName).Info("No record for this entity folder, creating")

			entityType, errEntity := GetEntityType(folderName)
			if errEntity != nil {
				slog.With("err", errEntity).Error("failed to create zip file")
				return false, errEntity
			}

			newEntityFolder := EntityFolderSQL{
				SnapshotId:       sh.currentSnapshotSQL.ID,
				Identifier:       identifier,
				EntityFolderName: folderName,
				EntityType:       entityType,
				Done:             false,
				Info:             "New Entity Folder Process Started",
				FullPath:         filepath.Join(sh.SnapshotZipPath, folderName),
			}
			errCreate := sh.db.Create(&newEntityFolder).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create entity folder entry")
			}
			sh.currentEntityFolderSQL = newEntityFolder
			return false, nil //new processing project
		} else {
			return false, errEntityFolder
		}
	}

	//won't be registered if already done (skip)
	if !entityFolder.Done {
		sh.currentEntityFolderSQL = entityFolder
	}

	return sh.currentEntityFolderSQL.Done, nil
}

func (sh *StateHandler) RegisterOrSkipDateFolder(dateFolderName string) (bool, error) {
	identifier := sh.currentEntityFolderSQL.Identifier + "::" + dateFolderName

	var dateFolder DateFolderSQL
	errDateFolder := sh.db.Where("identifier = ?", identifier).First(&dateFolder).Error
	if errDateFolder != nil {
		if errors.Is(errDateFolder, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newDateFolderSQL := DateFolderSQL{
				EntityFolderId: sh.currentEntityFolderSQL.ID,
				Identifier:     identifier,
				DateFolderName: dateFolderName,
				FullPath:       filepath.Join(sh.currentEntityFolderSQL.FullPath, dateFolderName),
				Info:           "New Date Folder Process Started",
				Done:           false,
			}
			errCreate := sh.db.Create(&newDateFolderSQL).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create xml file")
			}
			sh.currentDateFolderSQL = newDateFolderSQL
			return false, nil
		} else {
			return false, errDateFolder
		}
	}

	// won't be registered if already done
	if !dateFolder.Done {
		sh.currentDateFolderSQL = dateFolder
	}

	return dateFolder.Done, nil
}

func (sh *StateHandler) RegisterOrSkipEntityZip(zipName string) (bool, error) {
	identifier := sh.currentDateFolderSQL.Identifier + "::" + zipName

	var entityZip EntityZipSQL
	errEntityZip := sh.db.Where("identifier = ?", identifier).First(&entityZip).Error
	if errEntityZip != nil {
		if errors.Is(errEntityZip, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newEntityZip := EntityZipSQL{
				DateFolderId:  sh.currentDateFolderSQL.ID,
				Identifier:    identifier,
				EntityZipName: zipName,
				FullPath:      filepath.Join(sh.currentDateFolderSQL.FullPath, zipName),
				Done:          false,
				Info:          "new entity zip process started",
			}
			errCreate := sh.db.Create(&newEntityZip).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create entity zip entry in db")
			}
			sh.currentEntityZipSQL = newEntityZip
			return false, nil
		} else {
			return false, errEntityZip
		}
	}

	// won't be registered if already done
	if !entityZip.Done {
		sh.currentEntityZipSQL = entityZip
	}

	return entityZip.Done, nil
}

func (sh *StateHandler) RegisterOrSkipEntityLine(line_info string) (bool, error) {
	identifier := sh.currentEntityZipSQL.Identifier + "::" + line_info

	var entityLine EntityLineSQL
	errEntityLine := sh.db.Where("identifier = ?", identifier).First(&entityLine).Error
	if errEntityLine != nil {
		if errors.Is(errEntityLine, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newEntityLine := EntityLineSQL{
				EntityZipId: sh.currentEntityZipSQL.ID,
				Identifier:  identifier,
				LineInfo:    line_info,
				Done:        false,
				FullPath:    sh.currentEntityZipSQL.FullPath + "::" + line_info,
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

func (sh *StateHandler) RegisterOrSkipEntityFile(filePath string) (bool, error) {
	logger := slog.With("filePath", filePath)

	entityType, err := GetEntityType(filePath)
	if err != nil {
		logger.With("err", err).Error("error getting entity type")
	}

	entityFolderDone, err := sh.RegisterOrSkipEntityFolder(string(entityType))
	if err != nil {
		logger.With("err", err).Error("error by registering EntityFolder")
	}

	date := getUpdatedDate(filePath)
	dateFolderDone, err := sh.RegisterOrSkipDateFolder(date)
	if err != nil {
		logger.With("err", err).Error("error by registering DateFolder")
	}

	lastElement := filepath.Base(filePath)
	entityZipDone, err := sh.RegisterOrSkipEntityZip(lastElement)
	if err != nil {
		logger.With("err", err).Error("error by registering zip")
	}

	return entityFolderDone && dateFolderDone && entityZipDone, nil
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

func (sh *StateHandler) MarkEntityFolderAsFinished() {
	// set the status of the directory as finished
	err := sh.db.Model(&sh.currentEntityFolderSQL).Update("done", true).Error
	if err != nil {
		panic(err)
	}
	// set the info of the directory as finished
	err = sh.db.Model(&sh.currentEntityFolderSQL).Update("info", "Finished").Error
	if err != nil {
		panic(err)
	}
}

func (sh *StateHandler) MarkDateFolderAsFinished() {
	//Delete EntityZip Files belonging to the current DateFolder File
	var resultDateFolderDelete *gorm.DB

	if sh.SafeDeleteOnly {
		resultDateFolderDelete = sh.db.Where("date_folder_id = ?", sh.currentDateFolderSQL.ID).Delete(&EntityZipSQL{})
	} else {
		resultDateFolderDelete = sh.db.Unscoped().Where("date_folder_id = ?", sh.currentDateFolderSQL.ID).Delete(&EntityZipSQL{})
	}

	if resultDateFolderDelete.Error != nil {
		panic(resultDateFolderDelete.Error)
	}

	// Check deleted records
	fmt.Printf("Deleted %v record(s) with ID '%v'", resultDateFolderDelete.RowsAffected, sh.currentDateFolderSQL.ID)
	// set current Entity Zip to empty
	sh.currentEntityZipSQL = EntityZipSQL{}

	//We're keeping the Date Folder Entry, only delete downwards
	resultInfo := sh.db.Model(&sh.currentDateFolderSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}

	resultStatus := sh.db.Model(&sh.currentDateFolderSQL).Update("done", true)
	if resultStatus.Error != nil {
		panic(resultStatus.Error)
	}
}

func (sh *StateHandler) MarkEntityZipAsFinished() {
	//Delete EntityLines belonging to the current EntityZip File
	var resultEntityZipDelete *gorm.DB

	if sh.SafeDeleteOnly {
		resultEntityZipDelete = sh.db.Where("entity_zip_id = ?", sh.currentEntityZipSQL.ID).Delete(&EntityLineSQL{})
	} else {
		resultEntityZipDelete = sh.db.Unscoped().Where("entity_zip_id = ?", sh.currentEntityZipSQL.ID).Delete(&EntityLineSQL{})
	}

	if resultEntityZipDelete.Error != nil {
		panic(resultEntityZipDelete.Error)
	}

	// Check deleted records
	fmt.Printf("Deleted %v record(s) with ID '%v'", resultEntityZipDelete.RowsAffected, sh.currentEntityZipSQL.ID)
	// set current EntityLine to empty
	sh.currentEntityLineSQL = EntityLineSQL{}

	//We're keeping the Date Folder Entry, only delete downwards
	resultInfo := sh.db.Model(&sh.currentEntityZipSQL).Update("info", "finished")
	if resultInfo.Error != nil {
		panic(resultInfo.Error)
	}

	resultStatus := sh.db.Model(&sh.currentEntityZipSQL).Update("done", true)
	if resultStatus.Error != nil {
		panic(resultStatus.Error)
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
