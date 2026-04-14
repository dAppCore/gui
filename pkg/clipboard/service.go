// pkg/clipboard/service.go
package clipboard

import (
	"context"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register(p) binds the clipboard service to a Core instance.
// c.WithService(clipboard.Register(wailsClipboard))
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
	return nil
}

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryText:
		text, ok := s.platform.Text()
		return ClipboardContent{Text: text, HasContent: ok && text != ""}, true, nil
	case QueryImage:
		imgPlatform, ok := s.platform.(ImagePlatform)
		if !ok {
			return ImageContent{}, true, nil
		}
		data, hasImage := imgPlatform.Image()
		return ImageContent{Data: data, HasImage: hasImage && len(data) > 0}, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetText:
		return s.platform.SetText(t.Text), true, nil
	case TaskSetImage:
		imgPlatform, ok := s.platform.(ImagePlatform)
		if !ok {
			return nil, true, coreerr.E("clipboard.handleTask", "clipboard image operations are not supported by this platform", nil)
		}
		return imgPlatform.SetImage(t.Data), true, nil
	case TaskClear:
		success := s.platform.SetText("")
		if imgPlatform, ok := s.platform.(ImagePlatform); ok {
			success = imgPlatform.SetImage(nil) && success
		}
		return success, true, nil
	default:
		return nil, false, nil
	}
}
