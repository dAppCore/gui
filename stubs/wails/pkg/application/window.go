package application

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// Window is the interface satisfied by both WebviewWindow and BrowserWindow.
// All methods mirror the Wails v3 Window interface exactly, including unexported
// platform hooks required for internal dispatch.
//
//	var w Window = app.Window.NewWithOptions(opts)
//	w.SetTitle("My App").Show()
type Window interface {
	// Identity

	// w.ID() → 1
	ID() uint
	// w.Name() → "main"
	Name() string

	// Lifecycle

	// w.Show().Focus()
	Show() Window
	// w.Hide()
	Hide() Window
	// w.Close()
	Close()
	// w.Focus()
	Focus()
	// w.Run()
	Run()
	// w.Restore()
	Restore()

	// Geometry

	// x, y := w.Position()
	Position() (int, int)
	// x, y := w.RelativePosition()
	RelativePosition() (int, int)
	// w, h := w.Size()
	Size() (width int, height int)
	// px := w.Width()
	Width() int
	// px := w.Height()
	Height() int
	// r := w.Bounds()
	Bounds() Rect
	// w.SetBounds(Rect{X: 0, Y: 0, Width: 1280, Height: 800})
	SetBounds(bounds Rect)
	// w.SetPosition(100, 200)
	SetPosition(x, y int)
	// w.SetRelativePosition(0, 0)
	SetRelativePosition(x, y int) Window
	// w.SetSize(1280, 800)
	SetSize(width, height int) Window
	// w.SetMinSize(640, 480)
	SetMinSize(minWidth, minHeight int) Window
	// w.SetMaxSize(3840, 2160)
	SetMaxSize(maxWidth, maxHeight int) Window
	// w.Center()
	Center()

	// Title and content

	// t := w.Title()  — not in interface, provided by concrete types
	// w.SetTitle("My App")
	SetTitle(title string) Window
	// w.SetURL("https://example.com")
	SetURL(url string) Window
	// w.SetHTML("<h1>Hello</h1>")
	SetHTML(html string) Window

	// Visibility states

	// w.IsVisible()
	IsVisible() bool
	// w.IsFullscreen()
	IsFullscreen() bool
	// w.IsMaximised()
	IsMaximised() bool
	// w.IsMinimised()
	IsMinimised() bool
	// w.IsFocused()
	IsFocused() bool
	// w.IsIgnoreMouseEvents()
	IsIgnoreMouseEvents() bool
	// w.Resizable()
	Resizable() bool

	// Window state transitions

	// w.Fullscreen()
	Fullscreen() Window
	// w.UnFullscreen()
	UnFullscreen()
	// w.Maximise()
	Maximise() Window
	// w.UnMaximise()
	UnMaximise()
	// w.Minimise()
	Minimise() Window
	// w.UnMinimise()
	UnMinimise()
	// w.ToggleFullscreen()
	ToggleFullscreen()
	// w.ToggleMaximise()
	ToggleMaximise()
	// w.SnapAssist()
	SnapAssist()

	// Style

	// w.SetAlwaysOnTop(true)
	SetAlwaysOnTop(b bool) Window
	// w.SetBackgroundColour(application.NewRGBA(0, 0, 0, 255))
	SetBackgroundColour(colour RGBA) Window
	// w.SetFrameless(true)
	SetFrameless(frameless bool) Window
	// w.SetResizable(false)
	SetResizable(b bool) Window
	// w.SetIgnoreMouseEvents(true)
	SetIgnoreMouseEvents(ignore bool) Window
	// w.SetMinimiseButtonState(ButtonHidden)
	SetMinimiseButtonState(state ButtonState) Window
	// w.SetMaximiseButtonState(ButtonDisabled)
	SetMaximiseButtonState(state ButtonState) Window
	// w.SetCloseButtonState(ButtonEnabled)
	SetCloseButtonState(state ButtonState) Window
	// w.SetEnabled(false)
	SetEnabled(enabled bool)
	// w.SetContentProtection(true)
	SetContentProtection(protection bool) Window

	// Menu

	// w.SetMenu(myMenu)
	SetMenu(menu *Menu)
	// w.ShowMenuBar()
	ShowMenuBar()
	// w.HideMenuBar()
	HideMenuBar()
	// w.ToggleMenuBar()
	ToggleMenuBar()
	// w.ToggleFrameless()
	ToggleFrameless()

	// WebView

	// w.ExecJS("document.title = 'Hello'")
	ExecJS(js string)
	// w.Reload()
	Reload()
	// w.ForceReload()
	ForceReload()
	// w.OpenDevTools()
	OpenDevTools()
	// w.OpenContextMenu(&ContextMenuData{Name: "main"})
	OpenContextMenu(data *ContextMenuData)

	// Zoom

	// w.Zoom()
	Zoom()
	// w.ZoomIn()
	ZoomIn()
	// w.ZoomOut()
	ZoomOut()
	// w.ZoomReset()
	ZoomReset() Window
	// z := w.GetZoom()
	GetZoom() float64
	// w.SetZoom(1.5)
	SetZoom(magnification float64) Window

	// Events

	// cancel := w.OnWindowEvent(events.Common.WindowFocus, func(e *WindowEvent) { ... })
	OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func()
	// cancel := w.RegisterHook(events.Common.WindowClose, func(e *WindowEvent) { ... })
	RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func()
	// w.EmitEvent("user:login", payload)
	EmitEvent(name string, data ...any) bool
	// w.DispatchWailsEvent(&CustomEvent{Name: "init"})
	DispatchWailsEvent(event *CustomEvent)

	// Screen and display

	// screen, err := w.GetScreen()
	GetScreen() (*Screen, error)
	// borders := w.GetBorderSizes()
	GetBorderSizes() *LRTB

	// Constraints

	// w.EnableSizeConstraints()
	EnableSizeConstraints()
	// w.DisableSizeConstraints()
	DisableSizeConstraints()

	// Modal

	// w.AttachModal(modalWindow)
	AttachModal(modalWindow Window)

	// Utilities

	// w.Flash(true)
	Flash(enabled bool)
	// err := w.Print()
	Print() error
	// w.Error("something went wrong: %s", details)
	Error(message string, args ...any)
	// w.Info("window ready")
	Info(message string, args ...any)

	// Platform

	// ptr := w.NativeWindow()
	NativeWindow() unsafe.Pointer

	// Internal platform hooks — implemented by concrete window types.

	handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails)
	InitiateFrontendDropProcessing(filenames []string, x int, y int)
	HandleMessage(message string)
	HandleWindowEvent(id uint)
	HandleKeyEvent(acceleratorString string)
	shouldUnconditionallyClose() bool
	cut()
	copy()
	paste()
	undo()
	redo()
	delete()
	selectAll()
}
