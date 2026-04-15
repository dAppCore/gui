// pkg/window/window_test.go
package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowDefaults_Good(t *testing.T) {
	w := &Window{}
	assert.Equal(t, "", w.Name)
	assert.Equal(t, 0, w.Width)
}

func TestWindowOption_Name_Good(t *testing.T) {
	w := &Window{}
	err := WithName("main")(w)
	require.NoError(t, err)
	assert.Equal(t, "main", w.Name)
}

func TestWindowOption_Title_Good(t *testing.T) {
	w := &Window{}
	err := WithTitle("My App")(w)
	require.NoError(t, err)
	assert.Equal(t, "My App", w.Title)
}

func TestWindowOption_URL_Good(t *testing.T) {
	w := &Window{}
	err := WithURL("/dashboard")(w)
	require.NoError(t, err)
	assert.Equal(t, "/dashboard", w.URL)
}

func TestWindowOption_Size_Good(t *testing.T) {
	w := &Window{}
	err := WithSize(1280, 720)(w)
	require.NoError(t, err)
	assert.Equal(t, 1280, w.Width)
	assert.Equal(t, 720, w.Height)
}

func TestWindowOption_Position_Good(t *testing.T) {
	w := &Window{}
	err := WithPosition(100, 200)(w)
	require.NoError(t, err)
	assert.Equal(t, 100, w.X)
	assert.Equal(t, 200, w.Y)
}

func TestApplyOptions_Good(t *testing.T) {
	w, err := ApplyOptions(
		WithName("test"),
		WithTitle("Test Window"),
		WithURL("/test"),
		WithSize(800, 600),
	)
	require.NoError(t, err)
	assert.Equal(t, "test", w.Name)
	assert.Equal(t, "Test Window", w.Title)
	assert.Equal(t, "/test", w.URL)
	assert.Equal(t, 800, w.Width)
	assert.Equal(t, 600, w.Height)
}

func TestApplyOptions_Bad(t *testing.T) {
	_, err := ApplyOptions(func(w *Window) error {
		return assert.AnError
	})
	assert.Error(t, err)
}

func TestApplyOptions_Empty_Good(t *testing.T) {
	w, err := ApplyOptions()
	require.NoError(t, err)
	assert.NotNil(t, w)
}

// newTestManager creates a Manager with a mock platform and clean state for testing.
func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	m := &Manager{
		platform: p,
		state:    &StateManager{states: make(map[string]WindowState)},
		layout:   &LayoutManager{layouts: make(map[string]Layout)},
		windows:  make(map[string]PlatformWindow),
	}
	return m, p
}

func TestManager_Open_Good(t *testing.T) {
	m, p := newTestManager()
	pw, err := m.Open(WithName("test"), WithTitle("Test"), WithURL("/test"), WithSize(800, 600))
	require.NoError(t, err)
	assert.NotNil(t, pw)
	assert.Equal(t, "test", pw.Name())
	assert.Len(t, p.windows, 1)
}

func TestManager_Open_Defaults_Good(t *testing.T) {
	m, _ := newTestManager()
	pw, err := m.Open()
	require.NoError(t, err)
	assert.Equal(t, "main", pw.Name())
	w, h := pw.Size()
	assert.Equal(t, 1280, w)
	assert.Equal(t, 800, h)
}

func TestManager_Open_CustomDefaults_Good(t *testing.T) {
	m, _ := newTestManager()
	m.SetDefaultWidth(1440)
	m.SetDefaultHeight(900)

	pw, err := m.Open()
	require.NoError(t, err)

	w, h := pw.Size()
	assert.Equal(t, 1440, w)
	assert.Equal(t, 900, h)
}

func TestManager_Open_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, err := m.Open(func(w *Window) error { return assert.AnError })
	assert.Error(t, err)
}

func TestManager_Get_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("findme"))
	pw, ok := m.Get("findme")
	assert.True(t, ok)
	assert.Equal(t, "findme", pw.Name())
}

func TestManager_Get_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, ok := m.Get("nonexistent")
	assert.False(t, ok)
}

func TestManager_List_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"))
	_, _ = m.Open(WithName("b"))
	names := m.List()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "a")
	assert.Contains(t, names, "b")
}

func TestManager_Remove_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("temp"))
	m.Remove("temp")
	_, ok := m.Get("temp")
	assert.False(t, ok)
}

// --- Tiling Tests ---

func TestTileMode_String_Good(t *testing.T) {
	assert.Equal(t, "left-half", TileModeLeftHalf.String())
	assert.Equal(t, "grid", TileModeGrid.String())
}

func TestManager_TileWindows_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"), WithSize(800, 600))
	_, _ = m.Open(WithName("b"), WithSize(800, 600))
	err := m.TileWindows(TileModeLeftRight, []string{"a", "b"}, 1920, 1080)
	require.NoError(t, err)
	a, _ := m.Get("a")
	b, _ := m.Get("b")
	aw, _ := a.Size()
	bw, _ := b.Size()
	assert.Equal(t, 960, aw)
	assert.Equal(t, 960, bw)
}

func TestManager_TileWindows_Bad(t *testing.T) {
	m, _ := newTestManager()
	err := m.TileWindows(TileModeLeftRight, []string{"nonexistent"}, 1920, 1080)
	assert.Error(t, err)
}

func TestManager_SnapWindow_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("snap"), WithSize(800, 600))
	err := m.SnapWindow("snap", SnapLeft, 1920, 1080)
	require.NoError(t, err)
	w, _ := m.Get("snap")
	x, _ := w.Position()
	assert.Equal(t, 0, x)
	sw, _ := w.Size()
	assert.Equal(t, 960, sw)
}

func TestManager_StackWindows_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("s1"), WithSize(800, 600))
	_, _ = m.Open(WithName("s2"), WithSize(800, 600))
	err := m.StackWindows([]string{"s1", "s2"}, 30, 30)
	require.NoError(t, err)
	s2, _ := m.Get("s2")
	x, y := s2.Position()
	assert.Equal(t, 30, x)
	assert.Equal(t, 30, y)
}

func TestWorkflowLayout_Good(t *testing.T) {
	assert.Equal(t, "coding", WorkflowCoding.String())
	assert.Equal(t, "debugging", WorkflowDebugging.String())
}

// --- Comprehensive Tiling Tests ---

func TestTileWindows_AllModes_Good(t *testing.T) {
	const screenW, screenH = 1920, 1080
	halfW, halfH := screenW/2, screenH/2

	tests := []struct {
		name       string
		mode       TileMode
		wantX      int
		wantY      int
		wantWidth  int
		wantHeight int
	}{
		{"LeftHalf", TileModeLeftHalf, 0, 0, halfW, screenH},
		{"RightHalf", TileModeRightHalf, halfW, 0, halfW, screenH},
		{"TopHalf", TileModeTopHalf, 0, 0, screenW, halfH},
		{"BottomHalf", TileModeBottomHalf, 0, halfH, screenW, halfH},
		{"TopLeft", TileModeTopLeft, 0, 0, halfW, halfH},
		{"TopRight", TileModeTopRight, halfW, 0, halfW, halfH},
		{"BottomLeft", TileModeBottomLeft, 0, halfH, halfW, halfH},
		{"BottomRight", TileModeBottomRight, halfW, halfH, halfW, halfH},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, _ := newTestManager()
			_, err := m.Open(WithName("win"), WithSize(800, 600))
			require.NoError(t, err)

			err = m.TileWindows(tc.mode, []string{"win"}, screenW, screenH)
			require.NoError(t, err)

			pw, ok := m.Get("win")
			require.True(t, ok)

			x, y := pw.Position()
			w, h := pw.Size()
			assert.Equal(t, tc.wantX, x, "x position")
			assert.Equal(t, tc.wantY, y, "y position")
			assert.Equal(t, tc.wantWidth, w, "width")
			assert.Equal(t, tc.wantHeight, h, "height")
		})
	}
}

func TestSnapWindow_AllPositions_Good(t *testing.T) {
	const screenW, screenH = 1920, 1080
	halfW, halfH := screenW/2, screenH/2

	tests := []struct {
		name       string
		pos        SnapPosition
		initW      int
		initH      int
		wantX      int
		wantY      int
		wantWidth  int
		wantHeight int
	}{
		{"Right", SnapRight, 800, 600, halfW, 0, halfW, screenH},
		{"Top", SnapTop, 800, 600, 0, 0, screenW, halfH},
		{"Bottom", SnapBottom, 800, 600, 0, halfH, screenW, halfH},
		{"TopLeft", SnapTopLeft, 800, 600, 0, 0, halfW, halfH},
		{"TopRight", SnapTopRight, 800, 600, halfW, 0, halfW, halfH},
		{"BottomLeft", SnapBottomLeft, 800, 600, 0, halfH, halfW, halfH},
		{"BottomRight", SnapBottomRight, 800, 600, halfW, halfH, halfW, halfH},
		{"Center", SnapCenter, 800, 600, (screenW - 800) / 2, (screenH - 600) / 2, 800, 600},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, _ := newTestManager()
			_, err := m.Open(WithName("snap"), WithSize(tc.initW, tc.initH))
			require.NoError(t, err)

			err = m.SnapWindow("snap", tc.pos, screenW, screenH)
			require.NoError(t, err)

			pw, ok := m.Get("snap")
			require.True(t, ok)

			x, y := pw.Position()
			w, h := pw.Size()
			assert.Equal(t, tc.wantX, x, "x position")
			assert.Equal(t, tc.wantY, y, "y position")
			assert.Equal(t, tc.wantWidth, w, "width")
			assert.Equal(t, tc.wantHeight, h, "height")
		})
	}
}

func TestStackWindows_ThreeWindows_Good(t *testing.T) {
	m, _ := newTestManager()
	names := []string{"s1", "s2", "s3"}
	for _, name := range names {
		_, err := m.Open(WithName(name), WithSize(800, 600))
		require.NoError(t, err)
	}

	err := m.StackWindows(names, 30, 30)
	require.NoError(t, err)

	for i, name := range names {
		pw, ok := m.Get(name)
		require.True(t, ok, "window %s should exist", name)
		x, y := pw.Position()
		assert.Equal(t, i*30, x, "window %s x position", name)
		assert.Equal(t, i*30, y, "window %s y position", name)
	}
}

func TestApplyWorkflow_AllLayouts_Good(t *testing.T) {
	const screenW, screenH = 1920, 1080

	tests := []struct {
		name     string
		workflow WorkflowLayout
		// Expected positions/sizes for the first two windows.
		// For WorkflowSideBySide, TileWindows(LeftRight) divides equally.
		win0X, win0Y, win0W, win0H int
		win1X, win1Y, win1W, win1H int
	}{
		{
			"Coding",
			WorkflowCoding,
			0, 0, 1344, screenH, // 70% of 1920 = 1344
			1344, 0, screenW - 1344, screenH, // remaining 30%
		},
		{
			"Debugging",
			WorkflowDebugging,
			0, 0, 1152, screenH, // 60% of 1920 = 1152
			1152, 0, screenW - 1152, screenH, // remaining 40%
		},
		{
			"Presenting",
			WorkflowPresenting,
			0, 0, screenW, screenH, // maximised
			0, 0, 800, 600, // second window untouched
		},
		{
			"SideBySide",
			WorkflowSideBySide,
			0, 0, 960, screenH, // left half (1920/2)
			960, 0, 960, screenH, // right half
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, _ := newTestManager()
			_, err := m.Open(WithName("editor"), WithSize(800, 600))
			require.NoError(t, err)
			_, err = m.Open(WithName("terminal"), WithSize(800, 600))
			require.NoError(t, err)

			err = m.ApplyWorkflow(tc.workflow, []string{"editor", "terminal"}, screenW, screenH)
			require.NoError(t, err)

			pw0, ok := m.Get("editor")
			require.True(t, ok)
			x0, y0 := pw0.Position()
			w0, h0 := pw0.Size()
			assert.Equal(t, tc.win0X, x0, "editor x")
			assert.Equal(t, tc.win0Y, y0, "editor y")
			assert.Equal(t, tc.win0W, w0, "editor width")
			assert.Equal(t, tc.win0H, h0, "editor height")

			pw1, ok := m.Get("terminal")
			require.True(t, ok)
			x1, y1 := pw1.Position()
			w1, h1 := pw1.Size()
			assert.Equal(t, tc.win1X, x1, "terminal x")
			assert.Equal(t, tc.win1Y, y1, "terminal y")
			assert.Equal(t, tc.win1W, w1, "terminal width")
			assert.Equal(t, tc.win1H, h1, "terminal height")
		})
	}
}

func TestApplyWorkflow_Empty_Bad(t *testing.T) {
	m, _ := newTestManager()
	err := m.ApplyWorkflow(WorkflowCoding, []string{}, 1920, 1080)
	assert.Error(t, err)
}
