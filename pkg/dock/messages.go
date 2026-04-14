// pkg/dock/messages.go
package dock

// --- Queries (read-only) ---

// QueryVisible returns whether the dock icon is visible. Result: bool
type QueryVisible struct{}

// --- Tasks (side-effects) ---

// TaskShowIcon shows the dock/taskbar icon. Result: nil
type TaskShowIcon struct{}

// TaskHideIcon hides the dock/taskbar icon. Result: nil
type TaskHideIcon struct{}

// TaskSetBadge sets the dock/taskbar badge label.
// Empty string "" shows the default system badge indicator.
// Numeric "3", "99" shows unread count. Text "New", "Paused" shows brief status.
// Result: nil
type TaskSetBadge struct{ Label string }

// TaskRemoveBadge removes the dock/taskbar badge. Result: nil
type TaskRemoveBadge struct{}

// TaskSetProgressBar sets the progress indicator on the dock/taskbar icon.
// Value must be in the range [0.0, 1.0]. Use -1.0 to remove the bar.
// Result: nil
// _, _, err := c.PERFORM(TaskSetProgressBar{Value: 0.75})
type TaskSetProgressBar struct{ Value float64 }

// TaskBounce animates the dock icon to attract the user's attention. Result: int (bounce ID)
// _, result, err := c.PERFORM(TaskBounce{Type: BounceCritical})
// bounceID := result.(int)
type TaskBounce struct{ Type BounceType }

// TaskStopBounce cancels a running dock bounce animation. Result: nil
// _, _, err := c.PERFORM(TaskStopBounce{BounceID: bounceID})
type TaskStopBounce struct{ BounceID int }

// --- Actions (broadcasts) ---

// ActionVisibilityChanged is broadcast after a successful TaskShowIcon or TaskHideIcon.
type ActionVisibilityChanged struct{ Visible bool }

// ActionProgressChanged is broadcast after a successful TaskSetProgressBar.
type ActionProgressChanged struct{ Value float64 }

// ActionBounceStarted is broadcast after a successful TaskBounce.
type ActionBounceStarted struct {
	BounceID int        `json:"bounceId"`
	Type     BounceType `json:"type"`
}
