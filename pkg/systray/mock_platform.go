package systray

// MockPlatform is an exported mock for cross-package integration tests.
type MockPlatform struct{}

func NewMockPlatform() *MockPlatform { return &MockPlatform{} }

func (m *MockPlatform) NewTray() PlatformTray { return &exportedMockTray{} }
func (m *MockPlatform) NewMenu() PlatformMenu { return &exportedMockMenu{} }

type exportedMockTray struct {
	icon, templateIcon []byte
	tooltip, label     string
}

func (t *exportedMockTray) SetIcon(data []byte)         { t.icon = data }
func (t *exportedMockTray) SetTemplateIcon(data []byte) { t.templateIcon = data }
func (t *exportedMockTray) SetTooltip(text string)      { t.tooltip = text }
func (t *exportedMockTray) SetLabel(text string)        { t.label = text }
func (t *exportedMockTray) SetMenu(menu PlatformMenu)   {}
func (t *exportedMockTray) AttachWindow(w WindowHandle) {}

type exportedMockMenu struct {
	items    []exportedMockMenuItem
	submenus []*exportedMockMenu
}

func (m *exportedMockMenu) Add(label string) PlatformMenuItem {
	mi := &exportedMockMenuItem{label: label}
	m.items = append(m.items, *mi)
	return mi
}
func (m *exportedMockMenu) AddSeparator() {}
func (m *exportedMockMenu) AddSubmenu(label string) PlatformMenu {
	sub := &exportedMockMenu{}
	m.items = append(m.items, exportedMockMenuItem{label: label})
	m.submenus = append(m.submenus, sub)
	return sub
}

type exportedMockMenuItem struct {
	label, tooltip   string
	checked, enabled bool
	onClick          func()
}

func (mi *exportedMockMenuItem) SetTooltip(tip string)    { mi.tooltip = tip }
func (mi *exportedMockMenuItem) SetChecked(checked bool)  { mi.checked = checked }
func (mi *exportedMockMenuItem) SetEnabled(enabled bool)  { mi.enabled = enabled }
func (mi *exportedMockMenuItem) OnClick(fn func())        { mi.onClick = fn }
func (mi *exportedMockMenuItem) AddSubmenu() PlatformMenu { return &exportedMockMenu{} }
