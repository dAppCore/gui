// pkg/events/service_test.go
package events

import (
	"context"
	"sync"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Platform ---

type mockPlatform struct {
	mu          sync.Mutex
	listeners   map[string][]*mockListener
	emitted     []CustomEvent
	resetCalled bool
	nilCancel   bool
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
	if m.nilCancel {
		return nil
	}
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
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "events")
	return svc, c, mock
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

// --- Good path tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestService(t)
	assert.NotNil(t, svc)
}

func TestTaskEmit_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: "user:login", Data: "alice"})
	require.True(t, r.OK)
	assert.Equal(t, false, r.Value) // not cancelled

	assert.Len(t, mock.emitted, 1)
	assert.Equal(t, "user:login", mock.emitted[0].Name)
	assert.Equal(t, "alice", mock.emitted[0].Data)
}

func TestTaskEmit_NoData_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: "ping"})
	require.True(t, r.OK)
	assert.Len(t, mock.emitted, 1)
	assert.Nil(t, mock.emitted[0].Data)
}

func TestTaskOn_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	var received []ActionEventFired
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if fired, ok := msg.(ActionEventFired); ok {
			received = append(received, fired)
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "events.on", TaskOn{Name: "theme:changed"})
	require.True(t, r.OK)

	mock.simulateEvent("theme:changed", "dark")

	assert.Len(t, received, 1)
	assert.Equal(t, "theme:changed", received[0].Event.Name)
	assert.Equal(t, "dark", received[0].Event.Data)
}

func TestTaskOff_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	// Register via IPC then remove
	r := taskRun(c, "events.on", TaskOn{Name: "file:saved"})
	require.True(t, r.OK)
	assert.Equal(t, 1, mock.listenerCount())

	r2 := taskRun(c, "events.off", TaskOff{Name: "file:saved"})
	require.True(t, r2.OK)
	assert.Equal(t, 0, mock.listenerCount())
}

func TestQueryListeners_Good(t *testing.T) {
	_, c, _ := newTestService(t)

	require.True(t, taskRun(c, "events.on", TaskOn{Name: "user:login"}).OK)
	require.True(t, taskRun(c, "events.on", TaskOn{Name: "user:login"}).OK)
	require.True(t, taskRun(c, "events.on", TaskOn{Name: "theme:changed"}).OK)

	r := c.QUERY(QueryListeners{})
	require.True(t, r.OK)

	infos := r.Value.([]ListenerInfo)
	counts := make(map[string]int)
	for _, info := range infos {
		counts[info.EventName] = info.Count
	}
	assert.Equal(t, 2, counts["user:login"])
	assert.Equal(t, 1, counts["theme:changed"])
}

func TestQueryListeners_Empty_Good(t *testing.T) {
	_, c, _ := newTestService(t)

	r := c.QUERY(QueryListeners{})
	require.True(t, r.OK)

	infos := r.Value.([]ListenerInfo)
	assert.Empty(t, infos)
}

func TestOnShutdown_CancelsAll_Good(t *testing.T) {
	svc, _, mock := newTestService(t)

	require.True(t, taskRun(svc.Core(), "events.on", TaskOn{Name: "a:b"}).OK)
	require.True(t, taskRun(svc.Core(), "events.on", TaskOn{Name: "c:d"}).OK)
	assert.Equal(t, 2, mock.listenerCount())

	require.True(t, svc.OnShutdown(context.Background()).OK)
	assert.Equal(t, 0, mock.listenerCount())
}

func TestOnShutdown_IgnoresNilCancels_Good(t *testing.T) {
	mock := newMockPlatform()
	mock.nilCancel = true
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "events")

	require.True(t, taskRun(c, "events.on", TaskOn{Name: "a:b"}).OK)
	assert.Equal(t, 1, mock.listenerCount())

	require.NotPanics(t, func() {
		assert.True(t, svc.OnShutdown(context.Background()).OK)
	})
}

func TestActionEventFired_BroadcastOnSimulate_Good(t *testing.T) {
	_, c, mock := newTestService(t)

	var receivedEvents []CustomEvent
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if fired, ok := msg.(ActionEventFired); ok {
			receivedEvents = append(receivedEvents, fired.Event)
		}
		return core.Result{OK: true}
	})

	require.True(t, taskRun(c, "events.on", TaskOn{Name: "data:ready"}).OK)

	mock.simulateEvent("data:ready", map[string]any{"rows": 42})

	require.Len(t, receivedEvents, 1)
	assert.Equal(t, "data:ready", receivedEvents[0].Name)
}

// --- Bad path tests ---

func TestTaskOn_EmptyName_Bad(t *testing.T) {
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.on", TaskOn{Name: ""})
	assert.False(t, r.OK)
}

func TestTaskEmit_UnknownEvent_Bad(t *testing.T) {
	// Emitting an event with no listeners is valid — returns not-cancelled.
	_, c, mock := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: "no:listeners"})
	require.True(t, r.OK)
	assert.Equal(t, false, r.Value)
	assert.Len(t, mock.emitted, 1) // still recorded as emitted
}

func TestTaskEmit_EmptyName_Bad(t *testing.T) {
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.emit", TaskEmit{Name: ""})
	assert.False(t, r.OK)
}

func TestQueryListeners_NoService_Bad(t *testing.T) {
	// No events service registered — query is not handled.
	c := core.New(core.WithServiceLock())

	r := c.QUERY(QueryListeners{})
	assert.False(t, r.OK)
}

func TestTaskEmit_NoService_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())

	r := c.Action("events.emit").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestTaskEmit_PlatformUnavailable_Bad(t *testing.T) {
	c := core.New(
		core.WithService(Register(nil)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "events")

	r := taskRun(c, "events.emit", TaskEmit{Name: "user:login"})
	require.False(t, r.OK)
	err, ok := r.Value.(error)
	require.True(t, ok)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "event platform unavailable")

	r = taskRun(c, "events.on", TaskOn{Name: "user:login"})
	require.False(t, r.OK)
	err, ok = r.Value.(error)
	require.True(t, ok)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "event platform unavailable")

	r = taskRun(c, "events.off", TaskOff{Name: "user:login"})
	require.False(t, r.OK)
	err, ok = r.Value.(error)
	require.True(t, ok)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "event platform unavailable")

	require.NotPanics(t, func() {
		assert.True(t, svc.OnShutdown(context.Background()).OK)
	})
}

// --- Ugly path tests ---

func TestTaskOff_NeverRegistered_Ugly(t *testing.T) {
	// Off on a name that was never registered is a no-op — must not panic.
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.off", TaskOff{Name: "nonexistent:event"})
	assert.True(t, r.OK)
}

func TestTaskOff_EmptyName_Bad(t *testing.T) {
	_, c, _ := newTestService(t)

	r := taskRun(c, "events.off", TaskOff{Name: ""})
	assert.False(t, r.OK)
}

func TestTaskOn_MultipleListeners_Ugly(t *testing.T) {
	// Multiple IPC listeners for the same event each receive ActionEventFired.
	_, c, mock := newTestService(t)

	var mu sync.Mutex
	var fireCount int
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionEventFired); ok {
			mu.Lock()
			fireCount++
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	taskRun(c, "events.on", TaskOn{Name: "flood"})
	taskRun(c, "events.on", TaskOn{Name: "flood"})
	taskRun(c, "events.on", TaskOn{Name: "flood"})

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
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if _, ok := msg.(ActionEventFired); ok {
			received = true
		}
		return core.Result{OK: true}
	})

	taskRun(c, "events.on", TaskOn{Name: "transient"})
	taskRun(c, "events.off", TaskOff{Name: "transient"})

	mock.simulateEvent("transient", "late-data")
	assert.False(t, received)
}

func TestQueryListeners_AfterOff_Ugly(t *testing.T) {
	// After Off, the event name must not appear in QueryListeners results.
	_, c, _ := newTestService(t)

	taskRun(c, "events.on", TaskOn{Name: "ephemeral"})
	taskRun(c, "events.off", TaskOff{Name: "ephemeral"})

	r := c.QUERY(QueryListeners{})
	infos := r.Value.([]ListenerInfo)

	for _, info := range infos {
		assert.NotEqual(t, "ephemeral", info.EventName)
	}
}
