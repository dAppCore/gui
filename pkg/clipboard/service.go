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
	setText := func(_ context.Context, opts core.Options) core.Result {
		success := s.platform.SetText(opts.String("text"))
		return core.Result{Value: success, OK: true}
	}
	setImage := func(_ context.Context, opts core.Options) core.Result {
		imgPlatform, ok := s.platform.(ImagePlatform)
		if !ok {
			return core.Result{Value: false, OK: true}
		}
		data, _ := opts.Get("data").Value.([]byte)
		success := imgPlatform.SetImage(data)
		return core.Result{Value: success, OK: true}
	}
	clear := func(_ context.Context, _ core.Options) core.Result {
		success := s.platform.SetText("")
		if imgPlatform, ok := s.platform.(ImagePlatform); ok {
			success = imgPlatform.SetImage(nil) && success
		}
		return core.Result{Value: success, OK: true}
	}
	read := func(_ context.Context, _ core.Options) core.Result {
		text, ok := s.platform.Text()
		return core.Result{Value: ClipboardContent{Text: text, HasContent: ok && text != ""}, OK: true}
	}
	s.Core().Action("clipboard.setText", setText)
	s.Core().Action("gui.clipboard.write", setText)
	s.Core().Action("clipboard.setImage", setImage)
	s.Core().Action("clipboard.clear", clear)
	s.Core().Action("gui.clipboard.clear", clear)
	s.Core().Action("gui.clipboard.read", read)
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
	case QueryImage:
		imgPlatform, ok := s.platform.(ImagePlatform)
		if !ok {
			return core.Result{Value: ImageContent{}, OK: true}
		}
		data, hasImage := imgPlatform.Image()
		return core.Result{Value: ImageContent{Data: append([]byte(nil), data...), HasImage: hasImage && len(data) > 0}, OK: true}
	default:
		return core.Result{}
	}
}
