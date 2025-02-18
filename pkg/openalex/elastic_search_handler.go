package openalex

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/SbstnErhrdt/env"
	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

/*
BOOTSTRAP
1. kompletten snapshot laden
2. entpacken
3. processdirectory (422gb -> unpacked = 1.6TB) anwenden
3. & cache progress, savepoint from zip file level if failed

UPDATE
4. manifest laden
5. filter new date folders
6. download all new date folders
7. processdirectory
7. & cache progress, savepoint from zip file level if failed
*/

// flags
// targetIndex is the index to write to
var targetIndex = "docdb_cos"

// esClient is the client for the elasticsearch 8
var esClient *elasticsearch.Client

// bi is the bulk indexer
var bi esutil.BulkIndexer

// countSuccessful is the count of successful uploads
var countSuccessful uint64

// countFailed is the count of failed uploads
var countFailed uint64

func init() {
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

// ElasticSearchHandler is the handler for the content
// TODO Refactor for
func ElasticSearchHandler() ParsedEntityLineHandler {

	//
	return func(fileEntityType FileEntityType, entity any) error {
		// Parse to according entity

		// save to db
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				DocumentID: dbDoc.Id,

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(payload),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
					// for every 1000 successful uploads, print the count
					if countSuccessful > 0 && countSuccessful%1000 == 0 {
						slog.With("countSuccessful", countSuccessful).Info("count successful")
					}
				},
				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					// count failed
					atomic.AddUint64(&countFailed, 1)
					if err != nil {
						slog.With("err", err).Error("could not index document")
					} else {
						slog.With("err", res.Error.Type).Error(res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			slog.With("err", err).Error("could not add document to bulk")
		}
		return
	}
}
