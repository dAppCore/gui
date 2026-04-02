// pkg/environment/messages.go
package environment

// QueryTheme returns the current theme. Result: ThemeInfo
type QueryTheme struct{}

// QueryInfo returns environment information. Result: EnvironmentInfo
type QueryInfo struct{}

// QueryAccentColour returns the system accent colour. Result: string
type QueryAccentColour struct{}

// TaskOpenFileManager opens the system file manager. Result: error only
type TaskOpenFileManager struct {
	Path   string `json:"path"`
	Select bool   `json:"select"`
}

// TaskSetTheme applies an application theme override when supported.
type TaskSetTheme struct {
	IsDark bool `json:"isDark"`
}

// ActionThemeChanged is broadcast when the system theme changes.
type ActionThemeChanged struct {
	IsDark bool `json:"isDark"`
}
