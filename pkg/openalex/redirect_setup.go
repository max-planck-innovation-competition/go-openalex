package openalex

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/SbstnErhrdt/env"
	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var esClient *elasticsearch.Client

func InitElasticSearch() {
	// load env
	env.LoadEnvFiles()
	// create es client
	esUrl := fmt.Sprintf("%s://%s:%s",
		os.Getenv("ELASTICSEARCH_PROTOCOL"),
		os.Getenv("ELASTICSEARCH_HOST"),
		os.Getenv("ELASTICSEARCH_PORT"))
	// add retry
	retryBackoff := backoff.NewExponentialBackOff()
	es8, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{esUrl},
			Username:  os.Getenv("ELASTICSEARCH_USERNAME"),
			Password:  os.Getenv("ELASTICSEARCH_PASSWORD"),
			// Retry on 429 TooManyRequests statuses
			RetryOnStatus: []int{502, 503, 504, 429},
			// Configure the backoff function
			RetryBackoff: func(i int) time.Duration {
				if i == 1 {
					retryBackoff.Reset()
				}
				return retryBackoff.NextBackOff()
			},
			// Retry up to 5 attempts
			//
			MaxRetries: 5,
		},
	)
	if err != nil {
		slog.With("err", err).Error("could not create es client")
		panic(err)
	}
	esClient = es8
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	PublishDate string `json:"publish_date"`
	RedirectTo  string `json:"redirect_to"`
}

func addBook(book Book) {
	// Convert book to JSON
	bookJSON, err := json.Marshal(book)
	if err != nil {
		log.Fatalf("Error marshaling the book: %s", err)
	}

	// Index the book (create or update)
	req := esapi.IndexRequest{
		Index:      "books",
		DocumentID: book.ID,
		Body:       bytes.NewReader(bookJSON),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Fatalf("Error indexing the book: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error response from Elasticsearch: %s", res.Status())
	} else {
		fmt.Println("Book entry created successfully!")
	}
}

func findBookByID(es *elasticsearch.Client, id string) (*Book, error) {
	// Construct the search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": id,
			},
		},
	}

	// Encode the query to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding book query: %s", err)
	}

	// Perform the search request
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("books"), // Replace with your index name
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	// Check for errors in the search response
	if res.IsError() {
		return nil, fmt.Errorf("error response: %s", res.Status())
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	var foundBook Book

	// Extract the book details
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if total["value"].(float64) == 0 {
				return nil, nil
			}
		}
		if hitArray, ok := hits["hits"].([]interface{}); ok && len(hitArray) > 0 {
			for _, hit := range hitArray {
				if hitMap, ok := hit.(map[string]interface{}); ok {
					if source, ok := hitMap["_source"].(map[string]interface{}); ok {
						sourceBytes, err := json.Marshal(source)
						if err != nil {
							return nil, fmt.Errorf("error marshaling source: %s", err)
						}
						if err := json.Unmarshal(sourceBytes, &foundBook); err != nil {
							return nil, fmt.Errorf("error unmarshaling book: %s", err)
						}
					}
				}
			}
		}
	}

	if foundBook.RedirectTo != "" {
		return findBookByID(es, foundBook.RedirectTo)
	}

	return &foundBook, nil
}
