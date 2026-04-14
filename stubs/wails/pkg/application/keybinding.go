package application

import "sync"

// KeyBinding pairs an accelerator string with its callback.
// binding := &KeyBinding{Accelerator: "Ctrl+K", Callback: fn}
type KeyBinding struct {
	Accelerator string
	Callback    func(window Window)
}

// KeyBindingManager holds all registered global key bindings in memory.
// manager.Add("Ctrl+K", fn) — manager.Remove("Ctrl+K") — manager.GetAll()
type KeyBindingManager struct {
	mu       sync.RWMutex
	bindings map[string]func(window Window)
}

// NewKeyBindingManager constructs an empty KeyBindingManager.
// manager := NewKeyBindingManager()
func NewKeyBindingManager() *KeyBindingManager {
	return &KeyBindingManager{
		bindings: make(map[string]func(window Window)),
	}
}

// Add registers a callback for the given accelerator string.
// manager.Add("Ctrl+Shift+P", func(w Window) { w.Focus() })
func (m *KeyBindingManager) Add(accelerator string, callback func(window Window)) {
	m.mu.Lock()
	m.bindings[accelerator] = callback
	m.mu.Unlock()
}

// Remove deletes the binding for the given accelerator.
// manager.Remove("Ctrl+Shift+P")
func (m *KeyBindingManager) Remove(accelerator string) {
	m.mu.Lock()
	delete(m.bindings, accelerator)
	m.mu.Unlock()
}

// Process fires the callback for accelerator if registered, returning true when handled.
// handled := manager.Process("Ctrl+K", window)
func (m *KeyBindingManager) Process(accelerator string, window Window) bool {
	m.mu.RLock()
	callback, exists := m.bindings[accelerator]
	m.mu.RUnlock()
	if exists && callback != nil {
		callback(window)
		return true
	}
	return false
}

// GetAll returns a snapshot of all registered bindings.
// for _, b := range manager.GetAll() { fmt.Println(b.Accelerator) }
func (m *KeyBindingManager) GetAll() []*KeyBinding {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*KeyBinding, 0, len(m.bindings))
	for accelerator, callback := range m.bindings {
		result = append(result, &KeyBinding{
			Accelerator: accelerator,
			Callback:    callback,
		})
	}
	return result
}
