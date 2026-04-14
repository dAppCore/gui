package browser

import "forge.lthn.ai/core/go/pkg/core"

// Register(p) binds the browser service to a Core instance.
// core.WithService(browser.Register(wailsBrowser))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}
