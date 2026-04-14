// pkg/events/messages.go
package events

// TaskEmit fires a custom event by name with optional data.
// c.PERFORM(events.TaskEmit{Name: "build:done", Data: result})
type TaskEmit struct {
	Name string
	Data any
}

// TaskOn registers a persistent listener for a named event.
// The listener ID returned in the result can be used with TaskOff.
// c.PERFORM(events.TaskOn{Name: "build:done"})
type TaskOn struct {
	Name string
}

// TaskOff removes all listeners for a named event.
// c.PERFORM(events.TaskOff{Name: "build:done"})
type TaskOff struct {
	Name string
}

// QueryListeners returns the count of listeners registered for a named event.
// count := c.QUERY(events.QueryListeners{Name: "build:done"})
type QueryListeners struct {
	Name string
}

// ActionEventFired is broadcast when a custom event is emitted via TaskEmit.
type ActionEventFired struct {
	Name string
	Data any
}
