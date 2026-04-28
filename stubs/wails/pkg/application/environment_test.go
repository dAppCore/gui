package application

import (
	core "dappco.re/go"
)

func TestEnvironmentManager_IsDarkMode_Good(t *core.T) {
	manager := &EnvironmentManager{}
	manager.SetDarkMode(true)

	core.AssertTrue(t, manager.IsDarkMode())
}

func TestEnvironmentManager_IsDarkMode_Bad(t *core.T) {
	manager := &EnvironmentManager{}

	core.AssertFalse(t, manager.IsDarkMode())
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestEnvironmentManager_IsDarkMode_Ugly(t *core.T) {
	manager := &EnvironmentManager{}
	manager.SetDarkMode(true)
	manager.SetDarkMode(false)

	core.AssertFalse(t, manager.IsDarkMode())
}

func TestEnvironmentManager_GetAccentColor_Good(t *core.T) {
	manager := &EnvironmentManager{}
	manager.SetAccentColour("rgb(1,2,3)")

	core.AssertEqual(t, "rgb(1,2,3)", manager.GetAccentColor())
}

func TestEnvironmentManager_GetAccentColor_Bad(t *core.T) {
	manager := &EnvironmentManager{}

	core.AssertEqual(t, "rgb(0,122,255)", manager.GetAccentColor())
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestEnvironmentManager_GetAccentColor_Ugly(t *core.T) {
	manager := &EnvironmentManager{}
	manager.SetAccentColour("")

	core.AssertEqual(t, "rgb(0,122,255)", manager.GetAccentColor())
}

func TestEnvironmentManager_Info_Good(t *core.T) {
	manager := &EnvironmentManager{}
	manager.SetDarkMode(true)
	manager.SetAccentColour("rgb(1,2,3)")
	manager.operatingSystem = "linux"
	manager.architecture = "amd64"
	manager.debugMode = true

	got := manager.Info()

	core.AssertEqual(t, "linux", got.OS)
	core.AssertEqual(t, "amd64", got.Arch)
	core.AssertTrue(t, got.Debug)
	core.AssertTrue(t, got.IsDarkMode)
	core.AssertEqual(t, "rgb(1,2,3)", got.AccentColour)
}

func TestEnvironmentManager_Info_Bad(t *core.T) {
	manager := &EnvironmentManager{}

	got := manager.Info()

	core.AssertEmpty(t, got.OS)
	core.AssertEmpty(t, got.Arch)
	core.AssertFalse(t, got.Debug)
	core.AssertFalse(t, got.IsDarkMode)
	core.AssertEmpty(t, got.AccentColour)
}

func TestEnvironmentManager_Info_Ugly(t *core.T) {
	manager := &EnvironmentManager{}
	manager.operatingSystem = "plan9"
	manager.architecture = "riscv64"
	manager.debugMode = true
	manager.SetDarkMode(true)

	got := manager.Info()

	core.AssertEqual(t, "plan9", got.OS)
	core.AssertEqual(t, "riscv64", got.Arch)
	core.AssertTrue(t, got.Debug)
	core.AssertTrue(t, got.IsDarkMode)
}
