package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	data := `{"Version":1337}`
	tmpDir, err := os.MkdirTemp("", "post-")
	require.Nil(t, err)

	path := ConfigPath(tmpDir)

	require.Nil(t, os.WriteFile(path, []byte(data), os.ModePerm))

	cfg, err := Load(tmpDir)
	require.Nil(t, err)
	require.Equal(t, 1337, cfg.Version)
}

func TestLoadBadPath(t *testing.T) {
	cfg, err := Load("does-not-exist")
	require.Nil(t, cfg)
	require.ErrorContains(t, err, "Error reading config")
}

func TestLoadBadConfig(t *testing.T) {
	data := `{`
	tmpDir, err := os.MkdirTemp("", "post-")
	require.Nil(t, err)

	path := ConfigPath(tmpDir)

	require.Nil(t, os.WriteFile(path, []byte(data), os.ModePerm))

	cfg, err := Load(tmpDir)
	require.Nil(t, cfg)
	require.ErrorContains(t, err, "Error parsing config")
}
