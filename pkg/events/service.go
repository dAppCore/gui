// pkg/events/service.go
package events

import (
	"context"
	"sync"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the events service (currently none).
type Options struct{}

// Service bridges the platform custom event system to the Core IPC bus.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	mu       sync.RWMutex
	counts   map[string]int // listener count per event name
	cancels  []func()
}

func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) OnShutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cancel := range s.cancels {
		cancel()
	}
	s.cancels = nil
	s.counts = make(map[string]int)
	s.platform.Reset()
	return nil
}

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryListeners:
		s.mu.RLock()
		count := s.counts[q.Name]
		s.mu.RUnlock()
		return count, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskEmit:
		return s.taskEmit(t)
	case TaskOn:
		return s.taskOn(t)
	case TaskOff:
		return nil, true, s.taskOff(t)
	default:
		return nil, false, nil
	}
}

func (s *Service) taskEmit(t TaskEmit) (any, bool, error) {
	if t.Name == "" {
		return nil, true, coreerr.E("events.taskEmit", "event name is required", nil)
	}
	cancelled := s.platform.Emit(t.Name, t.Data)
	_ = s.Core().ACTION(ActionEventFired{Name: t.Name, Data: t.Data})
	return cancelled, true, nil
}

func (s *Service) taskOn(t TaskOn) (any, bool, error) {
	if t.Name == "" {
		return nil, true, coreerr.E("events.taskOn", "event name is required", nil)
	}
	cancel := s.platform.On(t.Name, func(event *CustomEvent) {
		_ = s.Core().ACTION(ActionEventFired{Name: event.Name, Data: event.Data})
	})
	s.mu.Lock()
	s.cancels = append(s.cancels, cancel)
	s.counts[t.Name]++
	s.mu.Unlock()
	return nil, true, nil
}

func (s *Service) taskOff(t TaskOff) error {
	if t.Name == "" {
		return coreerr.E("events.taskOff", "event name is required", nil)
	}
	s.platform.Off(t.Name)
	s.mu.Lock()
	delete(s.counts, t.Name)
	s.mu.Unlock()
	return nil
}
