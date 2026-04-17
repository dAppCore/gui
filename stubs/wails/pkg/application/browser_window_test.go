package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrowserWindow_NewBrowserWindow_Good(t *testing.T) {
	window := NewBrowserWindow(7, "client-abc")

	require.NotNil(t, window)
	assert.Equal(t, uint(7), window.ID())
	assert.Equal(t, "browser-7", window.Name())
	assert.Equal(t, "client-abc", window.ClientID())
	assert.True(t, window.IsVisible())
	assert.False(t, window.IsFullscreen())
	assert.False(t, window.IsMaximised())
	assert.False(t, window.IsMinimised())
	assert.False(t, window.IsFocused())
	assert.False(t, window.IsIgnoreMouseEvents())
	assert.False(t, window.Resizable())
	assert.Equal(t, 1.0, window.GetZoom())
	assert.Nil(t, window.NativeWindow())
	assert.True(t, window.shouldUnconditionallyClose())
}

func TestBrowserWindow_NewBrowserWindow_Bad(t *testing.T) {
	window := NewBrowserWindow(0, "")

	assert.Equal(t, uint(0), window.ID())
	assert.Equal(t, "browser-0", window.Name())
	assert.Empty(t, window.ClientID())
	assert.True(t, window.IsVisible())
}

func TestBrowserWindow_NewBrowserWindow_Ugly(t *testing.T) {
	window := NewBrowserWindow(99, "client")

	assert.Same(t, window, window.Show())
	assert.True(t, window.IsVisible())
	assert.Same(t, window, window.Hide())
	assert.False(t, window.IsVisible())
	assert.Same(t, window, window.Fullscreen())
	assert.Same(t, window, window.Maximise())
	assert.Same(t, window, window.Minimise())
	assert.Same(t, window, window.SetAlwaysOnTop(true))
	assert.Same(t, window, window.SetBackgroundColour(NewRGBA(1, 2, 3, 4)))
	assert.Same(t, window, window.SetFrameless(true))
	assert.Same(t, window, window.SetHTML("<b>hi</b>"))
	assert.Same(t, window, window.SetMinSize(10, 20))
	assert.Same(t, window, window.SetMaxSize(30, 40))
	assert.Same(t, window, window.SetRelativePosition(1, 2))
	assert.Same(t, window, window.SetResizable(true))
	assert.Same(t, window, window.SetIgnoreMouseEvents(true))
	assert.Same(t, window, window.SetSize(100, 200))
	assert.Same(t, window, window.SetTitle("Title"))
	assert.Same(t, window, window.SetURL("https://example.com"))
	assert.Same(t, window, window.SetZoom(1.5))
	assert.Same(t, window, window.ZoomReset())
	assert.NoError(t, window.Print())
}
