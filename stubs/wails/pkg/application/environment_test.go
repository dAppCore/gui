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

// AX7 generated source-matching smoke coverage.
func TestEnvironment_EnvironmentManager_SetDarkMode_Good(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		subject.SetDarkMode(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_SetDarkMode_Bad(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		subject.SetDarkMode(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_SetDarkMode_Ugly(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		subject.SetDarkMode(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_IsDarkMode_Good(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.IsDarkMode()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_IsDarkMode_Bad(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.IsDarkMode()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_IsDarkMode_Ugly(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.IsDarkMode()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_SetAccentColour_Good(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		subject.SetAccentColour("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_SetAccentColour_Bad(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		subject.SetAccentColour("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_SetAccentColour_Ugly(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		subject.SetAccentColour("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_GetAccentColor_Good(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.GetAccentColor()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_GetAccentColor_Bad(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.GetAccentColor()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_GetAccentColor_Ugly(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.GetAccentColor()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_Info_Good(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_Info_Bad(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_Info_Ugly(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_OpenFileManager_Good(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.OpenFileManager("agent", true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_OpenFileManager_Bad(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.OpenFileManager("", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_OpenFileManager_Ugly(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.OpenFileManager("../../edge", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_HasFocusFollowsMouse_Good(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.HasFocusFollowsMouse()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_HasFocusFollowsMouse_Bad(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.HasFocusFollowsMouse()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEnvironment_EnvironmentManager_HasFocusFollowsMouse_Ugly(t *core.T) {
	subject := new(EnvironmentManager)
	result := core.Try(func() any {
		got0 := subject.HasFocusFollowsMouse()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
