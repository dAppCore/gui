// pkg/dialog/service.go
package dialog

import (
	"context"

	core "dappco.re/go/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// Register(p) binds the dialog service to a Core instance.
//
//	c.WithService(dialog.Register(wailsDialog))
func Register(p Platform) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Result{Value: &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().Action("dialog.openFile", func(_ context.Context, opts core.Options) core.Result {
		var openOpts OpenFileOptions
		switch v := opts.Get("task").Value.(type) {
		case TaskOpenFile:
			openOpts = v.Options
		case TaskOpenFileWithOptions:
			if v.Options != nil {
				openOpts = *v.Options
			}
		}
		paths, err := s.platform.OpenFile(openOpts)
		return core.Result{}.New(paths, err)
	})
	s.Core().Action("dialog.saveFile", func(_ context.Context, opts core.Options) core.Result {
		var saveOpts SaveFileOptions
		switch v := opts.Get("task").Value.(type) {
		case TaskSaveFile:
			saveOpts = v.Options
		case TaskSaveFileWithOptions:
			if v.Options != nil {
				saveOpts = *v.Options
			}
		}
		path, err := s.platform.SaveFile(saveOpts)
		return core.Result{}.New(path, err)
	})
	s.Core().Action("dialog.openDirectory", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskOpenDirectory)
		path, err := s.platform.OpenDirectory(t.Options)
		return core.Result{}.New(path, err)
	})
	s.Core().Action("dialog.message", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskMessageDialog)
		button, err := s.platform.MessageDialog(t.Options)
		return core.Result{}.New(button, err)
	})
	s.Core().Action("dialog.info", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskInfo)
		button, err := s.platform.MessageDialog(MessageDialogOptions{
			Type: DialogInfo, Title: t.Title, Message: t.Message, Buttons: t.Buttons,
		})
		return core.Result{}.New(button, err)
	})
	s.Core().Action("dialog.question", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskQuestion)
		button, err := s.platform.MessageDialog(MessageDialogOptions{
			Type: DialogQuestion, Title: t.Title, Message: t.Message, Buttons: t.Buttons,
		})
		return core.Result{}.New(button, err)
	})
	s.Core().Action("dialog.warning", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskWarning)
		button, err := s.platform.MessageDialog(MessageDialogOptions{
			Type: DialogWarning, Title: t.Title, Message: t.Message, Buttons: t.Buttons,
		})
		return core.Result{}.New(button, err)
	})
	s.Core().Action("dialog.error", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskError)
		button, err := s.platform.MessageDialog(MessageDialogOptions{
			Type: DialogError, Title: t.Title, Message: t.Message, Buttons: t.Buttons,
		})
		return core.Result{}.New(button, err)
	})
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}
