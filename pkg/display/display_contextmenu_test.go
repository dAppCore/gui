package display

import (
	"context"
	"sync"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/contextmenu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type wsContextMenuPlatform struct {
	mu    sync.Mutex
	menus map[string]contextmenu.ContextMenuDef
}

func newWSContextMenuPlatform() *wsContextMenuPlatform {
	return &wsContextMenuPlatform{
		menus: make(map[string]contextmenu.ContextMenuDef),
	}
}

func (m *wsContextMenuPlatform) Add(name string, menu contextmenu.ContextMenuDef, _ func(string, string, string)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.menus[name] = menu
	return nil
}

func (m *wsContextMenuPlatform) Remove(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.menus, name)
	return nil
}

func (m *wsContextMenuPlatform) Get(name string) (*contextmenu.ContextMenuDef, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	menu, ok := m.menus[name]
	if !ok {
		return nil, false
	}
	return &menu, true
}

func (m *wsContextMenuPlatform) GetAll() map[string]contextmenu.ContextMenuDef {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[string]contextmenu.ContextMenuDef, len(m.menus))
	for name, menu := range m.menus {
		out[name] = menu
	}
	return out
}

func newDisplayWithContextMenu(t *testing.T, platform *wsContextMenuPlatform) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(contextmenu.Register(platform)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	return core.MustServiceFor[*Service](c, "display"), c
}

func TestDisplay_handleWSMessage_ContextMenuAdd_MissingMenu(t *testing.T) {
	platform := newWSContextMenuPlatform()
	svc, _ := newDisplayWithContextMenu(t, platform)

	result := svc.handleWSMessage(WSMessage{
		Action: "contextmenu:add",
		Data: map[string]any{
			"name": "menu",
		},
	})

	require.False(t, result.OK)
	assert.Contains(t, result.Value.(error).Error(), `missing required field "menu"`)
	_, ok := platform.Get("menu")
	assert.False(t, ok)
}
