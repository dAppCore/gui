package application

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// EventListener mirrors the public listener handle exposed by Wails.
type EventListener struct {
	callback func(app *ApplicationEvent)
}

// WailsEventListener receives custom Wails events from the dispatcher.
type WailsEventListener interface {
	DispatchWailsEvent(event *CustomEvent)
}

type hook struct {
	callback func(*CustomEvent)
}

type eventListener struct {
	callback func(*CustomEvent)
	counter  int
	delete   bool
}

// EventProcessor dispatches custom events to listeners and hooks.
type EventProcessor struct {
	listeners              map[string][]*eventListener
	notifyLock             sync.RWMutex
	dispatchEventToWindows func(*CustomEvent)
	hooks                  map[string][]*hook
	hookLock               sync.RWMutex
}

// Void is used for registered events that must not carry associated data.
type Void interface {
	sentinel()
}

type applicationHookState struct {
	mu    sync.RWMutex
	hooks map[uint][]*EventListener
}

var (
	applicationEventHooks sync.Map
	registeredEvents      sync.Map
	voidType              = reflect.TypeFor[Void]()
)

// NewWailsEventProcessor creates an in-memory custom event processor.
func NewWailsEventProcessor(dispatchEventToWindows func(*CustomEvent)) *EventProcessor {
	return &EventProcessor{
		listeners:              make(map[string][]*eventListener),
		dispatchEventToWindows: dispatchEventToWindows,
		hooks:                  make(map[string][]*hook),
	}
}

// RegisterApplicationEventHook registers a hook for application events.
func (em *EventManager) RegisterApplicationEventHook(eventType events.ApplicationEventType, callback func(event *ApplicationEvent)) func() {
	if callback == nil {
		return func() {}
	}
	stateAny, _ := applicationEventHooks.LoadOrStore(em, &applicationHookState{
		hooks: make(map[uint][]*EventListener),
	})
	state := stateAny.(*applicationHookState)
	listener := &EventListener{callback: callback}
	eventID := uint(eventType)

	state.mu.Lock()
	state.hooks[eventID] = append(state.hooks[eventID], listener)
	state.mu.Unlock()

	return func() {
		state.mu.Lock()
		defer state.mu.Unlock()
		state.hooks[eventID] = slices.DeleteFunc(state.hooks[eventID], func(existing *EventListener) bool {
			return existing == listener
		})
	}
}

// ToJSON marshals the custom event into JSON.
func (e *CustomEvent) ToJSON() string {
	payload, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(payload)
}

// RegisterEvent records a custom event's expected data type.
func RegisterEvent[Data any](name string) {
	if _, ok := registeredEvents.Load(name); ok {
		panic(fmt.Errorf("event '%s' is already registered", name))
	}
	registeredEvents.Store(name, reflect.TypeFor[Data]())
}

// On registers a persistent event listener.
func (e *EventProcessor) On(eventName string, callback func(event *CustomEvent)) func() {
	return e.registerListener(eventName, callback, -1)
}

// OnMultiple registers an event listener with a limited invocation count.
func (e *EventProcessor) OnMultiple(eventName string, callback func(event *CustomEvent), counter int) func() {
	return e.registerListener(eventName, callback, counter)
}

// Once registers an event listener that will fire only once.
func (e *EventProcessor) Once(eventName string, callback func(event *CustomEvent)) func() {
	return e.registerListener(eventName, callback, 1)
}

// Emit dispatches a custom event to listeners and hooks.
func (e *EventProcessor) Emit(event *CustomEvent) error {
	if event == nil {
		return nil
	}
	if err := validateCustomEvent(event); err != nil {
		event.Cancel()
		return err
	}
	e.hookLock.RLock()
	hooks := append([]*hook(nil), e.hooks[event.Name]...)
	e.hookLock.RUnlock()
	for _, registeredHook := range hooks {
		registeredHook.callback(event)
		if event.IsCancelled() {
			return nil
		}
	}
	e.dispatchEventToListeners(event)
	if e.dispatchEventToWindows != nil && !event.IsCancelled() {
		e.dispatchEventToWindows(event)
	}
	return nil
}

// Off removes all listeners for an event name.
func (e *EventProcessor) Off(eventName string) {
	e.notifyLock.Lock()
	delete(e.listeners, eventName)
	e.notifyLock.Unlock()
}

// OffAll removes all registered listeners.
func (e *EventProcessor) OffAll() {
	e.notifyLock.Lock()
	e.listeners = make(map[string][]*eventListener)
	e.notifyLock.Unlock()
}

// RegisterHook registers a hook that runs before listeners.
func (e *EventProcessor) RegisterHook(eventName string, callback func(*CustomEvent)) func() {
	thisHook := &hook{callback: callback}
	e.hookLock.Lock()
	e.hooks[eventName] = append(e.hooks[eventName], thisHook)
	e.hookLock.Unlock()
	return func() {
		e.hookLock.Lock()
		defer e.hookLock.Unlock()
		e.hooks[eventName] = slices.DeleteFunc(e.hooks[eventName], func(existing *hook) bool {
			return existing == thisHook
		})
	}
}

func (e *EventProcessor) registerListener(eventName string, callback func(*CustomEvent), counter int) func() {
	listener := &eventListener{callback: callback, counter: counter}
	e.notifyLock.Lock()
	e.listeners[eventName] = append(e.listeners[eventName], listener)
	e.notifyLock.Unlock()
	return func() {
		e.notifyLock.Lock()
		defer e.notifyLock.Unlock()
		e.listeners[eventName] = slices.DeleteFunc(e.listeners[eventName], func(existing *eventListener) bool {
			return existing == listener
		})
	}
}

func (e *EventProcessor) dispatchEventToListeners(event *CustomEvent) {
	e.notifyLock.Lock()
	defer e.notifyLock.Unlock()
	listeners := e.listeners[event.Name]
	if len(listeners) == 0 {
		return
	}
	for _, listener := range listeners {
		if event.IsCancelled() {
			return
		}
		listener.callback(event)
		if listener.counter > 0 {
			listener.counter--
			if listener.counter == 0 {
				listener.delete = true
			}
		}
	}
	e.listeners[event.Name] = slices.DeleteFunc(e.listeners[event.Name], func(listener *eventListener) bool {
		return listener.delete
	})
}

func validateCustomEvent(event *CustomEvent) error {
	registered, ok := registeredEvents.Load(event.Name)
	if !ok {
		return nil
	}
	expected := registered.(reflect.Type)
	if expected == voidType {
		if event.Data == nil {
			return nil
		}
		return fmt.Errorf("data of type %T for event '%s' does not match registered data type %s", event.Data, event.Name, expected)
	}
	if event.Data == nil {
		return fmt.Errorf("data of type <nil> for event '%s' does not match registered data type %s", event.Name, expected)
	}
	actual := reflect.TypeOf(event.Data)
	if expected.Kind() == reflect.Interface {
		if actual.Implements(expected) {
			return nil
		}
	} else if actual == expected {
		return nil
	}
	return fmt.Errorf("data of type %s for event '%s' does not match registered data type %s", actual, event.Name, expected)
}
