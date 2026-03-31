package contextmenu

import "forge.lthn.ai/core/go/pkg/core"

// Register(p) binds the context menu service to a Core instance.
// core.WithService(contextmenu.Register(wailsContextMenu))
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime:  core.NewServiceRuntime[Options](c, Options{}),
			platform:        p,
			registeredMenus: make(map[string]ContextMenuDef),
		}, nil
	}
}
