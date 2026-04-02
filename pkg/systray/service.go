// pkg/systray/service.go
package systray

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/notification"
)

// Options configures the systray service.
//
// Example:
//
//	core.WithService(systray.Register(platform))
type Options struct{}

// Service manages system tray operations via Core tasks.
// Use: svc := &systray.Service{}
type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	iconPath string
}

// OnStartup loads tray config and registers task handlers.
// Use: _ = svc.OnStartup(context.Background())
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

// HandleIPCEvents satisfies Core's IPC hook.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetTrayIcon:
		return nil, true, s.manager.SetIcon(t.Data)
	case TaskSetTooltip:
		return nil, true, s.manager.SetTooltip(t.Tooltip)
	case TaskSetLabel:
		return nil, true, s.manager.SetLabel(t.Label)
	case TaskSetTrayMenu:
		return nil, true, s.taskSetTrayMenu(t)
	case TaskShowPanel:
		return nil, true, s.manager.ShowPanel()
	case TaskHidePanel:
		return nil, true, s.manager.HidePanel()
	case TaskShowMessage:
		return nil, true, s.showTrayMessage(t.Title, t.Message)
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

func (s *Service) showTrayMessage(title, message string) error {
	if s.manager == nil || !s.manager.IsActive() {
		_, _, err := s.Core().PERFORM(notification.TaskSend{
			Opts: notification.NotificationOptions{Title: title, Message: message},
		})
		return err
	}
	tray := s.manager.Tray()
	if tray == nil {
		return core.E("systray.showTrayMessage", "tray not initialised", nil)
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
// Use: manager := svc.Manager()
func (s *Service) Manager() *Manager {
	return s.manager
}
