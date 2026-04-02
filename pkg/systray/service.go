package systray

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/notification"
)

// Options holds configuration for the systray service.
type Options struct{}

// Service is a core.Service managing the system tray via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	iconPath string
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
	tooltip, _ := cfg["tooltip"].(string)
	if tooltip == "" {
		tooltip = "Core"
	}
	_ = s.manager.Setup(tooltip, tooltip)

	if iconPath, ok := cfg["icon"].(string); ok && iconPath != "" {
		// Icon loading is deferred to when assets are available.
		// Store the path for later use.
		s.iconPath = iconPath
	}
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
	case TaskShowMessage:
		return nil, true, s.taskShowMessage(t.Title, t.Message)
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

func (s *Service) taskShowMessage(title, message string) error {
	if s.manager == nil || !s.manager.IsActive() {
		_, _, err := s.Core().PERFORM(notification.TaskSend{
			Opts: notification.NotificationOptions{Title: title, Message: message},
		})
		return err
	}
	tray := s.manager.Tray()
	if tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	if messenger, ok := tray.(interface{ ShowMessage(title, message string) }); ok {
		messenger.ShowMessage(title, message)
		return nil
	}
	_, _, err := s.Core().PERFORM(notification.TaskSend{
		Opts: notification.NotificationOptions{Title: title, Message: message},
	})
	return err
}

// Manager returns the underlying systray Manager.
func (s *Service) Manager() *Manager {
	return s.manager
}
