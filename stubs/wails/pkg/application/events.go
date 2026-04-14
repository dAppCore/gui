package application

import (
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// ApplicationEventContext carries structured data for an ApplicationEvent.
//
//	files := event.Context().OpenedFiles()
//	dark := event.Context().IsDarkMode()
type ApplicationEventContext struct {
	data map[string]any
}

// OpenedFiles returns the list of files provided via the event context, or nil.
//
//	for _, path := range event.Context().OpenedFiles() { open(path) }
func (c *ApplicationEventContext) OpenedFiles() []string {
	if c.data == nil {
		return nil
	}
	files, ok := c.data["openedFiles"]
	if !ok {
		return nil
	}
	result, ok := files.([]string)
	if !ok {
		return nil
	}
	return result
}

// IsDarkMode returns true when the event context reports dark mode active.
//
//	if event.Context().IsDarkMode() { applyDark() }
func (c *ApplicationEventContext) IsDarkMode() bool {
	return c.getBool("isDarkMode")
}

// HasVisibleWindows returns true when the event context reports at least one visible window.
//
//	if event.Context().HasVisibleWindows() { ... }
func (c *ApplicationEventContext) HasVisibleWindows() bool {
	return c.getBool("hasVisibleWindows")
}

// Filename returns the filename value from the event context, or "".
//
//	path := event.Context().Filename()
func (c *ApplicationEventContext) Filename() string {
	if c.data == nil {
		return ""
	}
	v, ok := c.data["filename"]
	if !ok {
		return ""
	}
	result, ok := v.(string)
	if !ok {
		return ""
	}
	return result
}

// URL returns the URL value from the event context, or "".
//
//	url := event.Context().URL()
func (c *ApplicationEventContext) URL() string {
	if c.data == nil {
		return ""
	}
	v, ok := c.data["url"]
	if !ok {
		return ""
	}
	result, ok := v.(string)
	if !ok {
		return ""
	}
	return result
}

func (c *ApplicationEventContext) getBool(key string) bool {
	if c.data == nil {
		return false
	}
	v, ok := c.data[key]
	if !ok {
		return false
	}
	result, ok := v.(bool)
	if !ok {
		return false
	}
	return result
}

// ApplicationEvent is the event object delivered to OnApplicationEvent listeners.
//
//	em.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) {
//	    dark := e.Context().IsDarkMode()
//	})
type ApplicationEvent struct {
	Id        uint
	ctx       *ApplicationEventContext
	cancelled atomic.Bool
}

// Context returns the ApplicationEventContext attached to the event.
func (e *ApplicationEvent) Context() *ApplicationEventContext {
	if e.ctx == nil {
		e.ctx = &ApplicationEventContext{data: make(map[string]any)}
	}
	return e.ctx
}

// Cancel marks the event as cancelled, preventing further listener dispatch.
func (e *ApplicationEvent) Cancel() {
	e.cancelled.Store(true)
}

// IsCancelled reports whether the event has been cancelled.
func (e *ApplicationEvent) IsCancelled() bool {
	return e.cancelled.Load()
}

// customEventListener is an internal listener registration.
type customEventListener struct {
	callback func(*CustomEvent)
	counter  int // -1 = unlimited
}

// applicationEventListener is an internal listener registration.
type applicationEventListener struct {
	callback func(*ApplicationEvent)
}

// EventManager manages custom and application-level event subscriptions.
//
//	em.Emit("build:done", result)
//	cancel := em.On("build:done", func(e *application.CustomEvent) { ... })
//	em.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) { ... })
type EventManager struct {
	mu sync.RWMutex

	customListeners      map[string][]*customEventListener
	applicationListeners map[uint][]*applicationEventListener
}

func (em *EventManager) ensureCustomListeners() {
	if em.customListeners == nil {
		em.customListeners = make(map[string][]*customEventListener)
	}
}

func (em *EventManager) ensureApplicationListeners() {
	if em.applicationListeners == nil {
		em.applicationListeners = make(map[uint][]*applicationEventListener)
	}
}

// Emit emits a custom event by name with optional data arguments.
// Returns true if the event was cancelled by a listener.
//
//	cancelled := em.Emit("build:done", buildResult)
func (em *EventManager) Emit(name string, data ...any) bool {
	event := &CustomEvent{Name: name}
	if len(data) == 1 {
		event.Data = data[0]
	} else if len(data) > 1 {
		event.Data = data
	}
	return em.EmitEvent(event)
}

// EmitEvent emits a pre-constructed CustomEvent.
// Returns true if the event was cancelled by a listener.
//
//	cancelled := em.EmitEvent(&application.CustomEvent{Name: "ping"})
func (em *EventManager) EmitEvent(event *CustomEvent) bool {
	em.mu.RLock()
	listeners := append([]*customEventListener(nil), em.customListeners[event.Name]...)
	em.mu.RUnlock()

	toRemove := []int{}
	for index, listener := range listeners {
		if event.IsCancelled() {
			break
		}
		listener.callback(event)
		if listener.counter > 0 {
			listener.counter--
			if listener.counter == 0 {
				toRemove = append(toRemove, index)
			}
		}
	}

	if len(toRemove) > 0 {
		em.mu.Lock()
		remaining := em.customListeners[event.Name]
		for _, index := range toRemove {
			if index < len(remaining) {
				remaining = append(remaining[:index], remaining[index+1:]...)
			}
		}
		em.customListeners[event.Name] = remaining
		em.mu.Unlock()
	}

	return event.IsCancelled()
}

// On registers a persistent listener for the named custom event.
// Returns a cancel function that deregisters the listener.
//
//	cancel := em.On("build:done", func(e *application.CustomEvent) { ... })
//	defer cancel()
func (em *EventManager) On(name string, callback func(event *CustomEvent)) func() {
	em.mu.Lock()
	em.ensureCustomListeners()
	listener := &customEventListener{callback: callback, counter: -1}
	em.customListeners[name] = append(em.customListeners[name], listener)
	em.mu.Unlock()

	return func() {
		em.mu.Lock()
		defer em.mu.Unlock()
		slice := em.customListeners[name]
		for i, l := range slice {
			if l == listener {
				em.customListeners[name] = append(slice[:i], slice[i+1:]...)
				return
			}
		}
	}
}

// Off removes all listeners for the named custom event.
//
//	em.Off("build:done")
func (em *EventManager) Off(name string) {
	em.mu.Lock()
	delete(em.customListeners, name)
	em.mu.Unlock()
}

// OnMultiple registers a listener that fires at most counter times, then auto-deregisters.
//
//	em.OnMultiple("ping", func(e *application.CustomEvent) { ... }, 3)
func (em *EventManager) OnMultiple(name string, callback func(event *CustomEvent), counter int) {
	em.mu.Lock()
	em.ensureCustomListeners()
	em.customListeners[name] = append(em.customListeners[name], &customEventListener{
		callback: callback,
		counter:  counter,
	})
	em.mu.Unlock()
}

// Reset removes all custom event listeners.
//
//	em.Reset()
func (em *EventManager) Reset() {
	em.mu.Lock()
	em.customListeners = make(map[string][]*customEventListener)
	em.mu.Unlock()
}

// OnApplicationEvent registers a listener for a platform application event.
// Returns a cancel function that deregisters the listener.
//
//	cancel := em.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) { ... })
//	defer cancel()
func (em *EventManager) OnApplicationEvent(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func() {
	eventID := uint(eventType)
	em.mu.Lock()
	em.ensureApplicationListeners()
	listener := &applicationEventListener{callback: callback}
	em.applicationListeners[eventID] = append(em.applicationListeners[eventID], listener)
	em.mu.Unlock()

	return func() {
		em.mu.Lock()
		defer em.mu.Unlock()
		slice := em.applicationListeners[eventID]
		for i, l := range slice {
			if l == listener {
				em.applicationListeners[eventID] = append(slice[:i], slice[i+1:]...)
				return
			}
		}
	}
}
