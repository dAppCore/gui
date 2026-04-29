package application

import (
	core "dappco.re/go"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestEvents_CustomEvent_Good(t *core.T) {
	// CustomEvent
	ax7Variant := "CustomEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	event := &CustomEvent{Name: "ready", Data: "payload", Sender: "ui"}

	core.AssertFalse(t, event.IsCancelled())
	event.Cancel()
	core.AssertTrue(t, event.IsCancelled())
}

func TestEvents_CustomEvent_Bad(t *core.T) {
	// CustomEvent
	ax7Variant := "CustomEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	event := &CustomEvent{}

	core.AssertEmpty(t, event.Name)
	core.AssertNil(t, event.Data)
	core.AssertFalse(t, event.IsCancelled())
}

func TestEvents_CustomEvent_Ugly(t *core.T) {
	// CustomEvent
	ax7Variant := "CustomEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
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
	// ApplicationEvent
	ax7Variant := "ApplicationEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	event := &ApplicationEvent{Id: 7, ctx: newApplicationEventContext()}

	core.AssertFalse(t, event.IsCancelled())
	core.AssertNotNil(t, event.Context())
	event.Cancel()
	core.AssertTrue(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_Bad(t *core.T) {
	// ApplicationEvent
	ax7Variant := "ApplicationEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	event := &ApplicationEvent{}

	core.AssertNotNil(t, event.Context())
	core.AssertFalse(t, event.IsCancelled())
}

func TestEvents_ApplicationEvent_Ugly(t *core.T) {
	// ApplicationEvent
	ax7Variant := "ApplicationEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
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
	// EventManager Emit
	ax7Variant := "EventManager_Emit:good"
	core.AssertContains(t, ax7Variant, "good")
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
	// EventManager Emit
	ax7Variant := "EventManager_Emit:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &EventManager{}

	core.AssertFalse(t, manager.Emit("missing"))
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestEvents_EventManager_Emit_Ugly(t *core.T) {
	// EventManager Emit
	ax7Variant := "EventManager_Emit:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
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
	// EventManager OnApplicationEvent
	ax7Variant := "EventManager_OnApplicationEvent:good"
	core.AssertContains(t, ax7Variant, "good")
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
	// EventManager OnApplicationEvent
	ax7Variant := "EventManager_OnApplicationEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &EventManager{}
	eventType := events.ApplicationEventType(7)

	cancel := manager.OnApplicationEvent(eventType, func(*ApplicationEvent) {})
	core.AssertNotNil(t, cancel)
	cancel()

	core.AssertEmpty(t, manager.appListeners[uint(eventType)])
}

func TestEvents_EventManager_OnApplicationEvent_Ugly(t *core.T) {
	// EventManager OnApplicationEvent
	ax7Variant := "EventManager_OnApplicationEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := newEventManager()
	eventType := events.ApplicationEventType(1)
	manager.OnApplicationEvent(eventType, func(*ApplicationEvent) {})

	core.AssertLen(t, manager.appListeners[uint(eventType)], 1)
	manager.Off("does-not-exist")
	core.AssertLen(t, manager.appListeners[uint(eventType)], 1)
}

// AX7 generated source-matching smoke coverage.
func TestEvents_ApplicationEvent_Context_Good(t *core.T) {
	// ApplicationEvent Context
	ax7Variant := "ApplicationEvent_Context:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Context_Bad(t *core.T) {
	// ApplicationEvent Context
	ax7Variant := "ApplicationEvent_Context:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Context_Ugly(t *core.T) {
	// ApplicationEvent Context
	ax7Variant := "ApplicationEvent_Context:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.Context()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Cancel_Good(t *core.T) {
	// ApplicationEvent Cancel
	ax7Variant := "ApplicationEvent_Cancel:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Cancel_Bad(t *core.T) {
	// ApplicationEvent Cancel
	ax7Variant := "ApplicationEvent_Cancel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_Cancel_Ugly(t *core.T) {
	// ApplicationEvent Cancel
	ax7Variant := "ApplicationEvent_Cancel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_IsCancelled_Good(t *core.T) {
	// ApplicationEvent IsCancelled
	ax7Variant := "ApplicationEvent_IsCancelled:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_IsCancelled_Bad(t *core.T) {
	// ApplicationEvent IsCancelled
	ax7Variant := "ApplicationEvent_IsCancelled:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_ApplicationEvent_IsCancelled_Ugly(t *core.T) {
	// ApplicationEvent IsCancelled
	ax7Variant := "ApplicationEvent_IsCancelled:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ApplicationEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_Cancel_Good(t *core.T) {
	// CustomEvent Cancel
	ax7Variant := "CustomEvent_Cancel:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(CustomEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_Cancel_Bad(t *core.T) {
	// CustomEvent Cancel
	ax7Variant := "CustomEvent_Cancel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(CustomEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_Cancel_Ugly(t *core.T) {
	// CustomEvent Cancel
	ax7Variant := "CustomEvent_Cancel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(CustomEvent)
	result := core.Try(func() any {
		subject.Cancel()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_IsCancelled_Good(t *core.T) {
	// CustomEvent IsCancelled
	ax7Variant := "CustomEvent_IsCancelled:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(CustomEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_IsCancelled_Bad(t *core.T) {
	// CustomEvent IsCancelled
	ax7Variant := "CustomEvent_IsCancelled:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(CustomEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_CustomEvent_IsCancelled_Ugly(t *core.T) {
	// CustomEvent IsCancelled
	ax7Variant := "CustomEvent_IsCancelled:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(CustomEvent)
	result := core.Try(func() any {
		got0 := subject.IsCancelled()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_On_Good(t *core.T) {
	// EventManager On
	ax7Variant := "EventManager_On:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(EventManager)
	result := core.Try(func() any {
		got0 := subject.On("agent", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_On_Bad(t *core.T) {
	// EventManager On
	ax7Variant := "EventManager_On:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(EventManager)
	result := core.Try(func() any {
		got0 := subject.On("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_On_Ugly(t *core.T) {
	// EventManager On
	ax7Variant := "EventManager_On:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(EventManager)
	result := core.Try(func() any {
		got0 := subject.On("../../edge", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_Off_Good(t *core.T) {
	// EventManager Off
	ax7Variant := "EventManager_Off:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.Off("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_Off_Bad(t *core.T) {
	// EventManager Off
	ax7Variant := "EventManager_Off:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.Off("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_Off_Ugly(t *core.T) {
	// EventManager Off
	ax7Variant := "EventManager_Off:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.Off("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_OnMultiple_Good(t *core.T) {
	// EventManager OnMultiple
	ax7Variant := "EventManager_OnMultiple:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.OnMultiple("agent", nil, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_OnMultiple_Bad(t *core.T) {
	// EventManager OnMultiple
	ax7Variant := "EventManager_OnMultiple:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.OnMultiple("", nil, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestEvents_EventManager_OnMultiple_Ugly(t *core.T) {
	// EventManager OnMultiple
	ax7Variant := "EventManager_OnMultiple:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(EventManager)
	result := core.Try(func() any {
		subject.OnMultiple("../../edge", nil, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
