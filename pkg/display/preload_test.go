package display

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type capturePreloadWindow struct {
	script string
}

func (c *capturePreloadWindow) ExecJS(script string) {
	c.script = script
}

func TestPreloadScript_Good(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".core")
	configPath := filepath.Join(root, "gui", "config.yaml")
	viewPath := filepath.Join(root, "view.yaml")

	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0o755))
	require.NoError(t, os.WriteFile(viewPath, []byte(`preloads:
  - origin: https://app.example.com
    script: |
      window.__appPreloadLoaded = true;
`), 0o644))

	svc, err := NewService()
	require.NoError(t, err)
	svc.loadConfigFrom(configPath)
	svc.registerBuiltinSchemes()

	require.NoError(t, svc.SaveBrowserStorageState("https://app.example.com", OriginStorageState{
		LocalStorage:   map[string]string{"theme": "dark"},
		SessionStorage: map[string]string{"session_id": "abc123"},
	}))

	script, err := svc.PreloadScript("https://app.example.com/dashboard")
	require.NoError(t, err)

	assert.Contains(t, script, "window.__appPreloadLoaded = true;")
	assert.Contains(t, script, "root.core.ml.generate")
	assert.Contains(t, script, "root.electron = root.electron ||")
	assert.Contains(t, script, "root.electron.Notification")
	assert.Contains(t, script, "notification:show")
	assert.Contains(t, script, `"theme":"dark"`)
	assert.Contains(t, script, `"session_id":"abc123"`)
}

func TestInjectPreload_Good(t *testing.T) {
	svc, err := NewService()
	require.NoError(t, err)

	target := &capturePreloadWindow{}
	err = svc.InjectPreload(target, "https://app.example.com")
	require.NoError(t, err)

	assert.Contains(t, target.script, "root.core.ml.generate")
	assert.Contains(t, target.script, "root.core.storage")
	assert.Contains(t, target.script, "root.electron.Notification")
}

func TestBrowserStoragePersistenceAndSearch_Good(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".core")
	configPath := filepath.Join(root, "gui", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(configPath), 0o755))

	svc, err := NewService()
	require.NoError(t, err)
	svc.loadConfigFrom(configPath)
	svc.registerBuiltinSchemes()

	require.NoError(t, svc.SaveBrowserStorageState("https://app.example.com", OriginStorageState{
		LocalStorage: map[string]string{
			"invoice_status": "paid",
		},
		CacheStorage: map[string]string{
			"GET https://api.example.com/invoices/42": `{"invoice":"paid"}`,
		},
	}))

	response, err := svc.ResolveScheme(context.Background(), "core://store?q=invoice")
	require.NoError(t, err)

	results, ok := response.Data["results"].([]StoreSearchResult)
	require.True(t, ok)
	require.NotEmpty(t, results)
	assert.Equal(t, "https://app.example.com", results[0].Origin)
	assert.NotEmpty(t, results[0].StorageArea)
	assert.NotEmpty(t, results[0].Key)

	reloaded, err := NewService()
	require.NoError(t, err)
	reloaded.loadConfigFrom(configPath)

	state, ok := reloaded.BrowserStorageState("https://app.example.com")
	require.True(t, ok)
	assert.Equal(t, "paid", state.LocalStorage["invoice_status"])
	assert.Contains(t, state.CacheStorage["GET https://api.example.com/invoices/42"], "paid")
}

func TestResolveScheme_CoreRoutes_Good(t *testing.T) {
	svc, err := NewService()
	require.NoError(t, err)
	svc.registerBuiltinSchemes()

	for _, rawURL := range []string{
		"core://network",
		"core://agent",
		"core://wallet",
		"core://identity",
	} {
		response, err := svc.ResolveScheme(context.Background(), rawURL)
		require.NoError(t, err)
		assert.Equal(t, "core", response.Scheme)
		assert.Equal(t, 200, response.StatusCode)
		assert.Contains(t, response.Data, "available")
	}
}
