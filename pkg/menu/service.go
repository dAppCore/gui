package menu

import (
	"context"

	core "dappco.re/go/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager      *Manager
	platform     Platform
	menuItems    []MenuItem
	showDevTools bool
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	r := s.Core().QUERY(QueryConfig{})
	if r.OK {
		if menuConfig, ok := r.Value.(map[string]any); ok {
			s.applyConfig(menuConfig)
		}
	}
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("menu.setAppMenu", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetAppMenu)
		s.menuItems = t.Items
		s.manager.SetApplicationMenu(t.Items)
		return core.Result{OK: true}
	})
	return core.Result{OK: true}
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

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryGetAppMenu:
		return core.Result{Value: s.menuItems, OK: true}
	default:
		return core.Result{}
	}
}

func (s *Service) Manager() *Manager {
	return s.manager
}
