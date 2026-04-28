package display

import (
	core "dappco.re/go"
	"os"
	"path/filepath"
	"strings"
	"sync"

	coreio "dappco.re/go/io"
)

func TestInjectAppPreloads_FromManifest(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "preloads:\n  - path: preload.js\n", 0o644))

	svc, err := New()
	core.RequireNoError(t, err)

	script, err := svc.injectAppPreloads(filepath.Join(root, "index.html"))
	core.RequireNoError(t, err)
	core.RequireTrue(t, strings.Contains(script, "__manifestLoaded"))
}

func TestInjectAppPreloads_RejectsTraversal(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "preloads:\n  - path: ../preload.js\n", 0o644))

	svc, err := New()
	core.RequireNoError(t, err)

	_, err = svc.injectAppPreloads(filepath.Join(root, "index.html"))
	core.AssertError(t, err)
}

func TestManifest_SafeManifestPreloadPath_Good(t *core.T) {
	root := t.TempDir()
	target := filepath.Join(root, "preload.js")
	core.RequireNoError(t, coreio.Local.WriteMode(target, "globalThis.ready = true;", 0o644))
	expected, err := filepath.EvalSymlinks(target)
	core.RequireNoError(t, err)
	got, err := safeManifestPreloadPath(root, "preload.js")

	core.RequireNoError(t, err)
	core.AssertEqual(t, expected, got)
}

func TestManifest_SafeManifestPreloadPath_Bad(t *core.T) {
	root := t.TempDir()
	_, err := safeManifestPreloadPath(root, "")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "empty")
}

func TestManifest_SafeManifestPreloadPath_Ugly(t *core.T) {
	root := t.TempDir()
	_, err := safeManifestPreloadPath(root, "../preload.js")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "escapes")
}

func TestManifest_SafeManifestPreloadPath_RejectsSymlinkEscape(t *core.T) {
	root := t.TempDir()
	outside := t.TempDir()
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(outside, "preload.js"), "globalThis.__outside = true;", 0o644))
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, "assets")))
	if err := os.Symlink(outside, filepath.Join(root, "assets", "linked")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	_, err := safeManifestPreloadPath(filepath.Join(root, "assets"), filepath.Join("linked", "preload.js"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "escapes")
}

func TestManifest_DiscoverManifestPath_Good(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	manifestPath := filepath.Join(root, ".core", "view.yaml")
	core.RequireNoError(t, coreio.Local.WriteMode(manifestPath, "name: demo\n", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))

	got, err := discoverManifestPath(filepath.Join(root, "index.html"))

	core.RequireNoError(t, err)
	core.AssertEqual(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_Bad(t *core.T) {
	_, err := discoverManifestPath(filepath.Join(t.TempDir(), "missing.html"))

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "not found")
}

func TestManifest_DiscoverManifestPath_Ugly(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	manifestPath := filepath.Join(root, ".core", "view.yaml")
	core.RequireNoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath(root)

	core.RequireNoError(t, err)
	core.AssertEqual(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_Good(t *core.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)
	manifestPath := filepath.Join(home, ".core", "apps", "example.com", ".core", "view.yaml")
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Dir(manifestPath)))
	core.RequireNoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath("https://example.com/index.html")

	core.RequireNoError(t, err)
	core.AssertEqual(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_StripsPort(t *core.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)
	manifestPath := filepath.Join(home, ".core", "apps", "example.com", ".core", "view.yaml")
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Dir(manifestPath)))
	core.RequireNoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath("https://example.com:8080/x")

	core.RequireNoError(t, err)
	core.AssertEqual(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_IPv6Literal(t *core.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)
	manifestPath := filepath.Join(home, ".core", "apps", "::1", ".core", "view.yaml")
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Dir(manifestPath)))
	core.RequireNoError(t, coreio.Local.WriteMode(manifestPath, "name: remote\n", 0o644))

	got, err := discoverManifestPath("https://[::1]/x")

	core.RequireNoError(t, err)
	core.AssertEqual(t, manifestPath, got)
}

func TestManifest_DiscoverManifestPath_RemoteHost_RejectsControlCharacter(t *core.T) {
	t.Setenv("DIR_HOME", t.TempDir())

	_, err := discoverManifestPath("https://bad\nhost/x")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "control")
}

func TestManifest_DiscoverManifestPath_RemoteHost_RejectsTraversalHost(t *core.T) {
	t.Setenv("DIR_HOME", t.TempDir())

	_, err := discoverManifestPath("https://../x")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "relative path")
}

func TestManifest_ManifestWindowConfig_Good(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"windows:",
		"  main:",
		"    title: Core GUI",
		"    width: 1280",
		"    height: 720",
		"    preload: true",
	}, "\n"), 0o644))

	svc, err := New()
	core.RequireNoError(t, err)

	got := svc.manifestWindowConfig(filepath.Join(root, "index.html"))

	core.AssertNotNil(t, got)
	core.AssertContains(t, got, "main")
	core.AssertEqual(t, "Core GUI", got["main"].Title)
	core.AssertEqual(t, 1280, got["main"].Width)
	core.AssertEqual(t, 720, got["main"].Height)
	core.AssertTrue(t, got["main"].Preload)
}

func TestManifest_ManifestWindowConfig_Bad(t *core.T) {
	svc, err := New()
	core.RequireNoError(t, err)

	got := svc.manifestWindowConfig(filepath.Join(t.TempDir(), "missing.html"))

	core.AssertNil(t, got)
}

func TestManifest_ManifestWindowConfig_Ugly(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "windows: [\n", 0o644))

	svc, err := New()
	core.RequireNoError(t, err)

	got := svc.manifestWindowConfig(filepath.Join(root, "index.html"))

	core.AssertNil(t, got)
}

func TestManifest_ManifestWindowConfig_ReturnsCopy(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"windows:",
		"  main:",
		"    title: Core GUI",
		"    width: 1280",
		"    height: 720",
	}, "\n"), 0o644))

	svc, err := New()
	core.RequireNoError(t, err)

	first := svc.manifestWindowConfig(filepath.Join(root, "index.html"))
	core.AssertNotNil(t, first)
	first["main"] = ManifestWindow{Title: "mutated"}

	second := svc.manifestWindowConfig(filepath.Join(root, "index.html"))
	core.AssertNotNil(t, second)
	core.AssertEqual(t, "Core GUI", second["main"].Title)
}

func TestManifest_LoadManifestForOrigin_RejectsOversizedFile(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "name: "+strings.Repeat("a", maxViewManifestBytes), 0o644))

	svc, err := New()
	core.RequireNoError(t, err)

	_, err = svc.loadManifestForOrigin(filepath.Join(root, "index.html"))
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "exceeds")
}

func TestManifest_LoadManifestForOrigin_Concurrent(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"name: demo",
		"windows:",
		"  main:",
		"    title: Core GUI",
	}, "\n"), 0o644))

	svc, err := New()
	core.RequireNoError(t, err)

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
				errs <- core.AnError
			}
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		core.RequireNoError(t, err)
	}
}

func TestManifest_ManifestBaseDir_Good(t *core.T) {
	core.AssertEqual(t, "/tmp/app", manifestBaseDir("/tmp/app/.core/view.yaml"))
	core.AssertEqual(t, "/tmp/app/assets", manifestBaseDir("/tmp/app/assets/view.yaml"))
	core.AssertNotEmpty(t, core.Sprintf("%T", manifestBaseDir("/tmp/app/.core/view.yaml")))
}

func TestManifest_ManifestBaseDir_Bad(t *core.T) {
	core.AssertEqual(t, ".", manifestBaseDir(".core/view.yaml"))
	observedType := core.Sprintf("%T", manifestBaseDir(".core/view.yaml"))
	core.AssertNotEmpty(t, observedType)
}

func TestManifest_ManifestBaseDir_Ugly(t *core.T) {
	core.AssertEqual(t, "/", manifestBaseDir("/.core/view.yaml"))
	observedType := core.Sprintf("%T", manifestBaseDir("/.core/view.yaml"))
	core.AssertNotEmpty(t, observedType)
}

func TestManifest_SafeManifestRelativePath_Good(t *core.T) {
	root := t.TempDir()
	target := filepath.Join(root, "preload.js")
	core.RequireNoError(t, coreio.Local.WriteMode(target, "globalThis.ready = true;", 0o644))
	expected, err := filepath.EvalSymlinks(target)
	core.RequireNoError(t, err)

	got, err := safeManifestRelativePath(root, "preload.js", "preload path")

	core.RequireNoError(t, err)
	core.AssertEqual(t, expected, got)
}

func TestManifest_SafeManifestRelativePath_Bad(t *core.T) {
	root := t.TempDir()

	_, err := safeManifestRelativePath(root, "", "preload path")
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "empty")
}

func TestManifest_SafeManifestRelativePath_Bad_MissingFile(t *core.T) {
	root := t.TempDir()

	_, err := safeManifestRelativePath(root, "missing.js", "preload path")

	core.AssertError(t, err)
}

func TestManifest_SafeManifestRelativePath_Ugly(t *core.T) {
	root := t.TempDir()

	_, err := safeManifestRelativePath(root, "../escape.js", "preload path")
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "escapes")
}
