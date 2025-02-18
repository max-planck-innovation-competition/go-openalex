package openalex

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// e.g. C:/openalex/data/
type SnapshotSQL struct {
	gorm.Model
	SnapshotId    string `gorm:"unique"`
	ZipPath       string
	DatabasePath  string
	Done          bool
	Info          string
	EntityFolders []EntityFolderSQL `gorm:"foreignkey:FolderID"`
}

// From here on everything is a subdirectory of the SnapshotDirectory
// /works, /authors, etc...
type EntityFolderSQL struct {
	gorm.Model
	SnapshotId       uint `gorm:"index"` //foreign key fo snapshot
	EntityFolderName string
	EntityType       FileEntityType
	Done             bool
	Info             string
	FullPath         string
	DateFolders      []DateFolderSQL `gorm:"foreignkey:DateFolderID"`
}

// /works/updated_date=2024-06-30, /authors/updated_date=2024-06-30 etc.
type DateFolderSQL struct {
	gorm.Model
	EntityFolderId uint `gorm:"index"` //foreign key of entityfolder
	DateFolderName string
	FullPath       string
	Info           string
	Done           bool
	EntityFiles    []EntityZipSQL `gorm:"foreignkey:ZipID"`
}

// /works/updated_date=2024-06-30/part_043.gz
type EntityZipSQL struct {
	gorm.Model
	DateFolderId  uint `gorm:"index"` // Foreign key to DateFolderSQL
	EntityZipName string
	FullPath      string
	Done          bool
	Info          string
	EntityLines   []EntityLineSQL `gorm:"foreignkey:LineID"`
}

type EntityLineSQL struct {
	gorm.Model
	EntityZipId uint `gorm:"index"` // Foreign key to EntityZip
	LineId      string
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

	dirResult := sh.db.Where("processing_dir = ?", sh.SnapshotZipPath).First(&currentSnapshotSQL)
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
			sh.currentSnapshotSQL = currentSnapshotSQL
			return
		} else {
			panic(dirResult.Error)
		}
	}

	// loaded successfully, so cache it in the SqlLogger struct
	sh.currentSnapshotSQL = currentSnapshotSQL
}

// GetDirectoryProcessStatus no if directory.done = false or no entry exists
func (sh *StateHandler) IsSnapshotFinished() (bool, error) {
	return sh.currentSnapshotSQL.Done, nil
}

// RegisterOrEntityFolder returns True if the entity folder is already processed
// If the EntityFolder entry does not exist,
// creates a new one (using the current processDir as foreign key)
// or loads the existing bulk file information if the entry exists but is not done
func (sh *StateHandler) RegisterOrSkipEntityFolder(folderName string) (bool, error) {
	var entityFolder EntityFolderSQL
	errEntityFolder := sh.db.Where("entity_folder_name = ?", folderName).First(&entityFolder).Error
	if errEntityFolder != nil { //if entry does not exist
		if errors.Is(errEntityFolder, gorm.ErrRecordNotFound) {
			slog.With("fileName", folderName).Info("No record for this entity folder, creating")

			entityType, errEntity := getEntityType(folderName)
			if errEntity != nil {
				slog.With("err", errEntity).Error("failed to create zip file")
				return false, errEntity
			}

			newEntityFolder := EntityFolderSQL{
				SnapshotId:       sh.currentSnapshotSQL.ID,
				EntityFolderName: folderName,
				EntityType:       entityType,
				Done:             false,
				Info:             "New Entity Folder Process Started",
				FullPath:         filepath.Join(sh.SnapshotZipPath, "::", folderName),
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
	var dateFolder DateFolderSQL
	errDateFolder := sh.db.Where("date_folder_name = ?", dateFolderName).First(&dateFolder).Error
	if errDateFolder.Error != nil {
		if errors.Is(errDateFolder, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newDateFolderSQL := DateFolderSQL{
				EntityFolderId: sh.currentEntityFolderSQL.ID,
				DateFolderName: dateFolderName,
				FullPath:       sh.currentEntityFolderSQL.FullPath + "/" + dateFolderName,
				Info:           "New Date Folder Process Started",
				Done:           false,
			}
			errCreate := sh.db.Create(&newDateFolderSQL).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create xml file")
			}
			sh.currentDateFolderSQL = dateFolder
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
	var entityZip EntityZipSQL
	errEntityZip := sh.db.Where("entity_zip_name = ?", zipName).First(&entityZip).Error
	if errEntityZip != nil {
		if errors.Is(errEntityZip, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newEntityZip := EntityZipSQL{
				DateFolderId:  sh.currentDateFolderSQL.ID,
				EntityZipName: zipName,
				FullPath:      sh.currentDateFolderSQL.FullPath + "/" + zipName,
				Done:          false,
				Info:          "new entity zip process started",
			}
			errCreate := sh.db.Create(&newEntityZip).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create entity zip entry in db")
			}
			sh.currentEntityZipSQL = entityZip
			return false, nil //new processing project
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

func (sh *StateHandler) RegisterOrSkipEntityLine(line_id string) (bool, error) {
	var entityLine EntityLineSQL
	errEntityLine := sh.db.Where("line_id = ?", line_id).First(&entityLine).Error
	if errEntityLine != nil {
		if errors.Is(errEntityLine, gorm.ErrRecordNotFound) {
			fmt.Println("No record for this xml file, creating")
			newEntityZip := EntityLineSQL{
				EntityZipId: sh.currentEntityZipSQL.ID,
				LineId:      line_id,
				Done:        false,
				FullPath:    sh.currentEntityZipSQL.FullPath + "::" + line_id,
			}
			errCreate := sh.db.Create(&newEntityZip).Error
			if errCreate != nil {
				slog.With("err", errCreate).Error("failed to create entity zip entry in db")
			}
			sh.currentEntityLineSQL = entityLine
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
	// set current Exchange File to empty
	sh.currentDateFolderSQL = DateFolderSQL{}

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
	// set current Exchange File to empty
	sh.currentEntityZipSQL = EntityZipSQL{}

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
