package openalex

import "log/slog"

// ParseEntity is a sample implementation of a function that parses a JSON line into a struct
func ParseEntity(filePath string, line string) (data *Work, err error) {
	logger := slog.With("filePath", filePath)

	// determine the struct type based on the filePath
	entityType, err := GetEntityType(filePath)
	if err != nil {
		logger.With("err", err).Error("error getting entity type")
		return nil, err
	}

	//switch entityType {
	//case AuthorsFileEntityType:
	//	data = &Author{}
	//case ConceptsFileEntityType:
	//	data = &Concept{}
	//case FundersFileEntityType:
	//	data = &Funder{}
	//case InstitutionsFileEntityType:
	//	data = &Institution{}
	//case SourcesFileEntityType:
	//	data = &Source{}
	//case PublishersFileEntityType:
	//	data = &Publisher{}
	//case WorksFileEntityType:
	//	data = &Work{}
	//case TopicsFileEntityType:
	//	data = &Topic{}
	//case DomainsFileEntityType:
	//	data = &Domain{}
	//}

	// convert the inverted abstract
	if entityType == WorksFileEntityType {
		// determine the struct type based on the filePath
		var data Work

		// Unmarshal the JSON line into the determined struct using jsoniter
		err = json.UnmarshalFromString(line, &data)
		if err != nil {
			logger.With("err", err).Error("error unmarshalling line")
			return nil, err
		}

		data.Abstract = data.ToAbstract()

		return &data, err
	}

	// todo: add MergeIdEntityType

	return nil, err
}
