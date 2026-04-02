// pkg/display/interfaces.go
package display

import "github.com/wailsapp/wails/v3/pkg/application"

// App abstracts the Wails application for the display orchestrator.
// The service uses Logger() for diagnostics and Quit() for shutdown.
type App interface {
	Logger() Logger
	Quit()
}

// Logger wraps Wails logging.
type Logger interface {
	Info(message string, args ...any)
}

// wailsApp wraps *application.App for the App interface.
type wailsApp struct {
	app *application.App
}

func newWailsApp(app *application.App) App {
	return &wailsApp{app: app}
}

func (w *wailsApp) Logger() Logger { return w.app.Logger }
func (w *wailsApp) Quit()          { w.app.Quit() }
