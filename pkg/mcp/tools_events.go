// pkg/mcp/tools_events.go
package mcp

import (
	"context"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/gui/pkg/events"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- event_emit ---

type EventEmitInput struct {
	Name string `json:"name"`
	Data any    `json:"data,omitempty"`
}
type EventEmitOutput struct {
	Cancelled bool `json:"cancelled"`
}

func (s *Subsystem) eventEmit(_ context.Context, _ *mcp.CallToolRequest, input EventEmitInput) (*mcp.CallToolResult, EventEmitOutput, error) {
	// result := s.core.PERFORM(events.TaskEmit{Name: "user:login", Data: payload})
	result, _, err := s.core.PERFORM(events.TaskEmit{Name: input.Name, Data: input.Data})
	if err != nil {
		return nil, EventEmitOutput{}, err
	}
	cancelled, ok := result.(bool)
	if !ok {
		return nil, EventEmitOutput{}, coreerr.E("mcp.eventEmit", "unexpected result type", nil)
	}
	return nil, EventEmitOutput{Cancelled: cancelled}, nil
}

// --- event_on ---

type EventOnInput struct {
	Name string `json:"name"`
}
type EventOnOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) eventOn(_ context.Context, _ *mcp.CallToolRequest, input EventOnInput) (*mcp.CallToolResult, EventOnOutput, error) {
	// s.core.PERFORM(events.TaskOn{Name: "theme:changed"})
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

func (s *Subsystem) eventOff(_ context.Context, _ *mcp.CallToolRequest, input EventOffInput) (*mcp.CallToolResult, EventOffOutput, error) {
	// s.core.PERFORM(events.TaskOff{Name: "theme:changed"})
	_, _, err := s.core.PERFORM(events.TaskOff{Name: input.Name})
	if err != nil {
		return nil, EventOffOutput{}, err
	}
	return nil, EventOffOutput{Success: true}, nil
}

// --- event_list ---

type EventListInput struct{}
type EventListOutput struct {
	Listeners []events.ListenerInfo `json:"listeners"`
}

func (s *Subsystem) eventList(_ context.Context, _ *mcp.CallToolRequest, _ EventListInput) (*mcp.CallToolResult, EventListOutput, error) {
	// result, _, _ := s.core.QUERY(events.QueryListeners{})
	result, _, err := s.core.QUERY(events.QueryListeners{})
	if err != nil {
		return nil, EventListOutput{}, err
	}
	listenerInfos, ok := result.([]events.ListenerInfo)
	if !ok {
		return nil, EventListOutput{}, coreerr.E("mcp.eventList", "unexpected result type", nil)
	}
	return nil, EventListOutput{Listeners: listenerInfos}, nil
}

// --- Registration ---

func (s *Subsystem) registerEventsTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "event_emit", Description: "Fire a named custom event with optional data"}, s.eventEmit)
	mcp.AddTool(server, &mcp.Tool{Name: "event_on", Description: "Register a listener for a named custom event"}, s.eventOn)
	mcp.AddTool(server, &mcp.Tool{Name: "event_off", Description: "Remove all listeners for a named custom event"}, s.eventOff)
	mcp.AddTool(server, &mcp.Tool{Name: "event_list", Description: "Query all registered event listeners"}, s.eventList)
}
