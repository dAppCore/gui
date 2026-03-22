// pkg/systray/tray.go
package systray

import (
	_ "embed"
	"fmt"
	"sync"
)

//go:embed assets/apptray.png
var defaultIcon []byte

// Manager manages the system tray lifecycle.
// State that was previously in package-level vars is now on the Manager.
type Manager struct {
	platform  Platform
	tray      PlatformTray
	callbacks map[string]func()
	mu        sync.RWMutex
}

// NewManager creates a systray Manager.
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
		return fmt.Errorf("platform returned nil tray")
	}
	m.tray.SetTemplateIcon(defaultIcon)
	m.tray.SetTooltip(tooltip)
	m.tray.SetLabel(label)
	return nil
}

// SetIcon sets the tray icon.
func (m *Manager) SetIcon(data []byte) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetIcon(data)
	return nil
}

// SetTemplateIcon sets the template icon (macOS).
func (m *Manager) SetTemplateIcon(data []byte) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetTemplateIcon(data)
	return nil
}

// SetTooltip sets the tray tooltip.
func (m *Manager) SetTooltip(text string) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetTooltip(text)
	return nil
}

// SetLabel sets the tray label.
func (m *Manager) SetLabel(text string) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetLabel(text)
	return nil
}

// AttachWindow attaches a panel window to the tray.
func (m *Manager) AttachWindow(w WindowHandle) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.AttachWindow(w)
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
