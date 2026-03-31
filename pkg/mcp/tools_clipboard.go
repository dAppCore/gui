// pkg/mcp/tools_clipboard.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
	"forge.lthn.ai/core/gui/pkg/clipboard"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- clipboard_read ---

type ClipboardReadInput struct{}
type ClipboardReadOutput struct {
	Content string `json:"content"`
}

func (s *Subsystem) clipboardRead(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardReadInput) (*mcp.CallToolResult, ClipboardReadOutput, error) {
	r := s.core.QUERY(clipboard.QueryText{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardReadOutput{}, e
		}
		return nil, ClipboardReadOutput{}, nil
	}
	content, ok := r.Value.(clipboard.ClipboardContent)
	if !ok {
		return nil, ClipboardReadOutput{}, coreerr.E("mcp.clipboardRead", "unexpected result type", nil)
	}
	return nil, ClipboardReadOutput{Content: content.Text}, nil
}

// --- clipboard_write ---

type ClipboardWriteInput struct {
	Text string `json:"text"`
}
type ClipboardWriteOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardWrite(_ context.Context, _ *mcp.CallToolRequest, input ClipboardWriteInput) (*mcp.CallToolResult, ClipboardWriteOutput, error) {
	r := s.core.Action("clipboard.setText").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: clipboard.TaskSetText{Text: input.Text}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardWriteOutput{}, e
		}
		return nil, ClipboardWriteOutput{}, nil
	}
	return nil, ClipboardWriteOutput{Success: true}, nil
}

// --- clipboard_has ---

type ClipboardHasInput struct{}
type ClipboardHasOutput struct {
	HasContent bool `json:"hasContent"`
}

func (s *Subsystem) clipboardHas(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardHasInput) (*mcp.CallToolResult, ClipboardHasOutput, error) {
	r := s.core.QUERY(clipboard.QueryText{})
	if !r.OK {
		return nil, ClipboardHasOutput{}, nil
	}
	content, ok := r.Value.(clipboard.ClipboardContent)
	if !ok {
		return nil, ClipboardHasOutput{}, coreerr.E("mcp.clipboardHas", "unexpected result type", nil)
	}
	return nil, ClipboardHasOutput{HasContent: content.HasContent}, nil
}

// --- clipboard_clear ---

type ClipboardClearInput struct{}
type ClipboardClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardClear(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardClearInput) (*mcp.CallToolResult, ClipboardClearOutput, error) {
	r := s.core.Action("clipboard.clear").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ClipboardClearOutput{}, e
		}
		return nil, ClipboardClearOutput{}, nil
	}
	return nil, ClipboardClearOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerClipboardTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_read", Description: "Read the current clipboard content"}, s.clipboardRead)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_write", Description: "Write text to the clipboard"}, s.clipboardWrite)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_has", Description: "Check if the clipboard has content"}, s.clipboardHas)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_clear", Description: "Clear the clipboard"}, s.clipboardClear)
}
