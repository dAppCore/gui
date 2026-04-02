// pkg/systray/register.go
package systray

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, nil
	}
}
