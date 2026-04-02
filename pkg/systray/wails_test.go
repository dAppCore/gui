package systray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestWailsTray_AttachWindow_Good(t *testing.T) {
	app := application.NewApp()
	platform := NewWailsPlatform(app)

	tray, ok := platform.NewTray().(*wailsTray)
	require.True(t, ok)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "panel",
		Title:  "Panel",
		Hidden: true,
	})

	tray.AttachWindow(window)

	assert.False(t, window.IsVisible())

	tray.tray.Click()
	assert.True(t, window.IsVisible())
	assert.True(t, window.IsFocused())

	tray.tray.Click()
	assert.False(t, window.IsVisible())
}
