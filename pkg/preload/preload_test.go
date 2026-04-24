package preload

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type captureWebview struct {
	scripts []string
}

func (c *captureWebview) ExecJS(script string) {
	c.scripts = append(c.scripts, script)
}

func TestInjectPreload_Good(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(
		filepath.Join(root, ".core", "view.yaml"),
		[]byte("manifest:\n  preloads:\n    - path: preload.js\n"),
		0o644,
	))
	require.NoError(t, os.WriteFile(
		filepath.Join(root, "preload.js"),
		[]byte("globalThis.__manifestPreloadLoaded = true;"),
		0o644,
	))

	target := &captureWebview{}
	err := InjectPreload(target, "file://"+filepath.ToSlash(filepath.Join(root, "index.html")))
	require.NoError(t, err)
	require.Len(t, target.scripts, 1)

	script := target.scripts[0]
	assert.Contains(t, script, "globalThis.core.storage.local")
	assert.Contains(t, script, "globalThis.core.ml = globalThis.core.ml ||")
	assert.Contains(t, script, "globalThis.electron = electron")
	assert.Contains(t, script, "globalThis.__manifestPreloadLoaded = true;")
}

func TestInjectPreload_Bad(t *testing.T) {
	err := InjectPreload(nil, "http://localhost:3000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "preload target is required")
}

func TestInjectPreload_Ugly(t *testing.T) {
	target := &captureWebview{}
	err := InjectPreload(target, "https://example.com/app")
	require.NoError(t, err)
	require.Len(t, target.scripts, 1)

	script := target.scripts[0]
	assert.Contains(t, script, "globalThis.core.storage.local")
	assert.Contains(t, script, "globalThis.core.ml = globalThis.core.ml ||")
	assert.NotContains(t, script, "globalThis.electron = electron")
	assert.NotContains(t, script, "ipcRenderer")
}
