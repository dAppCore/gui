// pkg/dock/platform.go
package dock

// BounceType controls the style of the macOS dock icon bounce animation.
type BounceType int

const (
	// BounceInformational bounces the dock icon once briefly.
	BounceInformational BounceType = iota
	// BounceCritical bounces the dock icon continuously until the app is focused.
	BounceCritical
)

// Platform abstracts the dock/taskbar backend (Wails v3).
// macOS: dock icon show/hide + badge + bounce + progress.
// Windows: taskbar badge and progress only (show/hide/bounce not supported).
// Linux: not supported — adapter returns nil for all operations.
type Platform interface {
	ShowIcon() error
	HideIcon() error
	SetBadge(label string) error
	RemoveBadge() error
	IsVisible() bool
	// SetProgressBar sets a progress indicator on the dock/taskbar icon.
	// value is in the range [0.0, 1.0]. Use -1.0 to remove the progress bar.
	// p.SetProgressBar(0.5) // 50% progress
	SetProgressBar(value float64) error
	// Bounce animates the dock icon to attract the user's attention.
	// bounceType controls whether the animation loops (BounceCritical) or plays once (BounceInformational).
	// Returns a bounce ID that can be passed to StopBounce.
	Bounce(bounceType BounceType) (int, error)
	// StopBounce cancels a running bounce animation identified by bounceID.
	// p.StopBounce(id)
	StopBounce(bounceID int) error
}
