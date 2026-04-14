// pkg/notification/service.go
package notification

import (
	"context"
	"strconv"
	"time"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

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
		return nil, true, s.send(t.Options)
	case TaskRequestPermission:
		granted, err := s.platform.RequestPermission()
		return granted, true, err
	case TaskRevokePermission:
		return nil, true, s.platform.RevokePermission()
	case TaskRegisterCategory:
		return nil, true, s.platform.RegisterCategory(t.Category)
	case TaskClear:
		clearPlatform, ok := s.platform.(ClearPlatform)
		if !ok {
			return nil, true, coreerr.E("notification.handleTask", "notification clearing is not supported by this platform", nil)
		}
		return nil, true, clearPlatform.Clear(t.ID)
	default:
		return nil, false, nil
	}
}

// send attempts native notification, falls back to dialog via IPC.
func (s *Service) send(options NotificationOptions) error {
	// Generate ID if not provided
	if options.ID == "" {
		options.ID = "core-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	if len(options.Actions) > 0 {
		categoryID := options.CategoryID
		if categoryID == "" {
			categoryID = options.ID
		}
		if err := s.platform.RegisterCategory(NotificationCategory{
			ID:      categoryID,
			Actions: options.Actions,
		}); err != nil {
			return err
		}
		options.CategoryID = categoryID
	}

	if err := s.platform.Send(options); err != nil {
		// Fallback: show as dialog via IPC
		return s.fallbackDialog(options)
	}
	return nil
}

// fallbackDialog shows a dialog via IPC when native notifications fail.
func (s *Service) fallbackDialog(options NotificationOptions) error {
	// Map severity to dialog type
	var dt dialog.DialogType
	switch options.Severity {
	case SeverityWarning:
		dt = dialog.DialogWarning
	case SeverityError:
		dt = dialog.DialogError
	default:
		dt = dialog.DialogInfo
	}

	msg := options.Message
	if options.Subtitle != "" {
		msg = options.Subtitle + "\n\n" + msg
	}

	_, _, err := s.Core().PERFORM(dialog.TaskMessageDialog{
		Options: dialog.MessageDialogOptions{
			Type:    dt,
			Title:   options.Title,
			Message: msg,
			Buttons: []string{"OK"},
		},
	})
	return err
}
