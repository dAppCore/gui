package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockPlatform_CreateWindow_Good(t *testing.T) {
	p := NewMockPlatform()
	w := p.CreateWindow(PlatformWindowOptions{
		Name:   "main",
		Title:  "Core GUI",
		URL:    "/home",
		HTML:   "<main>Ready</main>",
		JS:     "globalThis.ready = true",
		Width:  1280,
		Height: 800,
		X:      10,
		Y:      20,
	})

	require.Len(t, p.Windows, 1)
	got := w.(*MockWindow)
	assert.Equal(t, "main", got.Name())
	assert.Equal(t, []string{"globalThis.ready = true"}, got.ExecJSCalls())
	assert.Equal(t, "Core GUI", got.Title())
	assert.Equal(t, 10, got.x)
	assert.Equal(t, 20, got.y)

	got.SetPosition(30, 40)
	got.SetSize(1920, 1080)
	got.SetVisibility(true)
	got.SetAlwaysOnTop(true)
	got.SetOpacity(0.75)
	got.SetBounds(1, 2, 3, 4)
	got.SetURL("/dashboard")
	got.SetHTML("<main>Updated</main>")
	got.SetZoom(1.25)
	got.SetContentProtection(true)
	got.Maximise()
	got.Restore()
	got.Minimise()
	got.Focus()
	got.Show()
	got.Hide()
	got.Fullscreen()
	got.UnFullscreen()
	got.ToggleFullscreen()
	got.ToggleMaximise()
	got.ExecJS("alert(1)")
	got.Flash(true)
	got.OpenDevTools()
	got.CloseDevTools()

	assert.Equal(t, 1, got.x)
	assert.Equal(t, 2, got.y)
	assert.Equal(t, 3, got.width)
	assert.Equal(t, 4, got.height)
	assert.True(t, got.maximised)
	assert.True(t, got.focused)
	assert.False(t, got.visible)
	assert.True(t, got.fullscreened)
	assert.True(t, got.minimised)
	assert.Equal(t, 0.75, got.opacity)
	assert.Equal(t, []string{"globalThis.ready = true", "alert(1)"}, got.ExecJSCalls())
	assert.True(t, got.flashed)
	assert.False(t, got.DevToolsOpen())
}

func TestMockPlatform_GetWindows_Bad(t *testing.T) {
	p := NewMockPlatform()
	assert.Empty(t, p.GetWindows())
}

func TestMockWindow_FileDrop_Ugly(t *testing.T) {
	w := &mockWindow{}
	calls := 0
	w.OnFileDrop(func(paths []string, targetID string) {
		calls++
		assert.Equal(t, []string{"a.txt"}, paths)
		assert.Equal(t, "drop-zone", targetID)
	})
	w.emitFileDrop([]string{"a.txt"}, "drop-zone")

	assert.Equal(t, 1, calls)
}
