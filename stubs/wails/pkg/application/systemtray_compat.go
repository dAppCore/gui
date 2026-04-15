package application

import (
	"errors"
	"time"
)

// IconPosition matches the exported system tray icon placement enum from Wails.
type IconPosition int

const (
	NSImageNone IconPosition = iota
	NSImageOnly
	NSImageLeft
	NSImageRight
	NSImageBelow
	NSImageAbove
	NSImageOverlaps
	NSImageLeading
	NSImageTrailing
)

// WindowAttachConfig stores the attached window and basic positioning hints.
type WindowAttachConfig struct {
	Window   Window
	Offset   int
	Debounce time.Duration
}

// Label returns the current tray label.
func (s *SystemTray) Label() string {
	return s.label
}

// Run initialises default tray state for the stub.
func (s *SystemTray) Run() {
	ensureTrayCompatState(s)
}

// ToggleWindow toggles the visibility of the attached window.
func (s *SystemTray) ToggleWindow() {
	if s.attachedWindow == nil {
		return
	}
	if s.attachedWindow.IsVisible() {
		s.attachedWindow.Hide()
		return
	}
	state := ensureTrayCompatState(s)
	state.mu.RLock()
	offset := state.windowOffset
	state.mu.RUnlock()
	_ = s.PositionWindow(s.attachedWindow, offset)
	s.attachedWindow.Show()
	s.attachedWindow.Focus()
}

// ShowMenu invokes the tray's menu-open hook.
func (s *SystemTray) ShowMenu() {
	s.OpenMenu()
}

// ShowWindow shows and focuses the attached window.
func (s *SystemTray) ShowWindow() {
	if s.attachedWindow == nil {
		return
	}
	state := ensureTrayCompatState(s)
	state.mu.RLock()
	offset := state.windowOffset
	state.mu.RUnlock()
	_ = s.PositionWindow(s.attachedWindow, offset)
	s.attachedWindow.Show()
	s.attachedWindow.Focus()
}

// HideWindow hides the attached window.
func (s *SystemTray) HideWindow() {
	if s.attachedWindow != nil {
		s.attachedWindow.Hide()
	}
}

// PositionWindow positions the window near the tray entry in a deterministic stub-friendly way.
func (s *SystemTray) PositionWindow(window Window, offset int) error {
	if window == nil {
		return errors.New("system tray has no attached window")
	}
	x, y := window.Position()
	window.SetPosition(x+offset, y+offset)
	return nil
}

// SetDarkModeIcon stores an alternate dark-mode icon.
func (s *SystemTray) SetDarkModeIcon(icon []byte) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.darkModeIcon = append([]byte(nil), icon...)
	state.mu.Unlock()
	return s
}

// SetIconPosition stores the tray icon placement preference.
func (s *SystemTray) SetIconPosition(iconPosition IconPosition) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.iconPosition = iconPosition
	state.mu.Unlock()
	return s
}

// Destroy clears the tray state.
func (s *SystemTray) Destroy() {
	trayCompatStates.Delete(s)
}

// OnRightClick stores the right-click callback.
func (s *SystemTray) OnRightClick(handler func()) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.rightClickHandler = handler
	state.mu.Unlock()
	return s
}

// OnDoubleClick stores the double-click callback.
func (s *SystemTray) OnDoubleClick(handler func()) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.doubleClickHandler = handler
	state.mu.Unlock()
	return s
}

// OnRightDoubleClick stores the right-double-click callback.
func (s *SystemTray) OnRightDoubleClick(handler func()) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.rightDoubleClick = handler
	state.mu.Unlock()
	return s
}

// OnMouseEnter stores the mouse-enter callback.
func (s *SystemTray) OnMouseEnter(handler func()) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.mouseEnterHandler = handler
	state.mu.Unlock()
	return s
}

// OnMouseLeave stores the mouse-leave callback.
func (s *SystemTray) OnMouseLeave(handler func()) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.mouseLeaveHandler = handler
	state.mu.Unlock()
	return s
}

// Show marks the tray as visible.
func (s *SystemTray) Show() {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.visible = true
	state.mu.Unlock()
}

// Hide marks the tray as hidden.
func (s *SystemTray) Hide() {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.visible = false
	state.mu.Unlock()
}

// WindowOffset sets the offset used when positioning attached windows.
func (s *SystemTray) WindowOffset(offset int) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.windowOffset = offset
	state.mu.Unlock()
	return s
}

// WindowDebounce stores the debounce duration for window attach toggling.
func (s *SystemTray) WindowDebounce(debounce time.Duration) *SystemTray {
	state := ensureTrayCompatState(s)
	state.mu.Lock()
	state.windowDebounce = debounce
	state.mu.Unlock()
	return s
}

// OpenMenu triggers the menu-open callback when set.
func (s *SystemTray) OpenMenu() {
	state := ensureTrayCompatState(s)
	state.mu.RLock()
	handler := state.onMenuOpen
	state.mu.RUnlock()
	if handler != nil {
		handler()
	}
}
