// pkg/events/register.go
package events

import "forge.lthn.ai/core/go/pkg/core"

// Register(p) binds the events service to a Core instance.
// core.WithService(events.Register(wailsEventManager))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			counts:         make(map[string]int),
		}, nil
	}
}
