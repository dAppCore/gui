package systray

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestWailsPlatform_NewTray_Good(t *testing.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)

	tray := platform.NewTray()
	require.NotNil(t, tray)
	wtray, ok := tray.(*wailsTray)
	require.True(t, ok)

	wtray.SetIcon([]byte{1, 2, 3})
	wtray.SetTemplateIcon([]byte{4, 5, 6})
	wtray.SetTooltip("Core")
	wtray.SetLabel("Ready")
	wtray.SetMenu(platform.NewMenu())
	wtray.AttachWindow(windowHandleStub{name: "panel"})

	trayValue := reflect.ValueOf(wtray.tray).Elem()
	assert.Equal(t, []byte{1, 2, 3}, trayValue.FieldByName("icon").Bytes())
	assert.Equal(t, []byte{4, 5, 6}, trayValue.FieldByName("templateIcon").Bytes())
	assert.Equal(t, "Core", trayValue.FieldByName("tooltip").String())
	assert.Equal(t, "Ready", trayValue.FieldByName("label").String())
	assert.True(t, trayValue.FieldByName("attachedWindow").IsNil())

	err := wtray.ShowMessage("Title", "Body")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func TestWailsPlatform_NewMenu_Good(t *testing.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	menu := platform.NewMenu()
	require.NotNil(t, menu)
	wmenu, ok := menu.(*wailsTrayMenu)
	require.True(t, ok)

	clicked := false
	item := wmenu.Add("Open").(*wailsTrayMenuItem)
	item.SetTooltip("open")
	item.SetChecked(true)
	item.SetEnabled(false)
	item.OnClick(func() { clicked = true })
	onClickField := reflect.ValueOf(item.item).Elem().FieldByName("onClick")
	require.True(t, onClickField.IsValid())
	onClick := reflect.NewAt(onClickField.Type(), unsafe.Pointer(onClickField.UnsafeAddr())).Elem().Interface().(func(*application.Context))
	onClick(&application.Context{})

	assert.True(t, clicked)
	assert.Equal(t, "Open", wmenu.menu.Items[0].Label)
	assert.Equal(t, "open", wmenu.menu.Items[0].Tooltip)
	assert.True(t, wmenu.menu.Items[0].Checked)
	assert.False(t, wmenu.menu.Items[0].Enabled)
}

func TestWailsPlatform_SetMenu_Bad(t *testing.T) {
	app := &application.App{}
	platform := NewWailsPlatform(app)
	tray := platform.NewTray().(*wailsTray)

	tray.SetMenu(&mockTrayMenu{})
	assert.True(t, reflect.ValueOf(tray.tray).Elem().FieldByName("menu").IsNil())
}
