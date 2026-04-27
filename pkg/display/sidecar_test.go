package display

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/deno"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stderr
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	defer func() {
		os.Stderr = original
		_ = reader.Close()
		_ = writer.Close()
	}()

	os.Stderr = writer
	fn()
	os.Stderr = original
	require.NoError(t, writer.Close())

	var output bytes.Buffer
	_, err = output.ReadFrom(reader)
	require.NoError(t, err)
	return output.String()
}

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

func TestSidecar_ValidateArgs_Good(t *testing.T) {
	output := captureStderr(t, func() {
		assert.NoError(t, validateSidecarArgs(splitCommandArgs(""), nil))
		assert.NoError(t, validateSidecarArgs(splitCommandArgs("   "), nil))
	})

	assert.Empty(t, strings.TrimSpace(output))
}

func TestSidecar_LaunchOptions_Good_EmptyEnv(t *testing.T) {
	t.Setenv("CORE_DENO_ARGS", "")
	t.Setenv("CORE_DENO_BINARY", "")
	t.Setenv("CORE_DENO_DIR", "")

	var options deno.Options
	output := captureStderr(t, func() {
		var err error
		options, err = sidecarLaunchOptions(nil)
		require.NoError(t, err)
	})

	assert.Nil(t, options.Args)
	assert.Empty(t, options.Binary)
	assert.Empty(t, options.Dir)
	assert.Empty(t, strings.TrimSpace(output))
}

func TestSidecar_ValidateArgs_Good_Unstable(t *testing.T) {
	output := captureStderr(t, func() {
		assert.NoError(t, validateSidecarArgs(splitCommandArgs("--unstable"), nil))
	})

	assert.Empty(t, strings.TrimSpace(output))
}

func TestSidecar_ValidateArgs_Bad_PermissionFlags(t *testing.T) {
	tests := []struct {
		name string
		args string
		flag string
	}{
		{name: "allow-all", args: "run --allow-all attacker.ts", flag: "--allow-all"},
		{name: "allow-all-short", args: "-A", flag: "-A"},
		{name: "multiple-allow-flags", args: "--allow-net --allow-read=/", flag: "--allow-net"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CORE_DENO_ALLOW_PERMISSIONS", "")

			err := validateSidecarArgs(splitCommandArgs(tt.args), nil)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.flag)
			assert.Contains(t, err.Error(), "deno sandbox is being weakened")
		})
	}
}

func TestSidecar_ValidateArgs_Good_OverrideWarns(t *testing.T) {
	t.Setenv("CORE_DENO_ALLOW_PERMISSIONS", "true")

	output := captureStderr(t, func() {
		assert.NoError(t, validateSidecarArgs(splitCommandArgs("run --allow-all attacker.ts"), nil))
	})

	assert.Contains(t, output, "CORE_DENO_ARGS contains permission flag --allow-all")
	assert.Contains(t, output, "deno sandbox is being weakened")
}

func TestSidecar_StartAction_Bad_RefusesPermissionArgs(t *testing.T) {
	t.Setenv("CORE_DENO_ARGS", "run --allow-all attacker.ts")

	svc, c := newTestDisplayService(t)
	result := c.Action("display.sidecar.start").Run(context.Background(), core.Options{})

	require.False(t, result.OK)
	err, ok := result.Value.(error)
	require.True(t, ok)
	assert.Contains(t, err.Error(), "--allow-all")
	assert.Nil(t, svc.sidecar)
}

func TestSidecar_ValidateBinary_Good(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "deno")
	require.NoError(t, os.WriteFile(binary, []byte("#!/bin/sh\n"), 0o755))
	expected, err := filepath.EvalSymlinks(binary)
	require.NoError(t, err)

	actual, err := validateSidecarBinary(binary)

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSidecar_ValidateBinary_Bad(t *testing.T) {
	customBinary := filepath.Join(t.TempDir(), "deno-custom")
	require.NoError(t, os.WriteFile(customBinary, []byte("#!/bin/sh\n"), 0o755))

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "relative", value: "deno", want: "absolute"},
		{name: "missing", value: filepath.Join(t.TempDir(), "deno"), want: "does not exist"},
		{name: "custom-name", value: customBinary, want: "named deno"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateSidecarBinary(tt.value)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestSidecar_ValidateDir_Good(t *testing.T) {
	dir := canonicalTempDir(t)

	actual, err := validateSidecarDir(dir)

	require.NoError(t, err)
	assert.Equal(t, dir, actual)
}

func TestSidecar_ValidateDir_Bad(t *testing.T) {
	base := canonicalTempDir(t)
	child := filepath.Join(base, "child")
	require.NoError(t, os.Mkdir(child, 0o755))
	file := filepath.Join(base, "not-a-dir")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o644))
	target := filepath.Join(base, "target")
	require.NoError(t, os.Mkdir(target, 0o755))
	link := filepath.Join(base, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "parent-component", value: child + string(filepath.Separator) + ".." + string(filepath.Separator) + "child", want: ".."},
		{name: "file", value: file, want: "directory"},
		{name: "symlink", value: link, want: "symlink"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateSidecarDir(tt.value)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.want)
		})
	}
}

func canonicalTempDir(t *testing.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return dir
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
