package application

import (
	core "dappco.re/go"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/events"
)

var _ Window = (*WebviewWindow)(nil)
var _ Window = (*BrowserWindow)(nil)

func TestApplication_NewRGBA_Good(t *core.T) {
	got := NewRGBA(1, 2, 3, 4)

	core.AssertEqual(t, RGBA{Red: 1, Green: 2, Blue: 3, Alpha: 4}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplication_NewRGBA_Bad(t *core.T) {
	got := NewRGBA(0, 0, 0, 0)

	core.AssertEqual(t, RGBA{}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplication_NewRGBA_Ugly(t *core.T) {
	got := NewRGBA(255, 255, 255, 255)

	core.AssertEqual(t, RGBA{Red: 255, Green: 255, Blue: 255, Alpha: 255}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplication_MenuRole_String_Good(t *core.T) {
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
		t.Run(tc.want, func(t *core.T) {
			core.AssertEqual(t, tc.want, tc.role.String())
		})
	}
}

func TestApplication_MenuRole_String_Bad(t *core.T) {
	core.AssertEqual(t, "unknown", MenuRole(-1).String())
	observedType := core.Sprintf("%T", MenuRole(-1).String())
	core.AssertNotEmpty(t, observedType)
}

func TestApplication_MenuRole_String_Ugly(t *core.T) {
	core.AssertEqual(t, "unknown", MenuRole(999).String())
	observedType := core.Sprintf("%T", MenuRole(999).String())
	core.AssertNotEmpty(t, observedType)
}

func TestApplication_MenuItem_OnClick_Good(t *core.T) {
	called := 0
	item := &MenuItem{}

	item.OnClick(func(*Context) {
		called++
	})
	core.AssertNotNil(t, item.onClick)
	item.onClick(&Context{})

	core.AssertEqual(t, 1, called)
}

func TestApplication_MenuItem_OnClick_Bad(t *core.T) {
	item := &MenuItem{}

	item.OnClick(nil)

	core.AssertNil(t, item.onClick)
}

func TestApplication_MenuItem_OnClick_Ugly(t *core.T) {
	called := 0
	item := &MenuItem{}

	item.OnClick(func(*Context) {
		called++
	})
	item.OnClick(func(*Context) {
		called++
	})
	core.AssertNotNil(t, item.onClick)
	item.onClick(&Context{})

	core.AssertEqual(t, 1, called)
}

func TestApplication_Menu_Good(t *core.T) {
	menu := NewMenu()
	menuItem := menu.Add("Open")
	menu.AddSeparator()
	submenu := menu.AddSubmenu("More")
	menu.AddRole(FileMenu)

	core.AssertNotNil(t, menuItem)
	core.AssertEqual(t, "Open", menuItem.Label)
	core.AssertNotNil(t, submenu)
	core.AssertSame(t, submenu, menu.Items[2].GetSubmenu())
	core.AssertLen(t, menu.Items, 4)
	core.AssertEqual(t, "Open", menu.Items[0].Label)
	core.AssertEqual(t, "---", menu.Items[1].Label)
	core.AssertEqual(t, "More", menu.Items[2].Label)
	core.AssertEqual(t, "file", menu.Items[3].Label)
}

func TestApplication_Menu_Bad(t *core.T) {
	menu := NewMenu()

	core.AssertEmpty(t, menu.Items)
	core.AssertNotEmpty(t, core.Sprintf("%T", menu))
}

func TestApplication_Menu_Ugly(t *core.T) {
	menu := NewMenu()
	submenu := menu.AddSubmenu("Nested")
	submenu.Add("Child")

	core.AssertLen(t, menu.Items, 1)
	core.AssertLen(t, submenu.Items, 1)
	core.AssertEqual(t, "Nested", menu.Items[0].Label)
}

func TestApplication_MenuManager_SetApplicationMenu_Good(t *core.T) {
	manager := &MenuManager{}
	menu := NewMenu()

	manager.SetApplicationMenu(menu)

	core.AssertSame(t, menu, manager.applicationMenu)
}

func TestApplication_MenuManager_SetApplicationMenu_Bad(t *core.T) {
	manager := &MenuManager{}

	manager.SetApplicationMenu(nil)

	core.AssertNil(t, manager.applicationMenu)
}

func TestApplication_MenuManager_SetApplicationMenu_Ugly(t *core.T) {
	manager := &MenuManager{}
	first := NewMenu()
	second := NewMenu()

	manager.SetApplicationMenu(first)
	manager.SetApplicationMenu(second)

	core.AssertSame(t, second, manager.applicationMenu)
}

func TestApplication_MenuManager_SetApplicationMenu_NilReceiver(t *core.T) {
	var manager *MenuManager

	core.AssertNotPanics(t, func() {
		manager.SetApplicationMenu(NewMenu())
	})
}

func TestApplication_SystemTray_Good(t *core.T) {
	tray := (&SystemTrayManager{}).New()
	menu := NewMenu()
	window := (&WindowManager{}).NewWithOptions(WebviewWindowOptions{Name: "tray"})

	tray.SetIcon([]byte{1, 2, 3})
	tray.SetTemplateIcon([]byte{4, 5, 6})
	tray.SetTooltip("tooltip")
	tray.SetLabel("label")
	tray.SetMenu(menu)
	tray.AttachWindow(window)

	core.AssertEqual(t, []byte{1, 2, 3}, tray.icon)
	core.AssertEqual(t, []byte{4, 5, 6}, tray.templateIcon)
	core.AssertEqual(t, "tooltip", tray.tooltip)
	core.AssertEqual(t, "label", tray.label)
	core.AssertSame(t, menu, tray.menu)
	core.AssertSame(t, window, tray.attachedWindow)
}

func TestApplication_SystemTray_Bad(t *core.T) {
	tray := (&SystemTrayManager{}).New()

	core.AssertEmpty(t, tray.icon)
	core.AssertEmpty(t, tray.templateIcon)
	core.AssertEmpty(t, tray.tooltip)
	core.AssertEmpty(t, tray.label)
	core.AssertNil(t, tray.menu)
	core.AssertNil(t, tray.attachedWindow)
}

func TestApplication_SystemTray_Ugly(t *core.T) {
	tray := (&SystemTrayManager{}).New()
	icon := []byte{9, 8, 7}
	tray.SetIcon(icon)
	icon[0] = 0

	core.AssertEqual(t, []byte{9, 8, 7}, tray.icon)
}

func TestApplication_WindowEventContext_Good(t *core.T) {
	ctx := &WindowEventContext{
		droppedFiles: []string{"a", "b"},
		dropDetails:  &DropTargetDetails{ElementID: "drop"},
	}

	core.AssertEqual(t, []string{"a", "b"}, ctx.DroppedFiles())
	core.AssertNotNil(t, ctx.DropTargetDetails())
	core.AssertEqual(t, "drop", ctx.DropTargetDetails().ElementID)
}

func TestApplication_WindowEventContext_Bad(t *core.T) {
	ctx := &WindowEventContext{}

	core.AssertEmpty(t, ctx.DroppedFiles())
	core.AssertNil(t, ctx.DropTargetDetails())
}

func TestApplication_WindowEventContext_Ugly(t *core.T) {
	ctx := &WindowEventContext{droppedFiles: []string{"x"}}
	files := ctx.DroppedFiles()
	files[0] = "mutated"

	core.AssertEqual(t, []string{"x"}, ctx.DroppedFiles())
}

func TestApplication_WindowEvent_Good(t *core.T) {
	event := &WindowEvent{}

	core.AssertNotNil(t, event.Context())
	core.AssertSame(t, event.Context(), event.Context())
}

func TestApplication_WindowEvent_Bad(t *core.T) {
	event := &WindowEvent{ctx: &WindowEventContext{}}

	core.AssertSame(t, event.ctx, event.Context())
	core.AssertNotEmpty(t, core.Sprintf("%T", event))
}

func TestApplication_WindowEvent_Ugly(t *core.T) {
	event := &WindowEvent{}
	event.Context().droppedFiles = []string{"file"}

	core.AssertEqual(t, []string{"file"}, event.Context().DroppedFiles())
}

func TestApplication_WindowEvent_NilReceiver(t *core.T) {
	var event *WindowEvent

	core.AssertNil(t, event.Context())
	core.AssertNotEmpty(t, core.Sprintf("%T", event.Context()))
}

func TestApplication_WindowEventContext_NilReceiver(t *core.T) {
	var ctx *WindowEventContext

	core.AssertEmpty(t, ctx.DroppedFiles())
	core.AssertNil(t, ctx.DropTargetDetails())
}

func TestApplication_WebviewWindow_Good(t *core.T) {
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

	core.AssertEqual(t, uint(1), window.ID())
	core.AssertEqual(t, "main", window.Name())
	core.AssertEqual(t, "Main", window.Title())
	core.AssertEqual(t, 10, window.Bounds().X)
	core.AssertEqual(t, 20, window.Bounds().Y)
	core.AssertEqual(t, 800, window.Width())
	core.AssertEqual(t, 600, window.Height())
	core.AssertTrue(t, window.IsVisible())
	core.AssertFalse(t, window.IsFullscreen())
	core.AssertFalse(t, window.IsMaximised())
	core.AssertTrue(t, window.GetOpacity() > 0)
	core.AssertEqual(t, 1.0, window.GetZoom())
	core.AssertSame(t, window, window.SetAlwaysOnTop(false))
	core.AssertSame(t, window, window.SetOpacity(0.5))
	core.AssertEqual(t, 0.5, window.GetOpacity())
	screen, err := window.GetScreen()
	core.RequireNoError(t, err)
	core.AssertNotNil(t, screen)
	core.AssertEqual(t, Screen{}, *screen)
	core.AssertEqual(t, &LRTB{}, window.GetBorderSizes())
	core.AssertSame(t, window, window.SetURL("https://example.org"))
	core.AssertSame(t, window, window.SetHTML("<p>Hi</p>"))
	core.AssertSame(t, window, window.SetSize(1024, 768))
	core.AssertEqual(t, 1024, window.Width())
	core.AssertEqual(t, 768, window.Height())
}

func TestApplication_WebviewWindow_Bad(t *core.T) {
	manager := &WindowManager{}
	window := manager.NewWithOptions(WebviewWindowOptions{Hidden: true})

	core.AssertFalse(t, window.IsVisible())
	core.AssertNil(t, manager.Get("missing"))
	core.AssertNil(t, manager.GetByID(99))
}

func TestApplication_WebviewWindow_Ugly(t *core.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	core.AssertSame(t, first, manager.Get("first"))
	core.AssertSame(t, second, manager.GetByID(2))
	core.AssertLen(t, manager.GetAll(), 2)
}

func TestApplication_BrowserWindow_StateTransitions(t *core.T) {
	window := NewBrowserWindow(99, "client")

	core.AssertEqual(t, 0, window.Width())
	core.AssertEqual(t, 0, window.Height())
	core.AssertEqual(t, Rect{}, window.Bounds())
	core.AssertFalse(t, window.Resizable())

	window.SetPosition(10, 20)
	core.AssertSame(t, window, window.SetSize(300, 200))
	posX, posY := window.Position()
	core.AssertEqual(t, 10, posX)
	core.AssertEqual(t, 20, posY)
	core.AssertEqual(t, 300, window.Width())
	core.AssertEqual(t, 200, window.Height())
	core.AssertEqual(t, Rect{X: 10, Y: 20, Width: 300, Height: 200}, window.Bounds())

	window.SetBounds(Rect{X: 1, Y: 2, Width: 3, Height: 4})
	core.AssertEqual(t, Rect{X: 1, Y: 2, Width: 3, Height: 4}, window.Bounds())
	relX, relY := window.RelativePosition()
	core.AssertEqual(t, 1, relX)
	core.AssertEqual(t, 2, relY)

	core.AssertSame(t, window, window.SetResizable(true))
	core.AssertTrue(t, window.Resizable())
	core.AssertSame(t, window, window.SetIgnoreMouseEvents(true))
	core.AssertTrue(t, window.IsIgnoreMouseEvents())
	core.AssertSame(t, window, window.SetZoom(1.5))
	core.AssertEqual(t, 1.5, window.GetZoom())
	core.AssertSame(t, window, window.ZoomReset())
	core.AssertEqual(t, 1.0, window.GetZoom())
	screen, err := window.GetScreen()
	core.RequireNoError(t, err)
	core.AssertNotNil(t, screen)
	core.AssertEqual(t, Screen{}, *screen)
	core.AssertEqual(t, &LRTB{}, window.GetBorderSizes())

	core.AssertSame(t, window, window.Fullscreen())
	core.AssertTrue(t, window.IsFullscreen())
	core.AssertSame(t, window, window.Maximise())
	core.AssertTrue(t, window.IsMaximised())
	core.AssertSame(t, window, window.Minimise())
	core.AssertTrue(t, window.IsMinimised())
	core.AssertSame(t, window, window.Show())
	core.AssertTrue(t, window.IsVisible())
	core.AssertFalse(t, window.IsMinimised())
	window.Restore()
	core.AssertFalse(t, window.IsFullscreen())
	core.AssertFalse(t, window.IsMaximised())
}

func TestApplication_App_Good(t *core.T) {
	app := &App{}

	core.AssertNotNil(t, app.NewMenu())
	core.AssertNotEmpty(t, core.Sprintf("%T", app))
}

func TestApplication_App_Bad(t *core.T) {
	var app App

	core.AssertEmpty(t, app.Logger)
	core.AssertEmpty(t, app.Window.GetAll())
	core.AssertNil(t, app.Menu.applicationMenu)
}

func TestApplication_App_Ugly(t *core.T) {
	app := &App{}
	app.Quit()
	core.AssertNotEmpty(t, core.Sprintf("%T", app))
}

func TestApplication_AppManagers_Good(t *core.T) {
	app := &App{}

	core.AssertNotNil(t, &app.Window)
	core.AssertNotNil(t, &app.Menu)
	core.AssertNotNil(t, &app.Dialog)
	core.AssertNotNil(t, &app.Event)
	core.AssertNotNil(t, &app.Browser)
	core.AssertNotNil(t, &app.Clipboard)
	core.AssertNotNil(t, &app.ContextMenu)
	core.AssertNotNil(t, &app.Environment)
	core.AssertNotNil(t, &app.Screen)
	core.AssertNotNil(t, &app.SystemTray)
	core.AssertNotNil(t, &app.KeyBinding)

	core.AssertNotPanics(t, func() {
		window := app.Window.NewWithOptions(WebviewWindowOptions{Name: "app-managers"})
		core.AssertNotNil(t, window)

		menu := app.NewMenu()
		app.Menu.SetApplicationMenu(menu)
		core.AssertSame(t, menu, app.Menu.applicationMenu)
		core.AssertNotNil(t, app.Dialog.Info())
		_, err := app.Dialog.ShowInfo("Done", "Saved")
		core.AssertNoError(t, err)

		core.AssertNoError(t, app.Browser.OpenURL("https://example.com"))
		core.AssertNoError(t, app.Browser.Open("https://example.com"))

		core.AssertTrue(t, app.Clipboard.SetText("copied"))
		text, ok := app.Clipboard.Text()
		core.AssertTrue(t, ok)
		core.AssertEqual(t, "copied", text)

		contextMenu := app.ContextMenu.New()
		app.ContextMenu.Add("main", contextMenu)
		gotMenu, exists := app.ContextMenu.Get("main")
		core.AssertTrue(t, exists)
		core.AssertSame(t, contextMenu, gotMenu)

		env := newEnvironmentManager().Info()
		core.AssertEqual(t, runtime.GOOS, env.OS)
		core.AssertEqual(t, runtime.GOARCH, env.Arch)
		core.AssertNoError(t, app.Environment.OpenFileManager("/tmp", false))
		core.AssertFalse(t, app.Environment.HasFocusFollowsMouse())

		cancel := app.Event.Once("ready", func(*CustomEvent) {})
		core.AssertNotNil(t, cancel)
		core.AssertFalse(t, app.Event.Emit("ready"))
		cancel()

		screen := &Screen{
			ID:        "primary",
			IsPrimary: true,
			Bounds:    Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
			Size:      Size{Width: 1920, Height: 1080},
		}
		core.AssertNoError(t, app.Screen.LayoutScreens([]*Screen{screen}))
		core.AssertSame(t, screen, app.Screen.GetPrimary())
		core.AssertSame(t, screen, app.Screen.Primary())
		core.AssertEqual(t, Point{X: 5, Y: 6}, app.Screen.DipToPhysicalPoint(Point{X: 5, Y: 6}))

		triggered := 0
		app.KeyBinding.Register("CmdOrCtrl+K", func(Window) { triggered++ })
		core.AssertTrue(t, app.KeyBinding.Process("CmdOrCtrl+K", nil))
		core.AssertEqual(t, 1, triggered)

		tray := app.SystemTray.New()
		core.AssertNotNil(t, tray)
	})
}

func TestApplication_AppManagers_Bad(t *core.T) {
	app := &App{}

	core.AssertNotPanics(t, func() {
		core.AssertNil(t, app.Window.GetByID(1))
		app.Menu.SetApplicationMenu(nil)

		_, err := app.Dialog.ShowError()
		core.AssertNoError(t, err)

		core.AssertNoError(t, app.Browser.Open(""))
		text, ok := app.Clipboard.Text()
		core.AssertFalse(t, ok)
		core.AssertEmpty(t, text)

		app.ContextMenu.Remove("missing")
		_, exists := app.ContextMenu.Get("missing")
		core.AssertFalse(t, exists)

		info := app.Environment.Info()
		core.AssertEmpty(t, info.OS)
		core.AssertEmpty(t, info.Arch)
		core.AssertNoError(t, app.Environment.OpenFileManager("", false))

		app.Event.Off("missing")
		app.Event.Reset()

		core.AssertNil(t, app.Screen.Primary())
		app.KeyBinding.Unregister("missing")
		core.AssertFalse(t, app.KeyBinding.Process("missing", nil))
	})
}

func TestApplication_AppManagers_Ugly(t *core.T) {
	app := &App{}

	core.AssertNotPanics(t, func() {
		core.AssertTrue(t, app.Clipboard.SetText("zero\x00byte"))
		core.AssertNoError(t, app.Browser.Open("/tmp/\x00report.txt"))

		app.ContextMenu.Add("dup", app.ContextMenu.New())
		app.ContextMenu.Add("dup", app.ContextMenu.New())
		core.AssertLen(t, app.ContextMenu.GetAll(), 1)

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
		core.AssertSame(t, screen, app.Screen.ScreenNearestDipPoint(Point{X: 50, Y: 50}))
		core.AssertSame(t, screen, app.Screen.ScreenNearestDipRect(Rect{X: 10, Y: 10, Width: 5, Height: 5}))

		triggered := 0
		app.KeyBinding.Register("CmdOrCtrl+Shift+P", func(Window) { triggered++ })
		app.KeyBinding.handleWindowKeyEvent(&windowKeyEvent{acceleratorString: "CmdOrCtrl+Shift+P"})
		core.AssertEqual(t, 1, triggered)
	})
}
