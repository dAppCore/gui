package window

type mockPlatform struct {
	windows []*mockWindow
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{}
}

func (m *mockPlatform) CreateWindow(options PlatformWindowOptions) PlatformWindow {
	w := &mockWindow{
		name: options.Name, title: options.Title, url: options.URL,
		width: options.Width, height: options.Height,
		x: options.X, y: options.Y,
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
	name, title, url      string
	width, height, x, y   int
	maximised, focused    bool
	visible, alwaysOnTop  bool
	backgroundColour      [4]uint8
	closed                bool
	minimised             bool
	fullscreened          bool
	zoom                  float64
	html                  string
	lastJS                string
	flashing              bool
	printCalled           bool
	toggleFullscreenCount int
	eventHandlers         []func(WindowEvent)
	fileDropHandlers      []func(paths []string, targetID string)
}

func (w *mockWindow) Name() string                         { return w.name }
func (w *mockWindow) Title() string                        { return w.title }
func (w *mockWindow) Position() (int, int)                 { return w.x, w.y }
func (w *mockWindow) Size() (int, int)                     { return w.width, w.height }
func (w *mockWindow) IsMaximised() bool                    { return w.maximised }
func (w *mockWindow) IsFocused() bool                      { return w.focused }
func (w *mockWindow) SetTitle(title string)                { w.title = title }
func (w *mockWindow) SetPosition(x, y int)                 { w.x = x; w.y = y }
func (w *mockWindow) SetSize(width, height int)            { w.width = width; w.height = height }
func (w *mockWindow) SetBackgroundColour(r, g, b, a uint8) { w.backgroundColour = [4]uint8{r, g, b, a} }
func (w *mockWindow) SetVisibility(visible bool)           { w.visible = visible }
func (w *mockWindow) SetAlwaysOnTop(alwaysOnTop bool)      { w.alwaysOnTop = alwaysOnTop }
func (w *mockWindow) Maximise()                            { w.maximised = true }
func (w *mockWindow) Restore()                             { w.maximised = false }
func (w *mockWindow) Minimise()                            { w.minimised = true }
func (w *mockWindow) Focus()                               { w.focused = true }
func (w *mockWindow) Close()                               { w.closed = true }
func (w *mockWindow) Show()                                { w.visible = true }
func (w *mockWindow) Hide()                                { w.visible = false }
func (w *mockWindow) Fullscreen()                          { w.fullscreened = true }
func (w *mockWindow) UnFullscreen()                        { w.fullscreened = false }
func (w *mockWindow) GetZoom() float64                     { return w.zoom }
func (w *mockWindow) SetZoom(factor float64)               { w.zoom = factor }
func (w *mockWindow) ZoomIn()                              { w.zoom += 0.1 }
func (w *mockWindow) ZoomOut()                             { w.zoom -= 0.1 }
func (w *mockWindow) SetURL(url string)                    { w.url = url }
func (w *mockWindow) SetHTML(html string)                  { w.html = html }
func (w *mockWindow) ExecJS(js string)                     { w.lastJS = js }
func (w *mockWindow) GetBounds() Bounds {
	return Bounds{X: w.x, Y: w.y, Width: w.width, Height: w.height}
}
func (w *mockWindow) SetBounds(bounds Bounds) {
	w.x = bounds.X
	w.y = bounds.Y
	w.width = bounds.Width
	w.height = bounds.Height
}
func (w *mockWindow) ToggleFullscreen() { w.toggleFullscreenCount++ }
func (w *mockWindow) Print() error      { w.printCalled = true; return nil }
func (w *mockWindow) Flash(enabled bool) { w.flashing = enabled }
func (w *mockWindow) OnWindowEvent(handler func(WindowEvent)) {
	w.eventHandlers = append(w.eventHandlers, handler)
}
func (w *mockWindow) OnFileDrop(handler func(paths []string, targetID string)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}

// emit fires a test event to all registered handlers.
func (w *mockWindow) emit(e WindowEvent) {
	for _, h := range w.eventHandlers {
		h(e)
	}
}

// emitFileDrop simulates a file drop on the window.
func (w *mockWindow) emitFileDrop(paths []string, targetID string) {
	for _, h := range w.fileDropHandlers {
		h(paths, targetID)
	}
}
