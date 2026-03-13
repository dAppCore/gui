// pkg/keybinding/service.go
package keybinding

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the keybinding service.
type Options struct{}

// Service is a core.Service managing keyboard shortcuts via IPC.
// It maintains an in-memory registry of bindings and delegates
// platform-level registration to the Platform interface.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	bindings map[string]BindingInfo
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

// --- Query Handlers ---

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryList:
		return s.queryList(), true, nil
	default:
		return nil, false, nil
	}
}

// queryList reads from the in-memory registry (not platform.GetAll()).
func (s *Service) queryList() []BindingInfo {
	result := make([]BindingInfo, 0, len(s.bindings))
	for _, info := range s.bindings {
		result = append(result, info)
	}
	return result
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskAdd:
		return nil, true, s.taskAdd(t)
	case TaskRemove:
		return nil, true, s.taskRemove(t)
	default:
		return nil, false, nil
	}
}

func (s *Service) taskAdd(t TaskAdd) error {
	if _, exists := s.bindings[t.Accelerator]; exists {
		return ErrAlreadyRegistered
	}

	// Register on platform with a callback that broadcasts ActionTriggered
	err := s.platform.Add(t.Accelerator, func() {
		_ = s.Core().ACTION(ActionTriggered{Accelerator: t.Accelerator})
	})
	if err != nil {
		return fmt.Errorf("keybinding: platform add failed: %w", err)
	}

	s.bindings[t.Accelerator] = BindingInfo{
		Accelerator: t.Accelerator,
		Description: t.Description,
	}
	return nil
}

func (s *Service) taskRemove(t TaskRemove) error {
	if _, exists := s.bindings[t.Accelerator]; !exists {
		return fmt.Errorf("keybinding: not registered: %s", t.Accelerator)
	}

	err := s.platform.Remove(t.Accelerator)
	if err != nil {
		return fmt.Errorf("keybinding: platform remove failed: %w", err)
	}

	delete(s.bindings, t.Accelerator)
	return nil
}
