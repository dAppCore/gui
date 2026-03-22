// pkg/dock/service.go
package dock

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the dock service.
type Options struct{}

// Service is a core.Service managing dock/taskbar operations via IPC.
// It embeds ServiceRuntime for Core access and delegates to Platform.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

// --- Query Handlers ---

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryVisible:
		return s.platform.IsVisible(), true, nil
	default:
		return nil, false, nil
	}
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskShowIcon:
		if err := s.platform.ShowIcon(); err != nil {
			return nil, true, err
		}
		_ = s.Core().ACTION(ActionVisibilityChanged{Visible: true})
		return nil, true, nil
	case TaskHideIcon:
		if err := s.platform.HideIcon(); err != nil {
			return nil, true, err
		}
		_ = s.Core().ACTION(ActionVisibilityChanged{Visible: false})
		return nil, true, nil
	case TaskSetBadge:
		if err := s.platform.SetBadge(t.Label); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	case TaskRemoveBadge:
		if err := s.platform.RemoveBadge(); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	default:
		return nil, false, nil
	}
}
