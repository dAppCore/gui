package lifecycle

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	cancels  []func()
}

func (s *Service) OnStartup(ctx context.Context) error {
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

	cancel := s.platform.OnOpenedWithFile(func(path string) {
		_ = s.Core().ACTION(ActionOpenedWithFile{Path: path})
	})
	s.cancels = append(s.cancels, cancel)

	return nil
}

func (s *Service) OnShutdown(ctx context.Context) error {
	for _, cancel := range s.cancels {
		cancel()
	}
	s.cancels = nil
	return nil
}

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}
