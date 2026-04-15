package application

import (
	"context"
	"log/slog"
	"testing"
)

type noopTransport struct{}

func (noopTransport) Start(context.Context, *MessageProcessor) error { return nil }
func (noopTransport) Stop() error                                    { return nil }

func TestNewAndGetExposeSingletonState(t *testing.T) {
	globalApplication = nil
	logger := slog.Default()
	transport := noopTransport{}
	app := New(Options{
		Name:      "Parity",
		Logger:    logger,
		LogLevel:  slog.LevelDebug,
		Transport: transport,
	})
	if Get() != app {
		t.Fatalf("Get() did not return singleton app")
	}
	if app.Config().Name != "Parity" {
		t.Fatalf("unexpected app name: %q", app.Config().Name)
	}
	if app.Config().Logger != logger {
		t.Fatalf("expected logger to round-trip through config")
	}
	if app.Config().LogLevel != slog.LevelDebug {
		t.Fatalf("expected log level to round-trip through config")
	}
	if app.Config().Transport == nil {
		t.Fatalf("expected transport to round-trip through config")
	}
	if app.Logger == nil {
		t.Fatalf("expected app logger to be initialised")
	}
	if app.Env == nil {
		t.Fatalf("expected Env manager to be initialised")
	}
	if app.Environment != app.Env {
		t.Fatalf("expected Environment alias to point at Env manager")
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

func TestEnvironmentAndDropTargetParity(t *testing.T) {
	info := (&EnvironmentManager{}).Info()
	if info.OS == "" || info.Arch == "" {
		t.Fatalf("expected runtime environment metadata, got %+v", info)
	}
	if info.OSInfo == nil {
		t.Fatalf("expected OSInfo to be populated")
	}

	ctx := &WindowEventContext{
		dropDetails: &DropTargetDetails{
			X:         10,
			Y:         20,
			ElementID: "dropzone",
			ClassList: []string{"primary", "drop-target"},
			Attributes: map[string]string{
				"data-file-drop-target": "true",
			},
		},
	}
	details := ctx.DropTargetDetails()
	if details == nil {
		t.Fatalf("expected drop target details")
	}
	if details.X != 10 || details.Y != 20 || details.ElementID != "dropzone" {
		t.Fatalf("unexpected drop details: %+v", details)
	}
	if len(details.ClassList) != 2 || details.Attributes["data-file-drop-target"] != "true" {
		t.Fatalf("drop target metadata was not preserved: %+v", details)
	}
}
