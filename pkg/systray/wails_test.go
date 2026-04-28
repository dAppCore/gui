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
