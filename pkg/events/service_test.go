// pkg/events/service_test.go
package events

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Platform ---

type mockPlatform struct {
	mu        sync.Mutex
	listeners map[string][]*mockListener
	emitted   []CustomEvent
	resetCalled bool
}

type mockListener struct {
	callback func(*CustomEvent)
	counter  int // -1 = persistent
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		listeners: make(map[string][]*mockListener),
	}
}

func (m *mockPlatform) Emit(name string, data ...any) bool {
	event := &CustomEvent{Name: name}
	if len(data) == 1 {
		event.Data = data[0]
	} else if len(data) > 1 {
		event.Data = data
	}

	m.mu.Lock()
	m.emitted = append(m.emitted, *event)
	active := make([]*mockListener, len(m.listeners[name]))
	copy(active, m.listeners[name])
	m.mu.Unlock()

	for _, listener := range active {
		listener.callback(event)
	}
	return false
}

func (m *mockPlatform) On(name string, callback func(*CustomEvent)) func() {
	listener := &mockListener{callback: callback, counter: -1}
	m.mu.Lock()
	m.listeners[name] = append(m.listeners[name], listener)
	m.mu.Unlock()
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		updated := m.listeners[name][:0]
		for _, existing := range m.listeners[name] {
			if existing != listener {
				updated = append(updated, existing)
			}
		}
		m.listeners[name] = updated
	}
}

func (m *mockPlatform) Off(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.listeners, name)
}

func (m *mockPlatform) OnMultiple(name string, callback func(*CustomEvent), counter int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners[name] = append(m.listeners[name], &mockListener{callback: callback, counter: counter})
}

func (m *mockPlatform) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = make(map[string][]*mockListener)
	m.resetCalled = true
}

// simulateEvent fires all registered listeners for the given event name with optional data.
func (m *mockPlatform) simulateEvent(name string, data any) {
	event := &CustomEvent{Name: name, Data: data}
	m.mu.Lock()
	active := make([]*mockListener, len(m.listeners[name]))
	copy(active, m.listeners[name])
	m.mu.Unlock()
	for _, listener := range active {
		listener.callback(event)
	}
}

// listenerCount returns the total number of registered listeners across all event names.
func (m *mockPlatform) listenerCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	total := 0
	for _, listeners := range m.listeners {
		total += len(listeners)
	}
	return total
}

// --- Test helpers ---

func newTestService(t *testing.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := newMockPlatform()
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "events")
	return svc, c, mock
}

// --- Good path tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestService(t)
	assert.NotNil(t, svc)
}

func TestTaskEmit_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	result, handled, err := c.PERFORM(TaskEmit{Name: "user:login", Data: "alice"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, false, result) // not cancelled

	assert.Len(t, mock.emitted, 1)
	assert.Equal(t, "user:login", mock.emitted[0].Name)
	assert.Equal(t, "alice", mock.emitted[0].Data)
}

func TestTaskEmit_NoData_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	_, handled, err := c.PERFORM(TaskEmit{Name: "ping"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Len(t, mock.emitted, 1)
	assert.Nil(t, mock.emitted[0].Data)
}

func TestTaskOn_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	var received []ActionEventFired
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if fired, ok := msg.(ActionEventFired); ok {
			received = append(received, fired)
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskOn{Name: "theme:changed"})
	require.NoError(t, err)
	assert.True(t, handled)

	mock.simulateEvent("theme:changed", "dark")

	assert.Len(t, received, 1)
	assert.Equal(t, "theme:changed", received[0].Event.Name)
	assert.Equal(t, "dark", received[0].Event.Data)
}

func TestTaskOff_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	// Register via IPC then remove
	_, _, err := c.PERFORM(TaskOn{Name: "file:saved"})
	require.NoError(t, err)
	assert.Equal(t, 1, mock.listenerCount())

	_, handled, err := c.PERFORM(TaskOff{Name: "file:saved"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 0, mock.listenerCount())
}

func TestQueryListeners_Good(t *testing.T) {
	_, c, _ := newTestService(t)

	_, _, err := c.PERFORM(TaskOn{Name: "user:login"})
	require.NoError(t, err)
	_, _, err = c.PERFORM(TaskOn{Name: "user:login"})
	require.NoError(t, err)
	_, _, err = c.PERFORM(TaskOn{Name: "theme:changed"})
	require.NoError(t, err)

	result, handled, err := c.QUERY(QueryListeners{})
	require.NoError(t, err)
	assert.True(t, handled)

	infos := result.([]ListenerInfo)
	counts := make(map[string]int)
	for _, info := range infos {
		counts[info.EventName] = info.Count
	}
	assert.Equal(t, 2, counts["user:login"])
	assert.Equal(t, 1, counts["theme:changed"])
}

func TestQueryListeners_Empty_Good(t *testing.T) {
	_, c, _ := newTestService(t)

	result, handled, err := c.QUERY(QueryListeners{})
	require.NoError(t, err)
	assert.True(t, handled)

	infos := result.([]ListenerInfo)
	assert.Empty(t, infos)
}

func TestOnShutdown_CancelsAll_Good(t *testing.T) {
	svc, _, mock := newTestService(t)

	_, _, err := svc.Core().PERFORM(TaskOn{Name: "a:b"})
	require.NoError(t, err)
	_, _, err = svc.Core().PERFORM(TaskOn{Name: "c:d"})
	require.NoError(t, err)
	assert.Equal(t, 2, mock.listenerCount())

	err = svc.OnShutdown(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, mock.listenerCount())
}

func TestActionEventFired_BroadcastOnSimulate_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	var receivedEvents []CustomEvent
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if fired, ok := msg.(ActionEventFired); ok {
			receivedEvents = append(receivedEvents, fired.Event)
		}
		return nil
	})

	_, _, err := c.PERFORM(TaskOn{Name: "data:ready"})
	require.NoError(t, err)

	mock.simulateEvent("data:ready", map[string]any{"rows": 42})

	require.Len(t, receivedEvents, 1)
	assert.Equal(t, "data:ready", receivedEvents[0].Name)
}

// --- Bad path tests ---

func TestTaskOn_EmptyName_Bad(t *testing.T) {
	_, c, _ := newTestService(t)

	_, handled, err := c.PERFORM(TaskOn{Name: ""})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskEmit_UnknownEvent_Bad(t *testing.T) {
	// Emitting an event with no listeners is valid — returns not-cancelled.
	_, c, mock := newTestService(t)

	result, handled, err := c.PERFORM(TaskEmit{Name: "no:listeners"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, false, result)
	assert.Len(t, mock.emitted, 1) // still recorded as emitted
}

func TestQueryListeners_NoService_Bad(t *testing.T) {
	// No events service registered — query is not handled.
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)

	_, handled, _ := c.QUERY(QueryListeners{})
	assert.False(t, handled)
}

func TestTaskEmit_NoService_Bad(t *testing.T) {
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)

	_, handled, _ := c.PERFORM(TaskEmit{Name: "x"})
	assert.False(t, handled)
}

// --- Ugly path tests ---

func TestTaskOff_NeverRegistered_Ugly(t *testing.T) {
	// Off on a name that was never registered is a no-op — must not panic.
	_, c, _ := newTestService(t)

	_, handled, err := c.PERFORM(TaskOff{Name: "nonexistent:event"})
	assert.True(t, handled)
	assert.NoError(t, err)
}

func TestTaskOn_MultipleListeners_Ugly(t *testing.T) {
	// Multiple IPC listeners for the same event each receive ActionEventFired.
	_, c, mock := newTestService(t)

	var mu sync.Mutex
	var fireCount int
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionEventFired); ok {
			mu.Lock()
			fireCount++
			mu.Unlock()
		}
		return nil
	})

	_, _, _ = c.PERFORM(TaskOn{Name: "flood"})
	_, _, _ = c.PERFORM(TaskOn{Name: "flood"})
	_, _, _ = c.PERFORM(TaskOn{Name: "flood"})

	mock.simulateEvent("flood", nil)

	mu.Lock()
	count := fireCount
	mu.Unlock()
	assert.Equal(t, 3, count)
}

func TestTaskOff_ThenEmit_Ugly(t *testing.T) {
	// After Off, simulating the event must not trigger any IPC actions.
	_, c, mock := newTestService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionEventFired); ok {
			received = true
		}
		return nil
	})

	_, _, _ = c.PERFORM(TaskOn{Name: "transient"})
	_, _, _ = c.PERFORM(TaskOff{Name: "transient"})

	mock.simulateEvent("transient", "late-data")
	assert.False(t, received)
}

func TestQueryListeners_AfterOff_Ugly(t *testing.T) {
	// After Off, the event name must not appear in QueryListeners results.
	_, c, _ := newTestService(t)

	_, _, _ = c.PERFORM(TaskOn{Name: "ephemeral"})
	_, _, _ = c.PERFORM(TaskOff{Name: "ephemeral"})

	result, _, _ := c.QUERY(QueryListeners{})
	infos := result.([]ListenerInfo)

	for _, info := range infos {
		assert.NotEqual(t, "ephemeral", info.EventName)
	}
}
