package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func applyWindowOptions(t *testing.T, options ...WindowOption) *Window {
	t.Helper()
	w, err := ApplyOptions(options...)
	require.NoError(t, err)
	return w
}

func TestOptions_WindowOptionSetters_Good(t *testing.T) {
	w := applyWindowOptions(t,
		WithName("main"),
		WithTitle("Core GUI"),
		WithURL("/dashboard"),
		WithHTML("<main>Ready</main>"),
		WithJS("globalThis.__CORE_READY__ = true"),
		WithSize(1280, 800),
		WithPosition(160, 120),
		WithMinSize(640, 480),
		WithMaxSize(1920, 1080),
		WithFrameless(true),
		WithHidden(false),
		WithAlwaysOnTop(true),
		WithBackgroundColour(12, 34, 56, 78),
		WithFileDrop(true),
	)

	assert.Equal(t, "main", w.Name)
	assert.Equal(t, "Core GUI", w.Title)
	assert.Equal(t, "/dashboard", w.URL)
	assert.Equal(t, "<main>Ready</main>", w.HTML)
	assert.Equal(t, "globalThis.__CORE_READY__ = true", w.JS)
	assert.Equal(t, 1280, w.Width)
	assert.Equal(t, 800, w.Height)
	assert.Equal(t, 160, w.X)
	assert.Equal(t, 120, w.Y)
	assert.Equal(t, 640, w.MinWidth)
	assert.Equal(t, 480, w.MinHeight)
	assert.Equal(t, 1920, w.MaxWidth)
	assert.Equal(t, 1080, w.MaxHeight)
	assert.True(t, w.Frameless)
	assert.False(t, w.Hidden)
	assert.True(t, w.AlwaysOnTop)
	assert.Equal(t, [4]uint8{12, 34, 56, 78}, w.BackgroundColour)
	assert.True(t, w.EnableFileDrop)
}

func TestOptions_WindowOptionSetters_Bad(t *testing.T) {
	w := applyWindowOptions(t,
		WithName(""),
		WithTitle(""),
		WithURL(""),
		WithHTML(""),
		WithJS(""),
		WithSize(0, 0),
		WithPosition(0, 0),
		WithMinSize(0, 0),
		WithMaxSize(0, 0),
		WithFrameless(false),
		WithHidden(false),
		WithAlwaysOnTop(false),
		WithBackgroundColour(0, 0, 0, 0),
		WithFileDrop(false),
	)

	assert.Equal(t, "", w.Name)
	assert.Equal(t, "", w.Title)
	assert.Equal(t, "", w.URL)
	assert.Equal(t, "", w.HTML)
	assert.Equal(t, "", w.JS)
	assert.Equal(t, 0, w.Width)
	assert.Equal(t, 0, w.Height)
	assert.Equal(t, 0, w.X)
	assert.Equal(t, 0, w.Y)
	assert.Equal(t, 0, w.MinWidth)
	assert.Equal(t, 0, w.MinHeight)
	assert.Equal(t, 0, w.MaxWidth)
	assert.Equal(t, 0, w.MaxHeight)
	assert.False(t, w.Frameless)
	assert.False(t, w.Hidden)
	assert.False(t, w.AlwaysOnTop)
	assert.Equal(t, [4]uint8{0, 0, 0, 0}, w.BackgroundColour)
	assert.False(t, w.EnableFileDrop)
}

func TestOptions_WindowOptionSetters_Ugly(t *testing.T) {
	w := applyWindowOptions(t,
		WithName("⚙︎core-window"),
		WithTitle("A very long title that stays intact"),
		WithURL("core://settings?tab=%F0%9F%93%81"),
		WithHTML("<section data-id=\"αβγ\">unsafe-looking but literal</section>"),
		WithJS("globalThis.__CORE_STATE__ = { mode: 'worker', value: -1 };"),
		WithSize(-1920, -1080),
		WithPosition(-42, 99999),
		WithMinSize(-1, -2),
		WithMaxSize(32767, 32767),
		WithFrameless(true),
		WithHidden(true),
		WithAlwaysOnTop(true),
		WithBackgroundColour(255, 254, 253, 252),
		WithFileDrop(true),
	)

	assert.Equal(t, "⚙︎core-window", w.Name)
	assert.Equal(t, "A very long title that stays intact", w.Title)
	assert.Equal(t, "core://settings?tab=%F0%9F%93%81", w.URL)
	assert.Equal(t, "<section data-id=\"αβγ\">unsafe-looking but literal</section>", w.HTML)
	assert.Equal(t, "globalThis.__CORE_STATE__ = { mode: 'worker', value: -1 };", w.JS)
	assert.Equal(t, -1920, w.Width)
	assert.Equal(t, -1080, w.Height)
	assert.Equal(t, -42, w.X)
	assert.Equal(t, 99999, w.Y)
	assert.Equal(t, -1, w.MinWidth)
	assert.Equal(t, -2, w.MinHeight)
	assert.Equal(t, 32767, w.MaxWidth)
	assert.Equal(t, 32767, w.MaxHeight)
	assert.True(t, w.Frameless)
	assert.True(t, w.Hidden)
	assert.True(t, w.AlwaysOnTop)
	assert.Equal(t, [4]uint8{255, 254, 253, 252}, w.BackgroundColour)
	assert.True(t, w.EnableFileDrop)
}

func TestOptions_ApplyOptions_Good(t *testing.T) {
	w, err := ApplyOptions(
		nil,
		WithName("main"),
		WithTitle("Core"),
	)

	require.NoError(t, err)
	require.NotNil(t, w)
	assert.Equal(t, "main", w.Name)
	assert.Equal(t, "Core", w.Title)
}

func TestOptions_ApplyOptions_Bad(t *testing.T) {
	boom := assert.AnError

	w, err := ApplyOptions(
		WithName("before"),
		func(*Window) error { return boom },
		WithTitle("after"),
	)

	require.ErrorIs(t, err, boom)
	assert.Nil(t, w)
}

func TestOptions_ApplyOptions_Ugly(t *testing.T) {
	w, err := ApplyOptions()

	require.NoError(t, err)
	require.NotNil(t, w)
	assert.Equal(t, &Window{}, w)
}
