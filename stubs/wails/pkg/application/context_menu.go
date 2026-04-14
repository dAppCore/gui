package application

import "sync"

// ContextMenu is a named Menu used as a right-click context menu.
//
//	cm := manager.New()
//	cm.Add("Cut").OnClick(func(*Context) { ... })
type ContextMenu struct {
	*Menu
	name string
}

// ContextMenuManager manages named context menus for the application.
//
//	manager.Add("fileList", cm)
//	menu, ok := manager.Get("fileList")
type ContextMenuManager struct {
	mu    sync.RWMutex
	menus map[string]*ContextMenu
}

// New creates an empty, unnamed ContextMenu ready for population.
//
//	cm := manager.New()
//	cm.Add("Open")
func (cmm *ContextMenuManager) New() *ContextMenu {
	return &ContextMenu{Menu: NewMenu()}
}

// Add registers a ContextMenu under the given name, replacing any existing entry.
//
//	manager.Add("fileList", cm)
func (cmm *ContextMenuManager) Add(name string, menu *ContextMenu) {
	cmm.mu.Lock()
	defer cmm.mu.Unlock()
	if cmm.menus == nil {
		cmm.menus = make(map[string]*ContextMenu)
	}
	cmm.menus[name] = menu
}

// Remove unregisters the context menu with the given name.
//
//	manager.Remove("fileList")
func (cmm *ContextMenuManager) Remove(name string) {
	cmm.mu.Lock()
	defer cmm.mu.Unlock()
	delete(cmm.menus, name)
}

// Get retrieves a registered context menu by name.
//
//	menu, ok := manager.Get("fileList")
func (cmm *ContextMenuManager) Get(name string) (*ContextMenu, bool) {
	cmm.mu.RLock()
	defer cmm.mu.RUnlock()
	menu, exists := cmm.menus[name]
	return menu, exists
}

// GetAll returns all registered context menus as a slice.
//
//	for _, cm := range manager.GetAll() { ... }
func (cmm *ContextMenuManager) GetAll() []*ContextMenu {
	cmm.mu.RLock()
	defer cmm.mu.RUnlock()
	result := make([]*ContextMenu, 0, len(cmm.menus))
	for _, menu := range cmm.menus {
		result = append(result, menu)
	}
	return result
}
