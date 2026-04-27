package container

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ContainerRuntime describes the preferred isolated workload runtime for CoreGUI.
//
//	runtime := container.Detect()
//	if runtime == container.RuntimeApple { /* prefer Apple Containers */ }
type ContainerRuntime string

const (
	RuntimeNone   ContainerRuntime = ""
	RuntimeApple  ContainerRuntime = "apple"
	RuntimeDocker ContainerRuntime = "docker"
	RuntimePodman ContainerRuntime = "podman"
)

// DetectEnvironment controls runtime detection for tests and other callers.
//
//	runtime := container.DetectWithEnvironment(container.DetectEnvironment{
//	    GOOS:           "darwin",
//	    ProductVersion: "26.0",
//	})
type DetectEnvironment struct {
	GOOS           string
	ProductVersion string
	LookPath       func(file string) (string, error)
}

// Detect prefers Apple Containers on macOS 26+, then Docker, then Podman.
//
//	runtime := container.Detect()
func Detect() ContainerRuntime {
	environment := DetectEnvironment{
		GOOS:           runtime.GOOS,
		ProductVersion: "",
		LookPath:       exec.LookPath,
	}
	if runtime.GOOS == "darwin" {
		environment.ProductVersion = productVersion()
	}
	return DetectWithEnvironment(environment)
}

// DetectWithEnvironment applies the RFC runtime ordering using an explicit environment.
//
//	runtime := container.DetectWithEnvironment(env)
func DetectWithEnvironment(environment DetectEnvironment) ContainerRuntime {
	lookPath := environment.LookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}

	goos := strings.ToLower(strings.TrimSpace(environment.GOOS))
	if goos == "darwin" && majorVersion(environment.ProductVersion) >= 26 {
		if hasBinary(lookPath, "container") || hasBinary(lookPath, "apple-container") || hasBinary(lookPath, "containerctl") {
			return RuntimeApple
		}
	}
	if hasBinary(lookPath, "docker") {
		return RuntimeDocker
	}
	if hasBinary(lookPath, "podman") {
		return RuntimePodman
	}
	return RuntimeNone
}

func hasBinary(lookPath func(string) (string, error), binary string) bool {
	if strings.TrimSpace(binary) == "" {
		return false
	}
	_, err := lookPath(binary)
	return err == nil
}

func majorVersion(productVersion string) int {
	productVersion = strings.TrimSpace(productVersion)
	if productVersion == "" {
		return 0
	}
	major, _, _ := strings.Cut(productVersion, ".")
	value, err := strconv.Atoi(major)
	if err != nil {
		return 0
	}
	return value
}

func productVersion() string {
	output, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
