// pkg/screen/service.go
package screen

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the screen service.
type Options struct{}

// Service is a core.Service providing screen/display queries via IPC.
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
	return nil
}

// HandleIPCEvents is auto-discovered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryAll:
		return s.platform.GetAll(), true, nil
	case QueryPrimary:
		return s.platform.GetPrimary(), true, nil
	case QueryByID:
		return s.queryByID(q.ID), true, nil
	case QueryAtPoint:
		return s.queryAtPoint(q.X, q.Y), true, nil
	case QueryWorkAreas:
		return s.queryWorkAreas(), true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) queryByID(id string) *Screen {
	for _, scr := range s.platform.GetAll() {
		if scr.ID == id {
			return &scr
		}
	}
	return nil
}

func (s *Service) queryAtPoint(x, y int) *Screen {
	for _, scr := range s.platform.GetAll() {
		b := scr.Bounds
		if x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height {
			return &scr
		}
	}
	return nil
}

func (s *Service) queryWorkAreas() []Rect {
	screens := s.platform.GetAll()
	areas := make([]Rect, len(screens))
	for i, scr := range screens {
		areas[i] = scr.WorkArea
	}
	return areas
}
