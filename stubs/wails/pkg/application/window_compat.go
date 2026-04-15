package application

import "sync"

// Enabled/Disabled map to the tristate values used by the stub macOS webview preferences.
var Enabled = 1
var Disabled = 2

// WindowEventListener mirrors the Wails listener handle type.
type WindowEventListener struct {
	callback func(event *WindowEvent)
}

var windowEventCancellation sync.Map

// NewWindowEvent constructs a WindowEvent with an empty context.
func NewWindowEvent() *WindowEvent {
	return &WindowEvent{}
}

// IsCancelled reports whether the event was cancelled.
func (w *WindowEvent) IsCancelled() bool {
	cancelled, _ := windowEventCancellation.Load(w)
	value, _ := cancelled.(bool)
	return value
}

// Cancel marks the event as cancelled.
func (w *WindowEvent) Cancel() {
	windowEventCancellation.Store(w, true)
}

// NewWindow creates a new in-memory webview window.
func NewWindow(options WebviewWindowOptions) *WebviewWindow {
	return newWebviewWindow(options)
}

// PhysicalBounds returns the stored physical bounds when available.
func (w *WebviewWindow) PhysicalBounds() Rect {
	state := ensureWebviewCompatState(w)
	state.mu.RLock()
	defer state.mu.RUnlock()
	if state.physicalBoundsOK {
		return state.physicalBounds
	}
	return w.Bounds()
}

// SetPhysicalBounds stores the physical bounds of the window.
func (w *WebviewWindow) SetPhysicalBounds(physicalBounds Rect) {
	state := ensureWebviewCompatState(w)
	state.mu.Lock()
	state.physicalBounds = physicalBounds
	state.physicalBoundsOK = true
	state.mu.Unlock()
}

// RegisterKeyBinding registers a window-scoped keybinding callback.
func (w *WebviewWindow) RegisterKeyBinding(binding string, callback func(window Window)) *WebviewWindow {
	state := ensureWebviewCompatState(w)
	state.mu.Lock()
	state.keyBindings[binding] = callback
	state.mu.Unlock()
	return w
}

// HandleDragEnter is a no-op in the in-memory stub.
func (w *WebviewWindow) HandleDragEnter() {}

// HandleDragOver is a no-op in the in-memory stub.
func (w *WebviewWindow) HandleDragOver(x int, y int) {}

// HandleDragLeave is a no-op in the in-memory stub.
func (w *WebviewWindow) HandleDragLeave() {}

// CloseDevTools closes the devtools pane in the stub.
func (w *WebviewWindow) CloseDevTools() {
	state := ensureWebviewCompatState(w)
	state.mu.Lock()
	state.devtoolsOpen = false
	state.mu.Unlock()
}

// DevToolsOpen reports whether devtools are currently marked open.
func (w *WebviewWindow) DevToolsOpen() bool {
	state := ensureWebviewCompatState(w)
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.devtoolsOpen
}
