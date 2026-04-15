package browser

import (
	"context"

	core "dappco.re/go/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	openURL := func(_ context.Context, opts core.Options) core.Result {
		if err := s.platform.OpenURL(opts.String("url")); err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{OK: true}
	}
	openFile := func(_ context.Context, opts core.Options) core.Result {
		if err := s.platform.OpenFile(opts.String("path")); err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{OK: true}
	}
	s.Core().Action("browser.openURL", openURL)
	s.Core().Action("gui.browser.open", openURL)
	s.Core().Action("browser.openFile", openFile)
	s.Core().Action("gui.browser.openFile", openFile)
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}
