package application

import "sync"

// EnvironmentInfo holds information about the host environment.
//
//	info := manager.Info()
//	if info.IsDarkMode { applyDarkTheme() }
type EnvironmentInfo struct {
	OS           string
	Arch         string
	Debug        bool
	IsDarkMode   bool
	AccentColour string
	PlatformInfo map[string]any
}

// EnvironmentManager tracks environment state in-memory.
//
//	manager := &EnvironmentManager{}
//	manager.SetDarkMode(true)
//	dark := manager.IsDarkMode() // true
type EnvironmentManager struct {
	mu           sync.RWMutex
	darkMode     bool
	accentColour string
	operatingSystem   string
	architecture string
	debugMode    bool
}

// SetDarkMode sets the dark mode state used by IsDarkMode.
//
//	manager.SetDarkMode(true)
func (em *EnvironmentManager) SetDarkMode(darkMode bool) {
	em.mu.Lock()
	em.darkMode = darkMode
	em.mu.Unlock()
}

// IsDarkMode returns true when the environment is in dark mode.
//
//	if manager.IsDarkMode() { applyDarkTheme() }
func (em *EnvironmentManager) IsDarkMode() bool {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return em.darkMode
}

// SetAccentColour sets the accent colour returned by GetAccentColor.
//
//	manager.SetAccentColour("rgb(0,122,255)")
func (em *EnvironmentManager) SetAccentColour(colour string) {
	em.mu.Lock()
	em.accentColour = colour
	em.mu.Unlock()
}

// GetAccentColor returns the stored accent colour, or the default blue if unset.
//
//	colour := manager.GetAccentColor() // "rgb(0,122,255)"
func (em *EnvironmentManager) GetAccentColor() string {
	em.mu.RLock()
	defer em.mu.RUnlock()
	if em.accentColour == "" {
		return "rgb(0,122,255)"
	}
	return em.accentColour
}

// Info returns a snapshot of the current environment state.
//
//	info := manager.Info()
//	_ = info.OS // e.g. "linux"
func (em *EnvironmentManager) Info() EnvironmentInfo {
	em.mu.RLock()
	defer em.mu.RUnlock()
	return EnvironmentInfo{
		OS:           em.operatingSystem,
		Arch:         em.architecture,
		Debug:        em.debugMode,
		IsDarkMode:   em.darkMode,
		AccentColour: em.accentColour,
	}
}
