package dock

import "forge.lthn.ai/core/go/pkg/core"

// Register(p) binds the dock service to a Core instance.
// core.WithService(dock.Register(wailsDock))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}
