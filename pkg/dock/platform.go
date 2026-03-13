// pkg/dock/platform.go
package dock

// Platform abstracts the dock/taskbar backend (Wails v3).
// macOS: dock icon show/hide + badge.
// Windows: taskbar badge only (show/hide not supported).
// Linux: not supported — adapter returns nil for all operations.
type Platform interface {
	ShowIcon() error
	HideIcon() error
	SetBadge(label string) error
	RemoveBadge() error
	IsVisible() bool
}
