package deno

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type Options struct {
	Binary string
	Args   []string
	Dir    string
	Env    []string
}

type Status struct {
	Running bool   `json:"running"`
	PID     int    `json:"pid,omitempty"`
	Binary  string `json:"binary,omitempty"`
}

type Manager struct {
	options Options
	mu      sync.Mutex
	cmd     *exec.Cmd
}

func New(options Options) *Manager {
	if strings.TrimSpace(options.Binary) == "" {
		options.Binary = "deno"
	}
	return &Manager{options: options}
}

func (m *Manager) Start(ctx context.Context) (Status, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cmd != nil && m.cmd.Process != nil {
		return m.statusLocked(), nil
	}

	cmd := exec.CommandContext(ctx, m.options.Binary, m.options.Args...)
	cmd.Dir = m.options.Dir
	cmd.Env = append(os.Environ(), m.options.Env...)
	if err := cmd.Start(); err != nil {
		return Status{}, err
	}
	m.cmd = cmd
	return m.statusLocked(), nil
}

func (m *Manager) Stop(context.Context) (Status, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cmd == nil || m.cmd.Process == nil {
		return Status{}, nil
	}
	err := m.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		return m.statusLocked(), err
	}
	m.cmd = nil
	return Status{}, nil
}

func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.statusLocked()
}

func (m *Manager) statusLocked() Status {
	status := Status{
		Binary: m.options.Binary,
	}
	if m.cmd != nil && m.cmd.Process != nil {
		status.Running = true
		status.PID = m.cmd.Process.Pid
	}
	return status
}
