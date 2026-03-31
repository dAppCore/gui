// pkg/mcp/subsystem.go
package mcp

import (
	"forge.lthn.ai/core/go/pkg/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Subsystem translates MCP tool calls to Core IPC messages for GUI operations.
type Subsystem struct {
	core *core.Core
}

// New(c) creates a display MCP subsystem backed by a Core instance.
// sub := mcp.New(c); sub.RegisterTools(server)
func New(c *core.Core) *Subsystem {
	return &Subsystem{core: c}
}

func (s *Subsystem) Name() string { return "display" }

func (s *Subsystem) RegisterTools(server *mcp.Server) {
	s.registerWebviewTools(server)
	s.registerWindowTools(server)
	s.registerLayoutTools(server)
	s.registerScreenTools(server)
	s.registerClipboardTools(server)
	s.registerDialogTools(server)
	s.registerNotificationTools(server)
	s.registerTrayTools(server)
	s.registerEnvironmentTools(server)
	s.registerBrowserTools(server)
	s.registerContextMenuTools(server)
	s.registerKeybindingTools(server)
	s.registerDockTools(server)
	s.registerLifecycleTools(server)
}
