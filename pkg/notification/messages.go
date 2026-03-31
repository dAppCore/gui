package notification

// QueryPermission returns current notification permission status. Result: PermissionStatus
type QueryPermission struct{}

// TaskSend sends a native notification, falling back to dialog on failure.
type TaskSend struct{ Options NotificationOptions }

// TaskRequestPermission requests notification permission from the OS. Result: bool (granted)
type TaskRequestPermission struct{}

// ActionNotificationClicked is broadcast when the user clicks a notification.
type ActionNotificationClicked struct{ ID string }
