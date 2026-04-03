// pkg/window/tiling.go
package window

import "fmt"

// normalizeWindowForLayout clears transient maximise/minimise state before
// applying a new geometry. This keeps layout helpers effective even when a
// window was previously maximised.
func normalizeWindowForLayout(pw PlatformWindow) {
	if pw == nil {
		return
	}
	if pw.IsMaximised() || pw.IsMinimised() {
		pw.Restore()
	}
}

// TileMode defines how windows are arranged.
// Use: mode := window.TileModeLeftRight
type TileMode int

const (
	TileModeLeftHalf TileMode = iota
	TileModeRightHalf
	TileModeTopHalf
	TileModeBottomHalf
	TileModeTopLeft
	TileModeTopRight
	TileModeBottomLeft
	TileModeBottomRight
	TileModeLeftRight
	TileModeGrid
)

var tileModeNames = map[TileMode]string{
	TileModeLeftHalf: "left-half", TileModeRightHalf: "right-half",
	TileModeTopHalf: "top-half", TileModeBottomHalf: "bottom-half",
	TileModeTopLeft: "top-left", TileModeTopRight: "top-right",
	TileModeBottomLeft: "bottom-left", TileModeBottomRight: "bottom-right",
	TileModeLeftRight: "left-right", TileModeGrid: "grid",
}

// String returns the canonical layout name for the tile mode.
// Use: label := window.TileModeGrid.String()
func (m TileMode) String() string { return tileModeNames[m] }

// SnapPosition defines where a window snaps to.
// Use: pos := window.SnapRight
type SnapPosition int

const (
	SnapLeft SnapPosition = iota
	SnapRight
	SnapTop
	SnapBottom
	SnapTopLeft
	SnapTopRight
	SnapBottomLeft
	SnapBottomRight
	SnapCenter
)

// WorkflowLayout is a predefined arrangement for common tasks.
// Use: workflow := window.WorkflowCoding
type WorkflowLayout int

const (
	WorkflowCoding     WorkflowLayout = iota // 70/30 split
	WorkflowDebugging                        // 60/40 split
	WorkflowPresenting                       // maximised
	WorkflowSideBySide                       // 50/50 split
)

var workflowNames = map[WorkflowLayout]string{
	WorkflowCoding: "coding", WorkflowDebugging: "debugging",
	WorkflowPresenting: "presenting", WorkflowSideBySide: "side-by-side",
}

// String returns the canonical workflow name.
// Use: label := window.WorkflowCoding.String()
func (w WorkflowLayout) String() string { return workflowNames[w] }

// ParseWorkflowLayout converts a workflow name into its enum value.
// Use: workflow, ok := window.ParseWorkflowLayout("coding")
func ParseWorkflowLayout(name string) (WorkflowLayout, bool) {
	for workflow, workflowName := range workflowNames {
		if workflowName == name {
			return workflow, true
		}
	}
	return WorkflowCoding, false
}

// TileWindows arranges the named windows in the given mode across the screen area.
func (m *Manager) TileWindows(mode TileMode, names []string, screenW, screenH int) error {
	windows := make([]PlatformWindow, 0, len(names))
	for _, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return fmt.Errorf("window %q not found", name)
		}
		windows = append(windows, pw)
	}
	if len(windows) == 0 {
		return fmt.Errorf("no windows to tile")
	}

	halfW, halfH := screenW/2, screenH/2

	switch mode {
	case TileModeLeftRight:
		w := screenW / len(windows)
		for i, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(i*w, 0)
			pw.SetSize(w, screenH)
		}
	case TileModeGrid:
		cols := 2
		if len(windows) > 4 {
			cols = 3
		}
		cellW := screenW / cols
		for i, pw := range windows {
			normalizeWindowForLayout(pw)
			row := i / cols
			col := i % cols
			rows := (len(windows) + cols - 1) / cols
			cellH := screenH / rows
			pw.SetPosition(col*cellW, row*cellH)
			pw.SetSize(cellW, cellH)
		}
	case TileModeLeftHalf:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, 0)
			pw.SetSize(halfW, screenH)
		}
	case TileModeRightHalf:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(halfW, 0)
			pw.SetSize(halfW, screenH)
		}
	case TileModeTopHalf:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, 0)
			pw.SetSize(screenW, halfH)
		}
	case TileModeBottomHalf:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, halfH)
			pw.SetSize(screenW, halfH)
		}
	case TileModeTopLeft:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, 0)
			pw.SetSize(halfW, halfH)
		}
	case TileModeTopRight:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(halfW, 0)
			pw.SetSize(halfW, halfH)
		}
	case TileModeBottomLeft:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, halfH)
			pw.SetSize(halfW, halfH)
		}
	case TileModeBottomRight:
		for _, pw := range windows {
			normalizeWindowForLayout(pw)
			pw.SetPosition(halfW, halfH)
			pw.SetSize(halfW, halfH)
		}
	}
	return nil
}

// SnapWindow snaps a window to a screen edge/corner/centre.
func (m *Manager) SnapWindow(name string, pos SnapPosition, screenW, screenH int) error {
	pw, ok := m.Get(name)
	if !ok {
		return fmt.Errorf("window %q not found", name)
	}

	halfW, halfH := screenW/2, screenH/2

	switch pos {
	case SnapLeft:
		normalizeWindowForLayout(pw)
		pw.SetPosition(0, 0)
		pw.SetSize(halfW, screenH)
	case SnapRight:
		normalizeWindowForLayout(pw)
		pw.SetPosition(halfW, 0)
		pw.SetSize(halfW, screenH)
	case SnapTop:
		normalizeWindowForLayout(pw)
		pw.SetPosition(0, 0)
		pw.SetSize(screenW, halfH)
	case SnapBottom:
		normalizeWindowForLayout(pw)
		pw.SetPosition(0, halfH)
		pw.SetSize(screenW, halfH)
	case SnapTopLeft:
		normalizeWindowForLayout(pw)
		pw.SetPosition(0, 0)
		pw.SetSize(halfW, halfH)
	case SnapTopRight:
		normalizeWindowForLayout(pw)
		pw.SetPosition(halfW, 0)
		pw.SetSize(halfW, halfH)
	case SnapBottomLeft:
		normalizeWindowForLayout(pw)
		pw.SetPosition(0, halfH)
		pw.SetSize(halfW, halfH)
	case SnapBottomRight:
		normalizeWindowForLayout(pw)
		pw.SetPosition(halfW, halfH)
		pw.SetSize(halfW, halfH)
	case SnapCenter:
		normalizeWindowForLayout(pw)
		cw, ch := pw.Size()
		pw.SetPosition((screenW-cw)/2, (screenH-ch)/2)
	}
	return nil
}

// StackWindows cascades windows with an offset.
func (m *Manager) StackWindows(names []string, offsetX, offsetY int) error {
	for i, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return fmt.Errorf("window %q not found", name)
		}
		normalizeWindowForLayout(pw)
		pw.SetPosition(i*offsetX, i*offsetY)
	}
	return nil
}

// ApplyWorkflow arranges windows in a predefined workflow layout.
func (m *Manager) ApplyWorkflow(workflow WorkflowLayout, names []string, screenW, screenH int) error {
	if len(names) == 0 {
		return fmt.Errorf("no windows for workflow")
	}

	switch workflow {
	case WorkflowCoding:
		// 70/30 split — main editor + terminal
		mainW := screenW * 70 / 100
		if pw, ok := m.Get(names[0]); ok {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, 0)
			pw.SetSize(mainW, screenH)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				normalizeWindowForLayout(pw)
				pw.SetPosition(mainW, 0)
				pw.SetSize(screenW-mainW, screenH)
			}
		}
	case WorkflowDebugging:
		// 60/40 split
		mainW := screenW * 60 / 100
		if pw, ok := m.Get(names[0]); ok {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, 0)
			pw.SetSize(mainW, screenH)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				normalizeWindowForLayout(pw)
				pw.SetPosition(mainW, 0)
				pw.SetSize(screenW-mainW, screenH)
			}
		}
	case WorkflowPresenting:
		// Maximise first window
		if pw, ok := m.Get(names[0]); ok {
			normalizeWindowForLayout(pw)
			pw.SetPosition(0, 0)
			pw.SetSize(screenW, screenH)
		}
	case WorkflowSideBySide:
		return m.TileWindows(TileModeLeftRight, names, screenW, screenH)
	}
	return nil
}
