package application

import (
	core "dappco.re/go"
)

func TestWebviewWindowOptions_Constants_Good(t *core.T) {
	core.AssertEqual(t, WindowStateNormal, WindowState(0))
	core.AssertEqual(t, WindowStateMinimised, WindowState(1))
	core.AssertEqual(t, WindowStateMaximised, WindowState(2))
	core.AssertEqual(t, WindowStateFullscreen, WindowState(3))
	core.AssertEqual(t, WindowCentered, WindowStartPosition(0))
	core.AssertEqual(t, WindowXY, WindowStartPosition(1))
	core.AssertEqual(t, BackgroundTypeSolid, BackgroundType(0))
	core.AssertEqual(t, BackgroundTypeTransparent, BackgroundType(1))
	core.AssertEqual(t, BackgroundTypeTranslucent, BackgroundType(2))
	core.AssertEqual(t, Auto, BackdropType(0))
	core.AssertEqual(t, None, BackdropType(1))
	core.AssertEqual(t, Mica, BackdropType(2))
	core.AssertEqual(t, Acrylic, BackdropType(3))
	core.AssertEqual(t, Tabbed, BackdropType(4))
	core.AssertEqual(t, SystemDefault, Theme(0))
	core.AssertEqual(t, Dark, Theme(1))
	core.AssertEqual(t, Light, Theme(2))
	core.AssertEqual(t, MacBackdropNormal, MacBackdrop(0))
	core.AssertEqual(t, MacBackdropTransparent, MacBackdrop(1))
	core.AssertEqual(t, MacBackdropTranslucent, MacBackdrop(2))
	core.AssertEqual(t, MacBackdropLiquidGlass, MacBackdrop(3))
	core.AssertEqual(t, MacToolbarStyleAutomatic, MacToolbarStyle(0))
	core.AssertEqual(t, MacToolbarStyleUnifiedCompact, MacToolbarStyle(4))
	core.AssertEqual(t, LiquidGlassStyleAutomatic, MacLiquidGlassStyle(0))
	core.AssertEqual(t, LiquidGlassStyleVibrant, MacLiquidGlassStyle(3))
	core.AssertEqual(t, WebviewGpuPolicyAlways, WebviewGpuPolicy(0))
	core.AssertEqual(t, WebviewGpuPolicyNever, WebviewGpuPolicy(2))
	core.AssertEqual(t, LinuxMenuStyleMenuBar, LinuxMenuStyle(0))
	core.AssertEqual(t, LinuxMenuStylePrimaryMenu, LinuxMenuStyle(1))
	core.AssertEqual(t, DefaultAppearance, MacAppearanceType(""))
	core.AssertEqual(t, NSAppearanceNameDarkAqua, MacAppearanceType("NSAppearanceNameDarkAqua"))
}

func TestWebviewWindowOptions_Constants_Bad(t *core.T) {
	core.AssertEqual(t, MacTitleBarDefault, MacTitleBar{})
	core.AssertEqual(t, MacTitleBarHidden, MacTitleBar{
		AppearsTransparent: true,
		HideTitle:          true,
		FullSizeContent:    true,
	})
	core.AssertEqual(t, MacTitleBarHiddenInset, MacTitleBar{
		AppearsTransparent:   true,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           true,
		HideToolbarSeparator: true,
	})
	core.AssertEqual(t, MacTitleBarHiddenInsetUnified, MacTitleBar{
		AppearsTransparent:   true,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           true,
		HideToolbarSeparator: true,
		ToolbarStyle:         MacToolbarStyleUnified,
	})
	core.AssertEqual(t, NSVisualEffectMaterialAuto, NSVisualEffectMaterial(-1))
	core.AssertEqual(t, MacWindowLevelNormal, MacWindowLevel("normal"))
	core.AssertEqual(t, MacWindowCollectionBehaviorCanJoinAllSpaces, MacWindowCollectionBehavior(1))
	core.AssertEqual(t, MacWindowCollectionBehaviorFullScreenAuxiliary, MacWindowCollectionBehavior(1<<8))
}

func TestWebviewWindowOptions_Constants_Ugly(t *core.T) {
	options := WebviewWindowOptions{
		Name:             "main",
		Title:            "Main",
		URL:              "https://example.invalid",
		HTML:             "<h1>Hello</h1>",
		JS:               "window.__ready = true",
		Width:            800,
		Height:           600,
		X:                10,
		Y:                20,
		MinWidth:         320,
		MinHeight:        240,
		MaxWidth:         1920,
		MaxHeight:        1080,
		Frameless:        true,
		Hidden:           true,
		AlwaysOnTop:      true,
		DisableResize:    true,
		EnableFileDrop:   true,
		BackgroundColour: NewRGBA(1, 2, 3, 4),
	}

	core.AssertEqual(t, "main", options.Name)
	core.AssertEqual(t, "Main", options.Title)
	core.AssertEqual(t, 800, options.Width)
	core.AssertEqual(t, 600, options.Height)
	core.AssertTrue(t, options.Frameless)
	core.AssertTrue(t, options.Hidden)
	core.AssertTrue(t, options.AlwaysOnTop)
	core.AssertTrue(t, options.DisableResize)
	core.AssertTrue(t, options.EnableFileDrop)
	core.AssertEqual(t, RGBA{Red: 1, Green: 2, Blue: 3, Alpha: 4}, options.BackgroundColour)
}
