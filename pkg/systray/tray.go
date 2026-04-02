// pkg/systray/tray.go
package systray

import (
	_ "embed"
	"sync"

	"forge.lthn.ai/core/go/pkg/core"
)

//go:embed assets/apptray.png
var defaultIcon []byte

// Manager manages the system tray lifecycle.
//
// Example:
//
//	manager := systray.NewManager(platform)
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
//
// Example:
//
//	manager := systray.NewManager(platform)
func NewManager(platform Platform) *Manager {
	return &Manager{
		platform:  platform,
		callbacks: make(map[string]func()),
	}
}

// Setup creates the system tray with default icon and tooltip.
func (m *Manager) Setup(tooltip, label string) error {
	m.tray = m.platform.NewTray()
	if m.tray == nil {
		return core.E("systray.Setup", "platform returned nil tray", nil)
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
func (m *Manager) SetIcon(data []byte) error {
	if m.tray == nil {
		return core.E("systray.SetIcon", "tray not initialised", nil)
	}
	m.tray.SetIcon(data)
	m.hasIcon = len(data) > 0
	return nil
}

// SetTemplateIcon sets the template icon (macOS).
func (m *Manager) SetTemplateIcon(data []byte) error {
	if m.tray == nil {
		return core.E("systray.SetTemplateIcon", "tray not initialised", nil)
	}
	m.tray.SetTemplateIcon(data)
	m.hasTemplateIcon = len(data) > 0
	return nil
}

// SetTooltip sets the tray tooltip.
func (m *Manager) SetTooltip(text string) error {
	if m.tray == nil {
		return core.E("systray.SetTooltip", "tray not initialised", nil)
	}
	m.tray.SetTooltip(text)
	m.tooltip = text
	return nil
}

// SetLabel sets the tray label.
func (m *Manager) SetLabel(text string) error {
	if m.tray == nil {
		return core.E("systray.SetLabel", "tray not initialised", nil)
	}
	m.tray.SetLabel(text)
	m.label = text
	return nil
}

// AttachWindow attaches a panel window to the tray.
func (m *Manager) AttachWindow(w WindowHandle) error {
	if m.tray == nil {
		return core.E("systray.AttachWindow", "tray not initialised", nil)
	}
	m.mu.Lock()
	m.panelWindow = w
	m.mu.Unlock()
	m.tray.AttachWindow(w)
	return nil
}

// ShowPanel shows the attached tray panel window if one is configured.
func (m *Manager) ShowPanel() error {
	m.mu.RLock()
	w := m.panelWindow
	m.mu.RUnlock()
	if w == nil {
		return nil
	}
	w.Show()
	return nil
}

// HidePanel hides the attached tray panel window if one is configured.
func (m *Manager) HidePanel() error {
	m.mu.RLock()
	w := m.panelWindow
	m.mu.RUnlock()
	if w == nil {
		return nil
	}
	w.Hide()
	return nil
}

// Tray returns the underlying platform tray for direct access.
func (m *Manager) Tray() PlatformTray {
	return m.tray
}

// IsActive returns whether a tray has been created.
func (m *Manager) IsActive() bool {
	return m.tray != nil
}
