// pkg/keybinding/service_test.go
package keybinding

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlatform records Add/Remove/Process calls and allows triggering shortcuts.
type mockPlatform struct {
	mu          sync.Mutex
	handlers    map[string]func()
	removed     []string
	processed   []string
	processErr  error
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

func (m *mockPlatform) Process(accelerator string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.processErr != nil {
		return m.processErr
	}
	m.processed = append(m.processed, accelerator)
	return nil
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
	c, err := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "keybinding")
	return svc, c
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

	_, handled, err := c.PERFORM(TaskAdd{
		Accelerator: "Ctrl+S", Description: "Save",
	})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify binding registered on platform
	assert.Contains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskAdd_Bad_Duplicate(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Second add with same accelerator should fail
	_, handled, err := c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save Again"})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrorAlreadyRegistered)
}

func TestTaskRemove_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	_, handled, err := c.PERFORM(TaskRemove{Accelerator: "Ctrl+S"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify removed from platform
	assert.NotContains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskRemove_Bad_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, handled, err := c.PERFORM(TaskRemove{Accelerator: "Ctrl+X"})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrorNotRegistered)
}

func TestQueryList_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+Z", Description: "Undo"})

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]BindingInfo)
	assert.Len(t, list, 2)
}

func TestQueryList_Good_Empty(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]BindingInfo)
	assert.Len(t, list, 0)
}

func TestTaskAdd_Good_TriggerBroadcast(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	// Capture broadcast actions
	var triggered ActionTriggered
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionTriggered); ok {
			mu.Lock()
			triggered = a
			mu.Unlock()
		}
		return nil
	})

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Simulate shortcut trigger via mock
	mp.trigger("Ctrl+S")

	mu.Lock()
	assert.Equal(t, "Ctrl+S", triggered.Accelerator)
	mu.Unlock()
}

func TestTaskAdd_Good_RebindAfterRemove(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	_, _, _ = c.PERFORM(TaskRemove{Accelerator: "Ctrl+S"})

	// Should succeed after remove
	_, handled, err := c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save v2"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify new description
	result, _, _ := c.QUERY(QueryList{})
	list := result.([]BindingInfo)
	assert.Len(t, list, 1)
	assert.Equal(t, "Save v2", list[0].Description)
}

func TestQueryList_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryList{})
	assert.False(t, handled)
}

func TestTaskProcess_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+P", Description: "Print"})

	_, handled, err := c.PERFORM(TaskProcess{Accelerator: "Ctrl+P"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Contains(t, mp.processed, "Ctrl+P")
}

func TestTaskProcess_Bad_NotRegistered(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, handled, err := c.PERFORM(TaskProcess{Accelerator: "Ctrl+Z"})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrorNotRegistered)
}

func TestTaskProcess_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskProcess{})
	assert.False(t, handled)
}

func TestTaskProcess_Bad_PlatformError(t *testing.T) {
	mp := newMockPlatform()
	mp.processErr = assert.AnError
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+P", Description: "Print"})
	_, handled, err := c.PERFORM(TaskProcess{Accelerator: "Ctrl+P"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestErrorNotRegistered_Ugly(t *testing.T) {
	// ErrorNotRegistered and ErrorAlreadyRegistered must be distinct sentinels.
	assert.NotEqual(t, ErrorNotRegistered, ErrorAlreadyRegistered)
	assert.NotErrorIs(t, ErrorNotRegistered, ErrorAlreadyRegistered)
}
