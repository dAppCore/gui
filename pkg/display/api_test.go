package display

import (
	"bytes"
	"context"

	core "dappco.re/go"
	"dappco.re/go/gui/pkg/clipboard"
	"dappco.re/go/gui/pkg/dialog"
	"dappco.re/go/gui/pkg/environment"
	"dappco.re/go/gui/pkg/notification"
	"dappco.re/go/gui/pkg/screen"
	"dappco.re/go/gui/pkg/systray"
	"dappco.re/go/gui/pkg/window"
)

func newTestDisplayAPIService(t *core.T) (*Service, *core.Core) {
	t.Helper()
	return newTestDisplayService(t)
}

func TestDisplayAPI_screenToDisplay_Good(t *core.T) {
	got := screenToDisplay(&screen.Screen{
		ID:          "screen-1",
		Name:        "Primary",
		ScaleFactor: 2,
		Bounds:      screen.Rect{X: 10, Y: 20, Width: 1920, Height: 1080},
		IsPrimary:   true,
	})

	core.AssertNotNil(t, got)
	core.AssertEqual(t, "screen-1", got.ID)
	core.AssertEqual(t, "Primary", got.Name)
	core.AssertEqual(t, 10, got.X)
	core.AssertEqual(t, 20, got.Y)
	core.AssertEqual(t, 1920, got.Width)
	core.AssertEqual(t, 1080, got.Height)
	core.AssertEqual(t, 2.0, got.ScaleFactor)
	core.AssertTrue(t, got.IsPrimary)
}

func TestDisplayAPI_screenToDisplay_Bad(t *core.T) {
	core.AssertNil(t, screenToDisplay(nil))
	observedType := core.Sprintf("%T", screenToDisplay(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestDisplayAPI_screenToDisplay_Ugly(t *core.T) {
	got := screenToDisplay(&screen.Screen{})

	core.AssertNotNil(t, got)
	core.AssertEmpty(t, got.ID)
	core.AssertEmpty(t, got.Name)
	core.AssertEmpty(t, got.Width)
	core.AssertEmpty(t, got.Height)
}

func TestDisplayAPI_toDialogOpenFileOptions_Good(t *core.T) {
	got := toDialogOpenFileOptions(OpenFileOptions{
		Title:            "Pick",
		DefaultDirectory: "/tmp",
		DefaultFilename:  "report.csv",
		AllowMultiple:    true,
		Filters: []FileFilter{
			{DisplayName: "CSV", Pattern: "*.csv"},
		},
	})

	core.AssertEqual(t, "Pick", got.Title)
	core.AssertEqual(t, "/tmp", got.Directory)
	core.AssertEqual(t, "report.csv", got.Filename)
	core.AssertTrue(t, got.AllowMultiple)
	core.AssertLen(t, got.Filters, 1)
	core.AssertEqual(t, "CSV", got.Filters[0].DisplayName)
	core.AssertEqual(t, "*.csv", got.Filters[0].Pattern)
}

func TestDisplayAPI_toDialogOpenFileOptions_Bad(t *core.T) {
	got := toDialogOpenFileOptions(OpenFileOptions{})

	core.AssertEmpty(t, got.Title)
	core.AssertEmpty(t, got.Directory)
	core.AssertEmpty(t, got.Filename)
	core.AssertFalse(t, got.AllowMultiple)
	core.AssertNil(t, got.Filters)
}

func TestDisplayAPI_toDialogOpenFileOptions_Ugly(t *core.T) {
	got := toDialogOpenFileOptions(OpenFileOptions{
		Filters: []FileFilter{
			{DisplayName: "All", Pattern: "*.*"},
			{DisplayName: "Media", Pattern: "*.png;*.jpg"},
		},
	})

	core.AssertLen(t, got.Filters, 2)
	core.AssertEqual(t, "All", got.Filters[0].DisplayName)
	core.AssertEqual(t, "*.png;*.jpg", got.Filters[1].Pattern)
}

func TestDisplayAPI_trayMenuItemsToSystray_Good(t *core.T) {
	got := trayMenuItemsToSystray([]TrayMenuItem{
		{Label: "Open", ActionID: "open"},
		{IsSeparator: true},
		{
			Label:    "More",
			ActionID: "more",
			Children: []TrayMenuItem{{Label: "Nested", ActionID: "nested"}},
		},
	})

	core.AssertLen(t, got, 3)
	core.AssertEqual(t, "Open", got[0].Label)
	core.AssertEqual(t, "separator", got[1].Type)
	core.AssertLen(t, got[2].Submenu, 1)
	core.AssertEqual(t, "nested", got[2].Submenu[0].ActionID)
}

func TestDisplayAPI_trayMenuItemsToSystray_Bad(t *core.T) {
	core.AssertNil(t, trayMenuItemsToSystray(nil))
	observedType := core.Sprintf("%T", trayMenuItemsToSystray(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestDisplayAPI_trayMenuItemsToSystray_Ugly(t *core.T) {
	got := trayMenuItemsToSystray([]TrayMenuItem{{Children: []TrayMenuItem{{IsSeparator: true}}}})

	core.AssertLen(t, got, 1)
	core.AssertLen(t, got[0].Submenu, 1)
	core.AssertEqual(t, "separator", got[0].Submenu[0].Type)
}

func TestDisplayAPI_GetScreens_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []screen.Screen{
				{
					ID:          "screen-1",
					Name:        "Primary",
					Bounds:      screen.Rect{X: 10, Y: 20, Width: 1920, Height: 1080},
					ScaleFactor: 2,
					IsPrimary:   true,
				},
			}, OK: true}
		default:
			return core.Result{}
		}
	})

	screens := svc.GetScreens()

	core.AssertLen(t, screens, 1)
	core.AssertEqual(t, "screen-1", screens[0].ID)
	core.AssertEqual(t, 10, screens[0].X)
	core.AssertEqual(t, 1920, screens[0].Width)
}

func TestDisplayAPI_GetScreens_Empty(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []screen.Screen{}, OK: true}
		default:
			return core.Result{}
		}
	})

	core.AssertEmpty(t, svc.GetScreens())
}

func TestDisplayAPI_GetScreens_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []string{"unexpected"}, OK: true}
		default:
			return core.Result{}
		}
	})

	screens := svc.GetScreens()
	core.AssertNotNil(t, screens)
	core.AssertEmpty(t, screens)
}

func TestDisplayAPI_GetScreens_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	screens := svc.GetScreens()
	core.AssertNotNil(t, screens)
	core.AssertEmpty(t, screens)
}

func TestDisplayAPI_GetWorkAreas_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryWorkAreas:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	areas := svc.GetWorkAreas()

	core.AssertNotNil(t, areas)
	core.AssertEmpty(t, areas)
}

func TestDisplayAPI_GetScreen_BadType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryByID:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreen("screen-1")

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_CreateWindow_UglyResultType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("window.open", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: true}
	})

	got, err := svc.CreateWindow(CreateWindowOptions{
		Name: "broken-window",
	})

	core.AssertError(t, err)
	core.AssertNil(t, got)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_GetScreen_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryByID:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreen("screen-1")

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_GetPrimaryScreen_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryPrimary:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetPrimaryScreen()

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_GetScreenAtPoint_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAtPoint:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenAtPoint(10, 20)

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_OpenFileDialog_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskOpenFile)
		core.AssertEqual(t, "Pick file", task.Options.Title)
		core.AssertTrue(t, task.Options.AllowMultiple)
		return core.Result{Value: []string{"/tmp/a.txt", "/tmp/b.txt"}, OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{
		Title:         "Pick file",
		AllowMultiple: true,
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, []string{"/tmp/a.txt", "/tmp/b.txt"}, paths)
}

func TestDisplayAPI_OpenFileDialog_BadType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertNil(t, paths)
}

func TestDisplayAPI_OpenFileDialog_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertNil(t, paths)
}

func TestDisplayAPI_OpenFileDialog_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertNil(t, paths)
}

func TestDisplayAPI_RequestNotificationPermission_BadType(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("notification.requestPermission", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: true}
	})

	granted, err := svc.RequestNotificationPermission()

	core.AssertError(t, err)
	core.AssertFalse(t, granted)
}

func TestDisplayAPI_GetTheme_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{Value: environment.ThemeInfo{IsDark: true, Theme: "dark"}, OK: true}
		default:
			return core.Result{}
		}
	})

	theme := svc.GetTheme()
	core.AssertNotNil(t, theme)
	core.AssertTrue(t, theme.IsDark)
	core.AssertEqual(t, "dark", svc.GetSystemTheme())
}

func TestDisplayAPI_GetTheme_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	theme := svc.GetTheme()
	core.AssertNil(t, theme)
	core.AssertEmpty(t, svc.GetSystemTheme())
}

func TestDisplayAPI_GetTheme_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	core.AssertNil(t, svc.GetTheme())
	core.AssertEmpty(t, svc.GetSystemTheme())
}

func TestDisplayAPI_SaveFileDialog_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.saveFile", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskSaveFile)
		core.AssertEqual(t, "Export", task.Options.Title)
		core.AssertEqual(t, "/tmp", task.Options.Directory)
		core.AssertEqual(t, "data.json", task.Options.Filename)
		core.AssertLen(t, task.Options.Filters, 1)
		core.AssertEqual(t, "JSON", task.Options.Filters[0].DisplayName)
		return core.Result{Value: "/exports/data.json", OK: true}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{
		Title:            "Export",
		DefaultDirectory: "/tmp",
		DefaultFilename:  "data.json",
		Filters:          []FileFilter{{DisplayName: "JSON", Pattern: "*.json"}},
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/exports/data.json", path)
}

func TestDisplayAPI_SaveFileDialog_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.saveFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_SaveFileDialog_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.saveFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_OpenDirectoryDialog_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openDirectory", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskOpenDirectory)
		core.AssertEqual(t, "Choose", task.Options.Title)
		core.AssertEqual(t, "/var", task.Options.Directory)
		return core.Result{Value: "/var/data", OK: true}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{
		Title:            "Choose",
		DefaultDirectory: "/var",
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/var/data", path)
}

func TestDisplayAPI_OpenDirectoryDialog_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openDirectory", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_OpenDirectoryDialog_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openDirectory", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_PromptDialog_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskPrompt)
		core.AssertEqual(t, "Rename", task.Title)
		core.AssertEqual(t, "Enter a new name", task.Message)
		return core.Result{Value: dialog.PromptResult{Value: "draft", Confirmed: true}, OK: true}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	core.RequireNoError(t, err)
	core.AssertTrue(t, confirmed)
	core.AssertEqual(t, "draft", value)
}

func TestDisplayAPI_PromptDialog_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertEmpty(t, value)
}

func TestDisplayAPI_PromptDialog_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertEmpty(t, value)
}

func TestDisplayAPI_ReadClipboardImage_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	payload := []byte{1, 2, 3}
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{Value: clipboard.ImageContent{Data: payload, HasImage: true}, OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.RequireNoError(t, err)
	core.AssertEqual(t, []byte{1, 2, 3}, got)
	payload[0] = 9
	core.AssertEqual(t, byte(1), got[0])
}

func TestDisplayAPI_ReadClipboardImage_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{Value: clipboard.ImageContent{HasImage: false}, OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.RequireNoError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_ReadClipboardImage_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_ReadClipboardImage_Ugly_BackendFailure(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryImage:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	got, err := svc.ReadClipboardImage()

	core.AssertError(t, err)
	core.AssertNil(t, got)
}

func TestDisplayAPI_WriteClipboardImage_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	var got []byte
	c.Action("clipboard.setImage", func(_ context.Context, opts core.Options) core.Result {
		got = append([]byte(nil), opts.Get("data").Value.([]byte)...)
		return core.Result{OK: true}
	})

	input := []byte{4, 5, 6}
	err := svc.WriteClipboardImage(input)

	core.RequireNoError(t, err)
	input[0] = 9
	core.AssertTrue(t, bytes.Equal([]byte{4, 5, 6}, got))
}

func TestDisplayAPI_WriteClipboardImage_Bad(t *core.T) {
	svc, _ := newTestDisplayAPIService(t)

	err := svc.WriteClipboardImage(nil)

	core.AssertError(t, err)
}

func TestDisplayAPI_WriteClipboardImage_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("clipboard.setImage", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	err := svc.WriteClipboardImage([]byte{1})

	core.AssertError(t, err)
}

func TestDisplayAPI_GetScreenForWindow_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	var gotX, gotY int
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch typed := q.(type) {
		case window.QueryWindowByName:
			core.AssertEqual(t, "editor", typed.Name)
			return core.Result{
				Value: &window.WindowInfo{
					Name:   typed.Name,
					X:      100,
					Y:      200,
					Width:  300,
					Height: 400,
				},
				OK: true,
			}
		case screen.QueryAtPoint:
			gotX, gotY = typed.X, typed.Y
			return core.Result{
				Value: &screen.Screen{
					ID:          "screen-1",
					Name:        "Primary",
					ScaleFactor: 2,
					Bounds:      screen.Rect{X: 10, Y: 20, Width: 1920, Height: 1080},
					IsPrimary:   true,
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenForWindow("editor")

	core.RequireNoError(t, err)
	core.AssertNotNil(t, got)
	core.AssertEqual(t, "screen-1", got.ID)
	core.AssertEqual(t, 250, gotX)
	core.AssertEqual(t, 400, gotY)
	core.AssertEqual(t, 10, got.X)
	core.AssertEqual(t, 20, got.Y)
	core.AssertEqual(t, 1920, got.Width)
	core.AssertEqual(t, 1080, got.Height)
}

func TestDisplayAPI_GetScreenForWindow_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	var screenQueried bool
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowByName:
			return core.Result{Value: (*window.WindowInfo)(nil), OK: true}
		case screen.QueryAtPoint:
			screenQueried = true
			return core.Result{OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenForWindow("missing")

	core.RequireNoError(t, err)
	core.AssertNil(t, got)
	core.AssertFalse(t, screenQueried)
}

func TestDisplayAPI_GetScreenForWindow_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowByName:
			return core.Result{
				Value: &window.WindowInfo{X: 1, Y: 2, Width: 3, Height: 4},
				OK:    true,
			}
		case screen.QueryAtPoint:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.GetScreenForWindow("editor")

	core.AssertError(t, err)
	core.AssertNil(t, got)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_OpenSingleFileDialog_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	var task dialog.TaskOpenFile
	c.Action("dialog.openFile", func(_ context.Context, opts core.Options) core.Result {
		task = opts.Get("task").Value.(dialog.TaskOpenFile)
		return core.Result{Value: []string{"/tmp/report.csv"}, OK: true}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{
		Title:           "Pick report",
		DefaultFilename: "report.csv",
	})

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/tmp/report.csv", path)
	core.AssertEqual(t, "Pick report", task.Options.Title)
	core.AssertEqual(t, "report.csv", task.Options.Filename)
}

func TestDisplayAPI_OpenSingleFileDialog_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: []string{}, OK: true}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{})

	core.RequireNoError(t, err)
	core.AssertEmpty(t, path)
}

func TestDisplayAPI_OpenSingleFileDialog_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{})

	core.AssertError(t, err)
	core.AssertEmpty(t, path)
	core.AssertContains(t, err.Error(), "dialog.openFile action failed")
}

func TestDisplayAPI_ConfirmDialog_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	var task dialog.TaskQuestion
	c.Action("dialog.question", func(_ context.Context, opts core.Options) core.Result {
		task = opts.Get("task").Value.(dialog.TaskQuestion)
		return core.Result{Value: "Yes", OK: true}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	core.RequireNoError(t, err)
	core.AssertTrue(t, confirmed)
	core.AssertEqual(t, "Confirm", task.Title)
	core.AssertEqual(t, []string{"Yes", "No"}, task.Buttons)
}

func TestDisplayAPI_ConfirmDialog_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.question", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertEqual(t, core.AnError, err)
}

func TestDisplayAPI_ConfirmDialog_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.question", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	core.AssertError(t, err)
	core.AssertFalse(t, confirmed)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_ReadClipboard_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryText:
			return core.Result{
				Value: clipboard.ClipboardContent{
					Text:       "hello clipboard",
					HasContent: true,
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	text, err := svc.ReadClipboard()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "hello clipboard", text)
}

func TestDisplayAPI_ReadClipboard_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryText:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	text, err := svc.ReadClipboard()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, text)
	// Missing seam: QUERY drops non-OK backend errors, so propagation is not observable here.
}

func TestDisplayAPI_ReadClipboard_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case clipboard.QueryText:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	text, err := svc.ReadClipboard()

	core.AssertError(t, err)
	core.AssertEmpty(t, text)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_CheckNotificationPermission_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case notification.QueryPermission:
			return core.Result{Value: notification.PermissionStatus{Granted: true}, OK: true}
		default:
			return core.Result{}
		}
	})

	granted, err := svc.CheckNotificationPermission()

	core.RequireNoError(t, err)
	core.AssertTrue(t, granted)
}

func TestDisplayAPI_CheckNotificationPermission_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case notification.QueryPermission:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	granted, err := svc.CheckNotificationPermission()

	core.AssertError(t, err)
	core.AssertFalse(t, granted)
	core.AssertContains(t, err.Error(), "notification query failed")
}

func TestDisplayAPI_CheckNotificationPermission_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case notification.QueryPermission:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	granted, err := svc.CheckNotificationPermission()

	core.AssertError(t, err)
	core.AssertFalse(t, granted)
	core.AssertContains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_WriteClipboard_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	var gotText string
	c.Action("clipboard.setText", func(_ context.Context, opts core.Options) core.Result {
		gotText = opts.Get("task").Value.(clipboard.TaskSetText).Text
		return core.Result{OK: true}
	})

	err := svc.WriteClipboard("hello")

	core.RequireNoError(t, err)
	core.AssertEqual(t, "hello", gotText)
}

func TestDisplayAPI_WriteClipboard_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("clipboard.setText", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	err := svc.WriteClipboard("hello")

	core.AssertError(t, err)
	core.AssertEqual(t, core.AnError, err)
}

func TestDisplayAPI_WriteClipboard_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("clipboard.setText", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	err := svc.WriteClipboard("")

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "clipboard.setText")
}

func TestDisplayAPI_SetTrayIcon_Good(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	var got []byte
	c.Action("systray.setIcon", func(_ context.Context, opts core.Options) core.Result {
		got = append([]byte(nil), opts.Get("task").Value.(systray.TaskSetTrayIcon).Data...)
		return core.Result{OK: true}
	})

	err := svc.SetTrayIcon([]byte{1, 2, 3})

	core.RequireNoError(t, err)
	core.AssertEqual(t, []byte{1, 2, 3}, got)
}

func TestDisplayAPI_SetTrayIcon_Bad(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("systray.setIcon", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: core.AnError, OK: false}
	})

	err := svc.SetTrayIcon([]byte{1})

	core.AssertError(t, err)
	core.AssertEqual(t, core.AnError, err)
}

func TestDisplayAPI_SetTrayIcon_Ugly(t *core.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("systray.setIcon", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	err := svc.SetTrayIcon(nil)

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "systray.setIcon")
}
