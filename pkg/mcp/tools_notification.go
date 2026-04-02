// pkg/mcp/tools_notification.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/notification"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- notification_show ---

type NotificationShowInput struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Subtitle string `json:"subtitle,omitempty"`
}
type NotificationShowOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) notificationShow(_ context.Context, _ *mcp.CallToolRequest, input NotificationShowInput) (*mcp.CallToolResult, NotificationShowOutput, error) {
	_, _, err := s.core.PERFORM(notification.TaskSend{Opts: notification.NotificationOptions{
		Title:    input.Title,
		Message:  input.Message,
		Subtitle: input.Subtitle,
	}})
	if err != nil {
		return nil, NotificationShowOutput{}, err
	}
	return nil, NotificationShowOutput{Success: true}, nil
}

// --- notification_with_actions ---

type NotificationWithActionsInput struct {
	Title    string                            `json:"title"`
	Message  string                            `json:"message"`
	Subtitle string                            `json:"subtitle,omitempty"`
	Actions  []notification.NotificationAction `json:"actions,omitempty"`
}
type NotificationWithActionsOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) notificationWithActions(_ context.Context, _ *mcp.CallToolRequest, input NotificationWithActionsInput) (*mcp.CallToolResult, NotificationWithActionsOutput, error) {
	_, _, err := s.core.PERFORM(notification.TaskSend{Opts: notification.NotificationOptions{
		Title:    input.Title,
		Message:  input.Message,
		Subtitle: input.Subtitle,
		Actions:  input.Actions,
	}})
	if err != nil {
		return nil, NotificationWithActionsOutput{}, err
	}
	return nil, NotificationWithActionsOutput{Success: true}, nil
}

// --- notification_permission_request ---

type NotificationPermissionRequestInput struct{}
type NotificationPermissionRequestOutput struct {
	Granted bool `json:"granted"`
}

func (s *Subsystem) notificationPermissionRequest(_ context.Context, _ *mcp.CallToolRequest, _ NotificationPermissionRequestInput) (*mcp.CallToolResult, NotificationPermissionRequestOutput, error) {
	result, _, err := s.core.PERFORM(notification.TaskRequestPermission{})
	if err != nil {
		return nil, NotificationPermissionRequestOutput{}, err
	}
	granted, ok := result.(bool)
	if !ok {
		return nil, NotificationPermissionRequestOutput{}, fmt.Errorf("unexpected result type from notification permission request")
	}
	return nil, NotificationPermissionRequestOutput{Granted: granted}, nil
}

// --- notification_permission_check ---

type NotificationPermissionCheckInput struct{}
type NotificationPermissionCheckOutput struct {
	Granted bool `json:"granted"`
}

func (s *Subsystem) notificationPermissionCheck(_ context.Context, _ *mcp.CallToolRequest, _ NotificationPermissionCheckInput) (*mcp.CallToolResult, NotificationPermissionCheckOutput, error) {
	result, _, err := s.core.QUERY(notification.QueryPermission{})
	if err != nil {
		return nil, NotificationPermissionCheckOutput{}, err
	}
	status, ok := result.(notification.PermissionStatus)
	if !ok {
		return nil, NotificationPermissionCheckOutput{}, fmt.Errorf("unexpected result type from notification permission check")
	}
	return nil, NotificationPermissionCheckOutput{Granted: status.Granted}, nil
}

// --- notification_clear ---

type NotificationClearInput struct{}
type NotificationClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) notificationClear(_ context.Context, _ *mcp.CallToolRequest, _ NotificationClearInput) (*mcp.CallToolResult, NotificationClearOutput, error) {
	_, _, err := s.core.PERFORM(notification.TaskClear{})
	if err != nil {
		return nil, NotificationClearOutput{}, err
	}
	return nil, NotificationClearOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerNotificationTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "notification_show", Description: "Show a desktop notification"}, s.notificationShow)
	mcp.AddTool(server, &mcp.Tool{Name: "notification_with_actions", Description: "Show a desktop notification with actions"}, s.notificationWithActions)
	mcp.AddTool(server, &mcp.Tool{Name: "notification_permission_request", Description: "Request notification permission"}, s.notificationPermissionRequest)
	mcp.AddTool(server, &mcp.Tool{Name: "notification_permission_check", Description: "Check notification permission status"}, s.notificationPermissionCheck)
	mcp.AddTool(server, &mcp.Tool{Name: "notification_clear", Description: "Clear notifications when supported"}, s.notificationClear)
}
