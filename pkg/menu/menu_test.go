// pkg/menu/menu_test.go
package menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	return NewManager(p), p
}

func TestManager_Build_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "File"},
		{Label: "Edit"},
	}
	menu := m.Build(items)
	assert.NotNil(t, menu)
	assert.Len(t, p.menus, 1)
	assert.Len(t, p.menus[0].items, 2)
	assert.Equal(t, "File", p.menus[0].items[0].label)
}

func TestManager_Build_Separator_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Above"},
		{Type: "separator"},
		{Label: "Below"},
	}
	m.Build(items)
	assert.Len(t, p.menus[0].items, 3)
	assert.Equal(t, "---", p.menus[0].items[1].label)
}

func TestManager_Build_Submenu_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Parent", Children: []MenuItem{
			{Label: "Child 1"},
			{Label: "Child 2"},
		}},
	}
	m.Build(items)
	assert.Len(t, p.menus[0].subs, 1)
	assert.Len(t, p.menus[0].subs[0].items, 2)
}

func TestManager_Build_Accelerator_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Save", Accelerator: "CmdOrCtrl+S"},
	}
	m.Build(items)
	assert.Equal(t, "CmdOrCtrl+S", p.menus[0].items[0].accel)
}

func TestManager_Build_OnClick_Good(t *testing.T) {
	m, p := newTestManager()
	called := false
	items := []MenuItem{
		{Label: "Action", OnClick: func() { called = true }},
	}
	m.Build(items)
	p.menus[0].items[0].onClick()
	assert.True(t, called)
}

func TestManager_Build_Role_Good(t *testing.T) {
	m, p := newTestManager()
	appMenu := RoleAppMenu
	items := []MenuItem{
		{Role: &appMenu},
	}
	m.Build(items)
	assert.Contains(t, p.menus[0].roles, RoleAppMenu)
}

func TestManager_SetApplicationMenu_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{{Label: "Test"}}
	m.SetApplicationMenu(items)
	assert.NotNil(t, p.appMenu)
}

func TestManager_Build_Empty_Good(t *testing.T) {
	m, _ := newTestManager()
	menu := m.Build(nil)
	assert.NotNil(t, menu)
}
