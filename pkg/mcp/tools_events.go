// pkg/mcp/tools_events.go
package mcp

import (
	"context"

	"dappco.re/go/core/gui/pkg/events"
	coreerr "forge.lthn.ai/core/go-log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- event_emit ---

type EventEmitInput struct {
	Name string `json:"name"`
	Data any    `json:"data,omitempty"`
}
type EventEmitOutput struct {
	Success bool `json:"success"`
}

// eventEmit fires a custom event by name with optional data.
// c.PERFORM(events.TaskEmit{Name: "build:done", Data: result})
func (s *Subsystem) eventEmit(_ context.Context, _ *mcp.CallToolRequest, input EventEmitInput) (*mcp.CallToolResult, EventEmitOutput, error) {
	_, _, err := s.core.PERFORM(events.TaskEmit{Name: input.Name, Data: input.Data})
	if err != nil {
		return nil, EventEmitOutput{}, err
	}
	return nil, EventEmitOutput{Success: true}, nil
}

// --- event_on ---

type EventOnInput struct {
	Name string `json:"name"`
}
type EventOnOutput struct {
	Success bool `json:"success"`
}

// eventOn registers a persistent listener for a named event.
// c.PERFORM(events.TaskOn{Name: "build:done"})
func (s *Subsystem) eventOn(_ context.Context, _ *mcp.CallToolRequest, input EventOnInput) (*mcp.CallToolResult, EventOnOutput, error) {
	_, _, err := s.core.PERFORM(events.TaskOn{Name: input.Name})
	if err != nil {
		return nil, EventOnOutput{}, err
	}
	return nil, EventOnOutput{Success: true}, nil
}

// --- event_off ---

type EventOffInput struct {
	Name string `json:"name"`
}
type EventOffOutput struct {
	Success bool `json:"success"`
}

// eventOff removes all listeners for a named event.
// c.PERFORM(events.TaskOff{Name: "build:done"})
func (s *Subsystem) eventOff(_ context.Context, _ *mcp.CallToolRequest, input EventOffInput) (*mcp.CallToolResult, EventOffOutput, error) {
	_, _, err := s.core.PERFORM(events.TaskOff{Name: input.Name})
	if err != nil {
		return nil, EventOffOutput{}, err
	}
	return nil, EventOffOutput{Success: true}, nil
}

// --- event_list ---

type EventListInput struct {
	Name string `json:"name"`
}
type EventListOutput struct {
	Count int `json:"count"`
}

// eventList returns the number of listeners registered for a named event.
// count := c.QUERY(events.QueryListeners{Name: "build:done"})
func (s *Subsystem) eventList(_ context.Context, _ *mcp.CallToolRequest, input EventListInput) (*mcp.CallToolResult, EventListOutput, error) {
	result, _, err := s.core.QUERY(events.QueryListeners{Name: input.Name})
	if err != nil {
		return nil, EventListOutput{}, err
	}
	count, ok := result.(int)
	if !ok {
		return nil, EventListOutput{}, coreerr.E("mcp.eventList", "unexpected result type", nil)
	}
	return nil, EventListOutput{Count: count}, nil
}

// --- Registration ---

func (s *Subsystem) registerEventTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "event_emit", Description: "Fire a custom event by name with optional data"}, s.eventEmit)
	mcp.AddTool(server, &mcp.Tool{Name: "event_on", Description: "Register a persistent listener for a named event"}, s.eventOn)
	mcp.AddTool(server, &mcp.Tool{Name: "event_off", Description: "Remove all listeners for a named event"}, s.eventOff)
	mcp.AddTool(server, &mcp.Tool{Name: "event_list", Description: "Return the number of listeners registered for a named event"}, s.eventList)
}
