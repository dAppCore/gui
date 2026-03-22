// pkg/menu/platform.go
package menu

// Platform abstracts the menu backend.
type Platform interface {
	NewMenu() PlatformMenu
	SetApplicationMenu(menu PlatformMenu)
}

// PlatformMenu is a live menu handle.
type PlatformMenu interface {
	Add(label string) PlatformMenuItem
	AddSeparator()
	AddSubmenu(label string) PlatformMenu
	// Roles — macOS menu roles
	AddRole(role MenuRole)
}

// PlatformMenuItem is a single menu item.
type PlatformMenuItem interface {
	SetAccelerator(accel string) PlatformMenuItem
	SetTooltip(text string) PlatformMenuItem
	SetChecked(checked bool) PlatformMenuItem
	SetEnabled(enabled bool) PlatformMenuItem
	OnClick(fn func()) PlatformMenuItem
}

// MenuRole is a predefined platform menu role.
type MenuRole int

const (
	RoleAppMenu MenuRole = iota
	RoleFileMenu
	RoleEditMenu
	RoleViewMenu
	RoleWindowMenu
	RoleHelpMenu
)
