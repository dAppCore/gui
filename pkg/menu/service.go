package menu

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the menu service.
type Options struct{}

// Service is a core.Service managing application menus via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	items    []MenuItem // last-set menu items for QueryGetAppMenu
}

// OnStartup queries config and registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	cfg, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if mCfg, ok := cfg.(map[string]any); ok {
			s.applyConfig(mCfg)
		}
	}
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) applyConfig(cfg map[string]any) {
	// Apply config — e.g., show_dev_tools
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryGetAppMenu:
		return s.items, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetAppMenu:
		s.items = t.Items
		s.manager.SetApplicationMenu(t.Items)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

// Manager returns the underlying menu Manager.
func (s *Service) Manager() *Manager {
	return s.manager
}
