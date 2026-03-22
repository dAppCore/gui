// pkg/mcp/tools_clipboard.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/clipboard"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- clipboard_read ---

type ClipboardReadInput struct{}
type ClipboardReadOutput struct {
	Content string `json:"content"`
}

func (s *Subsystem) clipboardRead(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardReadInput) (*mcp.CallToolResult, ClipboardReadOutput, error) {
	result, _, err := s.core.QUERY(clipboard.QueryText{})
	if err != nil {
		return nil, ClipboardReadOutput{}, err
	}
	content, ok := result.(clipboard.ClipboardContent)
	if !ok {
		return nil, ClipboardReadOutput{}, fmt.Errorf("unexpected result type from clipboard read query")
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
	result, _, err := s.core.PERFORM(clipboard.TaskSetText{Text: input.Text})
	if err != nil {
		return nil, ClipboardWriteOutput{}, err
	}
	success, ok := result.(bool)
	if !ok {
		return nil, ClipboardWriteOutput{}, fmt.Errorf("unexpected result type from clipboard write task")
	}
	return nil, ClipboardWriteOutput{Success: success}, nil
}

// --- clipboard_has ---

type ClipboardHasInput struct{}
type ClipboardHasOutput struct {
	HasContent bool `json:"hasContent"`
}

func (s *Subsystem) clipboardHas(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardHasInput) (*mcp.CallToolResult, ClipboardHasOutput, error) {
	result, _, err := s.core.QUERY(clipboard.QueryText{})
	if err != nil {
		return nil, ClipboardHasOutput{}, err
	}
	content, ok := result.(clipboard.ClipboardContent)
	if !ok {
		return nil, ClipboardHasOutput{}, fmt.Errorf("unexpected result type from clipboard has query")
	}
	return nil, ClipboardHasOutput{HasContent: content.HasContent}, nil
}

// --- clipboard_clear ---

type ClipboardClearInput struct{}
type ClipboardClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardClear(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardClearInput) (*mcp.CallToolResult, ClipboardClearOutput, error) {
	result, _, err := s.core.PERFORM(clipboard.TaskClear{})
	if err != nil {
		return nil, ClipboardClearOutput{}, err
	}
	success, ok := result.(bool)
	if !ok {
		return nil, ClipboardClearOutput{}, fmt.Errorf("unexpected result type from clipboard clear task")
	}
	return nil, ClipboardClearOutput{Success: success}, nil
}

// --- Registration ---

func (s *Subsystem) registerClipboardTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_read", Description: "Read the current clipboard content"}, s.clipboardRead)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_write", Description: "Write text to the clipboard"}, s.clipboardWrite)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_has", Description: "Check if the clipboard has content"}, s.clipboardHas)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_clear", Description: "Clear the clipboard"}, s.clipboardClear)
}
