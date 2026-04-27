package application

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wailsapp/wails/v3/pkg/events"
)

var _ Window = (*WebviewWindow)(nil)
var _ Window = (*BrowserWindow)(nil)

func TestApplication_NewRGBA_Good(t *testing.T) {
	got := NewRGBA(1, 2, 3, 4)

	assert.Equal(t, RGBA{Red: 1, Green: 2, Blue: 3, Alpha: 4}, got)
}

func TestApplication_NewRGBA_Bad(t *testing.T) {
	got := NewRGBA(0, 0, 0, 0)

	assert.Equal(t, RGBA{}, got)
}

func TestApplication_NewRGBA_Ugly(t *testing.T) {
	got := NewRGBA(255, 255, 255, 255)

	assert.Equal(t, RGBA{Red: 255, Green: 255, Blue: 255, Alpha: 255}, got)
}

func TestApplication_MenuRole_String_Good(t *testing.T) {
	cases := []struct {
		role MenuRole
		want string
	}{
		{AppMenu, "app"},
		{FileMenu, "file"},
		{EditMenu, "edit"},
		{ViewMenu, "view"},
		{WindowMenu, "window"},
		{HelpMenu, "help"},
	}

	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.role.String())
		})
	}
}

func TestApplication_MenuRole_String_Bad(t *testing.T) {
	assert.Equal(t, "unknown", MenuRole(-1).String())
}

func TestApplication_MenuRole_String_Ugly(t *testing.T) {
	assert.Equal(t, "unknown", MenuRole(999).String())
}

func TestApplication_MenuItem_OnClick_Good(t *testing.T) {
	called := 0
	item := &MenuItem{}

	item.OnClick(func(*Context) {
		called++
	})
	require.NotNil(t, item.onClick)
	item.onClick(&Context{})

	assert.Equal(t, 1, called)
}

func TestApplication_MenuItem_OnClick_Bad(t *testing.T) {
	item := &MenuItem{}

	item.OnClick(nil)

	assert.Nil(t, item.onClick)
}

func TestApplication_MenuItem_OnClick_Ugly(t *testing.T) {
	called := 0
	item := &MenuItem{}

	item.OnClick(func(*Context) {
		called++
	})
	item.OnClick(func(*Context) {
		called++
	})
	require.NotNil(t, item.onClick)
	item.onClick(&Context{})

	assert.Equal(t, 1, called)
}

func TestApplication_Menu_Good(t *testing.T) {
	menu := NewMenu()
	menuItem := menu.Add("Open")
	menu.AddSeparator()
	submenu := menu.AddSubmenu("More")
	menu.AddRole(FileMenu)

	require.NotNil(t, menuItem)
	assert.Equal(t, "Open", menuItem.Label)
	require.NotNil(t, submenu)
	assert.Same(t, submenu, menu.Items[2].GetSubmenu())
	assert.Len(t, menu.Items, 4)
	assert.Equal(t, "Open", menu.Items[0].Label)
	assert.Equal(t, "---", menu.Items[1].Label)
	assert.Equal(t, "More", menu.Items[2].Label)
	assert.Equal(t, "file", menu.Items[3].Label)
}

func TestApplication_Menu_Bad(t *testing.T) {
	menu := NewMenu()

	assert.Empty(t, menu.Items)
}

func TestApplication_Menu_Ugly(t *testing.T) {
	menu := NewMenu()
	submenu := menu.AddSubmenu("Nested")
	submenu.Add("Child")

	assert.Len(t, menu.Items, 1)
	assert.Len(t, submenu.Items, 1)
	assert.Equal(t, "Nested", menu.Items[0].Label)
}

func TestApplication_MenuManager_SetApplicationMenu_Good(t *testing.T) {
	manager := &MenuManager{}
	menu := NewMenu()

	manager.SetApplicationMenu(menu)

	assert.Same(t, menu, manager.applicationMenu)
}

func TestApplication_MenuManager_SetApplicationMenu_Bad(t *testing.T) {
	manager := &MenuManager{}

	manager.SetApplicationMenu(nil)

	assert.Nil(t, manager.applicationMenu)
}

func TestApplication_MenuManager_SetApplicationMenu_Ugly(t *testing.T) {
	manager := &MenuManager{}
	first := NewMenu()
	second := NewMenu()

	manager.SetApplicationMenu(first)
	manager.SetApplicationMenu(second)

	assert.Same(t, second, manager.applicationMenu)
}

func TestApplication_MenuManager_SetApplicationMenu_NilReceiver(t *testing.T) {
	var manager *MenuManager

	assert.NotPanics(t, func() {
		manager.SetApplicationMenu(NewMenu())
	})
}

func TestApplication_SystemTray_Good(t *testing.T) {
	tray := (&SystemTrayManager{}).New()
	menu := NewMenu()
	window := (&WindowManager{}).NewWithOptions(WebviewWindowOptions{Name: "tray"})

	tray.SetIcon([]byte{1, 2, 3})
	tray.SetTemplateIcon([]byte{4, 5, 6})
	tray.SetTooltip("tooltip")
	tray.SetLabel("label")
	tray.SetMenu(menu)
	tray.AttachWindow(window)

	assert.Equal(t, []byte{1, 2, 3}, tray.icon)
	assert.Equal(t, []byte{4, 5, 6}, tray.templateIcon)
	assert.Equal(t, "tooltip", tray.tooltip)
	assert.Equal(t, "label", tray.label)
	assert.Same(t, menu, tray.menu)
	assert.Same(t, window, tray.attachedWindow)
}

func TestApplication_SystemTray_Bad(t *testing.T) {
	tray := (&SystemTrayManager{}).New()

	assert.Empty(t, tray.icon)
	assert.Empty(t, tray.templateIcon)
	assert.Empty(t, tray.tooltip)
	assert.Empty(t, tray.label)
	assert.Nil(t, tray.menu)
	assert.Nil(t, tray.attachedWindow)
}

func TestApplication_SystemTray_Ugly(t *testing.T) {
	tray := (&SystemTrayManager{}).New()
	icon := []byte{9, 8, 7}
	tray.SetIcon(icon)
	icon[0] = 0

	assert.Equal(t, []byte{9, 8, 7}, tray.icon)
}

func TestApplication_WindowEventContext_Good(t *testing.T) {
	ctx := &WindowEventContext{
		droppedFiles: []string{"a", "b"},
		dropDetails:  &DropTargetDetails{ElementID: "drop"},
	}

	assert.Equal(t, []string{"a", "b"}, ctx.DroppedFiles())
	require.NotNil(t, ctx.DropTargetDetails())
	assert.Equal(t, "drop", ctx.DropTargetDetails().ElementID)
}

func TestApplication_WindowEventContext_Bad(t *testing.T) {
	ctx := &WindowEventContext{}

	assert.Empty(t, ctx.DroppedFiles())
	assert.Nil(t, ctx.DropTargetDetails())
}

func TestApplication_WindowEventContext_Ugly(t *testing.T) {
	ctx := &WindowEventContext{droppedFiles: []string{"x"}}
	files := ctx.DroppedFiles()
	files[0] = "mutated"

	assert.Equal(t, []string{"x"}, ctx.DroppedFiles())
}

func TestApplication_WindowEvent_Good(t *testing.T) {
	event := &WindowEvent{}

	require.NotNil(t, event.Context())
	require.Same(t, event.Context(), event.Context())
}

func TestApplication_WindowEvent_Bad(t *testing.T) {
	event := &WindowEvent{ctx: &WindowEventContext{}}

	assert.Same(t, event.ctx, event.Context())
}

func TestApplication_WindowEvent_Ugly(t *testing.T) {
	event := &WindowEvent{}
	event.Context().droppedFiles = []string{"file"}

	assert.Equal(t, []string{"file"}, event.Context().DroppedFiles())
}

func TestApplication_WindowEvent_NilReceiver(t *testing.T) {
	var event *WindowEvent

	assert.Nil(t, event.Context())
}

func TestApplication_WindowEventContext_NilReceiver(t *testing.T) {
	var ctx *WindowEventContext

	assert.Empty(t, ctx.DroppedFiles())
	assert.Nil(t, ctx.DropTargetDetails())
}

func TestApplication_WebviewWindow_Good(t *testing.T) {
	manager := &WindowManager{}
	window := manager.NewWithOptions(WebviewWindowOptions{
		Name:        "main",
		Title:       "Main",
		URL:         "https://example.com",
		HTML:        "<h1>Hello</h1>",
		X:           10,
		Y:           20,
		Width:       800,
		Height:      600,
		Hidden:      false,
		AlwaysOnTop: true,
	})

	assert.Equal(t, uint(1), window.ID())
	assert.Equal(t, "main", window.Name())
	assert.Equal(t, "Main", window.Title())
	assert.Equal(t, 10, window.Bounds().X)
	assert.Equal(t, 20, window.Bounds().Y)
	assert.Equal(t, 800, window.Width())
	assert.Equal(t, 600, window.Height())
	assert.True(t, window.IsVisible())
	assert.False(t, window.IsFullscreen())
	assert.False(t, window.IsMaximised())
	assert.True(t, window.GetOpacity() > 0)
	assert.Equal(t, 1.0, window.GetZoom())
	assert.Same(t, window, window.SetAlwaysOnTop(false))
	assert.Same(t, window, window.SetOpacity(0.5))
	assert.Equal(t, 0.5, window.GetOpacity())
	screen, err := window.GetScreen()
	require.NoError(t, err)
	require.NotNil(t, screen)
	assert.Equal(t, Screen{}, *screen)
	assert.Equal(t, &LRTB{}, window.GetBorderSizes())
	assert.Same(t, window, window.SetURL("https://example.org"))
	assert.Same(t, window, window.SetHTML("<p>Hi</p>"))
	assert.Same(t, window, window.SetSize(1024, 768))
	assert.Equal(t, 1024, window.Width())
	assert.Equal(t, 768, window.Height())
}

func TestApplication_WebviewWindow_Bad(t *testing.T) {
	manager := &WindowManager{}
	window := manager.NewWithOptions(WebviewWindowOptions{Hidden: true})

	assert.False(t, window.IsVisible())
	assert.Nil(t, manager.Get("missing"))
	assert.Nil(t, manager.GetByID(99))
}

func TestApplication_WebviewWindow_Ugly(t *testing.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	assert.Same(t, first, manager.Get("first"))
	assert.Same(t, second, manager.GetByID(2))
	assert.Len(t, manager.GetAll(), 2)
}

func TestApplication_BrowserWindow_StateTransitions(t *testing.T) {
	window := NewBrowserWindow(99, "client")

	assert.Equal(t, 0, window.Width())
	assert.Equal(t, 0, window.Height())
	assert.Equal(t, Rect{}, window.Bounds())
	assert.False(t, window.Resizable())

	window.SetPosition(10, 20)
	assert.Same(t, window, window.SetSize(300, 200))
	posX, posY := window.Position()
	assert.Equal(t, 10, posX)
	assert.Equal(t, 20, posY)
	assert.Equal(t, 300, window.Width())
	assert.Equal(t, 200, window.Height())
	assert.Equal(t, Rect{X: 10, Y: 20, Width: 300, Height: 200}, window.Bounds())

	window.SetBounds(Rect{X: 1, Y: 2, Width: 3, Height: 4})
	assert.Equal(t, Rect{X: 1, Y: 2, Width: 3, Height: 4}, window.Bounds())
	relX, relY := window.RelativePosition()
	assert.Equal(t, 1, relX)
	assert.Equal(t, 2, relY)

	assert.Same(t, window, window.SetResizable(true))
	assert.True(t, window.Resizable())
	assert.Same(t, window, window.SetIgnoreMouseEvents(true))
	assert.True(t, window.IsIgnoreMouseEvents())
	assert.Same(t, window, window.SetZoom(1.5))
	assert.Equal(t, 1.5, window.GetZoom())
	assert.Same(t, window, window.ZoomReset())
	assert.Equal(t, 1.0, window.GetZoom())
	screen, err := window.GetScreen()
	require.NoError(t, err)
	require.NotNil(t, screen)
	assert.Equal(t, Screen{}, *screen)
	assert.Equal(t, &LRTB{}, window.GetBorderSizes())

	assert.Same(t, window, window.Fullscreen())
	assert.True(t, window.IsFullscreen())
	assert.Same(t, window, window.Maximise())
	assert.True(t, window.IsMaximised())
	assert.Same(t, window, window.Minimise())
	assert.True(t, window.IsMinimised())
	assert.Same(t, window, window.Show())
	assert.True(t, window.IsVisible())
	assert.False(t, window.IsMinimised())
	window.Restore()
	assert.False(t, window.IsFullscreen())
	assert.False(t, window.IsMaximised())
}

func TestApplication_App_Good(t *testing.T) {
	app := &App{}

	assert.NotNil(t, app.NewMenu())
}

func TestApplication_App_Bad(t *testing.T) {
	var app App

	assert.Zero(t, app.Logger)
	assert.Empty(t, app.Window.GetAll())
	assert.Nil(t, app.Menu.applicationMenu)
}

func TestApplication_App_Ugly(t *testing.T) {
	app := &App{}
	app.Quit()
}

func TestApplication_AppManagers_Good(t *testing.T) {
	app := &App{}

	require.NotNil(t, &app.Window)
	require.NotNil(t, &app.Menu)
	require.NotNil(t, &app.Dialog)
	require.NotNil(t, &app.Event)
	require.NotNil(t, &app.Browser)
	require.NotNil(t, &app.Clipboard)
	require.NotNil(t, &app.ContextMenu)
	require.NotNil(t, &app.Environment)
	require.NotNil(t, &app.Screen)
	require.NotNil(t, &app.SystemTray)
	require.NotNil(t, &app.KeyBinding)

	assert.NotPanics(t, func() {
		window := app.Window.NewWithOptions(WebviewWindowOptions{Name: "app-managers"})
		require.NotNil(t, window)

		menu := app.NewMenu()
		app.Menu.SetApplicationMenu(menu)
		assert.Same(t, menu, app.Menu.applicationMenu)
		assert.NotNil(t, app.Dialog.Info())
		_, err := app.Dialog.ShowInfo("Done", "Saved")
		assert.NoError(t, err)

		assert.NoError(t, app.Browser.OpenURL("https://example.com"))
		assert.NoError(t, app.Browser.Open("https://example.com"))

		assert.True(t, app.Clipboard.SetText("copied"))
		text, ok := app.Clipboard.Text()
		assert.True(t, ok)
		assert.Equal(t, "copied", text)

		contextMenu := app.ContextMenu.New()
		app.ContextMenu.Add("main", contextMenu)
		gotMenu, exists := app.ContextMenu.Get("main")
		assert.True(t, exists)
		assert.Same(t, contextMenu, gotMenu)

		env := newEnvironmentManager().Info()
		assert.Equal(t, runtime.GOOS, env.OS)
		assert.Equal(t, runtime.GOARCH, env.Arch)
		assert.NoError(t, app.Environment.OpenFileManager("/tmp", false))
		assert.False(t, app.Environment.HasFocusFollowsMouse())

		cancel := app.Event.Once("ready", func(*CustomEvent) {})
		require.NotNil(t, cancel)
		assert.False(t, app.Event.Emit("ready"))
		cancel()

		screen := &Screen{
			ID:        "primary",
			IsPrimary: true,
			Bounds:    Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
			Size:      Size{Width: 1920, Height: 1080},
		}
		assert.NoError(t, app.Screen.LayoutScreens([]*Screen{screen}))
		assert.Same(t, screen, app.Screen.GetPrimary())
		assert.Same(t, screen, app.Screen.Primary())
		assert.Equal(t, Point{X: 5, Y: 6}, app.Screen.DipToPhysicalPoint(Point{X: 5, Y: 6}))

		triggered := 0
		app.KeyBinding.Register("CmdOrCtrl+K", func(Window) { triggered++ })
		assert.True(t, app.KeyBinding.Process("CmdOrCtrl+K", nil))
		assert.Equal(t, 1, triggered)

		tray := app.SystemTray.New()
		assert.NotNil(t, tray)
	})
}

func TestApplication_AppManagers_Bad(t *testing.T) {
	app := &App{}

	assert.NotPanics(t, func() {
		assert.Nil(t, app.Window.GetByID(1))
		app.Menu.SetApplicationMenu(nil)

		_, err := app.Dialog.ShowError()
		assert.NoError(t, err)

		assert.NoError(t, app.Browser.Open(""))
		text, ok := app.Clipboard.Text()
		assert.False(t, ok)
		assert.Empty(t, text)

		app.ContextMenu.Remove("missing")
		_, exists := app.ContextMenu.Get("missing")
		assert.False(t, exists)

		info := app.Environment.Info()
		assert.Empty(t, info.OS)
		assert.Empty(t, info.Arch)
		assert.NoError(t, app.Environment.OpenFileManager("", false))

		app.Event.Off("missing")
		app.Event.Reset()

		assert.Nil(t, app.Screen.Primary())
		app.KeyBinding.Unregister("missing")
		assert.False(t, app.KeyBinding.Process("missing", nil))
	})
}

func TestApplication_AppManagers_Ugly(t *testing.T) {
	app := &App{}

	assert.NotPanics(t, func() {
		assert.True(t, app.Clipboard.SetText("zero\x00byte"))
		assert.NoError(t, app.Browser.Open("/tmp/\x00report.txt"))

		app.ContextMenu.Add("dup", app.ContextMenu.New())
		app.ContextMenu.Add("dup", app.ContextMenu.New())
		assert.Len(t, app.ContextMenu.GetAll(), 1)

		cancelHook := app.Event.RegisterApplicationEventHook(events.ApplicationEventType(9), func(event *ApplicationEvent) {
			event.Cancel()
		})
		cancelListener := app.Event.OnApplicationEvent(events.ApplicationEventType(9), func(*ApplicationEvent) {
			t.Fatal("cancelled event should not reach listeners")
		})
		app.Event.handleApplicationEvent(&ApplicationEvent{Id: 9})
		cancelListener()
		cancelHook()

		screen := &Screen{
			ID:        "primary",
			IsPrimary: true,
			Bounds:    Rect{X: 0, Y: 0, Width: 100, Height: 100},
			Size:      Size{Width: 100, Height: 100},
		}
		app.Screen.SetScreens([]*Screen{screen})
		assert.Same(t, screen, app.Screen.ScreenNearestDipPoint(Point{X: 50, Y: 50}))
		assert.Same(t, screen, app.Screen.ScreenNearestDipRect(Rect{X: 10, Y: 10, Width: 5, Height: 5}))

		triggered := 0
		app.KeyBinding.Register("CmdOrCtrl+Shift+P", func(Window) { triggered++ })
		app.KeyBinding.handleWindowKeyEvent(&windowKeyEvent{acceleratorString: "CmdOrCtrl+Shift+P"})
		assert.Equal(t, 1, triggered)
	})
}
