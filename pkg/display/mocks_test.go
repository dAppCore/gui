package display

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// mockApp is a mock implementation of the App interface for testing.
type mockApp struct {
	dialogManager *mockDialogManager
	envManager    *mockEnvManager
	eventManager  *mockEventManager
	logger        *mockLogger
	quitCalled    bool
}

func newMockApp() *mockApp {
	return &mockApp{
		dialogManager: newMockDialogManager(),
		envManager:    newMockEnvManager(),
		eventManager:  newMockEventManager(),
		logger:        &mockLogger{},
	}
}

func (m *mockApp) Dialog() DialogManager { return m.dialogManager }
func (m *mockApp) Env() EnvManager       { return m.envManager }
func (m *mockApp) Event() EventManager   { return m.eventManager }
func (m *mockApp) Logger() Logger        { return m.logger }
func (m *mockApp) Quit()                 { m.quitCalled = true }

// mockDialogManager tracks dialog creation calls.
type mockDialogManager struct {
	infoDialogsCreated    int
	warningDialogsCreated int
}

func newMockDialogManager() *mockDialogManager {
	return &mockDialogManager{}
}

func (m *mockDialogManager) Info() *application.MessageDialog {
	m.infoDialogsCreated++
	return nil // Can't create real dialog without Wails runtime
}

func (m *mockDialogManager) Warning() *application.MessageDialog {
	m.warningDialogsCreated++
	return nil // Can't create real dialog without Wails runtime
}

func (m *mockDialogManager) OpenFile() *application.OpenFileDialogStruct {
	return nil // Can't create real dialog without Wails runtime
}

// mockEnvManager provides mock environment info.
type mockEnvManager struct {
	envInfo  application.EnvironmentInfo
	darkMode bool
}

func newMockEnvManager() *mockEnvManager {
	return &mockEnvManager{
		envInfo: application.EnvironmentInfo{
			OS:           "test-os",
			Arch:         "test-arch",
			Debug:        true,
			PlatformInfo: map[string]any{"test": "value"},
		},
		darkMode: false,
	}
}

func (m *mockEnvManager) Info() application.EnvironmentInfo {
	return m.envInfo
}

func (m *mockEnvManager) IsDarkMode() bool {
	return m.darkMode
}

// mockEventManager tracks event registration.
type mockEventManager struct {
	registeredEvents []events.ApplicationEventType
	emittedEvents    []string
}

func newMockEventManager() *mockEventManager {
	return &mockEventManager{
		registeredEvents: make([]events.ApplicationEventType, 0),
		emittedEvents:    make([]string, 0),
	}
}

func (m *mockEventManager) OnApplicationEvent(eventType events.ApplicationEventType, handler func(*application.ApplicationEvent)) func() {
	m.registeredEvents = append(m.registeredEvents, eventType)
	return func() {} // Return a no-op unsubscribe function
}

func (m *mockEventManager) Emit(name string, data ...any) bool {
	m.emittedEvents = append(m.emittedEvents, name)
	return true // Pretend emission succeeded
}

// mockLogger tracks log calls.
type mockLogger struct {
	infoMessages []string
}

func (m *mockLogger) Info(message string, args ...any) {
	m.infoMessages = append(m.infoMessages, message)
}

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
