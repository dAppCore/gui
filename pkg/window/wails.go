// pkg/window/wails.go
package window

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// WailsPlatform implements Platform using Wails v3.
type WailsPlatform struct {
	app *application.App
}

// NewWailsPlatform creates a Wails-backed Platform.
func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) CreateWindow(options PlatformWindowOptions) PlatformWindow {
	wOpts := application.WebviewWindowOptions{
		Name:             options.Name,
		Title:            options.Title,
		URL:              options.URL,
		Width:            options.Width,
		Height:           options.Height,
		X:                options.X,
		Y:                options.Y,
		MinWidth:         options.MinWidth,
		MinHeight:        options.MinHeight,
		MaxWidth:         options.MaxWidth,
		MaxHeight:        options.MaxHeight,
		Frameless:        options.Frameless,
		Hidden:           options.Hidden,
		AlwaysOnTop:      options.AlwaysOnTop,
		DisableResize:    options.DisableResize,
		EnableFileDrop:   options.EnableFileDrop,
		BackgroundColour: application.NewRGBA(options.BackgroundColour[0], options.BackgroundColour[1], options.BackgroundColour[2], options.BackgroundColour[3]),
	}
	w := wp.app.Window.NewWithOptions(wOpts)
	return &wailsWindow{w: w, title: options.Title}
}

func (wp *WailsPlatform) GetWindows() []PlatformWindow {
	all := wp.app.Window.GetAll()
	out := make([]PlatformWindow, 0, len(all))
	for _, w := range all {
		if wv, ok := w.(*application.WebviewWindow); ok {
			out = append(out, &wailsWindow{w: wv})
		}
	}
	return out
}

// wailsWindow wraps *application.WebviewWindow to implement PlatformWindow.
// It stores the title locally because Wails v3 does not expose a title getter.
type wailsWindow struct {
	w     *application.WebviewWindow
	title string
}

func (ww *wailsWindow) Name() string              { return ww.w.Name() }
func (ww *wailsWindow) Title() string             { return ww.title }
func (ww *wailsWindow) Position() (int, int)      { return ww.w.Position() }
func (ww *wailsWindow) Size() (int, int)          { return ww.w.Size() }
func (ww *wailsWindow) IsMaximised() bool         { return ww.w.IsMaximised() }
func (ww *wailsWindow) IsFocused() bool           { return ww.w.IsFocused() }
func (ww *wailsWindow) SetTitle(title string)     { ww.title = title; ww.w.SetTitle(title) }
func (ww *wailsWindow) SetPosition(x, y int)      { ww.w.SetPosition(x, y) }
func (ww *wailsWindow) SetSize(width, height int) { ww.w.SetSize(width, height) }
func (ww *wailsWindow) SetBackgroundColour(r, g, b, a uint8) {
	ww.w.SetBackgroundColour(application.NewRGBA(r, g, b, a))
}
func (ww *wailsWindow) SetVisibility(visible bool) {
	if visible {
		ww.w.Show()
	} else {
		ww.w.Hide()
	}
}
func (ww *wailsWindow) SetAlwaysOnTop(alwaysOnTop bool) { ww.w.SetAlwaysOnTop(alwaysOnTop) }
func (ww *wailsWindow) Maximise()                       { ww.w.Maximise() }
func (ww *wailsWindow) Restore()                        { ww.w.Restore() }
func (ww *wailsWindow) Minimise()                       { ww.w.Minimise() }
func (ww *wailsWindow) Focus()                          { ww.w.Focus() }
func (ww *wailsWindow) Close()                          { ww.w.Close() }
func (ww *wailsWindow) Show()                           { ww.w.Show() }
func (ww *wailsWindow) Hide()                           { ww.w.Hide() }
func (ww *wailsWindow) Fullscreen()         { ww.w.Fullscreen() }
func (ww *wailsWindow) UnFullscreen()       { ww.w.UnFullscreen() }
func (ww *wailsWindow) GetZoom() float64    { return ww.w.GetZoom() }
func (ww *wailsWindow) SetZoom(factor float64) { ww.w.SetZoom(factor) }
func (ww *wailsWindow) ZoomIn()             { ww.w.ZoomIn() }
func (ww *wailsWindow) ZoomOut()            { ww.w.ZoomOut() }
func (ww *wailsWindow) SetURL(url string)   { ww.w.SetURL(url) }
func (ww *wailsWindow) SetHTML(html string) { ww.w.SetHTML(html) }
func (ww *wailsWindow) ExecJS(js string)    { ww.w.ExecJS(js) }
func (ww *wailsWindow) GetBounds() Bounds {
	r := ww.w.Bounds()
	return Bounds{X: r.X, Y: r.Y, Width: r.Width, Height: r.Height}
}
func (ww *wailsWindow) SetBounds(b Bounds) {
	ww.w.SetBounds(application.Rect{X: b.X, Y: b.Y, Width: b.Width, Height: b.Height})
}
func (ww *wailsWindow) ToggleFullscreen() { ww.w.ToggleFullscreen() }
func (ww *wailsWindow) Print() error      { return ww.w.Print() }
func (ww *wailsWindow) Flash(enabled bool) { ww.w.Flash(enabled) }

func (ww *wailsWindow) OnWindowEvent(handler func(event WindowEvent)) {
	name := ww.w.Name()

	// Map common Wails window events to our WindowEvent type.
	eventMap := map[events.WindowEventType]string{
		events.Common.WindowFocus:     "focus",
		events.Common.WindowLostFocus: "blur",
		events.Common.WindowDidMove:   "move",
		events.Common.WindowDidResize: "resize",
		events.Common.WindowClosing:   "close",
	}

	for eventType, eventName := range eventMap {
		typeName := eventName // capture for closure
		ww.w.OnWindowEvent(eventType, func(event *application.WindowEvent) {
			data := make(map[string]any)
			switch typeName {
			case "move":
				x, y := ww.w.Position()
				data["x"] = x
				data["y"] = y
			case "resize":
				w, h := ww.w.Size()
				data["width"] = w
				data["height"] = h
			}
			handler(WindowEvent{
				Type: typeName,
				Name: name,
				Data: data,
			})
		})
	}
}

func (ww *wailsWindow) OnFileDrop(handler func(paths []string, targetID string)) {
	ww.w.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		details := event.Context().DropTargetDetails()
		targetID := ""
		if details != nil {
			targetID = details.ElementID
		}
		handler(files, targetID)
	})
}

// Ensure wailsWindow satisfies PlatformWindow at compile time.
var _ PlatformWindow = (*wailsWindow)(nil)

// Ensure WailsPlatform satisfies Platform at compile time.
var _ Platform = (*WailsPlatform)(nil)
