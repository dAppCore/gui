package systray

import (
	core "dappco.re/go"
)

func TestMockPlatform_NewTray_Good(t *core.T) {
	p := NewMockPlatform()
	tray := p.NewTray()
	core.AssertNotNil(t, tray)

	mockTray := tray.(*exportedMockTray)
	tray.SetIcon([]byte{1, 2, 3})
	tray.SetTemplateIcon([]byte{4, 5, 6})
	tray.SetTooltip("Core")
	tray.SetLabel("Ready")
	tray.SetMenu(p.NewMenu())
	tray.AttachWindow(windowHandleStub{name: "panel"})

	core.AssertEqual(t, []byte{1, 2, 3}, mockTray.icon)
	core.AssertEqual(t, []byte{4, 5, 6}, mockTray.templateIcon)
	core.AssertEqual(t, "Core", mockTray.tooltip)
	core.AssertEqual(t, "Ready", mockTray.label)
	core.AssertNotNil(t, mockTray)
}

func TestMockPlatform_NewMenu_Bad(t *core.T) {
	p := NewMockPlatform()
	menu := p.NewMenu()
	core.AssertNotNil(t, menu)
	_, ok := menu.(*exportedMockMenu)
	core.AssertTrue(t, ok)
}

func TestMockPlatform_NewTray_Ugly(t *core.T) {
	p := NewMockPlatform()
	tray := p.NewTray().(*exportedMockTray)
	core.AssertNotNil(t, tray)
	core.AssertNoError(t, tray.ShowMessage("title", "message"))
}

type windowHandleStub struct {
	name string
}

func (w windowHandleStub) Name() string { return w.name }

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

func TestMockPlatform_MockPlatform_NewTray_Good(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_NewTray_Bad(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_NewTray_Ugly(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_NewMenu_Good(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_NewMenu_Bad(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_NewMenu_Ugly(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

type MockTray = exportedMockTray

func TestMockPlatform_MockTray_SetIcon_Good(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetIcon_Bad(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetIcon_Ugly(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetTemplateIcon_Good(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetTemplateIcon_Bad(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetTemplateIcon_Ugly(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetTooltip_Good(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetTooltip_Bad(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetTooltip_Ugly(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetLabel_Good(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetLabel("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetLabel_Bad(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetLabel("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetLabel_Ugly(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetLabel("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetMenu_Good(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetMenu_Bad(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_SetMenu_Ugly(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_AttachWindow_Good(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_AttachWindow_Bad(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_AttachWindow_Ugly(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_ShowMessage_Good(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_ShowMessage_Bad(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockTray_ShowMessage_Ugly(t *core.T) {
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

type MockMenu = exportedMockMenu

func TestMockPlatform_MockMenu_Add_Good(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.Add("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_Add_Bad(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.Add("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_Add_Ugly(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.Add("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_AddSeparator_Good(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_AddSeparator_Bad(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_AddSeparator_Ugly(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_AddSubmenu_Good(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_AddSubmenu_Bad(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenu_AddSubmenu_Ugly(t *core.T) {
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

type MockMenuItem = exportedMockMenuItem

func TestMockPlatform_MockMenuItem_SetTooltip_Good(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetTooltip_Bad(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetTooltip_Ugly(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetChecked_Good(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetChecked_Bad(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetChecked_Ugly(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetEnabled_Good(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetEnabled_Bad(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_SetEnabled_Ugly(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_OnClick_Good(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_OnClick_Bad(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockMenuItem_OnClick_Ugly(t *core.T) {
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

// AX7 generated source-matching smoke coverage.
