package application

// Get returns the first window whose Name() matches the given name, or nil if
// no match is found.
//
//	w := app.Window.Get("main")
//	if w == nil { panic("main window not registered") }
func (wm *WindowManager) Get(name string) Window {
	if wm == nil {
		return nil
	}
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	for _, window := range wm.windows {
		if window.Name() == name {
			return window
		}
	}
	return nil
}

// GetByID returns the window with the given numeric ID, or nil if not found.
//
//	w := app.Window.GetByID(1)
//	if w != nil { w.Focus() }
func (wm *WindowManager) GetByID(id uint) Window {
	if wm == nil {
		return nil
	}
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	for _, window := range wm.windows {
		if window.ID() == id {
			return window
		}
	}
	return nil
}
