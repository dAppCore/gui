// pkg/events/platform.go
package events

// Platform abstracts the custom event backend (Wails v3 EventManager).
// Emit fires an event by name with optional data arguments.
// On registers a persistent listener; returns a cancel function.
// Off removes all listeners for the named event.
// OnMultiple registers a listener that auto-deregisters after counter firings.
// Reset removes all custom event listeners.
//
//	cancel := platform.On("build:done", func(e *CustomEvent) { ... })
//	defer cancel()
//	platform.Emit("build:done", result)
type Platform interface {
	Emit(name string, data ...any) bool
	On(name string, callback func(event *CustomEvent)) func()
	Off(name string)
	OnMultiple(name string, callback func(event *CustomEvent), counter int)
	Reset()
}

// CustomEvent is the event object delivered to On/OnMultiple listeners.
type CustomEvent struct {
	Name string
	Data any
}
