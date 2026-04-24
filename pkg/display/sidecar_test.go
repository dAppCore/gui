package display

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/deno"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSidecar_SplitCommandArgs_Good(t *testing.T) {
	assert.Equal(t, []string{"--import-map", "map.json", "--watch"}, splitCommandArgs("--import-map map.json --watch"))
}

func TestSidecar_SplitCommandArgs_Bad(t *testing.T) {
	assert.Nil(t, splitCommandArgs(""))
	assert.Nil(t, splitCommandArgs("   "))
}

func TestSidecar_SplitCommandArgs_Ugly(t *testing.T) {
	assert.Equal(t, []string{"--flag", "--another", "value"}, splitCommandArgs("\t--flag\n--another   value\t"))
}

func TestSidecar_EnsureSidecar_Good(t *testing.T) {
	t.Setenv("CORE_DENO_BINARY", "/usr/local/bin/deno-custom")
	t.Setenv("CORE_DENO_DIR", "/tmp/core-deno")
	t.Setenv("CORE_DENO_ARGS", "--import-map map.json")

	svc := &Service{}
	manager := svc.ensureSidecar()

	require.NotNil(t, manager)
	status := manager.Status()
	assert.Equal(t, "/usr/local/bin/deno-custom", status.Binary)
	assert.False(t, status.Running)
}

func TestSidecar_EnsureSidecar_Bad(t *testing.T) {
	svc := &Service{sidecar: deno.New(deno.Options{Binary: "custom-deno"})}

	manager := svc.ensureSidecar()

	require.Same(t, svc.sidecar, manager)
	assert.Equal(t, "custom-deno", manager.Status().Binary)
}

func TestSidecar_EnsureSidecar_Ugly(t *testing.T) {
	t.Setenv("CORE_DENO_BINARY", "   ")
	t.Setenv("CORE_DENO_DIR", "")
	t.Setenv("CORE_DENO_ARGS", "   ")

	svc := &Service{}
	manager := svc.ensureSidecar()

	require.NotNil(t, manager)
	assert.Equal(t, "deno", manager.Status().Binary)
}

func TestSidecar_RegisterActions_StartFailureClearsSidecar(t *testing.T) {
	t.Setenv("CORE_DENO_ENABLE", "1")
	t.Setenv("CORE_DENO_BINARY", "/definitely/not/a/real/deno")

	c := core.New(core.WithServiceLock())
	svc := &Service{ServiceRuntime: core.NewServiceRuntime(c, Options{})}

	svc.registerSidecarActions()

	assert.Nil(t, svc.sidecar)
}

func TestSidecar_StatusAction_Good(t *testing.T) {
	t.Setenv("CORE_DENO_BINARY", "/opt/core/deno")

	_, c := newTestDisplayService(t)
	result := c.Action("display.sidecar.status").Run(context.Background(), core.Options{})

	require.True(t, result.OK)
	status, ok := result.Value.(deno.Status)
	require.True(t, ok)
	assert.Equal(t, "/opt/core/deno", status.Binary)
	assert.False(t, status.Running)
}
