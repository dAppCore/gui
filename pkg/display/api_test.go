package display

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/screen"
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

	assert.Nil(t, svc.GetScreens())
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

	assert.Nil(t, svc.GetScreens())
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
	require.NotNil(t, theme)
	assert.False(t, theme.IsDark)
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
