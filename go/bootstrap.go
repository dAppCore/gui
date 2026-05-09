// SPDX-License-Identifier: EUPL-1.2

package gui

import (
	core "dappco.re/go"
	"dappco.re/go/gui/pkg/browser"
	"dappco.re/go/gui/pkg/contextmenu"
	"dappco.re/go/gui/pkg/dialog"
	"dappco.re/go/gui/pkg/display"
	"dappco.re/go/gui/pkg/dock"
	"dappco.re/go/gui/pkg/environment"
	"dappco.re/go/gui/pkg/keybinding"
	"dappco.re/go/gui/pkg/lifecycle"
	"dappco.re/go/gui/pkg/menu"
	"dappco.re/go/gui/pkg/notification"
	"dappco.re/go/gui/pkg/screen"
	"dappco.re/go/gui/pkg/systray"
	"dappco.re/go/gui/pkg/webview"
	"dappco.re/go/gui/pkg/window"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Bootstrap returns the full set of [core.CoreOption] a desktop consumer
// needs to wire the GUI service stack in one append, instead of hand-wiring
// every sub-service.
//
// What's included:
//   - "gui"          — the top-level gui shell service
//   - "window"       — window manager + layout/state persistence (NewLayoutManager
//     uses DIR_CONFIG / Core/layouts.json by default)
//   - "display"      — screens, work areas, dialogs, clipboard, system tray
//   - "webview"      — JS evaluation, console capture, DOM queries (CDP-backed)
//   - "menu"         — native application menu
//   - "systray"      — system tray icon + menu
//   - "browser"      — open external URLs/files in the OS default app
//     (forge/mantis/docs links; native target=_blank replacement)
//   - "notification" — native OS notifications (macOS Notification Center,
//     Windows Toast, Linux D-Bus). macOS requires bundle+sign to fire.
//   - "lifecycle"    — app lifecycle events (started, will-terminate,
//     did-become-active, did-resign-active, opened-with-file). Subscribers
//     receive the corresponding c.Action dispatches.
//   - "dialog"       — native file open/save/directory pickers + info/
//     warning/error/question message dialogs. File pickers block until
//     the user resolves the dialog; message dialogs adapt Wails's async
//     button-callback model to a synchronous "which label was clicked"
//     return.
//   - "contextmenu"  — register named native context menus. Frontend
//     opts an element in via the CSS custom property
//     --custom-contextmenu: <name>; (with optional
//     --custom-contextmenu-data for per-element payload). Item clicks
//     dispatch ActionItemClicked on the consumer's core.
//   - "keybinding"   — register global accelerators at runtime
//     (Cmd+S, Ctrl+P, F1, etc.). Shares the same Wails key-binding map
//     as application.Options.KeyBindings — runtime registration wins
//     over boot-time when accelerators collide.
//   - "dock"         — macOS dock + Windows taskbar icon visibility +
//     badge label. Progress bar / bounce are accepted but no-op until
//     Wails exposes them upstream.
//   - "environment"  — OS / arch / debug / platform info, dark-mode
//     query + ThemeChanged subscription, accent colour, OpenFileManager,
//     focus-follows-mouse (Linux).
//   - "screen"       — multi-monitor info: GetAll / GetPrimary /
//     GetCurrent (containing-window fallback to primary).
//
// The wails [*application.App] is the only boundary the consumer touches —
// after that, everything runs through the canonical Core IPC pattern
// (c.Action/c.QUERY) so consumers don't need direct wails imports.
//
//	app := application.New(opts) // consumer creates the app
//	coreOpts := []core.CoreOption{ /* your services */ }
//	coreOpts = append(coreOpts, gui.Bootstrap(app)...)
//	c, _ := core.New(coreOpts...)
//
// More gui sub-services (clipboard-as-service, etc.) can be added to
// this Bootstrap as the desktop surface needs them — this is the
// single point to extend.
func Bootstrap(app *application.App) []core.CoreOption {
	if app == nil {
		return nil
	}
	return []core.CoreOption{
		core.WithName("gui", NewService(GuiConfig{})),
		core.WithService(window.Register(window.NewWailsPlatform(app))),
		core.WithService(display.Register(app)),
		core.WithService(webview.Register()),
		core.WithService(menu.Register(menu.NewWailsPlatform(app))),
		core.WithService(systray.Register(systray.NewWailsPlatform(app))),
		core.WithService(browser.Register(browser.NewWailsPlatform(app))),
		core.WithService(notification.Register(notification.NewWailsPlatform(app))),
		core.WithService(lifecycle.Register(lifecycle.NewWailsPlatform(app))),
		core.WithService(dialog.Register(dialog.NewWailsPlatform(app))),
		core.WithService(contextmenu.Register(contextmenu.NewWailsPlatform(app))),
		core.WithService(keybinding.Register(keybinding.NewWailsPlatform(app))),
		core.WithService(dock.Register(dock.NewWailsPlatform(app))),
		core.WithService(environment.Register(environment.NewWailsPlatform(app))),
		core.WithService(screen.Register(screen.NewWailsPlatform(app))),
	}
}
