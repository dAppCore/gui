package systray

import (
	core "dappco.re/go"
)

func TestMockPlatform_NewTray_Good(t *core.T) {
	p := NewMockPlatform()
	tray := p.NewTray()
	core.AssertNotNil(t, tray)

	mockTray := tray.(*exportedMockTray)
	tray.SetIcon([]byte{1, 2, 3})
	tray.SetTemplateIcon([]byte{4, 5, 6})
	tray.SetTooltip("Core")
	tray.SetLabel("Ready")
	tray.SetMenu(p.NewMenu())
	tray.AttachWindow(windowHandleStub{name: "panel"})

	core.AssertEqual(t, []byte{1, 2, 3}, mockTray.icon)
	core.AssertEqual(t, []byte{4, 5, 6}, mockTray.templateIcon)
	core.AssertEqual(t, "Core", mockTray.tooltip)
	core.AssertEqual(t, "Ready", mockTray.label)
	core.AssertNotNil(t, mockTray)
}

func TestMockPlatform_NewMenu_Bad(t *core.T) {
	p := NewMockPlatform()
	menu := p.NewMenu()
	core.AssertNotNil(t, menu)
	_, ok := menu.(*exportedMockMenu)
	core.AssertTrue(t, ok)
}

func TestMockPlatform_NewTray_Ugly(t *core.T) {
	p := NewMockPlatform()
	tray := p.NewTray().(*exportedMockTray)
	core.AssertNotNil(t, tray)
	core.AssertNoError(t, tray.ShowMessage("title", "message"))
}

type windowHandleStub struct {
	name string
}

func (w windowHandleStub) Name() string { return w.name }
