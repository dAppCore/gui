package keybinding

import "forge.lthn.ai/core/go/pkg/core"

// Register(p) binds the keybinding service to a Core instance.
// core.WithService(keybinding.Register(wailsKeybinding))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime:     core.NewServiceRuntime[Options](c, Options{}),
			platform:           p,
			registeredBindings: make(map[string]BindingInfo),
		}, nil
	}
}
