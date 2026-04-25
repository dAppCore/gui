package container

import (
	"context"
	"errors"
	"testing"
	"time"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestContainerService(t *testing.T, options TIMOptions) (*Service, *core.Core) {
	t.Helper()

	var svc *Service
	c := core.New(
		core.WithService(func(c *core.Core) core.Result {
			svc = NewService(c, options)
			return core.Result{Value: svc, OK: true}
		}),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	require.NotNil(t, svc)
	return svc, c
}

func newInvalidTestContainerService(t *testing.T, options TIMOptions) (*Service, *core.Core, error) {
	t.Helper()

	if options.Detect == nil {
		options.Detect = func() ContainerRuntime {
			return RuntimeDocker
		}
	}

	var svc *Service
	c := core.New(
		core.WithService(func(c *core.Core) core.Result {
			svc = NewService(c, options)
			return core.Result{Value: svc, OK: true}
		}),
		core.WithServiceLock(),
	)
	result := c.ServiceStartup(context.Background(), nil)
	require.False(t, result.OK)
	require.NotNil(t, svc)
	err, ok := result.Value.(error)
	require.True(t, ok)
	return svc, c, err
}

func TestService_OptionsFromEnv_Good(t *testing.T) {
	t.Setenv("CORE_TIM_NAME", "  worker ")
	t.Setenv("CORE_TIM_IMAGE", " ghcr.io/example/tim:latest ")
	t.Setenv("CORE_TIM_COMMAND", "run,  --flag, , value ")
	t.Setenv("CORE_TIM_DATA_DIR", " /var/lib/core-tim ")
	t.Setenv("CORE_TIM_GPU", " all ")

	opts, err := OptionsFromEnvValidated()

	require.NoError(t, err)
	assert.Equal(t, "worker", opts.Name)
	assert.Equal(t, "ghcr.io/example/tim:latest", opts.Image)
	assert.Equal(t, []string{"run", "--flag", "value"}, opts.Command)
	assert.Equal(t, "/var/lib/core-tim", opts.DataDir)
	assert.Equal(t, "all", opts.Resources.GPU)
}

func TestService_OptionsFromEnv_Bad(t *testing.T) {
	t.Setenv("CORE_TIM_NAME", "")
	t.Setenv("CORE_TIM_IMAGE", "")
	t.Setenv("CORE_TIM_COMMAND", "")
	t.Setenv("CORE_TIM_DATA_DIR", "")
	t.Setenv("CORE_TIM_GPU", "")

	opts, err := OptionsFromEnvValidated()

	require.NoError(t, err)
	assert.Empty(t, opts.Name)
	assert.Empty(t, opts.Image)
	assert.Nil(t, opts.Command)
	assert.Empty(t, opts.DataDir)
	assert.Empty(t, opts.Resources.GPU)
}

func TestService_OptionsFromEnv_Ugly(t *testing.T) {
	t.Setenv("CORE_TIM_NAME", " \t\n ")
	t.Setenv("CORE_TIM_IMAGE", " \t ghcr.io/example/tim:latest \n")
	t.Setenv("CORE_TIM_COMMAND", " , first ,, second , ")
	t.Setenv("CORE_TIM_DATA_DIR", "\t /tmp/core-tim \n")
	t.Setenv("CORE_TIM_GPU", "\t device=0 \n")

	opts, err := OptionsFromEnvValidated()

	require.NoError(t, err)
	assert.Empty(t, opts.Name)
	assert.Equal(t, "ghcr.io/example/tim:latest", opts.Image)
	assert.Equal(t, []string{"first", "second"}, opts.Command)
	assert.Equal(t, "/tmp/core-tim", opts.DataDir)
	assert.Equal(t, "device=0", opts.Resources.GPU)
}

func TestService_OptionsFromEnv_RejectsLeadingDash(t *testing.T) {
	cases := []struct {
		name   string
		envKey string
		value  string
		want   string
	}{
		{name: "name", envKey: "CORE_TIM_NAME", value: "-rm", want: "name cannot start with -"},
		{name: "image", envKey: "CORE_TIM_IMAGE", value: "--privileged", want: "image cannot start with -"},
		{name: "gpu", envKey: "CORE_TIM_GPU", value: "-it", want: "gpu cannot start with -"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("CORE_TIM_NAME", "")
			t.Setenv("CORE_TIM_IMAGE", "")
			t.Setenv("CORE_TIM_COMMAND", "")
			t.Setenv("CORE_TIM_DATA_DIR", "")
			t.Setenv("CORE_TIM_GPU", "")
			t.Setenv(tc.envKey, tc.value)

			_, err := OptionsFromEnvValidated()

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestService_NewService_Good(t *testing.T) {
	svc, _ := newTestContainerService(t, TIMOptions{
		Name:  "normal-container",
		Image: "alpine:3.19",
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
	})

	assert.NotNil(t, svc.manager)
	state := svc.State()
	assert.Equal(t, "normal-container", state.Name)
	assert.Equal(t, "alpine:3.19", state.Image)
	assert.Equal(t, RuntimeDocker, state.Runtime)
	assert.Equal(t, "stopped", state.Status)
}

func TestService_NewService_Bad(t *testing.T) {
	svc, _ := newTestContainerService(t, TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	state := svc.State()
	assert.Equal(t, "coregui-tim", state.Name)
	assert.Equal(t, "ghcr.io/lthn/core/tim:latest", state.Image)
	assert.Equal(t, RuntimeNone, state.Runtime)
}

func TestService_NewService_Ugly(t *testing.T) {
	svc, _ := newTestContainerService(t, TIMOptions{
		Name:    "  worker.node  ",
		Image:   "  ghcr.io/example/tim:edge  ",
		Command: []string{"alpha", "beta"},
		DataDir: " /tmp/data ",
		Resources: TIMResources{
			CPUCores: 2,
			MemoryMB: 512,
			GPU:      "all",
		},
		Detect: func() ContainerRuntime {
			return RuntimePodman
		},
	})

	state := svc.State()
	assert.Equal(t, "worker.node", state.Name)
	assert.Equal(t, "ghcr.io/example/tim:edge", state.Image)
	assert.Equal(t, []string{"alpha", "beta"}, state.Command)
	assert.Equal(t, " /tmp/data ", state.DataDir)
	assert.Equal(t, TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"}, state.Resources)
	assert.Equal(t, RuntimePodman, state.Runtime)
}

func TestService_NewService_RejectsInvalidTIMOptions(t *testing.T) {
	cases := []struct {
		name    string
		options TIMOptions
		want    string
	}{
		{
			name:    "name leading dash",
			options: TIMOptions{Name: "-rm", Image: "alpine:3.19"},
			want:    "name cannot start with -",
		},
		{
			name:    "image leading dash",
			options: TIMOptions{Name: "normal-container", Image: "--privileged"},
			want:    "image cannot start with -",
		},
		{
			name: "gpu leading dash",
			options: TIMOptions{
				Name:      "normal-container",
				Image:     "alpine:3.19",
				Resources: TIMResources{GPU: "-it"},
			},
			want: "gpu cannot start with -",
		},
		{
			name:    "name starts with dot",
			options: TIMOptions{Name: ".hidden", Image: "alpine:3.19"},
			want:    "name cannot start with .",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := newInvalidTestContainerService(t, tc.options)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestService_OnStartup_Good(t *testing.T) {
	var calls []string
	svc, c := newTestContainerService(t, TIMOptions{
		Name:  "coregui-tim",
		Image: "ghcr.io/example/tim:latest",
		Command: []string{
			"sleep",
			"1",
		},
		Resources: TIMResources{CPUCores: 2, MemoryMB: 512, GPU: "all"},
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(_ context.Context, name string, args ...string) error {
			calls = append(calls, append([]string{name}, args...)...)
			return nil
		},
		Now: func() time.Time {
			return time.Unix(123, 0).UTC()
		},
	})

	runtime := c.Action("container.runtime.detect").Run(context.Background(), core.NewOptions())
	require.True(t, runtime.OK)
	assert.Equal(t, RuntimeDocker, runtime.Value)

	status := c.Action("tim.status").Run(context.Background(), core.NewOptions())
	require.True(t, status.OK)
	initial := status.Value.(TIMState)
	assert.Equal(t, "stopped", initial.Status)

	started := c.Action("tim.start").Run(context.Background(), core.NewOptions())
	require.True(t, started.OK)
	startState := started.Value.(TIMState)
	assert.Equal(t, "running", startState.Status)
	assert.Equal(t, time.Unix(123, 0).UTC(), startState.StartedAt)
	require.NotEmpty(t, calls)
	assert.Equal(t, "docker", calls[0])
	assert.Contains(t, calls, "run")
	assert.Contains(t, calls, "--name")
	assert.Contains(t, calls, "coregui-tim")
	assert.Contains(t, calls, "--cpus")
	assert.Contains(t, calls, "2")
	assert.Contains(t, calls, "--memory")
	assert.Contains(t, calls, "512m")
	assert.Contains(t, calls, "--gpus")
	assert.Contains(t, calls, "all")

	stopped := c.Action("tim.stop").Run(context.Background(), core.NewOptions())
	require.True(t, stopped.OK)
	stopState := stopped.Value.(TIMState)
	assert.Equal(t, "stopped", stopState.Status)
	assert.Equal(t, "coregui-tim", svc.State().Name)
}

func TestService_OnStartup_Bad(t *testing.T) {
	_, c := newTestContainerService(t, TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeNone
		},
	})

	result := c.Action("tim.start").Run(context.Background(), core.NewOptions())

	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "no supported container runtime detected")
}

func TestService_OnStartup_Ugly(t *testing.T) {
	_, c := newTestContainerService(t, TIMOptions{
		Detect: func() ContainerRuntime {
			return RuntimeDocker
		},
		Exec: func(context.Context, string, ...string) error {
			return errors.New("boom")
		},
	})

	result := c.Action("tim.start").Run(context.Background(), core.NewOptions())

	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "boom")

	status := c.Action("tim.status").Run(context.Background(), core.NewOptions())
	require.True(t, status.OK)
	assert.Equal(t, "error", status.Value.(TIMState).Status)
}
