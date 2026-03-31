// pkg/window/options.go
package window

// WindowOption is the compatibility layer for option-chain callers.
// Prefer a Window literal with Manager.CreateWindow.
type WindowOption func(*Window) error

// Deprecated: use Manager.CreateWindow(Window{Name: "settings", URL: "/", Width: 800, Height: 600}).
func ApplyOptions(options ...WindowOption) (*Window, error) {
	w := &Window{}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

// Compatibility helpers for callers still using option chains.
func WithName(name string) WindowOption {
	return func(w *Window) error { w.Name = name; return nil }
}

func WithTitle(title string) WindowOption {
	return func(w *Window) error { w.Title = title; return nil }
}

func WithURL(url string) WindowOption {
	return func(w *Window) error { w.URL = url; return nil }
}

func WithSize(width, height int) WindowOption {
	return func(w *Window) error { w.Width = width; w.Height = height; return nil }
}

func WithPosition(x, y int) WindowOption {
	return func(w *Window) error { w.X = x; w.Y = y; return nil }
}

func WithMinSize(width, height int) WindowOption {
	return func(w *Window) error { w.MinWidth = width; w.MinHeight = height; return nil }
}

func WithMaxSize(width, height int) WindowOption {
	return func(w *Window) error { w.MaxWidth = width; w.MaxHeight = height; return nil }
}

func WithFrameless(frameless bool) WindowOption {
	return func(w *Window) error { w.Frameless = frameless; return nil }
}

func WithHidden(hidden bool) WindowOption {
	return func(w *Window) error { w.Hidden = hidden; return nil }
}

func WithAlwaysOnTop(alwaysOnTop bool) WindowOption {
	return func(w *Window) error { w.AlwaysOnTop = alwaysOnTop; return nil }
}

func WithBackgroundColour(r, g, b, a uint8) WindowOption {
	return func(w *Window) error { w.BackgroundColour = [4]uint8{r, g, b, a}; return nil }
}

func WithFileDrop(enabled bool) WindowOption {
	return func(w *Window) error { w.EnableFileDrop = enabled; return nil }
}
