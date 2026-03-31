// pkg/environment/service.go
package environment

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform    Platform
	cancelTheme func() // returned by Platform.OnThemeChange — called on shutdown
}

// Register(p) binds the environment service to a Core instance.
// core.WithService(environment.Register(wailsEnvironment))
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

	// Register theme change callback — broadcasts ActionThemeChanged via IPC
	s.cancelTheme = s.platform.OnThemeChange(func(isDark bool) {
		_ = s.Core().ACTION(ActionThemeChanged{IsDark: isDark})
	})
	return nil
}

func (s *Service) OnShutdown(ctx context.Context) error {
	if s.cancelTheme != nil {
		s.cancelTheme()
	}
	return nil
}

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryTheme:
		isDark := s.platform.IsDarkMode()
		theme := "light"
		if isDark {
			theme = "dark"
		}
		return ThemeInfo{IsDark: isDark, Theme: theme}, true, nil
	case QueryInfo:
		return s.platform.Info(), true, nil
	case QueryAccentColour:
		return s.platform.AccentColour(), true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenFileManager:
		return nil, true, s.platform.OpenFileManager(t.Path, t.Select)
	default:
		return nil, false, nil
	}
}
