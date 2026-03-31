package notification

// QueryPermission returns current notification permission status. Result: PermissionStatus
type QueryPermission struct{}

// TaskSend sends a native notification, falling back to dialog on failure.
type TaskSend struct{ Options NotificationOptions }

// TaskRequestPermission requests notification permission from the OS. Result: bool (granted)
type TaskRequestPermission struct{}

// TaskRevokePermission revokes previously granted notification permission.
// Result: nil
// _, _, err := c.PERFORM(TaskRevokePermission{})
type TaskRevokePermission struct{}

// TaskRegisterCategory registers a notification category with its action buttons.
// Must be called before sending notifications that use that category.
// _, _, err := c.PERFORM(TaskRegisterCategory{Category: NotificationCategory{ID: "message", Actions: [...]}})
type TaskRegisterCategory struct{ Category NotificationCategory }

// ActionNotificationClicked is broadcast when the user clicks a notification.
type ActionNotificationClicked struct{ ID string }

// ActionNotificationActionTriggered is broadcast when the user taps an action button on a notification.
type ActionNotificationActionTriggered struct {
	NotificationID string `json:"notificationId"`
	ActionID       string `json:"actionId"`
}

// ActionNotificationDismissed is broadcast when the user dismisses a notification without acting on it.
type ActionNotificationDismissed struct {
	NotificationID string `json:"notificationId"`
}
