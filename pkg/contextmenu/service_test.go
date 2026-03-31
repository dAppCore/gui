// pkg/contextmenu/service_test.go
package contextmenu

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlatform records Add/Remove calls and allows simulating clicks.
type mockPlatform struct {
	mu            sync.Mutex
	menus         map[string]ContextMenuDef
	clickHandlers map[string]func(menuName, actionID, data string)
	removed       []string
	addErr        error
	removeErr     error
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		menus:         make(map[string]ContextMenuDef),
		clickHandlers: make(map[string]func(menuName, actionID, data string)),
	}
}

func (m *mockPlatform) Add(name string, menu ContextMenuDef, onItemClick func(string, string, string)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.addErr != nil {
		return m.addErr
	}
	m.menus[name] = menu
	m.clickHandlers[name] = onItemClick
	return nil
}

func (m *mockPlatform) Remove(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.removeErr != nil {
		return m.removeErr
	}
	delete(m.menus, name)
	delete(m.clickHandlers, name)
	m.removed = append(m.removed, name)
	return nil
}

func (m *mockPlatform) Get(name string) (*ContextMenuDef, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	menu, ok := m.menus[name]
	if !ok {
		return nil, false
	}
	return &menu, true
}

func (m *mockPlatform) GetAll() map[string]ContextMenuDef {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[string]ContextMenuDef, len(m.menus))
	for k, v := range m.menus {
		out[k] = v
	}
	return out
}

// simulateClick simulates a context menu item click by calling the registered handler.
func (m *mockPlatform) simulateClick(menuName, actionID, data string) {
	m.mu.Lock()
	h, ok := m.clickHandlers[menuName]
	m.mu.Unlock()
	if ok {
		h(menuName, actionID, data)
	}
}

func newTestContextMenuService(t *testing.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "contextmenu")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	mp := newMockPlatform()
	svc, _ := newTestContextMenuService(t, mp)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.platform)
}

func TestTaskAdd_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskAdd{
		Name: "file-menu",
		Menu: ContextMenuDef{
			Name: "file-menu",
			Items: []MenuItemDef{
				{Label: "Open", ActionID: "open"},
				{Label: "Delete", ActionID: "delete"},
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify menu registered on platform
	_, ok := mp.Get("file-menu")
	assert.True(t, ok)
}

func TestTaskAdd_Good_ReplaceExisting(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Add initial menu
	_, _, _ = c.PERFORM(TaskAdd{
		Name: "ctx",
		Menu: ContextMenuDef{Name: "ctx", Items: []MenuItemDef{{Label: "A", ActionID: "a"}}},
	})

	// Replace with new menu
	_, handled, err := c.PERFORM(TaskAdd{
		Name: "ctx",
		Menu: ContextMenuDef{Name: "ctx", Items: []MenuItemDef{{Label: "B", ActionID: "b"}}},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify registry has new menu
	result, _, _ := c.QUERY(QueryGet{Name: "ctx"})
	def := result.(*ContextMenuDef)
	require.Len(t, def.Items, 1)
	assert.Equal(t, "B", def.Items[0].Label)
}

func TestTaskRemove_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Add then remove
	_, _, _ = c.PERFORM(TaskAdd{
		Name: "test",
		Menu: ContextMenuDef{Name: "test"},
	})
	_, handled, err := c.PERFORM(TaskRemove{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify removed from registry
	result, _, _ := c.QUERY(QueryGet{Name: "test"})
	assert.Nil(t, result)
}

func TestTaskRemove_Bad_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskRemove{Name: "nonexistent"})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrorMenuNotFound)
}

func TestQueryGet_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{
		Name: "my-menu",
		Menu: ContextMenuDef{
			Name:  "my-menu",
			Items: []MenuItemDef{{Label: "Edit", ActionID: "edit"}},
		},
	})

	result, handled, err := c.QUERY(QueryGet{Name: "my-menu"})
	require.NoError(t, err)
	assert.True(t, handled)
	def := result.(*ContextMenuDef)
	assert.Equal(t, "my-menu", def.Name)
	assert.Len(t, def.Items, 1)
}

func TestQueryGet_Good_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	result, handled, err := c.QUERY(QueryGet{Name: "missing"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
}

func TestQueryList_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Name: "a", Menu: ContextMenuDef{Name: "a"}})
	_, _, _ = c.PERFORM(TaskAdd{Name: "b", Menu: ContextMenuDef{Name: "b"}})

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.(map[string]ContextMenuDef)
	assert.Len(t, list, 2)
}

func TestQueryList_Good_Empty(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.(map[string]ContextMenuDef)
	assert.Len(t, list, 0)
}

func TestTaskAdd_Good_ClickBroadcast(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Capture broadcast actions
	var clicked ActionItemClicked
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionItemClicked); ok {
			mu.Lock()
			clicked = a
			mu.Unlock()
		}
		return nil
	})

	_, _, _ = c.PERFORM(TaskAdd{
		Name: "file-menu",
		Menu: ContextMenuDef{
			Name: "file-menu",
			Items: []MenuItemDef{
				{Label: "Open", ActionID: "open"},
			},
		},
	})

	// Simulate click via mock
	mp.simulateClick("file-menu", "open", "file-123")

	mu.Lock()
	assert.Equal(t, "file-menu", clicked.MenuName)
	assert.Equal(t, "open", clicked.ActionID)
	assert.Equal(t, "file-123", clicked.Data)
	mu.Unlock()
}

func TestTaskAdd_Good_SubmenuItems(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskAdd{
		Name: "nested",
		Menu: ContextMenuDef{
			Name: "nested",
			Items: []MenuItemDef{
				{Label: "File", Type: "submenu", Items: []MenuItemDef{
					{Label: "New", ActionID: "new"},
					{Label: "Open", ActionID: "open"},
				}},
				{Type: "separator"},
				{Label: "Quit", ActionID: "quit"},
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryGet{Name: "nested"})
	def := result.(*ContextMenuDef)
	assert.Len(t, def.Items, 3)
	assert.Len(t, def.Items[0].Items, 2) // submenu children
}

func TestQueryList_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryList{})
	assert.False(t, handled)
}

// --- TaskUpdate ---

func TestTaskUpdate_Good(t *testing.T) {
	// Update replaces items on an existing menu
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{
		Name: "edit-menu",
		Menu: ContextMenuDef{Name: "edit-menu", Items: []MenuItemDef{{Label: "Cut", ActionID: "cut"}}},
	})

	_, handled, err := c.PERFORM(TaskUpdate{
		Name: "edit-menu",
		Menu: ContextMenuDef{Name: "edit-menu", Items: []MenuItemDef{
			{Label: "Cut", ActionID: "cut"},
			{Label: "Copy", ActionID: "copy"},
		}},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryGet{Name: "edit-menu"})
	def := result.(*ContextMenuDef)
	assert.Len(t, def.Items, 2)
	assert.Equal(t, "Copy", def.Items[1].Label)
}

func TestTaskUpdate_Bad_NotFound(t *testing.T) {
	// Update on a non-existent menu returns ErrorMenuNotFound
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskUpdate{
		Name: "ghost",
		Menu: ContextMenuDef{Name: "ghost"},
	})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrorMenuNotFound)
}

func TestTaskUpdate_Ugly_PlatformRemoveError(t *testing.T) {
	// Platform Remove fails mid-update — error is propagated
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{
		Name: "tricky",
		Menu: ContextMenuDef{Name: "tricky"},
	})

	mp.mu.Lock()
	mp.removeErr = ErrorMenuNotFound // reuse sentinel as a platform-level error
	mp.mu.Unlock()

	_, handled, err := c.PERFORM(TaskUpdate{
		Name: "tricky",
		Menu: ContextMenuDef{Name: "tricky", Items: []MenuItemDef{{Label: "X", ActionID: "x"}}},
	})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskDestroy ---

func TestTaskDestroy_Good(t *testing.T) {
	// Destroy removes the menu and releases platform resources
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Name: "doomed", Menu: ContextMenuDef{Name: "doomed"}})

	_, handled, err := c.PERFORM(TaskDestroy{Name: "doomed"})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryGet{Name: "doomed"})
	assert.Nil(t, result)

	_, ok := mp.Get("doomed")
	assert.False(t, ok)
}

func TestTaskDestroy_Bad_NotFound(t *testing.T) {
	// Destroy on a non-existent menu returns ErrorMenuNotFound
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskDestroy{Name: "nonexistent"})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrorMenuNotFound)
}

func TestTaskDestroy_Ugly_PlatformError(t *testing.T) {
	// Platform Remove fails — error is propagated but service remains consistent
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Name: "frail", Menu: ContextMenuDef{Name: "frail"}})

	mp.mu.Lock()
	mp.removeErr = ErrorMenuNotFound
	mp.mu.Unlock()

	_, handled, err := c.PERFORM(TaskDestroy{Name: "frail"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- QueryGetAll ---

func TestQueryGetAll_Good(t *testing.T) {
	// QueryGetAll returns all registered menus (equivalent to QueryList)
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Name: "x", Menu: ContextMenuDef{Name: "x"}})
	_, _, _ = c.PERFORM(TaskAdd{Name: "y", Menu: ContextMenuDef{Name: "y"}})

	result, handled, err := c.QUERY(QueryGetAll{})
	require.NoError(t, err)
	assert.True(t, handled)
	all := result.(map[string]ContextMenuDef)
	assert.Len(t, all, 2)
	assert.Contains(t, all, "x")
	assert.Contains(t, all, "y")
}

func TestQueryGetAll_Bad_Empty(t *testing.T) {
	// QueryGetAll on an empty registry returns an empty map
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	result, handled, err := c.QUERY(QueryGetAll{})
	require.NoError(t, err)
	assert.True(t, handled)
	all := result.(map[string]ContextMenuDef)
	assert.Len(t, all, 0)
}

func TestQueryGetAll_Ugly_NoService(t *testing.T) {
	// No contextmenu service — query is unhandled
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryGetAll{})
	assert.False(t, handled)
}

// --- OnShutdown ---

func TestOnShutdown_Good_CleansUpMenus(t *testing.T) {
	// OnShutdown removes all registered menus from the platform
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Name: "alpha", Menu: ContextMenuDef{Name: "alpha"}})
	_, _, _ = c.PERFORM(TaskAdd{Name: "beta", Menu: ContextMenuDef{Name: "beta"}})

	require.NoError(t, c.ServiceShutdown(t.Context()))

	assert.Len(t, mp.menus, 0)
}

func TestOnShutdown_Bad_NothingRegistered(t *testing.T) {
	// OnShutdown with no menus — no-op, no error
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	assert.NoError(t, c.ServiceShutdown(t.Context()))
}

func TestOnShutdown_Ugly_PlatformRemoveErrors(t *testing.T) {
	// Platform Remove errors during shutdown are silently swallowed
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Name: "stubborn", Menu: ContextMenuDef{Name: "stubborn"}})

	mp.mu.Lock()
	mp.removeErr = ErrorMenuNotFound
	mp.mu.Unlock()

	// Shutdown must not return an error even if platform Remove fails
	assert.NoError(t, c.ServiceShutdown(t.Context()))
}
