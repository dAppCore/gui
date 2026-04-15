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

func TestManifest_SafeManifestPreloadPath_RejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "preload.js"), []byte("globalThis.__outside = true;"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "assets"), 0o755))
	require.NoError(t, os.Symlink(outside, filepath.Join(root, "assets", "linked")))

	_, err := safeManifestPreloadPath(filepath.Join(root, "assets"), filepath.Join("linked", "preload.js"))

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

func TestManifest_ManifestWindowConfig_Good(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte(strings.Join([]string{
		"windows:",
		"  main:",
		"    title: Core GUI",
		"    width: 1280",
		"    height: 720",
		"    preload: true",
	}, "\n")), 0o644))

	svc, err := New()
	require.NoError(t, err)

	got := svc.manifestWindowConfig(filepath.Join(root, "index.html"))

	require.NotNil(t, got)
	require.Contains(t, got, "main")
	assert.Equal(t, "Core GUI", got["main"].Title)
	assert.Equal(t, 1280, got["main"].Width)
	assert.Equal(t, 720, got["main"].Height)
	assert.True(t, got["main"].Preload)
}

func TestManifest_ManifestWindowConfig_Bad(t *testing.T) {
	svc, err := New()
	require.NoError(t, err)

	got := svc.manifestWindowConfig(filepath.Join(t.TempDir(), "missing.html"))

	assert.Nil(t, got)
}

func TestManifest_ManifestWindowConfig_Ugly(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte("windows: [\n"), 0o644))

	svc, err := New()
	require.NoError(t, err)

	got := svc.manifestWindowConfig(filepath.Join(root, "index.html"))

	assert.Nil(t, got)
}
