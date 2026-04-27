package display

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	coreio "dappco.re/go/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInjectAppPreloads_FromManifest(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "preloads:\n  - path: preload.js\n", 0o644))

	svc, err := New()
	require.NoError(t, err)

	script, err := svc.injectAppPreloads(filepath.Join(root, "index.html"))
	require.NoError(t, err)
	require.True(t, strings.Contains(script, "__manifestLoaded"))
}

func TestInjectAppPreloads_RejectsTraversal(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "preloads:\n  - path: ../preload.js\n", 0o644))

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.injectAppPreloads(filepath.Join(root, "index.html"))
	require.Error(t, err)
}

func TestManifest_SafeManifestPreloadPath_Good(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "preload.js")
	require.NoError(t, coreio.Local.WriteMode(target, "globalThis.ready = true;", 0o644))
	expected, err := filepath.EvalSymlinks(target)
	require.NoError(t, err)
	got, err := safeManifestPreloadPath(root, "preload.js")

	require.NoError(t, err)
	assert.Equal(t, expected, got)
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
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(outside, "preload.js"), "globalThis.__outside = true;", 0o644))
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, "assets")))
	if err := os.Symlink(outside, filepath.Join(root, "assets", "linked")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	_, err := safeManifestPreloadPath(filepath.Join(root, "assets"), filepath.Join("linked", "preload.js"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "escapes")
}

func TestManifest_DiscoverManifestPath_Good(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	manifestPath := filepath.Join(root, ".core", "view.yaml")
	require.NoError(t, coreio.Local.WriteMode(manifestPath, "name: demo\n", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))

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
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	manifestPath := filepath.Join(root, ".core", "view.yaml")
	require.NoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath(root)

	require.NoError(t, err)
	assert.Equal(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_Good(t *testing.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)
	manifestPath := filepath.Join(home, ".core", "apps", "example.com", ".core", "view.yaml")
	require.NoError(t, coreio.Local.EnsureDir(filepath.Dir(manifestPath)))
	require.NoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath("https://example.com/index.html")

	require.NoError(t, err)
	assert.Equal(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_StripsPort(t *testing.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)
	manifestPath := filepath.Join(home, ".core", "apps", "example.com", ".core", "view.yaml")
	require.NoError(t, coreio.Local.EnsureDir(filepath.Dir(manifestPath)))
	require.NoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath("https://example.com:8080/x")

	require.NoError(t, err)
	assert.Equal(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_IPv6Literal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)
	manifestPath := filepath.Join(home, ".core", "apps", "::1", ".core", "view.yaml")
	require.NoError(t, coreio.Local.EnsureDir(filepath.Dir(manifestPath)))
	require.NoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath("https://[::1]/x")

	require.NoError(t, err)
	assert.Equal(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_RejectsControlCharacter(t *testing.T) {
	t.Setenv("DIR_HOME", t.TempDir())

	_, err := discoverManifestPath("https://bad\nhost/x")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "control")
}

func TestManifest_DiscoverManifestPath_RemoteHost_RejectsTraversalHost(t *testing.T) {
	t.Setenv("DIR_HOME", t.TempDir())

	_, err := discoverManifestPath("https://../x")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "relative path")
}

func TestManifest_ManifestWindowConfig_Good(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"windows:",
		"  main:",
		"    title: Core GUI",
		"    width: 1280",
		"    height: 720",
		"    preload: true",
	}, "\n"), 0o644))

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
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "windows: [\n", 0o644))

	svc, err := New()
	require.NoError(t, err)

	got := svc.manifestWindowConfig(filepath.Join(root, "index.html"))

	assert.Nil(t, got)
}

func TestManifest_ManifestWindowConfig_ReturnsCopy(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"windows:",
		"  main:",
		"    title: Core GUI",
		"    width: 1280",
		"    height: 720",
	}, "\n"), 0o644))

	svc, err := New()
	require.NoError(t, err)

	first := svc.manifestWindowConfig(filepath.Join(root, "index.html"))
	require.NotNil(t, first)
	first["main"] = ManifestWindow{Title: "mutated"}

	second := svc.manifestWindowConfig(filepath.Join(root, "index.html"))
	require.NotNil(t, second)
	assert.Equal(t, "Core GUI", second["main"].Title)
}

func TestManifest_LoadManifestForOrigin_RejectsOversizedFile(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "name: "+strings.Repeat("a", maxViewManifestBytes), 0o644))

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.loadManifestForOrigin(filepath.Join(root, "index.html"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds")
}

func TestManifest_LoadManifestForOrigin_Concurrent(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"name: demo",
		"windows:",
		"  main:",
		"    title: Core GUI",
	}, "\n"), 0o644))

	svc, err := New()
	require.NoError(t, err)

	var wg sync.WaitGroup
	errs := make(chan error, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loaded, loadErr := svc.loadManifestForOrigin(filepath.Join(root, "index.html"))
			if loadErr != nil {
				errs <- loadErr
				return
			}
			if loaded == nil || loaded.Manifest.Name != "demo" {
				errs <- assert.AnError
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
}

func TestManifest_ManifestBaseDir_Good(t *testing.T) {
	assert.Equal(t, "/tmp/app", manifestBaseDir("/tmp/app/.core/view.yaml"))
	assert.Equal(t, "/tmp/app/assets", manifestBaseDir("/tmp/app/assets/view.yaml"))
}

func TestManifest_ManifestBaseDir_Bad(t *testing.T) {
	assert.Equal(t, ".", manifestBaseDir(".core/view.yaml"))
}

func TestManifest_ManifestBaseDir_Ugly(t *testing.T) {
	assert.Equal(t, "/", manifestBaseDir("/.core/view.yaml"))
}

func TestManifest_SafeManifestRelativePath_Good(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "preload.js")
	require.NoError(t, coreio.Local.WriteMode(target, "globalThis.ready = true;", 0o644))
	expected, err := filepath.EvalSymlinks(target)
	require.NoError(t, err)

	got, err := safeManifestRelativePath(root, "preload.js", "preload path")

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestManifest_SafeManifestRelativePath_Bad(t *testing.T) {
	root := t.TempDir()

	_, err := safeManifestRelativePath(root, "", "preload path")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestManifest_SafeManifestRelativePath_Bad_MissingFile(t *testing.T) {
	root := t.TempDir()

	_, err := safeManifestRelativePath(root, "missing.js", "preload path")

	require.Error(t, err)
}

func TestManifest_SafeManifestRelativePath_Ugly(t *testing.T) {
	root := t.TempDir()

	_, err := safeManifestRelativePath(root, "../escape.js", "preload path")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "escapes")
}
