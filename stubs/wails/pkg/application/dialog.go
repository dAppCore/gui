package application

// DialogType identifies the visual style of a message dialog.
type DialogType int

const (
	InfoDialogType     DialogType = iota
	QuestionDialogType DialogType = iota
	WarningDialogType  DialogType = iota
	ErrorDialogType    DialogType = iota
)

// FileFilter describes a file type filter for open/save dialogs.
//
//	filter := FileFilter{DisplayName: "Images (*.png;*.jpg)", Pattern: "*.png;*.jpg"}
type FileFilter struct {
	DisplayName string
	Pattern     string
}

// Button is a labelled action in a MessageDialog.
//
//	btn := dialog.AddButton("OK")
//	btn.SetAsDefault().OnClick(func() { ... })
type Button struct {
	Label     string
	IsCancel  bool
	IsDefault bool
	Callback  func()
}

// OnClick registers a click handler on the button and returns itself for chaining.
//
//	btn.OnClick(func() { saveFile() })
func (b *Button) OnClick(callback func()) *Button {
	b.Callback = callback
	return b
}

// SetAsDefault marks this button as the default (Enter key) action.
func (b *Button) SetAsDefault() *Button {
	b.IsDefault = true
	return b
}

// SetAsCancel marks this button as the cancel (Escape key) action.
func (b *Button) SetAsCancel() *Button {
	b.IsCancel = true
	return b
}

// MessageDialogOptions holds configuration for a MessageDialog.
type MessageDialogOptions struct {
	DialogType DialogType
	Title      string
	Message    string
	Buttons    []*Button
	Icon       []byte
}

// MessageDialog is an in-memory message dialog (info / question / warning / error).
//
//	dialog.Info().SetTitle("Done").SetMessage("File saved.").Show()
type MessageDialog struct {
	MessageDialogOptions
}

// SetTitle sets the dialog window title.
//
//	dialog.SetTitle("Confirm Delete")
func (d *MessageDialog) SetTitle(title string) *MessageDialog {
	d.Title = title
	return d
}

// SetMessage sets the body text shown in the dialog.
//
//	dialog.SetMessage("Are you sure?")
func (d *MessageDialog) SetMessage(message string) *MessageDialog {
	d.Message = message
	return d
}

// SetIcon sets the icon bytes shown in the dialog.
func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog {
	d.Icon = icon
	return d
}

// AddButton appends a labelled button and returns it for further configuration.
//
//	btn := dialog.AddButton("Yes")
func (d *MessageDialog) AddButton(label string) *Button {
	btn := &Button{Label: label}
	d.Buttons = append(d.Buttons, btn)
	return btn
}

// AddButtons replaces the button list in bulk.
func (d *MessageDialog) AddButtons(buttons []*Button) *MessageDialog {
	d.Buttons = buttons
	return d
}

// SetDefaultButton marks the given button as the default action.
func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog {
	for _, b := range d.Buttons {
		b.IsDefault = false
	}
	button.IsDefault = true
	return d
}

// SetCancelButton marks the given button as the cancel action.
func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog {
	for _, b := range d.Buttons {
		b.IsCancel = false
	}
	button.IsCancel = true
	return d
}

// AttachToWindow associates the dialog with a parent window (no-op in the stub).
func (d *MessageDialog) AttachToWindow(window *WebviewWindow) *MessageDialog {
	return d
}

// Show presents the dialog. No-op in the stub.
func (d *MessageDialog) Show() {}

func newMessageDialog(dialogType DialogType) *MessageDialog {
	return &MessageDialog{
		MessageDialogOptions: MessageDialogOptions{DialogType: dialogType},
	}
}

// OpenFileDialogOptions configures an OpenFileDialogStruct.
type OpenFileDialogOptions struct {
	CanChooseDirectories            bool
	CanChooseFiles                  bool
	CanCreateDirectories            bool
	ShowHiddenFiles                 bool
	ResolvesAliases                 bool
	AllowsMultipleSelection         bool
	HideExtension                   bool
	CanSelectHiddenExtension        bool
	TreatsFilePackagesAsDirectories bool
	AllowsOtherFileTypes            bool
	Filters                         []FileFilter
	Title                           string
	Message                         string
	ButtonText                      string
	Directory                       string
}

// OpenFileDialogStruct is a builder for file-open dialogs.
//
//	path, err := manager.OpenFile().SetTitle("Pick a file").PromptForSingleSelection()
type OpenFileDialogStruct struct {
	canChooseDirectories            bool
	canChooseFiles                  bool
	canCreateDirectories            bool
	showHiddenFiles                 bool
	resolvesAliases                 bool
	allowsMultipleSelection         bool
	hideExtension                   bool
	canSelectHiddenExtension        bool
	treatsFilePackagesAsDirectories bool
	allowsOtherFileTypes            bool
	filters                         []FileFilter
	title                           string
	message                         string
	buttonText                      string
	directory                       string
}

func newOpenFileDialog() *OpenFileDialogStruct {
	return &OpenFileDialogStruct{
		canChooseFiles:       true,
		canCreateDirectories: true,
	}
}

// SetOptions applies all fields from OpenFileDialogOptions to the dialog.
func (d *OpenFileDialogStruct) SetOptions(options *OpenFileDialogOptions) {
	d.title = options.Title
	d.message = options.Message
	d.buttonText = options.ButtonText
	d.directory = options.Directory
	d.canChooseDirectories = options.CanChooseDirectories
	d.canChooseFiles = options.CanChooseFiles
	d.canCreateDirectories = options.CanCreateDirectories
	d.showHiddenFiles = options.ShowHiddenFiles
	d.resolvesAliases = options.ResolvesAliases
	d.allowsMultipleSelection = options.AllowsMultipleSelection
	d.hideExtension = options.HideExtension
	d.canSelectHiddenExtension = options.CanSelectHiddenExtension
	d.treatsFilePackagesAsDirectories = options.TreatsFilePackagesAsDirectories
	d.allowsOtherFileTypes = options.AllowsOtherFileTypes
	d.filters = options.Filters
}

func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct {
	d.title = title
	return d
}

func (d *OpenFileDialogStruct) SetMessage(message string) *OpenFileDialogStruct {
	d.message = message
	return d
}

func (d *OpenFileDialogStruct) SetButtonText(text string) *OpenFileDialogStruct {
	d.buttonText = text
	return d
}

func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct {
	d.directory = directory
	return d
}

func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct {
	d.canChooseFiles = canChooseFiles
	return d
}

func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct {
	d.canChooseDirectories = canChooseDirectories
	return d
}

func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct {
	d.canCreateDirectories = canCreateDirectories
	return d
}

func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct {
	d.showHiddenFiles = showHiddenFiles
	return d
}

func (d *OpenFileDialogStruct) HideExtension(hideExtension bool) *OpenFileDialogStruct {
	d.hideExtension = hideExtension
	return d
}

func (d *OpenFileDialogStruct) CanSelectHiddenExtension(canSelectHiddenExtension bool) *OpenFileDialogStruct {
	d.canSelectHiddenExtension = canSelectHiddenExtension
	return d
}

func (d *OpenFileDialogStruct) ResolvesAliases(resolvesAliases bool) *OpenFileDialogStruct {
	d.resolvesAliases = resolvesAliases
	return d
}

func (d *OpenFileDialogStruct) AllowsOtherFileTypes(allowsOtherFileTypes bool) *OpenFileDialogStruct {
	d.allowsOtherFileTypes = allowsOtherFileTypes
	return d
}

func (d *OpenFileDialogStruct) TreatsFilePackagesAsDirectories(treats bool) *OpenFileDialogStruct {
	d.treatsFilePackagesAsDirectories = treats
	return d
}

// AddFilter appends a file type filter to the dialog.
//
//	dialog.AddFilter("Images", "*.png;*.jpg")
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct {
	d.filters = append(d.filters, FileFilter{DisplayName: displayName, Pattern: pattern})
	return d
}

func (d *OpenFileDialogStruct) AttachToWindow(window *WebviewWindow) *OpenFileDialogStruct {
	return d
}

// PromptForSingleSelection shows the dialog and returns the chosen path.
// Always returns ("", nil) in the stub.
//
//	path, err := dialog.PromptForSingleSelection()
func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error) {
	return "", nil
}

// PromptForMultipleSelection shows the dialog and returns all chosen paths.
// Always returns (nil, nil) in the stub.
//
//	paths, err := dialog.PromptForMultipleSelection()
func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error) {
	return nil, nil
}

// SaveFileDialogOptions configures a SaveFileDialogStruct.
type SaveFileDialogOptions struct {
	CanCreateDirectories            bool
	ShowHiddenFiles                 bool
	CanSelectHiddenExtension        bool
	AllowOtherFileTypes             bool
	HideExtension                   bool
	TreatsFilePackagesAsDirectories bool
	Title                           string
	Message                         string
	Directory                       string
	Filename                        string
	ButtonText                      string
	Filters                         []FileFilter
}

// SaveFileDialogStruct is a builder for file-save dialogs.
//
//	path, err := manager.SaveFile().SetTitle("Save As").PromptForSingleSelection()
type SaveFileDialogStruct struct {
	canCreateDirectories            bool
	showHiddenFiles                 bool
	canSelectHiddenExtension        bool
	allowOtherFileTypes             bool
	hideExtension                   bool
	treatsFilePackagesAsDirectories bool
	title                           string
	message                         string
	directory                       string
	filename                        string
	buttonText                      string
	filters                         []FileFilter
}

func newSaveFileDialog() *SaveFileDialogStruct {
	return &SaveFileDialogStruct{canCreateDirectories: true}
}

// SetOptions applies all fields from SaveFileDialogOptions to the dialog.
func (d *SaveFileDialogStruct) SetOptions(options *SaveFileDialogOptions) {
	d.title = options.Title
	d.canCreateDirectories = options.CanCreateDirectories
	d.showHiddenFiles = options.ShowHiddenFiles
	d.canSelectHiddenExtension = options.CanSelectHiddenExtension
	d.allowOtherFileTypes = options.AllowOtherFileTypes
	d.hideExtension = options.HideExtension
	d.treatsFilePackagesAsDirectories = options.TreatsFilePackagesAsDirectories
	d.message = options.Message
	d.directory = options.Directory
	d.filename = options.Filename
	d.buttonText = options.ButtonText
	d.filters = options.Filters
}

func (d *SaveFileDialogStruct) SetTitle(title string) *SaveFileDialogStruct {
	d.title = title
	return d
}

func (d *SaveFileDialogStruct) SetMessage(message string) *SaveFileDialogStruct {
	d.message = message
	return d
}

func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct {
	d.directory = directory
	return d
}

func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct {
	d.filename = filename
	return d
}

func (d *SaveFileDialogStruct) SetButtonText(text string) *SaveFileDialogStruct {
	d.buttonText = text
	return d
}

func (d *SaveFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *SaveFileDialogStruct {
	d.canCreateDirectories = canCreateDirectories
	return d
}

func (d *SaveFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *SaveFileDialogStruct {
	d.showHiddenFiles = showHiddenFiles
	return d
}

func (d *SaveFileDialogStruct) CanSelectHiddenExtension(canSelectHiddenExtension bool) *SaveFileDialogStruct {
	d.canSelectHiddenExtension = canSelectHiddenExtension
	return d
}

func (d *SaveFileDialogStruct) AllowsOtherFileTypes(allowOtherFileTypes bool) *SaveFileDialogStruct {
	d.allowOtherFileTypes = allowOtherFileTypes
	return d
}

func (d *SaveFileDialogStruct) HideExtension(hideExtension bool) *SaveFileDialogStruct {
	d.hideExtension = hideExtension
	return d
}

func (d *SaveFileDialogStruct) TreatsFilePackagesAsDirectories(treats bool) *SaveFileDialogStruct {
	d.treatsFilePackagesAsDirectories = treats
	return d
}

// AddFilter appends a file type filter to the dialog.
//
//	dialog.AddFilter("Text files", "*.txt")
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct {
	d.filters = append(d.filters, FileFilter{DisplayName: displayName, Pattern: pattern})
	return d
}

func (d *SaveFileDialogStruct) AttachToWindow(window *WebviewWindow) *SaveFileDialogStruct {
	return d
}

// PromptForSingleSelection shows the save dialog and returns the chosen path.
// Always returns ("", nil) in the stub.
//
//	path, err := dialog.PromptForSingleSelection()
func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error) {
	return "", nil
}

// DialogManager exposes factory methods for all dialog types.
//
//	manager.Info().SetMessage("Saved!").Show()
//	path, _ := manager.OpenFile().PromptForSingleSelection()
type DialogManager struct{}

// OpenFile creates a file-open dialog builder.
func (dm *DialogManager) OpenFile() *OpenFileDialogStruct {
	return newOpenFileDialog()
}

// OpenFileWithOptions creates a file-open dialog builder pre-populated from options.
func (dm *DialogManager) OpenFileWithOptions(options *OpenFileDialogOptions) *OpenFileDialogStruct {
	result := newOpenFileDialog()
	result.SetOptions(options)
	return result
}

// SaveFile creates a file-save dialog builder.
func (dm *DialogManager) SaveFile() *SaveFileDialogStruct {
	return newSaveFileDialog()
}

// SaveFileWithOptions creates a file-save dialog builder pre-populated from options.
func (dm *DialogManager) SaveFileWithOptions(options *SaveFileDialogOptions) *SaveFileDialogStruct {
	result := newSaveFileDialog()
	result.SetOptions(options)
	return result
}

// Info creates an information message dialog.
//
//	manager.Info().SetMessage("Done").Show()
func (dm *DialogManager) Info() *MessageDialog {
	return newMessageDialog(InfoDialogType)
}

// Question creates a question message dialog.
//
//	manager.Question().SetMessage("Continue?").Show()
func (dm *DialogManager) Question() *MessageDialog {
	return newMessageDialog(QuestionDialogType)
}

// Warning creates a warning message dialog.
//
//	manager.Warning().SetMessage("Low disk space").Show()
func (dm *DialogManager) Warning() *MessageDialog {
	return newMessageDialog(WarningDialogType)
}

// Error creates an error message dialog.
//
//	manager.Error().SetMessage("Write failed").Show()
func (dm *DialogManager) Error() *MessageDialog {
	return newMessageDialog(ErrorDialogType)
}
