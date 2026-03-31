package systray

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	iconPath string
}

func (s *Service) OnStartup(ctx context.Context) error {
	configValue, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if trayConfig, ok := configValue.(map[string]any); ok {
			s.applyConfig(trayConfig)
		}
	}
	s.Core().RegisterTask(s.handleTask)
	return nil
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

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetTrayIcon:
		return nil, true, s.manager.SetIcon(t.Data)
	case TaskSetTrayMenu:
		return nil, true, s.taskSetTrayMenu(t)
	case TaskShowPanel:
		// Panel show — deferred (requires WindowHandle integration)
		return nil, true, nil
	case TaskHidePanel:
		// Panel hide — deferred (requires WindowHandle integration)
		return nil, true, nil
	default:
		return nil, false, nil
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
