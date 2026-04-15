// pkg/mcp/tools_notification.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
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
	r := s.core.Action("notification.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: notification.TaskSend{Options: notification.NotificationOptions{
			Title:    input.Title,
			Message:  input.Message,
			Subtitle: input.Subtitle,
		}}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, NotificationShowOutput{}, e
		}
		return nil, NotificationShowOutput{}, nil
	}
	return nil, NotificationShowOutput{Success: true}, nil
}

// --- notification_permission_request ---

type NotificationPermissionRequestInput struct{}
type NotificationPermissionRequestOutput struct {
	Granted bool `json:"granted"`
}

func (s *Subsystem) notificationPermissionRequest(_ context.Context, _ *mcp.CallToolRequest, _ NotificationPermissionRequestInput) (*mcp.CallToolResult, NotificationPermissionRequestOutput, error) {
	r := s.core.Action("notification.requestPermission").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, NotificationPermissionRequestOutput{}, e
		}
		return nil, NotificationPermissionRequestOutput{}, nil
	}
	granted, ok := r.Value.(bool)
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
	r := s.core.QUERY(notification.QueryPermission{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, NotificationPermissionCheckOutput{}, e
		}
		return nil, NotificationPermissionCheckOutput{}, nil
	}
	status, ok := r.Value.(notification.PermissionStatus)
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
