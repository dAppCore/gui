// pkg/notification/platform.go
package notification

// Platform abstracts the native notification backend.
type Platform interface {
	Send(opts NotificationOptions) error
	RequestPermission() (bool, error)
	CheckPermission() (bool, error)
}

// NotificationSeverity indicates the severity for dialog fallback.
type NotificationSeverity int

const (
	SeverityInfo NotificationSeverity = iota
	SeverityWarning
	SeverityError
)

// NotificationOptions contains options for sending a notification.
type NotificationOptions struct {
	ID       string               `json:"id,omitempty"`
	Title    string               `json:"title"`
	Message  string               `json:"message"`
	Subtitle string               `json:"subtitle,omitempty"`
	Severity NotificationSeverity `json:"severity,omitempty"`
}

// PermissionStatus indicates whether notifications are authorised.
type PermissionStatus struct {
	Granted bool `json:"granted"`
}
