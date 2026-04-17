package display

import (
	"bytes"
	"context"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/notification"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/systray"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDisplayAPIService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	return newTestDisplayService(t)
}

func TestDisplayAPI_screenToDisplay_Good(t *testing.T) {
	got := screenToDisplay(&screen.Screen{
		ID:          "screen-1",
		Name:        "Primary",
		ScaleFactor: 2,
		Bounds:      screen.Rect{X: 10, Y: 20, Width: 1920, Height: 1080},
		IsPrimary:   true,
	})

	require.NotNil(t, got)
	assert.Equal(t, "screen-1", got.ID)
	assert.Equal(t, "Primary", got.Name)
	assert.Equal(t, 10, got.X)
	assert.Equal(t, 20, got.Y)
	assert.Equal(t, 1920, got.Width)
	assert.Equal(t, 1080, got.Height)
	assert.Equal(t, 2.0, got.ScaleFactor)
	assert.True(t, got.IsPrimary)
}

func TestDisplayAPI_screenToDisplay_Bad(t *testing.T) {
	assert.Nil(t, screenToDisplay(nil))
}

func TestDisplayAPI_screenToDisplay_Ugly(t *testing.T) {
	got := screenToDisplay(&screen.Screen{})

	require.NotNil(t, got)
	assert.Zero(t, got.ID)
	assert.Zero(t, got.Name)
	assert.Zero(t, got.Width)
	assert.Zero(t, got.Height)
}

func TestDisplayAPI_toDialogOpenFileOptions_Good(t *testing.T) {
	got := toDialogOpenFileOptions(OpenFileOptions{
		Title:            "Pick",
		DefaultDirectory: "/tmp",
		DefaultFilename:  "report.csv",
		AllowMultiple:    true,
		Filters: []FileFilter{
			{DisplayName: "CSV", Pattern: "*.csv"},
		},
	})

	assert.Equal(t, "Pick", got.Title)
	assert.Equal(t, "/tmp", got.Directory)
	assert.Equal(t, "report.csv", got.Filename)
	assert.True(t, got.AllowMultiple)
	require.Len(t, got.Filters, 1)
	assert.Equal(t, "CSV", got.Filters[0].DisplayName)
	assert.Equal(t, "*.csv", got.Filters[0].Pattern)
}

func TestDisplayAPI_toDialogOpenFileOptions_Bad(t *testing.T) {
	got := toDialogOpenFileOptions(OpenFileOptions{})

	assert.Empty(t, got.Title)
	assert.Empty(t, got.Directory)
	assert.Empty(t, got.Filename)
	assert.False(t, got.AllowMultiple)
	assert.Nil(t, got.Filters)
}

func TestDisplayAPI_toDialogOpenFileOptions_Ugly(t *testing.T) {
	got := toDialogOpenFileOptions(OpenFileOptions{
		Filters: []FileFilter{
			{DisplayName: "All", Pattern: "*.*"},
			{DisplayName: "Media", Pattern: "*.png;*.jpg"},
		},
	})

	require.Len(t, got.Filters, 2)
	assert.Equal(t, "All", got.Filters[0].DisplayName)
	assert.Equal(t, "*.png;*.jpg", got.Filters[1].Pattern)
}

func TestDisplayAPI_trayMenuItemsToSystray_Good(t *testing.T) {
	got := trayMenuItemsToSystray([]TrayMenuItem{
		{Label: "Open", ActionID: "open"},
		{IsSeparator: true},
		{
			Label:    "More",
			ActionID: "more",
			Children: []TrayMenuItem{{Label: "Nested", ActionID: "nested"}},
		},
	})

	require.Len(t, got, 3)
	assert.Equal(t, "Open", got[0].Label)
	assert.Equal(t, "separator", got[1].Type)
	require.Len(t, got[2].Submenu, 1)
	assert.Equal(t, "nested", got[2].Submenu[0].ActionID)
}

func TestDisplayAPI_trayMenuItemsToSystray_Bad(t *testing.T) {
	assert.Nil(t, trayMenuItemsToSystray(nil))
}

func TestDisplayAPI_trayMenuItemsToSystray_Ugly(t *testing.T) {
	got := trayMenuItemsToSystray([]TrayMenuItem{{Children: []TrayMenuItem{{IsSeparator: true}}}})

	require.Len(t, got, 1)
	require.Len(t, got[0].Submenu, 1)
	assert.Equal(t, "separator", got[0].Submenu[0].Type)
}

func TestDisplayAPI_GetScreens_Good(t *testing.T) {
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

	require.Len(t, screens, 1)
	assert.Equal(t, "screen-1", screens[0].ID)
	assert.Equal(t, 10, screens[0].X)
	assert.Equal(t, 1920, screens[0].Width)
}

func TestDisplayAPI_GetScreens_Empty(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []screen.Screen{}, OK: true}
		default:
			return core.Result{}
		}
	})

	assert.Empty(t, svc.GetScreens())
}

func TestDisplayAPI_GetScreens_Bad(t *testing.T) {
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
	require.NotNil(t, screens)
	assert.Empty(t, screens)
}

func TestDisplayAPI_GetScreens_Ugly(t *testing.T) {
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
	require.NotNil(t, screens)
	assert.Empty(t, screens)
}

func TestDisplayAPI_GetWorkAreas_Ugly(t *testing.T) {
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

	require.NotNil(t, areas)
	assert.Empty(t, areas)
}

func TestDisplayAPI_GetScreen_BadType(t *testing.T) {
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

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestDisplayAPI_CreateWindow_UglyResultType(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("window.open", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: true}
	})

	got, err := svc.CreateWindow(CreateWindowOptions{
		Name: "broken-window",
	})

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_GetScreen_Ugly(t *testing.T) {
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

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestDisplayAPI_GetPrimaryScreen_Ugly(t *testing.T) {
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

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestDisplayAPI_GetScreenAtPoint_Ugly(t *testing.T) {
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

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestDisplayAPI_OpenFileDialog_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskOpenFile)
		assert.Equal(t, "Pick file", task.Options.Title)
		assert.True(t, task.Options.AllowMultiple)
		return core.Result{Value: []string{"/tmp/a.txt", "/tmp/b.txt"}, OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{
		Title:         "Pick file",
		AllowMultiple: true,
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"/tmp/a.txt", "/tmp/b.txt"}, paths)
}

func TestDisplayAPI_OpenFileDialog_BadType(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	require.Error(t, err)
	assert.Nil(t, paths)
}

func TestDisplayAPI_OpenFileDialog_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	require.Error(t, err)
	assert.Nil(t, paths)
}

func TestDisplayAPI_OpenFileDialog_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{OK: true}
	})

	paths, err := svc.OpenFileDialog(OpenFileOptions{})

	require.Error(t, err)
	assert.Nil(t, paths)
}

func TestDisplayAPI_RequestNotificationPermission_BadType(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("notification.requestPermission", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: true}
	})

	granted, err := svc.RequestNotificationPermission()

	require.Error(t, err)
	assert.False(t, granted)
}

func TestDisplayAPI_GetTheme_Good(t *testing.T) {
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
	require.NotNil(t, theme)
	assert.True(t, theme.IsDark)
	assert.Equal(t, "dark", svc.GetSystemTheme())
}

func TestDisplayAPI_GetTheme_Bad(t *testing.T) {
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
	assert.Nil(t, theme)
	assert.Empty(t, svc.GetSystemTheme())
}

func TestDisplayAPI_GetTheme_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})

	assert.Nil(t, svc.GetTheme())
	assert.Empty(t, svc.GetSystemTheme())
}

func TestDisplayAPI_SaveFileDialog_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.saveFile", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskSaveFile)
		assert.Equal(t, "Export", task.Options.Title)
		assert.Equal(t, "/tmp", task.Options.Directory)
		assert.Equal(t, "data.json", task.Options.Filename)
		require.Len(t, task.Options.Filters, 1)
		assert.Equal(t, "JSON", task.Options.Filters[0].DisplayName)
		return core.Result{Value: "/exports/data.json", OK: true}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{
		Title:            "Export",
		DefaultDirectory: "/tmp",
		DefaultFilename:  "data.json",
		Filters:          []FileFilter{{DisplayName: "JSON", Pattern: "*.json"}},
	})

	require.NoError(t, err)
	assert.Equal(t, "/exports/data.json", path)
}

func TestDisplayAPI_SaveFileDialog_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.saveFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{})

	require.Error(t, err)
	assert.Empty(t, path)
}

func TestDisplayAPI_SaveFileDialog_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.saveFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	path, err := svc.SaveFileDialog(SaveFileOptions{})

	require.Error(t, err)
	assert.Empty(t, path)
}

func TestDisplayAPI_OpenDirectoryDialog_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openDirectory", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskOpenDirectory)
		assert.Equal(t, "Choose", task.Options.Title)
		assert.Equal(t, "/var", task.Options.Directory)
		return core.Result{Value: "/var/data", OK: true}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{
		Title:            "Choose",
		DefaultDirectory: "/var",
	})

	require.NoError(t, err)
	assert.Equal(t, "/var/data", path)
}

func TestDisplayAPI_OpenDirectoryDialog_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openDirectory", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{})

	require.Error(t, err)
	assert.Empty(t, path)
}

func TestDisplayAPI_OpenDirectoryDialog_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.openDirectory", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	path, err := svc.OpenDirectoryDialog(OpenDirectoryOptions{})

	require.Error(t, err)
	assert.Empty(t, path)
}

func TestDisplayAPI_PromptDialog_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(dialog.TaskPrompt)
		assert.Equal(t, "Rename", task.Title)
		assert.Equal(t, "Enter a new name", task.Message)
		return core.Result{Value: dialog.PromptResult{Value: "draft", Confirmed: true}, OK: true}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	require.NoError(t, err)
	assert.True(t, confirmed)
	assert.Equal(t, "draft", value)
}

func TestDisplayAPI_PromptDialog_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	require.Error(t, err)
	assert.False(t, confirmed)
	assert.Empty(t, value)
}

func TestDisplayAPI_PromptDialog_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("dialog.prompt", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	value, confirmed, err := svc.PromptDialog("Rename", "Enter a new name")

	require.Error(t, err)
	assert.False(t, confirmed)
	assert.Empty(t, value)
}

func TestDisplayAPI_ReadClipboardImage_Good(t *testing.T) {
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

	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3}, got)
	payload[0] = 9
	assert.Equal(t, byte(1), got[0])
}

func TestDisplayAPI_ReadClipboardImage_Bad(t *testing.T) {
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

	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestDisplayAPI_ReadClipboardImage_Ugly(t *testing.T) {
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

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestDisplayAPI_ReadClipboardImage_Ugly_BackendFailure(t *testing.T) {
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

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestDisplayAPI_WriteClipboardImage_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	var got []byte
	c.Action("clipboard.setImage", func(_ context.Context, opts core.Options) core.Result {
		got = append([]byte(nil), opts.Get("data").Value.([]byte)...)
		return core.Result{OK: true}
	})

	input := []byte{4, 5, 6}
	err := svc.WriteClipboardImage(input)

	require.NoError(t, err)
	input[0] = 9
	assert.True(t, bytes.Equal([]byte{4, 5, 6}, got))
}

func TestDisplayAPI_WriteClipboardImage_Bad(t *testing.T) {
	svc, _ := newTestDisplayAPIService(t)

	err := svc.WriteClipboardImage(nil)

	require.Error(t, err)
}

func TestDisplayAPI_WriteClipboardImage_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.Action("clipboard.setImage", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	err := svc.WriteClipboardImage([]byte{1})

	require.Error(t, err)
}

func TestDisplayAPI_GetScreenForWindow_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	var gotX, gotY int
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch typed := q.(type) {
		case window.QueryWindowByName:
			assert.Equal(t, "editor", typed.Name)
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

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "screen-1", got.ID)
	assert.Equal(t, 250, gotX)
	assert.Equal(t, 400, gotY)
	assert.Equal(t, 10, got.X)
	assert.Equal(t, 20, got.Y)
	assert.Equal(t, 1920, got.Width)
	assert.Equal(t, 1080, got.Height)
}

func TestDisplayAPI_GetScreenForWindow_Bad(t *testing.T) {
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

	require.NoError(t, err)
	assert.Nil(t, got)
	assert.False(t, screenQueried)
}

func TestDisplayAPI_GetScreenForWindow_Ugly(t *testing.T) {
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

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_OpenSingleFileDialog_Good(t *testing.T) {
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

	require.NoError(t, err)
	assert.Equal(t, "/tmp/report.csv", path)
	assert.Equal(t, "Pick report", task.Options.Title)
	assert.Equal(t, "report.csv", task.Options.Filename)
}

func TestDisplayAPI_OpenSingleFileDialog_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: []string{}, OK: true}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{})

	require.NoError(t, err)
	assert.Empty(t, path)
}

func TestDisplayAPI_OpenSingleFileDialog_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.openFile", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	path, err := svc.OpenSingleFileDialog(OpenFileOptions{})

	require.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "dialog.openFile action failed")
}

func TestDisplayAPI_ConfirmDialog_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	var task dialog.TaskQuestion
	c.Action("dialog.question", func(_ context.Context, opts core.Options) core.Result {
		task = opts.Get("task").Value.(dialog.TaskQuestion)
		return core.Result{Value: "Yes", OK: true}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	require.NoError(t, err)
	assert.True(t, confirmed)
	assert.Equal(t, "Confirm", task.Title)
	assert.Equal(t, []string{"Yes", "No"}, task.Buttons)
}

func TestDisplayAPI_ConfirmDialog_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.question", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	require.Error(t, err)
	assert.False(t, confirmed)
	assert.Equal(t, assert.AnError, err)
}

func TestDisplayAPI_ConfirmDialog_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("dialog.question", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: 42, OK: true}
	})

	confirmed, err := svc.ConfirmDialog("Confirm", "Delete this file?")

	require.Error(t, err)
	assert.False(t, confirmed)
	assert.Contains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_ReadClipboard_Good(t *testing.T) {
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

	require.NoError(t, err)
	assert.Equal(t, "hello clipboard", text)
}

func TestDisplayAPI_ReadClipboard_Bad(t *testing.T) {
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

	require.NoError(t, err)
	assert.Empty(t, text)
	// Missing seam: QUERY drops non-OK backend errors, so propagation is not observable here.
}

func TestDisplayAPI_ReadClipboard_Ugly(t *testing.T) {
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

	require.Error(t, err)
	assert.Empty(t, text)
	assert.Contains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_CheckNotificationPermission_Good(t *testing.T) {
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

	require.NoError(t, err)
	assert.True(t, granted)
}

func TestDisplayAPI_CheckNotificationPermission_Bad(t *testing.T) {
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

	require.Error(t, err)
	assert.False(t, granted)
	assert.Contains(t, err.Error(), "notification query failed")
}

func TestDisplayAPI_CheckNotificationPermission_Ugly(t *testing.T) {
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

	require.Error(t, err)
	assert.False(t, granted)
	assert.Contains(t, err.Error(), "unexpected result type")
}

func TestDisplayAPI_WriteClipboard_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	var gotText string
	c.Action("clipboard.setText", func(_ context.Context, opts core.Options) core.Result {
		gotText = opts.Get("task").Value.(clipboard.TaskSetText).Text
		return core.Result{OK: true}
	})

	err := svc.WriteClipboard("hello")

	require.NoError(t, err)
	assert.Equal(t, "hello", gotText)
}

func TestDisplayAPI_WriteClipboard_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("clipboard.setText", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	err := svc.WriteClipboard("hello")

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestDisplayAPI_WriteClipboard_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("clipboard.setText", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	err := svc.WriteClipboard("")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "clipboard.setText")
}

func TestDisplayAPI_SetTrayIcon_Good(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	var got []byte
	c.Action("systray.setIcon", func(_ context.Context, opts core.Options) core.Result {
		got = append([]byte(nil), opts.Get("task").Value.(systray.TaskSetTrayIcon).Data...)
		return core.Result{OK: true}
	})

	err := svc.SetTrayIcon([]byte{1, 2, 3})

	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, got)
}

func TestDisplayAPI_SetTrayIcon_Bad(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("systray.setIcon", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: assert.AnError, OK: false}
	})

	err := svc.SetTrayIcon([]byte{1})

	require.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestDisplayAPI_SetTrayIcon_Ugly(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)

	c.Action("systray.setIcon", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "unexpected", OK: false}
	})

	err := svc.SetTrayIcon(nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "systray.setIcon")
}
