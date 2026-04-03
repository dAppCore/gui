// pkg/menu/menu.go
package menu

// MenuItem describes a menu item for construction (structure only — no handlers).
// Use: item := menu.MenuItem{Label: "Quit", OnClick: func() {}}
type MenuItem struct {
	Label       string
	Accelerator string
	Type        string // "normal", "separator", "checkbox", "radio", "submenu"
	Checked     bool
	Disabled    bool
	Tooltip     string
	Children    []MenuItem
	Role        *MenuRole
	OnClick     func() // Injected by orchestrator, not by menu package consumer
}

// Manager builds application menus via a Platform backend.
// Use: mgr := menu.NewManager(platform)
type Manager struct {
	platform Platform
}

// NewManager creates a menu Manager.
// Use: mgr := menu.NewManager(platform)
func NewManager(platform Platform) *Manager {
	return &Manager{platform: platform}
}

// Build constructs a PlatformMenu from a tree of MenuItems.
// Use: built := mgr.Build([]menu.MenuItem{{Label: "File"}})
func (m *Manager) Build(items []MenuItem) PlatformMenu {
	menu := m.platform.NewMenu()
	m.buildItems(menu, items)
	return menu
}

func (m *Manager) buildItems(menu PlatformMenu, items []MenuItem) {
	for _, item := range items {
		if item.Role != nil {
			menu.AddRole(*item.Role)
			continue
		}
		if item.Type == "separator" {
			menu.AddSeparator()
			continue
		}
		if len(item.Children) > 0 {
			sub := menu.AddSubmenu(item.Label)
			m.buildItems(sub, item.Children)
			continue
		}
		mi := menu.Add(item.Label)
		if item.Accelerator != "" {
			mi.SetAccelerator(item.Accelerator)
		}
		if item.Tooltip != "" {
			mi.SetTooltip(item.Tooltip)
		}
		if item.OnClick != nil {
			mi.OnClick(item.OnClick)
		}
	}
}

// SetApplicationMenu builds and sets the application menu.
// Use: mgr.SetApplicationMenu([]menu.MenuItem{{Label: "Quit"}})
func (m *Manager) SetApplicationMenu(items []MenuItem) {
	menu := m.Build(items)
	m.platform.SetApplicationMenu(menu)
}

// Platform returns the underlying platform.
// Use: backend := mgr.Platform()
func (m *Manager) Platform() Platform {
	return m.platform
}
