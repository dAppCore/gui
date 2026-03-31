package notification

type QueryPermission struct{}

type TaskSend struct{ Options NotificationOptions }

type TaskRequestPermission struct{}

type ActionNotificationClicked struct{ ID string }
