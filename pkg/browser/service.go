// pkg/browser/service.go
package browser

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the browser service.
type Options struct{}

// Service is a core.Service that delegates browser/file-open operations
// to the platform. It is stateless — no queries, no actions.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenURL:
		return nil, true, s.platform.OpenURL(t.URL)
	case TaskOpenFile:
		return nil, true, s.platform.OpenFile(t.Path)
	default:
		return nil, false, nil
	}
}
