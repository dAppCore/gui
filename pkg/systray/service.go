package systray

import (
	"context"

	core "dappco.re/go/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	iconPath string
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	r := s.Core().QUERY(QueryConfig{})
	if r.OK {
		if trayConfig, ok := r.Value.(map[string]any); ok {
			s.applyConfig(trayConfig)
		}
	}
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("systray.setIcon", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetTrayIcon)
		return core.Result{Value: nil, OK: true}.New(s.manager.SetIcon(t.Data))
	})
	s.Core().Action("systray.setTooltip", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetTrayTooltip)
		return core.Result{Value: nil, OK: true}.New(s.manager.SetTooltip(t.Tooltip))
	})
	s.Core().Action("systray.setLabel", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetTrayLabel)
		return core.Result{Value: nil, OK: true}.New(s.manager.SetLabel(t.Label))
	})
	s.Core().Action("systray.setMenu", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetTrayMenu)
		return core.Result{Value: nil, OK: true}.New(s.taskSetTrayMenu(t))
	})
	s.Core().Action("systray.showMessage", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskShowMessage)
		return core.Result{Value: nil, OK: true}.New(s.manager.ShowMessage(t.Title, t.Message))
	})
	s.Core().Action("systray.showPanel", func(_ context.Context, _ core.Options) core.Result {
		// Panel show — deferred (requires WindowHandle integration)
		return core.Result{OK: true}
	})
	s.Core().Action("systray.hidePanel", func(_ context.Context, _ core.Options) core.Result {
		// Panel hide — deferred (requires WindowHandle integration)
		return core.Result{OK: true}
	})
	return core.Result{OK: true}
}

func (s *Service) applyConfig(configData map[string]any) {
	tooltip, _ := configData["tooltip"].(string)
	if tooltip == "" {
		tooltip = "Core"
	}
	_ = s.manager.Setup(tooltip, tooltip)

	if iconPath, ok := configData["icon"].(string); ok && iconPath != "" {
		// Icon loading is deferred to when assets are available.
		// Store the path for later use.
		s.iconPath = iconPath
	}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryInfo:
		return core.Result{Value: s.manager.GetInfo(), OK: true}
	default:
		return core.Result{}
	}
}

func (s *Service) taskSetTrayMenu(t TaskSetTrayMenu) error {
	// Register IPC-emitting callbacks for each menu item
	for _, item := range t.Items {
		if item.ActionID != "" {
			actionID := item.ActionID
			s.manager.RegisterCallback(actionID, func() {
				_ = s.Core().ACTION(ActionTrayMenuItemClicked{ActionID: actionID})
			})
		}
	}
	return s.manager.SetMenu(t.Items)
}

func (s *Service) Manager() *Manager {
	return s.manager
}
