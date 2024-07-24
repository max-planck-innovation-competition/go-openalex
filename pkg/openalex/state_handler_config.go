package openalex

import (
	"path/filepath"

	"gorm.io/gorm"
)

// StateHandler contains the config for the state handler
type StateHandler struct {
	//initialize these
	DatabaseName    string //e.g. log.db, for the initializer
	DatabaseDir     string //path of the .db, e.g. C:\docdb\ or .\ for relative path
	SnapshotZipPath string //full path to the snapshot zip
	SafeDeleteOnly  bool
	//for the state
	//these are initialized in NewSqlLogger(...)
	currentSnapshotSQL     SnapshotSQL
	currentEntityFolderSQL EntityFolderSQL
	currentDateFolderSQL   DateFolderSQL
	currentEntityZipSQL    EntityZipSQL
	currentEntityLineSQL   EntityLineSQL
	DatabasePath           string //Database Dir + Database Name
	db                     *gorm.DB
}

// New creates a new state handler
func NewStateHandler(databaseName string, databaseDir string, snapshotZipPath string) *StateHandler {
	stateHandler := StateHandler{
		DatabaseName:    databaseName,
		DatabaseDir:     databaseDir,
		SnapshotZipPath: snapshotZipPath,
		SafeDeleteOnly:  true,
		DatabasePath:    filepath.Join(databaseDir, databaseName),
	}
	stateHandler.Initialize() //Initializes the other fields
	return &stateHandler
}
