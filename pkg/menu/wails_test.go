package menu

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestWailsPlatform_NewMenu_Good(t *testing.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)

	menu := platform.NewMenu()
	require.NotNil(t, menu)
	root, ok := menu.(*wailsMenu)
	require.True(t, ok)

	clicked := false
	item := root.Add("Open").(*wailsMenuItem)
	item.SetAccelerator("Cmd+O").SetTooltip("open").SetChecked(true).SetEnabled(false).OnClick(func() {
		clicked = true
	})
	root.AddSeparator()
	sub := root.AddSubmenu("More").(*wailsMenu)
	sub.AddRole(RoleAppMenu)
	sub.AddRole(RoleFileMenu)
	sub.AddRole(RoleEditMenu)
	sub.AddRole(RoleViewMenu)
	sub.AddRole(RoleWindowMenu)
	sub.AddRole(RoleHelpMenu)

	platform.SetApplicationMenu(root)

	menuField := reflect.ValueOf(&app.Menu).Elem().FieldByName("applicationMenu")
	require.True(t, menuField.IsValid())
	assert.Equal(t, reflect.ValueOf(root.menu).Pointer(), menuField.Pointer())

	onClickField := reflect.ValueOf(item.item).Elem().FieldByName("onClick")
	require.True(t, onClickField.IsValid())
	onClick := reflect.NewAt(onClickField.Type(), unsafe.Pointer(onClickField.UnsafeAddr())).Elem().Interface().(func(*application.Context))
	onClick(&application.Context{})
	assert.True(t, clicked)

	assert.Equal(t, "Open", root.menu.Items[0].Label)
	assert.Equal(t, "Cmd+O", root.menu.Items[0].Accelerator)
	assert.Equal(t, "open", root.menu.Items[0].Tooltip)
	assert.True(t, root.menu.Items[0].Checked)
	assert.False(t, root.menu.Items[0].Enabled)
	assert.Len(t, sub.menu.Items, 6)
}

func TestWailsPlatform_SetApplicationMenu_Bad(t *testing.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	platform.SetApplicationMenu(newMockPlatform().NewMenu())

	menuField := reflect.ValueOf(&app.Menu).Elem().FieldByName("applicationMenu")
	require.True(t, menuField.IsValid())
	assert.True(t, menuField.IsNil())
}

func TestWailsPlatform_NewMenu_Ugly(t *testing.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	menu := platform.NewMenu().(*wailsMenu)

	menu.AddRole(MenuRole(99))
	assert.NotNil(t, menu)
}
