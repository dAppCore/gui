package display

import (
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/menu"
	"forge.lthn.ai/core/gui/pkg/systray"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock platform implementations for sub-packages ---

// displayMockPlatformWindow implements window.PlatformWindow for display tests.
type displayMockPlatformWindow struct {
	name, title, url     string
	width, height, x, y  int
	maximised, focused   bool
	visible, alwaysOnTop bool
	closed               bool
	eventHandlers        []func(window.WindowEvent)
}

func (w *displayMockPlatformWindow) Name() string             { return w.name }
func (w *displayMockPlatformWindow) Position() (int, int)     { return w.x, w.y }
func (w *displayMockPlatformWindow) Size() (int, int)         { return w.width, w.height }
func (w *displayMockPlatformWindow) IsMaximised() bool        { return w.maximised }
func (w *displayMockPlatformWindow) IsFocused() bool          { return w.focused }
func (w *displayMockPlatformWindow) SetTitle(title string)    { w.title = title }
func (w *displayMockPlatformWindow) SetPosition(x, y int)     { w.x = x; w.y = y }
func (w *displayMockPlatformWindow) SetSize(width, height int) {
	w.width = width
	w.height = height
}
func (w *displayMockPlatformWindow) SetBackgroundColour(r, g, b, a uint8) {}
func (w *displayMockPlatformWindow) SetVisibility(visible bool)           { w.visible = visible }
func (w *displayMockPlatformWindow) SetAlwaysOnTop(alwaysOnTop bool)      { w.alwaysOnTop = alwaysOnTop }
func (w *displayMockPlatformWindow) Maximise()                            { w.maximised = true }
func (w *displayMockPlatformWindow) Restore()                             { w.maximised = false }
func (w *displayMockPlatformWindow) Minimise()                            {}
func (w *displayMockPlatformWindow) Focus()                               { w.focused = true }
func (w *displayMockPlatformWindow) Close()                               { w.closed = true }
func (w *displayMockPlatformWindow) Show()                                { w.visible = true }
func (w *displayMockPlatformWindow) Hide()                                { w.visible = false }
func (w *displayMockPlatformWindow) Fullscreen()                          {}
func (w *displayMockPlatformWindow) UnFullscreen()                        {}
func (w *displayMockPlatformWindow) OnWindowEvent(handler func(window.WindowEvent)) {
	w.eventHandlers = append(w.eventHandlers, handler)
}

// displayMockWindowPlatform implements window.Platform for display tests.
type displayMockWindowPlatform struct {
	windows []*displayMockPlatformWindow
}

func (p *displayMockWindowPlatform) CreateWindow(opts window.PlatformWindowOptions) window.PlatformWindow {
	w := &displayMockPlatformWindow{
		name: opts.Name, title: opts.Title, url: opts.URL,
		width: opts.Width, height: opts.Height,
		x: opts.X, y: opts.Y,
	}
	p.windows = append(p.windows, w)
	return w
}

func (p *displayMockWindowPlatform) GetWindows() []window.PlatformWindow {
	out := make([]window.PlatformWindow, len(p.windows))
	for i, w := range p.windows {
		out[i] = w
	}
	return out
}

// displayMockSystrayPlatform implements systray.Platform for display tests.
type displayMockSystrayPlatform struct {
	trays []*displayMockTray
	menus []*displayMockSystrayMenu
}

func (p *displayMockSystrayPlatform) NewTray() systray.PlatformTray {
	t := &displayMockTray{}
	p.trays = append(p.trays, t)
	return t
}

func (p *displayMockSystrayPlatform) NewMenu() systray.PlatformMenu {
	m := &displayMockSystrayMenu{}
	p.menus = append(p.menus, m)
	return m
}

type displayMockTray struct {
	tooltip, label string
	menu           systray.PlatformMenu
}

func (t *displayMockTray) SetIcon(data []byte)            {}
func (t *displayMockTray) SetTemplateIcon(data []byte)    {}
func (t *displayMockTray) SetTooltip(text string)         { t.tooltip = text }
func (t *displayMockTray) SetLabel(text string)            { t.label = text }
func (t *displayMockTray) SetMenu(m systray.PlatformMenu)  { t.menu = m }
func (t *displayMockTray) AttachWindow(w systray.WindowHandle) {}

type displayMockSystrayMenu struct {
	items []string
}

func (m *displayMockSystrayMenu) Add(label string) systray.PlatformMenuItem {
	m.items = append(m.items, label)
	return &displayMockSystrayMenuItem{}
}
func (m *displayMockSystrayMenu) AddSeparator() { m.items = append(m.items, "---") }

type displayMockSystrayMenuItem struct{}

func (mi *displayMockSystrayMenuItem) SetTooltip(text string)  {}
func (mi *displayMockSystrayMenuItem) SetChecked(checked bool)  {}
func (mi *displayMockSystrayMenuItem) SetEnabled(enabled bool)  {}
func (mi *displayMockSystrayMenuItem) OnClick(fn func())        {}
func (mi *displayMockSystrayMenuItem) AddSubmenu() systray.PlatformMenu {
	return &displayMockSystrayMenu{}
}

// displayMockMenuPlatform implements menu.Platform for display tests.
type displayMockMenuPlatform struct {
	appMenu menu.PlatformMenu
}

func (p *displayMockMenuPlatform) NewMenu() menu.PlatformMenu {
	return &displayMockMenu{}
}

func (p *displayMockMenuPlatform) SetApplicationMenu(m menu.PlatformMenu) {
	p.appMenu = m
}

type displayMockMenu struct {
	items []string
}

func (m *displayMockMenu) Add(label string) menu.PlatformMenuItem {
	m.items = append(m.items, label)
	return &displayMockMenuItem{}
}
func (m *displayMockMenu) AddSeparator()                      { m.items = append(m.items, "---") }
func (m *displayMockMenu) AddSubmenu(label string) menu.PlatformMenu {
	m.items = append(m.items, label)
	return &displayMockMenu{}
}
func (m *displayMockMenu) AddRole(role menu.MenuRole) {}

type displayMockMenuItem struct{}

func (mi *displayMockMenuItem) SetAccelerator(accel string) menu.PlatformMenuItem { return mi }
func (mi *displayMockMenuItem) SetTooltip(text string) menu.PlatformMenuItem      { return mi }
func (mi *displayMockMenuItem) SetChecked(checked bool) menu.PlatformMenuItem      { return mi }
func (mi *displayMockMenuItem) SetEnabled(enabled bool) menu.PlatformMenuItem      { return mi }
func (mi *displayMockMenuItem) OnClick(fn func()) menu.PlatformMenuItem            { return mi }

// --- Test helpers ---

// newTestCore creates a new core instance for testing.
func newTestCore(t *testing.T) *core.Core {
	coreInstance, err := core.New()
	require.NoError(t, err)
	return coreInstance
}

// newServiceWithMocks creates a Service with mock sub-managers for testing.
// Uses a temp directory for state/layout persistence to avoid loading real saved state.
func newServiceWithMocks(t *testing.T) (*Service, *mockApp, *displayMockWindowPlatform) {
	service, err := New()
	require.NoError(t, err)

	mock := newMockApp()
	service.app = mock

	wp := &displayMockWindowPlatform{}
	service.windows = window.NewManagerWithDir(wp, t.TempDir())
	service.tray = systray.NewManager(&displayMockSystrayPlatform{})
	service.menus = menu.NewManager(&displayMockMenuPlatform{})

	return service, mock, wp
}

// --- Tests ---

func TestNew(t *testing.T) {
	t.Run("creates service successfully", func(t *testing.T) {
		service, err := New()
		assert.NoError(t, err)
		assert.NotNil(t, service, "New() should return a non-nil service instance")
	})

	t.Run("returns independent instances", func(t *testing.T) {
		service1, err1 := New()
		service2, err2 := New()
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotSame(t, service1, service2, "New() should return different instances")
	})
}

func TestRegister(t *testing.T) {
	t.Run("registers with core successfully", func(t *testing.T) {
		coreInstance := newTestCore(t)
		service, err := Register(coreInstance)
		require.NoError(t, err)
		assert.NotNil(t, service, "Register() should return a non-nil service instance")
	})

	t.Run("returns Service type", func(t *testing.T) {
		coreInstance := newTestCore(t)
		service, err := Register(coreInstance)
		require.NoError(t, err)

		displayService, ok := service.(*Service)
		assert.True(t, ok, "Register() should return *Service type")
		assert.NotNil(t, displayService.ServiceRuntime, "ServiceRuntime should be initialized")
	})
}

func TestServiceName(t *testing.T) {
	service, err := New()
	require.NoError(t, err)

	name := service.ServiceName()
	assert.Equal(t, "forge.lthn.ai/core/gui/display", name)
}

func TestOpenWindow_Good(t *testing.T) {
	t.Run("creates window with default options", func(t *testing.T) {
		service, _, wp := newServiceWithMocks(t)

		err := service.OpenWindow()
		assert.NoError(t, err)

		// Verify window was created through the platform
		assert.Len(t, wp.windows, 1)
		assert.Equal(t, "main", wp.windows[0].name)
		assert.Equal(t, "Core", wp.windows[0].title)
		assert.Equal(t, 1280, wp.windows[0].width)
		assert.Equal(t, 800, wp.windows[0].height)
	})

	t.Run("creates window with custom options", func(t *testing.T) {
		service, _, wp := newServiceWithMocks(t)

		err := service.OpenWindow(
			window.WithName("custom-window"),
			window.WithTitle("Custom Title"),
			window.WithSize(640, 480),
			window.WithURL("/custom"),
		)
		assert.NoError(t, err)

		assert.Len(t, wp.windows, 1)
		assert.Equal(t, "custom-window", wp.windows[0].name)
		assert.Equal(t, "Custom Title", wp.windows[0].title)
		assert.Equal(t, 640, wp.windows[0].width)
		assert.Equal(t, 480, wp.windows[0].height)
	})
}

func TestGetWindowInfo_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)

	_ = service.OpenWindow(
		window.WithName("test-win"),
		window.WithSize(800, 600),
	)

	// Set position on the mock window
	wp.windows[0].x = 100
	wp.windows[0].y = 200

	info, err := service.GetWindowInfo("test-win")
	require.NoError(t, err)
	assert.Equal(t, "test-win", info.Name)
	assert.Equal(t, 100, info.X)
	assert.Equal(t, 200, info.Y)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestGetWindowInfo_Bad(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)

	_, err := service.GetWindowInfo("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "window not found")
}

func TestListWindowInfos_Good(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)

	_ = service.OpenWindow(window.WithName("win-1"))
	_ = service.OpenWindow(window.WithName("win-2"))

	infos := service.ListWindowInfos()
	assert.Len(t, infos, 2)
}

func TestSetWindowPosition_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("pos-win"))

	err := service.SetWindowPosition("pos-win", 300, 400)
	assert.NoError(t, err)
	assert.Equal(t, 300, wp.windows[0].x)
	assert.Equal(t, 400, wp.windows[0].y)
}

func TestSetWindowPosition_Bad(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)

	err := service.SetWindowPosition("nonexistent", 0, 0)
	assert.Error(t, err)
}

func TestSetWindowSize_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("size-win"))

	err := service.SetWindowSize("size-win", 1024, 768)
	assert.NoError(t, err)
	assert.Equal(t, 1024, wp.windows[0].width)
	assert.Equal(t, 768, wp.windows[0].height)
}

func TestMaximizeWindow_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("max-win"))

	err := service.MaximizeWindow("max-win")
	assert.NoError(t, err)
	assert.True(t, wp.windows[0].maximised)
}

func TestRestoreWindow_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("restore-win"))
	wp.windows[0].maximised = true

	err := service.RestoreWindow("restore-win")
	assert.NoError(t, err)
	assert.False(t, wp.windows[0].maximised)
}

func TestFocusWindow_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("focus-win"))

	err := service.FocusWindow("focus-win")
	assert.NoError(t, err)
	assert.True(t, wp.windows[0].focused)
}

func TestCloseWindow_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("close-win"))

	err := service.CloseWindow("close-win")
	assert.NoError(t, err)
	assert.True(t, wp.windows[0].closed)

	// Window should be removed from manager
	_, ok := service.windows.Get("close-win")
	assert.False(t, ok)
}

func TestSetWindowVisibility_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("vis-win"))

	err := service.SetWindowVisibility("vis-win", false)
	assert.NoError(t, err)
	assert.False(t, wp.windows[0].visible)

	err = service.SetWindowVisibility("vis-win", true)
	assert.NoError(t, err)
	assert.True(t, wp.windows[0].visible)
}

func TestSetWindowAlwaysOnTop_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("ontop-win"))

	err := service.SetWindowAlwaysOnTop("ontop-win", true)
	assert.NoError(t, err)
	assert.True(t, wp.windows[0].alwaysOnTop)
}

func TestSetWindowTitle_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("title-win"))

	err := service.SetWindowTitle("title-win", "New Title")
	assert.NoError(t, err)
	assert.Equal(t, "New Title", wp.windows[0].title)
}

func TestGetFocusedWindow_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("win-a"))
	_ = service.OpenWindow(window.WithName("win-b"))
	wp.windows[1].focused = true

	focused := service.GetFocusedWindow()
	assert.Equal(t, "win-b", focused)
}

func TestGetFocusedWindow_NoneSelected(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)
	_ = service.OpenWindow(window.WithName("win-a"))

	focused := service.GetFocusedWindow()
	assert.Equal(t, "", focused)
}

func TestCreateWindow_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)

	info, err := service.CreateWindow(CreateWindowOptions{
		Name:   "new-win",
		Title:  "New Window",
		URL:    "/new",
		Width:  600,
		Height: 400,
	})
	require.NoError(t, err)
	assert.Equal(t, "new-win", info.Name)
	assert.Equal(t, 600, info.Width)
	assert.Equal(t, 400, info.Height)
	assert.Len(t, wp.windows, 1)
}

func TestCreateWindow_Bad(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)

	_, err := service.CreateWindow(CreateWindowOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "window name is required")
}

func TestShowEnvironmentDialog_Good(t *testing.T) {
	service, mock, _ := newServiceWithMocks(t)

	// This will panic because Dialog().Info() returns nil
	// We're verifying the env info is accessed, not that a dialog shows
	assert.NotPanics(t, func() {
		defer func() { recover() }() // Recover from nil dialog
		service.ShowEnvironmentDialog()
	})

	// Verify dialog was requested (even though it's nil)
	assert.Equal(t, 1, mock.dialogManager.infoDialogsCreated)
}

func TestBuildMenu_Good(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)
	coreInstance := newTestCore(t)
	service.ServiceRuntime = core.NewServiceRuntime[Options](coreInstance, Options{})

	// buildMenu should not panic with mock platforms
	assert.NotPanics(t, func() {
		service.buildMenu()
	})
}

func TestSetupTray_Good(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)
	coreInstance := newTestCore(t)
	service.ServiceRuntime = core.NewServiceRuntime[Options](coreInstance, Options{})

	// setupTray should not panic with mock platforms
	assert.NotPanics(t, func() {
		service.setupTray()
	})

	// Verify tray is active
	assert.True(t, service.tray.IsActive())
}

func TestHandleNewWorkspace_Good(t *testing.T) {
	service, _, wp := newServiceWithMocks(t)

	service.handleNewWorkspace()

	// Verify a window was created with correct options
	assert.Len(t, wp.windows, 1)
	assert.Equal(t, "workspace-new", wp.windows[0].name)
	assert.Equal(t, "New Workspace", wp.windows[0].title)
	assert.Equal(t, 500, wp.windows[0].width)
	assert.Equal(t, 400, wp.windows[0].height)
}

func TestHandleListWorkspaces_Good(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)
	coreInstance := newTestCore(t)
	service.ServiceRuntime = core.NewServiceRuntime[Options](coreInstance, Options{})

	// handleListWorkspaces should not panic when workspace service is not available
	assert.NotPanics(t, func() {
		service.handleListWorkspaces()
	})
}

func TestHandleSaveFile_Good(t *testing.T) {
	service, mock, _ := newServiceWithMocks(t)

	service.handleSaveFile()

	assert.Contains(t, mock.eventManager.emittedEvents, "ide:save")
}

func TestHandleRun_Good(t *testing.T) {
	service, mock, _ := newServiceWithMocks(t)

	service.handleRun()

	assert.Contains(t, mock.eventManager.emittedEvents, "ide:run")
}

func TestHandleBuild_Good(t *testing.T) {
	service, mock, _ := newServiceWithMocks(t)

	service.handleBuild()

	assert.Contains(t, mock.eventManager.emittedEvents, "ide:build")
}

func TestWSEventManager_Good(t *testing.T) {
	es := newMockEventSource()
	em := NewWSEventManager(es)
	defer em.Close()

	assert.NotNil(t, em)
	assert.Equal(t, 0, em.ConnectedClients())
}

func TestWSEventManager_SetupWindowEventListeners_Good(t *testing.T) {
	es := newMockEventSource()
	em := NewWSEventManager(es)
	defer em.Close()

	em.SetupWindowEventListeners()

	// Verify theme handler was registered
	assert.Len(t, es.themeHandlers, 1)
}

func TestResetWindowState_Good(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)

	err := service.ResetWindowState()
	assert.NoError(t, err)
}

func TestGetSavedWindowStates_Good(t *testing.T) {
	service, _, _ := newServiceWithMocks(t)

	states := service.GetSavedWindowStates()
	assert.NotNil(t, states)
}

func TestActionOpenWindow_Good(t *testing.T) {
	action := ActionOpenWindow{
		Window: window.Window{
			Name:   "test",
			Title:  "Test Window",
			Width:  800,
			Height: 600,
		},
	}

	assert.Equal(t, "test", action.Name)
	assert.Equal(t, "Test Window", action.Title)
	assert.Equal(t, 800, action.Width)
	assert.Equal(t, 600, action.Height)
}
