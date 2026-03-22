// pkg/mcp/tools_dialog.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/dialog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- dialog_open_file ---

type DialogOpenFileInput struct {
	Title         string              `json:"title,omitempty"`
	Directory     string              `json:"directory,omitempty"`
	Filters       []dialog.FileFilter `json:"filters,omitempty"`
	AllowMultiple bool                `json:"allowMultiple,omitempty"`
}
type DialogOpenFileOutput struct {
	Paths []string `json:"paths"`
}

func (s *Subsystem) dialogOpenFile(_ context.Context, _ *mcp.CallToolRequest, input DialogOpenFileInput) (*mcp.CallToolResult, DialogOpenFileOutput, error) {
	result, _, err := s.core.PERFORM(dialog.TaskOpenFile{Opts: dialog.OpenFileOptions{
		Title:         input.Title,
		Directory:     input.Directory,
		Filters:       input.Filters,
		AllowMultiple: input.AllowMultiple,
	}})
	if err != nil {
		return nil, DialogOpenFileOutput{}, err
	}
	paths, ok := result.([]string)
	if !ok {
		return nil, DialogOpenFileOutput{}, fmt.Errorf("unexpected result type from open file dialog")
	}
	return nil, DialogOpenFileOutput{Paths: paths}, nil
}

// --- dialog_save_file ---

type DialogSaveFileInput struct {
	Title     string              `json:"title,omitempty"`
	Directory string              `json:"directory,omitempty"`
	Filename  string              `json:"filename,omitempty"`
	Filters   []dialog.FileFilter `json:"filters,omitempty"`
}
type DialogSaveFileOutput struct {
	Path string `json:"path"`
}

func (s *Subsystem) dialogSaveFile(_ context.Context, _ *mcp.CallToolRequest, input DialogSaveFileInput) (*mcp.CallToolResult, DialogSaveFileOutput, error) {
	result, _, err := s.core.PERFORM(dialog.TaskSaveFile{Opts: dialog.SaveFileOptions{
		Title:     input.Title,
		Directory: input.Directory,
		Filename:  input.Filename,
		Filters:   input.Filters,
	}})
	if err != nil {
		return nil, DialogSaveFileOutput{}, err
	}
	path, ok := result.(string)
	if !ok {
		return nil, DialogSaveFileOutput{}, fmt.Errorf("unexpected result type from save file dialog")
	}
	return nil, DialogSaveFileOutput{Path: path}, nil
}

// --- dialog_open_directory ---

type DialogOpenDirectoryInput struct {
	Title     string `json:"title,omitempty"`
	Directory string `json:"directory,omitempty"`
}
type DialogOpenDirectoryOutput struct {
	Path string `json:"path"`
}

func (s *Subsystem) dialogOpenDirectory(_ context.Context, _ *mcp.CallToolRequest, input DialogOpenDirectoryInput) (*mcp.CallToolResult, DialogOpenDirectoryOutput, error) {
	result, _, err := s.core.PERFORM(dialog.TaskOpenDirectory{Opts: dialog.OpenDirectoryOptions{
		Title:     input.Title,
		Directory: input.Directory,
	}})
	if err != nil {
		return nil, DialogOpenDirectoryOutput{}, err
	}
	path, ok := result.(string)
	if !ok {
		return nil, DialogOpenDirectoryOutput{}, fmt.Errorf("unexpected result type from open directory dialog")
	}
	return nil, DialogOpenDirectoryOutput{Path: path}, nil
}

// --- dialog_confirm ---

type DialogConfirmInput struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Buttons []string `json:"buttons,omitempty"`
}
type DialogConfirmOutput struct {
	Button string `json:"button"`
}

func (s *Subsystem) dialogConfirm(_ context.Context, _ *mcp.CallToolRequest, input DialogConfirmInput) (*mcp.CallToolResult, DialogConfirmOutput, error) {
	result, _, err := s.core.PERFORM(dialog.TaskMessageDialog{Opts: dialog.MessageDialogOptions{
		Type:    dialog.DialogQuestion,
		Title:   input.Title,
		Message: input.Message,
		Buttons: input.Buttons,
	}})
	if err != nil {
		return nil, DialogConfirmOutput{}, err
	}
	button, ok := result.(string)
	if !ok {
		return nil, DialogConfirmOutput{}, fmt.Errorf("unexpected result type from confirm dialog")
	}
	return nil, DialogConfirmOutput{Button: button}, nil
}

// --- dialog_prompt ---

type DialogPromptInput struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
type DialogPromptOutput struct {
	Button string `json:"button"`
}

func (s *Subsystem) dialogPrompt(_ context.Context, _ *mcp.CallToolRequest, input DialogPromptInput) (*mcp.CallToolResult, DialogPromptOutput, error) {
	result, _, err := s.core.PERFORM(dialog.TaskMessageDialog{Opts: dialog.MessageDialogOptions{
		Type:    dialog.DialogInfo,
		Title:   input.Title,
		Message: input.Message,
		Buttons: []string{"OK", "Cancel"},
	}})
	if err != nil {
		return nil, DialogPromptOutput{}, err
	}
	button, ok := result.(string)
	if !ok {
		return nil, DialogPromptOutput{}, fmt.Errorf("unexpected result type from prompt dialog")
	}
	return nil, DialogPromptOutput{Button: button}, nil
}

// --- Registration ---

func (s *Subsystem) registerDialogTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "dialog_open_file", Description: "Show an open file dialog"}, s.dialogOpenFile)
	mcp.AddTool(server, &mcp.Tool{Name: "dialog_save_file", Description: "Show a save file dialog"}, s.dialogSaveFile)
	mcp.AddTool(server, &mcp.Tool{Name: "dialog_open_directory", Description: "Show a directory picker dialog"}, s.dialogOpenDirectory)
	mcp.AddTool(server, &mcp.Tool{Name: "dialog_confirm", Description: "Show a confirmation dialog"}, s.dialogConfirm)
	mcp.AddTool(server, &mcp.Tool{Name: "dialog_prompt", Description: "Show a prompt dialog"}, s.dialogPrompt)
}
