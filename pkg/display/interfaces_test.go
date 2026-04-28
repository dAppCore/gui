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
