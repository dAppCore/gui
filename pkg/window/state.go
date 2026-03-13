// pkg/window/state.go
package window

// StateManager persists window positions to disk.
// Full implementation in Task 3.
type StateManager struct{}

func NewStateManager() *StateManager { return &StateManager{} }

// ApplyState restores saved position/size to a Window descriptor.
func (sm *StateManager) ApplyState(w *Window) {}
