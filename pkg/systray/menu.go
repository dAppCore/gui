// pkg/systray/menu.go
package systray

import "forge.lthn.ai/core/go/pkg/core"

// SetMenu sets a dynamic menu on the tray from TrayMenuItem descriptors.
func (m *Manager) SetMenu(items []TrayMenuItem) error {
	if m.tray == nil {
		return core.E("systray.SetMenu", "tray not initialised", nil)
	}
	m.menuItems = append([]TrayMenuItem(nil), items...)
	menu := m.buildMenu(items)
	m.tray.SetMenu(menu)
	return nil
}

// buildMenu recursively builds a PlatformMenu from TrayMenuItem descriptors.
func (m *Manager) buildMenu(items []TrayMenuItem) PlatformMenu {
	menu := m.platform.NewMenu()
	m.buildMenuInto(menu, items)
	return menu
}

func (m *Manager) buildMenuInto(menu PlatformMenu, items []TrayMenuItem) {
	for _, item := range items {
		if item.Type == "separator" {
			menu.AddSeparator()
			continue
		}
		if len(item.Submenu) > 0 {
			sub := menu.AddSubmenu(item.Label)
			m.buildMenuInto(sub, item.Submenu)
			continue
		}
		mi := menu.Add(item.Label)
		if item.Tooltip != "" {
			mi.SetTooltip(item.Tooltip)
		}
		if item.Disabled {
			mi.SetEnabled(false)
		}
		if item.Checked {
			mi.SetChecked(true)
		}
		if item.ActionID != "" {
			actionID := item.ActionID
			mi.OnClick(func() {
				if cb, ok := m.GetCallback(actionID); ok {
					cb()
				}
			})
		}
	}
}

// RegisterCallback registers a callback for a menu action ID.
func (m *Manager) RegisterCallback(actionID string, callback func()) {
	m.mu.Lock()
	m.callbacks[actionID] = callback
	m.mu.Unlock()
}

// UnregisterCallback removes a callback.
func (m *Manager) UnregisterCallback(actionID string) {
	m.mu.Lock()
	delete(m.callbacks, actionID)
	m.mu.Unlock()
}

// GetCallback returns the callback for an action ID.
func (m *Manager) GetCallback(actionID string) (func(), bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cb, ok := m.callbacks[actionID]
	return cb, ok
}

// GetInfo returns tray status information.
func (m *Manager) GetInfo() map[string]any {
	return map[string]any{
		"active":          m.IsActive(),
		"tooltip":         m.tooltip,
		"label":           m.label,
		"hasIcon":         m.hasIcon,
		"hasTemplateIcon": m.hasTemplateIcon,
		"menuItems":       append([]TrayMenuItem(nil), m.menuItems...),
	}
}
