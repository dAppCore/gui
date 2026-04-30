//go:build compliance

package application

import core "dappco.re/go"

func ExampleEnvironmentManager_Info() {
	core.Println("EnvironmentManager_Info")
	// Output:
	// EnvironmentManager_Info
}

func ExampleEnvironmentManager_IsDarkMode() {
	core.Println("EnvironmentManager_IsDarkMode")
	// Output:
	// EnvironmentManager_IsDarkMode
}

func ExampleEnvironmentManager_GetAccentColor() {
	core.Println("EnvironmentManager_GetAccentColor")
	// Output:
	// EnvironmentManager_GetAccentColor
}

func ExampleEnvironmentManager_OpenFileManager() {
	core.Println("EnvironmentManager_OpenFileManager")
	// Output:
	// EnvironmentManager_OpenFileManager
}

func ExampleEnvironmentManager_HasFocusFollowsMouse() {
	core.Println("EnvironmentManager_HasFocusFollowsMouse")
	// Output:
	// EnvironmentManager_HasFocusFollowsMouse
}
