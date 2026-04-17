package display

import (
	"context"
	"os"
	"reflect"
	"strings"
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

func TestStorageTask_Bad(t *testing.T) {
	_, c := newTestDisplayService(t)

	r := c.Action("display.storage.set").Run(context.Background(), core.NewOptions(
		core.Option{Key: "origin", Value: "core://settings"},
		core.Option{Key: "bucket", Value: "localStorage"},
		core.Option{Key: "key", Value: strings.Repeat("k", maxStorageKeyBytes+1)},
		core.Option{Key: "value", Value: "dark"},
	))

	require.False(t, r.OK)
	assert.Contains(t, r.Value.(error).Error(), "invalid storage entry")
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

func TestGetWindowInfo_BadType(t *testing.T) {
	svc, c := newTestDisplayService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowByName:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	info, err := svc.GetWindowInfo("broken")

	require.Error(t, err)
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

func TestSetWindowPosition_ActionFailure(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	c.Action("window.setPosition", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: false}
	})

	err := svc.SetWindowPosition("pos-win", 300, 400)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "window.setPosition")
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

func TestSetWindowBounds_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("bounds-win"))

	err := svc.SetWindowBounds("bounds-win", 10, 20, 640, 480)
	assert.NoError(t, err)

	info, _ := svc.GetWindowInfo("bounds-win")
	assert.Equal(t, 10, info.X)
	assert.Equal(t, 20, info.Y)
	assert.Equal(t, 640, info.Width)
	assert.Equal(t, 480, info.Height)
}

func TestSetWindowBounds_Bad(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	err := svc.SetWindowBounds("missing", 1, 2, 3, 4)

	assert.Error(t, err)
}

func TestSetWindowBounds_Ugly(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("bounds-win"))

	err := svc.SetWindowBounds("bounds-win", -10, -20, 0, 1)
	assert.NoError(t, err)

	info, _ := svc.GetWindowInfo("bounds-win")
	assert.Equal(t, -10, info.X)
	assert.Equal(t, -20, info.Y)
	assert.Equal(t, 0, info.Width)
	assert.Equal(t, 1, info.Height)
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

func TestSetWindowFullscreen_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	windowSvc := core.MustServiceFor[*window.Service](c, "window")
	_ = svc.OpenWindow(window.WithName("full-win"))

	err := svc.SetWindowFullscreen("full-win", true)

	require.NoError(t, err)
	pw, ok := windowSvc.Manager().Get("full-win")
	require.True(t, ok)
	assert.True(t, pw.IsFullscreen())
}

func TestSetWindowFullscreen_Bad(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	err := svc.SetWindowFullscreen("missing", true)

	require.Error(t, err)
}

func TestLayoutBesideEditor_ActionFailure(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	c.Action("window.layoutBesideEditor", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: false}
	})

	result, err := svc.LayoutBesideEditor("preview", "code", "right", 0.62)

	require.Error(t, err)
	assert.Zero(t, result)
	assert.Contains(t, err.Error(), "window.layoutBesideEditor")
}

func TestSetWindowFullscreen_Ugly(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	windowSvc := core.MustServiceFor[*window.Service](c, "window")
	_ = svc.OpenWindow(window.WithName("full-win"))

	require.NoError(t, svc.SetWindowFullscreen("full-win", true))
	err := svc.SetWindowFullscreen("full-win", false)

	require.NoError(t, err)
	pw, ok := windowSvc.Manager().Get("full-win")
	require.True(t, ok)
	assert.False(t, pw.IsFullscreen())
}

func TestGetWindowTitle_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("title-win"), window.WithTitle("Inspector"))

	title, err := svc.GetWindowTitle("title-win")

	require.NoError(t, err)
	assert.Equal(t, "Inspector", title)
}

func TestGetWindowTitle_Bad(t *testing.T) {
	svc, _ := newTestDisplayService(t)

	title, err := svc.GetWindowTitle("missing")

	require.Error(t, err)
	assert.Empty(t, title)
}

func TestGetWindowTitle_Ugly(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("title-win"), window.WithTitle("Line 1 <Line 2>\nTabbed"))

	title, err := svc.GetWindowTitle("title-win")

	require.NoError(t, err)
	assert.Equal(t, "Line 1 <Line 2>\nTabbed", title)
}

func TestMinimizeWindow_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	windowSvc := core.MustServiceFor[*window.Service](c, "window")
	_ = svc.OpenWindow(window.WithName("min-win"))

	err := svc.MinimizeWindow("min-win")

	require.NoError(t, err)
	pw, ok := windowSvc.Manager().Get("min-win")
	require.True(t, ok)
	assert.True(t, pw.IsMinimised())
}

func TestMinimizeWindow_Bad(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")

	err := svc.MinimizeWindow("missing")

	require.Error(t, err)
}

func TestMinimizeWindow_Ugly(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	windowSvc := core.MustServiceFor[*window.Service](c, "window")
	_ = svc.OpenWindow(window.WithName("min-win"))

	require.NoError(t, svc.MinimizeWindow("min-win"))
	err := svc.MinimizeWindow("min-win")

	require.NoError(t, err)
	pw, ok := windowSvc.Manager().Get("min-win")
	require.True(t, ok)
	assert.True(t, pw.IsMinimised())
}

func TestHandleWSMessage_SetWindowOpacity_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("opacity-win"))

	r := svc.handleWSMessage(WSMessage{
		Action: "window:set-opacity",
		Data: map[string]any{
			"name":    "opacity-win",
			"opacity": 0.35,
		},
	})
	require.True(t, r.OK)

	info, err := svc.GetWindowInfo("opacity-win")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.InDelta(t, 0.35, info.Opacity, 0.0001)
}

func TestDisplay_requireStringField_Good(t *testing.T) {
	value, err := requireStringField(map[string]any{"window": "main"}, "window")

	require.NoError(t, err)
	assert.Equal(t, "main", value)
}

func TestDisplay_requireStringField_Bad(t *testing.T) {
	value, err := requireStringField(map[string]any{"window": ""}, "window")

	require.Error(t, err)
	assert.Empty(t, value)
}

func TestDisplay_requireStringField_Ugly(t *testing.T) {
	value, err := requireStringField(map[string]any{"window": 42}, "window")

	require.Error(t, err)
	assert.Empty(t, value)
}

func TestDisplay_optionsFromMap_Good(t *testing.T) {
	opts := optionsFromMap(map[string]any{"alpha": "one", "beta": 2})

	require.Equal(t, 2, opts.Len())
	got := map[string]any{}
	for _, opt := range opts.Items() {
		got[opt.Key] = opt.Value
	}
	assert.True(t, reflect.DeepEqual(map[string]any{"alpha": "one", "beta": 2}, got))
}

func TestDisplay_optionsFromMap_Bad(t *testing.T) {
	opts := optionsFromMap(nil)

	require.NotNil(t, opts)
	assert.Equal(t, 0, opts.Len())
}

func TestDisplay_optionsFromMap_Ugly(t *testing.T) {
	opts := wsOptions(map[string]any{"nested": map[string]any{"value": "x"}})

	require.Equal(t, 1, opts.Len())
	item := opts.Items()[0]
	assert.Equal(t, "nested", item.Key)
	assert.Equal(t, map[string]any{"value": "x"}, item.Value)
}

func TestDisplay_handleWSMessage_Good(t *testing.T) {
	c := newTestConclave(t)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("opacity-win"))

	result := svc.handleWSMessage(WSMessage{
		Action: "window:set-opacity",
		Data: map[string]any{
			"name":    "opacity-win",
			"opacity": 0.55,
		},
	})
	require.True(t, result.OK)

	info, err := svc.GetWindowInfo("opacity-win")
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.InDelta(t, 0.55, info.Opacity, 0.0001)
}

func TestDisplay_handleWSMessage_Bad(t *testing.T) {
	svc, _ := newTestDisplayService(t)
	result := svc.handleWSMessage(WSMessage{Action: "unknown:action"})

	require.False(t, result.OK)
	assert.Contains(t, result.Value.(error).Error(), "unknown websocket action")
}

func TestDisplay_handleWSMessage_Ugly(t *testing.T) {
	svc, _ := newTestDisplayService(t)
	result := svc.handleWSMessage(WSMessage{
		Action: "window:set-opacity",
		Data: map[string]any{
			"name": "main",
		},
	})

	require.False(t, result.OK)
	assert.Contains(t, result.Value.(error).Error(), "missing required field \"opacity\"")
}

func TestDisplay_handleTrayAction_Good(t *testing.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithService(systray.Register(systray.NewMockPlatform())),
		core.WithService(menu.Register(menu.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "display")
	_ = svc.OpenWindow(window.WithName("one"))
	_ = svc.OpenWindow(window.WithName("two"))

	svc.handleTrayAction("open-desktop")
	require.Len(t, platform.Windows, 2)
	assert.True(t, platform.Windows[0].IsFocused())
	assert.True(t, platform.Windows[1].IsFocused())

	svc.handleTrayAction("close-desktop")
	assert.False(t, platform.Windows[0].IsVisible())
	assert.False(t, platform.Windows[1].IsVisible())
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
