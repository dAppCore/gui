// pkg/notification/messages.go
package notification

// QueryPermission checks notification authorisation. Result: PermissionStatus
type QueryPermission struct{}

// TaskSend sends a notification. Falls back to dialog if platform fails.
type TaskSend struct{ Opts NotificationOptions }

// TaskRequestPermission requests notification authorisation. Result: bool (granted)
type TaskRequestPermission struct{}

// ActionNotificationClicked is broadcast when a notification is clicked (future).
type ActionNotificationClicked struct{ ID string }
