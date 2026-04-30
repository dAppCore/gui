package application

import "sync"

// KeyBinding pairs an accelerator string with its registered callback.
//
//	binding := &KeyBinding{Accelerator: "CmdOrCtrl+K", Callback: handler}
type KeyBinding struct {
	Accelerator string
	Callback    func(window Window)
}

// KeyBindingManager stores and dispatches global key bindings.
//
//	manager.Add("CmdOrCtrl+K", func(w Window) { w.Focus() })
//	handled := manager.Process("CmdOrCtrl+K", currentWindow)
type KeyBindingManager struct {
	mu       sync.RWMutex
	bindings map[string]func(window Window)
}

// Add registers a callback for the given accelerator string.
//
//	manager.Add("CmdOrCtrl+Shift+P", func(w Window) { launchCommandPalette(w) })
func (m *KeyBindingManager) Add(accelerator string, callback func(window Window)) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.bindings == nil {
		m.bindings = make(map[string]func(window Window))
	}
	m.bindings[accelerator] = callback
}

// Remove deregisters the callback for the given accelerator string.
//
//	manager.Remove("CmdOrCtrl+Shift+P")
func (m *KeyBindingManager) Remove(accelerator string) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.bindings, accelerator)
}

// Process fires the callback for the given accelerator and returns true if handled.
//
//	if manager.Process("CmdOrCtrl+K", window) { return }
func (m *KeyBindingManager) Process(accelerator string, window Window) (handled bool) {
	if m == nil {
		return false
	}
	m.mu.RLock()
	callback, exists := m.bindings[accelerator]
	m.mu.RUnlock()
	if exists && callback != nil {
		handled = true
		defer func() {
			if recover() != nil {
				handled = false
			}
		}()
		callback(window)
		return
	}
	return false
}

// GetAll returns a snapshot of all registered key bindings.
//
//	for _, kb := range manager.GetAll() { log(kb.Accelerator) }
func (m *KeyBindingManager) GetAll() []*KeyBinding {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	bindings := make([]*KeyBinding, 0, len(m.bindings))
	for accelerator, callback := range m.bindings {
		bindings = append(bindings, &KeyBinding{
			Accelerator: accelerator,
			Callback:    callback,
		})
	}
	return bindings
}
