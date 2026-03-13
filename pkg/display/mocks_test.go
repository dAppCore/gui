package display

// mockEventSource implements EventSource for testing.
type mockEventSource struct {
	themeHandlers []func(isDark bool)
}

func newMockEventSource() *mockEventSource {
	return &mockEventSource{}
}

func (m *mockEventSource) OnThemeChange(handler func(isDark bool)) func() {
	m.themeHandlers = append(m.themeHandlers, handler)
	return func() {}
}

func (m *mockEventSource) Emit(name string, data ...any) bool {
	return true
}
