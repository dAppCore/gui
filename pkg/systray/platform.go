// pkg/systray/platform.go
package systray

// Platform abstracts the system tray backend.
type Platform interface {
	NewTray() PlatformTray
	NewMenu() PlatformMenu // Menu factory for building tray menus
}

// PlatformTray is a live tray handle from the backend.
type PlatformTray interface {
	SetIcon(data []byte)
	SetTemplateIcon(data []byte)
	SetTooltip(text string)
	SetLabel(text string)
	SetMenu(menu PlatformMenu)
	AttachWindow(w WindowHandle)
}

// PlatformMenu is a tray menu built by the backend.
type PlatformMenu interface {
	Add(label string) PlatformMenuItem
	AddSeparator()
}

// PlatformMenuItem is a single item in a tray menu.
type PlatformMenuItem interface {
	SetTooltip(text string)
	SetChecked(checked bool)
	SetEnabled(enabled bool)
	OnClick(fn func())
	AddSubmenu() PlatformMenu
}

// WindowHandle is a cross-package interface for window operations.
// Defined locally to avoid circular imports (display imports systray).
// pkg/window.PlatformWindow satisfies this implicitly.
type WindowHandle interface {
	Name() string
	Show()
	Hide()
	SetPosition(x, y int)
	SetSize(width, height int)
}
