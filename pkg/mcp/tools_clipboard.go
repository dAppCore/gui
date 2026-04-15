// pkg/mcp/tools_clipboard.go
package mcp

import (
	"context"
	"encoding/base64"

	"dappco.re/go/core/gui/pkg/clipboard"
	coreerr "forge.lthn.ai/core/go-log"
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
	result, _, err := s.core.PERFORM(clipboard.TaskSetText{Text: input.Text})
	if err != nil {
		return nil, ClipboardWriteOutput{}, err
	}
	success, ok := result.(bool)
	if !ok {
		return nil, ClipboardWriteOutput{}, coreerr.E("mcp.clipboardWrite", "unexpected result type", nil)
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
	result, _, err := s.core.PERFORM(clipboard.TaskClear{})
	if err != nil {
		return nil, ClipboardClearOutput{}, err
	}
	success, ok := result.(bool)
	if !ok {
		return nil, ClipboardClearOutput{}, coreerr.E("mcp.clipboardClear", "unexpected result type", nil)
	}
	return nil, ClipboardClearOutput{Success: success}, nil
}

// --- clipboard_read_image ---

type ClipboardReadImageInput struct{}
type ClipboardReadImageOutput struct {
	Base64 string `json:"base64"`
}

func (s *Subsystem) clipboardReadImage(_ context.Context, _ *mcp.CallToolRequest, _ ClipboardReadImageInput) (*mcp.CallToolResult, ClipboardReadImageOutput, error) {
	result, _, err := s.core.QUERY(clipboard.QueryImage{})
	if err != nil {
		return nil, ClipboardReadImageOutput{}, err
	}
	content, ok := result.(clipboard.ImageContent)
	if !ok {
		return nil, ClipboardReadImageOutput{}, coreerr.E("mcp.clipboardReadImage", "unexpected result type", nil)
	}
	if !content.HasImage {
		return nil, ClipboardReadImageOutput{}, nil
	}
	return nil, ClipboardReadImageOutput{Base64: base64.StdEncoding.EncodeToString(content.Data)}, nil
}

// --- clipboard_write_image ---

type ClipboardWriteImageInput struct {
	Base64 string `json:"base64"`
}
type ClipboardWriteImageOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) clipboardWriteImage(_ context.Context, _ *mcp.CallToolRequest, input ClipboardWriteImageInput) (*mcp.CallToolResult, ClipboardWriteImageOutput, error) {
	data, err := base64.StdEncoding.DecodeString(input.Base64)
	if err != nil {
		return nil, ClipboardWriteImageOutput{}, err
	}
	result, _, err := s.core.PERFORM(clipboard.TaskSetImage{Data: data})
	if err != nil {
		return nil, ClipboardWriteImageOutput{}, err
	}
	success, ok := result.(bool)
	if !ok {
		return nil, ClipboardWriteImageOutput{}, coreerr.E("mcp.clipboardWriteImage", "unexpected result type", nil)
	}
	return nil, ClipboardWriteImageOutput{Success: success}, nil
}

// --- Registration ---

func (s *Subsystem) registerClipboardTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_read", Description: "Read the current clipboard content"}, s.clipboardRead)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_write", Description: "Write text to the clipboard"}, s.clipboardWrite)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_has", Description: "Check if the clipboard has content"}, s.clipboardHas)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_read_image", Description: "Read image data from the clipboard as base64"}, s.clipboardReadImage)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_write_image", Description: "Write base64 image data to the clipboard"}, s.clipboardWriteImage)
	mcp.AddTool(server, &mcp.Tool{Name: "clipboard_clear", Description: "Clear the clipboard"}, s.clipboardClear)
}
