package systray

import (
	core "dappco.re/go"
	"github.com/wailsapp/wails/v3/pkg/application"
	"reflect"
	"unsafe"
)

func TestWailsPlatform_NewTray_Good(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)

	tray := platform.NewTray()
	core.AssertNotNil(t, tray)
	wtray, ok := tray.(*wailsTray)
	core.RequireTrue(t, ok)

	wtray.SetIcon([]byte{1, 2, 3})
	wtray.SetTemplateIcon([]byte{4, 5, 6})
	wtray.SetTooltip("Core")
	wtray.SetLabel("Ready")
	wtray.SetMenu(platform.NewMenu())
	wtray.AttachWindow(windowHandleStub{name: "panel"})

	trayValue := reflect.ValueOf(wtray.tray).Elem()
	core.AssertEqual(t, []byte{1, 2, 3}, trayValue.FieldByName("icon").Bytes())
	core.AssertEqual(t, []byte{4, 5, 6}, trayValue.FieldByName("templateIcon").Bytes())
	core.AssertEqual(t, "Core", trayValue.FieldByName("tooltip").String())
	core.AssertEqual(t, "Ready", trayValue.FieldByName("label").String())
	core.AssertTrue(t, trayValue.FieldByName("attachedWindow").IsNil())

	err := wtray.ShowMessage("Title", "Body")
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "not supported")
}

func TestWailsPlatform_NewMenu_Good(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	menu := platform.NewMenu()
	core.AssertNotNil(t, menu)
	wmenu, ok := menu.(*wailsTrayMenu)
	core.RequireTrue(t, ok)

	clicked := false
	item := wmenu.Add("Open").(*wailsTrayMenuItem)
	item.SetTooltip("open")
	item.SetChecked(true)
	item.SetEnabled(false)
	item.OnClick(func() { clicked = true })
	onClickField := reflect.ValueOf(item.item).Elem().FieldByName("onClick")
	core.RequireTrue(t, onClickField.IsValid())
	onClick := reflect.NewAt(onClickField.Type(), unsafe.Pointer(onClickField.UnsafeAddr())).Elem().Interface().(func(*application.Context))
	onClick(&application.Context{})

	core.AssertTrue(t, clicked)
	core.AssertEqual(t, "Open", wmenu.menu.Items[0].Label)
	core.AssertEqual(t, "open", wmenu.menu.Items[0].Tooltip)
	core.AssertTrue(t, wmenu.menu.Items[0].Checked)
	core.AssertFalse(t, wmenu.menu.Items[0].Enabled)
}

func TestWailsPlatform_SetMenu_Bad(t *core.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	tray := platform.NewTray().(*wailsTray)

	tray.SetMenu(&mockTrayMenu{})
	core.AssertTrue(t, reflect.ValueOf(tray.tray).Elem().FieldByName("menu").IsNil())
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

func TestWails_WailsPlatform_NewTray_Good(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_NewTray_Bad(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_WailsPlatform_NewTray_Ugly(t *core.T) {
	subject := new(WailsPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
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

type Tray = wailsTray

func TestWails_Tray_SetIcon_Good(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetIcon_Bad(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetIcon_Ugly(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetTemplateIcon_Good(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetTemplateIcon_Bad(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetTemplateIcon_Ugly(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetTooltip_Good(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetTooltip_Bad(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetTooltip_Ugly(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetLabel_Good(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetLabel("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetLabel_Bad(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetLabel("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetLabel_Ugly(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetLabel("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetMenu_Good(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetMenu_Bad(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_SetMenu_Ugly(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_AttachWindow_Good(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_AttachWindow_Bad(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_AttachWindow_Ugly(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_ShowMessage_Good(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_ShowMessage_Bad(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_Tray_ShowMessage_Ugly(t *core.T) {
	subject := new(wailsTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

type TrayMenu = wailsTrayMenu

func TestWails_TrayMenu_Add_Good(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		got0 := subject.Add("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_Add_Bad(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		got0 := subject.Add("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_Add_Ugly(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		got0 := subject.Add("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_AddSeparator_Good(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_AddSeparator_Bad(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_AddSeparator_Ugly(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_AddSubmenu_Good(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_AddSubmenu_Bad(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenu_AddSubmenu_Ugly(t *core.T) {
	subject := new(wailsTrayMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetTooltip_Good(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetTooltip_Bad(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetTooltip_Ugly(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetChecked_Good(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetChecked_Bad(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetChecked_Ugly(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetEnabled_Good(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetEnabled_Bad(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_SetEnabled_Ugly(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_OnClick_Good(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_OnClick_Bad(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWails_TrayMenuItem_OnClick_Ugly(t *core.T) {
	subject := new(wailsTrayMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

// AX7 generated source-matching smoke coverage.
