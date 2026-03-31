package application

import "sync"

// ContextMenu is a named menu used for right-click context menus.
//
//	menu := &ContextMenu{name: "file-list", Menu: NewMenu()}
//	menu.Add("Open")
type ContextMenu struct {
	*Menu
	name string
}

// ContextMenuManager manages all context menu operations in-memory.
//
//	manager := &ContextMenuManager{}
//	menu := manager.New()
//	manager.Add("file-list", menu)
//	retrieved, ok := manager.Get("file-list")
type ContextMenuManager struct {
	mu           sync.RWMutex
	contextMenus map[string]*ContextMenu
}

// New creates a new context menu.
//
//	menu := manager.New()
//	menu.Add("Delete")
func (cmm *ContextMenuManager) New() *ContextMenu {
	return &ContextMenu{
		Menu: NewMenu(),
	}
}

// Add registers a context menu under the given name.
//
//	manager.Add("item-actions", menu)
func (cmm *ContextMenuManager) Add(name string, menu *ContextMenu) {
	cmm.mu.Lock()
	defer cmm.mu.Unlock()
	if cmm.contextMenus == nil {
		cmm.contextMenus = make(map[string]*ContextMenu)
	}
	cmm.contextMenus[name] = menu
}

// Remove removes a context menu by name.
//
//	manager.Remove("item-actions")
func (cmm *ContextMenuManager) Remove(name string) {
	cmm.mu.Lock()
	defer cmm.mu.Unlock()
	if cmm.contextMenus != nil {
		delete(cmm.contextMenus, name)
	}
}

// Get retrieves a context menu by name.
//
//	menu, ok := manager.Get("item-actions")
//	if !ok { return nil }
func (cmm *ContextMenuManager) Get(name string) (*ContextMenu, bool) {
	cmm.mu.RLock()
	defer cmm.mu.RUnlock()
	if cmm.contextMenus == nil {
		return nil, false
	}
	menu, exists := cmm.contextMenus[name]
	return menu, exists
}

// GetAll returns all registered context menus as a slice.
//
//	menus := manager.GetAll()
//	for _, m := range menus { _ = m }
func (cmm *ContextMenuManager) GetAll() []*ContextMenu {
	cmm.mu.RLock()
	defer cmm.mu.RUnlock()
	result := make([]*ContextMenu, 0, len(cmm.contextMenus))
	for _, menu := range cmm.contextMenus {
		result = append(result, menu)
	}
	return result
}
