package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyBindingManager_Add_Good(t *testing.T) {
	manager := &KeyBindingManager{}
	calls := 0

	manager.Add("CmdOrCtrl+K", func(window Window) {
		calls++
		assert.Nil(t, window)
	})

	handled := manager.Process("CmdOrCtrl+K", nil)

	assert.True(t, handled)
	assert.Equal(t, 1, calls)
	assert.Len(t, manager.GetAll(), 1)
}

func TestKeyBindingManager_Add_Bad(t *testing.T) {
	manager := &KeyBindingManager{}

	handled := manager.Process("missing", nil)

	assert.False(t, handled)
	assert.Empty(t, manager.GetAll())
}

func TestKeyBindingManager_Add_Ugly(t *testing.T) {
	manager := &KeyBindingManager{}
	calls := 0

	manager.Add("CmdOrCtrl+K", func(Window) { calls++ })
	manager.Add("CmdOrCtrl+K", func(Window) { calls += 10 })

	handled := manager.Process("CmdOrCtrl+K", nil)

	assert.True(t, handled)
	assert.Equal(t, 10, calls)
}

func TestKeyBindingManager_Process_RecoversFromPanic(t *testing.T) {
	manager := &KeyBindingManager{}

	manager.Add("CmdOrCtrl+K", func(Window) {
		panic("boom")
	})

	assert.False(t, manager.Process("CmdOrCtrl+K", nil))
}

func TestKeyBindingManager_Remove_Good(t *testing.T) {
	manager := &KeyBindingManager{}
	manager.Add("CmdOrCtrl+K", func(Window) {})

	manager.Remove("CmdOrCtrl+K")

	assert.False(t, manager.Process("CmdOrCtrl+K", nil))
	assert.Empty(t, manager.GetAll())
}

func TestKeyBindingManager_Remove_Bad(t *testing.T) {
	manager := &KeyBindingManager{}

	manager.Remove("missing")

	assert.Empty(t, manager.GetAll())
}

func TestKeyBindingManager_Remove_Ugly(t *testing.T) {
	manager := &KeyBindingManager{}

	manager.Remove("")

	assert.Empty(t, manager.GetAll())
}
