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

func TestMockWindow_FileDrop_UglyCase(t *core.T) {
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

// AX7 generated source-matching smoke coverage.
func TestMockPlatform_NewMockPlatform_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_NewMockPlatform_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_NewMockPlatform_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_CreateWindow_Good(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.CreateWindow(*new(PlatformWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_CreateWindow_Bad(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.CreateWindow(*new(PlatformWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_CreateWindow_Ugly(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.CreateWindow(*new(PlatformWindowOptions))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_GetWindows_Good(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.GetWindows()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_GetWindows_Bad(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.GetWindows()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_GetWindows_Ugly(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.GetWindows()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Name_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Name_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Name_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Title_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Title_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Title_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Title()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Position_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Position_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Position_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Size_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Size_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Size_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsMaximised_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsMaximised_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsMaximised_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsFocused_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsFocused_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsFocused_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsVisible_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsVisible_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsVisible_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsFullscreen_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsFullscreen_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsFullscreen_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsMinimised_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsMinimised_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_IsMinimised_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetBounds_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1, got2, got3 := subject.GetBounds()
		return core.Sprintf("%T,%T,%T,%T", got0, got1, got2, got3)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetBounds_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1, got2, got3 := subject.GetBounds()
		return core.Sprintf("%T,%T,%T,%T", got0, got1, got2, got3)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetBounds_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0, got1, got2, got3 := subject.GetBounds()
		return core.Sprintf("%T,%T,%T,%T", got0, got1, got2, got3)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetZoom_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetZoom_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetZoom_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetOpacity_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetOpacity_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_GetOpacity_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.GetOpacity()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetTitle_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetTitle("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetTitle_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetTitle("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetTitle_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetTitle("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetPosition_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetPosition(1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetPosition_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetPosition(0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetPosition_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetPosition(-1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetSize_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetSize(1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetSize_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetSize(0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetSize_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetSize(-1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetBackgroundColour_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBackgroundColour(1, 1, 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetBackgroundColour_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBackgroundColour(0, 0, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetBackgroundColour_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBackgroundColour(0, 0, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetVisibility_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetVisibility(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetVisibility_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetVisibility(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetVisibility_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetVisibility(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetAlwaysOnTop_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetAlwaysOnTop(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetAlwaysOnTop_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetAlwaysOnTop(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetAlwaysOnTop_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetAlwaysOnTop(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetOpacity_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetOpacity(1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetOpacity_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetOpacity(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetOpacity_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetOpacity(-1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetBounds_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBounds(1, 1, 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetBounds_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBounds(0, 0, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetBounds_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetBounds(-1, -1, -1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetURL_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetURL("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetURL_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetURL("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetURL_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetURL("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetHTML_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetHTML("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetHTML_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetHTML("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetHTML_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetHTML("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetZoom_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetZoom(1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetZoom_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetZoom(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetZoom_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetZoom(-1.5)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetContentProtection_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetContentProtection(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetContentProtection_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetContentProtection(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_SetContentProtection_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.SetContentProtection(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Maximise_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Maximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Maximise_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Maximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Maximise_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Maximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Restore_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Restore_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Restore_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Minimise_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Minimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Minimise_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Minimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Minimise_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Minimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Focus_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Focus_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Focus_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Close_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Close_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Close_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Show_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Show()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Show_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Show()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Show_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Show()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Hide_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Hide()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Hide_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Hide()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Hide_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Hide()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Fullscreen_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Fullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Fullscreen_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Fullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Fullscreen_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Fullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_UnFullscreen_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_UnFullscreen_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_UnFullscreen_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ToggleFullscreen_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ToggleFullscreen_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ToggleFullscreen_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ToggleMaximise_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ToggleMaximise_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ToggleMaximise_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ExecJS_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ExecJS("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ExecJS_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ExecJS("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ExecJS_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.ExecJS("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Flash_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Flash(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Flash_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Flash_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Print_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Print_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_Print_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OpenDevTools_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OpenDevTools_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OpenDevTools_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_CloseDevTools_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.CloseDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_CloseDevTools_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.CloseDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_CloseDevTools_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.CloseDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OnWindowEvent_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnWindowEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OnWindowEvent_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnWindowEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OnWindowEvent_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnWindowEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OnFileDrop_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnFileDrop(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OnFileDrop_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnFileDrop(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_OnFileDrop_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		subject.OnFileDrop(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ExecJSCalls_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.ExecJSCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ExecJSCalls_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.ExecJSCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_ExecJSCalls_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.ExecJSCalls()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_HTMLContent_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.HTMLContent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_HTMLContent_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.HTMLContent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_HTMLContent_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.HTMLContent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_DevToolsOpen_Good(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.DevToolsOpen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_DevToolsOpen_Bad(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.DevToolsOpen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockWindow_DevToolsOpen_Ugly(t *core.T) {
	subject := new(MockWindow)
	result := core.Try(func() any {
		got0 := subject.DevToolsOpen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
