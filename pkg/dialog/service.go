// pkg/dialog/service.go
package dialog

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the dialog service.
type Options struct{}

// Service is a core.Service managing native dialogs via IPC.
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
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenFile:
		paths, err := s.platform.OpenFile(t.Opts)
		return paths, true, err
	case TaskSaveFile:
		path, err := s.platform.SaveFile(t.Opts)
		return path, true, err
	case TaskOpenDirectory:
		path, err := s.platform.OpenDirectory(t.Opts)
		return path, true, err
	case TaskMessageDialog:
		button, err := s.platform.MessageDialog(t.Opts)
		return button, true, err
	default:
		return nil, false, nil
	}
}
