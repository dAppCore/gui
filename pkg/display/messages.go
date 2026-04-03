// pkg/display/messages.go
package display

// ActionIDECommand is broadcast when a menu handler triggers an IDE command
// (save, run, build). Replaces direct s.app.Event().Emit("ide:*") calls.
// Listeners (e.g. editor windows) handle this via HandleIPCEvents.
// Use: _ = c.ACTION(display.ActionIDECommand{Command: "save"})
type ActionIDECommand struct {
	Command string `json:"command"` // "save", "run", "build"
}

// EventIDECommand is the WS event type for IDE commands.
// Use: eventType := display.EventIDECommand
const EventIDECommand EventType = "ide.command"

// Theme is the display-level theme summary exposed by the service API.
// Use: theme := display.Theme{IsDark: true}
type Theme struct {
	IsDark bool `json:"isDark"`
}
