// pkg/clipboard/service.go
package clipboard

import (
	"context"

	core "dappco.re/go/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register(p) binds the clipboard service to a Core instance.
// c.WithService(clipboard.Register(wailsClipboard))
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
	s.Core().Action("clipboard.setText", func(_ context.Context, opts core.Options) core.Result {
		success := s.platform.SetText(opts.String("text"))
		return core.Result{Value: success, OK: true}
	})
	s.Core().Action("clipboard.clear", func(_ context.Context, _ core.Options) core.Result {
		success := s.platform.SetText("")
		return core.Result{Value: success, OK: true}
	})
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case QueryText:
		text, ok := s.platform.Text()
		return core.Result{Value: ClipboardContent{Text: text, HasContent: ok && text != ""}, OK: true}
	default:
		return core.Result{}
	}
}
