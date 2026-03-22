// pkg/lifecycle/service.go
package lifecycle

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the lifecycle service.
type Options struct{}

// Service is a core.Service that registers platform lifecycle callbacks
// and broadcasts corresponding IPC Actions. It implements both Startable
// and Stoppable: OnStartup registers all callbacks, OnShutdown cancels them.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	cancels  []func()
}

// OnStartup registers a platform callback for each EventType and for file-open.
// Each callback broadcasts the corresponding Action via s.Core().ACTION().
func (s *Service) OnStartup(ctx context.Context) error {
	// Register fire-and-forget event callbacks
	eventActions := map[EventType]func(){
		EventApplicationStarted: func() { _ = s.Core().ACTION(ActionApplicationStarted{}) },
		EventWillTerminate:      func() { _ = s.Core().ACTION(ActionWillTerminate{}) },
		EventDidBecomeActive:    func() { _ = s.Core().ACTION(ActionDidBecomeActive{}) },
		EventDidResignActive:    func() { _ = s.Core().ACTION(ActionDidResignActive{}) },
		EventPowerStatusChanged: func() { _ = s.Core().ACTION(ActionPowerStatusChanged{}) },
		EventSystemSuspend:      func() { _ = s.Core().ACTION(ActionSystemSuspend{}) },
		EventSystemResume:       func() { _ = s.Core().ACTION(ActionSystemResume{}) },
	}

	for eventType, handler := range eventActions {
		cancel := s.platform.OnApplicationEvent(eventType, handler)
		s.cancels = append(s.cancels, cancel)
	}

	// Register file-open callback (carries data)
	cancel := s.platform.OnOpenedWithFile(func(path string) {
		_ = s.Core().ACTION(ActionOpenedWithFile{Path: path})
	})
	s.cancels = append(s.cancels, cancel)

	return nil
}

// OnShutdown cancels all registered platform callbacks.
func (s *Service) OnShutdown(ctx context.Context) error {
	for _, cancel := range s.cancels {
		cancel()
	}
	s.cancels = nil
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
// Lifecycle events are all outbound (platform -> IPC) so there is nothing to handle here.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}
