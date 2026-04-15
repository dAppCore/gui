package application

import "testing"

func TestNewAndGetExposeSingletonState(t *testing.T) {
	globalApplication = nil
	app := New(Options{Name: "Parity"})
	if Get() != app {
		t.Fatalf("Get() did not return singleton app")
	}
	if app.Config().Name != "Parity" {
		t.Fatalf("unexpected app name: %q", app.Config().Name)
	}
	if err := app.Run(); err != nil {
		t.Fatalf("Run() failed: %v", err)
	}
}

func TestEventProcessorAndRegisterEvent(t *testing.T) {
	RegisterEvent[string]("parity.event")

	processor := NewWailsEventProcessor(nil)
	called := false
	cancel := processor.On("parity.event", func(event *CustomEvent) {
		called = event.Data == "ok"
	})
	defer cancel()

	if err := processor.Emit(&CustomEvent{Name: "parity.event", Data: "ok"}); err != nil {
		t.Fatalf("Emit() failed: %v", err)
	}
	if !called {
		t.Fatalf("event listener was not called")
	}
	if err := processor.Emit(&CustomEvent{Name: "parity.event", Data: 7}); err == nil {
		t.Fatalf("expected registered type mismatch to fail")
	}
}

func TestMenuAndWindowParityHelpers(t *testing.T) {
	menu := NewMenu()
	menu.AddCheckbox("A", true)
	sub := menu.AddSubmenu("More")
	sub.Add("Child")
	if menu.FindByLabel("Child") == nil {
		t.Fatalf("FindByLabel did not traverse submenu")
	}

	window := NewWindow(WebviewWindowOptions{Name: "main", Width: 100, Height: 80})
	window.RegisterKeyBinding("CmdOrCtrl+K", func(Window) {})
	window.SetPhysicalBounds(Rect{X: 5, Y: 6, Width: 120, Height: 90})
	if got := window.PhysicalBounds(); got.Width != 120 || got.Height != 90 {
		t.Fatalf("unexpected physical bounds: %+v", got)
	}
	window.OpenDevTools()
	if !window.DevToolsOpen() {
		t.Fatalf("expected devtools to be marked open")
	}
	window.CloseDevTools()
	if window.DevToolsOpen() {
		t.Fatalf("expected devtools to be marked closed")
	}
}
