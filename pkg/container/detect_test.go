package container

import (
	core "dappco.re/go"
	"errors"
	"path/filepath"

	coreio "dappco.re/go/io"
)

func TestDetectWithEnvironment_PrefersAppleContainersOnMacOS26(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "darwin",
		ProductVersion: "26.0",
		LookPath: func(file string) (string, error) {
			if file == "container" {
				return "/usr/bin/container", nil
			}
			return "", errors.New("not found")
		},
	})

	core.AssertEqual(t, RuntimeApple, runtime)
}

func TestDetectWithEnvironment_FallsBackToDockerWhenAppleUnavailable(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "darwin",
		ProductVersion: "26.1",
		LookPath: func(file string) (string, error) {
			if file == "docker" {
				return "/usr/local/bin/docker", nil
			}
			return "", errors.New("not found")
		},
	})

	core.AssertEqual(t, RuntimeDocker, runtime)
}

func TestDetectWithEnvironment_UsesDockerOnNonMacHosts(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "linux",
		ProductVersion: "",
		LookPath: func(file string) (string, error) {
			if file == "docker" {
				return "/usr/bin/docker", nil
			}
			return "", errors.New("not found")
		},
	})

	core.AssertEqual(t, RuntimeDocker, runtime)
}

func TestDetectWithEnvironment_UsesPodmanWhenDockerMissing(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "linux",
		ProductVersion: "",
		LookPath: func(file string) (string, error) {
			if file == "podman" {
				return "/usr/bin/podman", nil
			}
			return "", errors.New("not found")
		},
	})

	core.AssertEqual(t, RuntimePodman, runtime)
}

func TestDetectWithEnvironment_ReturnsNoneWhenNoRuntimeIsAvailable(t *core.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "linux",
		ProductVersion: "",
		LookPath: func(string) (string, error) {
			return "", errors.New("not found")
		},
	})

	core.AssertEqual(t, RuntimeNone, runtime)
}

func TestMajorVersion(t *core.T) {
	core.AssertEqual(t, 26, majorVersion("26.0"))
	core.AssertEqual(t, 0, majorVersion("bogus"))
	core.AssertEqual(t, 0, majorVersion(""))
}

func TestDetect_Good(t *core.T) {
	binDir := t.TempDir()
	containerPath := writeExecutable(t, binDir, "container", "#!/bin/sh\nexit 0\n")

	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "darwin",
		ProductVersion: "26.0",
		LookPath: func(file string) (string, error) {
			if file == "container" {
				return containerPath, nil
			}
			return "", errors.New("not found")
		},
	})

	core.AssertEqual(t, RuntimeApple, runtime)
}

func TestDetect_Bad(t *core.T) {
	binDir := t.TempDir()
	writeExecutable(t, binDir, "sw_vers", "#!/bin/sh\nprintf '25.0\\n'\n")
	writeExecutable(t, binDir, "docker", "#!/bin/sh\nexit 0\n")
	t.Setenv("PATH", binDir)

	core.AssertEqual(t, RuntimeDocker, Detect())
}

func TestDetect_Ugly(t *core.T) {
	binDir := t.TempDir()
	writeExecutable(t, binDir, "sw_vers", "#!/bin/sh\nprintf 'not-a-version\\n'\n")
	t.Setenv("PATH", binDir)

	core.AssertEqual(t, RuntimeNone, Detect())
}

func writeExecutable(t *core.T, dir, name, script string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	core.RequireNoError(t, coreio.Local.WriteMode(path, script, 0o755))
	return path
}
