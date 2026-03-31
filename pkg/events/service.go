// pkg/events/service.go
package events

import (
	"context"
	"sync"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the events service (currently empty).
type Options struct{}

// Service bridges Wails custom events into Core IPC.
// Emit/On/Off/OnMultiple/Reset are available as Tasks; QueryListeners reads state.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform

	mu        sync.Mutex
	listeners map[string][]func()      // IPC-registered cancels per event name
	counts    map[string]int           // listener counts per event name
}

// OnStartup registers query and task handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// OnShutdown cancels all IPC-registered platform listeners.
func (s *Service) OnShutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cancels := range s.listeners {
		for _, cancel := range cancels {
			cancel()
		}
	}
	s.listeners = make(map[string][]func())
	s.counts = make(map[string]int)
	return nil
}

// HandleIPCEvents satisfies the core.Service interface (no-op for now).
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryListeners:
		return s.listenerSnapshot(), true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskEmit:
		cancelled := s.platform.Emit(t.Name, t.Data)
		return cancelled, true, nil

	case TaskOn:
		if t.Name == "" {
			return nil, true, coreerr.E("events.taskOn", "event name must not be empty", nil)
		}
		cancel := s.platform.On(t.Name, func(event *CustomEvent) {
			_ = c.ACTION(ActionEventFired{Event: *event})
		})
		s.mu.Lock()
		s.listeners[t.Name] = append(s.listeners[t.Name], cancel)
		s.counts[t.Name]++
		s.mu.Unlock()
		return nil, true, nil

	case TaskOff:
		s.platform.Off(t.Name)
		s.mu.Lock()
		for _, cancel := range s.listeners[t.Name] {
			cancel()
		}
		delete(s.listeners, t.Name)
		delete(s.counts, t.Name)
		s.mu.Unlock()
		return nil, true, nil

	default:
		return nil, false, nil
	}
}

// listenerSnapshot returns a sorted slice of ListenerInfo for all known event names.
//
//	snapshot := s.listenerSnapshot()
//	for _, info := range snapshot { log(info.EventName, info.Count) }
func (s *Service) listenerSnapshot() []ListenerInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	snapshot := make([]ListenerInfo, 0, len(s.counts))
	for name, count := range s.counts {
		snapshot = append(snapshot, ListenerInfo{EventName: name, Count: count})
	}
	return snapshot
}
