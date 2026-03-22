// pkg/dialog/messages.go
package dialog

// TaskOpenFile shows an open file dialog. Result: []string (paths)
type TaskOpenFile struct{ Opts OpenFileOptions }

// TaskSaveFile shows a save file dialog. Result: string (path)
type TaskSaveFile struct{ Opts SaveFileOptions }

// TaskOpenDirectory shows a directory picker. Result: string (path)
type TaskOpenDirectory struct{ Opts OpenDirectoryOptions }

// TaskMessageDialog shows a message dialog. Result: string (button clicked)
type TaskMessageDialog struct{ Opts MessageDialogOptions }
