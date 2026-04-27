package application

// Role identifies a platform-specific menu role.
// It is the same underlying type as MenuRole so existing constants (AppMenu,
// FileMenu, EditMenu, ViewMenu, WindowMenu, HelpMenu) are valid Role values.
//
//	quitItem := application.NewRole(application.Quit)
type Role = MenuRole

const (
	// NoRole indicates no special platform role.
	NoRole MenuRole = iota + 100

	// ServicesMenu is the macOS Services sub-menu.
	ServicesMenu

	// SpeechMenu is the macOS Speech sub-menu.
	SpeechMenu

	// Hide hides the current application.
	Hide

	// HideOthers hides all other applications.
	HideOthers

	// UnHide shows all hidden applications.
	UnHide

	// Front brings all windows to front.
	Front

	// Undo triggers the standard Undo action.
	Undo

	// Redo triggers the standard Redo action.
	Redo

	// Cut triggers the standard Cut action.
	Cut

	// Copy triggers the standard Copy action.
	Copy

	// Paste triggers the standard Paste action.
	Paste

	// PasteAndMatchStyle pastes without source formatting.
	PasteAndMatchStyle

	// SelectAll triggers the standard Select All action.
	SelectAll

	// Delete triggers the standard Delete action.
	Delete

	// Quit quits the application.
	Quit

	// CloseWindow closes the focused window.
	CloseWindow

	// About opens the About panel.
	About

	// Reload reloads the current webview.
	Reload

	// ForceReload force-reloads the current webview.
	ForceReload

	// ToggleFullscreen toggles fullscreen mode.
	ToggleFullscreen

	// OpenDevTools opens the developer tools panel.
	OpenDevTools

	// ResetZoom resets the webview zoom level.
	ResetZoom

	// ZoomIn increases the webview zoom level.
	ZoomIn

	// ZoomOut decreases the webview zoom level.
	ZoomOut

	// Minimise minimises the focused window.
	Minimise

	// Zoom zooms the focused window (macOS).
	Zoom

	// FullScreen enters fullscreen (macOS).
	FullScreen

	// Print opens the print dialog.
	Print

	// PageLayout opens the page layout dialog.
	PageLayout

	// ShowAll shows all windows.
	ShowAll

	// BringAllToFront brings all windows to front.
	BringAllToFront

	// NewFile triggers the New File action.
	NewFile

	// Open triggers the Open action.
	Open

	// Save triggers the Save action.
	Save

	// SaveAs triggers the Save As action.
	SaveAs

	// StartSpeaking starts text-to-speech on selected text.
	StartSpeaking

	// StopSpeaking stops text-to-speech.
	StopSpeaking

	// Revert triggers the Revert action.
	Revert

	// Find triggers the Find action.
	Find

	// FindAndReplace triggers the Find and Replace action.
	FindAndReplace

	// FindNext finds the next match.
	FindNext

	// FindPrevious finds the previous match.
	FindPrevious

	// Help opens the application help.
	Help
)

// NewMenuItem creates a new text menu item with the given label.
//
//	item := application.NewMenuItem("Open File").SetAccelerator("CmdOrCtrl+O")
func NewMenuItem(label string) *MenuItem {
	return &MenuItem{Label: label, Enabled: true}
}

// NewMenuItemSeparator creates a horizontal separator for use in menus.
//
//	menu.Items = append(menu.Items, application.NewMenuItemSeparator())
func NewMenuItemSeparator() *MenuItem {
	return &MenuItem{Label: "---"}
}

// NewMenuItemCheckbox creates a checkable menu item.
//
//	item := application.NewMenuItemCheckbox("Show Toolbar", true)
func NewMenuItemCheckbox(label string, checked bool) *MenuItem {
	return &MenuItem{Label: label, Checked: checked, Enabled: true}
}

// NewMenuItemRadio creates a radio-group menu item.
//
//	item := application.NewMenuItemRadio("Small", false)
func NewMenuItemRadio(label string, checked bool) *MenuItem {
	return &MenuItem{Label: label, Checked: checked, Enabled: true}
}

// NewSubMenuItem creates a menu item that opens a sub-menu.
//
//	sub := application.NewSubMenuItem("Recent Files")
func NewSubMenuItem(label string) *MenuItem {
	return &MenuItem{
		Label:   label,
		Enabled: true,
		submenu: &Menu{},
	}
}

// NewRole creates a menu item pre-configured for a platform role.
//
//	quitItem := application.NewRole(application.Quit)
func NewRole(role Role) *MenuItem {
	return &MenuItem{Label: roleLabel(role), Enabled: true}
}

// NewServicesMenu creates the macOS Services sub-menu item.
//
//	servicesMenu := application.NewServicesMenu()
func NewServicesMenu() *MenuItem {
	return NewSubMenuItem("Services")
}

// GetAccelerator returns the accelerator string currently set on the item.
//
//	accel := item.GetAccelerator() // e.g. "CmdOrCtrl+S"
func (mi *MenuItem) GetAccelerator() string {
	return mi.Accelerator
}

// GetSubmenu returns the submenu associated with the item, if any.
func (mi *MenuItem) GetSubmenu() *Menu {
	if mi == nil {
		return nil
	}
	return mi.submenu
}

// roleLabel maps a Role constant to a human-readable label.
func roleLabel(role Role) string {
	switch role {
	case AppMenu:
		return "App"
	case EditMenu:
		return "Edit"
	case FileMenu:
		return "File"
	case ViewMenu:
		return "View"
	case ServicesMenu:
		return "Services"
	case SpeechMenu:
		return "Speech"
	case WindowMenu:
		return "Window"
	case HelpMenu:
		return "Help"
	case Hide:
		return "Hide"
	case HideOthers:
		return "Hide Others"
	case UnHide:
		return "Show All"
	case Front:
		return "Bring All to Front"
	case Undo:
		return "Undo"
	case Redo:
		return "Redo"
	case Cut:
		return "Cut"
	case Copy:
		return "Copy"
	case Paste:
		return "Paste"
	case PasteAndMatchStyle:
		return "Paste and Match Style"
	case SelectAll:
		return "Select All"
	case Delete:
		return "Delete"
	case Quit:
		return "Quit"
	case CloseWindow:
		return "Close Window"
	case About:
		return "About"
	case Reload:
		return "Reload"
	case ForceReload:
		return "Force Reload"
	case ToggleFullscreen:
		return "Toggle Full Screen"
	case OpenDevTools:
		return "Open Developer Tools"
	case ResetZoom:
		return "Reset Zoom"
	case ZoomIn:
		return "Zoom In"
	case ZoomOut:
		return "Zoom Out"
	case Minimise:
		return "Minimise"
	case Zoom:
		return "Zoom"
	case FullScreen:
		return "Full Screen"
	case Print:
		return "Print"
	case PageLayout:
		return "Page Layout"
	case ShowAll:
		return "Show All"
	case BringAllToFront:
		return "Bring All to Front"
	case NewFile:
		return "New"
	case Open:
		return "Open"
	case Save:
		return "Save"
	case SaveAs:
		return "Save As"
	case StartSpeaking:
		return "Start Speaking"
	case StopSpeaking:
		return "Stop Speaking"
	case Revert:
		return "Revert"
	case Find:
		return "Find"
	case FindAndReplace:
		return "Find and Replace"
	case FindNext:
		return "Find Next"
	case FindPrevious:
		return "Find Previous"
	case Help:
		return "Help"
	default:
		return ""
	}
}
