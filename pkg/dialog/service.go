// pkg/dialog/service.go
package dialog

import (
	"context"
	"fmt"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
	"forge.lthn.ai/core/gui/pkg/webview"
	"forge.lthn.ai/core/gui/pkg/window"
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
	s.Core().Action("dialog.prompt", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskPrompt)
		windowName, err := s.promptWindowName()
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		script := promptScript(t.Title, t.Message, t.DefaultValue)
		result := s.Core().Action("webview.evaluate").Run(context.Background(), core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskEvaluate{Window: windowName, Script: script}},
		))
		if !result.OK {
			if e, ok := result.Value.(error); ok {
				return core.Result{Value: e, OK: false}
			}
			return core.Result{OK: false}
		}
		switch value := result.Value.(type) {
		case nil:
			return core.Result{Value: PromptResult{Confirmed: false}, OK: true}
		case string:
			return core.Result{Value: PromptResult{Value: value, Confirmed: true}, OK: true}
		default:
			return core.Result{Value: PromptResult{Value: fmt.Sprint(value), Confirmed: true}, OK: true}
		}
	})
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) promptWindowName() (string, error) {
	r := s.Core().QUERY(window.QueryWindowList{})
	if !r.OK {
		return "", coreerr.E("dialog.promptWindowName", "window service unavailable", nil)
	}
	windows, ok := r.Value.([]window.WindowInfo)
	if !ok {
		return "", coreerr.E("dialog.promptWindowName", "unexpected window list result type", nil)
	}
	for _, info := range windows {
		if info.Focused {
			return info.Name, nil
		}
	}
	if len(windows) > 0 {
		return windows[0].Name, nil
	}
	return "", coreerr.E("dialog.promptWindowName", "no application window available for prompt", nil)
}

func promptScript(title, message, defaultValue string) string {
	promptText := title
	if message != "" {
		if promptText != "" {
			promptText += "\n\n"
		}
		promptText += message
	}
	return core.Sprintf(`(() => {
  const value = window.prompt(%s, %s);
  return value === null ? null : value;
})()`, core.JSONMarshalString(promptText), core.JSONMarshalString(defaultValue))
}
