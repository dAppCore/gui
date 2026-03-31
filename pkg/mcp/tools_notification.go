// pkg/mcp/tools_notification.go
package mcp

import (
	"context"

	coreerr "forge.lthn.ai/core/go-log"
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
	_, _, err := s.core.PERFORM(notification.TaskSend{Options: notification.NotificationOptions{
		Title:    input.Title,
		Message:  input.Message,
		Subtitle: input.Subtitle,
	}})
	if err != nil {
		return nil, NotificationShowOutput{}, err
	}
	return nil, NotificationShowOutput{Success: true}, nil
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
		return nil, NotificationPermissionRequestOutput{}, coreerr.E("mcp.notificationPermissionRequest", "unexpected result type", nil)
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
		return nil, NotificationPermissionCheckOutput{}, coreerr.E("mcp.notificationPermissionCheck", "unexpected result type", nil)
	}
	return nil, NotificationPermissionCheckOutput{Granted: status.Granted}, nil
}

// --- Registration ---

func (s *Subsystem) registerNotificationTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "notification_show", Description: "Show a desktop notification"}, s.notificationShow)
	mcp.AddTool(server, &mcp.Tool{Name: "notification_permission_request", Description: "Request notification permission"}, s.notificationPermissionRequest)
	mcp.AddTool(server, &mcp.Tool{Name: "notification_permission_check", Description: "Check notification permission status"}, s.notificationPermissionCheck)
}
