# Open Alex
This package is interacting with the Open Alex API.
It downloads the data into strongly typed structs.

# Status
Work in progress

## Requirements

* AWS CLI for syncing the data from S3

## Install

```
go get -u github.com/max-planck-innovation-competition/go-openalex
```

## Usage

### Downloading the data

```go
dirPath := "./path/to/folder"
openalex.Sync(dirPath)
```

### Process the directory

```go
// write yor own handlers
err = openalex.ProcessDirectory(dirPath, openalex.PrintEntityHandler, openalex.PrintMergedIdRecordHandler)
if err != nil {
    panic(err)
}
```

### Handlers

#### EntityHandler
Every line that is parsed from the data is passed to the EntityHandler.
The EntityHandler is called with the FileEntityType and the entity.
You can pass your own handler to upload the data to a database.
```go
func EntityHandler(fileEntityType FileEntityType, entity any) error {
	// TODO
} 
```

#### MergedIdRecordHandler
The MergedIdRecordHandler is called with the FileEntityType and the mergedIdRecord.
You can pass your own handler to modify the data in your database.
```go
func MergedIdRecordHandler(fileEntityType FileEntityType, mergedIdRecord any) error {
    // TODO
} 
```