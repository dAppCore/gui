package application

import (
	"context"
	"io/fs"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v3/internal/capabilities"
)

var globalApplication *App

// AlphaAssets mirrors the default alpha asset configuration in the real Wails API.
var AlphaAssets = AssetOptions{
	Handler: http.NotFoundHandler(),
}

// Get returns the singleton application created via New, if any.
func Get() *App {
	return globalApplication
}

// New returns a singleton in-memory application stub.
func New(appOptions Options) *App {
	if globalApplication != nil {
		return globalApplication
	}
	mergeApplicationDefaults(&appOptions)
	globalApplication = newStubApp(appOptions)
	if appOptions.OnShutdown != nil {
		globalApplication.OnShutdown(appOptions.OnShutdown)
	}
	return globalApplication
}

// Config returns the application options used to build this stub app.
func (a *App) Config() Options {
	state := ensureAppCompatState(a)
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.options
}

// Context returns the app context.
func (a *App) Context() context.Context {
	state := ensureAppCompatState(a)
	state.mu.RLock()
	defer state.mu.RUnlock()
	if state.ctx == nil {
		return context.Background()
	}
	return state.ctx
}

// RegisterService records a service registration in the in-memory stub.
func (a *App) RegisterService(service Service) {
	state := ensureAppCompatState(a)
	state.mu.Lock()
	state.services = append(state.services, service)
	options := state.options
	options.Services = append(options.Services, service)
	state.options = options
	state.mu.Unlock()
}

// Capabilities returns an empty capability set for the stub runtime.
func (a *App) Capabilities() capabilities.Capabilities {
	return capabilities.Capabilities{}
}

// GetPID returns the current process ID.
func (a *App) GetPID() int {
	return os.Getpid()
}

// Run marks the stub app as running.
func (a *App) Run() error {
	state := ensureAppCompatState(a)
	state.mu.Lock()
	if state.running {
		state.mu.Unlock()
		return nil
	}
	state.running = true
	services := append([]Service(nil), state.services...)
	state.mu.Unlock()

	for _, service := range services {
		if err := a.startupService(service); err != nil {
			state.mu.Lock()
			state.running = false
			state.mu.Unlock()
			return err
		}
	}
	return nil
}

// OnShutdown registers a callback that can be triggered by Quit.
func (a *App) OnShutdown(f func()) {
	if f == nil {
		return
	}
	state := ensureAppCompatState(a)
	state.mu.Lock()
	state.shutdown = append(state.shutdown, f)
	state.mu.Unlock()
}

// SetIcon stores the application icon bytes.
func (a *App) SetIcon(icon []byte) {
	state := ensureAppCompatState(a)
	state.mu.Lock()
	state.icon = append([]byte(nil), icon...)
	state.mu.Unlock()
}

// Hide records that the application is hidden.
func (a *App) Hide() {
	state := ensureAppCompatState(a)
	state.mu.Lock()
	state.hidden = true
	state.mu.Unlock()
}

// Show records that the application is visible.
func (a *App) Show() {
	state := ensureAppCompatState(a)
	state.mu.Lock()
	state.hidden = false
	state.mu.Unlock()
}

// AssetFileServerFS serves an fs.FS through the standard library file server.
func AssetFileServerFS(assets fs.FS) http.Handler {
	return http.FileServer(http.FS(assets))
}

// BundledAssetFileServer is the stub equivalent of Wails' bundled asset server.
func BundledAssetFileServer(assets fs.FS) http.Handler {
	return AssetFileServerFS(assets)
}
