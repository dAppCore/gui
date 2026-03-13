// pkg/mcp/subsystem.go
package mcp

import (
	"forge.lthn.ai/core/go/pkg/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Subsystem implements the MCP Subsystem interface via structural typing.
// It registers GUI tools that translate MCP tool calls to IPC messages.
type Subsystem struct {
	core *core.Core
}

// New creates a display MCP subsystem backed by the given Core instance.
func New(c *core.Core) *Subsystem {
	return &Subsystem{core: c}
}

// Name returns the subsystem identifier.
func (s *Subsystem) Name() string { return "display" }

// RegisterTools registers all GUI tools with the MCP server.
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
