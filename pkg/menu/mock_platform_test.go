package menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockPlatform_NewMenu_Good(t *testing.T) {
	p := NewMockPlatform()
	menu := p.NewMenu()

	require.NotNil(t, menu)
	root, ok := menu.(*exportedMockPlatformMenu)
	require.True(t, ok)

	item := root.Add("Open")
	require.NotNil(t, item)
	item.SetAccelerator("Cmd+O").SetTooltip("open").SetChecked(true).SetEnabled(false).OnClick(func() {})
	root.AddSeparator()
	sub := root.AddSubmenu("More")
	sub.AddRole(RoleHelpMenu)

	assert.NotNil(t, root)
	assert.NotNil(t, sub)
}

func TestMockPlatform_SetApplicationMenu_Bad(t *testing.T) {
	p := NewMockPlatform()
	menu := p.NewMenu()
	p.SetApplicationMenu(menu)

	assert.NotNil(t, menu)
}

func TestMockPlatform_NewMenu_Ugly(t *testing.T) {
	p := NewMockPlatform()
	root := p.NewMenu().(*exportedMockPlatformMenu)
	root.AddRole(RoleAppMenu)
	root.AddRole(RoleFileMenu)
	root.AddRole(RoleEditMenu)
	root.AddRole(RoleViewMenu)
	root.AddRole(RoleWindowMenu)
	root.AddRole(RoleHelpMenu)
	assert.NotNil(t, root)
}
