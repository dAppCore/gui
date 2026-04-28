package application

import (
	core "dappco.re/go"
)

func TestBrowserWindow_NewBrowserWindow_Good(t *core.T) {
	window := NewBrowserWindow(7, "client-abc")

	core.AssertNotNil(t, window)
	core.AssertEqual(t, uint(7), window.ID())
	core.AssertEqual(t, "browser-7", window.Name())
	core.AssertEqual(t, "client-abc", window.ClientID())
	core.AssertTrue(t, window.IsVisible())
	core.AssertFalse(t, window.IsFullscreen())
	core.AssertFalse(t, window.IsMaximised())
	core.AssertFalse(t, window.IsMinimised())
	core.AssertFalse(t, window.IsFocused())
	core.AssertFalse(t, window.IsIgnoreMouseEvents())
	core.AssertFalse(t, window.Resizable())
	core.AssertEqual(t, 1.0, window.GetZoom())
	screen, err := window.GetScreen()
	core.RequireNoError(t, err)
	core.AssertNotNil(t, screen)
	core.AssertEqual(t, Screen{}, *screen)
	core.AssertEqual(t, &LRTB{}, window.GetBorderSizes())
	core.AssertNil(t, window.NativeWindow())
	core.AssertTrue(t, window.shouldUnconditionallyClose())
}

func TestBrowserWindow_NewBrowserWindow_Bad(t *core.T) {
	window := NewBrowserWindow(0, "")

	core.AssertEqual(t, uint(0), window.ID())
	core.AssertEqual(t, "browser-0", window.Name())
	core.AssertEmpty(t, window.ClientID())
	core.AssertTrue(t, window.IsVisible())
}

func TestBrowserWindow_NewBrowserWindow_Ugly(t *core.T) {
	window := NewBrowserWindow(99, "client")

	core.AssertSame(t, window, window.Show())
	core.AssertTrue(t, window.IsVisible())
	core.AssertSame(t, window, window.Hide())
	core.AssertFalse(t, window.IsVisible())
	core.AssertSame(t, window, window.Fullscreen())
	core.AssertSame(t, window, window.Maximise())
	core.AssertSame(t, window, window.Minimise())
	core.AssertSame(t, window, window.SetAlwaysOnTop(true))
	core.AssertSame(t, window, window.SetBackgroundColour(NewRGBA(1, 2, 3, 4)))
	core.AssertSame(t, window, window.SetFrameless(true))
	core.AssertSame(t, window, window.SetHTML("<b>hi</b>"))
	core.AssertSame(t, window, window.SetMinSize(10, 20))
	core.AssertSame(t, window, window.SetMaxSize(30, 40))
	core.AssertSame(t, window, window.SetRelativePosition(1, 2))
	core.AssertSame(t, window, window.SetResizable(true))
	core.AssertSame(t, window, window.SetIgnoreMouseEvents(true))
	core.AssertSame(t, window, window.SetSize(100, 200))
	core.AssertSame(t, window, window.SetTitle("Title"))
	core.AssertSame(t, window, window.SetURL("https://example.com"))
	core.AssertSame(t, window, window.SetZoom(1.5))
	core.AssertSame(t, window, window.ZoomReset())
	core.AssertNoError(t, window.Print())
}
