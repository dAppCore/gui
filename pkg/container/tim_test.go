package container

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTIMManager_NewTIMManager_Good(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	state := manager.State()
	assert.Equal(t, "coregui-tim", state.Name)
	assert.Equal(t, "ghcr.io/lthn/core/tim:latest", state.Image)
	assert.Equal(t, RuntimeDocker, state.Runtime)
	assert.Equal(t, "stopped", state.Status)
	assert.Empty(t, state.StartedAt)
}

func TestTIMManager_NewTIMManager_Bad(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Name:  " ",
		Image: " ",
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state := manager.State()
	assert.Equal(t, "coregui-tim", state.Name)
	assert.Equal(t, "ghcr.io/lthn/core/tim:latest", state.Image)
	assert.Equal(t, RuntimeNone, state.Runtime)
	assert.Equal(t, "stopped", state.Status)
}

func TestTIMManager_NewTIMManager_Ugly(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Name:      "worker-node",
		Image:     "ghcr.io/example/tim:edge",
		Command:   []string{"alpha", "beta"},
		DataDir:   "/var/lib/tim",
		Runtime:   RuntimePodman,
		Resources: TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	state := manager.State()
	assert.Equal(t, "worker-node", state.Name)
	assert.Equal(t, "ghcr.io/example/tim:edge", state.Image)
	assert.Equal(t, RuntimePodman, state.Runtime)
	assert.Equal(t, []string{"alpha", "beta"}, state.Command)
	assert.Equal(t, "/var/lib/tim", state.DataDir)
	assert.Equal(t, TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"}, state.Resources)
}

func TestTIMManager_Start_Good(t *testing.T) {
	var calls []string
	manager := NewTIMManager(TIMOptions{
		Name:    "coregui-tim",
		Image:   "ghcr.io/example/tim:latest",
		Command: []string{"sleep", "1"},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(_ context.Context, name string, args ...string) error {
			calls = append(calls, append([]string{name}, args...)...)
			return nil
		},
		Now: func() time.Time {
			return time.Unix(456, 0).UTC()
		},
	})

	state, err := manager.Start(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "running", state.Status)
	assert.Equal(t, time.Unix(456, 0).UTC(), state.StartedAt)
	assert.Equal(t, "docker", calls[0])
	assert.Contains(t, calls, "run")
	assert.Contains(t, calls, "--rm")
	assert.Contains(t, calls, "--name")
	assert.Contains(t, calls, "coregui-tim")
	assert.Contains(t, calls, "ghcr.io/example/tim:latest")
}

func TestTIMManager_Start_Bad(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state, err := manager.Start(context.Background())
	require.Error(t, err)
	assert.Equal(t, RuntimeNone, state.Runtime)
	assert.Equal(t, "stopped", state.Status)
	assert.Contains(t, err.Error(), "no supported container runtime detected")
}

func TestTIMManager_Start_Ugly(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(context.Context, string, ...string) error {
			return errors.New("docker failed")
		},
	})

	state, err := manager.Start(context.Background())
	require.Error(t, err)
	assert.Equal(t, "error", state.Status)
	assert.Contains(t, err.Error(), "docker failed")
}

func TestTIMManager_Stop_Good(t *testing.T) {
	var calls []string
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(_ context.Context, name string, args ...string) error {
			calls = append(calls, append([]string{name}, args...)...)
			return nil
		},
	})

	started, err := manager.Start(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "running", started.Status)

	stopped, err := manager.Stop(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "stopped", stopped.Status)
	require.GreaterOrEqual(t, len(calls), 3)
	assert.Equal(t, "docker", calls[len(calls)-3])
	assert.Equal(t, "stop", calls[len(calls)-2])
	assert.Equal(t, "coregui-tim", calls[len(calls)-1])
}

func TestTIMManager_Stop_Bad(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(context.Context, string, ...string) error {
			return errors.New("stop failed")
		},
	})

	_, err := manager.Start(context.Background())
	require.Error(t, err)

	state, err := manager.Stop(context.Background())
	require.Error(t, err)
	assert.Equal(t, "error", state.Status)
	assert.Contains(t, err.Error(), "stop failed")
}

func TestTIMManager_Stop_Ugly(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state, err := manager.Stop(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "stopped", state.Status)
}

func TestTIMManager_runtimeCommand_Good(t *testing.T) {
	cases := []struct {
		name     string
		runtime  ContainerRuntime
		verb     string
		wantBin  string
		wantArgs []string
	}{
		{
			name:    "docker run",
			runtime: RuntimeDocker,
			verb:    "run",
			wantBin: "docker",
			wantArgs: []string{"run", "-d", "--rm", "--name", "tim"},
		},
		{
			name:    "apple stop",
			runtime: RuntimeApple,
			verb:    "stop",
			wantBin: "container",
			wantArgs: []string{"stop", "tim"},
		},
		{
			name:    "podman run",
			runtime: RuntimePodman,
			verb:    "run",
			wantBin: "podman",
			wantArgs: []string{"run", "-d", "--replace", "--name", "tim"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			manager := NewTIMManager(TIMOptions{
				Name:    "tim",
				Image:   "image",
				Runtime: tc.runtime,
				Detect: func() ContainerRuntime {
					return RuntimeDocker
				},
			})

			gotBin, gotArgs := manager.runtimeCommand(tc.runtime, tc.verb)
			assert.Equal(t, tc.wantBin, gotBin)
			assert.Equal(t, tc.wantArgs, gotArgs[:len(tc.wantArgs)])
			if tc.verb == "run" {
				assert.Contains(t, gotArgs, "image")
			}
		})
	}
}

func TestTIMManager_runtimeCommand_Bad(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Name:    "tim",
		Image:   "image",
		Runtime: RuntimeNone,
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	bin, args := manager.runtimeCommand(RuntimeNone, "run")
	assert.Equal(t, "docker", bin)
	assert.Contains(t, args, "--rm")
}

func TestTIMManager_runtimeCommand_Ugly(t *testing.T) {
	manager := NewTIMManager(TIMOptions{
		Name: "tim",
		Image: "image",
		Resources: TIMResources{
			CPUCores: 2,
			MemoryMB: 512,
			GPU:      "all",
		},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	bin, args := manager.runtimeCommand(RuntimeDocker, "run")
	require.Equal(t, "docker", bin)
	assert.Contains(t, args, "--cpus")
	assert.Contains(t, args, "2")
	assert.Contains(t, args, "--memory")
	assert.Contains(t, args, "512m")
	assert.Contains(t, args, "--gpus")
	assert.Contains(t, args, "all")
}

func TestTIMManager_resourceArgs_Good(t *testing.T) {
	assert.Nil(t, resourceArgs(TIMResources{}))
}

func TestTIMManager_resourceArgs_Bad(t *testing.T) {
	args := resourceArgs(TIMResources{CPUCores: 2})

	assert.Equal(t, []string{"--cpus", "2"}, args)
}

func TestTIMManager_resourceArgs_Ugly(t *testing.T) {
	args := resourceArgs(TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"})

	assert.Equal(t, []string{"--cpus", "2", "--memory", "512m", "--gpus", "all"}, args)
}

func TestTIMManager_coalesceRuntime_Good(t *testing.T) {
	assert.Equal(t, RuntimeApple, coalesceRuntime(RuntimeApple, RuntimeDocker))
}

func TestTIMManager_coalesceRuntime_Bad(t *testing.T) {
	assert.Equal(t, RuntimeDocker, coalesceRuntime(RuntimeNone, RuntimeDocker))
}

func TestTIMManager_coalesceRuntime_Ugly(t *testing.T) {
	assert.Equal(t, RuntimeNone, coalesceRuntime(RuntimeNone, RuntimeNone))
}
