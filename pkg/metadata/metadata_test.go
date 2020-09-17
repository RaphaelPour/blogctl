package metadata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadMetadata(t *testing.T) {
}

func TestNewInvalidMetadata(t *testing.T) {
}

func TestSaveMetadata(t *testing.T) {
}

func TestMetadataFile(t *testing.T) {
	require.Equal(t, GetMetadataPath("/tmp"), "/tmp/metadata.json")
}
