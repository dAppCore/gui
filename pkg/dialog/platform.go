// pkg/dialog/platform.go
package dialog

// Platform abstracts the native dialog backend.
type Platform interface {
	OpenFile(opts OpenFileOptions) ([]string, error)
	SaveFile(opts SaveFileOptions) (string, error)
	OpenDirectory(opts OpenDirectoryOptions) (string, error)
	MessageDialog(opts MessageDialogOptions) (string, error)
}

// DialogType represents the type of message dialog.
type DialogType int

const (
	DialogInfo DialogType = iota
	DialogWarning
	DialogError
	DialogQuestion
)

// OpenFileOptions contains options for the open file dialog.
type OpenFileOptions struct {
	Title         string       `json:"title,omitempty"`
	Directory     string       `json:"directory,omitempty"`
	Filename      string       `json:"filename,omitempty"`
	Filters       []FileFilter `json:"filters,omitempty"`
	AllowMultiple bool         `json:"allowMultiple,omitempty"`
}

// SaveFileOptions contains options for the save file dialog.
type SaveFileOptions struct {
	Title     string       `json:"title,omitempty"`
	Directory string       `json:"directory,omitempty"`
	Filename  string       `json:"filename,omitempty"`
	Filters   []FileFilter `json:"filters,omitempty"`
}

// OpenDirectoryOptions contains options for the directory picker.
type OpenDirectoryOptions struct {
	Title         string `json:"title,omitempty"`
	Directory     string `json:"directory,omitempty"`
	AllowMultiple bool   `json:"allowMultiple,omitempty"`
}

// MessageDialogOptions contains options for a message dialog.
type MessageDialogOptions struct {
	Type    DialogType `json:"type"`
	Title   string     `json:"title"`
	Message string     `json:"message"`
	Buttons []string   `json:"buttons,omitempty"`
}

// FileFilter represents a file type filter for dialogs.
type FileFilter struct {
	DisplayName string   `json:"displayName"`
	Pattern     string   `json:"pattern"`
	Extensions  []string `json:"extensions,omitempty"`
}
