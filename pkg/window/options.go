// pkg/window/options.go
package window

// WindowOption is a functional option applied to a Window descriptor.
type WindowOption func(*Window) error

// ApplyOptions creates a Window and applies all options in order.
// Use: w, err := window.ApplyOptions(window.WithName("editor"), window.WithURL("/editor"))
func ApplyOptions(opts ...WindowOption) (*Window, error) {
	w := &Window{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

// WithName sets the window name.
// Use: window.WithName("editor")
func WithName(name string) WindowOption {
	return func(w *Window) error { w.Name = name; return nil }
}

// WithTitle sets the window title.
// Use: window.WithTitle("Core Editor")
func WithTitle(title string) WindowOption {
	return func(w *Window) error { w.Title = title; return nil }
}

// WithURL sets the initial window URL.
// Use: window.WithURL("/editor")
func WithURL(url string) WindowOption {
	return func(w *Window) error { w.URL = url; return nil }
}

// WithSize sets the initial window size.
// Use: window.WithSize(1280, 800)
func WithSize(width, height int) WindowOption {
	return func(w *Window) error { w.Width = width; w.Height = height; return nil }
}

// WithPosition sets the initial window position.
// Use: window.WithPosition(160, 120)
func WithPosition(x, y int) WindowOption {
	return func(w *Window) error { w.X = x; w.Y = y; return nil }
}

// WithMinSize sets the minimum window size.
// Use: window.WithMinSize(640, 480)
func WithMinSize(width, height int) WindowOption {
	return func(w *Window) error { w.MinWidth = width; w.MinHeight = height; return nil }
}

// WithMaxSize sets the maximum window size.
// Use: window.WithMaxSize(1920, 1080)
func WithMaxSize(width, height int) WindowOption {
	return func(w *Window) error { w.MaxWidth = width; w.MaxHeight = height; return nil }
}

// WithFrameless toggles the native window frame.
// Use: window.WithFrameless(true)
func WithFrameless(frameless bool) WindowOption {
	return func(w *Window) error { w.Frameless = frameless; return nil }
}

// WithHidden starts the window hidden.
// Use: window.WithHidden(true)
func WithHidden(hidden bool) WindowOption {
	return func(w *Window) error { w.Hidden = hidden; return nil }
}

// WithAlwaysOnTop keeps the window above other windows.
// Use: window.WithAlwaysOnTop(true)
func WithAlwaysOnTop(alwaysOnTop bool) WindowOption {
	return func(w *Window) error { w.AlwaysOnTop = alwaysOnTop; return nil }
}

// WithBackgroundColour sets the window background colour with alpha.
// Use: window.WithBackgroundColour(0, 0, 0, 0)
func WithBackgroundColour(r, g, b, a uint8) WindowOption {
	return func(w *Window) error { w.BackgroundColour = [4]uint8{r, g, b, a}; return nil }
}

// WithFileDrop enables drag-and-drop file handling.
// Use: window.WithFileDrop(true)
func WithFileDrop(enabled bool) WindowOption {
	return func(w *Window) error { w.EnableFileDrop = enabled; return nil }
}
