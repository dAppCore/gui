package deno

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSidecar_New_Good(t *testing.T) {
	manager := New(Options{})

	status := manager.Status()
	assert.Equal(t, "deno", status.Binary)
	assert.False(t, status.Running)
	assert.Zero(t, status.PID)
}

func TestSidecar_New_Bad(t *testing.T) {
	manager := New(Options{Binary: "/usr/local/bin/deno-custom", Args: []string{"fmt"}})

	status := manager.Status()
	assert.Equal(t, "/usr/local/bin/deno-custom", status.Binary)
	assert.False(t, status.Running)
}

func TestSidecar_New_Ugly(t *testing.T) {
	manager := New(Options{Binary: "   "})

	status := manager.Status()
	assert.Equal(t, "deno", status.Binary)
	assert.False(t, status.Running)
}
