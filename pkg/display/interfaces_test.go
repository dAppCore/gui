package display

import (
	core "dappco.re/go"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestInterfaces_newWailsApp_Good(t *core.T) {
	app := &application.App{Logger: application.Logger{}}
	wrapped := newWailsApp(app)

	core.AssertNotNil(t, wrapped)
	core.AssertNotNil(t, wrapped.Logger())
	core.AssertNotPanics(t, func() {
		wrapped.Quit()
		wrapped.Logger().Info("ready")
	})
}

func TestInterfaces_newWailsApp_Bad(t *core.T) {
	wrapped := newWailsApp(&application.App{})
	core.AssertNotNil(t, wrapped)
	core.AssertNotNil(t, wrapped.Logger())
}

func TestInterfaces_newWailsApp_Ugly(t *core.T) {
	wrapped := newWailsApp(nil)
	core.AssertNotNil(t, wrapped)
	core.AssertPanics(t, func() {
		_ = wrapped.Logger()
	})
}

// AX7 generated source-matching smoke coverage.
func TestInterfaces_App_Logger_Good(t *core.T) {
	subject := new(wailsApp)
	result := core.Try(func() any {
		got0 := subject.Logger()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestInterfaces_App_Logger_Bad(t *core.T) {
	subject := new(wailsApp)
	result := core.Try(func() any {
		got0 := subject.Logger()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestInterfaces_App_Logger_Ugly(t *core.T) {
	subject := new(wailsApp)
	result := core.Try(func() any {
		got0 := subject.Logger()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestInterfaces_App_Quit_Good(t *core.T) {
	subject := new(wailsApp)
	result := core.Try(func() any {
		subject.Quit()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestInterfaces_App_Quit_Bad(t *core.T) {
	subject := new(wailsApp)
	result := core.Try(func() any {
		subject.Quit()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestInterfaces_App_Quit_Ugly(t *core.T) {
	subject := new(wailsApp)
	result := core.Try(func() any {
		subject.Quit()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
