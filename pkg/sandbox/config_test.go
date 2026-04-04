package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func TestConfigMarshal(t *testing.T) {
	cfg := Config{
		Command:        "echo",
		IncludeParents: ReadMeta | ReadData,
		FoldersInParents: []FolderPermission{
			{
				Path:       "/path/to/folder",
				Permission: ReadMeta,
			},
			{
				Path:       "/other",
				Permission: ReadMeta | ReadData | WriteMeta | WriteData,
			},
		},
		ExportEnvironment: []string{"PWD", "HOME"},
	}

	yamlData, err := yaml.Marshal(cfg)
	require.NoError(t, err)

	expected := map[any]any{
		"command":         "echo",
		"include_parents": []any{"read-meta", "read-data"},
		"folders_in_parents": []any{
			map[string]any{
				"path":        "/path/to/folder",
				"permissions": []any{"read-meta"},
			},
			map[string]any{
				"path":        "/other",
				"permissions": []any{"read-meta", "read-data", "write-meta", "write-data"},
			},
		},
		"export_env": []any{"PWD", "HOME"},
	}

	var result map[any]any
	err = yaml.Unmarshal(yamlData, &result)
	require.NoError(t, err)

	assert.Equal(t, expected, result)

	restored := Config{}
	err = yaml.Unmarshal(yamlData, &restored)
	require.NoError(t, err)
	assert.Equal(t, cfg, restored)
}
