// pkg/keybinding/service_test.go
package keybinding

import (
	"context"
	"sync"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlatform records Add/Remove calls and allows triggering shortcuts.
type mockPlatform struct {
	mu       sync.Mutex
	handlers map[string]func()
	removed  []string
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{handlers: make(map[string]func())}
}

func (m *mockPlatform) Add(accelerator string, handler func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[accelerator] = handler
	return nil
}

func (m *mockPlatform) Remove(accelerator string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, accelerator)
	m.removed = append(m.removed, accelerator)
	return nil
}

func (m *mockPlatform) Process(accelerator string) bool {
	m.mu.Lock()
	h, ok := m.handlers[accelerator]
	m.mu.Unlock()
	if ok && h != nil {
		h()
		return true
	}
	return false
}

func (m *mockPlatform) GetAll() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, 0, len(m.handlers))
	for k := range m.handlers {
		out = append(out, k)
	}
	return out
}

// trigger simulates a shortcut keypress by calling the registered handler.
func (m *mockPlatform) trigger(accelerator string) {
	m.mu.Lock()
	h, ok := m.handlers[accelerator]
	m.mu.Unlock()
	if ok {
		h()
	}
}

func newTestKeybindingService(t *testing.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "keybinding")
	return svc, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *testing.T) {
	mp := newMockPlatform()
	svc, _ := newTestKeybindingService(t, mp)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.platform)
}

func TestTaskAdd_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.add", TaskAdd{
		Accelerator: "Ctrl+S", Description: "Save",
	})
	require.True(t, r.OK)

	// Verify binding registered on platform
	assert.Contains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskAdd_Bad_Duplicate(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Second add with same accelerator should fail
	r := taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save Again"})
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, ErrorAlreadyRegistered)
}

func TestTaskRemove_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	r := taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+S"})
	require.True(t, r.OK)

	// Verify removed from platform
	assert.NotContains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskRemove_Bad_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+X"})
	assert.False(t, r.OK)
}

func TestQueryList_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+Z", Description: "Undo"})

	r := c.QUERY(QueryList{})
	require.True(t, r.OK)
	list := r.Value.([]BindingInfo)
	assert.Len(t, list, 2)
}

func TestQueryList_Good_Empty(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := c.QUERY(QueryList{})
	require.True(t, r.OK)
	list := r.Value.([]BindingInfo)
	assert.Len(t, list, 0)
}

func TestTaskAdd_Good_TriggerBroadcast(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	// Capture broadcast actions
	var triggered ActionTriggered
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionTriggered); ok {
			mu.Lock()
			triggered = a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Simulate shortcut trigger via mock
	mp.trigger("Ctrl+S")

	mu.Lock()
	assert.Equal(t, "Ctrl+S", triggered.Accelerator)
	mu.Unlock()
}

func TestTaskAdd_Good_RebindAfterRemove(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+S"})

	// Should succeed after remove
	r := taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+S", Description: "Save v2"})
	require.True(t, r.OK)

	// Verify new description
	r2 := c.QUERY(QueryList{})
	list := r2.Value.([]BindingInfo)
	assert.Len(t, list, 1)
	assert.Equal(t, "Save v2", list[0].Description)
}

func TestQueryList_Bad_NoService(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryList{})
	assert.False(t, r.OK)
}

// --- TaskProcess tests ---

func TestTaskProcess_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+P", Description: "Print"})

	var triggered ActionTriggered
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionTriggered); ok {
			mu.Lock()
			triggered = a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "keybinding.process", TaskProcess{Accelerator: "Ctrl+P"})
	require.True(t, r.OK)

	mu.Lock()
	assert.Equal(t, "Ctrl+P", triggered.Accelerator)
	mu.Unlock()
}

func TestTaskProcess_Bad_NotRegistered(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.process", TaskProcess{Accelerator: "Ctrl+P"})
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, ErrorNotRegistered)
}

func TestTaskProcess_Ugly_RemovedBinding(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	taskRun(c, "keybinding.add", TaskAdd{Accelerator: "Ctrl+P", Description: "Print"})
	taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+P"})

	// After remove, process should fail with ErrorNotRegistered
	r := taskRun(c, "keybinding.process", TaskProcess{Accelerator: "Ctrl+P"})
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, ErrorNotRegistered)
}

// --- TaskRemove ErrorNotRegistered sentinel tests ---

func TestTaskRemove_Bad_ErrorSentinel(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	r := taskRun(c, "keybinding.remove", TaskRemove{Accelerator: "Ctrl+X"})
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, ErrorNotRegistered)
}

// --- QueryList Ugly: concurrent adds ---

func TestQueryList_Ugly_ConcurrentAdds(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	accelerators := []string{"Ctrl+1", "Ctrl+2", "Ctrl+3", "Ctrl+4", "Ctrl+5"}
	var wg sync.WaitGroup
	for _, accelerator := range accelerators {
		wg.Add(1)
		go func(acc string) {
			defer wg.Done()
			taskRun(c, "keybinding.add", TaskAdd{Accelerator: acc, Description: acc})
		}(accelerator)
	}
	wg.Wait()

	r := c.QUERY(QueryList{})
	require.True(t, r.OK)
	list := r.Value.([]BindingInfo)
	assert.Len(t, list, len(accelerators))
}
