package application

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// Window is the interface satisfied by all window types in the application.
// Fluent mutating methods return Window so callers can chain calls:
//
//	app.Window.NewWithOptions(opts).SetTitle("Main").Show()
type Window interface {
	// Identity
	ID() uint
	Name() string

	// Visibility
	Show() Window
	Hide() Window
	IsVisible() bool

	// Lifecycle
	Close()
	Focus()
	Run()

	// Geometry
	Center()
	Position() (x int, y int)
	RelativePosition() (x int, y int)
	Size() (width int, height int)
	Width() int
	Height() int
	Bounds() Rect
	SetPosition(x, y int)
	SetRelativePosition(x, y int) Window
	SetSize(width, height int) Window
	SetBounds(bounds Rect)
	SetMaxSize(maxWidth, maxHeight int) Window
	SetMinSize(minWidth, minHeight int) Window
	EnableSizeConstraints()
	DisableSizeConstraints()
	Resizable() bool
	SetResizable(b bool) Window

	// State
	Maximise() Window
	UnMaximise()
	ToggleMaximise()
	IsMaximised() bool
	Minimise() Window
	UnMinimise()
	IsMinimised() bool
	Fullscreen() Window
	UnFullscreen()
	ToggleFullscreen()
	IsFullscreen() bool
	Restore()
	SnapAssist()

	// Title and content
	SetTitle(title string) Window
	SetURL(s string) Window
	SetHTML(html string) Window

	// Titlebar buttons (macOS / Windows)
	SetMinimiseButtonState(state ButtonState) Window
	SetMaximiseButtonState(state ButtonState) Window
	SetCloseButtonState(state ButtonState) Window

	// Menu bar
	SetMenu(menu *Menu)
	ShowMenuBar()
	HideMenuBar()
	ToggleMenuBar()

	// Appearance
	SetBackgroundColour(colour RGBA) Window
	SetAlwaysOnTop(b bool) Window
	SetFrameless(frameless bool) Window
	ToggleFrameless()
	SetIgnoreMouseEvents(ignore bool) Window
	IsIgnoreMouseEvents() bool
	SetContentProtection(protection bool) Window

	// Zoom
	GetZoom() float64
	SetZoom(magnification float64) Window
	Zoom()
	ZoomIn()
	ZoomOut()
	ZoomReset() Window

	// Border sizes (Windows)
	GetBorderSizes() *LRTB

	// Screen
	GetScreen() (*Screen, error)

	// JavaScript / events
	ExecJS(js string)
	EmitEvent(name string, data ...any) bool
	DispatchWailsEvent(event *CustomEvent)
	OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func()
	RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func()

	// Drag-and-drop (internal message bus)
	handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails)
	InitiateFrontendDropProcessing(filenames []string, x int, y int)

	// Message handling (internal)
	HandleMessage(message string)
	HandleWindowEvent(id uint)
	HandleKeyEvent(acceleratorString string)

	// Context menu
	OpenContextMenu(data *ContextMenuData)

	// Modal
	AttachModal(modalWindow Window)

	// DevTools
	OpenDevTools()

	// Print
	Print() error

	// Flash (Windows taskbar flash)
	Flash(enabled bool)

	// Focus tracking
	IsFocused() bool

	// Native handle (platform-specific, use with care)
	NativeWindow() unsafe.Pointer

	// Enabled state
	SetEnabled(enabled bool)

	// Reload
	Reload()
	ForceReload()

	// Logging (routed to the application logger)
	Info(message string, args ...any)
	Error(message string, args ...any)

	// Internal hooks
	shouldUnconditionallyClose() bool

	// Editing operations (routed to focused element)
	cut()
	copy()
	paste()
	undo()
	redo()
	delete()
	selectAll()
}
