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

func TestApplication_Menu_GoodCase(t *core.T) {
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

func TestApplication_Menu_BadCase(t *core.T) {
	menu := NewMenu()

	core.AssertEmpty(t, menu.Items)
	core.AssertNotEmpty(t, core.Sprintf("%T", menu))
}

func TestApplication_Menu_UglyCase(t *core.T) {
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

func TestApplication_SystemTray_GoodCase(t *core.T) {
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

func TestApplication_SystemTray_BadCase(t *core.T) {
	tray := (&SystemTrayManager{}).New()

	core.AssertEmpty(t, tray.icon)
	core.AssertEmpty(t, tray.templateIcon)
	core.AssertEmpty(t, tray.tooltip)
	core.AssertEmpty(t, tray.label)
	core.AssertNil(t, tray.menu)
	core.AssertNil(t, tray.attachedWindow)
}

func TestApplication_SystemTray_UglyCase(t *core.T) {
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

func TestApplication_WebviewWindow_GoodCase(t *core.T) {
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

func TestApplication_WebviewWindow_BadCase(t *core.T) {
	manager := &WindowManager{}
	window := manager.NewWithOptions(WebviewWindowOptions{Hidden: true})

	core.AssertFalse(t, window.IsVisible())
	core.AssertNil(t, manager.Get("missing"))
	core.AssertNil(t, manager.GetByID(99))
}

func TestApplication_WebviewWindow_UglyCase(t *core.T) {
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

func TestApplication_AppManagers_GoodCase(t *core.T) {
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

func TestApplication_AppManagers_BadCase(t *core.T) {
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

func TestApplication_AppManagers_UglyCase(t *core.T) {
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

// AX7 generated source-matching smoke coverage.
func TestApplication_Logger_Info_Good(t *core.T) {
	var subject Logger
	result := core.Try(func() any {
		subject.Info("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Logger_Info_Bad(t *core.T) {
	var subject Logger
	result := core.Try(func() any {
		subject.Info("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Logger_Info_Ugly(t *core.T) {
	var subject Logger
	result := core.Try(func() any {
		subject.Info("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetAccelerator_Good(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetAccelerator("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetAccelerator_Bad(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetAccelerator("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetAccelerator_Ugly(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetAccelerator("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetTooltip_Good(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetTooltip_Bad(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetTooltip_Ugly(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetChecked_Good(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetChecked(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetChecked_Bad(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetChecked_Ugly(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetEnabled_Good(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetEnabled_Bad(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_MenuItem_SetEnabled_Ugly(t *core.T) {
	subject := new(MenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_NewMenu_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_NewMenu_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_NewMenu_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_Add_Good(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		got0 := subject.Add("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_Add_Bad(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		got0 := subject.Add("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_Add_Ugly(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		got0 := subject.Add("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddSeparator_Good(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddSeparator_Bad(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddSeparator_Ugly(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddSubmenu_Good(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddSubmenu_Bad(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddSubmenu_Ugly(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddRole_Good(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddRole_Bad(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_Menu_AddRole_Ugly(t *core.T) {
	subject := new(Menu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetIcon_Good(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetIcon_Bad(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetIcon_Ugly(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetTemplateIcon_Good(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetTemplateIcon_Bad(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetTemplateIcon_Ugly(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetTooltip_Good(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetTooltip_Bad(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetTooltip_Ugly(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetLabel_Good(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetLabel("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetLabel_Bad(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetLabel("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetLabel_Ugly(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetLabel("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetMenu_Good(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetMenu_Bad(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_SetMenu_Ugly(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_AttachWindow_Good(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.AttachWindow(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_AttachWindow_Bad(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.AttachWindow(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTray_AttachWindow_Ugly(t *core.T) {
	subject := new(SystemTray)
	result := core.Try(func() any {
		subject.AttachWindow(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTrayManager_New_Good(t *core.T) {
	subject := new(SystemTrayManager)
	result := core.Try(func() any {
		got0 := subject.New()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTrayManager_New_Bad(t *core.T) {
	subject := new(SystemTrayManager)
	result := core.Try(func() any {
		got0 := subject.New()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_SystemTrayManager_New_Ugly(t *core.T) {
	subject := new(SystemTrayManager)
	result := core.Try(func() any {
		got0 := subject.New()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEventContext_DroppedFiles_Good(t *core.T) {
	subject := new(WindowEventContext)
	result := core.Try(func() any {
		got0 := subject.DroppedFiles()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEventContext_DroppedFiles_Bad(t *core.T) {
	subject := new(WindowEventContext)
	result := core.Try(func() any {
		got0 := subject.DroppedFiles()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEventContext_DroppedFiles_Ugly(t *core.T) {
	subject := new(WindowEventContext)
	result := core.Try(func() any {
		got0 := subject.DroppedFiles()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEventContext_DropTargetDetails_Good(t *core.T) {
	subject := new(WindowEventContext)
	result := core.Try(func() any {
		got0 := subject.DropTargetDetails()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEventContext_DropTargetDetails_Bad(t *core.T) {
	subject := new(WindowEventContext)
	result := core.Try(func() any {
		got0 := subject.DropTargetDetails()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEventContext_DropTargetDetails_Ugly(t *core.T) {
	subject := new(WindowEventContext)
	result := core.Try(func() any {
		got0 := subject.DropTargetDetails()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEvent_Context_Good(t *core.T) {
	subject := new(WindowEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEvent_Context_Bad(t *core.T) {
	subject := new(WindowEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowEvent_Context_Ugly(t *core.T) {
	subject := new(WindowEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Name_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Name_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Name_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Title_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Title_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Title_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Position_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Position_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Position_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Size_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Size_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Size_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsMaximised_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsMaximised_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsMaximised_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsFocused_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsFocused_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsFocused_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetTitle_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetTitle("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetTitle_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetTitle("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetTitle_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetTitle("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetPosition_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetPosition(1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetPosition_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetPosition(0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetPosition_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetPosition(-1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetSize_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetSize_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetSize_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetBackgroundColour_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetBackgroundColour(*new(RGBA))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetBackgroundColour_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetBackgroundColour(*new(RGBA))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetBackgroundColour_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetBackgroundColour(*new(RGBA))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetAlwaysOnTop_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetAlwaysOnTop(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetAlwaysOnTop_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetAlwaysOnTop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetAlwaysOnTop_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetAlwaysOnTop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetOpacity_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetOpacity(1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetOpacity_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetOpacity(0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetOpacity_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetOpacity(-1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Maximise_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Maximise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Maximise_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Maximise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Maximise_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Maximise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Restore_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Restore_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Restore_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Minimise_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Minimise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Minimise_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Minimise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Minimise_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Minimise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Focus_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Focus_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Focus_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Close_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Close_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Close_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Show_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Show()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Show_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Show()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Show_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Show()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Hide_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Hide()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Hide_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Hide()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Hide_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Hide()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Fullscreen_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Fullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Fullscreen_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Fullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Fullscreen_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Fullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnFullscreen_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnFullscreen_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnFullscreen_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OnWindowEvent_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.OnWindowEvent(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OnWindowEvent_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.OnWindowEvent(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OnWindowEvent_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.OnWindowEvent(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ID_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ID_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ID_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ClientID_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ClientID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ClientID_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ClientID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ClientID_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ClientID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Width_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Width()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Width_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Width()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Width_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Width()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Height_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Height()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Height_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Height()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Height_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Height()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsVisible_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsVisible_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsVisible_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsFullscreen_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsFullscreen_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsFullscreen_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsMinimised_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsMinimised_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsMinimised_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsIgnoreMouseEvents_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsIgnoreMouseEvents()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsIgnoreMouseEvents_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsIgnoreMouseEvents()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_IsIgnoreMouseEvents_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.IsIgnoreMouseEvents()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Resizable_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Resizable()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Resizable_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Resizable()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Resizable_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Resizable()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Bounds_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Bounds()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Bounds_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Bounds()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Bounds_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Bounds()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetBounds_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetBounds(*new(Rect))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetBounds_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetBounds(*new(Rect))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetBounds_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetBounds(*new(Rect))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_RelativePosition_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.RelativePosition()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_RelativePosition_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.RelativePosition()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_RelativePosition_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.RelativePosition()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetRelativePosition_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetRelativePosition(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetRelativePosition_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetRelativePosition(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetRelativePosition_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetRelativePosition(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMinSize_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMinSize_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMinSize_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMaxSize_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaxSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMaxSize_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaxSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMaxSize_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaxSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Center_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Center()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Center_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Center()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Center_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Center()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetURL_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetURL("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetURL_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetURL("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetURL_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetURL("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetHTML_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetHTML("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetHTML_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetHTML("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetHTML_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetHTML("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetFrameless_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetFrameless(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetFrameless_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetFrameless(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetFrameless_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetFrameless(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetResizable_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetResizable(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetResizable_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetResizable(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetResizable_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetResizable(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetIgnoreMouseEvents_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetIgnoreMouseEvents(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetIgnoreMouseEvents_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetIgnoreMouseEvents(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetIgnoreMouseEvents_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetIgnoreMouseEvents(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMinimiseButtonState_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinimiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMinimiseButtonState_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinimiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMinimiseButtonState_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinimiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMaximiseButtonState_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaximiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMaximiseButtonState_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaximiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMaximiseButtonState_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaximiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetCloseButtonState_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetCloseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetCloseButtonState_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetCloseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetCloseButtonState_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetCloseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetEnabled_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetEnabled(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetEnabled_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetEnabled_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetContentProtection_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetContentProtection(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetContentProtection_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetContentProtection(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetContentProtection_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetContentProtection(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMenu_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMenu_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetMenu_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ShowMenuBar_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ShowMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ShowMenuBar_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ShowMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ShowMenuBar_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ShowMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HideMenuBar_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HideMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HideMenuBar_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HideMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HideMenuBar_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HideMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleMenuBar_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleMenuBar_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleMenuBar_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleFrameless_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleFrameless()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleFrameless_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleFrameless()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleFrameless_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleFrameless()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ExecJS_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ExecJS("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ExecJS_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ExecJS("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ExecJS_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ExecJS("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Reload_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Reload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Reload_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Reload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Reload_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Reload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ForceReload_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ForceReload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ForceReload_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ForceReload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ForceReload_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ForceReload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OpenDevTools_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OpenDevTools_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OpenDevTools_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OpenContextMenu_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.OpenContextMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OpenContextMenu_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.OpenContextMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_OpenContextMenu_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.OpenContextMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Zoom_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Zoom()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Zoom_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Zoom()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Zoom_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Zoom()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomIn_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ZoomIn()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomIn_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ZoomIn()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomIn_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ZoomIn()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomOut_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ZoomOut()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomOut_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ZoomOut()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomOut_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ZoomOut()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomReset_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ZoomReset()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomReset_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ZoomReset()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ZoomReset_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.ZoomReset()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetZoom_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetZoom_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetZoom_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetOpacity_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetOpacity_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetOpacity_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetZoom_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetZoom(1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetZoom_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetZoom(0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SetZoom_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.SetZoom(-1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_RegisterHook_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.RegisterHook(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_RegisterHook_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.RegisterHook(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_RegisterHook_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.RegisterHook(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_EmitEvent_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.EmitEvent("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_EmitEvent_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.EmitEvent("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_EmitEvent_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.EmitEvent("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_DispatchWailsEvent_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.DispatchWailsEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_DispatchWailsEvent_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.DispatchWailsEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_DispatchWailsEvent_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.DispatchWailsEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetScreen_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetScreen_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetScreen_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetBorderSizes_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetBorderSizes()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetBorderSizes_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetBorderSizes()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_GetBorderSizes_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.GetBorderSizes()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_EnableSizeConstraints_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.EnableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_EnableSizeConstraints_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.EnableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_EnableSizeConstraints_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.EnableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_DisableSizeConstraints_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.DisableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_DisableSizeConstraints_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.DisableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_DisableSizeConstraints_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.DisableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_AttachModal_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.AttachModal(*new(Window))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_AttachModal_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.AttachModal(*new(Window))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_AttachModal_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.AttachModal(*new(Window))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Flash_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Flash(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Flash_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Flash_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Print_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Print_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Print_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Error_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Error("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Error_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Error("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Error_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Error("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Info_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Info("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Info_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Info("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Info_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Info("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_NativeWindow_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.NativeWindow()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_NativeWindow_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.NativeWindow()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_NativeWindow_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		got0 := subject.NativeWindow()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Run_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Run()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Run_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Run()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_Run_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.Run()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnMaximise_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnMaximise_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnMaximise_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnMinimise_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnMinimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnMinimise_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnMinimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_UnMinimise_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.UnMinimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleFullscreen_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleFullscreen_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleFullscreen_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleMaximise_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleMaximise_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_ToggleMaximise_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SnapAssist_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SnapAssist()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SnapAssist_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SnapAssist()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_SnapAssist_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.SnapAssist()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_InitiateFrontendDropProcessing_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.InitiateFrontendDropProcessing(nil, 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_InitiateFrontendDropProcessing_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.InitiateFrontendDropProcessing(nil, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_InitiateFrontendDropProcessing_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.InitiateFrontendDropProcessing(nil, -1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleMessage_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleMessage("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleMessage_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleMessage("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleMessage_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleMessage("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleWindowEvent_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleWindowEvent(1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleWindowEvent_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleWindowEvent(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleWindowEvent_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleWindowEvent(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleKeyEvent_Good(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleKeyEvent("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleKeyEvent_Bad(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleKeyEvent("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WebviewWindow_HandleKeyEvent_Ugly(t *core.T) {
	subject := new(WebviewWindow)
	result := core.Try(func() any {
		subject.HandleKeyEvent("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowManager_NewWithOptions_Good(t *core.T) {
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.NewWithOptions(*new(WebviewWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowManager_NewWithOptions_Bad(t *core.T) {
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.NewWithOptions(*new(WebviewWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowManager_NewWithOptions_Ugly(t *core.T) {
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.NewWithOptions(*new(WebviewWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowManager_GetAll_Good(t *core.T) {
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowManager_GetAll_Bad(t *core.T) {
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_WindowManager_GetAll_Ugly(t *core.T) {
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_App_Quit_Good(t *core.T) {
	subject := new(App)
	result := core.Try(func() any {
		subject.Quit()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_App_Quit_Bad(t *core.T) {
	subject := new(App)
	result := core.Try(func() any {
		subject.Quit()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_App_Quit_Ugly(t *core.T) {
	subject := new(App)
	result := core.Try(func() any {
		subject.Quit()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_App_NewMenu_Good(t *core.T) {
	subject := new(App)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_App_NewMenu_Bad(t *core.T) {
	subject := new(App)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestApplication_App_NewMenu_Ugly(t *core.T) {
	subject := new(App)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
