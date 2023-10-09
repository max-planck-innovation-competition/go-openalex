package openalex

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
)

// ErrStatusNotOK is returned when the status code is not OK
var ErrStatusNotOK = errors.New("status code is not OK")

// ReadManifestFromS3Url reads the manifest file from S3 and returns a Manifest struct
func ReadManifestFromS3Url(s3Url ManifestUrl) (result *Manifest, err error) {
	logger := slog.With("s3Url", s3Url)
	// Create an HTTP GET request
	req, err := http.NewRequest("GET", string(s3Url), nil)
	if err != nil {
		logger.With("error", err).Error("Failed to create HTTP request")
		return nil, err
	}

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logger.With("error", err).Error("Failed to send HTTP request")
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		err = ErrStatusNotOK
		logger.
			With("statusCode", resp.StatusCode).
			With("err", err).
			Error("Failed to fetch S3 object")
		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into the manifest struct
	var manifest Manifest
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		logger.With("error", err).Error("Failed to unmarshal JSON")
		return nil, err
	}

	return &manifest, nil
}
