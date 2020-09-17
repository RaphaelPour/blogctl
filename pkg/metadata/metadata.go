package metadata // import "github.com/RaphaelPour/blogctl/metadata"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

/**
 * Metadata
 *
 * The metadata.json stores additional information about
 * a blog post.
 */
type Metadata struct {
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

const METADATA_FILE = "metadata.json"

func GetMetadataPath(postPath string) string {
	return filepath.Join(postPath, METADATA_FILE)
}

func LoadMetadata(postPath string) (*Metadata, error) {

	raw, err := ioutil.ReadFile(GetMetadataPath(postPath))
	if err != nil {
		return nil, fmt.Errorf("Error reading metadata: %s", err)
	}

	metadata := new(Metadata)
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return nil, fmt.Errorf("Error parsing metadata: %s", err)
	}

	return metadata, nil
}

func (m Metadata) Save(postPath string) error {

	raw, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Error converting metadata to json: %s", err)
	}

	if err := ioutil.WriteFile(GetMetadataPath(postPath), raw, os.ModePerm); err != nil {
		return fmt.Errorf("Error writing metadata: %s", err)
	}

	return nil
}
