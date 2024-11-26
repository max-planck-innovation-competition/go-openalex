package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/SbstnErhrdt/env"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/max-planck-innovation-competition/go-openalex/pkg/openalex"
	"log/slog"
	"strings"
	"time"
)

var esClient *elasticsearch.Client

func init() {
	env.LoadEnvFiles()
	c, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{
				"http://localhost:9200",
			},
			Username: "elastic",
			Password: env.FallbackEnvVariable("ELASTICSEARCH_PASSWORD", "TF2vEJMsWzn82"),
		},
	)
	if err != nil {
		panic(err)
	}
	// ping the client
	_, err = c.Info()
	if err != nil {
		panic(err)
	}

	esClient = c
}

var indexers = map[openalex.FileEntityType]esutil.BulkIndexer{}

func createBulkIndexer(entityType openalex.FileEntityType) {
	bulkIndex, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        esClient,
		Index:         string(entityType),
		NumWorkers:    4,
		FlushBytes:    1e+6, // 5MB
		FlushInterval: 10 * time.Second,
		OnError: func(ctx context.Context, err error) {
			slog.With("err", err).With("bi", entityType).Error("error in bulk indexer")
		},
	})
	if err != nil {
		panic(err)
	}
	indexers[entityType] = bulkIndex
}

/*
"mappings": {
    "properties": {
      "primary_topic": {
        "properties": {
          "field": {
            "properties": {
              "id": {
                "type": "keyword"
              }
            }
          }
        }
      }
    }
  }
*/

var mappings = map[openalex.FileEntityType]interface{}{
	openalex.WorksFileEntityType: map[string]interface{}{
		"properties": map[string]interface{}{
			"abstract": map[string]interface{}{
				"type": "text",
			},
			"primary_topic": map[string]interface{}{
				"properties": map[string]interface{}{
					"domain": map[string]interface{}{
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type": "keyword",
							},
						},
					},
					"field": map[string]interface{}{
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type": "keyword",
							},
						},
					},
					"subfield": map[string]interface{}{
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type": "keyword",
							},
						},
					},
				},
			},
			"topics": map[string]interface{}{
				"properties": map[string]interface{}{
					"domain": map[string]interface{}{
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type": "keyword",
							},
						},
					},
					"field": map[string]interface{}{
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type": "keyword",
							},
						},
					},
					"subfield": map[string]interface{}{
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type": "keyword",
							},
						},
					},
				},
			},
			"concepts": map[string]interface{}{
				"properties": map[string]interface{}{
					"score": map[string]interface{}{
						"type": "float",
					},
				},
			},
		},
	},
}

func createSettings(entityType openalex.FileEntityType) {
	// create a mapping for the entity type
	// get the index
	index := string(entityType)
	// create a setting map where there are only 1 shard and 0 replicas
	// also set the refresh interval to -1
	payload := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"refresh_interval":   "-1",
		},
	}
	// add the mappings to the payload if they exist
	mapping, ok := mappings[entityType]
	if ok {
		payload["mappings"] = mapping
	}
	// convert the settings to bytes
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		slog.With("err", err).Error("could not marshal settings")
	}
	// create the request
	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  bytes.NewReader(payloadBytes),
	}
	// send the request
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		slog.With("err", err).Error("could not create mapping")
	}
	defer res.Body.Close()
}

var counter = 0

func ElasticLineHandler(filePath string, line string) error {
	// remove all "https://openalex.org/domains/" from line
	line = strings.ReplaceAll(line, "https://openalex.org/domains/", "")
	line = strings.ReplaceAll(line, "https://openalex.org/", "")

	fileEntityType, err := openalex.GetEntityType(filePath)
	if err != nil {
		slog.With("err", err).Error("could not get entity type")
	}

	// check if the entity type is in the indexers
	_, ok := indexers[fileEntityType]
	if !ok {
		slog.With("fileEntityType", fileEntityType).Error("bulk indexer for entity type not found")
		return fmt.Errorf("bulk indexer for entity type not found")
	}

	// parse the line into a map[string]interface{}
	lineContent := make(map[string]interface{})
	err = json.Unmarshal([]byte(line), &lineContent)
	if err != nil {
		slog.With("err", err).Error("could not unmarshal line")
		return err
	}
	// extract the id from the line
	entityID, ok := lineContent["id"].(string)
	if !ok {
		slog.Error("line has no id")
		return fmt.Errorf("line has no id")
	}
	// remove https://openalex.org/ from the id
	entityID = strings.Replace(entityID, "https://openalex.org/", "", 1)

	bulkIndexer := indexers[fileEntityType]

	// delete the id from the line
	delete(lineContent, "id")

	if fileEntityType == openalex.WorksFileEntityType {
		field, ok := lineContent["abstract_inverted_index"]
		if ok {
			// Assert field as map[string]interface{} first
			if invertedMap, ok := field.(map[string]interface{}); ok {
				transformedMap := make(map[string][]int)

				for key, value := range invertedMap {
					// Check if the value is []interface{} and convert to []int
					if intSlice, ok := value.([]interface{}); ok {
						intValues := make([]int, len(intSlice))
						for i, v := range intSlice {
							if intValue, ok := v.(float64); ok { // JSON numbers are float64
								intValues[i] = int(intValue)
							} else {
								// Handle unexpected type
								intValues = nil
								break
							}
						}
						if intValues != nil {
							transformedMap[key] = intValues
						}
					}
				}

				// Use the transformed map if all conversions succeeded
				lineContent["abstract"] = openalex.TransformInvertedIndexToAbstract(transformedMap)
			} else {
				// If it isn't a map[string]interface{}, set to nil
				lineContent["abstract"] = nil
			}
		} else {
			// No abstract_inverted_index field found
			lineContent["abstract"] = nil
		}
		delete(lineContent, "abstract_inverted_index")
	}

	// marshal the line to a string
	payload, err := json.Marshal(lineContent)
	if err != nil {
		slog.With("err", err).Error("could not marshal line")
	}

	// add the entity to the bulk indexer
	err = bulkIndexer.Add(
		context.Background(),
		esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: entityID,
			Body:       bytes.NewReader(payload),
		},
	)
	if err != nil {
		slog.With("err", err).Error("could not add entity to bulk")
		return err
	}

	counter++

	if counter%10000 == 0 {
		s := bulkIndexer.Stats()
		slog.
			With("NumAdded", s.NumAdded).
			With("NumFlushed", s.NumFlushed).
			With("NumFailed", s.NumFailed).
			With("NumIndexed", s.NumIndexed).
			With("NumCreated", s.NumCreated).
			With("NumUpdated", s.NumUpdated).
			With("NumDeleted", s.NumDeleted).
			With("NumRequests", s.NumRequests).
			Info("bulk indexer stats")
	}

	return nil
}

func main() {

	openAlexDir := env.FallbackEnvVariable("OPENALEX_DIR", "/media/seb/T18-1/openalex-data/data")

	entityTypes := []openalex.FileEntityType{
		openalex.AuthorsFileEntityType,
		openalex.ConceptsFileEntityType,
		openalex.FundersFileEntityType,
		openalex.InstitutionsFileEntityType,
		openalex.PublishersFileEntityType,
		openalex.SourcesFileEntityType,
		openalex.WorksFileEntityType,
		openalex.TopicsFileEntityType,
		openalex.DomainsFileEntityType,
	}
	for _, entityType := range entityTypes {
		createSettings(entityType)
		createBulkIndexer(entityType)
	}

	p := openalex.Processor{
		DirectoryPath:   openAlexDir,
		StateHandler:    nil,
		LineHandler:     ElasticLineHandler,
		MergedIdHandler: nil,
	}

	err := p.ProcessDirectory()
	if err != nil {
		slog.With("err", err).Error("error processing directory")
	}

	for _, entityType := range entityTypes {
		bulkIndexer := indexers[entityType]
		err = bulkIndexer.Close(context.Background())
		if err != nil {
			slog.With("err", err).Error("error closing bulk indexer")
		}
	}

	slog.Info("finished processing directory")

}
