// SPDX-License-Identifier: EUPL-1.2

package gui

import (
	core "dappco.re/go"
	"dappco.re/go/gui/pkg/window"
)

// OpenWindow shows + focuses a registered window by name. Composes the
// three actions a "tray click opens chat" / "Cmd+J toggles popover"
// dispatch normally fires:
//
//   1. window.restore — undoes a minimised state
//   2. window.set_visibility{Visible: true} — un-hides
//   3. window.focus — raises + activates
//
// Returns false (and is a no-op) if the named window is not registered.
// Returns true on success. Callers that need the wrapped Result types
// can call the underlying actions directly.
//
//	if !gui.OpenWindow(c, "chat") {
//	    // window not in registry — caller chooses fallback
//	}
//
// Consumer-side this replaces a ~3-action dispatch sequence per call
// site. The lthn-side openWindow() helper that did exactly this is the
// canonical example.
func OpenWindow(c *core.Core, name string) bool {
	if c == nil || name == "" {
		return false
	}
	if !WindowExists(c, name) {
		return false
	}
	ctx := core.Background()
	c.Action("window.restore").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: window.TaskRestore{Name: name}},
	))
	c.Action("window.set_visibility").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetVisibility{Name: name, Visible: true}},
	))
	c.Action("window.focus").Run(ctx, core.NewOptions(
		core.Option{Key: "task", Value: window.TaskFocus{Name: name}},
	))
	return true
}

// HideWindow hides a registered window. Composition of
// window.set_visibility{Visible: false}. Returns false if the name
// isn't in the registry, true on dispatch success.
//
//	gui.HideWindow(c, "chat")
func HideWindow(c *core.Core, name string) bool {
	if c == nil || name == "" {
		return false
	}
	if !WindowExists(c, name) {
		return false
	}
	r := c.Action("window.set_visibility").Run(core.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetVisibility{Name: name, Visible: false}},
	))
	return r.OK
}

// WindowExists reports whether a window with the given name has been
// registered (via GuiConfig.WindowRegistry or an ad-hoc window.open
// call). Wraps the QueryWindowByName query.
//
//	if gui.WindowExists(c, "chat") { … }
func WindowExists(c *core.Core, name string) bool {
	if c == nil || name == "" {
		return false
	}
	r := c.QUERY(window.QueryWindowByName{Name: name})
	return r.OK && r.Value != nil
}
