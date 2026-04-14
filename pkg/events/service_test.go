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
	emitted   []mockEmit
}

type mockListener struct {
	callback func(*CustomEvent)
}

type mockEmit struct {
	name string
	data any
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		listeners: make(map[string][]*mockListener),
	}
}

func (m *mockPlatform) Emit(name string, data ...any) bool {
	m.mu.Lock()
	var payload any
	if len(data) == 1 {
		payload = data[0]
	} else if len(data) > 1 {
		payload = data
	}
	m.emitted = append(m.emitted, mockEmit{name: name, data: payload})
	listeners := append([]*mockListener(nil), m.listeners[name]...)
	m.mu.Unlock()

	event := &CustomEvent{Name: name, Data: payload}
	for _, l := range listeners {
		l.callback(event)
	}
	return false
}

func (m *mockPlatform) On(name string, callback func(*CustomEvent)) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	listener := &mockListener{callback: callback}
	m.listeners[name] = append(m.listeners[name], listener)
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		slice := m.listeners[name]
		for i, l := range slice {
			if l == listener {
				m.listeners[name] = append(slice[:i], slice[i+1:]...)
				return
			}
		}
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
	listener := &mockListener{callback: callback}
	m.listeners[name] = append(m.listeners[name], listener)
}

func (m *mockPlatform) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = make(map[string][]*mockListener)
}

func (m *mockPlatform) listenerCount(name string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.listeners[name])
}

func (m *mockPlatform) emitCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.emitted)
}

// --- Test helpers ---

func newTestEventsService(t *testing.T) (*Service, *core.Core, *mockPlatform) {
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

// --- Good tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestEventsService(t)
	assert.NotNil(t, svc)
}

func TestTaskEmit_Good(t *testing.T) {
	_, c, mock := newTestEventsService(t)

	var fired ActionEventFired
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionEventFired); ok {
			fired = a
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskEmit{Name: "ping", Data: "hello"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 1, mock.emitCount())
	assert.Equal(t, "ping", fired.Name)
	assert.Equal(t, "hello", fired.Data)
}

func TestTaskOn_Good(t *testing.T) {
	_, c, mock := newTestEventsService(t)

	_, handled, err := c.PERFORM(TaskOn{Name: "build:done"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 1, mock.listenerCount("build:done"))
}

func TestTaskOff_Good(t *testing.T) {
	_, c, mock := newTestEventsService(t)

	_, _, _ = c.PERFORM(TaskOn{Name: "build:done"})
	assert.Equal(t, 1, mock.listenerCount("build:done"))

	_, handled, err := c.PERFORM(TaskOff{Name: "build:done"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 0, mock.listenerCount("build:done"))
}

func TestQueryListeners_Good(t *testing.T) {
	_, c, _ := newTestEventsService(t)

	_, _, _ = c.PERFORM(TaskOn{Name: "my:event"})
	_, _, _ = c.PERFORM(TaskOn{Name: "my:event"})

	result, handled, err := c.QUERY(QueryListeners{Name: "my:event"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 2, result.(int))
}

func TestOnShutdown_ResetsAll_Good(t *testing.T) {
	svc, _, mock := newTestEventsService(t)

	_, _, _ = svc.Core().PERFORM(TaskOn{Name: "a"})
	_, _, _ = svc.Core().PERFORM(TaskOn{Name: "b"})

	err := svc.OnShutdown(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 0, mock.listenerCount("a"))
	assert.Equal(t, 0, mock.listenerCount("b"))

	result, _, _ := svc.Core().QUERY(QueryListeners{Name: "a"})
	assert.Equal(t, 0, result.(int))
}

// --- Bad tests ---

func TestTaskEmit_Bad_EmptyName(t *testing.T) {
	_, c, _ := newTestEventsService(t)
	_, handled, err := c.PERFORM(TaskEmit{Name: ""})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskOn_Bad_EmptyName(t *testing.T) {
	_, c, _ := newTestEventsService(t)
	_, handled, err := c.PERFORM(TaskOn{Name: ""})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskOff_Bad_EmptyName(t *testing.T) {
	_, c, _ := newTestEventsService(t)
	_, handled, err := c.PERFORM(TaskOff{Name: ""})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestQueryListeners_Bad_NoService(t *testing.T) {
	// No events service registered — query is not handled
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.QUERY(QueryListeners{Name: "anything"})
	assert.False(t, handled)
}

// --- Ugly tests ---

func TestTaskEmit_Ugly_NilData(t *testing.T) {
	_, c, mock := newTestEventsService(t)
	_, handled, err := c.PERFORM(TaskEmit{Name: "null:event", Data: nil})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 1, mock.emitCount())
}

func TestTaskOn_Ugly_MultipleListenersSameEvent(t *testing.T) {
	_, c, mock := newTestEventsService(t)

	for i := 0; i < 5; i++ {
		_, _, _ = c.PERFORM(TaskOn{Name: "flood"})
	}

	result, _, _ := c.QUERY(QueryListeners{Name: "flood"})
	assert.Equal(t, 5, result.(int))
	assert.Equal(t, 5, mock.listenerCount("flood"))
}

func TestTaskOff_Ugly_NonExistentEvent(t *testing.T) {
	// Off on unknown event name should not error
	_, c, _ := newTestEventsService(t)
	_, handled, err := c.PERFORM(TaskOff{Name: "never-registered"})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskEmit_Ugly_ActionBroadcastOnEachEmit(t *testing.T) {
	_, c, _ := newTestEventsService(t)

	var count int
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionEventFired); ok {
			count++
		}
		return nil
	})

	for i := 0; i < 3; i++ {
		_, _, _ = c.PERFORM(TaskEmit{Name: "tick"})
	}

	assert.Equal(t, 3, count)
}
