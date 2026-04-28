package window

import (
	core "dappco.re/go"
)

func applyWindowOptions(t *core.T, options ...WindowOption) *Window {
	t.Helper()
	w, err := ApplyOptions(options...)
	core.RequireNoError(t, err)
	return w
}

func TestOptions_WindowOptionSetters_Good(t *core.T) {
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

	core.AssertEqual(t, "main", w.Name)
	core.AssertEqual(t, "Core GUI", w.Title)
	core.AssertEqual(t, "/dashboard", w.URL)
	core.AssertEqual(t, "<main>Ready</main>", w.HTML)
	core.AssertEqual(t, "globalThis.__CORE_READY__ = true", w.JS)
	core.AssertEqual(t, 1280, w.Width)
	core.AssertEqual(t, 800, w.Height)
	core.AssertEqual(t, 160, w.X)
	core.AssertEqual(t, 120, w.Y)
	core.AssertEqual(t, 640, w.MinWidth)
	core.AssertEqual(t, 480, w.MinHeight)
	core.AssertEqual(t, 1920, w.MaxWidth)
	core.AssertEqual(t, 1080, w.MaxHeight)
	core.AssertTrue(t, w.Frameless)
	core.AssertFalse(t, w.Hidden)
	core.AssertTrue(t, w.AlwaysOnTop)
	core.AssertEqual(t, [4]uint8{12, 34, 56, 78}, w.BackgroundColour)
	core.AssertTrue(t, w.EnableFileDrop)
}

func TestOptions_WindowOptionSetters_Bad(t *core.T) {
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

	core.AssertEqual(t, "", w.Name)
	core.AssertEqual(t, "", w.Title)
	core.AssertEqual(t, "", w.URL)
	core.AssertEqual(t, "", w.HTML)
	core.AssertEqual(t, "", w.JS)
	core.AssertEqual(t, 0, w.Width)
	core.AssertEqual(t, 0, w.Height)
	core.AssertEqual(t, 0, w.X)
	core.AssertEqual(t, 0, w.Y)
	core.AssertEqual(t, 0, w.MinWidth)
	core.AssertEqual(t, 0, w.MinHeight)
	core.AssertEqual(t, 0, w.MaxWidth)
	core.AssertEqual(t, 0, w.MaxHeight)
	core.AssertFalse(t, w.Frameless)
	core.AssertFalse(t, w.Hidden)
	core.AssertFalse(t, w.AlwaysOnTop)
	core.AssertEqual(t, [4]uint8{0, 0, 0, 0}, w.BackgroundColour)
	core.AssertFalse(t, w.EnableFileDrop)
}

func TestOptions_WindowOptionSetters_Ugly(t *core.T) {
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

	core.AssertEqual(t, "⚙︎core-window", w.Name)
	core.AssertEqual(t, "A very long title that stays intact", w.Title)
	core.AssertEqual(t, "core://settings?tab=%F0%9F%93%81", w.URL)
	core.AssertEqual(t, "<section data-id=\"αβγ\">unsafe-looking but literal</section>", w.HTML)
	core.AssertEqual(t, "globalThis.__CORE_STATE__ = { mode: 'worker', value: -1 };", w.JS)
	core.AssertEqual(t, -1920, w.Width)
	core.AssertEqual(t, -1080, w.Height)
	core.AssertEqual(t, -42, w.X)
	core.AssertEqual(t, 99999, w.Y)
	core.AssertEqual(t, -1, w.MinWidth)
	core.AssertEqual(t, -2, w.MinHeight)
	core.AssertEqual(t, 32767, w.MaxWidth)
	core.AssertEqual(t, 32767, w.MaxHeight)
	core.AssertTrue(t, w.Frameless)
	core.AssertTrue(t, w.Hidden)
	core.AssertTrue(t, w.AlwaysOnTop)
	core.AssertEqual(t, [4]uint8{255, 254, 253, 252}, w.BackgroundColour)
	core.AssertTrue(t, w.EnableFileDrop)
}

func TestOptions_ApplyOptions_Good(t *core.T) {
	w, err := ApplyOptions(
		nil,
		WithName("main"),
		WithTitle("Core"),
	)

	core.RequireNoError(t, err)
	core.AssertNotNil(t, w)
	core.AssertEqual(t, "main", w.Name)
	core.AssertEqual(t, "Core", w.Title)
}

func TestOptions_ApplyOptions_Bad(t *core.T) {
	boom := core.AnError

	w, err := ApplyOptions(
		WithName("before"),
		func(*Window) error { return boom },
		WithTitle("after"),
	)

	core.AssertErrorIs(t, err, boom)
	core.AssertNil(t, w)
}

func TestOptions_ApplyOptions_Ugly(t *core.T) {
	w, err := ApplyOptions()

	core.RequireNoError(t, err)
	core.AssertNotNil(t, w)
	core.AssertEqual(t, &Window{}, w)
}
