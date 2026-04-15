package display

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadModelCatalogFromConfigRoot_Good(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".core")
	configPath := filepath.Join(root, "gui", "config.yaml")
	modelsPath := filepath.Join(root, "models.yaml")

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0o755))
	require.NoError(t, os.WriteFile(modelsPath, []byte(`
default_model: lemma-lite
selected_model: lemma-lite
models:
  - name: lemma-lite
    architecture: gemma4
    quant_bits: 4
    size_bytes: 987654321
    loaded: true
    backend: mlx
    supports_vision: true
  - name: lemmy
    backend: llama.cpp
`), 0o644))

	svc, err := NewService()
	require.NoError(t, err)
	svc.loadConfigFrom(configPath)

	models := svc.chat.Models()
	require.NotEmpty(t, models)
	assert.Equal(t, "lemma-lite", svc.chat.SelectedModel())
	assert.Equal(t, "lemma-lite", svc.chat.Settings().DefaultModel)

	var custom *ModelEntry
	for index := range models {
		if models[index].Name == "lemma-lite" {
			custom = &models[index]
			break
		}
	}
	require.NotNil(t, custom)
	assert.Equal(t, "gemma4", custom.Architecture)
	assert.Equal(t, 4, custom.QuantBits)
	assert.Equal(t, int64(987654321), custom.SizeBytes)
	assert.Equal(t, "mlx", custom.Backend)
	assert.True(t, custom.Loaded)
	assert.True(t, custom.SupportsVision)
}

func TestLoadModelCatalog_MapShape_Good(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".core")
	configPath := filepath.Join(root, "gui", "config.yaml")
	modelsPath := filepath.Join(root, "models.yaml")

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0o755))
	require.NoError(t, os.WriteFile(modelsPath, []byte(`
models:
  lemrd:
    architecture: qwen3
    quant_bits: 8
    size_bytes: 123456789
    backend: ollama
    supports_vision: false
`), 0o644))

	svc, err := NewService()
	require.NoError(t, err)
	svc.loadConfigFrom(configPath)

	models := svc.chat.Models()
	assert.Contains(t, modelNames(models), "lemrd")
}

func TestResolveScheme_CoreModelsIncludesBackendSummary_Good(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".core")
	configPath := filepath.Join(root, "gui", "config.yaml")
	modelsPath := filepath.Join(root, "models.yaml")

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0o755))
	require.NoError(t, os.WriteFile(modelsPath, []byte(`
selected_model: lemer
models:
  - name: lemer
    architecture: gemma3
    quant_bits: 4
    size_bytes: 1500000000
    loaded: true
    backend: metal
    supports_vision: true
  - name: lemmy
    architecture: qwen3
    quant_bits: 4
    size_bytes: 1100000000
    loaded: false
    backend: ollama
    supports_vision: false
`), 0o644))

	svc, err := NewService()
	require.NoError(t, err)
	svc.loadConfigFrom(configPath)
	svc.registerBuiltinSchemes()

	response, err := svc.ResolveScheme(context.Background(), "core://models")
	require.NoError(t, err)

	assert.Equal(t, "lemer", response.Data["selected_model"])
	assert.Equal(t, 1, response.Data["loaded_count"])
	assert.Equal(t, 3, response.Data["model_count"])

	backends, ok := response.Data["backends"].(map[string]map[string]any)
	require.True(t, ok)
	assert.Contains(t, backends, "metal")
	assert.Contains(t, backends, "ollama")
	assert.Equal(t, 1, backends["metal"]["loaded"])
	assert.Equal(t, true, backends["metal"]["supports_vision"])
}

func modelNames(models []ModelEntry) []string {
	names := make([]string, 0, len(models))
	for _, model := range models {
		names = append(names, model.Name)
	}
	return names
}
