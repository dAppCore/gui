package container

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type TIMOptions struct {
	Name      string
	Image     string
	Command   []string
	WorkDir   string
	Env       map[string]string
	DataDir   string
	Runtime   ContainerRuntime
	Detect    func() ContainerRuntime
	Exec      func(context.Context, string, ...string) error
	Now       func() time.Time
	Resources TIMResources
}

type TIMResources struct {
	CPUCores int    `json:"cpu_cores,omitempty"`
	MemoryMB int    `json:"memory_mb,omitempty"`
	GPU      string `json:"gpu,omitempty"`
}

type TIMState struct {
	Name      string           `json:"name"`
	Image     string           `json:"image"`
	Runtime   ContainerRuntime `json:"runtime"`
	Status    string           `json:"status"`
	StartedAt time.Time        `json:"started_at,omitempty"`
	Command   []string         `json:"command,omitempty"`
	DataDir   string           `json:"data_dir,omitempty"`
	Resources TIMResources     `json:"resources,omitempty"`
}

type TIMManager struct {
	options TIMOptions
	mu      sync.Mutex
	state   TIMState
}

func NewTIMManager(options TIMOptions) *TIMManager {
	if strings.TrimSpace(options.Name) == "" {
		options.Name = "coregui-tim"
	}
	if strings.TrimSpace(options.Image) == "" {
		options.Image = "ghcr.io/lthn/core/tim:latest"
	}
	if options.Detect == nil {
		options.Detect = Detect
	}
	if options.Exec == nil {
		options.Exec = func(ctx context.Context, name string, args ...string) error {
			cmd := exec.CommandContext(ctx, name, args...)
			return cmd.Run()
		}
	}
	if options.Now == nil {
		options.Now = time.Now
	}
	return &TIMManager{
		options: options,
		state: TIMState{
			Name:      options.Name,
			Image:     options.Image,
			Runtime:   coalesceRuntime(options.Runtime, options.Detect()),
			Status:    "stopped",
			Command:   append([]string(nil), options.Command...),
			DataDir:   options.DataDir,
			Resources: options.Resources,
		},
	}
}

func (m *TIMManager) State() TIMState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

func (m *TIMManager) Start(ctx context.Context) (TIMState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	runtime := coalesceRuntime(m.options.Runtime, m.options.Detect())
	m.state.Runtime = runtime
	if runtime == RuntimeNone {
		return m.state, fmt.Errorf("no supported container runtime detected")
	}

	command, args := m.runtimeCommand(runtime, "run")
	if err := m.options.Exec(ctx, command, args...); err != nil {
		m.state.Status = "error"
		return m.state, err
	}

	m.state.Status = "running"
	m.state.StartedAt = m.options.Now()
	return m.state, nil
}

func (m *TIMManager) Stop(ctx context.Context) (TIMState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state.Runtime == RuntimeNone {
		m.state.Status = "stopped"
		return m.state, nil
	}
	command, args := m.runtimeCommand(m.state.Runtime, "stop")
	if err := m.options.Exec(ctx, command, args...); err != nil {
		m.state.Status = "error"
		return m.state, err
	}
	m.state.Status = "stopped"
	return m.state, nil
}

func (m *TIMManager) runtimeCommand(runtime ContainerRuntime, verb string) (string, []string) {
	name := m.options.Name
	image := m.options.Image
	switch runtime {
	case RuntimeApple:
		if verb == "run" {
			args := []string{"run", "--name", name, image}
			args = append(args, m.options.Command...)
			return "container", args
		}
		return "container", []string{"stop", name}
	case RuntimePodman:
		if verb == "run" {
			args := []string{"run", "-d", "--replace", "--name", name}
			args = append(args, resourceArgs(m.options.Resources)...)
			args = append(args, image)
			args = append(args, m.options.Command...)
			return "podman", args
		}
		return "podman", []string{"stop", name}
	default:
		if verb == "run" {
			args := []string{"run", "-d", "--rm", "--name", name}
			args = append(args, resourceArgs(m.options.Resources)...)
			args = append(args, image)
			args = append(args, m.options.Command...)
			return "docker", args
		}
		return "docker", []string{"stop", name}
	}
}

func resourceArgs(resources TIMResources) []string {
	var args []string
	if resources.CPUCores > 0 {
		args = append(args, "--cpus", fmt.Sprintf("%d", resources.CPUCores))
	}
	if resources.MemoryMB > 0 {
		args = append(args, "--memory", fmt.Sprintf("%dm", resources.MemoryMB))
	}
	if strings.TrimSpace(resources.GPU) != "" {
		args = append(args, "--gpus", resources.GPU)
	}
	return args
}

func coalesceRuntime(values ...ContainerRuntime) ContainerRuntime {
	for _, value := range values {
		if value != RuntimeNone {
			return value
		}
	}
	return RuntimeNone
}
