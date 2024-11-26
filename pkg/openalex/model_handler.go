package openalex

import "log/slog"

// ParseEntity is a sample implementation of a function that parses a JSON line into a struct
func ParseEntity(filePath string, line string) (err error) {
	logger := slog.With("filePath", filePath)

	// determine the struct type based on the filePath
	entityType, err := GetEntityType(filePath)
	if err != nil {
		logger.With("err", err).Error("error getting entity type")
		return err
	}

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
	case TopicsFileEntityType:
		data = &Topic{}
	case DomainsFileEntityType:
		data = &Domain{}
	}

	// Unmarshal the JSON line into the determined struct using jsoniter
	err = json.UnmarshalFromString(line, data)
	if err != nil {
		logger.With("err", err).Error("error unmarshalling line")
		return err
	}

	// convert the inverted abstract
	if entityType == WorksFileEntityType {
		work := data.(*Work)
		work.Abstract = work.ToAbstract()
		data = work
	}

	return nil
}
