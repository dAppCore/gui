package display

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
