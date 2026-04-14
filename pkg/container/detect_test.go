package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
