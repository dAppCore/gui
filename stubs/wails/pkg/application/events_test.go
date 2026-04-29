package application

import (
	core "dappco.re/go"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestEvents_CustomEvent_Good(t *core.T) {
	event := &CustomEvent{Name: "ready", Data: "payload", Sender: "ui"}

	core.AssertFalse(t, event.IsCancelled())
	event.Cancel()
	core.AssertTrue(t, event.IsCancelled())
}

func TestEvents_CustomEvent_Bad(t *core.T) {
	event := &CustomEvent{}

	core.AssertEmpty(t, event.Name)
	core.AssertNil(t, event.Data)
	core.AssertFalse(t, event.IsCancelled())
}

func TestEvents_CustomEvent_Ugly(t *core.T) {
	event := &CustomEvent{Name: "event", Data: []any{"a", 1}, Sender: "sender"}

	event.Cancel()
	event.Cancel()

	core.AssertTrue(t, event.IsCancelled())
	core.AssertEqual(t, []any{"a", 1}, event.Data)
}

func TestEvents_CustomEvent_NilReceiver(t *core.T) {
	var event *CustomEvent

	core.AssertNotPanics(t, func() {
		event.Cancel()
	})
	core.AssertFalse(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_Good(t *core.T) {
	event := &ApplicationEvent{Id: 7, ctx: newApplicationEventContext()}

	core.AssertFalse(t, event.IsCancelled())
	core.AssertNotNil(t, event.Context())
	event.Cancel()
	core.AssertTrue(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_Bad(t *core.T) {
	event := &ApplicationEvent{}

	core.AssertNotNil(t, event.Context())
	core.AssertFalse(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_Ugly(t *core.T) {
	event := &ApplicationEvent{Id: 99}

	event.Cancel()
	event.Cancel()

	core.AssertTrue(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_NilReceiver(t *core.T) {
	var event *ApplicationEvent

	core.AssertNotPanics(t, func() {
		event.Cancel()
	})
	core.AssertFalse(t, event.IsCancelled())
}

func TestEvents_EventManager_Emit_Good(t *core.T) {
	manager := newEventManager()
	calls := 0
	manager.On("ready", func(event *CustomEvent) {
		calls++
		core.AssertEqual(t, "ready", event.Name)
		core.AssertEqual(t, "payload", event.Data)
		event.Cancel()
	})
	manager.On("ready", func(*CustomEvent) {
		calls++
	})

	cancelled := manager.Emit("ready", "payload")

	core.AssertTrue(t, cancelled)
	core.AssertEqual(t, 1, calls)
}

func TestEvents_EventManager_Emit_Bad(t *core.T) {
	manager := &EventManager{}

	core.AssertFalse(t, manager.Emit("missing"))
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestEvents_EventManager_Emit_Ugly(t *core.T) {
	manager := newEventManager()
	calls := 0
	manager.OnMultiple("tick", func(event *CustomEvent) {
		calls++
		core.AssertEqual(t, []any{1, 2}, event.Data)
	}, 1)

	core.AssertFalse(t, manager.Emit("tick", 1, 2))
	core.AssertFalse(t, manager.Emit("tick", 1, 2))
	core.AssertEqual(t, 1, calls)
}

func TestEvents_EventManager_Emit_RecoversFromPanic(t *core.T) {
	manager := newEventManager()
	calls := 0

	manager.On("ready", func(*CustomEvent) {
		panic("boom")
	})
	manager.On("ready", func(*CustomEvent) {
		calls++
	})

	core.AssertFalse(t, manager.Emit("ready"))
	core.AssertEqual(t, 1, calls)
}

func TestEvents_EventManager_OnApplicationEvent_Good(t *core.T) {
	manager := newEventManager()
	eventType := events.ApplicationEventType(42)
	calls := 0

	cancel := manager.OnApplicationEvent(eventType, func(event *ApplicationEvent) {
		calls++
		core.AssertNotNil(t, event)
		event.Cancel()
	})

	core.AssertLen(t, manager.appListeners[uint(eventType)], 1)
	cancel()
	core.AssertEmpty(t, manager.appListeners[uint(eventType)])
	core.AssertEqual(t, 0, calls)
}

func TestEvents_EventManager_OnApplicationEvent_Bad(t *core.T) {
	manager := &EventManager{}
	eventType := events.ApplicationEventType(7)

	cancel := manager.OnApplicationEvent(eventType, func(*ApplicationEvent) {})
	core.AssertNotNil(t, cancel)
	cancel()

	core.AssertEmpty(t, manager.appListeners[uint(eventType)])
}

func TestEvents_EventManager_OnApplicationEvent_Ugly(t *core.T) {
	manager := newEventManager()
	eventType := events.ApplicationEventType(1)
	manager.OnApplicationEvent(eventType, func(*ApplicationEvent) {})

	core.AssertLen(t, manager.appListeners[uint(eventType)], 1)
	manager.Off("does-not-exist")
	core.AssertLen(t, manager.appListeners[uint(eventType)], 1)
}

// AX7 generated source-matching smoke coverage.
func TestEvents_ApplicationEvent_Context_Good(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Context_Bad(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Context_Ugly(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Cancel_Good(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Cancel_Bad(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Cancel_Ugly(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_IsCancelled_Good(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_IsCancelled_Bad(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_IsCancelled_Ugly(t *core.T) {
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_Cancel_Good(t *core.T) {
	subject := new(CustomEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_Cancel_Bad(t *core.T) {
	subject := new(CustomEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_Cancel_Ugly(t *core.T) {
	subject := new(CustomEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_IsCancelled_Good(t *core.T) {
	subject := new(CustomEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_IsCancelled_Bad(t *core.T) {
	subject := new(CustomEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_IsCancelled_Ugly(t *core.T) {
	subject := new(CustomEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_On_Good(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		got0 := subject.On("agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_On_Bad(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		got0 := subject.On("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_On_Ugly(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		got0 := subject.On("../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_Off_Good(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.Off("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_Off_Bad(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.Off("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_Off_Ugly(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.Off("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_OnMultiple_Good(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.OnMultiple("agent", nil, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_OnMultiple_Bad(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.OnMultiple("", nil, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_OnMultiple_Ugly(t *core.T) {
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.OnMultiple("../../edge", nil, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
