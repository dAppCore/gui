package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestEvents_CustomEvent_Good(t *testing.T) {
	event := &CustomEvent{Name: "ready", Data: "payload", Sender: "ui"}

	assert.False(t, event.IsCancelled())
	event.Cancel()
	assert.True(t, event.IsCancelled())
}

func TestEvents_CustomEvent_Bad(t *testing.T) {
	event := &CustomEvent{}

	assert.Empty(t, event.Name)
	assert.Nil(t, event.Data)
	assert.False(t, event.IsCancelled())
}

func TestEvents_CustomEvent_Ugly(t *testing.T) {
	event := &CustomEvent{Name: "event", Data: []any{"a", 1}, Sender: "sender"}

	event.Cancel()
	event.Cancel()

	assert.True(t, event.IsCancelled())
	assert.Equal(t, []any{"a", 1}, event.Data)
}

func TestEvents_ApplicationEvent_Good(t *testing.T) {
	event := &ApplicationEvent{Id: 7, ctx: newApplicationEventContext()}

	assert.False(t, event.IsCancelled())
	require.NotNil(t, event.Context())
	event.Cancel()
	assert.True(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_Bad(t *testing.T) {
	event := &ApplicationEvent{}

	assert.Nil(t, event.Context())
	assert.False(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_Ugly(t *testing.T) {
	event := &ApplicationEvent{Id: 99}

	event.Cancel()
	event.Cancel()

	assert.True(t, event.IsCancelled())
}

func TestEvents_EventManager_Emit_Good(t *testing.T) {
	manager := newEventManager()
	calls := 0
	manager.On("ready", func(event *CustomEvent) {
		calls++
		assert.Equal(t, "ready", event.Name)
		assert.Equal(t, "payload", event.Data)
		event.Cancel()
	})
	manager.On("ready", func(*CustomEvent) {
		calls++
	})

	cancelled := manager.Emit("ready", "payload")

	assert.True(t, cancelled)
	assert.Equal(t, 1, calls)
}

func TestEvents_EventManager_Emit_Bad(t *testing.T) {
	manager := &EventManager{}

	assert.False(t, manager.Emit("missing"))
}

func TestEvents_EventManager_Emit_Ugly(t *testing.T) {
	manager := newEventManager()
	calls := 0
	manager.OnMultiple("tick", func(event *CustomEvent) {
		calls++
		assert.Equal(t, []any{1, 2}, event.Data)
	}, 1)

	assert.False(t, manager.Emit("tick", 1, 2))
	assert.False(t, manager.Emit("tick", 1, 2))
	assert.Equal(t, 1, calls)
}

func TestEvents_EventManager_OnApplicationEvent_Good(t *testing.T) {
	manager := newEventManager()
	eventType := events.ApplicationEventType(42)
	calls := 0

	cancel := manager.OnApplicationEvent(eventType, func(event *ApplicationEvent) {
		calls++
		require.NotNil(t, event)
		event.Cancel()
	})

	require.Len(t, manager.appListeners[uint(eventType)], 1)
	cancel()
	assert.Empty(t, manager.appListeners[uint(eventType)])
	assert.Equal(t, 0, calls)
}

func TestEvents_EventManager_OnApplicationEvent_Bad(t *testing.T) {
	manager := &EventManager{}
	eventType := events.ApplicationEventType(7)

	cancel := manager.OnApplicationEvent(eventType, func(*ApplicationEvent) {})
	require.NotNil(t, cancel)
	cancel()

	assert.Empty(t, manager.appListeners[uint(eventType)])
}

func TestEvents_EventManager_OnApplicationEvent_Ugly(t *testing.T) {
	manager := newEventManager()
	eventType := events.ApplicationEventType(1)
	manager.OnApplicationEvent(eventType, func(*ApplicationEvent) {})

	assert.Len(t, manager.appListeners[uint(eventType)], 1)
	manager.Off("does-not-exist")
	assert.Len(t, manager.appListeners[uint(eventType)], 1)
}
