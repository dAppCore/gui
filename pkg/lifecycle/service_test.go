// pkg/lifecycle/service_test.go
package lifecycle

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
	mu           sync.Mutex
	handlers     map[EventType][]func()
	fileHandlers []func(string)
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		handlers: make(map[EventType][]func()),
	}
}

func (m *mockPlatform) OnApplicationEvent(eventType EventType, handler func()) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[eventType] = append(m.handlers[eventType], handler)
	idx := len(m.handlers[eventType]) - 1
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if idx < len(m.handlers[eventType]) {
			m.handlers[eventType] = append(m.handlers[eventType][:idx], m.handlers[eventType][idx+1:]...)
		}
	}
}

func (m *mockPlatform) OnOpenedWithFile(handler func(string)) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fileHandlers = append(m.fileHandlers, handler)
	idx := len(m.fileHandlers) - 1
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if idx < len(m.fileHandlers) {
			m.fileHandlers = append(m.fileHandlers[:idx], m.fileHandlers[idx+1:]...)
		}
	}
}

// simulateEvent fires all registered handlers for the given event type.
func (m *mockPlatform) simulateEvent(eventType EventType) {
	m.mu.Lock()
	handlers := make([]func(), len(m.handlers[eventType]))
	copy(handlers, m.handlers[eventType])
	m.mu.Unlock()
	for _, h := range handlers {
		h()
	}
}

// simulateFileOpen fires all registered file-open handlers.
func (m *mockPlatform) simulateFileOpen(path string) {
	m.mu.Lock()
	handlers := make([]func(string), len(m.fileHandlers))
	copy(handlers, m.fileHandlers)
	m.mu.Unlock()
	for _, h := range handlers {
		h(path)
	}
}

// handlerCount returns the number of registered handlers for event-based + file-based.
func (m *mockPlatform) handlerCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := len(m.fileHandlers)
	for _, handlers := range m.handlers {
		count += len(handlers)
	}
	return count
}

// --- Test helpers ---

func newTestLifecycleService(t *testing.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := newMockPlatform()
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "lifecycle")
	return svc, c, mock
}

// --- Tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestLifecycleService(t)
	assert.NotNil(t, svc)
}

func TestApplicationStarted_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionApplicationStarted); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventApplicationStarted)
	assert.True(t, received)
}

func TestDidBecomeActive_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionDidBecomeActive); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventDidBecomeActive)
	assert.True(t, received)
}

func TestDidResignActive_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionDidResignActive); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventDidResignActive)
	assert.True(t, received)
}

func TestWillTerminate_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionWillTerminate); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventWillTerminate)
	assert.True(t, received)
}

func TestPowerStatusChanged_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionPowerStatusChanged); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventPowerStatusChanged)
	assert.True(t, received)
}

func TestSystemSuspend_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionSystemSuspend); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventSystemSuspend)
	assert.True(t, received)
}

func TestSystemResume_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionSystemResume); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventSystemResume)
	assert.True(t, received)
}

func TestOpenedWithFile_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var receivedPath string
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionOpenedWithFile); ok {
			receivedPath = a.Path
		}
		return nil
	})

	mock.simulateFileOpen("/Users/snider/Documents/test.txt")
	assert.Equal(t, "/Users/snider/Documents/test.txt", receivedPath)
}

func TestOnShutdown_CancelsAll_Good(t *testing.T) {
	svc, _, mock := newTestLifecycleService(t)

	// Verify handlers were registered during OnStartup
	assert.Greater(t, mock.handlerCount(), 0, "handlers should be registered after OnStartup")

	// Shutdown should cancel all registrations
	err := svc.OnShutdown(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 0, mock.handlerCount(), "all handlers should be cancelled after OnShutdown")
}

func TestRegister_Bad(t *testing.T) {
	// No lifecycle service registered — actions are not received
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionApplicationStarted); ok {
			received = true
		}
		return nil
	})

	// No way to trigger events without the service
	assert.False(t, received)
}
