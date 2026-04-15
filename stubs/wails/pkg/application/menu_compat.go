package application

// NewContextMenu creates and registers a named context menu.
func NewContextMenu(name string) *ContextMenu {
	result := &ContextMenu{Menu: NewMenu(), name: name}
	result.Update()
	return result
}

// Update refreshes the menu and re-registers it with the app-level manager.
func (m *ContextMenu) Update() {
	if m == nil || m.Menu == nil {
		return
	}
	m.Menu.Update()
	if globalApplication != nil {
		globalApplication.ContextMenu.Add(m.name, m)
	}
}

// Destroy unregisters the context menu.
func (m *ContextMenu) Destroy() {
	if globalApplication != nil {
		globalApplication.ContextMenu.Remove(m.name)
	}
}

// AddCheckbox appends a checkbox item to the menu.
func (m *Menu) AddCheckbox(label string, enabled bool) *MenuItem {
	item := NewMenuItemCheckbox(label, enabled)
	m.Items = append(m.Items, item)
	return item
}

// AddRadio appends a radio item to the menu.
func (m *Menu) AddRadio(label string, enabled bool) *MenuItem {
	item := NewMenuItemRadio(label, enabled)
	m.Items = append(m.Items, item)
	return item
}

// Update normalises menu radio groups.
func (m *Menu) Update() {
	var radioGroup []*MenuItem
	flush := func() {
		if len(radioGroup) == 0 {
			return
		}
		for _, item := range radioGroup {
			item.radioGroupMembers = radioGroup
		}
		radioGroup = nil
	}
	for _, item := range m.Items {
		if item == nil {
			continue
		}
		if item.itemType != menuItemTypeRadio {
			flush()
		}
		if item.itemType == menuItemTypeSubmenu && item.submenu != nil {
			item.submenu.Update()
			continue
		}
		if item.itemType == menuItemTypeRadio {
			radioGroup = append(radioGroup, item)
		}
	}
	flush()
}

// Clear removes all menu items.
func (m *Menu) Clear() {
	for _, item := range m.Items {
		if item != nil {
			item.Destroy()
		}
	}
	m.Items = nil
}

// SetLabel sets the menu label.
func (m *Menu) SetLabel(label string) {
	m.label = label
}

// FindByLabel recursively searches for an item by label.
func (m *Menu) FindByLabel(label string) *MenuItem {
	for _, item := range m.Items {
		if item == nil {
			continue
		}
		if item.label == label {
			return item
		}
		if item.submenu != nil {
			if found := item.submenu.FindByLabel(label); found != nil {
				return found
			}
		}
	}
	return nil
}

// FindByRole recursively searches for an item by role.
func (m *Menu) FindByRole(role Role) *MenuItem {
	for _, item := range m.Items {
		if item == nil {
			continue
		}
		if item.role == role {
			return item
		}
		if item.submenu != nil {
			if found := item.submenu.FindByRole(role); found != nil {
				return found
			}
		}
	}
	return nil
}

// RemoveMenuItem removes an item from the menu tree.
func (m *Menu) RemoveMenuItem(target *MenuItem) {
	for index, item := range m.Items {
		if item == target {
			m.Items = append(m.Items[:index], m.Items[index+1:]...)
			return
		}
		if item != nil && item.submenu != nil {
			item.submenu.RemoveMenuItem(target)
		}
	}
}

// ItemAt returns the menu item at index, or nil when out of bounds.
func (m *Menu) ItemAt(index int) *MenuItem {
	if index < 0 || index >= len(m.Items) {
		return nil
	}
	return m.Items[index]
}

// Append appends another menu's items.
func (m *Menu) Append(in *Menu) {
	if in == nil {
		return
	}
	m.Items = append(m.Items, in.Items...)
}

// Prepend prepends another menu's items.
func (m *Menu) Prepend(in *Menu) {
	if in == nil {
		return
	}
	items := append([]*MenuItem(nil), in.Items...)
	m.Items = append(items, m.Items...)
}

// NewMenuFromItems creates a menu from an item list.
func NewMenuFromItems(item *MenuItem, items ...*MenuItem) *Menu {
	menu := &Menu{}
	if item != nil {
		menu.Items = append(menu.Items, item)
	}
	menu.Items = append(menu.Items, items...)
	return menu
}

// NewSubmenu creates a submenu item from an existing menu.
func NewSubmenu(s string, items *Menu) *MenuItem {
	result := NewSubMenuItem(s)
	result.submenu = items
	return result
}

// Set is an alias for SetApplicationMenu.
func (mm *MenuManager) Set(menu *Menu) {
	mm.SetApplicationMenu(menu)
}

// GetApplicationMenu returns the current application menu.
func (mm *MenuManager) GetApplicationMenu() *Menu {
	return mm.applicationMenu
}

// New creates a new empty menu.
func (mm *MenuManager) New() *Menu {
	return &Menu{}
}

// ShowAbout is a no-op in the stub.
func (mm *MenuManager) ShowAbout() {}

// NewServicesMenu returns the platform services role item.
func NewServicesMenu() *MenuItem {
	return NewRole(ServicesMenu)
}
