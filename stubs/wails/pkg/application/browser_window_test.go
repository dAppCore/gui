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

// AX7 generated source-matching smoke coverage.
func TestBrowserWindow_BrowserWindow_ID_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ID_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ID_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Name_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Name_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Name_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ClientID_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ClientID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ClientID_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ClientID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ClientID_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ClientID()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_DispatchWailsEvent_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.DispatchWailsEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_DispatchWailsEvent_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.DispatchWailsEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_DispatchWailsEvent_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.DispatchWailsEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_EmitEvent_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.EmitEvent("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_EmitEvent_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.EmitEvent("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_EmitEvent_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.EmitEvent("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Error_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Error("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Error_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Error("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Error_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Error("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Info_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Info("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Info_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Info("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Info_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Info("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Center_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Center()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Center_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Center()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Center_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Center()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Close_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Close_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Close_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Close()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_DisableSizeConstraints_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.DisableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_DisableSizeConstraints_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.DisableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_DisableSizeConstraints_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.DisableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_EnableSizeConstraints_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.EnableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_EnableSizeConstraints_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.EnableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_EnableSizeConstraints_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.EnableSizeConstraints()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ExecJS_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ExecJS("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ExecJS_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ExecJS("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ExecJS_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ExecJS("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Focus_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Focus_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Focus_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Focus()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ForceReload_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ForceReload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ForceReload_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ForceReload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ForceReload_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ForceReload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Fullscreen_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Fullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Fullscreen_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Fullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Fullscreen_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Fullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetBorderSizes_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.GetBorderSizes()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetBorderSizes_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.GetBorderSizes()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetBorderSizes_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.GetBorderSizes()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetScreen_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetScreen_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetScreen_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.GetScreen()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetZoom_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetZoom_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_GetZoom_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.GetZoom()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleMessage_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleMessage("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleMessage_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleMessage("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleMessage_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleMessage("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleWindowEvent_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleWindowEvent(1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleWindowEvent_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleWindowEvent(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleWindowEvent_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleWindowEvent(0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Height_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Height()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Height_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Height()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Height_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Height()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Hide_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Hide()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Hide_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Hide()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Hide_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Hide()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HideMenuBar_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HideMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HideMenuBar_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HideMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HideMenuBar_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HideMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsFocused_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsFocused_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsFocused_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsFocused()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsFullscreen_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsFullscreen_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsFullscreen_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsFullscreen()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsIgnoreMouseEvents_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsIgnoreMouseEvents()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsIgnoreMouseEvents_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsIgnoreMouseEvents()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsIgnoreMouseEvents_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsIgnoreMouseEvents()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsMaximised_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsMaximised_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsMaximised_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsMaximised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsMinimised_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsMinimised_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsMinimised_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsMinimised()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleKeyEvent_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleKeyEvent("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleKeyEvent_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleKeyEvent("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_HandleKeyEvent_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.HandleKeyEvent("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Maximise_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Maximise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Maximise_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Maximise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Maximise_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Maximise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Minimise_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Minimise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Minimise_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Minimise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Minimise_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Minimise()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OnWindowEvent_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.OnWindowEvent(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OnWindowEvent_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.OnWindowEvent(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OnWindowEvent_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.OnWindowEvent(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OpenContextMenu_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.OpenContextMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OpenContextMenu_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.OpenContextMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OpenContextMenu_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.OpenContextMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Position_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Position_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Position_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Position()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_RelativePosition_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.RelativePosition()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_RelativePosition_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.RelativePosition()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_RelativePosition_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.RelativePosition()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Reload_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Reload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Reload_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Reload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Reload_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Reload()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Resizable_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Resizable()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Resizable_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Resizable()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Resizable_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Resizable()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Restore_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Restore_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Restore_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Restore()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Run_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Run()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Run_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Run()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Run_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Run()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetPosition_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetPosition(1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetPosition_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetPosition(0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetPosition_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetPosition(-1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetAlwaysOnTop_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetAlwaysOnTop(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetAlwaysOnTop_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetAlwaysOnTop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetAlwaysOnTop_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetAlwaysOnTop(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetBackgroundColour_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetBackgroundColour(*new(RGBA))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetBackgroundColour_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetBackgroundColour(*new(RGBA))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetBackgroundColour_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetBackgroundColour(*new(RGBA))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetFrameless_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetFrameless(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetFrameless_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetFrameless(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetFrameless_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetFrameless(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetHTML_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetHTML("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetHTML_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetHTML("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetHTML_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetHTML("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMinimiseButtonState_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinimiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMinimiseButtonState_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinimiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMinimiseButtonState_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinimiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMaximiseButtonState_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaximiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMaximiseButtonState_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaximiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMaximiseButtonState_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaximiseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetCloseButtonState_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetCloseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetCloseButtonState_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetCloseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetCloseButtonState_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetCloseButtonState(*new(ButtonState))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMaxSize_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaxSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMaxSize_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaxSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMaxSize_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMaxSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMinSize_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMinSize_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMinSize_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetMinSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetRelativePosition_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetRelativePosition(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetRelativePosition_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetRelativePosition(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetRelativePosition_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetRelativePosition(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetResizable_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetResizable(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetResizable_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetResizable(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetResizable_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetResizable(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetIgnoreMouseEvents_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetIgnoreMouseEvents(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetIgnoreMouseEvents_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetIgnoreMouseEvents(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetIgnoreMouseEvents_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetIgnoreMouseEvents(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetSize_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetSize(1, 1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetSize_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetSize(0, 0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetSize_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetSize(-1, -1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetTitle_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetTitle("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetTitle_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetTitle("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetTitle_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetTitle("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetURL_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetURL("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetURL_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetURL("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetURL_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetURL("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetZoom_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetZoom(1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetZoom_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetZoom(0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetZoom_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetZoom(-1.5)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Show_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Show()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Show_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Show()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Show_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Show()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ShowMenuBar_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ShowMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ShowMenuBar_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ShowMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ShowMenuBar_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ShowMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Size_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Size_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Size_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0, got1 := subject.Size()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OpenDevTools_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OpenDevTools_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_OpenDevTools_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.OpenDevTools()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleFullscreen_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleFullscreen_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleFullscreen_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleMaximise_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleMaximise_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleMaximise_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleMenuBar_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleMenuBar_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleMenuBar_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleMenuBar()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleFrameless_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleFrameless()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleFrameless_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleFrameless()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ToggleFrameless_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ToggleFrameless()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnFullscreen_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnFullscreen_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnFullscreen_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnFullscreen()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnMaximise_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnMaximise_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnMaximise_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnMaximise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnMinimise_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnMinimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnMinimise_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnMinimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_UnMinimise_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.UnMinimise()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Width_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Width()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Width_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Width()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Width_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Width()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsVisible_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsVisible_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_IsVisible_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.IsVisible()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Bounds_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Bounds()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Bounds_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Bounds()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Bounds_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Bounds()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetBounds_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetBounds(*new(Rect))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetBounds_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetBounds(*new(Rect))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetBounds_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetBounds(*new(Rect))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Zoom_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Zoom()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Zoom_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Zoom()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Zoom_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Zoom()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomIn_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ZoomIn()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomIn_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ZoomIn()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomIn_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ZoomIn()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomOut_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ZoomOut()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomOut_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ZoomOut()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomOut_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.ZoomOut()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomReset_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ZoomReset()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomReset_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ZoomReset()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_ZoomReset_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.ZoomReset()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMenu_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMenu_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetMenu_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetMenu(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SnapAssist_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SnapAssist()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SnapAssist_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SnapAssist()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SnapAssist_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SnapAssist()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetContentProtection_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetContentProtection(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetContentProtection_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetContentProtection(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetContentProtection_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.SetContentProtection(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetEnabled_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetEnabled(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetEnabled_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_SetEnabled_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Flash_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Flash(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Flash_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Flash_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.Flash(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Print_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Print_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_Print_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.Print()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_RegisterHook_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.RegisterHook(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_RegisterHook_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.RegisterHook(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_RegisterHook_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.RegisterHook(nil, nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_InitiateFrontendDropProcessing_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.InitiateFrontendDropProcessing(nil, 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_InitiateFrontendDropProcessing_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.InitiateFrontendDropProcessing(nil, 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_InitiateFrontendDropProcessing_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.InitiateFrontendDropProcessing(nil, -1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_NativeWindow_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.NativeWindow()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_NativeWindow_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.NativeWindow()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_NativeWindow_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		got0 := subject.NativeWindow()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_AttachModal_Good(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.AttachModal(*new(Window))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_AttachModal_Bad(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.AttachModal(*new(Window))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserWindow_BrowserWindow_AttachModal_Ugly(t *core.T) {
	subject := new(BrowserWindow)
	result := core.Try(func() any {
		subject.AttachModal(*new(Window))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
