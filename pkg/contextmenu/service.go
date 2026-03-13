// pkg/contextmenu/service.go
package contextmenu

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the context menu service.
type Options struct{}

// Service is a core.Service managing context menus via IPC.
// It maintains an in-memory registry of menus (map[string]ContextMenuDef)
// and delegates platform-level registration to the Platform interface.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	menus    map[string]ContextMenuDef
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

// --- Query Handlers ---

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryGet:
		return s.queryGet(q), true, nil
	case QueryList:
		return s.queryList(), true, nil
	default:
		return nil, false, nil
	}
}

// queryGet returns a single menu definition by name, or nil if not found.
func (s *Service) queryGet(q QueryGet) *ContextMenuDef {
	menu, ok := s.menus[q.Name]
	if !ok {
		return nil
	}
	return &menu
}

// queryList returns a copy of all registered menus.
func (s *Service) queryList() map[string]ContextMenuDef {
	result := make(map[string]ContextMenuDef, len(s.menus))
	for k, v := range s.menus {
		result[k] = v
	}
	return result
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskAdd:
		return nil, true, s.taskAdd(t)
	case TaskRemove:
		return nil, true, s.taskRemove(t)
	default:
		return nil, false, nil
	}
}

func (s *Service) taskAdd(t TaskAdd) error {
	// If menu already exists, remove it first (replace semantics)
	if _, exists := s.menus[t.Name]; exists {
		_ = s.platform.Remove(t.Name)
		delete(s.menus, t.Name)
	}

	// Register on platform with a callback that broadcasts ActionItemClicked
	err := s.platform.Add(t.Name, t.Menu, func(menuName, actionID, data string) {
		_ = s.Core().ACTION(ActionItemClicked{
			MenuName: menuName,
			ActionID: actionID,
			Data:     data,
		})
	})
	if err != nil {
		return fmt.Errorf("contextmenu: platform add failed: %w", err)
	}

	s.menus[t.Name] = t.Menu
	return nil
}

func (s *Service) taskRemove(t TaskRemove) error {
	if _, exists := s.menus[t.Name]; !exists {
		return ErrMenuNotFound
	}

	err := s.platform.Remove(t.Name)
	if err != nil {
		return fmt.Errorf("contextmenu: platform remove failed: %w", err)
	}

	delete(s.menus, t.Name)
	return nil
}
