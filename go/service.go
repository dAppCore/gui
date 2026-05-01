// SPDX-License-Identifier: EUPL-1.2

// Package gui is the canonical entry-point for the core/gui repo.
// One service per repo — `c.Service("gui")` returns this.
//
//	c, _ := core.New(core.WithName("gui", gui.NewService(gui.GuiConfig{})))

package gui

import (
	"context"

	core "dappco.re/go"
)

// GuiConfig configures the gui service.
//
// Usage example: `cfg := gui.GuiConfig{}`
type GuiConfig struct{}

// Service is the registerable handle for the gui repo.
//
// Usage example: `svc := core.MustServiceFor[*gui.Service](c, "gui")`
type Service struct {
	*core.ServiceRuntime[GuiConfig]
	registrations core.Once
}

// NewService returns the factory for c.Service() registration.
//
// Usage example: `c, _ := core.New(core.WithName("gui", gui.NewService(gui.GuiConfig{})))`
func NewService(config GuiConfig) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		return core.Ok(&Service{
			ServiceRuntime: core.NewServiceRuntime(c, config),
		})
	}
}

// Register builds the gui service with default GuiConfig — imperative
// alternative to NewService.
//
// Usage example: `r := gui.Register(c)`
func Register(c *core.Core) core.Result {
	return NewService(GuiConfig{})(c)
}

// OnStartup implements core.Startable. Idempotent via core.Once.
//
// Usage example: `r := svc.OnStartup(ctx)`
func (s *Service) OnStartup(context.Context) core.Result {
	if s == nil {
		return core.Ok(nil)
	}
	s.registrations.Do(func() {})
	return core.Ok(nil)
}

// OnShutdown implements core.Stoppable.
//
// Usage example: `r := svc.OnShutdown(ctx)`
func (s *Service) OnShutdown(context.Context) core.Result {
	return core.Ok(nil)
}
