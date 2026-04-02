package window

// MockPlatform is an exported mock for cross-package integration tests.
// For internal tests, use the unexported mockPlatform in mock_test.go.
type MockPlatform struct {
	Windows []*MockWindow
}

func NewMockPlatform() *MockPlatform {
	return &MockPlatform{}
}

func (m *MockPlatform) CreateWindow(opts PlatformWindowOptions) PlatformWindow {
	w := &MockWindow{
		name: opts.Name, title: opts.Title, url: opts.URL,
		width: opts.Width, height: opts.Height,
		x: opts.X, y: opts.Y,
		alwaysOnTop:     opts.AlwaysOnTop,
		backgroundColor: opts.BackgroundColour,
		visible:         !opts.Hidden,
	}
	m.Windows = append(m.Windows, w)
	return w
}

func (m *MockPlatform) GetWindows() []PlatformWindow {
	out := make([]PlatformWindow, len(m.Windows))
	for i, w := range m.Windows {
		out[i] = w
	}
	return out
}

type MockWindow struct {
	name, title, url     string
	width, height, x, y  int
	maximised, minimised bool
	focused              bool
	visible, alwaysOnTop bool
	backgroundColor      [4]uint8
	closed               bool
	eventHandlers        []func(WindowEvent)
	fileDropHandlers     []func(paths []string, targetID string)
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
func (w *MockWindow) SetBackgroundColour(r, g, b, a uint8) { w.backgroundColor = [4]uint8{r, g, b, a} }
func (w *MockWindow) SetVisibility(visible bool)           { w.visible = visible }
func (w *MockWindow) SetAlwaysOnTop(alwaysOnTop bool)      { w.alwaysOnTop = alwaysOnTop }
func (w *MockWindow) Maximise()                            { w.maximised = true; w.minimised = false; w.visible = true }
func (w *MockWindow) Restore()                             { w.maximised = false; w.minimised = false; w.visible = true }
func (w *MockWindow) Minimise()                            { w.minimised = true; w.maximised = false; w.visible = false }
func (w *MockWindow) Focus()                               { w.focused = true }
func (w *MockWindow) Close()                               { w.closed = true }
func (w *MockWindow) Show()                                { w.visible = true }
func (w *MockWindow) Hide()                                { w.visible = false }
func (w *MockWindow) Fullscreen()                          {}
func (w *MockWindow) UnFullscreen()                        {}
func (w *MockWindow) OpenDevTools()                        {}
func (w *MockWindow) CloseDevTools()                       {}
func (w *MockWindow) OnWindowEvent(handler func(WindowEvent)) {
	w.eventHandlers = append(w.eventHandlers, handler)
}
func (w *MockWindow) OnFileDrop(handler func(paths []string, targetID string)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}
