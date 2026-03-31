package dialog

// TaskOpenFile opens a file picker dialog. Result: []string (selected paths)
type TaskOpenFile struct{ Options OpenFileOptions }

// TaskOpenFileWithOptions opens a file picker, applying caller-supplied options. Result: []string
// paths, _, err := c.PERFORM(TaskOpenFileWithOptions{Title: "Import", AllowMultiple: true})
type TaskOpenFileWithOptions struct {
	Title                string
	Directory            string
	Filename             string
	Filters              []FileFilter
	AllowMultiple        bool
	CanChooseDirectories bool
	CanChooseFiles       bool
	ShowHiddenFiles      bool
}

// TaskSaveFile opens a save file dialog. Result: string (chosen path)
type TaskSaveFile struct{ Options SaveFileOptions }

// TaskSaveFileWithOptions opens a save dialog with caller-supplied options. Result: string
// path, _, err := c.PERFORM(TaskSaveFileWithOptions{Title: "Export", Filename: "out.csv"})
type TaskSaveFileWithOptions struct {
	Title     string
	Directory string
	Filename  string
	Filters   []FileFilter
}

// TaskOpenDirectory opens a directory picker. Result: string
type TaskOpenDirectory struct{ Options OpenDirectoryOptions }

// TaskMessageDialog opens an arbitrary message dialog. Result: string (button clicked)
type TaskMessageDialog struct{ Options MessageDialogOptions }

// TaskInfo shows an informational dialog. Result: string (button clicked)
// result, _, err := c.PERFORM(TaskInfo{Title: "Done", Message: "File saved."})
type TaskInfo struct {
	Title   string
	Message string
	Buttons []string
}

// TaskQuestion shows a question dialog. Result: string (button clicked)
// result, _, err := c.PERFORM(TaskQuestion{Title: "Confirm", Message: "Delete?", Buttons: []string{"Yes","No"}})
type TaskQuestion struct {
	Title   string
	Message string
	Buttons []string
}

// TaskWarning shows a warning dialog. Result: string (button clicked)
// result, _, err := c.PERFORM(TaskWarning{Title: "Warning", Message: "File exists."})
type TaskWarning struct {
	Title   string
	Message string
	Buttons []string
}

// TaskError shows an error dialog. Result: string (button clicked)
// result, _, err := c.PERFORM(TaskError{Title: "Error", Message: "Write failed."})
type TaskError struct {
	Title   string
	Message string
	Buttons []string
}
