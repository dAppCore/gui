// pkg/dialog/service.go
package dialog

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
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
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenFile:
		paths, err := s.platform.OpenFile(t.Options)
		return paths, true, err
	case TaskSaveFile:
		path, err := s.platform.SaveFile(t.Options)
		return path, true, err
	case TaskOpenDirectory:
		path, err := s.platform.OpenDirectory(t.Options)
		return path, true, err
	case TaskMessageDialog:
		button, err := s.platform.MessageDialog(t.Options)
		return button, true, err
	default:
		return nil, false, nil
	}
}
