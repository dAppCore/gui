// pkg/systray/wails.go
package systray

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform using Wails v3.
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) NewTray() PlatformTray {
	return &wailsTray{tray: wp.app.SystemTray.New(), app: wp.app}
}

func (wp *WailsPlatform) NewMenu() PlatformMenu {
	return &wailsTrayMenu{menu: wp.app.NewMenu()}
}

type wailsTray struct {
	tray *application.SystemTray
	app  *application.App
}

func (wt *wailsTray) SetIcon(data []byte)         { wt.tray.SetIcon(data) }
func (wt *wailsTray) SetTemplateIcon(data []byte) { wt.tray.SetTemplateIcon(data) }
func (wt *wailsTray) SetTooltip(text string)      { wt.tray.SetTooltip(text) }
func (wt *wailsTray) SetLabel(text string)        { wt.tray.SetLabel(text) }

func (wt *wailsTray) SetMenu(menu PlatformMenu) {
	if wm, ok := menu.(*wailsTrayMenu); ok {
		wt.tray.SetMenu(wm.menu)
	}
}

func (wt *wailsTray) AttachWindow(w WindowHandle) {
	// Wails systray AttachWindow expects an application.Window interface.
	// The caller must pass an appropriate wrapper.
}

// wailsTrayMenu wraps *application.Menu for the PlatformMenu interface.
type wailsTrayMenu struct {
	menu *application.Menu
}

func (m *wailsTrayMenu) Add(label string) PlatformMenuItem {
	return &wailsTrayMenuItem{item: m.menu.Add(label)}
}

func (m *wailsTrayMenu) AddSeparator() {
	m.menu.AddSeparator()
}

func (m *wailsTrayMenu) AddSubmenu(label string) PlatformMenu {
	return &wailsTrayMenu{menu: m.menu.AddSubmenu(label)}
}

// wailsTrayMenuItem wraps *application.MenuItem for the PlatformMenuItem interface.
type wailsTrayMenuItem struct {
	item *application.MenuItem
}

func (mi *wailsTrayMenuItem) SetTooltip(text string)  { mi.item.SetTooltip(text) }
func (mi *wailsTrayMenuItem) SetChecked(checked bool) { mi.item.SetChecked(checked) }
func (mi *wailsTrayMenuItem) SetEnabled(enabled bool) { mi.item.SetEnabled(enabled) }
func (mi *wailsTrayMenuItem) OnClick(fn func()) {
	mi.item.OnClick(func(ctx *application.Context) { fn() })
}
func (mi *wailsTrayMenuItem) AddSubmenu() PlatformMenu {
	// Wails doesn't have a direct AddSubmenu on MenuItem — use Menu.AddSubmenu instead
	return &wailsTrayMenu{menu: application.NewMenu()}
}
