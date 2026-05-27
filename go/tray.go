// SPDX-License-Identifier: EUPL-1.2

package gui

import (
	"runtime"

	core "dappco.re/go"
	"dappco.re/go/gui/pkg/systray"
)

// TrayItem is an alias for the systray menu item shape. Re-exported so
// consumers can build TrayConfig without importing the systray package
// directly.
type TrayItem = systray.TrayMenuItem

// TrayConfig declares the system tray surface — icon, tooltip, label,
// menu items, and an optional popover window attachment. Applied by
// gui.Service.OnStartup after the systray sub-service has started.
//
// Click routing remains caller-owned: register an action handler with
// core.RegisterAction that switches on
// systray.ActionTrayMenuItemClicked.ActionID. Items declare their
// ActionID so the dispatch is decoupled from the menu shape.
//
// Platform branching: when IconTemplate is true, the icon is set via
// the systray.set_template_icon action on darwin (auto-inverted by the
// OS to match menu-bar light/dark), otherwise via systray.set_icon
// uniformly. Tooltip + Label are always set when non-empty.
type TrayConfig struct {
	// Icon is the PNG bytes for the tray glyph. Empty leaves the
	// systray sub-service's bootstrap default in place.
	Icon []byte
	// IconTemplate marks the icon as a macOS template image. The
	// OS auto-inverts it for light/dark menu bars. No-op on other
	// platforms.
	IconTemplate bool
	// Tooltip is the hover text. On macOS this renders as the
	// menu-bar label TEXT next to the icon (there is no separate
	// tooltip surface), so empty here on darwin keeps the tray
	// icon-only. Non-darwin platforms show it as a real tooltip.
	Tooltip string
	// Label is the menu-bar label text. Empty clears core/gui's
	// systray bootstrap default of "Core".
	Label string
	// Menu is the declarative tray menu. Items support Label, Type
	// ("normal" / "separator" / "checkbox" / "radio"), Checked,
	// Disabled, Tooltip, Submenu, and ActionID for click routing.
	Menu []TrayItem
	// PopoverWindow is the name of a registered window to attach as
	// the tray's popover (clicking the tray icon shows the named
	// window anchored under it). Empty disables the popover
	// attachment; the systray icon stays click-to-show-menu only.
	PopoverWindow string
	// PopoverOffsetY shifts the attached popover N pixels down from
	// the tray icon's anchor point. Useful when the chrome / arrow
	// region needs visual breathing room. Zero = no offset.
	PopoverOffsetY int
}

// applyTrayConfig fires the systray.set_* + attach_window actions for
// the declared config. Called by gui.Service.start() once the systray
// sub-service is registered + started.
func applyTrayConfig(c *core.Core, cfg *TrayConfig) {
	if c == nil || cfg == nil {
		return
	}
	ctx := core.Background()
	if cfg.Icon != nil {
		if cfg.IconTemplate && runtime.GOOS == "darwin" {
			c.Action("systray.set_template_icon").Run(ctx, core.NewOptions(
				core.Option{Key: "task", Value: systray.TaskSetTrayTemplateIcon{Data: cfg.Icon}},
			))
		} else {
			c.Action("systray.set_icon").Run(ctx, core.NewOptions(
				core.Option{Key: "task", Value: systray.TaskSetTrayIcon{Data: cfg.Icon}},
			))
		}
	}
	c.Action("systray.set_tooltip").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayTooltip{Tooltip: cfg.Tooltip}},
	))
	c.Action("systray.set_label").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayLabel{Label: cfg.Label}},
	))
	if len(cfg.Menu) > 0 {
		c.Action("systray.set_menu").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: systray.TaskSetTrayMenu{Items: cfg.Menu}},
		))
	}
	if cfg.PopoverWindow != "" {
		c.Action("systray.attach_window").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: systray.TaskAttachWindow{Name: cfg.PopoverWindow, OffsetY: cfg.PopoverOffsetY}},
		))
	}
}
