package menu

import (
	core "dappco.re/go"
	"github.com/wailsapp/wails/v3/pkg/application"
	"reflect"
	"unsafe"
)

func TestWailsPlatform_NewMenu_Good(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)

	menu := platform.NewMenu()
	core.AssertNotNil(t, menu)
	root, ok := menu.(*wailsMenu)
	core.RequireTrue(t, ok)

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
	core.RequireTrue(t, menuField.IsValid())
	core.AssertEqual(t, reflect.ValueOf(root.menu).Pointer(), menuField.Pointer())

	onClickField := reflect.ValueOf(item.item).Elem().FieldByName("onClick")
	core.RequireTrue(t, onClickField.IsValid())
	onClick := reflect.NewAt(onClickField.Type(), unsafe.Pointer(onClickField.UnsafeAddr())).Elem().Interface().(func(*application.Context))
	onClick(&application.Context{})
	core.AssertTrue(t, clicked)

	core.AssertEqual(t, "Open", root.menu.Items[0].Label)
	core.AssertEqual(t, "Cmd+O", root.menu.Items[0].Accelerator)
	core.AssertEqual(t, "open", root.menu.Items[0].Tooltip)
	core.AssertTrue(t, root.menu.Items[0].Checked)
	core.AssertFalse(t, root.menu.Items[0].Enabled)
	core.AssertLen(t, sub.menu.Items, 6)
}

func TestWailsPlatform_SetApplicationMenu_Bad(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	platform.SetApplicationMenu(newMockPlatform().NewMenu())

	menuField := reflect.ValueOf(&app.Menu).Elem().FieldByName("applicationMenu")
	core.RequireTrue(t, menuField.IsValid())
	core.AssertTrue(t, menuField.IsNil())
}

func TestWailsPlatform_SetApplicationMenu_NilReceiver_Good(t *core.T) {
	var platform *WailsPlatform
	core.AssertNotPanics(t, func() {
		platform.SetApplicationMenu(newMockPlatform().NewMenu())
	})
}

func TestWailsPlatform_NewMenu_Ugly(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	menu := platform.NewMenu().(*wailsMenu)

	menu.AddRole(MenuRole(99))
	core.AssertNotNil(t, menu)
}
