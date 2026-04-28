package window

import (
	core "dappco.re/go"
)

func TestMockPlatform_CreateWindow_Good(t *core.T) {
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

	core.AssertLen(t, p.Windows, 1)
	got := w.(*MockWindow)
	core.AssertEqual(t, "main", got.Name())
	core.AssertEqual(t, []string{"globalThis.ready = true"}, got.ExecJSCalls())
	core.AssertEqual(t, "Core GUI", got.Title())
	core.AssertEqual(t, 10, got.x)
	core.AssertEqual(t, 20, got.y)

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

	core.AssertEqual(t, 1, got.x)
	core.AssertEqual(t, 2, got.y)
	core.AssertEqual(t, 3, got.width)
	core.AssertEqual(t, 4, got.height)
	core.AssertTrue(t, got.maximised)
	core.AssertTrue(t, got.focused)
	core.AssertFalse(t, got.visible)
	core.AssertTrue(t, got.fullscreened)
	core.AssertTrue(t, got.minimised)
	core.AssertEqual(t, 0.75, got.opacity)
	core.AssertEqual(t, []string{"globalThis.ready = true", "alert(1)"}, got.ExecJSCalls())
	core.AssertTrue(t, got.flashed)
	core.AssertFalse(t, got.DevToolsOpen())
}

func TestMockPlatform_GetWindows_Bad(t *core.T) {
	p := NewMockPlatform()
	core.AssertEmpty(t, p.GetWindows())
	core.AssertNotEmpty(t, core.Sprintf("%T", p))
}

func TestMockWindow_FileDrop_Ugly(t *core.T) {
	w := &mockWindow{}
	calls := 0
	w.OnFileDrop(func(paths []string, targetID string) {
		calls++
		core.AssertEqual(t, []string{"a.txt"}, paths)
		core.AssertEqual(t, "drop-zone", targetID)
	})
	w.emitFileDrop([]string{"a.txt"}, "drop-zone")

	core.AssertEqual(t, 1, calls)
}
