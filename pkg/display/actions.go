// pkg/display/actions.go
package display

import "forge.lthn.ai/core/gui/pkg/window"

// ActionOpenWindow is an IPC message type requesting a new window.
type ActionOpenWindow struct {
	window.Window
}
