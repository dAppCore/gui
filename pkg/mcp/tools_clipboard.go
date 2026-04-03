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

// --- clipboard_read_image ---

type ClipboardReadImageInput struct{}
type ClipboardReadImageOutput struct {
	Image clipboard.ClipboardImageContent `json:"image"`
}

func (s *Subsystem) clipboardReadImage(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardReadImageInput) (*mcp.CallToolResult, ClipboardReadImageOutput, error) {
	result, _, err := s.core.QUERY(clipboard.QueryImage{})
	if err != nil {
		return nil, ClipboardReadImageOutput{}, err
	}
	image, ok := result.(clipboard.ClipboardImageContent)
	if !ok {
		return nil, ClipboardReadImageOutput{}, fmt.Errorf("unexpected result type from clipboard image query")
	}
	return nil, ClipboardReadImageOutput{Image: image}, nil
}

// --- clipboard_write_image ---

type ClipboardWriteImageInput struct {
	Data []byte `json:"data"`
}
type ClipboardWriteImageOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardWriteImage(_ context.Context, _ *mcp.CallToolRequest, input ClipboardWriteImageInput) (*mcp.CallToolResult, ClipboardWriteImageOutput, error) {
	_, _, err := s.core.PERFORM(clipboard.TaskSetImage{Data: input.Data})
	if err != nil {
		return nil, ClipboardWriteImageOutput{}, err
	}
	return nil, ClipboardWriteImageOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerClipboardTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_read", Description: "Read the current clipboard content"}, s.clipboardRead)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_write", Description: "Write text to the clipboard"}, s.clipboardWrite)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_has", Description: "Check if the clipboard has content"}, s.clipboardHas)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_clear", Description: "Clear the clipboard"}, s.clipboardClear)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_read_image", Description: "Read an image from the clipboard"}, s.clipboardReadImage)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_write_image", Description: "Write an image to the clipboard"}, s.clipboardWriteImage)
}
