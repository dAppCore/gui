package dock

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
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
	case QueryVisible:
		return s.platform.IsVisible(), true, nil
	default:
		return nil, false, nil
	}
}

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
	case TaskSetProgressBar:
		if err := s.platform.SetProgressBar(t.Value); err != nil {
			return nil, true, err
		}
		_ = s.Core().ACTION(ActionProgressChanged{Value: t.Value})
		return nil, true, nil
	case TaskBounce:
		bounceID, err := s.platform.Bounce(t.Type)
		if err != nil {
			return nil, true, err
		}
		_ = s.Core().ACTION(ActionBounceStarted{BounceID: bounceID, Type: t.Type})
		return bounceID, true, nil
	case TaskStopBounce:
		if err := s.platform.StopBounce(t.BounceID); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	default:
		return nil, false, nil
	}
}
