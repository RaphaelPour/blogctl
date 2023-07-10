package metadata

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadMetadata(t *testing.T) {

	/* Save example metadata json */
	title := "99 diets that make you actually fatter"
	date := int64(1234)
	data := fmt.Sprintf(`{"title": "%s", "createdAt": %d}`,
		title, date,
	)

	tmpDir, err := os.MkdirTemp("", "post-")
	require.Nil(t, err)

	path := GetMetadataPath(tmpDir)

	require.Nil(t, os.WriteFile(path, []byte(data), os.ModePerm))

	meta, err := Load(tmpDir)
	require.Nil(t, err)

	require.Equal(t, title, meta.Title)
	require.Equal(t, date, meta.CreatedAt)
}

func TestNewInvalidMetadata(t *testing.T) {
	/* Save example metadata json */
	data := `{"title": "How I lost -5 kgs in one week", "createdAt": 1234`

	tmpDir, err := os.MkdirTemp("", "post-")
	require.Nil(t, err)

	path := GetMetadataPath(tmpDir)

	require.Nil(t, os.WriteFile(path, []byte(data), os.ModePerm))

	meta, err := Load(tmpDir)
	require.NotNil(t, err)
	require.Regexp(t, "^Error parsing metadata:", err)
	require.Nil(t, meta)
}

func TestSaveMetadata(t *testing.T) {
	meta := &Metadata{
		Title:     "99 problems with PHP",
		CreatedAt: 788918400,
	}

	tmpDir, err := os.MkdirTemp("", "post-")
	require.Nil(t, err)

	err = meta.Save(tmpDir)
	require.Nil(t, err)
	require.FileExists(t, GetMetadataPath(tmpDir))
}

func TestMetadataFile(t *testing.T) {
	require.Equal(t, GetMetadataPath("/tmp"), "/tmp/metadata.json")
}
