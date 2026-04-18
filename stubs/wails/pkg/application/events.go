package application

import (
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// ApplicationEventContext carries optional context for an application event.
//
//	ctx := event.Context()
//	_ = ctx // reserved for future platform-specific fields
type ApplicationEventContext struct{}

func newApplicationEventContext() *ApplicationEventContext {
	return &ApplicationEventContext{}
}

// ApplicationEvent is emitted by the application layer for system-level events.
//
//	manager.OnApplicationEvent(events.Mac.ApplicationShouldTerminate, func(e *ApplicationEvent) {
//	    e.Cancel()
//	})
type ApplicationEvent struct {
	Id        uint
	ctx       *ApplicationEventContext
	cancelled atomic.Bool
}

// Context returns the context attached to this application event.
//
//	ctx := event.Context()
func (e *ApplicationEvent) Context() *ApplicationEventContext {
	if e == nil {
		return nil
	}
	if e.ctx == nil {
		e.ctx = newApplicationEventContext()
	}
	return e.ctx
}

// Cancel prevents further processing of this application event.
//
//	event.Cancel()
func (e *ApplicationEvent) Cancel() {
	if e == nil {
		return
	}
	e.cancelled.Store(true)
}

// IsCancelled reports whether Cancel has been called on this event.
//
//	if event.IsCancelled() { return }
func (e *ApplicationEvent) IsCancelled() bool {
	if e == nil {
		return false
	}
	return e.cancelled.Load()
}

// CustomEvent is a named application-level event carrying arbitrary data.
//
//	manager.Emit("user:login", userPayload)
type CustomEvent struct {
	Name      string `json:"name"`
	Data      any    `json:"data"`
	Sender    string `json:"sender,omitempty"`
	cancelled atomic.Bool
}

// Cancel prevents further processing of this custom event.
//
//	event.Cancel()
func (e *CustomEvent) Cancel() {
	if e == nil {
		return
	}
	e.cancelled.Store(true)
}

// IsCancelled reports whether Cancel has been called on this event.
//
//	if event.IsCancelled() { return }
func (e *CustomEvent) IsCancelled() bool {
	if e == nil {
		return false
	}
	return e.cancelled.Load()
}

// customEventListener holds a callback and its remaining invocation count.
type customEventListener struct {
	callback func(*CustomEvent)
	counter  int
}

// applicationEventListener holds a callback for an application event.
type applicationEventListener struct {
	callback func(*ApplicationEvent)
}

// EventManager manages custom and application events in-memory.
//
//	manager := &EventManager{}
//	cancel := manager.On("data:ready", func(e *CustomEvent) { process(e.Data) })
//	defer cancel()
type EventManager struct {
	mu              sync.RWMutex
	customListeners map[string][]*customEventListener
	appListeners    map[uint][]*applicationEventListener
}

func newEventManager() *EventManager {
	return &EventManager{
		customListeners: make(map[string][]*customEventListener),
		appListeners:    make(map[uint][]*applicationEventListener),
	}
}

func (em *EventManager) ensureMapsLocked() {
	if em.customListeners == nil {
		em.customListeners = make(map[string][]*customEventListener)
	}
	if em.appListeners == nil {
		em.appListeners = make(map[uint][]*applicationEventListener)
	}
}

// Emit fires a named custom event with optional data to all registered listeners.
// Returns true if the event was cancelled by a listener.
//
//	cancelled := manager.Emit("file:saved", "/home/user/doc.txt")
//	if cancelled { log("event cancelled") }
func (em *EventManager) Emit(name string, data ...any) bool {
	event := &CustomEvent{Name: name}
	switch len(data) {
	case 0:
		// no data
	case 1:
		event.Data = data[0]
	default:
		event.Data = data
	}

	em.mu.Lock()
	em.ensureMapsLocked()
	listeners := append([]*customEventListener(nil), em.customListeners[name]...)
	remaining := em.customListeners[name][:0]
	for _, listener := range em.customListeners[name] {
		if listener.counter < 0 {
			remaining = append(remaining, listener)
		} else {
			listener.counter--
			if listener.counter > 0 {
				remaining = append(remaining, listener)
			}
		}
	}
	em.customListeners[name] = remaining
	em.mu.Unlock()

	for _, listener := range listeners {
		if event.IsCancelled() {
			break
		}
		invokeCustomEventListener(listener, event)
	}
	return event.IsCancelled()
}

func invokeCustomEventListener(listener *customEventListener, event *CustomEvent) {
	if listener == nil || listener.callback == nil || event == nil {
		return
	}

	defer func() {
		_ = recover()
	}()
	listener.callback(event)
}

// On registers a persistent listener for the named custom event.
// Returns a cancellation function that removes the listener.
//
//	cancel := manager.On("theme:changed", func(e *CustomEvent) { applyTheme(e.Data) })
//	defer cancel()
func (em *EventManager) On(name string, callback func(*CustomEvent)) func() {
	listener := &customEventListener{callback: callback, counter: -1}
	em.mu.Lock()
	em.ensureMapsLocked()
	em.customListeners[name] = append(em.customListeners[name], listener)
	em.mu.Unlock()
	return func() {
		em.mu.Lock()
		defer em.mu.Unlock()
		if em.customListeners == nil {
			return
		}
		updated := em.customListeners[name][:0]
		for _, existing := range em.customListeners[name] {
			if existing != listener {
				updated = append(updated, existing)
			}
		}
		em.customListeners[name] = updated
	}
}

// Off removes all listeners for the named custom event.
//
//	manager.Off("theme:changed")
func (em *EventManager) Off(name string) {
	em.mu.Lock()
	delete(em.customListeners, name)
	em.mu.Unlock()
}

// OnMultiple registers a listener for the named custom event that fires at most counter times.
//
//	manager.OnMultiple("startup:phase", onPhase, 3)
func (em *EventManager) OnMultiple(name string, callback func(*CustomEvent), counter int) {
	listener := &customEventListener{callback: callback, counter: counter}
	em.mu.Lock()
	em.ensureMapsLocked()
	em.customListeners[name] = append(em.customListeners[name], listener)
	em.mu.Unlock()
}

// OnApplicationEvent registers a listener for application-level events.
// Returns a cancellation function that removes the listener.
//
//	cancel := manager.OnApplicationEvent(events.Mac.ApplicationShouldTerminate, func(e *ApplicationEvent) {
//	    saveState()
//	})
//	defer cancel()
func (em *EventManager) OnApplicationEvent(eventType events.ApplicationEventType, callback func(*ApplicationEvent)) func() {
	eventID := uint(eventType)
	listener := &applicationEventListener{callback: callback}
	em.mu.Lock()
	em.ensureMapsLocked()
	em.appListeners[eventID] = append(em.appListeners[eventID], listener)
	em.mu.Unlock()
	return func() {
		em.mu.Lock()
		defer em.mu.Unlock()
		if em.appListeners == nil {
			return
		}
		updated := em.appListeners[eventID][:0]
		for _, existing := range em.appListeners[eventID] {
			if existing != listener {
				updated = append(updated, existing)
			}
		}
		em.appListeners[eventID] = updated
	}
}
