package window

import (
	core "dappco.re/go"
	"os"
	"path/filepath"
	"time"
)

func TestLayoutManager_SaveLayout_Good(t *core.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	windows := map[string]WindowState{
		"editor":   {X: 0, Y: 0, Width: 960, Height: 1080},
		"terminal": {X: 960, Y: 0, Width: 960, Height: 540},
	}

	core.RequireNoError(t, lm.SaveLayout("coding", windows))

	layout, ok := lm.GetLayout("coding")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "coding", layout.Name)
	core.AssertLen(t, layout.Windows, 2)
	core.AssertNotEmpty(t, layout.CreatedAt)
	core.AssertNotEmpty(t, layout.UpdatedAt)

	infos := lm.ListLayouts()
	core.AssertLen(t, infos, 1)
	core.AssertEqual(t, "coding", infos[0].Name)
	core.AssertEqual(t, 2, infos[0].WindowCount)

	lm.DeleteLayout("coding")
	_, ok = lm.GetLayout("coding")
	core.AssertFalse(t, ok)
}

func TestLayoutManager_SaveLayout_Bad(t *core.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	err := lm.SaveLayout("", map[string]WindowState{"main": {Width: 1}})

	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "layout name cannot be empty")
}

func TestLayoutManager_SaveLayout_Ugly(t *core.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	core.RequireNoError(t, lm.SaveLayout("coding", map[string]WindowState{"main": {Width: 800}}))
	first, ok := lm.GetLayout("coding")
	core.RequireTrue(t, ok)

	time.Sleep(2 * time.Millisecond)
	core.RequireNoError(t, lm.SaveLayout("coding", map[string]WindowState{"main": {Width: 1024}}))
	second, ok := lm.GetLayout("coding")
	core.RequireTrue(t, ok)

	core.AssertEqual(t, first.CreatedAt, second.CreatedAt)
	core.AssertGreater(t, second.UpdatedAt, first.UpdatedAt)
	core.AssertEqual(t, 1024, second.Windows["main"].Width)
}

func TestLayoutManager_NewLayoutManagerWithPathEnv_Good(t *core.T) {
	path := filepath.Join(t.TempDir(), "custom", "layouts.json")
	t.Setenv(layoutFileEnv, path)

	lm := NewLayoutManager()

	core.AssertNotNil(t, lm)
	core.AssertEqual(t, path, lm.filePath())
	core.AssertEqual(t, filepath.Dir(path), lm.dataDir())

	core.RequireNoError(t, lm.SaveLayout("coding", map[string]WindowState{
		"main": {Width: 800, Height: 600},
	}))

	content, err := os.ReadFile(path)
	core.RequireNoError(t, err)
	core.AssertContains(t, string(content), `"coding"`)
}
