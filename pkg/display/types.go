// pkg/display/types.go
package display

// EventSource abstracts the application event system (Wails insulation for WSEventManager).
// WSEventManager receives this instead of calling application.Get() directly.
type EventSource interface {
	OnThemeChange(handler func(isDark bool)) func()
	Emit(name string, data ...any) bool
}
