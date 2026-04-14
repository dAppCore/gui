// pkg/window/tiling.go
package window

import coreerr "forge.lthn.ai/core/go-log"

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

var snapPositionNames = map[SnapPosition]string{
	SnapLeft: "left", SnapRight: "right",
	SnapTop: "top", SnapBottom: "bottom",
	SnapTopLeft: "top-left", SnapTopRight: "top-right",
	SnapBottomLeft: "bottom-left", SnapBottomRight: "bottom-right",
	SnapCenter: "center",
}

func (p SnapPosition) String() string { return snapPositionNames[p] }

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

func layoutOrigin(origin []int) (int, int) {
	if len(origin) == 0 {
		return 0, 0
	}
	if len(origin) == 1 {
		return origin[0], 0
	}
	return origin[0], origin[1]
}

func (m *Manager) captureState(pw PlatformWindow) {
	if m.state == nil || pw == nil {
		return
	}
	m.state.CaptureState(pw)
}

func normalizeWindowForLayout(pw PlatformWindow) {
	if pw == nil {
		return
	}
	if pw.IsMaximised() || pw.IsMinimised() {
		pw.Restore()
	}
}

// TileWindows arranges the named windows in the given mode across the screen area.
func (m *Manager) TileWindows(mode TileMode, names []string, screenW, screenH int, origin ...int) error {
	originX, originY := layoutOrigin(origin)
	windows := make([]PlatformWindow, 0, len(names))
	for _, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return coreerr.E("window.Manager.TileWindows", "window not found: "+name, nil)
		}
		windows = append(windows, pw)
	}
	if len(windows) == 0 {
		return coreerr.E("window.Manager.TileWindows", "no windows to tile", nil)
	}

	halfW, halfH := screenW/2, screenH/2

	switch mode {
	case TileModeLeftRight:
		w := screenW / len(windows)
		for i, pw := range windows {
			pw.SetPosition(originX+i*w, originY)
			pw.SetSize(w, screenH)
			m.captureState(pw)
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
			pw.SetPosition(originX+col*cellW, originY+row*cellH)
			pw.SetSize(cellW, cellH)
			m.captureState(pw)
		}
	case TileModeLeftHalf:
		for _, pw := range windows {
			pw.SetPosition(originX, originY)
			pw.SetSize(halfW, screenH)
			m.captureState(pw)
		}
	case TileModeRightHalf:
		for _, pw := range windows {
			pw.SetPosition(originX+halfW, originY)
			pw.SetSize(halfW, screenH)
			m.captureState(pw)
		}
	case TileModeTopHalf:
		for _, pw := range windows {
			pw.SetPosition(originX, originY)
			pw.SetSize(screenW, halfH)
			m.captureState(pw)
		}
	case TileModeBottomHalf:
		for _, pw := range windows {
			pw.SetPosition(originX, originY+halfH)
			pw.SetSize(screenW, halfH)
			m.captureState(pw)
		}
	case TileModeTopLeft:
		for _, pw := range windows {
			pw.SetPosition(originX, originY)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	case TileModeTopRight:
		for _, pw := range windows {
			pw.SetPosition(originX+halfW, originY)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	case TileModeBottomLeft:
		for _, pw := range windows {
			pw.SetPosition(originX, originY+halfH)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	case TileModeBottomRight:
		for _, pw := range windows {
			pw.SetPosition(originX+halfW, originY+halfH)
			pw.SetSize(halfW, halfH)
			m.captureState(pw)
		}
	}
	return nil
}

// SnapWindow snaps a window to a screen edge/corner/centre.
func (m *Manager) SnapWindow(name string, pos SnapPosition, screenW, screenH int, origin ...int) error {
	originX, originY := layoutOrigin(origin)
	pw, ok := m.Get(name)
	if !ok {
		return coreerr.E("window.Manager.SnapWindow", "window not found: "+name, nil)
	}

	halfW, halfH := screenW/2, screenH/2

	switch pos {
	case SnapLeft:
		pw.SetPosition(originX, originY)
		pw.SetSize(halfW, screenH)
	case SnapRight:
		pw.SetPosition(originX+halfW, originY)
		pw.SetSize(halfW, screenH)
	case SnapTop:
		pw.SetPosition(originX, originY)
		pw.SetSize(screenW, halfH)
	case SnapBottom:
		pw.SetPosition(originX, originY+halfH)
		pw.SetSize(screenW, halfH)
	case SnapTopLeft:
		pw.SetPosition(originX, originY)
		pw.SetSize(halfW, halfH)
	case SnapTopRight:
		pw.SetPosition(originX+halfW, originY)
		pw.SetSize(halfW, halfH)
	case SnapBottomLeft:
		pw.SetPosition(originX, originY+halfH)
		pw.SetSize(halfW, halfH)
	case SnapBottomRight:
		pw.SetPosition(originX+halfW, originY+halfH)
		pw.SetSize(halfW, halfH)
	case SnapCenter:
		normalizeWindowForLayout(pw)
		cw, ch := pw.Size()
		pw.SetPosition(originX+(screenW-cw)/2, originY+(screenH-ch)/2)
	}
	m.captureState(pw)
	return nil
}

// StackWindows cascades windows with an offset.
func (m *Manager) StackWindows(names []string, offsetX, offsetY int, origin ...int) error {
	originX, originY := layoutOrigin(origin)
	for i, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return coreerr.E("window.Manager.StackWindows", "window not found: "+name, nil)
		}
		pw.SetPosition(originX+i*offsetX, originY+i*offsetY)
		m.captureState(pw)
	}
	return nil
}

// ApplyWorkflow arranges windows in a predefined workflow layout.
func (m *Manager) ApplyWorkflow(workflow WorkflowLayout, names []string, screenW, screenH int, origin ...int) error {
	originX, originY := layoutOrigin(origin)
	if len(names) == 0 {
		return coreerr.E("window.Manager.ApplyWorkflow", "no windows for workflow", nil)
	}

	switch workflow {
	case WorkflowCoding:
		// 70/30 split — main editor + terminal
		mainW := screenW * 70 / 100
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(originX, originY)
			pw.SetSize(mainW, screenH)
			m.captureState(pw)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				pw.SetPosition(originX+mainW, originY)
				pw.SetSize(screenW-mainW, screenH)
				m.captureState(pw)
			}
		}
	case WorkflowDebugging:
		// 60/40 split
		mainW := screenW * 60 / 100
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(originX, originY)
			pw.SetSize(mainW, screenH)
			m.captureState(pw)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				pw.SetPosition(originX+mainW, originY)
				pw.SetSize(screenW-mainW, screenH)
				m.captureState(pw)
			}
		}
	case WorkflowPresenting:
		// Maximise first window
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(originX, originY)
			pw.SetSize(screenW, screenH)
			m.captureState(pw)
		}
	case WorkflowSideBySide:
		return m.TileWindows(TileModeLeftRight, names, screenW, screenH, originX, originY)
	}
	return nil
}
