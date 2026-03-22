// pkg/clipboard/service.go
package clipboard

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the clipboard service.
type Options struct{}

// Service is a core.Service managing clipboard operations via IPC.
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
	case QueryText:
		text, ok := s.platform.Text()
		return ClipboardContent{Text: text, HasContent: ok && text != ""}, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetText:
		return s.platform.SetText(t.Text), true, nil
	case TaskClear:
		return s.platform.SetText(""), true, nil
	default:
		return nil, false, nil
	}
}
