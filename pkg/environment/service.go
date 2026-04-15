// pkg/environment/service.go
package environment

import (
	"context"

	core "dappco.re/go/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform    Platform
	cancelTheme func() // returned by Platform.OnThemeChange — called on shutdown
}

// Register(p) binds the environment service to a Core instance.
// core.WithService(environment.Register(wailsEnvironment))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("environment.openFileManager", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskOpenFileManager)
		if err := s.platform.OpenFileManager(t.Path, t.Select); err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{OK: true}
	})

	// Register theme change callback — broadcasts ActionThemeChanged via IPC
	s.cancelTheme = s.platform.OnThemeChange(func(isDark bool) {
		_ = s.Core().ACTION(ActionThemeChanged{IsDark: isDark})
	})
	return core.Result{OK: true}
}

func (s *Service) OnShutdown(_ context.Context) core.Result {
	if s.cancelTheme != nil {
		s.cancelTheme()
	}
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryTheme:
		isDark := s.platform.IsDarkMode()
		theme := "light"
		if isDark {
			theme = "dark"
		}
		return core.Result{Value: ThemeInfo{IsDark: isDark, Theme: theme}, OK: true}
	case QueryInfo:
		return core.Result{Value: s.platform.Info(), OK: true}
	case QueryAccentColour:
		return core.Result{Value: s.platform.AccentColour(), OK: true}
	case QueryFocusFollowsMouse:
		return core.Result{Value: s.platform.HasFocusFollowsMouse(), OK: true}
	default:
		return core.Result{}
	}
}
