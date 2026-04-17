package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextMenuManager_Add_Good(t *testing.T) {
	manager := &ContextMenuManager{}
	menu := manager.New()
	menu.Add("Open")

	manager.Add("files", menu)

	got, ok := manager.Get("files")
	require.True(t, ok)
	require.Same(t, menu, got)
	assert.Len(t, manager.GetAll(), 1)
}

func TestContextMenuManager_Add_Bad(t *testing.T) {
	manager := &ContextMenuManager{}

	manager.Add("empty", nil)

	got, ok := manager.Get("empty")
	assert.True(t, ok)
	assert.Nil(t, got)
}

func TestContextMenuManager_Add_Ugly(t *testing.T) {
	manager := &ContextMenuManager{}
	first := manager.New()
	second := manager.New()

	manager.Add("dup", first)
	manager.Add("dup", second)

	got, ok := manager.Get("dup")
	require.True(t, ok)
	require.Same(t, second, got)
	assert.Len(t, manager.GetAll(), 1)
}

func TestContextMenuManager_Remove_Good(t *testing.T) {
	manager := &ContextMenuManager{}
	menu := manager.New()
	manager.Add("files", menu)

	manager.Remove("files")

	_, ok := manager.Get("files")
	assert.False(t, ok)
	assert.Empty(t, manager.GetAll())
}

func TestContextMenuManager_Remove_Bad(t *testing.T) {
	manager := &ContextMenuManager{}

	manager.Remove("missing")

	assert.Empty(t, manager.GetAll())
}

func TestContextMenuManager_Remove_Ugly(t *testing.T) {
	manager := &ContextMenuManager{}
	manager.Remove("")

	assert.Empty(t, manager.GetAll())
}
