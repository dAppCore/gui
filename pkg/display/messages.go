// pkg/display/messages.go
package display

// ActionIDECommand is broadcast when a menu handler triggers an IDE command
// (save, run, build). Replaces direct s.app.Event().Emit("ide:*") calls.
// Listeners (e.g. editor windows) handle this via HandleIPCEvents.
type ActionIDECommand struct {
	Command string `json:"command"` // "save", "run", "build"
}

// EventIDECommand is the WS event type for IDE commands.
const EventIDECommand EventType = "ide.command"
