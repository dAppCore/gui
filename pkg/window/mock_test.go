// pkg/window/mock_test.go
package window

type mockPlatform struct {
	windows []*mockWindow
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{}
}

func (m *mockPlatform) CreateWindow(opts PlatformWindowOptions) PlatformWindow {
	w := &mockWindow{
		name: opts.Name, title: opts.Title, url: opts.URL,
		width: opts.Width, height: opts.Height,
		x: opts.X, y: opts.Y,
	}
	m.windows = append(m.windows, w)
	return w
}

func (m *mockPlatform) GetWindows() []PlatformWindow {
	out := make([]PlatformWindow, len(m.windows))
	for i, w := range m.windows {
		out[i] = w
	}
	return out
}

type mockWindow struct {
	name, title, url     string
	width, height, x, y  int
	maximised, focused   bool
	visible, alwaysOnTop bool
	closed               bool
	eventHandlers        []func(WindowEvent)
}

func (w *mockWindow) Name() string                            { return w.name }
func (w *mockWindow) Position() (int, int)                    { return w.x, w.y }
func (w *mockWindow) Size() (int, int)                        { return w.width, w.height }
func (w *mockWindow) IsMaximised() bool                       { return w.maximised }
func (w *mockWindow) IsFocused() bool                         { return w.focused }
func (w *mockWindow) SetTitle(title string)                   { w.title = title }
func (w *mockWindow) SetPosition(x, y int)                    { w.x = x; w.y = y }
func (w *mockWindow) SetSize(width, height int)               { w.width = width; w.height = height }
func (w *mockWindow) SetBackgroundColour(r, g, b, a uint8)    {}
func (w *mockWindow) SetVisibility(visible bool)              { w.visible = visible }
func (w *mockWindow) SetAlwaysOnTop(alwaysOnTop bool)         { w.alwaysOnTop = alwaysOnTop }
func (w *mockWindow) Maximise()                               { w.maximised = true }
func (w *mockWindow) Restore()                                { w.maximised = false }
func (w *mockWindow) Minimise()                               {}
func (w *mockWindow) Focus()                                  { w.focused = true }
func (w *mockWindow) Close()                                  { w.closed = true }
func (w *mockWindow) Show()                                   { w.visible = true }
func (w *mockWindow) Hide()                                   { w.visible = false }
func (w *mockWindow) Fullscreen()                             {}
func (w *mockWindow) UnFullscreen()                           {}
func (w *mockWindow) OnWindowEvent(handler func(WindowEvent)) { w.eventHandlers = append(w.eventHandlers, handler) }

// emit fires a test event to all registered handlers.
func (w *mockWindow) emit(e WindowEvent) {
	for _, h := range w.eventHandlers {
		h(e)
	}
}
