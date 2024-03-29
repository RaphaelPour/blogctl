package metadata

import (
	"encoding/json"
	"fmt"
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
	Title         string `json:"title"`
	Status        string `json:"status"`
	Static        bool   `json:"static"`
	CreatedAt     int64  `json:"createdAt"`
	FeaturedImage string `json:"featuredImage"`
}

const (
	METADATA_FILE = "metadata.json"

	PUBLIC_STATUS = "public"
	DRAFT_STATUS  = "draft"
)

func GetMetadataPath(postPath string) string {
	return filepath.Join(postPath, METADATA_FILE)
}

func Load(postPath string) (*Metadata, error) {

	raw, err := os.ReadFile(GetMetadataPath(postPath))
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

	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("Error converting metadata to json: %s", err)
	}

	if err := os.WriteFile(GetMetadataPath(postPath), raw, os.ModePerm); err != nil {
		return fmt.Errorf("Error writing metadata: %s", err)
	}

	return nil
}

func (m Metadata) Date() string {
	return time.Unix(m.CreatedAt, 0).Format(time.DateOnly)
}
