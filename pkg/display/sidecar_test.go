package display

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	core "dappco.re/go"
	"dappco.re/go/gui/pkg/deno"
)

func captureStderr(t *core.T, fn func()) string {
	t.Helper()

	original := os.Stderr
	reader, writer, err := os.Pipe()
	core.RequireNoError(t, err)
	defer func() {
		os.Stderr = original
		_ = reader.Close()
		_ = writer.Close()
	}()

	os.Stderr = writer
	fn()
	os.Stderr = original
	core.RequireNoError(t, writer.Close())

	var output bytes.Buffer
	_, err = output.ReadFrom(reader)
	core.RequireNoError(t, err)
	return output.String()
}

func TestSidecar_SplitCommandArgs_Good(t *core.T) {
	core.AssertEqual(t, []string{"--import-map", "map.json", "--watch"}, splitCommandArgs("--import-map map.json --watch"))
	observedType := core.Sprintf("%T", splitCommandArgs("--import-map map.json --watch"))
	core.AssertNotEmpty(t, observedType)
}

func TestSidecar_SplitCommandArgs_Bad(t *core.T) {
	core.AssertNil(t, splitCommandArgs(""))
	core.AssertNil(t, splitCommandArgs("   "))
	core.AssertNotEmpty(t, core.Sprintf("%T", splitCommandArgs("")))
}

func TestSidecar_SplitCommandArgs_Ugly(t *core.T) {
	core.AssertEqual(t, []string{"--flag", "--another", "value"}, splitCommandArgs("\t--flag\n--another   value\t"))
	observedType := core.Sprintf("%T", splitCommandArgs("\t--flag\n--another   value\t"))
	core.AssertNotEmpty(t, observedType)
}

func TestSidecar_ValidateArgs_Good(t *core.T) {
	output := captureStderr(t, func() {
		core.AssertNoError(t, validateSidecarArgs(splitCommandArgs(""), nil))
		core.AssertNoError(t, validateSidecarArgs(splitCommandArgs("   "), nil))
	})

	core.AssertEmpty(t, strings.TrimSpace(output))
}

func TestSidecar_LaunchOptions_Good_EmptyEnv(t *core.T) {
	t.Setenv("CORE_DENO_ARGS", "")
	t.Setenv("CORE_DENO_BINARY", "")
	t.Setenv("CORE_DENO_DIR", "")

	var options deno.Options
	output := captureStderr(t, func() {
		var err error
		options, err = sidecarLaunchOptions(nil)
		core.RequireNoError(t, err)
	})

	core.AssertNil(t, options.Args)
	core.AssertEmpty(t, options.Binary)
	core.AssertEmpty(t, options.Dir)
	core.AssertEmpty(t, strings.TrimSpace(output))
}

func TestSidecar_ValidateArgs_Good_Unstable(t *core.T) {
	output := captureStderr(t, func() {
		core.AssertNoError(t, validateSidecarArgs(splitCommandArgs("--unstable"), nil))
	})

	core.AssertEmpty(t, strings.TrimSpace(output))
}

func TestSidecar_ValidateArgs_Bad_PermissionFlags(t *core.T) {
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
		t.Run(tt.name, func(t *core.T) {
			t.Setenv("CORE_DENO_ALLOW_PERMISSIONS", "")

			err := validateSidecarArgs(splitCommandArgs(tt.args), nil)

			core.AssertError(t, err)
			core.AssertContains(t, err.Error(), tt.flag)
			core.AssertContains(t, err.Error(), "deno sandbox is being weakened")
		})
	}
}

func TestSidecar_ValidateArgs_Good_OverrideWarns(t *core.T) {
	t.Setenv("CORE_DENO_ALLOW_PERMISSIONS", "true")

	output := captureStderr(t, func() {
		core.AssertNoError(t, validateSidecarArgs(splitCommandArgs("run --allow-all attacker.ts"), nil))
	})

	core.AssertContains(t, output, "CORE_DENO_ARGS contains permission flag --allow-all")
	core.AssertContains(t, output, "deno sandbox is being weakened")
}

func TestSidecar_StartAction_Bad_RefusesPermissionArgs(t *core.T) {
	t.Setenv("CORE_DENO_ARGS", "run --allow-all attacker.ts")

	svc, c := newTestDisplayService(t)
	result := c.Action("display.sidecar.start").Run(context.Background(), core.Options{})

	core.AssertFalse(t, result.OK)
	err, ok := result.Value.(error)
	core.RequireTrue(t, ok)
	core.AssertContains(t, err.Error(), "--allow-all")
	core.AssertNil(t, svc.sidecar)
}

func TestSidecar_ValidateBinary_Good(t *core.T) {
	binary := filepath.Join(t.TempDir(), "deno")
	core.RequireNoError(t, os.WriteFile(binary, []byte("#!/bin/sh\n"), 0o755))
	expected, err := filepath.EvalSymlinks(binary)
	core.RequireNoError(t, err)

	actual, err := validateSidecarBinary(binary)

	core.RequireNoError(t, err)
	core.AssertEqual(t, expected, actual)
}

func TestSidecar_ValidateBinary_Bad(t *core.T) {
	customBinary := filepath.Join(t.TempDir(), "deno-custom")
	core.RequireNoError(t, os.WriteFile(customBinary, []byte("#!/bin/sh\n"), 0o755))

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
		t.Run(tt.name, func(t *core.T) {
			_, err := validateSidecarBinary(tt.value)

			core.AssertError(t, err)
			core.AssertContains(t, err.Error(), tt.want)
		})
	}
}

func TestSidecar_ValidateDir_Good(t *core.T) {
	dir := canonicalTempDir(t)

	actual, err := validateSidecarDir(dir)

	core.RequireNoError(t, err)
	core.AssertEqual(t, dir, actual)
}

func TestSidecar_ValidateDir_Bad(t *core.T) {
	base := canonicalTempDir(t)
	child := filepath.Join(base, "child")
	core.RequireNoError(t, os.Mkdir(child, 0o755))
	file := filepath.Join(base, "not-a-dir")
	core.RequireNoError(t, os.WriteFile(file, []byte("x"), 0o644))
	target := filepath.Join(base, "target")
	core.RequireNoError(t, os.Mkdir(target, 0o755))
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
		t.Run(tt.name, func(t *core.T) {
			_, err := validateSidecarDir(tt.value)

			core.AssertError(t, err)
			core.AssertContains(t, err.Error(), tt.want)
		})
	}
}

func canonicalTempDir(t *core.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	core.RequireNoError(t, err)
	return dir
}

func TestSidecar_EnsureSidecar_Good(t *core.T) {
	t.Setenv("CORE_DENO_BINARY", "/usr/local/bin/deno-custom")
	t.Setenv("CORE_DENO_DIR", "/tmp/core-deno")
	t.Setenv("CORE_DENO_ARGS", "--import-map map.json")

	svc := &Service{}
	manager := svc.ensureSidecar()

	core.AssertNotNil(t, manager)
	status := manager.Status()
	core.AssertEqual(t, "/usr/local/bin/deno-custom", status.Binary)
	core.AssertFalse(t, status.Running)
}

func TestSidecar_EnsureSidecar_Bad(t *core.T) {
	svc := &Service{sidecar: deno.New(deno.Options{Binary: "custom-deno"})}

	manager := svc.ensureSidecar()

	core.AssertSame(t, svc.sidecar, manager)
	core.AssertEqual(t, "custom-deno", manager.Status().Binary)
}

func TestSidecar_EnsureSidecar_Ugly(t *core.T) {
	t.Setenv("CORE_DENO_BINARY", "   ")
	t.Setenv("CORE_DENO_DIR", "")
	t.Setenv("CORE_DENO_ARGS", "   ")

	svc := &Service{}
	manager := svc.ensureSidecar()

	core.AssertNotNil(t, manager)
	core.AssertEqual(t, "deno", manager.Status().Binary)
}

func TestSidecar_RegisterActions_StartFailureClearsSidecar(t *core.T) {
	t.Setenv("CORE_DENO_ENABLE", "1")
	t.Setenv("CORE_DENO_BINARY", "/definitely/not/a/real/deno")

	c := core.New(core.WithServiceLock())
	svc := &Service{ServiceRuntime: core.NewServiceRuntime(c, Options{})}

	svc.registerSidecarActions()

	core.AssertNil(t, svc.sidecar)
}

func TestSidecar_StatusAction_Good(t *core.T) {
	t.Setenv("CORE_DENO_BINARY", "/opt/core/deno")

	_, c := newTestDisplayService(t)
	result := c.Action("display.sidecar.status").Run(context.Background(), core.Options{})

	core.RequireTrue(t, result.OK)
	status, ok := result.Value.(deno.Status)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "/opt/core/deno", status.Binary)
	core.AssertFalse(t, status.Running)
}
