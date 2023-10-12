package openalex

import (
	"testing"
)

func TestReadManifest(t *testing.T) {
	manifestsHashes := make(map[string]struct{})
	for _, manifestUrl := range AllManifestUrls {
		manifest, err := ReadManifestFromS3Url(manifestUrl)
		if err != nil {
			t.Fatal(err)
		}
		hash, err := manifest.Hash()
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := manifestsHashes[hash]; ok {
			t.Fatal("duplicate manifest url", manifestUrl)
		}
	}
}
