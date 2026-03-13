// pkg/notification/service.go
package notification

import (
	"context"
	"fmt"
	"time"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
)

// Options holds configuration for the notification service.
type Options struct{}

// Service is a core.Service managing notifications via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryPermission:
		granted, err := s.platform.CheckPermission()
		return PermissionStatus{Granted: granted}, true, err
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSend:
		return nil, true, s.send(t.Opts)
	case TaskRequestPermission:
		granted, err := s.platform.RequestPermission()
		return granted, true, err
	default:
		return nil, false, nil
	}
}

// send attempts native notification, falls back to dialog via IPC.
func (s *Service) send(opts NotificationOptions) error {
	// Generate ID if not provided
	if opts.ID == "" {
		opts.ID = fmt.Sprintf("core-%d", time.Now().UnixNano())
	}

	if err := s.platform.Send(opts); err != nil {
		// Fallback: show as dialog via IPC
		return s.fallbackDialog(opts)
	}
	return nil
}

// fallbackDialog shows a dialog via IPC when native notifications fail.
func (s *Service) fallbackDialog(opts NotificationOptions) error {
	// Map severity to dialog type
	var dt dialog.DialogType
	switch opts.Severity {
	case SeverityWarning:
		dt = dialog.DialogWarning
	case SeverityError:
		dt = dialog.DialogError
	default:
		dt = dialog.DialogInfo
	}

	msg := opts.Message
	if opts.Subtitle != "" {
		msg = opts.Subtitle + "\n\n" + msg
	}

	_, _, err := s.Core().PERFORM(dialog.TaskMessageDialog{
		Opts: dialog.MessageDialogOptions{
			Type:    dt,
			Title:   opts.Title,
			Message: msg,
			Buttons: []string{"OK"},
		},
	})
	return err
}
