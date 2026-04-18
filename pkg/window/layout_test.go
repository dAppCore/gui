package window

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayoutManager_SaveLayout_Good(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	windows := map[string]WindowState{
		"editor":   {X: 0, Y: 0, Width: 960, Height: 1080},
		"terminal": {X: 960, Y: 0, Width: 960, Height: 540},
	}

	require.NoError(t, lm.SaveLayout("coding", windows))

	layout, ok := lm.GetLayout("coding")
	require.True(t, ok)
	assert.Equal(t, "coding", layout.Name)
	assert.Len(t, layout.Windows, 2)
	assert.NotZero(t, layout.CreatedAt)
	assert.NotZero(t, layout.UpdatedAt)

	infos := lm.ListLayouts()
	require.Len(t, infos, 1)
	assert.Equal(t, "coding", infos[0].Name)
	assert.Equal(t, 2, infos[0].WindowCount)

	lm.DeleteLayout("coding")
	_, ok = lm.GetLayout("coding")
	assert.False(t, ok)
}

func TestLayoutManager_SaveLayout_Bad(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	err := lm.SaveLayout("", map[string]WindowState{"main": {Width: 1}})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "layout name cannot be empty")
}

func TestLayoutManager_SaveLayout_Ugly(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	require.NoError(t, lm.SaveLayout("coding", map[string]WindowState{"main": {Width: 800}}))
	first, ok := lm.GetLayout("coding")
	require.True(t, ok)

	time.Sleep(2 * time.Millisecond)
	require.NoError(t, lm.SaveLayout("coding", map[string]WindowState{"main": {Width: 1024}}))
	second, ok := lm.GetLayout("coding")
	require.True(t, ok)

	assert.Equal(t, first.CreatedAt, second.CreatedAt)
	assert.Greater(t, second.UpdatedAt, first.UpdatedAt)
	assert.Equal(t, 1024, second.Windows["main"].Width)
}

func TestLayoutManager_NewLayoutManagerWithPathEnv_Good(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom", "layouts.json")
	t.Setenv(layoutFileEnv, path)

	lm := NewLayoutManager()

	require.NotNil(t, lm)
	assert.Equal(t, path, lm.filePath())
	assert.Equal(t, filepath.Dir(path), lm.dataDir())

	require.NoError(t, lm.SaveLayout("coding", map[string]WindowState{
		"main": {Width: 800, Height: 600},
	}))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"coding"`)
}
