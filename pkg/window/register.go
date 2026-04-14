// pkg/window/register.go
package window

import "forge.lthn.ai/core/go/pkg/core"

// Register(p) binds the window service to a Core instance.
// core.WithService(window.Register(window.NewWailsPlatform(app)))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, nil
	}
}
