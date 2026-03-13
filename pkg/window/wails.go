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

func (wp *WailsPlatform) CreateWindow(opts PlatformWindowOptions) PlatformWindow {
	wOpts := application.WebviewWindowOptions{
		Name:              opts.Name,
		Title:             opts.Title,
		URL:               opts.URL,
		Width:             opts.Width,
		Height:            opts.Height,
		X:                 opts.X,
		Y:                 opts.Y,
		MinWidth:          opts.MinWidth,
		MinHeight:         opts.MinHeight,
		MaxWidth:          opts.MaxWidth,
		MaxHeight:         opts.MaxHeight,
		Frameless:         opts.Frameless,
		Hidden:            opts.Hidden,
		AlwaysOnTop:       opts.AlwaysOnTop,
		DisableResize:  opts.DisableResize,
		EnableFileDrop: opts.EnableFileDrop,
		BackgroundColour:  application.NewRGBA(opts.BackgroundColour[0], opts.BackgroundColour[1], opts.BackgroundColour[2], opts.BackgroundColour[3]),
	}
	w := wp.app.Window.NewWithOptions(wOpts)
	return &wailsWindow{w: w}
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
type wailsWindow struct {
	w *application.WebviewWindow
}

func (ww *wailsWindow) Name() string             { return ww.w.Name() }
func (ww *wailsWindow) Position() (int, int)     { return ww.w.Position() }
func (ww *wailsWindow) Size() (int, int)         { return ww.w.Size() }
func (ww *wailsWindow) IsMaximised() bool        { return ww.w.IsMaximised() }
func (ww *wailsWindow) IsFocused() bool          { return ww.w.IsFocused() }
func (ww *wailsWindow) SetTitle(title string)    { ww.w.SetTitle(title) }
func (ww *wailsWindow) SetPosition(x, y int)     { ww.w.SetPosition(x, y) }
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
func (ww *wailsWindow) Fullscreen()                     { ww.w.Fullscreen() }
func (ww *wailsWindow) UnFullscreen()                   { ww.w.UnFullscreen() }

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

// Ensure wailsWindow satisfies PlatformWindow at compile time.
var _ PlatformWindow = (*wailsWindow)(nil)

// Ensure WailsPlatform satisfies Platform at compile time.
var _ Platform = (*WailsPlatform)(nil)

