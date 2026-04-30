package application

import (
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type applicationEventHook struct {
	callback func(*ApplicationEvent)
}

type eventHookRegistry struct {
	mu    sync.RWMutex
	hooks map[uint][]*applicationEventHook
}

var eventHookRegistries sync.Map

func (em *EventManager) Once(name string, callback func(*CustomEvent)) func() {
	listener := &customEventListener{callback: callback, counter: 1}
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

func (em *EventManager) EmitEvent(event *CustomEvent) bool {
	if event == nil {
		return false
	}

	em.mu.Lock()
	em.ensureMapsLocked()
	listeners := append([]*customEventListener(nil), em.customListeners[event.Name]...)
	remaining := em.customListeners[event.Name][:0]
	for _, listener := range em.customListeners[event.Name] {
		if listener.counter < 0 {
			remaining = append(remaining, listener)
			continue
		}
		listener.counter--
		if listener.counter > 0 {
			remaining = append(remaining, listener)
		}
	}
	em.customListeners[event.Name] = remaining
	em.mu.Unlock()

	for _, listener := range listeners {
		if event.IsCancelled() {
			break
		}
		invokeCustomEventListener(listener, event)
	}

	return event.IsCancelled()
}

func (em *EventManager) Reset() {
	em.mu.Lock()
	if em.customListeners != nil {
		clear(em.customListeners)
	}
	em.mu.Unlock()
}

func (em *EventManager) RegisterApplicationEventHook(eventType events.ApplicationEventType, callback func(*ApplicationEvent)) func() {
	registry := getEventHookRegistry(em)
	hook := &applicationEventHook{callback: callback}
	eventID := uint(eventType)

	registry.mu.Lock()
	registry.hooks[eventID] = append(registry.hooks[eventID], hook)
	registry.mu.Unlock()

	return func() {
		registry.mu.Lock()
		defer registry.mu.Unlock()
		updated := registry.hooks[eventID][:0]
		for _, existing := range registry.hooks[eventID] {
			if existing != hook {
				updated = append(updated, existing)
			}
		}
		registry.hooks[eventID] = updated
	}
}

func (em *EventManager) handleApplicationEvent(event *ApplicationEvent) {
	if em == nil || event == nil {
		return
	}

	registry := getEventHookRegistry(em)
	registry.mu.RLock()
	hooks := append([]*applicationEventHook(nil), registry.hooks[event.Id]...)
	registry.mu.RUnlock()

	for _, hook := range hooks {
		if event.IsCancelled() {
			return
		}
		if hook == nil || hook.callback == nil {
			continue
		}
		func() {
			defer func() {
				recover()
			}()
			hook.callback(event)
		}()
	}

	em.mu.RLock()
	listeners := append([]*applicationEventListener(nil), em.appListeners[event.Id]...)
	em.mu.RUnlock()
	for _, listener := range listeners {
		if event.IsCancelled() {
			return
		}
		if listener == nil || listener.callback == nil {
			continue
		}
		func() {
			defer func() {
				recover()
			}()
			listener.callback(event)
		}()
	}
}

func getEventHookRegistry(em *EventManager) *eventHookRegistry {
	if em == nil {
		return &eventHookRegistry{hooks: make(map[uint][]*applicationEventHook)}
	}

	registry, ok := eventHookRegistries.Load(em)
	if ok {
		return registry.(*eventHookRegistry)
	}

	newRegistry := &eventHookRegistry{hooks: make(map[uint][]*applicationEventHook)}
	actual, _ := eventHookRegistries.LoadOrStore(em, newRegistry)
	return actual.(*eventHookRegistry)
}
