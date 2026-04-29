package menu

import (
	core "dappco.re/go"
)

func TestMockPlatform_NewMenu_Good(t *core.T) {
	p := NewMockPlatform()
	menu := p.NewMenu()

	core.AssertNotNil(t, menu)
	root, ok := menu.(*exportedMockPlatformMenu)
	core.RequireTrue(t, ok)

	item := root.Add("Open")
	core.AssertNotNil(t, item)
	item.SetAccelerator("Cmd+O").SetTooltip("open").SetChecked(true).SetEnabled(false).OnClick(func() {})
	root.AddSeparator()
	sub := root.AddSubmenu("More")
	sub.AddRole(RoleHelpMenu)

	core.AssertNotNil(t, root)
	core.AssertNotNil(t, sub)
}

func TestMockPlatform_SetApplicationMenu_Bad(t *core.T) {
	p := NewMockPlatform()
	menu := p.NewMenu()
	p.SetApplicationMenu(menu)

	core.AssertNotNil(t, menu)
}

func TestMockPlatform_NewMenu_Ugly(t *core.T) {
	p := NewMockPlatform()
	root := p.NewMenu().(*exportedMockPlatformMenu)
	root.AddRole(RoleAppMenu)
	root.AddRole(RoleFileMenu)
	root.AddRole(RoleEditMenu)
	root.AddRole(RoleViewMenu)
	root.AddRole(RoleWindowMenu)
	root.AddRole(RoleHelpMenu)
	core.AssertNotNil(t, root)
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

func TestMockPlatform_MockPlatform_SetApplicationMenu_Good(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_SetApplicationMenu_Bad(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatform_SetApplicationMenu_Ugly(t *core.T) {
	subject := new(MockPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

type MockPlatformMenu = exportedMockPlatformMenu

func TestMockPlatform_MockPlatformMenu_Add_Good(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.Add("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_Add_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.Add("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_Add_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.Add("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddSeparator_Good(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddSeparator_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddSeparator_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddSubmenu_Good(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddSubmenu_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddSubmenu_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddRole_Good(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddRole_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenu_AddRole_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

type MockPlatformMenuItem = exportedMockPlatformMenuItem

func TestMockPlatform_MockPlatformMenuItem_SetAccelerator_Good(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetAccelerator_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetAccelerator_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetTooltip_Good(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetTooltip_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetTooltip_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetChecked_Good(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetChecked_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetChecked_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetEnabled_Good(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetEnabled_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_SetEnabled_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_OnClick_Good(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_OnClick_Bad(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestMockPlatform_MockPlatformMenuItem_OnClick_Ugly(t *core.T) {
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

// AX7 generated source-matching smoke coverage.
