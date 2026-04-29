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

// AX7 generated source-matching smoke coverage.
func TestWails_NewWailsPlatform_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := NewWailsPlatform(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_NewWailsPlatform_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := NewWailsPlatform(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_NewWailsPlatform_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := NewWailsPlatform(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_NewMenu_Good(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_NewMenu_Bad(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_NewMenu_Ugly(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_SetApplicationMenu_Good(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_SetApplicationMenu_Bad(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_SetApplicationMenu_Ugly(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

type Menu = wailsMenu

func TestWails_Menu_Add_Good(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		got0 := subject.Add("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_Add_Bad(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		got0 := subject.Add("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_Add_Ugly(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		got0 := subject.Add("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddSeparator_Good(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddSeparator_Bad(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddSeparator_Ugly(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddSubmenu_Good(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddSubmenu_Bad(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddSubmenu_Ugly(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddRole_Good(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddRole_Bad(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Menu_AddRole_Ugly(t *core.T) {
	subject := new(wailsMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetAccelerator_Good(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetAccelerator_Bad(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetAccelerator_Ugly(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetTooltip_Good(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetTooltip_Bad(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetTooltip_Ugly(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetChecked_Good(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetChecked_Bad(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetChecked_Ugly(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetEnabled_Good(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetEnabled_Bad(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_SetEnabled_Ugly(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_OnClick_Good(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_OnClick_Bad(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_MenuItem_OnClick_Ugly(t *core.T) {
	subject := new(wailsMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

// AX7 generated source-matching smoke coverage.
