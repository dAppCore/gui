package application

import (
	"runtime"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

// EnvironmentInfo holds runtime information about the host OS and build.
//
//	info := manager.Info()
//	fmt.Println(info.OS, info.Arch)
type EnvironmentInfo struct {
	OS           string
	Arch         string
	Debug        bool
	OSInfo       *operatingsystem.OS
	PlatformInfo map[string]any
}

// EnvironmentManager provides queries about the host OS environment.
//
//	if manager.IsDarkMode() { applyDarkTheme() }
//	accent := manager.GetAccentColor()
type EnvironmentManager struct {
	darkMode    bool
	accentColor string
}

// IsDarkMode returns true when the OS is using a dark colour scheme.
//
//	if manager.IsDarkMode() { applyDarkTheme() }
func (em *EnvironmentManager) IsDarkMode() bool {
	return em.darkMode
}

// GetAccentColor returns the OS accent colour as an rgb() string.
//
//	colour := manager.GetAccentColor() // e.g. "rgb(0,122,255)"
func (em *EnvironmentManager) GetAccentColor() string {
	if em.accentColor == "" {
		return "rgb(0,122,255)"
	}
	return em.accentColor
}

// Info returns a snapshot of OS and build environment information.
//
//	info := manager.Info()
func (em *EnvironmentManager) Info() EnvironmentInfo {
	return EnvironmentInfo{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		OSInfo:       &operatingsystem.OS{Name: runtime.GOOS, Version: "stub"},
		PlatformInfo: make(map[string]any),
	}
}

// OpenFileManager opens the file manager at the given path, optionally selecting the file.
// No-op in the stub.
//
//	err := manager.OpenFileManager("/home/user/docs", false)
func (em *EnvironmentManager) OpenFileManager(path string, selectFile bool) error {
	return nil
}

// HasFocusFollowsMouse reports whether the Linux desktop is configured for
// focus-follows-mouse. Always returns false in the stub.
//
//	if manager.HasFocusFollowsMouse() { adjustMouseBehaviour() }
func (em *EnvironmentManager) HasFocusFollowsMouse() bool {
	return false
}
