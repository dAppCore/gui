package lifecycle

import "forge.lthn.ai/core/go/pkg/core"

// Register(p) binds the lifecycle service to a Core instance.
// core.WithService(lifecycle.Register(wailsLifecycle))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}
