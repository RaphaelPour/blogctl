package metadata

import (
	"io/ioutil"
	"os"
	"testing"
"fmt"
"time"
	"github.com/stretchr/testify/require"
)

func TestLoadMetadata(t *testing.T) {

	/* Save example metadata json */
	title := "Exit the known"
	date := time.Now()
	data := fmt.Sprintf(
		`{"title": "%s", "createdAt": "%s"}`,
		title,
		date.String(),
	)

	tmpDir, err := ioutil.TempDir("", "post-")
	require.Nil(t, err)

	path := GetMetadataPath(tmpDir)

	require.Nil(t, ioutil.WriteFile(path, []byte(data), os.ModePerm))

	meta, err := Load(tmpDir)
	require.Nil(t, err)

	require.Equal(t, title, meta.Title)
	require.Equal(t, date, meta.CreatedAt)
}

func TestNewInvalidMetadata(t *testing.T) {
}

func TestSaveMetadata(t *testing.T) {
}

func TestMetadataFile(t *testing.T) {
	require.Equal(t, GetMetadataPath("/tmp"), "/tmp/metadata.json")
}
