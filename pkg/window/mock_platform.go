// pkg/window/mock_platform.go
package window

// MockPlatform is an exported mock for cross-package integration tests.
// Use: platform := window.NewMockPlatform()
// For internal tests, use the unexported mockPlatform in mock_test.go.
type MockPlatform struct {
	Windows []*MockWindow
}

// NewMockPlatform creates a window platform mock.
// Use: platform := window.NewMockPlatform()
func NewMockPlatform() *MockPlatform {
	return &MockPlatform{}
}

func (m *MockPlatform) CreateWindow(options PlatformWindowOptions) PlatformWindow {
	w := &MockWindow{
		name: options.Name, title: options.Title, url: options.URL,
		width: options.Width, height: options.Height,
		x: options.X, y: options.Y,
		visible: !options.Hidden,
	}
	m.Windows = append(m.Windows, w)
	return w
}

// GetWindows returns all tracked mock windows.
// Use: windows := platform.GetWindows()
func (m *MockPlatform) GetWindows() []PlatformWindow {
	out := make([]PlatformWindow, len(m.Windows))
	for i, w := range m.Windows {
		out[i] = w
	}
	return out
}

// MockWindow is an in-memory window handle used by tests.
// Use: w := &window.MockWindow{}
type MockWindow struct {
	name, title, url      string
	width, height, x, y   int
	maximised, minimised  bool
	focused, devtoolsOpen bool
	visible, alwaysOnTop  bool
	backgroundColour      [4]uint8
	closed                bool
	opacity               float32
	zoom                  float64
	html                  string
	lastJS                string
	flashing              bool
	printCalled           bool
	toggleFullscreenCount int
	eventHandlers         []func(WindowEvent)
	fileDropHandlers      []func(paths []string, targetID string)
}

func (w *MockWindow) Name() string                         { return w.name }
func (w *MockWindow) Title() string                        { return w.title }
func (w *MockWindow) Position() (int, int)                 { return w.x, w.y }
func (w *MockWindow) Size() (int, int)                     { return w.width, w.height }
func (w *MockWindow) IsVisible() bool                      { return w.visible }
func (w *MockWindow) IsMinimised() bool                    { return w.minimised }
func (w *MockWindow) IsMaximised() bool                    { return w.maximised }
func (w *MockWindow) IsFocused() bool                      { return w.focused }
func (w *MockWindow) SetTitle(title string)                { w.title = title }
func (w *MockWindow) SetPosition(x, y int)                 { w.x = x; w.y = y }
func (w *MockWindow) SetSize(width, height int)            { w.width = width; w.height = height }
func (w *MockWindow) SetBackgroundColour(r, g, b, a uint8) { w.backgroundColour = [4]uint8{r, g, b, a} }
func (w *MockWindow) SetOpacity(opacity float32)           { w.opacity = opacity }
func (w *MockWindow) SetVisibility(visible bool) {
	w.visible = visible
	if visible {
		w.minimised = false
	}
}
func (w *MockWindow) SetAlwaysOnTop(alwaysOnTop bool) { w.alwaysOnTop = alwaysOnTop }
func (w *MockWindow) Maximise()                       { w.maximised = true }
func (w *MockWindow) Restore()                        { w.maximised = false; w.minimised = false; w.visible = true }
func (w *MockWindow) Minimise()                       { w.minimised = true; w.visible = false }
func (w *MockWindow) Focus()                          { w.focused = true }
func (w *MockWindow) Close()                          { w.closed = true }
func (w *MockWindow) Show()                           { w.visible = true; w.minimised = false }
func (w *MockWindow) Hide()                           { w.visible = false }
func (w *MockWindow) Fullscreen()                     {}
func (w *MockWindow) UnFullscreen()                   {}
func (w *MockWindow) OpenDevTools()                   { w.devtoolsOpen = true }
func (w *MockWindow) CloseDevTools()                  { w.devtoolsOpen = false }
func (w *MockWindow) GetZoom() float64                { return w.zoom }
func (w *MockWindow) SetZoom(factor float64)          { w.zoom = factor }
func (w *MockWindow) ZoomIn()                         { w.zoom += 0.1 }
func (w *MockWindow) ZoomOut()                        { w.zoom -= 0.1 }
func (w *MockWindow) SetURL(url string)               { w.url = url }
func (w *MockWindow) SetHTML(html string)             { w.html = html }
func (w *MockWindow) ExecJS(js string)                { w.lastJS = js }
func (w *MockWindow) GetBounds() Bounds {
	return Bounds{X: w.x, Y: w.y, Width: w.width, Height: w.height}
}
func (w *MockWindow) SetBounds(bounds Bounds) {
	w.x = bounds.X
	w.y = bounds.Y
	w.width = bounds.Width
	w.height = bounds.Height
}
func (w *MockWindow) ToggleFullscreen()  { w.toggleFullscreenCount++ }
func (w *MockWindow) Print() error       { w.printCalled = true; return nil }
func (w *MockWindow) Flash(enabled bool) { w.flashing = enabled }
func (w *MockWindow) OnWindowEvent(handler func(WindowEvent)) {
	w.eventHandlers = append(w.eventHandlers, handler)
}
func (w *MockWindow) OnFileDrop(handler func(paths []string, targetID string)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}

func (w *MockWindow) emit(e WindowEvent) {
	for _, h := range w.eventHandlers {
		h(e)
	}
}

func (w *MockWindow) emitFileDrop(paths []string, targetID string) {
	for _, h := range w.fileDropHandlers {
		h(paths, targetID)
	}
}
