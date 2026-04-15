// pkg/mcp/tools_events.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
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
	r := s.core.Action("events.emit").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: events.TaskEmit{Name: input.Name, Data: input.Data}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventEmitOutput{}, e
		}
		return nil, EventEmitOutput{}, nil
	}
	cancelled, ok := r.Value.(bool)
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
	r := s.core.Action("events.on").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: events.TaskOn{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventOnOutput{}, e
		}
		return nil, EventOnOutput{}, nil
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
	r := s.core.Action("events.off").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: events.TaskOff{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventOffOutput{}, e
		}
		return nil, EventOffOutput{}, nil
	}
	return nil, EventOffOutput{Success: true}, nil
}

// --- event_list ---

type EventListInput struct{}
type EventListOutput struct {
	Listeners []events.ListenerInfo `json:"listeners"`
}

func (s *Subsystem) eventList(_ context.Context, _ *mcp.CallToolRequest, _ EventListInput) (*mcp.CallToolResult, EventListOutput, error) {
	r := s.core.QUERY(events.QueryListeners{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, EventListOutput{}, e
		}
		return nil, EventListOutput{}, nil
	}
	listenerInfos, ok := r.Value.([]events.ListenerInfo)
	if !ok {
		return nil, EventListOutput{}, coreerr.E("mcp.eventList", "unexpected result type", nil)
	}
	return nil, EventListOutput{Listeners: listenerInfos}, nil
}

// --- Registration ---

func (s *Subsystem) registerEventsTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "event_emit", Description: "Fire a named custom event with optional data"}, s.eventEmit)
	addTool(s, server, &mcp.Tool{Name: "event_on", Description: "Register a listener for a named custom event"}, s.eventOn)
	addTool(s, server, &mcp.Tool{Name: "event_off", Description: "Remove all listeners for a named custom event"}, s.eventOff)
	addTool(s, server, &mcp.Tool{Name: "event_list", Description: "Query all registered event listeners"}, s.eventList)
}
