package window

import "forge.lthn.ai/core/go/pkg/core"

func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, nil
	}
}
