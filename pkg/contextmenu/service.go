// pkg/contextmenu/service.go
package contextmenu

import (
	"context"
	"sync"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	platform        Platform
	mu              sync.RWMutex
	registeredMenus map[string]ContextMenuDef
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().Action("contextmenu.add", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskAdd)
		return core.Result{Value: nil, OK: true}.New(s.taskAdd(t))
	})
	s.Core().Action("contextmenu.remove", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskRemove)
		return core.Result{Value: nil, OK: true}.New(s.taskRemove(t))
	})
	s.Core().Action("contextmenu.update", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskUpdate)
		return core.Result{Value: nil, OK: true}.New(s.taskUpdate(t))
	})
	s.Core().Action("contextmenu.destroy", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskDestroy)
		return core.Result{Value: nil, OK: true}.New(s.taskDestroy(t))
	})
	return core.Result{OK: true}
}

func (s *Service) OnShutdown(_ context.Context) core.Result {
	// Destroy all registered menus on shutdown to release platform resources
	s.mu.Lock()
	defer s.mu.Unlock()
	for name := range s.registeredMenus {
		_ = s.platform.Remove(name)
	}
	s.registeredMenus = make(map[string]ContextMenuDef)
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

// --- Query Handlers ---

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q := q.(type) {
	case QueryGet:
		return core.Result{Value: s.queryGet(q), OK: true}
	case QueryList:
		return core.Result{Value: s.queryList(), OK: true}
	case QueryGetAll:
		return core.Result{Value: s.queryList(), OK: true}
	default:
		return core.Result{}
	}
}

func (s *Service) queryGet(q QueryGet) *ContextMenuDef {
	s.mu.RLock()
	defer s.mu.RUnlock()
	menu, ok := s.registeredMenus[q.Name]
	if !ok {
		return nil
	}
	return &menu
}

func (s *Service) queryList() map[string]ContextMenuDef {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]ContextMenuDef, len(s.registeredMenus))
	for k, v := range s.registeredMenus {
		result[k] = v
	}
	return result
}

func (s *Service) taskAdd(t TaskAdd) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()
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

func (s *Service) taskUpdate(t TaskUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.registeredMenus[t.Name]; !exists {
		return ErrorMenuNotFound
	}

	// Re-register with updated definition — remove then add
	if err := s.platform.Remove(t.Name); err != nil {
		return coreerr.E("contextmenu.taskUpdate", "platform remove failed", err)
	}

	err := s.platform.Add(t.Name, t.Menu, func(menuName, actionID, data string) {
		_ = s.Core().ACTION(ActionItemClicked{
			MenuName: menuName,
			ActionID: actionID,
			Data:     data,
		})
	})
	if err != nil {
		return coreerr.E("contextmenu.taskUpdate", "platform add failed", err)
	}

	s.registeredMenus[t.Name] = t.Menu
	return nil
}

func (s *Service) taskDestroy(t TaskDestroy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.registeredMenus[t.Name]; !exists {
		return ErrorMenuNotFound
	}

	if err := s.platform.Remove(t.Name); err != nil {
		return coreerr.E("contextmenu.taskDestroy", "platform remove failed", err)
	}

	delete(s.registeredMenus, t.Name)
	return nil
}
