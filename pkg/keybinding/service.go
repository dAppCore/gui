// pkg/keybinding/service.go
package keybinding

import (
	"context"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform           Platform
	registeredBindings map[string]BindingInfo
}

func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

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

func (s *Service) queryList() []BindingInfo {
	result := make([]BindingInfo, 0, len(s.registeredBindings))
	for _, info := range s.registeredBindings {
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
	if _, exists := s.registeredBindings[t.Accelerator]; exists {
		return ErrorAlreadyRegistered
	}

	// Register on platform with a callback that broadcasts ActionTriggered
	err := s.platform.Add(t.Accelerator, func() {
		_ = s.Core().ACTION(ActionTriggered{Accelerator: t.Accelerator})
	})
	if err != nil {
		return coreerr.E("keybinding.taskAdd", "platform add failed", err)
	}

	s.registeredBindings[t.Accelerator] = BindingInfo{
		Accelerator: t.Accelerator,
		Description: t.Description,
	}
	return nil
}

func (s *Service) taskRemove(t TaskRemove) error {
	if _, exists := s.registeredBindings[t.Accelerator]; !exists {
		return coreerr.E("keybinding.taskRemove", "not registered: "+t.Accelerator, nil)
	}

	err := s.platform.Remove(t.Accelerator)
	if err != nil {
		return coreerr.E("keybinding.taskRemove", "platform remove failed", err)
	}

	delete(s.registeredBindings, t.Accelerator)
	return nil
}
