// pkg/clipboard/service.go
package clipboard

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options configures the clipboard service.
// Use: core.WithService(clipboard.Register(platform))
type Options struct{}

// Service manages clipboard operations via Core queries and tasks.
// Use: svc := &clipboard.Service{}
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register creates a Core service factory for the clipboard backend.
// Use: core.New(core.WithService(clipboard.Register(platform)))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}

// OnStartup registers clipboard handlers with Core.
// Use: _ = svc.OnStartup(context.Background())
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents satisfies Core's IPC hook.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryText:
		text, ok := s.platform.Text()
		return ClipboardContent{Text: text, HasContent: ok && text != ""}, true, nil
	case QueryImage:
		if reader, ok := s.platform.(imageReader); ok {
			data, _ := reader.Image()
			return encodeImageContent(data), true, nil
		}
		return ClipboardImageContent{MimeType: "image/png"}, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetText:
		return s.platform.SetText(t.Text), true, nil
	case TaskClear:
		_ = s.platform.SetText("")
		if writer, ok := s.platform.(imageWriter); ok {
			// Best-effort clear for image-aware clipboard backends.
			_ = writer.SetImage(nil)
		}
		return true, true, nil
	case TaskSetImage:
		if writer, ok := s.platform.(imageWriter); ok {
			return writer.SetImage(t.Data), true, nil
		}
		return false, true, core.E("clipboard.handleTask", "clipboard image write not supported", nil)
	default:
		return nil, false, nil
	}
}
