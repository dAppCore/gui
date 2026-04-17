package application

import "sync"

// DialogType identifies the visual style of a message dialog.
type DialogType int

const (
	InfoDialogType DialogType = iota
	QuestionDialogType
	WarningDialogType
	ErrorDialogType
)

// FileFilter restricts which files are shown in a file dialog.
//
//	filter := FileFilter{DisplayName: "Images", Pattern: "*.png;*.jpg"}
type FileFilter struct {
	DisplayName string
	Pattern     string
}

// OpenFileDialogOptions configures an open-file dialog.
//
//	opts := &OpenFileDialogOptions{Title: "Choose file", AllowsMultipleSelection: true}
type OpenFileDialogOptions struct {
	Title                   string
	Directory               string
	Filters                 []FileFilter
	AllowsMultipleSelection bool
	CanChooseDirectories    bool
	CanChooseFiles          bool
	ShowHiddenFiles         bool
}

// SaveFileDialogOptions configures a save-file dialog.
//
//	opts := &SaveFileDialogOptions{Title: "Save report", Directory: "/tmp"}
type SaveFileDialogOptions struct {
	Title           string
	Directory       string
	Filename        string
	Filters         []FileFilter
	ShowHiddenFiles bool
}

// OpenFileDialogStruct is an in-memory open-file dialog.
//
//	dialog := &OpenFileDialogStruct{}
//	dialog.SetTitle("Pick a file")
//	path, _ := dialog.PromptForSingleSelection()
type OpenFileDialogStruct struct {
	mu              sync.RWMutex
	title           string
	directory       string
	filters         []FileFilter
	multipleAllowed bool
	canChooseDirs   bool
	canChooseFiles  bool
	showHidden      bool
	selectedFiles   []string
}

func newOpenFileDialog() *OpenFileDialogStruct {
	return &OpenFileDialogStruct{canChooseFiles: true}
}

// SetOptions applies OpenFileDialogOptions to the dialog.
//
//	dialog.SetOptions(&OpenFileDialogOptions{Title: "Open log", ShowHiddenFiles: true})
func (d *OpenFileDialogStruct) SetOptions(options *OpenFileDialogOptions) {
	if options == nil {
		return
	}
	d.mu.Lock()
	d.title = options.Title
	d.directory = options.Directory
	d.filters = append([]FileFilter(nil), options.Filters...)
	d.multipleAllowed = options.AllowsMultipleSelection
	d.canChooseDirs = options.CanChooseDirectories
	d.canChooseFiles = options.CanChooseFiles
	d.showHidden = options.ShowHiddenFiles
	d.mu.Unlock()
}

// SetTitle sets the dialog window title.
//
//	dialog.SetTitle("Select configuration file")
func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct {
	d.mu.Lock()
	d.title = title
	d.mu.Unlock()
	return d
}

// SetDirectory sets the initial directory shown in the dialog.
//
//	dialog.SetDirectory("/home/user/documents")
func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct {
	d.mu.Lock()
	d.directory = directory
	d.mu.Unlock()
	return d
}

// AddFilter appends a file filter to the dialog.
//
//	dialog.AddFilter("Go source", "*.go")
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct {
	d.mu.Lock()
	d.filters = append(d.filters, FileFilter{DisplayName: displayName, Pattern: pattern})
	d.mu.Unlock()
	return d
}

// SetAllowsMultipleSelection controls whether multiple files can be selected.
//
//	dialog.SetAllowsMultipleSelection(true)
func (d *OpenFileDialogStruct) SetAllowsMultipleSelection(allow bool) *OpenFileDialogStruct {
	d.mu.Lock()
	d.multipleAllowed = allow
	d.mu.Unlock()
	return d
}

// SetSelectedFiles injects pre-selected files for stub testing.
//
//	dialog.SetSelectedFiles([]string{"/tmp/a.txt", "/tmp/b.txt"})
func (d *OpenFileDialogStruct) SetSelectedFiles(paths []string) {
	d.mu.Lock()
	d.selectedFiles = append([]string(nil), paths...)
	d.mu.Unlock()
}

// PromptForSingleSelection returns the first injected file path, or "" if none.
//
//	path, err := dialog.PromptForSingleSelection()
//	if err != nil { return err }
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.selectedFiles) > 0 {
		return d.selectedFiles[0], nil
	}
	return "", nil
}

// PromptForMultipleSelection returns all injected file paths.
//
//	paths, err := dialog.PromptForMultipleSelection()
//	for _, p := range paths { process(p) }
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return append([]string(nil), d.selectedFiles...), nil
}

// SaveFileDialogStruct is an in-memory save-file dialog.
//
//	dialog := &SaveFileDialogStruct{}
//	dialog.SetFilename("report.pdf")
//	path, _ := dialog.PromptForSingleSelection()
type SaveFileDialogStruct struct {
	mu           sync.RWMutex
	title        string
	directory    string
	filename     string
	filters      []FileFilter
	showHidden   bool
	selectedPath string
}

func newSaveFileDialog() *SaveFileDialogStruct {
	return &SaveFileDialogStruct{}
}

// SetOptions applies SaveFileDialogOptions to the dialog.
//
//	dialog.SetOptions(&SaveFileDialogOptions{Title: "Export", Filename: "data.json"})
func (d *SaveFileDialogStruct) SetOptions(options *SaveFileDialogOptions) {
	if options == nil {
		return
	}
	d.mu.Lock()
	d.title = options.Title
	d.directory = options.Directory
	d.filename = options.Filename
	d.filters = append([]FileFilter(nil), options.Filters...)
	d.showHidden = options.ShowHiddenFiles
	d.mu.Unlock()
}

// SetTitle sets the dialog window title.
//
//	dialog.SetTitle("Export configuration")
func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct {
	d.mu.Lock()
	d.title = title
	d.mu.Unlock()
	return d
}

// SetDirectory sets the initial directory shown in the dialog.
//
//	dialog.SetDirectory("/home/user/exports")
func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct {
	d.mu.Lock()
	d.directory = directory
	d.mu.Unlock()
	return d
}

// SetFilename sets the default filename shown in the dialog.
//
//	dialog.SetFilename("backup-2026.tar.gz")
func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct {
	d.mu.Lock()
	d.filename = filename
	d.mu.Unlock()
	return d
}

// AddFilter appends a file filter to the dialog.
//
//	dialog.AddFilter("JSON files", "*.json")
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct {
	d.mu.Lock()
	d.filters = append(d.filters, FileFilter{DisplayName: displayName, Pattern: pattern})
	d.mu.Unlock()
	return d
}

// SetSelectedPath injects the path returned by PromptForSingleSelection for stub testing.
//
//	dialog.SetSelectedPath("/tmp/output.csv")
func (d *SaveFileDialogStruct) SetSelectedPath(path string) {
	d.mu.Lock()
	d.selectedPath = path
	d.mu.Unlock()
}

// PromptForSingleSelection returns the injected save path, or "" if none.
//
//	path, err := dialog.PromptForSingleSelection()
//	if err != nil { return err }
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.selectedPath, nil
}

// MessageButton represents a button in a message dialog.
type MessageButton struct {
	Label     string
	IsDefault bool
	IsCancel  bool
}

// MessageDialog is an in-memory message dialog (info, question, warning, error).
//
//	dialog := &MessageDialog{dialogType: InfoDialogType}
//	dialog.SetTitle("Done").SetMessage("File saved successfully.")
//	_ = dialog.Show()
type MessageDialog struct {
	mu            sync.RWMutex
	dialogType    DialogType
	title         string
	message       string
	buttons       []MessageButton
	clickedButton string
}

func newMessageDialog(dialogType DialogType) *MessageDialog {
	return &MessageDialog{dialogType: dialogType}
}

// SetTitle sets the dialog window title.
//
//	dialog.SetTitle("Confirm deletion")
func (d *MessageDialog) SetTitle(title string) *MessageDialog {
	d.mu.Lock()
	d.title = title
	d.mu.Unlock()
	return d
}

// SetMessage sets the body text of the dialog.
//
//	dialog.SetMessage("Are you sure you want to delete this file?")
func (d *MessageDialog) SetMessage(message string) *MessageDialog {
	d.mu.Lock()
	d.message = message
	d.mu.Unlock()
	return d
}

// AddButton appends a button to the dialog.
//
//	dialog.AddButton("Yes").AddButton("No")
func (d *MessageDialog) AddButton(label string) *MessageDialog {
	d.mu.Lock()
	d.buttons = append(d.buttons, MessageButton{Label: label})
	d.mu.Unlock()
	return d
}

// SetDefaultButton marks the named button as the default action.
//
//	dialog.SetDefaultButton("OK")
func (d *MessageDialog) SetDefaultButton(label string) *MessageDialog {
	d.mu.Lock()
	for index := range d.buttons {
		d.buttons[index].IsDefault = d.buttons[index].Label == label
	}
	d.mu.Unlock()
	return d
}

// SetCancelButton marks the named button as the cancel action.
//
//	dialog.SetCancelButton("Cancel")
func (d *MessageDialog) SetCancelButton(label string) *MessageDialog {
	d.mu.Lock()
	for index := range d.buttons {
		d.buttons[index].IsCancel = d.buttons[index].Label == label
	}
	d.mu.Unlock()
	return d
}

// SetButtonClickedForStub injects which button was clicked for stub testing.
//
//	dialog.SetButtonClickedForStub("Yes")
//	result, _ := dialog.Show()
func (d *MessageDialog) SetButtonClickedForStub(label string) {
	d.mu.Lock()
	d.clickedButton = label
	d.mu.Unlock()
}

// Show displays the dialog and returns the label of the clicked button.
//
//	clicked, err := dialog.Show()
//	if clicked == "Yes" { deleteFile() }
func (d *MessageDialog) Show() (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.clickedButton, nil
}

// DialogManager manages dialog operations in-memory.
//
//	manager := &DialogManager{}
//	dialog := manager.Info().SetTitle("Done").SetMessage("Saved.")
type DialogManager struct {
	mu sync.RWMutex
}

// OpenFile creates an open-file dialog.
//
//	dialog := manager.OpenFile()
//	dialog.SetTitle("Choose config")
//	path, _ := dialog.PromptForSingleSelection()
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct {
	return newOpenFileDialog()
}

// OpenFileWithOptions creates an open-file dialog pre-configured from options.
//
//	dialog := manager.OpenFileWithOptions(&OpenFileDialogOptions{Title: "Select log"})
func (dm *DialogManager) OpenFileWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct {
	dialog := newOpenFileDialog()
	dialog.SetOptions(options)
	return dialog
}

// SaveFile creates a save-file dialog.
//
//	dialog := manager.SaveFile()
//	dialog.SetFilename("export.csv")
//	path, _ := dialog.PromptForSingleSelection()
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct {
	return newSaveFileDialog()
}

// SaveFileWithOptions creates a save-file dialog pre-configured from options.
//
//	dialog := manager.SaveFileWithOptions(&SaveFileDialogOptions{Title: "Export data"})
func (dm *DialogManager) SaveFileWithOptions(options *SaveFileDialogOptions) *SaveFileDialogStruct {
	dialog := newSaveFileDialog()
	dialog.SetOptions(options)
	return dialog
}

// Info creates an information message dialog.
//
//	manager.Info().SetTitle("Done").SetMessage("File saved.").Show()
func (dm *DialogManager) Info() *MessageDialog {
	return newMessageDialog(InfoDialogType)
}

// Question creates a question message dialog.
//
//	manager.Question().SetTitle("Confirm").SetMessage("Delete file?").AddButton("Yes").AddButton("No").Show()
func (dm *DialogManager) Question() *MessageDialog {
	return newMessageDialog(QuestionDialogType)
}

// Warning creates a warning message dialog.
//
//	manager.Warning().SetTitle("Warning").SetMessage("Disk almost full.").Show()
func (dm *DialogManager) Warning() *MessageDialog {
	return newMessageDialog(WarningDialogType)
}

// Error creates an error message dialog.
//
//	manager.Error().SetTitle("Error").SetMessage("Operation failed.").Show()
func (dm *DialogManager) Error() *MessageDialog {
	return newMessageDialog(ErrorDialogType)
}
