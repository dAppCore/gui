// pkg/environment/service.go
package environment

import (
	"context"
	"strings"
	"sync"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform    Platform
	cancelTheme func() // returned by Platform.OnThemeChange — called on shutdown
	themeMu     sync.RWMutex
	override    string
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
	s.Core().Action("environment.setTheme", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetTheme)
		isDark, err := s.setThemeOverride(t.Theme)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: ThemeInfo{IsDark: isDark, Theme: themeName(isDark)}, OK: true}
	})
	s.Core().Action("gui.theme.set", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetTheme)
		isDark, err := s.setThemeOverride(t.Theme)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: ThemeInfo{IsDark: isDark, Theme: themeName(isDark)}, OK: true}
	})

	// Register theme change callback — broadcasts ActionThemeChanged via IPC
	s.cancelTheme = s.platform.OnThemeChange(func(isDark bool) {
		if s.hasThemeOverride() {
			return
		}
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
		isDark := s.currentTheme()
		return core.Result{Value: ThemeInfo{IsDark: isDark, Theme: themeName(isDark)}, OK: true}
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

func (s *Service) currentTheme() bool {
	s.themeMu.RLock()
	override := s.override
	s.themeMu.RUnlock()
	switch override {
	case "dark":
		return true
	case "light":
		return false
	default:
		return s.platform.IsDarkMode()
	}
}

func (s *Service) hasThemeOverride() bool {
	s.themeMu.RLock()
	defer s.themeMu.RUnlock()
	return s.override != ""
}

func (s *Service) setThemeOverride(theme string) (bool, error) {
	normalized, err := normalizeTheme(theme)
	if err != nil {
		return false, err
	}
	before := s.currentTheme()

	s.themeMu.Lock()
	s.override = normalized
	s.themeMu.Unlock()

	after := s.currentTheme()
	if before != after {
		_ = s.Core().ACTION(ActionThemeChanged{IsDark: after})
	}
	return after, nil
}

func normalizeTheme(theme string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case "", "system":
		return "", nil
	case "dark":
		return "dark", nil
	case "light":
		return "light", nil
	default:
		return "", coreerr.E("environment.normalizeTheme", "invalid theme: "+theme, nil)
	}
}

func themeName(isDark bool) string {
	if isDark {
		return "dark"
	}
	return "light"
}
