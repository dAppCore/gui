package window

import (
	core "dappco.re/go"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"reflect"
	"unsafe"
)

func TestWailsPlatform_CreateWindow_Good(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)

	w := platform.CreateWindow(PlatformWindowOptions{
		Name:             "main",
		Title:            "Core GUI",
		URL:              "/home",
		HTML:             "<main>Ready</main>",
		JS:               "globalThis.ready = true",
		Width:            1280,
		Height:           800,
		X:                10,
		Y:                20,
		MinWidth:         640,
		MinHeight:        480,
		MaxWidth:         1920,
		MaxHeight:        1080,
		Frameless:        true,
		Hidden:           true,
		AlwaysOnTop:      true,
		DisableResize:    true,
		EnableFileDrop:   true,
		BackgroundColour: [4]uint8{1, 2, 3, 4},
	})
	core.AssertNotNil(t, w)
	wails, ok := w.(*wailsWindow)
	core.RequireTrue(t, ok)

	core.AssertEqual(t, "main", wails.Name())
	core.AssertEqual(t, "Core GUI", wails.Title())
	x, y := wails.Position()
	core.AssertEqual(t, 10, x)
	core.AssertEqual(t, 20, y)

	underlying := app.Window.GetAll()[0].(*application.WebviewWindow)
	core.AssertEqual(t, "main", underlying.Name())
	core.AssertEqual(t, "Core GUI", underlying.Title())
	core.AssertEqual(t, 1280, underlying.Width())
	core.AssertEqual(t, 800, underlying.Height())
	core.AssertFalse(t, underlying.IsVisible())

	wails.SetTitle("Updated")
	wails.SetPosition(30, 40)
	wails.SetSize(1920, 1080)
	wails.SetBackgroundColour(10, 20, 30, 40)
	wails.SetVisibility(true)
	wails.SetVisibility(false)
	wails.SetAlwaysOnTop(false)
	wails.SetOpacity(0.85)
	wails.SetBounds(1, 2, 3, 4)
	wails.SetURL("/dashboard")
	wails.SetHTML("<main>Updated</main>")
	wails.SetZoom(1.25)
	wails.SetContentProtection(true)
	wails.Maximise()
	wails.Restore()
	wails.Minimise()
	wails.Focus()
	wails.Close()
	wails.Show()
	wails.Hide()
	wails.Fullscreen()
	wails.UnFullscreen()
	wails.ToggleFullscreen()
	wails.ToggleMaximise()
	wails.ExecJS("alert(1)")
	wails.Flash(true)
	wails.OpenDevTools()
	wails.CloseDevTools()
	core.RequireNoError(t, wails.Print())

	x, y = underlying.Position()
	core.AssertEqual(t, 1, x)
	core.AssertEqual(t, 2, y)
	width, height := underlying.Size()
	core.AssertEqual(t, 3, width)
	core.AssertEqual(t, 4, height)
	core.AssertTrue(t, underlying.IsMaximised())
	core.AssertTrue(t, underlying.IsFullscreen())
	core.AssertTrue(t, underlying.IsFocused())
	core.AssertFalse(t, underlying.IsVisible())
	core.AssertFalse(t, underlying.IsMinimised())
	core.AssertEqual(t, 0.85, wails.GetOpacity())
	execJSField := reflect.ValueOf(underlying).Elem().FieldByName("execJSCalls")
	core.RequireTrue(t, execJSField.IsValid())
	execJSCalls := reflect.NewAt(execJSField.Type(), unsafe.Pointer(execJSField.UnsafeAddr())).Elem().Interface().([]string)
	core.AssertEqual(t, []string{"globalThis.ready = true", "alert(1)"}, execJSCalls)

	handlers := reflect.ValueOf(underlying).Elem().FieldByName("eventHandlers")
	core.RequireTrue(t, handlers.IsValid())
	core.AssertEmpty(t, handlers.Len())

	var eventsSeen []WindowEvent
	wails.OnWindowEvent(func(event WindowEvent) {
		eventsSeen = append(eventsSeen, event)
	})

	handlers = reflect.ValueOf(underlying).Elem().FieldByName("eventHandlers")
	core.AssertEqual(t, 5, handlers.Len())
	handlerMap := reflect.NewAt(handlers.Type(), unsafe.Pointer(handlers.UnsafeAddr())).Elem().Interface().(map[events.WindowEventType][]func(*application.WindowEvent))
	moveHandlers := handlerMap[events.Common.WindowDidMove]
	core.AssertGreater(t, len(moveHandlers), 0)
	wails.SetPosition(77, 88)
	moveHandlers[0](&application.WindowEvent{})

	resizeHandlers := handlerMap[events.Common.WindowDidResize]
	core.AssertGreater(t, len(resizeHandlers), 0)
	wails.SetSize(640, 360)
	resizeHandlers[0](&application.WindowEvent{})

	core.AssertLen(t, eventsSeen, 2)
	core.AssertEqual(t, "move", eventsSeen[0].Type)
	core.AssertEqual(t, "main", eventsSeen[0].Name)
	core.AssertEqual(t, 77, eventsSeen[0].Data["x"])
	core.AssertEqual(t, 88, eventsSeen[0].Data["y"])
	core.AssertEqual(t, "resize", eventsSeen[1].Type)
	core.AssertEqual(t, 640, eventsSeen[1].Data["width"])
	core.AssertEqual(t, 360, eventsSeen[1].Data["height"])

	var filesSeen []string
	var targetSeen string
	wails.OnFileDrop(func(paths []string, targetID string) {
		filesSeen = append(filesSeen, paths...)
		targetSeen = targetID
	})

	dropHandlers := reflect.ValueOf(underlying).Elem().FieldByName("eventHandlers")
	dropHandlerMap := reflect.NewAt(dropHandlers.Type(), unsafe.Pointer(dropHandlers.UnsafeAddr())).Elem().Interface().(map[events.WindowEventType][]func(*application.WindowEvent))
	fileDropHandlers := dropHandlerMap[events.Common.WindowFilesDropped]
	core.AssertGreater(t, len(fileDropHandlers), 0)

	event := &application.WindowEvent{}
	ctx := event.Context()
	ctxValue := reflect.ValueOf(ctx).Elem()
	filesField := ctxValue.FieldByName("droppedFiles")
	reflect.NewAt(filesField.Type(), unsafe.Pointer(filesField.UnsafeAddr())).Elem().Set(reflect.ValueOf([]string{"a.txt", "b.txt"}))
	detailsField := ctxValue.FieldByName("dropDetails")
	reflect.NewAt(detailsField.Type(), unsafe.Pointer(detailsField.UnsafeAddr())).Elem().Set(reflect.ValueOf(&application.DropTargetDetails{ElementID: "drop-zone"}))
	fileDropHandlers[0](event)

	core.AssertEqual(t, []string{"a.txt", "b.txt"}, filesSeen)
	core.AssertEqual(t, "drop-zone", targetSeen)
}

func TestWailsPlatform_GetWindows_Bad(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	core.AssertEmpty(t, platform.GetWindows())
}
