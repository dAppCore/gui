package systray

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the systray service.
type Options struct{}

// Service is a core.Service managing the system tray via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
}

// OnStartup queries config and registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	cfg, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if tCfg, ok := cfg.(map[string]any); ok {
			s.applyConfig(tCfg)
		}
	}
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) applyConfig(cfg map[string]any) {
	// Apply config — tooltip, icon path, etc.
	tooltip, _ := cfg["tooltip"].(string)
	if tooltip == "" {
		tooltip = "Core"
	}
	_ = s.manager.Setup(tooltip, tooltip)
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
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

// Manager returns the underlying systray Manager.
func (s *Service) Manager() *Manager {
	return s.manager
}
