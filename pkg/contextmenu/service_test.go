// pkg/contextmenu/service_test.go
package contextmenu

import (
	"context"
	"sync"
	"testing"

	core "dappco.re/go/core"
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

type flakyAddPlatform struct {
	*mockPlatform
	mu          sync.Mutex
	failAddOnce bool
}

func newFlakyAddPlatform() *flakyAddPlatform {
	return &flakyAddPlatform{mockPlatform: newMockPlatform()}
}

func (m *flakyAddPlatform) Add(name string, menu ContextMenuDef, onItemClick func(string, string, string)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failAddOnce {
		m.failAddOnce = false
		return ErrorMenuNotFound
	}
	return m.mockPlatform.Add(name, menu, onItemClick)
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

func newTestContextMenuService(t *testing.T, mp Platform) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "contextmenu")
	return svc, c
}

// taskRun runs a named action with a task struct and returns the result.
func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
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

	r := taskRun(c, "contextmenu.add", TaskAdd{
		Name: "file-menu",
		Menu: ContextMenuDef{
			Name: "file-menu",
			Items: []MenuItemDef{
				{Label: "Open", ActionID: "open"},
				{Label: "Delete", ActionID: "delete"},
			},
		},
	})
	require.True(t, r.OK)

	// Verify menu registered on platform
	_, ok := mp.Get("file-menu")
	assert.True(t, ok)
}

func TestTaskAdd_Good_ReplaceExisting(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Add initial menu
	_ = taskRun(c, "contextmenu.add", TaskAdd{
		Name: "ctx",
		Menu: ContextMenuDef{Name: "ctx", Items: []MenuItemDef{{Label: "A", ActionID: "a"}}},
	})

	// Replace with new menu
	r := taskRun(c, "contextmenu.add", TaskAdd{
		Name: "ctx",
		Menu: ContextMenuDef{Name: "ctx", Items: []MenuItemDef{{Label: "B", ActionID: "b"}}},
	})
	require.True(t, r.OK)

	// Verify registry has new menu
	qr := c.QUERY(QueryGet{Name: "ctx"})
	require.True(t, qr.OK)
	def := qr.Value.(*ContextMenuDef)
	require.Len(t, def.Items, 1)
	assert.Equal(t, "B", def.Items[0].Label)
}

func TestTaskRemove_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Add then remove
	_ = taskRun(c, "contextmenu.add", TaskAdd{
		Name: "test",
		Menu: ContextMenuDef{Name: "test"},
	})
	r := taskRun(c, "contextmenu.remove", TaskRemove{Name: "test"})
	require.True(t, r.OK)

	// Verify removed from registry
	qr := c.QUERY(QueryGet{Name: "test"})
	assert.Nil(t, qr.Value)
}

func TestTaskRemove_Bad_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	r := taskRun(c, "contextmenu.remove", TaskRemove{Name: "nonexistent"})
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, ErrorMenuNotFound)
}

func TestQueryGet_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{
		Name: "my-menu",
		Menu: ContextMenuDef{
			Name:  "my-menu",
			Items: []MenuItemDef{{Label: "Edit", ActionID: "edit"}},
		},
	})

	r := c.QUERY(QueryGet{Name: "my-menu"})
	require.True(t, r.OK)
	def := r.Value.(*ContextMenuDef)
	assert.Equal(t, "my-menu", def.Name)
	assert.Len(t, def.Items, 1)
}

func TestQueryGet_Good_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	r := c.QUERY(QueryGet{Name: "missing"})
	require.True(t, r.OK)
	assert.Nil(t, r.Value)
}

func TestQueryList_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "a", Menu: ContextMenuDef{Name: "a"}})
	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "b", Menu: ContextMenuDef{Name: "b"}})

	r := c.QUERY(QueryList{})
	require.True(t, r.OK)
	list := r.Value.(map[string]ContextMenuDef)
	assert.Len(t, list, 2)
}

func TestQueryList_Good_Empty(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	r := c.QUERY(QueryList{})
	require.True(t, r.OK)
	list := r.Value.(map[string]ContextMenuDef)
	assert.Len(t, list, 0)
}

func TestTaskAdd_Good_ClickBroadcast(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Capture broadcast actions
	var clicked ActionItemClicked
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionItemClicked); ok {
			mu.Lock()
			clicked = a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	_ = taskRun(c, "contextmenu.add", TaskAdd{
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

func TestTaskAdd_Ugly_PlatformAddFailureRollsBackExistingMenu(t *testing.T) {
	mp := newFlakyAddPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{
		Name: "file-menu",
		Menu: ContextMenuDef{Name: "file-menu", Items: []MenuItemDef{{Label: "Open", ActionID: "open"}}},
	})

	mp.mu.Lock()
	mp.failAddOnce = true
	mp.mu.Unlock()

	r := taskRun(c, "contextmenu.add", TaskAdd{
		Name: "file-menu",
		Menu: ContextMenuDef{Name: "file-menu", Items: []MenuItemDef{{Label: "Delete", ActionID: "delete"}}},
	})
	assert.False(t, r.OK)

	qr := c.QUERY(QueryGet{Name: "file-menu"})
	require.True(t, qr.OK)
	def := qr.Value.(*ContextMenuDef)
	require.Len(t, def.Items, 1)
	assert.Equal(t, "Open", def.Items[0].Label)
	platformMenu, ok := mp.Get("file-menu")
	require.True(t, ok)
	require.Len(t, platformMenu.Items, 1)
	assert.Equal(t, "Open", platformMenu.Items[0].Label)
}

func TestTaskAdd_Good_SubmenuItems(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	r := taskRun(c, "contextmenu.add", TaskAdd{
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
	require.True(t, r.OK)

	qr := c.QUERY(QueryGet{Name: "nested"})
	def := qr.Value.(*ContextMenuDef)
	assert.Len(t, def.Items, 3)
	assert.Len(t, def.Items[0].Items, 2) // submenu children
}

func TestQueryList_Bad_NoService(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryList{})
	assert.False(t, r.OK)
}

// --- TaskUpdate ---

func TestTaskUpdate_Good(t *testing.T) {
	// Update replaces items on an existing menu
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{
		Name: "edit-menu",
		Menu: ContextMenuDef{Name: "edit-menu", Items: []MenuItemDef{{Label: "Cut", ActionID: "cut"}}},
	})

	r := taskRun(c, "contextmenu.update", TaskUpdate{
		Name: "edit-menu",
		Menu: ContextMenuDef{Name: "edit-menu", Items: []MenuItemDef{
			{Label: "Cut", ActionID: "cut"},
			{Label: "Copy", ActionID: "copy"},
		}},
	})
	require.True(t, r.OK)

	qr := c.QUERY(QueryGet{Name: "edit-menu"})
	def := qr.Value.(*ContextMenuDef)
	assert.Len(t, def.Items, 2)
	assert.Equal(t, "Copy", def.Items[1].Label)
}

func TestTaskUpdate_Bad_NotFound(t *testing.T) {
	// Update on a non-existent menu returns ErrorMenuNotFound
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	r := taskRun(c, "contextmenu.update", TaskUpdate{
		Name: "ghost",
		Menu: ContextMenuDef{Name: "ghost"},
	})
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, ErrorMenuNotFound)
}

func TestTaskUpdate_Ugly_PlatformRemoveError(t *testing.T) {
	// Platform Remove fails mid-update — error is propagated
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{
		Name: "tricky",
		Menu: ContextMenuDef{Name: "tricky"},
	})

	mp.mu.Lock()
	mp.removeErr = ErrorMenuNotFound // reuse sentinel as a platform-level error
	mp.mu.Unlock()

	r := taskRun(c, "contextmenu.update", TaskUpdate{
		Name: "tricky",
		Menu: ContextMenuDef{Name: "tricky", Items: []MenuItemDef{{Label: "X", ActionID: "x"}}},
	})
	assert.False(t, r.OK)
}

func TestTaskUpdate_Ugly_PlatformAddFailureRollsBackExistingMenu(t *testing.T) {
	mp := newFlakyAddPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{
		Name: "edit-menu",
		Menu: ContextMenuDef{Name: "edit-menu", Items: []MenuItemDef{{Label: "Cut", ActionID: "cut"}}},
	})

	mp.mu.Lock()
	mp.failAddOnce = true
	mp.mu.Unlock()

	r := taskRun(c, "contextmenu.update", TaskUpdate{
		Name: "edit-menu",
		Menu: ContextMenuDef{Name: "edit-menu", Items: []MenuItemDef{{Label: "Copy", ActionID: "copy"}}},
	})
	assert.False(t, r.OK)

	qr := c.QUERY(QueryGet{Name: "edit-menu"})
	require.True(t, qr.OK)
	def := qr.Value.(*ContextMenuDef)
	require.Len(t, def.Items, 1)
	assert.Equal(t, "Cut", def.Items[0].Label)
	platformMenu, ok := mp.Get("edit-menu")
	require.True(t, ok)
	require.Len(t, platformMenu.Items, 1)
	assert.Equal(t, "Cut", platformMenu.Items[0].Label)
}

// --- TaskDestroy ---

func TestTaskDestroy_Good(t *testing.T) {
	// Destroy removes the menu and releases platform resources
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "doomed", Menu: ContextMenuDef{Name: "doomed"}})

	r := taskRun(c, "contextmenu.destroy", TaskDestroy{Name: "doomed"})
	require.True(t, r.OK)

	qr := c.QUERY(QueryGet{Name: "doomed"})
	assert.Nil(t, qr.Value)

	_, ok := mp.Get("doomed")
	assert.False(t, ok)
}

func TestTaskDestroy_Bad_NotFound(t *testing.T) {
	// Destroy on a non-existent menu returns ErrorMenuNotFound
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	r := taskRun(c, "contextmenu.destroy", TaskDestroy{Name: "nonexistent"})
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, ErrorMenuNotFound)
}

func TestTaskDestroy_Ugly_PlatformError(t *testing.T) {
	// Platform Remove fails — error is propagated but service remains consistent
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "frail", Menu: ContextMenuDef{Name: "frail"}})

	mp.mu.Lock()
	mp.removeErr = ErrorMenuNotFound
	mp.mu.Unlock()

	r := taskRun(c, "contextmenu.destroy", TaskDestroy{Name: "frail"})
	assert.False(t, r.OK)
}

// --- QueryGetAll ---

func TestQueryGetAll_Good(t *testing.T) {
	// QueryGetAll returns all registered menus (equivalent to QueryList)
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "x", Menu: ContextMenuDef{Name: "x"}})
	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "y", Menu: ContextMenuDef{Name: "y"}})

	r := c.QUERY(QueryGetAll{})
	require.True(t, r.OK)
	all := r.Value.(map[string]ContextMenuDef)
	assert.Len(t, all, 2)
	assert.Contains(t, all, "x")
	assert.Contains(t, all, "y")
}

func TestQueryGetAll_Bad_Empty(t *testing.T) {
	// QueryGetAll on an empty registry returns an empty map
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	r := c.QUERY(QueryGetAll{})
	require.True(t, r.OK)
	all := r.Value.(map[string]ContextMenuDef)
	assert.Len(t, all, 0)
}

func TestQueryGetAll_Ugly_NoService(t *testing.T) {
	// No contextmenu service — query is unhandled
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryGetAll{})
	assert.False(t, r.OK)
}

// --- OnShutdown ---

func TestOnShutdown_Good_CleansUpMenus(t *testing.T) {
	// OnShutdown removes all registered menus from the platform
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "alpha", Menu: ContextMenuDef{Name: "alpha"}})
	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "beta", Menu: ContextMenuDef{Name: "beta"}})

	require.True(t, c.ServiceShutdown(t.Context()).OK)

	assert.Len(t, mp.menus, 0)
}

func TestOnShutdown_Bad_NothingRegistered(t *testing.T) {
	// OnShutdown with no menus — no-op, no error
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	assert.True(t, c.ServiceShutdown(t.Context()).OK)
}

func TestOnShutdown_Ugly_PlatformRemoveErrors(t *testing.T) {
	// Platform Remove errors during shutdown are silently swallowed
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_ = taskRun(c, "contextmenu.add", TaskAdd{Name: "stubborn", Menu: ContextMenuDef{Name: "stubborn"}})

	mp.mu.Lock()
	mp.removeErr = ErrorMenuNotFound
	mp.mu.Unlock()

	// Shutdown must not return an error even if platform Remove fails
	assert.True(t, c.ServiceShutdown(t.Context()).OK)
}
