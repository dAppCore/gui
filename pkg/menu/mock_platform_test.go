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
