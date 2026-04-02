// pkg/notification/service.go
package notification

import (
	"context"
	"fmt"
	"time"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
)

// Options configures the notification service.
// Use: core.WithService(notification.Register(platform))
type Options struct{}

// Service manages notifications via Core tasks and queries.
// Use: svc := &notification.Service{}
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register creates a Core service factory for the notification backend.
// Use: core.New(core.WithService(notification.Register(platform)))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers notification handlers with Core.
// Use: _ = svc.OnStartup(context.Background())
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents satisfies Core's IPC hook.
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
		return nil, true, s.sendNotification(t.Opts)
	case TaskRequestPermission:
		granted, err := s.platform.RequestPermission()
		return granted, true, err
	case TaskClear:
		if clr, ok := s.platform.(clearer); ok {
			return nil, true, clr.Clear()
		}
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

// sendNotification attempts a native notification and falls back to a dialog via IPC.
func (s *Service) sendNotification(opts NotificationOptions) error {
	// Generate an ID when the caller does not provide one.
	if opts.ID == "" {
		opts.ID = fmt.Sprintf("core-%d", time.Now().UnixNano())
	}

	if len(opts.Actions) > 0 {
		if sender, ok := s.platform.(actionSender); ok {
			if err := sender.SendWithActions(opts); err == nil {
				return nil
			}
		}
	}

	if err := s.platform.Send(opts); err != nil {
		// Fall back to a dialog when the native notification fails.
		return s.showFallbackDialog(opts)
	}
	return nil
}

// showFallbackDialog shows a dialog via IPC when native notifications fail.
func (s *Service) showFallbackDialog(opts NotificationOptions) error {
	// Map severity to dialog type.
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
