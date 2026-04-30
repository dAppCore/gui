package application

func newKeyBindingManager() *KeyBindingManager {
	return &KeyBindingManager{}
}

type windowKeyEvent struct {
	windowId          uint
	acceleratorString string
}

func (m *KeyBindingManager) Register(accelerator string, callback func(window Window)) {
	m.Add(accelerator, callback)
}

func (m *KeyBindingManager) Unregister(accelerator string) {
	m.Remove(accelerator)
}

func (m *KeyBindingManager) handleWindowKeyEvent(event *windowKeyEvent) {
	if event == nil {
		return
	}
	m.Process(event.acceleratorString, nil)
}
