package display

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/chat"
	"dappco.re/go/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisplay_Good_WindowOpenIncludesPreload(t *testing.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("preload"),
				window.WithURL("https://example.com"),
			},
		}},
	))
	require.True(t, result.OK)
	require.Len(t, platform.Windows, 1)
	assert.NotEmpty(t, platform.Windows[0].ExecJSCalls())
	assert.Contains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.core.ml")
	assert.Contains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.core.storage.cookies")
	assert.Contains(t, platform.Windows[0].ExecJSCalls()[0], "Document.prototype, 'cookie'")
	assert.NotContains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.electron")
	assert.NotContains(t, platform.Windows[0].ExecJSCalls()[0], "core.background.serviceWorker.register")
}

func TestDisplay_Good_WindowOpenTrustedOriginIncludesPrivilegedBridge(t *testing.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("preload"),
				window.WithURL("http://localhost:3000"),
			},
		}},
	))
	require.True(t, result.OK)
	require.Len(t, platform.Windows, 1)
	script := platform.Windows[0].ExecJSCalls()[0]
	assert.Contains(t, script, "globalThis.electron")
	assert.Contains(t, script, "core.background.serviceWorker.register")
	assert.Contains(t, script, "globalThis.core.ml")
	assert.Contains(t, script, "gui.notification.requestPermission")
	assert.Contains(t, script, "gui.notification.clear")
	assert.Contains(t, script, "systray.showMessage")
	assert.Contains(t, script, "webview.devtoolsOpen")
}

func TestDisplay_Good_WindowOpenManifestBackedOriginIncludesPrivilegedBridge(t *testing.T) {
	home := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(home, ".core", "apps", "example.com", ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(home, ".core", "apps", "example.com", ".core", "view.yaml"), []byte("name: example\n"), 0o644))
	t.Setenv("DIR_HOME", home)

	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("manifest-backed"),
				window.WithURL("https://example.com/app"),
			},
		}},
	))
	require.True(t, result.OK)
	require.Len(t, platform.Windows, 1)
	script := platform.Windows[0].ExecJSCalls()[0]
	assert.Contains(t, script, "globalThis.electron")
	assert.Contains(t, script, "core.background.serviceWorker.register")
	assert.Contains(t, script, "globalThis.core.ml")
}

func TestDisplay_Good_CoreSchemeRoutesThroughBackend(t *testing.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(chat.Register(func(o *chat.Options) { o.StorePath = filepath.Join(t.TempDir(), "chat.db") })),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	require.True(t, c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{window.WithName("settings")},
		}},
	)).OK)

	require.True(t, c.Action("window.setURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetURL{Name: "settings", URL: "core://settings"}},
	)).OK)

	require.Len(t, platform.Windows, 1)
	assert.True(t, strings.Contains(platform.Windows[0].HTMLContent(), "core://settings"))
}

func TestPreload_ValidatedLocalMLAPIURL_Good(t *testing.T) {
	assert.Equal(t, "http://localhost:8090", validatedLocalMLAPIURL("http://localhost:8090/"))
	assert.Equal(t, "https://127.0.0.1:9443", validatedLocalMLAPIURL("https://127.0.0.1:9443/"))
}

func TestPreload_ValidatedLocalMLAPIURL_Bad(t *testing.T) {
	assert.Equal(t, "http://localhost:8090", validatedLocalMLAPIURL("https://example.com"))
	assert.Equal(t, "http://localhost:8090", validatedLocalMLAPIURL("ftp://localhost:8090"))
}

func TestPreload_ValidatedLocalMLAPIURL_Ugly(t *testing.T) {
	assert.Equal(t, "http://localhost:8090", validatedLocalMLAPIURL(""))
	assert.Equal(t, "http://localhost:8090", validatedLocalMLAPIURL("not a url"))
}

func TestPreload_TrustedPreloadOrigin_Good(t *testing.T) {
	assert.True(t, trustedPreloadOrigin("core://store"))
	assert.True(t, trustedPreloadOrigin("http://localhost:3000"))
	assert.True(t, trustedPreloadOrigin("https://127.0.0.1:8443"))
}

func TestPreload_TrustedPreloadOrigin_Bad(t *testing.T) {
	assert.False(t, trustedPreloadOrigin("https://example.com"))
	assert.False(t, trustedPreloadOrigin("http://10.0.0.1:3000"))
	assert.False(t, trustedPreloadOrigin("file:///tmp/app/index.html"))
}

type preloadCapture struct {
	scripts []string
}

func (p *preloadCapture) ExecJS(script string) {
	p.scripts = append(p.scripts, script)
}

func TestPreload_InjectPreload_Good(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "preload.js"), []byte("globalThis.__manifestLoaded = true;"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte("preloads:\n  - path: preload.js\n"), 0o644))

	svc, err := New()
	require.NoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "file://"+filepath.ToSlash(filepath.Join(root, "index.html")))
	require.NoError(t, err)
	require.Len(t, target.scripts, 1)
	assert.Contains(t, target.scripts[0], "globalThis.core.ml")
	assert.Contains(t, target.scripts[0], "globalThis.electron")
	assert.Contains(t, target.scripts[0], "__manifestLoaded")
}

func TestPreload_InjectPreload_Bad(t *testing.T) {
	svc, err := New()
	require.NoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "https://example.com/app")
	require.NoError(t, err)
	require.Len(t, target.scripts, 1)
	assert.Contains(t, target.scripts[0], "globalThis.core.ml")
	assert.NotContains(t, target.scripts[0], "globalThis.electron")
	assert.NotContains(t, target.scripts[0], "core.background.serviceWorker.register")
}

func TestPreload_InjectPreload_Ugly(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte("preloads: [\n"), 0o644))

	svc, err := New()
	require.NoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "file://"+filepath.ToSlash(filepath.Join(root, "index.html")))
	require.Error(t, err)
	assert.Empty(t, target.scripts)
}
