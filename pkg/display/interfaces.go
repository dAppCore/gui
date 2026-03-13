// pkg/display/interfaces.go
package display

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// App abstracts the Wails application for the orchestrator.
type App interface {
	Dialog() DialogManager
	Env() EnvManager
	Event() EventManager
	Logger() Logger
	Quit()
}

// DialogManager wraps Wails dialog operations.
type DialogManager interface {
	Info() *application.MessageDialog
	Warning() *application.MessageDialog
	OpenFile() *application.OpenFileDialogStruct
}

// EnvManager wraps Wails environment queries.
type EnvManager interface {
	Info() application.EnvironmentInfo
	IsDarkMode() bool
}

// EventManager wraps Wails application events.
type EventManager interface {
	OnApplicationEvent(eventType events.ApplicationEventType, handler func(*application.ApplicationEvent)) func()
	Emit(name string, data ...any) bool
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

func (w *wailsApp) Dialog() DialogManager { return &wailsDialogManager{app: w.app} }
func (w *wailsApp) Env() EnvManager       { return &wailsEnvManager{app: w.app} }
func (w *wailsApp) Event() EventManager   { return &wailsEventManager{app: w.app} }
func (w *wailsApp) Logger() Logger        { return w.app.Logger }
func (w *wailsApp) Quit()                 { w.app.Quit() }

type wailsDialogManager struct{ app *application.App }

func (d *wailsDialogManager) Info() *application.MessageDialog    { return d.app.Dialog.Info() }
func (d *wailsDialogManager) Warning() *application.MessageDialog { return d.app.Dialog.Warning() }
func (d *wailsDialogManager) OpenFile() *application.OpenFileDialogStruct {
	return d.app.Dialog.OpenFile()
}

type wailsEnvManager struct{ app *application.App }

func (e *wailsEnvManager) Info() application.EnvironmentInfo { return e.app.Env.Info() }
func (e *wailsEnvManager) IsDarkMode() bool                  { return e.app.Env.IsDarkMode() }

type wailsEventManager struct{ app *application.App }

func (ev *wailsEventManager) OnApplicationEvent(eventType events.ApplicationEventType, handler func(*application.ApplicationEvent)) func() {
	return ev.app.Event.OnApplicationEvent(eventType, handler)
}
func (ev *wailsEventManager) Emit(name string, data ...any) bool {
	return ev.app.Event.Emit(name, data...)
}

// wailsEventSource implements EventSource using a Wails app.
type wailsEventSource struct{ app *application.App }

func newWailsEventSource(app *application.App) EventSource {
	return &wailsEventSource{app: app}
}

func (es *wailsEventSource) OnThemeChange(handler func(isDark bool)) func() {
	return es.app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(_ *application.ApplicationEvent) {
		handler(es.app.Env.IsDarkMode())
	})
}

func (es *wailsEventSource) Emit(name string, data ...any) bool {
	return es.app.Event.Emit(name, data...)
}
