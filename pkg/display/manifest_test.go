package display

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInjectAppPreloads_FromManifest(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "preload.js"), []byte("globalThis.__manifestLoaded = true;"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte("preloads:\n  - path: preload.js\n"), 0o644))

	svc, err := New()
	require.NoError(t, err)

	script, err := svc.injectAppPreloads(filepath.Join(root, "index.html"))
	require.NoError(t, err)
	require.True(t, strings.Contains(script, "__manifestLoaded"))
}

func TestInjectAppPreloads_RejectsTraversal(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "preload.js"), []byte("globalThis.__manifestLoaded = true;"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte("preloads:\n  - path: ../preload.js\n"), 0o644))

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.injectAppPreloads(filepath.Join(root, "index.html"))
	require.Error(t, err)
}

func TestManifest_SafeManifestPreloadPath_Good(t *testing.T) {
	root := t.TempDir()
	got, err := safeManifestPreloadPath(root, "preload.js")

	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "preload.js"), got)
}

func TestManifest_SafeManifestPreloadPath_Bad(t *testing.T) {
	root := t.TempDir()
	_, err := safeManifestPreloadPath(root, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestManifest_SafeManifestPreloadPath_Ugly(t *testing.T) {
	root := t.TempDir()
	_, err := safeManifestPreloadPath(root, "../preload.js")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "escapes")
}

func TestManifest_DiscoverManifestPath_Good(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	manifestPath := filepath.Join(root, ".core", "view.yaml")
	require.NoError(t, os.WriteFile(manifestPath, []byte("name: demo\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))

	got, err := discoverManifestPath(filepath.Join(root, "index.html"))

	require.NoError(t, err)
	assert.Equal(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_Bad(t *testing.T) {
	_, err := discoverManifestPath(filepath.Join(t.TempDir(), "missing.html"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestManifest_DiscoverManifestPath_Ugly(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	manifestPath := filepath.Join(root, ".core", "view.yaml")
	require.NoError(t, os.WriteFile(manifestPath, []byte("name: remote\n"), 0o644))

	got, err := discoverManifestPath(root)

	require.NoError(t, err)
	assert.Equal(t, manifestPath, got)
}
