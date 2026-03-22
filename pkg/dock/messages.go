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

// --- Actions (broadcasts) ---

// ActionVisibilityChanged is broadcast after a successful TaskShowIcon or TaskHideIcon.
type ActionVisibilityChanged struct{ Visible bool }
