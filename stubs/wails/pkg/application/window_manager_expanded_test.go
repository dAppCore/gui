package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowManagerExpanded_Get_Good(t *testing.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	require.Same(t, first, manager.Get("first"))
	require.Same(t, second, manager.Get("second"))
	assert.Len(t, manager.GetAll(), 2)
}

func TestWindowManagerExpanded_Get_Bad(t *testing.T) {
	manager := &WindowManager{}

	assert.Nil(t, manager.Get("missing"))
}

func TestWindowManagerExpanded_Get_Ugly(t *testing.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "dup"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "dup"})

	assert.Same(t, first, manager.Get("dup"))
	assert.NotSame(t, first, second)
}

func TestWindowManagerExpanded_GetByID_Good(t *testing.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	assert.Same(t, first, manager.GetByID(1))
	assert.Same(t, second, manager.GetByID(2))
}

func TestWindowManagerExpanded_GetByID_Bad(t *testing.T) {
	manager := &WindowManager{}

	assert.Nil(t, manager.GetByID(99))
}

func TestWindowManagerExpanded_GetByID_Ugly(t *testing.T) {
	manager := &WindowManager{}
	manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	assert.Nil(t, manager.GetByID(0))
	assert.Nil(t, manager.GetByID(3))
}

func TestWindowManagerExpanded_NilReceiver_IsSafe(t *testing.T) {
	var manager *WindowManager

	assert.NotPanics(t, func() {
		assert.Nil(t, manager.NewWithOptions(WebviewWindowOptions{Name: "ignored"}))
		assert.Nil(t, manager.Get("missing"))
		assert.Nil(t, manager.GetByID(1))
		assert.Nil(t, manager.GetAll())
	})
}
