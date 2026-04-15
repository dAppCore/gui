// pkg/notification/service.go
package notification

import (
	"context"
	"strconv"
	"time"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform   Platform
	categories map[string]NotificationCategory
}

func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			categories:     make(map[string]NotificationCategory),
		}, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSend)
		return core.Result{Value: nil, OK: true}.New(s.send(t.Options))
	})
	s.Core().Action("notification.requestPermission", func(_ context.Context, _ core.Options) core.Result {
		granted, err := s.platform.RequestPermission()
		return core.Result{}.New(granted, err)
	})
	s.Core().Action("notification.revokePermission", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: nil, OK: true}.New(s.platform.RevokePermission())
	})
	s.Core().Action("notification.registerCategory", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskRegisterCategory)
		s.categories[t.Category.ID] = t.Category
		return core.Result{OK: true}
	})
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryPermission:
		granted, err := s.platform.CheckPermission()
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: PermissionStatus{Granted: granted}, OK: true}
	default:
		return core.Result{}
	}
}

// send attempts native notification, falls back to dialog via IPC.
func (s *Service) send(options NotificationOptions) error {
	// Generate ID if not provided
	if options.ID == "" {
		options.ID = "core-" + strconv.FormatInt(time.Now().UnixNano(), 10)
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

	message := options.Message
	if options.Subtitle != "" {
		message = options.Subtitle + "\n\n" + message
	}

	r := s.Core().Action("dialog.message").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskMessageDialog{
			Options: dialog.MessageDialogOptions{
				Type:    dt,
				Title:   options.Title,
				Message: message,
				Buttons: []string{"OK"},
			},
		}},
	))
	if !r.OK {
		if err, ok := r.Value.(error); ok {
			return err
		}
	}
	return nil
}
