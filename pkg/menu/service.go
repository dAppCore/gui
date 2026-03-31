package menu

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager      *Manager
	platform     Platform
	menuItems    []MenuItem
	showDevTools bool
}

func (s *Service) OnStartup(ctx context.Context) error {
	configValue, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if menuConfig, ok := configValue.(map[string]any); ok {
			s.applyConfig(menuConfig)
		}
	}
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) applyConfig(configData map[string]any) {
	if v, ok := configData["show_dev_tools"]; ok {
		if show, ok := v.(bool); ok {
			s.showDevTools = show
		}
	}
}

func (s *Service) ShowDevTools() bool {
	return s.showDevTools
}

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryGetAppMenu:
		return s.menuItems, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetAppMenu:
		s.menuItems = t.Items
		s.manager.SetApplicationMenu(t.Items)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) Manager() *Manager {
	return s.manager
}
