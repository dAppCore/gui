package display

import (
	"context"
	"os"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/menu"
	"forge.lthn.ai/core/gui/pkg/systray"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test helpers ---

// newTestDisplayService creates a display service registered with Core for IPC testing.
func newTestDisplayService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "display")
	return svc, c
}

// newTestConclave creates a full 4-service conclave for integration testing.
func newTestConclave(t *testing.T) *core.Core {
	t.Helper()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithService(systray.Register(systray.NewMockPlatform())),
		core.WithService(menu.Register(menu.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	return c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

// --- Tests ---

func TestNew_Good(t *testing.T) {
	service, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestNew_Good_IndependentInstances(t *testing.T) {
	service1, err1 := New()
	service2, err2 := New()
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotSame(t, service1, service2)
}

func TestRegister_Good(t *testing.T) {
	factory := Register(nil) // nil wailsApp for testing
	assert.NotNil(t, factory)

	c := core.New(
		core.WithService(factory),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	svc := core.MustServiceFor[*Service](c, "display")
	assert.NotNil(t, svc)
}

func TestConfigQuery_Good(t *testing.T) {
	svc, c := newTestDisplayService(t)

	// Set window config
	svc.configData["window"] = map[string]any{
		"default_width": 1024,
	}

	r := c.QUERY(window.QueryConfig{})
	require.True(t, r.OK)
	cfg := r.Value.(map[string]any)
	assert.Equal(t, 1024, cfg["default_width"])
}

func TestConfigQuery_Bad(t *testing.T) {
	// No display service — window config query returns handled=false
	c := core.New(core.WithServiceLock())
	r := c.QUERY(window.QueryConfig{})
	assert.False(t, r.OK)
}

func TestConfigTask_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	newCfg := map[string]any{"default_width": 800}
	r := taskRun(c, "display.saveWindowConfig", window.TaskSaveConfig{Config: newCfg})
	require.True(t, r.OK)

	// Verify config was saved
	r2 := c.QUERY(window.QueryConfig{})
	cfg := r2.Value.(map[string]any)
	assert.Equal(t, 800, cfg["default_width"])
}

func TestResolveScheme_StoreRoute_Good(t *testing.T) {
	svc, _ := newTestDisplayService(t)

	result := svc.ResolveScheme(context.Background(), "core://store?q=alpha")
	require.True(t, result.OK)

	payload, ok := result.Value.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "text/html", payload["content_type"])

	body, ok := payload["body"].(string)
	require.True(t, ok)
	assert.Contains(t, body, "core://store")
	assert.Contains(t, body, "storage scopes")
	assert.Contains(t, body, "Search the in-memory storage scopes")
}

// --- Conclave integration tests ---

func TestServiceConclave_Good(t *testing.T) {
	c := newTestConclave(t)

	// Open a window via IPC
	r := taskRun(c, "window.open", window.TaskOpenWindow{
		Window: &window.Window{Name: "main"},
	})
	require.True(t, r.OK)
	info := r.Value.(window.WindowInfo)
	assert.Equal(t, "main", info.Name)

	// Query window config from display
	r2 := c.QUERY(window.QueryConfig{})
	require.True(t, r2.OK)
	assert.NotNil(t, r2.Value)

	// Set app menu via IPC
	r3 := taskRun(c, "menu.setAppMenu", menu.TaskSetAppMenu{Items: []menu.MenuItem{
		{Label: "File"},
	}})
	require.True(t, r3.OK)

	// Query app menu via IPC
	r4 := c.QUERY(menu.QueryGetAppMenu{})
	assert.True(t, r4.OK)
	items := r4.Value.([]menu.MenuItem)
	assert.Len(t, items, 1)
}

func TestServiceConclave_Bad(t *testing.T) {
	// Sub-service starts without display — config QUERY returns handled=false
	c := core.New(
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	r := c.QUERY(window.QueryConfig{})
	assert.False(t, r.OK, "no display service means no config handler")
}

// --- IPC delegation tests (full conclave) ---

func TestOpenWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	t.Run("creates window with default options", func(t *testing.T) {
		err := svc.OpenWindow()
		assert.NoError(t, err)

		// Verify via IPC query
		infos := svc.ListWindowInfos()
		assert.GreaterOrEqual(t, len(infos), 1)
	})

	t.Run("creates window with custom options", func(t *testing.T) {
		err := svc.OpenWindow(
			window.WithName("custom-window"),
			window.WithTitle("Custom Title"),
			window.WithSize(640, 480),
			window.WithURL("/custom"),
		)
		assert.NoError(t, err)

		r := c.QUERY(window.QueryWindowByName{Name: "custom-window"})
		info := r.Value.(*window.WindowInfo)
		assert.Equal(t, "custom-window", info.Name)
	})
}

func TestGetWindowInfo_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	_ = svc.OpenWindow(
		window.WithName("test-win"),
		window.WithSize(800, 600),
	)

	// Modify position via IPC
	taskRun(c, "window.setPosition", window.TaskSetPosition{Name: "test-win", X: 100, Y: 200})

	info, err := svc.GetWindowInfo("test-win")
	require.NoError(t, err)
	assert.Equal(t, "test-win", info.Name)
	assert.Equal(t, 100, info.X)
	assert.Equal(t, 200, info.Y)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestGetWindowInfo_Bad(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	info, err := svc.GetWindowInfo("nonexistent")
	// QueryWindowByName returns nil for nonexistent — handled=true, result=nil
	assert.NoError(t, err)
	assert.Nil(t, info)
}

func TestListWindowInfos_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	_ = svc.OpenWindow(window.WithName("win-1"))
	_ = svc.OpenWindow(window.WithName("win-2"))

	infos := svc.ListWindowInfos()
	assert.Len(t, infos, 2)
}

func TestSetWindowPosition_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("pos-win"))

	err := svc.SetWindowPosition("pos-win", 300, 400)
	assert.NoError(t, err)

	info, _ := svc.GetWindowInfo("pos-win")
	assert.Equal(t, 300, info.X)
	assert.Equal(t, 400, info.Y)
}

func TestSetWindowPosition_Bad(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	err := svc.SetWindowPosition("nonexistent", 0, 0)
	assert.Error(t, err)
}

func TestSetWindowSize_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("size-win"))

	err := svc.SetWindowSize("size-win", 1024, 768)
	assert.NoError(t, err)

	info, _ := svc.GetWindowInfo("size-win")
	assert.Equal(t, 1024, info.Width)
	assert.Equal(t, 768, info.Height)
}

func TestMaximizeWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("max-win"))

	err := svc.MaximizeWindow("max-win")
	assert.NoError(t, err)

	info, _ := svc.GetWindowInfo("max-win")
	assert.True(t, info.Maximized)
}

func TestRestoreWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("restore-win"))
	_ = svc.MaximizeWindow("restore-win")

	err := svc.RestoreWindow("restore-win")
	assert.NoError(t, err)

	info, _ := svc.GetWindowInfo("restore-win")
	assert.False(t, info.Maximized)
}

func TestFocusWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("focus-win"))

	err := svc.FocusWindow("focus-win")
	assert.NoError(t, err)

	info, _ := svc.GetWindowInfo("focus-win")
	assert.True(t, info.Focused)
}

func TestCloseWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("close-win"))

	err := svc.CloseWindow("close-win")
	assert.NoError(t, err)

	// Window should be removed
	info, _ := svc.GetWindowInfo("close-win")
	assert.Nil(t, info)
}

func TestSetWindowVisibility_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("vis-win"))

	err := svc.SetWindowVisibility("vis-win", false)
	assert.NoError(t, err)

	err = svc.SetWindowVisibility("vis-win", true)
	assert.NoError(t, err)
}

func TestSetWindowAlwaysOnTop_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("ontop-win"))

	err := svc.SetWindowAlwaysOnTop("ontop-win", true)
	assert.NoError(t, err)
}

func TestSetWindowTitle_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("title-win"))

	err := svc.SetWindowTitle("title-win", "New Title")
	assert.NoError(t, err)
}

func TestGetFocusedWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("win-a"))
	_ = svc.OpenWindow(window.WithName("win-b"))
	_ = svc.FocusWindow("win-b")

	focused := svc.GetFocusedWindow()
	assert.Equal(t, "win-b", focused)
}

func TestGetFocusedWindow_Good_NoneSelected(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("win-a"))

	focused := svc.GetFocusedWindow()
	assert.Equal(t, "", focused)
}

func TestCreateWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	info, err := svc.CreateWindow(CreateWindowOptions{
		Name:   "new-win",
		Title:  "New Window",
		URL:    "/new",
		Width:  600,
		Height: 400,
	})
	require.NoError(t, err)
	assert.Equal(t, "new-win", info.Name)
}

func TestCreateWindow_Bad(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	_, err := svc.CreateWindow(CreateWindowOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "window name is required")
}

func TestResetWindowState_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	err := svc.ResetWindowState()
	assert.NoError(t, err)
}

func TestGetSavedWindowStates_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	states := svc.GetSavedWindowStates()
	assert.NotNil(t, states)
}

func TestHandleIPCEvents_WindowOpened_Good(t *testing.T) {
	c := newTestConclave(t)

	// Open a window — this should trigger ActionWindowOpened
	// which HandleIPCEvents should convert to a WS event
	r := taskRun(c, "window.open", window.TaskOpenWindow{
		Window: &window.Window{Name: "test"},
	})
	require.True(t, r.OK)
	info := r.Value.(window.WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestHandleListWorkspaces_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	// handleListWorkspaces should not panic when workspace service is not available
	assert.NotPanics(t, func() {
		svc.handleListWorkspaces()
	})
}

func TestWSEventManager_Good(t *testing.T) {
	em := NewWSEventManager()
	defer em.Close()

	assert.NotNil(t, em)
	assert.Equal(t, 0, em.ConnectedClients())
}

// --- Config file loading tests ---

func TestLoadConfig_Good(t *testing.T) {
	// Create temp config file
	dir := t.TempDir()
	cfgPath := core.JoinPath(dir, ".core", "gui", "config.yaml")
	require.NoError(t, os.MkdirAll(core.PathDir(cfgPath), 0o755))
	require.NoError(t, os.WriteFile(cfgPath, []byte(`
window:
  default_width: 1280
  default_height: 720
systray:
  tooltip: "Test App"
menu:
  show_dev_tools: false
`), 0o644))

	s, _ := New()
	s.loadConfigFrom(cfgPath)

	// Verify configData was populated from file
	assert.Equal(t, 1280, s.configData["window"]["default_width"])
	assert.Equal(t, "Test App", s.configData["systray"]["tooltip"])
	assert.Equal(t, false, s.configData["menu"]["show_dev_tools"])
}

func TestLoadConfig_Bad_MissingFile(t *testing.T) {
	s, _ := New()
	s.loadConfigFrom(core.JoinPath(t.TempDir(), "nonexistent.yaml"))

	// Should not panic, configData stays at empty defaults
	assert.Empty(t, s.configData["window"])
	assert.Empty(t, s.configData["systray"])
	assert.Empty(t, s.configData["menu"])
}

func TestHandleConfigTask_Persists_Good(t *testing.T) {
	dir := t.TempDir()
	cfgPath := core.JoinPath(dir, "config.yaml")

	s, _ := New()
	s.loadConfigFrom(cfgPath) // Creates empty config (file doesn't exist yet)

	// Simulate a TaskSaveConfig through the handler
	c := core.New(
		core.WithService(func(c *core.Core) core.Result {
			s.ServiceRuntime = core.NewServiceRuntime[Options](c, Options{})
			return core.Result{Value: s, OK: true}
		}),
		core.WithServiceLock(),
	)
	c.ServiceStartup(context.Background(), nil)

	r := taskRun(c, "display.saveWindowConfig", window.TaskSaveConfig{
		Config: map[string]any{"default_width": 1920},
	})
	require.True(t, r.OK)

	// Verify file was written
	data, err := os.ReadFile(cfgPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "default_width")
}
