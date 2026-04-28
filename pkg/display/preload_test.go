package display

import (
	"context"
	"path/filepath"
	"strings"

	core "dappco.re/go"
	"dappco.re/go/gui/pkg/chat"
	"dappco.re/go/gui/pkg/window"
	coreio "dappco.re/go/io"
)

func TestDisplay_Good_WindowOpenIncludesPreload(t *core.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("preload"),
				window.WithURL("https://example.com"),
			},
		}},
	))
	core.RequireTrue(t, result.OK)
	core.AssertLen(t, platform.Windows, 1)
	core.AssertNotEmpty(t, platform.Windows[0].ExecJSCalls())
	core.AssertContains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.core.ml")
	core.AssertContains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.core.storage.cookies")
	core.AssertContains(t, platform.Windows[0].ExecJSCalls()[0], "Document.prototype, 'cookie'")
	core.AssertNotContains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.electron")
	core.AssertNotContains(t, platform.Windows[0].ExecJSCalls()[0], "core.background.serviceWorker.register")
}

func TestPreload_Good_TrustedOriginIncludesPrivilegedBridge(t *core.T) {
	svc, err := New()
	core.RequireNoError(t, err)

	script, err := svc.BuildPreloadScriptWithTrustedOriginPolicy(
		"core://app/",
		NewTrustedOriginPolicy([]string{"core://app/"}),
	)
	core.RequireNoError(t, err)

	core.AssertContains(t, script, "globalThis.electron")
	core.AssertContains(t, script, "core.background.serviceWorker.register")
	core.AssertContains(t, script, "globalThis.core.ml")
	core.AssertContains(t, script, "gui.notification.requestPermission")
	core.AssertContains(t, script, "gui.notification.clear")
	core.AssertContains(t, script, "systray.showMessage")
	core.AssertContains(t, script, "webview.devtoolsOpen")
}

func TestDisplay_Good_WindowOpenManifestBackedOriginIncludesManifestPreloadOnly(t *core.T) {
	home := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(home, ".core", "apps", "example.com", ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(home, ".core", "apps", "example.com", "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(home, ".core", "apps", "example.com", ".core", "view.yaml"), "name: example\npreloads:\n  - path: preload.js\n", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(
		filepath.Join(home, ".core", "preload-origins.yaml"),
		"origins:\n  - https://example.com/\n",
		0o644,
	))
	t.Setenv("DIR_HOME", home)

	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("manifest-backed"),
				window.WithURL("https://example.com/app"),
			},
		}},
	))
	core.RequireTrue(t, result.OK)
	core.AssertLen(t, platform.Windows, 1)
	script := platform.Windows[0].ExecJSCalls()[0]
	core.AssertContains(t, script, "__manifestLoaded")
	core.AssertContains(t, script, "globalThis.core.ml")
	core.AssertNotContains(t, script, "globalThis.electron")
	core.AssertNotContains(t, script, "core.background.serviceWorker.register")
}

func TestDisplay_Good_CoreSchemeRoutesThroughBackend(t *core.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(chat.Register(func(o *chat.Options) { o.StorePath = filepath.Join(t.TempDir(), "chat.db") })),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	core.RequireTrue(t, c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{window.WithName("settings")},
		}},
	)).OK)

	core.RequireTrue(t, c.Action("window.setURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetURL{Name: "settings", URL: "core://settings"}},
	)).OK)

	core.AssertLen(t, platform.Windows, 1)
	core.AssertTrue(t, strings.Contains(platform.Windows[0].HTMLContent(), "core://settings"))
}

func TestPreload_ValidatedLocalMLAPIURL_Good(t *core.T) {
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("http://localhost:8090/"))
	core.AssertEqual(t, "https://127.0.0.1:9443", validatedLocalMLAPIURL("https://127.0.0.1:9443/"))
	core.AssertNotEmpty(t, core.Sprintf("%T", validatedLocalMLAPIURL("http://localhost:8090/")))
}

func TestPreload_ValidatedLocalMLAPIURL_Bad(t *core.T) {
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("https://example.com"))
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("ftp://localhost:8090"))
	core.AssertNotEmpty(t, core.Sprintf("%T", validatedLocalMLAPIURL("https://example.com")))
}

func TestPreload_ValidatedLocalMLAPIURL_Ugly(t *core.T) {
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL(""))
	core.AssertEqual(t, "http://localhost:8090", validatedLocalMLAPIURL("not a url"))
	core.AssertNotEmpty(t, core.Sprintf("%T", validatedLocalMLAPIURL("")))
}

func TestPreload_TrustedPreloadOrigin_Good(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/"})

	core.AssertTrue(t, trustedPreloadOrigin("core://lab.lthn.sh/page", policy))
	core.AssertNotEmpty(t, core.Sprintf("%T", policy))
}

func TestPreload_TrustedPreloadOrigin_Bad(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/"})

	core.AssertFalse(t, trustedPreloadOrigin("core://attacker.com/x", policy))
	core.AssertFalse(t, trustedPreloadOrigin("wails://lab.lthn.sh/x", policy))
	core.AssertFalse(t, trustedPreloadOrigin("https://example.com", policy))
	core.AssertFalse(t, trustedPreloadOrigin("http://localhost:3000", policy))
	core.AssertFalse(t, trustedPreloadOrigin("file:///tmp/app/index.html", policy))
}

func TestPreload_TrustedPreloadOrigin_EmptyAllowListDeniesSchemeURLs(t *core.T) {
	policy := NewTrustedOriginPolicy(nil)

	core.AssertFalse(t, trustedPreloadOrigin("core://lab.lthn.sh/page", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://app/", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://attacker.com/x", policy))
}

func TestPreload_TrustedPreloadOrigin_PathPrefix(t *core.T) {
	policy := NewTrustedOriginPolicy([]string{"core://lab.lthn.sh/x"})

	core.AssertTrue(t, trustedPreloadOrigin("core://lab.lthn.sh/x/y", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://lab.lthn.sh/y", policy))
}

func TestPreload_BridgeActionAllowList(t *core.T) {
	policy := NewTrustedOriginPolicyWithActions(map[string][]string{
		"core://lab.lthn.sh/": {"display.sidecar.eval"},
		"core://empty/":       {},
	})

	core.AssertFalse(t, policy.AllowsActionURL("core://lab.lthn.sh/page", "marketplace.install"))
	core.AssertTrue(t, policy.AllowsActionURL("core://lab.lthn.sh/page", "display.sidecar.eval"))
	core.AssertFalse(t, policy.AllowsActionURL("core://attacker.com/page", "display.sidecar.eval"))
	core.AssertTrue(t, policy.AllowsURL("core://empty/page"))
	core.AssertFalse(t, policy.AllowsActionURL("core://empty/page", "display.sidecar.eval"))
}

func TestPreload_BridgeActionGuardScript(t *core.T) {
	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicyWithActions(map[string][]string{
		"core://lab.lthn.sh/": {"display.sidecar.eval"},
	})

	script, err := svc.BuildPreloadScriptWithTrustedOriginPolicy("core://lab.lthn.sh/page", policy)
	core.RequireNoError(t, err)

	core.AssertContains(t, script, "Core bridge action not permitted for this origin")
	core.AssertContains(t, script, `"display.sidecar.eval"`)
	core.AssertNotContains(t, script, `"marketplace.install"`)
}

func TestPreload_ManifestBackedPreloadOrigin_EmptyAllowListDeniesPlantedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "attacker.com")
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)

	core.AssertFalse(t, svc.manifestBackedPreloadOrigin(
		"https://attacker.com/app",
		NewTrustedOriginPolicy(nil),
	))
}

func TestPreload_ManifestBackedPreloadOrigin_AllowsListedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "lab.lthn.sh")
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertTrue(t, svc.manifestBackedPreloadOrigin("https://lab.lthn.sh/app", policy))
}

func TestPreload_ManifestBackedPreloadOrigin_DeniesUnlistedHTTPSManifest(t *core.T) {
	home := t.TempDir()
	writeMarketplaceViewManifest(t, home, "attacker.com")
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertFalse(t, svc.manifestBackedPreloadOrigin("https://attacker.com/app", policy))
}

func TestPreload_ManifestBackedPreloadOrigin_DeniesListedHTTPSOriginWithoutManifest(t *core.T) {
	home := t.TempDir()
	t.Setenv("DIR_HOME", home)

	svc, err := New()
	core.RequireNoError(t, err)
	policy := NewTrustedOriginPolicy([]string{"https://lab.lthn.sh/"})

	core.AssertFalse(t, svc.manifestBackedPreloadOrigin("https://lab.lthn.sh/app", policy))
}

func TestPreload_DefaultTrustedOriginPolicy_LoadsConfig(t *core.T) {
	home := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(home, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(
		filepath.Join(home, ".core", "preload-origins.yaml"),
		"origins:\n  - core://app/\n",
		0o644,
	))
	t.Setenv("DIR_HOME", home)
	t.Setenv("HOME", home)

	policy := DefaultTrustedOriginPolicy()

	core.AssertTrue(t, trustedPreloadOrigin("core://app/shell", policy))
	core.AssertFalse(t, trustedPreloadOrigin("core://attacker.com/shell", policy))
}

type preloadCapture struct {
	scripts []string
}

func (p *preloadCapture) ExecJS(script string) {
	p.scripts = append(p.scripts, script)
}

func writeMarketplaceViewManifest(t *core.T, home, host string) {
	t.Helper()
	dir := filepath.Join(home, ".core", "apps", host, ".core")
	core.RequireNoError(t, coreio.Local.EnsureDir(dir))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(dir, "view.yaml"), "name: "+host+"\n", 0o644))
}

func TestPreload_InjectPreload_Good(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "preload.js"), "globalThis.__manifestLoaded = true;", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "preloads:\n  - path: preload.js\n", 0o644))

	svc, err := New()
	core.RequireNoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "file://"+filepath.ToSlash(filepath.Join(root, "index.html")))
	core.RequireNoError(t, err)
	core.AssertLen(t, target.scripts, 1)
	core.AssertContains(t, target.scripts[0], "globalThis.core.ml")
	core.AssertNotContains(t, target.scripts[0], "globalThis.electron")
	core.AssertContains(t, target.scripts[0], "__manifestLoaded")
}

func TestPreload_InjectPreload_Bad(t *core.T) {
	svc, err := New()
	core.RequireNoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "https://example.com/app")
	core.RequireNoError(t, err)
	core.AssertLen(t, target.scripts, 1)
	core.AssertContains(t, target.scripts[0], "globalThis.core.ml")
	core.AssertNotContains(t, target.scripts[0], "globalThis.electron")
	core.AssertNotContains(t, target.scripts[0], "core.background.serviceWorker.register")
}

func TestPreload_InjectPreload_Ugly(t *core.T) {
	root := t.TempDir()
	core.RequireNoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	core.RequireNoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), "preloads: [\n", 0o644))

	svc, err := New()
	core.RequireNoError(t, err)
	target := &preloadCapture{}

	err = svc.InjectPreload(target, "file://"+filepath.ToSlash(filepath.Join(root, "index.html")))
	core.AssertError(t, err)
	core.AssertEmpty(t, target.scripts)
}
