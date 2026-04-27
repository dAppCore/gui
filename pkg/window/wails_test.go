package window

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestWailsPlatform_CreateWindow_Good(t *testing.T) {
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
	require.NotNil(t, w)
	wails, ok := w.(*wailsWindow)
	require.True(t, ok)

	assert.Equal(t, "main", wails.Name())
	assert.Equal(t, "Core GUI", wails.Title())
	x, y := wails.Position()
	assert.Equal(t, 10, x)
	assert.Equal(t, 20, y)

	underlying := app.Window.GetAll()[0].(*application.WebviewWindow)
	assert.Equal(t, "main", underlying.Name())
	assert.Equal(t, "Core GUI", underlying.Title())
	assert.Equal(t, 1280, underlying.Width())
	assert.Equal(t, 800, underlying.Height())
	assert.False(t, underlying.IsVisible())

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
	require.NoError(t, wails.Print())

	x, y = underlying.Position()
	assert.Equal(t, 1, x)
	assert.Equal(t, 2, y)
	width, height := underlying.Size()
	assert.Equal(t, 3, width)
	assert.Equal(t, 4, height)
	assert.True(t, underlying.IsMaximised())
	assert.True(t, underlying.IsFullscreen())
	assert.True(t, underlying.IsFocused())
	assert.False(t, underlying.IsVisible())
	assert.False(t, underlying.IsMinimised())
	assert.Equal(t, 0.85, wails.GetOpacity())
	execJSField := reflect.ValueOf(underlying).Elem().FieldByName("execJSCalls")
	require.True(t, execJSField.IsValid())
	execJSCalls := reflect.NewAt(execJSField.Type(), unsafe.Pointer(execJSField.UnsafeAddr())).Elem().Interface().([]string)
	assert.Equal(t, []string{"globalThis.ready = true", "alert(1)"}, execJSCalls)

	handlers := reflect.ValueOf(underlying).Elem().FieldByName("eventHandlers")
	require.True(t, handlers.IsValid())
	require.Zero(t, handlers.Len())

	var eventsSeen []WindowEvent
	wails.OnWindowEvent(func(event WindowEvent) {
		eventsSeen = append(eventsSeen, event)
	})

	handlers = reflect.ValueOf(underlying).Elem().FieldByName("eventHandlers")
	require.Equal(t, 5, handlers.Len())
	handlerMap := reflect.NewAt(handlers.Type(), unsafe.Pointer(handlers.UnsafeAddr())).Elem().Interface().(map[events.WindowEventType][]func(*application.WindowEvent))
	moveHandlers := handlerMap[events.Common.WindowDidMove]
	require.Greater(t, len(moveHandlers), 0)
	wails.SetPosition(77, 88)
	moveHandlers[0](&application.WindowEvent{})

	resizeHandlers := handlerMap[events.Common.WindowDidResize]
	require.Greater(t, len(resizeHandlers), 0)
	wails.SetSize(640, 360)
	resizeHandlers[0](&application.WindowEvent{})

	require.Len(t, eventsSeen, 2)
	assert.Equal(t, "move", eventsSeen[0].Type)
	assert.Equal(t, "main", eventsSeen[0].Name)
	assert.Equal(t, 77, eventsSeen[0].Data["x"])
	assert.Equal(t, 88, eventsSeen[0].Data["y"])
	assert.Equal(t, "resize", eventsSeen[1].Type)
	assert.Equal(t, 640, eventsSeen[1].Data["width"])
	assert.Equal(t, 360, eventsSeen[1].Data["height"])

	var filesSeen []string
	var targetSeen string
	wails.OnFileDrop(func(paths []string, targetID string) {
		filesSeen = append(filesSeen, paths...)
		targetSeen = targetID
	})

	dropHandlers := reflect.ValueOf(underlying).Elem().FieldByName("eventHandlers")
	dropHandlerMap := reflect.NewAt(dropHandlers.Type(), unsafe.Pointer(dropHandlers.UnsafeAddr())).Elem().Interface().(map[events.WindowEventType][]func(*application.WindowEvent))
	fileDropHandlers := dropHandlerMap[events.Common.WindowFilesDropped]
	require.Greater(t, len(fileDropHandlers), 0)

	event := &application.WindowEvent{}
	ctx := event.Context()
	ctxValue := reflect.ValueOf(ctx).Elem()
	filesField := ctxValue.FieldByName("droppedFiles")
	reflect.NewAt(filesField.Type(), unsafe.Pointer(filesField.UnsafeAddr())).Elem().Set(reflect.ValueOf([]string{"a.txt", "b.txt"}))
	detailsField := ctxValue.FieldByName("dropDetails")
	reflect.NewAt(detailsField.Type(), unsafe.Pointer(detailsField.UnsafeAddr())).Elem().Set(reflect.ValueOf(&application.DropTargetDetails{ElementID: "drop-zone"}))
	fileDropHandlers[0](event)

	assert.Equal(t, []string{"a.txt", "b.txt"}, filesSeen)
	assert.Equal(t, "drop-zone", targetSeen)
}

func TestWailsPlatform_GetWindows_Bad(t *testing.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	assert.Empty(t, platform.GetWindows())
}
