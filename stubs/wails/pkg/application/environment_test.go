package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentManager_IsDarkMode_Good(t *testing.T) {
	manager := &EnvironmentManager{}
	manager.SetDarkMode(true)

	assert.True(t, manager.IsDarkMode())
}

func TestEnvironmentManager_IsDarkMode_Bad(t *testing.T) {
	manager := &EnvironmentManager{}

	assert.False(t, manager.IsDarkMode())
}

func TestEnvironmentManager_IsDarkMode_Ugly(t *testing.T) {
	manager := &EnvironmentManager{}
	manager.SetDarkMode(true)
	manager.SetDarkMode(false)

	assert.False(t, manager.IsDarkMode())
}

func TestEnvironmentManager_GetAccentColor_Good(t *testing.T) {
	manager := &EnvironmentManager{}
	manager.SetAccentColour("rgb(1,2,3)")

	assert.Equal(t, "rgb(1,2,3)", manager.GetAccentColor())
}

func TestEnvironmentManager_GetAccentColor_Bad(t *testing.T) {
	manager := &EnvironmentManager{}

	assert.Equal(t, "rgb(0,122,255)", manager.GetAccentColor())
}

func TestEnvironmentManager_GetAccentColor_Ugly(t *testing.T) {
	manager := &EnvironmentManager{}
	manager.SetAccentColour("")

	assert.Equal(t, "rgb(0,122,255)", manager.GetAccentColor())
}

func TestEnvironmentManager_Info_Good(t *testing.T) {
	manager := &EnvironmentManager{}
	manager.SetDarkMode(true)
	manager.SetAccentColour("rgb(1,2,3)")
	manager.operatingSystem = "linux"
	manager.architecture = "amd64"
	manager.debugMode = true

	got := manager.Info()

	assert.Equal(t, "linux", got.OS)
	assert.Equal(t, "amd64", got.Arch)
	assert.True(t, got.Debug)
	assert.True(t, got.IsDarkMode)
	assert.Equal(t, "rgb(1,2,3)", got.AccentColour)
}

func TestEnvironmentManager_Info_Bad(t *testing.T) {
	manager := &EnvironmentManager{}

	got := manager.Info()

	assert.Empty(t, got.OS)
	assert.Empty(t, got.Arch)
	assert.False(t, got.Debug)
	assert.False(t, got.IsDarkMode)
	assert.Empty(t, got.AccentColour)
}

func TestEnvironmentManager_Info_Ugly(t *testing.T) {
	manager := &EnvironmentManager{}
	manager.operatingSystem = "plan9"
	manager.architecture = "riscv64"
	manager.debugMode = true
	manager.SetDarkMode(true)

	got := manager.Info()

	assert.Equal(t, "plan9", got.OS)
	assert.Equal(t, "riscv64", got.Arch)
	assert.True(t, got.Debug)
	assert.True(t, got.IsDarkMode)
}
