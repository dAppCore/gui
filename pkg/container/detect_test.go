package container

import (
	"errors"
	"path/filepath"
	"testing"

	coreio "dappco.re/go/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectWithEnvironment_PrefersAppleContainersOnMacOS26(t *testing.T) {
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

	assert.Equal(t, RuntimeApple, runtime)
}

func TestDetectWithEnvironment_FallsBackToDockerWhenAppleUnavailable(t *testing.T) {
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

	assert.Equal(t, RuntimeDocker, runtime)
}

func TestDetectWithEnvironment_UsesDockerOnNonMacHosts(t *testing.T) {
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

	assert.Equal(t, RuntimeDocker, runtime)
}

func TestDetectWithEnvironment_UsesPodmanWhenDockerMissing(t *testing.T) {
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

	assert.Equal(t, RuntimePodman, runtime)
}

func TestDetectWithEnvironment_ReturnsNoneWhenNoRuntimeIsAvailable(t *testing.T) {
	runtime := DetectWithEnvironment(DetectEnvironment{
		GOOS:           "linux",
		ProductVersion: "",
		LookPath: func(string) (string, error) {
			return "", errors.New("not found")
		},
	})

	assert.Equal(t, RuntimeNone, runtime)
}

func TestMajorVersion(t *testing.T) {
	assert.Equal(t, 26, majorVersion("26.0"))
	assert.Equal(t, 0, majorVersion("bogus"))
	assert.Equal(t, 0, majorVersion(""))
}

func TestDetect_Good(t *testing.T) {
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

	assert.Equal(t, RuntimeApple, runtime)
}

func TestDetect_Bad(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, binDir, "sw_vers", "#!/bin/sh\nprintf '25.0\\n'\n")
	writeExecutable(t, binDir, "docker", "#!/bin/sh\nexit 0\n")
	t.Setenv("PATH", binDir)

	assert.Equal(t, RuntimeDocker, Detect())
}

func TestDetect_Ugly(t *testing.T) {
	binDir := t.TempDir()
	writeExecutable(t, binDir, "sw_vers", "#!/bin/sh\nprintf 'not-a-version\\n'\n")
	t.Setenv("PATH", binDir)

	assert.Equal(t, RuntimeNone, Detect())
}

func writeExecutable(t *testing.T, dir, name, script string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, coreio.Local.WriteMode(path, script, 0o755))
	return path
}
