// pkg/contextmenu/service.go
package contextmenu

import (
	"context"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform        Platform
	registeredMenus map[string]ContextMenuDef
}

func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

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
	case QueryGetAll:
		menus := make([]ContextMenuDef, 0, len(s.registeredMenus))
		for _, menu := range s.registeredMenus {
			menus = append(menus, menu)
		}
		return menus, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) queryGet(q QueryGet) *ContextMenuDef {
	menu, ok := s.registeredMenus[q.Name]
	if !ok {
		return nil
	}
	return &menu
}

func (s *Service) queryList() map[string]ContextMenuDef {
	result := make(map[string]ContextMenuDef, len(s.registeredMenus))
	for k, v := range s.registeredMenus {
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
	case TaskUpdate:
		if _, exists := s.registeredMenus[t.Name]; !exists {
			return nil, true, ErrorMenuNotFound
		}
		_ = s.platform.Remove(t.Name)
		delete(s.registeredMenus, t.Name)
		return nil, true, s.taskAdd(TaskAdd{Name: t.Name, Menu: t.Menu})
	case TaskDestroy:
		if _, exists := s.registeredMenus[t.Name]; !exists {
			return nil, true, ErrorMenuNotFound
		}
		_ = s.platform.Remove(t.Name)
		delete(s.registeredMenus, t.Name)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) taskAdd(t TaskAdd) error {
	// If menu already exists, remove it first (replace semantics)
	if _, exists := s.registeredMenus[t.Name]; exists {
		_ = s.platform.Remove(t.Name)
		delete(s.registeredMenus, t.Name)
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
		return coreerr.E("contextmenu.taskAdd", "platform add failed", err)
	}

	s.registeredMenus[t.Name] = t.Menu
	return nil
}

func (s *Service) taskRemove(t TaskRemove) error {
	if _, exists := s.registeredMenus[t.Name]; !exists {
		return ErrorMenuNotFound
	}

	err := s.platform.Remove(t.Name)
	if err != nil {
		return coreerr.E("contextmenu.taskRemove", "platform remove failed", err)
	}

	delete(s.registeredMenus, t.Name)
	return nil
}
