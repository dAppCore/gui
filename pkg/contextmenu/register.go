// pkg/contextmenu/register.go
package contextmenu

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
// The returned function has the signature WithService requires: func(*Core) (any, error).
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			menus:          make(map[string]ContextMenuDef),
		}, nil
	}
}
