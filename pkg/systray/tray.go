// pkg/systray/tray.go
package systray

import (
	_ "embed"
	"sync"

	coreerr "forge.lthn.ai/core/go-log"
)

//go:embed assets/apptray.png
var defaultIcon []byte

// Manager manages the system tray lifecycle.
// Use: manager := systray.NewManager(platform)
type Manager struct {
	platform        Platform
	tray            PlatformTray
	panelWindow     WindowHandle
	callbacks       map[string]func()
	tooltip         string
	label           string
	hasIcon         bool
	hasTemplateIcon bool
	menuItems       []TrayMenuItem
	mu              sync.RWMutex
}

// NewManager creates a systray Manager.
// systray.NewManager(systray.NewWailsPlatform(app)).Setup("Core", "Core")
func NewManager(platform Platform) *Manager {
	return &Manager{
		platform:  platform,
		callbacks: make(map[string]func()),
	}
}

// Setup creates the system tray with default icon and tooltip.
// systray.NewManager(systray.NewWailsPlatform(app)).Setup("Core", "Core")
func (m *Manager) Setup(tooltip, label string) error {
	m.tray = m.platform.NewTray()
	if m.tray == nil {
		return coreerr.E("systray.Setup", "platform returned nil tray", nil)
	}
	m.tray.SetTemplateIcon(defaultIcon)
	m.tray.SetTooltip(tooltip)
	m.tray.SetLabel(label)
	m.tooltip = tooltip
	m.label = label
	m.hasTemplateIcon = true
	return nil
}

// SetIcon sets the tray icon.
// Use: _ = manager.SetIcon(iconBytes)
func (m *Manager) SetIcon(data []byte) error {
	if m.tray == nil {
		return coreerr.E("systray.SetIcon", "tray not initialised", nil)
	}
	m.tray.SetIcon(data)
	m.hasIcon = len(data) > 0
	return nil
}

// SetTemplateIcon sets the template icon (macOS).
// Use: _ = manager.SetTemplateIcon(iconBytes)
func (m *Manager) SetTemplateIcon(data []byte) error {
	if m.tray == nil {
		return coreerr.E("systray.SetTemplateIcon", "tray not initialised", nil)
	}
	m.tray.SetTemplateIcon(data)
	m.hasTemplateIcon = len(data) > 0
	return nil
}

// SetTooltip sets the tray tooltip.
// Use: _ = manager.SetTooltip("Core is ready")
func (m *Manager) SetTooltip(text string) error {
	if m.tray == nil {
		return coreerr.E("systray.SetTooltip", "tray not initialised", nil)
	}
	m.tray.SetTooltip(text)
	m.tooltip = text
	return nil
}

// SetLabel sets the tray label.
// Use: _ = manager.SetLabel("Core")
func (m *Manager) SetLabel(text string) error {
	if m.tray == nil {
		return coreerr.E("systray.SetLabel", "tray not initialised", nil)
	}
	m.tray.SetLabel(text)
	m.label = text
	return nil
}

// AttachWindow attaches a panel window to the tray.
// Use: _ = manager.AttachWindow(windowHandle)
func (m *Manager) AttachWindow(w WindowHandle) error {
	if m.tray == nil {
		return coreerr.E("systray.AttachWindow", "tray not initialised", nil)
	}
	m.mu.Lock()
	m.panelWindow = w
	m.mu.Unlock()
	m.tray.AttachWindow(w)
	return nil
}

// ShowMessage displays a tray message if the backend supports it.
func (m *Manager) ShowMessage(title, message string) error {
	if m.tray == nil {
		return coreerr.E("systray.ShowMessage", "tray not initialised", nil)
	}
	return m.tray.ShowMessage(title, message)
}

// ShowPanel reveals the attached tray panel window.
func (m *Manager) ShowPanel() error {
	m.mu.RLock()
	panel := m.panelWindow
	m.mu.RUnlock()
	if panel == nil {
		return coreerr.E("systray.ShowPanel", "panel window not attached", nil)
	}
	panel.Show()
	return nil
}

// HidePanel hides the attached tray panel window.
func (m *Manager) HidePanel() error {
	m.mu.RLock()
	panel := m.panelWindow
	m.mu.RUnlock()
	if panel == nil {
		return coreerr.E("systray.HidePanel", "panel window not attached", nil)
	}
	panel.Hide()
	return nil
}

// Tray returns the underlying platform tray for direct access.
// Use: tray := manager.Tray()
func (m *Manager) Tray() PlatformTray {
	return m.tray
}

// IsActive returns whether a tray has been created.
// Use: active := manager.IsActive()
func (m *Manager) IsActive() bool {
	return m.tray != nil
}
