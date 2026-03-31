package application

import "sync/atomic"

// Role identifies a platform-defined menu role.
//
//	item := NewRole(Quit)   // produces the platform quit item
type Role uint

const (
	NoRole       Role = iota
	AppMenu      Role = iota
	EditMenu      Role = iota
	ViewMenu      Role = iota
	WindowMenu    Role = iota
	ServicesMenu  Role = iota
	HelpMenu      Role = iota
	SpeechMenu    Role = iota
	FileMenu      Role = iota

	Hide               Role = iota
	HideOthers         Role = iota
	ShowAll            Role = iota
	BringAllToFront    Role = iota
	UnHide             Role = iota
	About              Role = iota
	Undo               Role = iota
	Redo               Role = iota
	Cut                Role = iota
	Copy               Role = iota
	Paste              Role = iota
	PasteAndMatchStyle Role = iota
	SelectAll          Role = iota
	Delete             Role = iota
	Quit               Role = iota
	CloseWindow        Role = iota
	Reload             Role = iota
	ForceReload        Role = iota
	OpenDevTools       Role = iota
	ResetZoom          Role = iota
	ZoomIn             Role = iota
	ZoomOut            Role = iota
	ToggleFullscreen   Role = iota
	Minimise           Role = iota
	Zoom               Role = iota
	FullScreen         Role = iota
	NewFile            Role = iota
	Open               Role = iota
	Save               Role = iota
	SaveAs             Role = iota
	StartSpeaking      Role = iota
	StopSpeaking       Role = iota
	Revert             Role = iota
	Print              Role = iota
	PageLayout         Role = iota
	Find               Role = iota
	FindAndReplace     Role = iota
	FindNext           Role = iota
	FindPrevious       Role = iota
	Front              Role = iota
	Help               Role = iota
)

// menuItemType classifies what kind of item a MenuItem is.
type menuItemType int

const (
	menuItemTypeText      menuItemType = iota
	menuItemTypeSeparator menuItemType = iota
	menuItemTypeCheckbox  menuItemType = iota
	menuItemTypeRadio     menuItemType = iota
	menuItemTypeSubmenu   menuItemType = iota
)

var globalMenuItemID uintptr

// MenuItem is a node in a menu tree.
//
//	item := NewMenuItem("Preferences").
//	    SetAccelerator("CmdOrCtrl+,").
//	    OnClick(func(ctx *Context) { openPrefs() })
type MenuItem struct {
	id              uint
	label           string
	tooltip         string
	disabled        bool
	checked         bool
	hidden          bool
	bitmap          []byte
	submenu         *Menu
	callback        func(*Context)
	itemType        menuItemType
	acceleratorStr  string
	role            Role
	contextMenuData *ContextMenuData

	radioGroupMembers []*MenuItem
}

func nextMenuItemID() uint {
	return uint(atomic.AddUintptr(&globalMenuItemID, 1))
}

// NewMenuItem creates a standard clickable menu item with the given label.
//
//	item := NewMenuItem("Save").OnClick(func(ctx *Context) { save() })
func NewMenuItem(label string) *MenuItem {
	return &MenuItem{
		id:       nextMenuItemID(),
		label:    label,
		itemType: menuItemTypeText,
		disabled: false,
	}
}

// NewMenuItemSeparator creates a horizontal separator.
//
//	menu.AppendItem(NewMenuItemSeparator())
func NewMenuItemSeparator() *MenuItem {
	return &MenuItem{
		id:       nextMenuItemID(),
		itemType: menuItemTypeSeparator,
	}
}

// NewMenuItemCheckbox creates a checkable menu item.
//
//	item := NewMenuItemCheckbox("Show Toolbar", true)
func NewMenuItemCheckbox(label string, checked bool) *MenuItem {
	return &MenuItem{
		id:       nextMenuItemID(),
		label:    label,
		checked:  checked,
		itemType: menuItemTypeCheckbox,
	}
}

// NewMenuItemRadio creates a radio-group menu item.
//
//	light := NewMenuItemRadio("Light Theme", true)
func NewMenuItemRadio(label string, checked bool) *MenuItem {
	return &MenuItem{
		id:       nextMenuItemID(),
		label:    label,
		checked:  checked,
		itemType: menuItemTypeRadio,
	}
}

// NewSubMenuItem creates an item that reveals a child menu on hover.
//
//	sub := NewSubMenuItem("Recent Files")
//	sub.GetSubmenu().Add("report.pdf")
func NewSubMenuItem(label string) *MenuItem {
	return &MenuItem{
		id:       nextMenuItemID(),
		label:    label,
		itemType: menuItemTypeSubmenu,
		submenu:  &Menu{label: label},
	}
}

// NewRole creates a platform-managed menu item for the given role.
//
//	menu.AppendItem(NewRole(Quit))
func NewRole(role Role) *MenuItem {
	item := &MenuItem{
		id:       nextMenuItemID(),
		label:    roleLabel(role),
		itemType: menuItemTypeText,
		role:     role,
	}
	return item
}

func roleLabel(role Role) string {
	switch role {
	case AppMenu:
		return "Application"
	case EditMenu:
		return "Edit"
	case ViewMenu:
		return "View"
	case WindowMenu:
		return "Window"
	case ServicesMenu:
		return "Services"
	case HelpMenu:
		return "Help"
	case SpeechMenu:
		return "Speech"
	case FileMenu:
		return "File"
	case Hide:
		return "Hide"
	case HideOthers:
		return "Hide Others"
	case ShowAll:
		return "Show All"
	case BringAllToFront:
		return "Bring All to Front"
	case UnHide:
		return "Unhide"
	case About:
		return "About"
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
	case Reload:
		return "Reload"
	case ForceReload:
		return "Force Reload"
	case OpenDevTools:
		return "Open Dev Tools"
	case ResetZoom:
		return "Reset Zoom"
	case ZoomIn:
		return "Zoom In"
	case ZoomOut:
		return "Zoom Out"
	case ToggleFullscreen:
		return "Toggle Fullscreen"
	case Minimise:
		return "Minimise"
	case Zoom:
		return "Zoom"
	case FullScreen:
		return "Fullscreen"
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
	case Print:
		return "Print"
	case PageLayout:
		return "Page Layout"
	case Find:
		return "Find"
	case FindAndReplace:
		return "Find and Replace"
	case FindNext:
		return "Find Next"
	case FindPrevious:
		return "Find Previous"
	case Front:
		return "Bring All to Front"
	case Help:
		return "Help"
	default:
		return ""
	}
}

// Fluent setters — all return *MenuItem for chaining.

// SetLabel updates the visible label.
func (m *MenuItem) SetLabel(s string) *MenuItem { m.label = s; return m }

// SetAccelerator sets the keyboard shortcut string (e.g. "CmdOrCtrl+S").
func (m *MenuItem) SetAccelerator(shortcut string) *MenuItem {
	m.acceleratorStr = shortcut
	return m
}

// GetAccelerator returns the raw accelerator string.
func (m *MenuItem) GetAccelerator() string { return m.acceleratorStr }

// RemoveAccelerator clears the keyboard shortcut.
func (m *MenuItem) RemoveAccelerator() { m.acceleratorStr = "" }

// SetTooltip sets the hover tooltip.
func (m *MenuItem) SetTooltip(s string) *MenuItem { m.tooltip = s; return m }

// SetRole assigns a platform role to the item.
func (m *MenuItem) SetRole(role Role) *MenuItem { m.role = role; return m }

// SetEnabled enables or disables the item.
func (m *MenuItem) SetEnabled(enabled bool) *MenuItem { m.disabled = !enabled; return m }

// SetBitmap sets a raw-bytes icon for the item (Windows).
func (m *MenuItem) SetBitmap(bitmap []byte) *MenuItem {
	m.bitmap = append([]byte(nil), bitmap...)
	return m
}

// SetChecked sets the checked state for checkbox/radio items.
func (m *MenuItem) SetChecked(checked bool) *MenuItem { m.checked = checked; return m }

// SetHidden hides or shows the item without removing it.
func (m *MenuItem) SetHidden(hidden bool) *MenuItem { m.hidden = hidden; return m }

// OnClick registers the callback invoked when the item is activated.
func (m *MenuItem) OnClick(f func(*Context)) *MenuItem { m.callback = f; return m }

// Getters

// Label returns the item's display label.
func (m *MenuItem) Label() string { return m.label }

// Tooltip returns the hover tooltip.
func (m *MenuItem) Tooltip() string { return m.tooltip }

// Enabled reports whether the item is interactive.
func (m *MenuItem) Enabled() bool { return !m.disabled }

// Checked reports whether the item is checked.
func (m *MenuItem) Checked() bool { return m.checked }

// Hidden reports whether the item is hidden.
func (m *MenuItem) Hidden() bool { return m.hidden }

// IsSeparator reports whether this is a separator item.
func (m *MenuItem) IsSeparator() bool { return m.itemType == menuItemTypeSeparator }

// IsSubmenu reports whether this item contains a child menu.
func (m *MenuItem) IsSubmenu() bool { return m.itemType == menuItemTypeSubmenu }

// IsCheckbox reports whether this is a checkbox item.
func (m *MenuItem) IsCheckbox() bool { return m.itemType == menuItemTypeCheckbox }

// IsRadio reports whether this is a radio item.
func (m *MenuItem) IsRadio() bool { return m.itemType == menuItemTypeRadio }

// GetSubmenu returns the submenu, or nil if this is not a submenu item.
func (m *MenuItem) GetSubmenu() *Menu { return m.submenu }

// Clone returns a shallow copy of the MenuItem.
func (m *MenuItem) Clone() *MenuItem {
	cloned := *m
	if m.submenu != nil {
		cloned.submenu = m.submenu.Clone()
	}
	if m.bitmap != nil {
		cloned.bitmap = append([]byte(nil), m.bitmap...)
	}
	if m.contextMenuData != nil {
		cloned.contextMenuData = m.contextMenuData.clone()
	}
	cloned.radioGroupMembers = nil
	return &cloned
}

// Destroy frees resources held by the item and its submenu.
func (m *MenuItem) Destroy() {
	if m.submenu != nil {
		m.submenu.Destroy()
		m.submenu = nil
	}
	m.callback = nil
	m.radioGroupMembers = nil
	m.contextMenuData = nil
}
