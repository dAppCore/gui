package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebviewWindowOptions_Constants_Good(t *testing.T) {
	assert.Equal(t, WindowStateNormal, WindowState(0))
	assert.Equal(t, WindowStateMinimised, WindowState(1))
	assert.Equal(t, WindowStateMaximised, WindowState(2))
	assert.Equal(t, WindowStateFullscreen, WindowState(3))
	assert.Equal(t, WindowCentered, WindowStartPosition(0))
	assert.Equal(t, WindowXY, WindowStartPosition(1))
	assert.Equal(t, BackgroundTypeSolid, BackgroundType(0))
	assert.Equal(t, BackgroundTypeTransparent, BackgroundType(1))
	assert.Equal(t, BackgroundTypeTranslucent, BackgroundType(2))
	assert.Equal(t, Auto, BackdropType(0))
	assert.Equal(t, None, BackdropType(1))
	assert.Equal(t, Mica, BackdropType(2))
	assert.Equal(t, Acrylic, BackdropType(3))
	assert.Equal(t, Tabbed, BackdropType(4))
	assert.Equal(t, SystemDefault, Theme(0))
	assert.Equal(t, Dark, Theme(1))
	assert.Equal(t, Light, Theme(2))
	assert.Equal(t, MacBackdropNormal, MacBackdrop(0))
	assert.Equal(t, MacBackdropTransparent, MacBackdrop(1))
	assert.Equal(t, MacBackdropTranslucent, MacBackdrop(2))
	assert.Equal(t, MacBackdropLiquidGlass, MacBackdrop(3))
	assert.Equal(t, MacToolbarStyleAutomatic, MacToolbarStyle(0))
	assert.Equal(t, MacToolbarStyleUnifiedCompact, MacToolbarStyle(4))
	assert.Equal(t, LiquidGlassStyleAutomatic, MacLiquidGlassStyle(0))
	assert.Equal(t, LiquidGlassStyleVibrant, MacLiquidGlassStyle(3))
	assert.Equal(t, WebviewGpuPolicyAlways, WebviewGpuPolicy(0))
	assert.Equal(t, WebviewGpuPolicyNever, WebviewGpuPolicy(2))
	assert.Equal(t, LinuxMenuStyleMenuBar, LinuxMenuStyle(0))
	assert.Equal(t, LinuxMenuStylePrimaryMenu, LinuxMenuStyle(1))
	assert.Equal(t, DefaultAppearance, MacAppearanceType(""))
	assert.Equal(t, NSAppearanceNameDarkAqua, MacAppearanceType("NSAppearanceNameDarkAqua"))
}

func TestWebviewWindowOptions_Constants_Bad(t *testing.T) {
	assert.Equal(t, MacTitleBarDefault, MacTitleBar{})
	assert.Equal(t, MacTitleBarHidden, MacTitleBar{
		AppearsTransparent: true,
		HideTitle:          true,
		FullSizeContent:    true,
	})
	assert.Equal(t, MacTitleBarHiddenInset, MacTitleBar{
		AppearsTransparent:   true,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           true,
		HideToolbarSeparator: true,
	})
	assert.Equal(t, MacTitleBarHiddenInsetUnified, MacTitleBar{
		AppearsTransparent:   true,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           true,
		HideToolbarSeparator: true,
		ToolbarStyle:         MacToolbarStyleUnified,
	})
	assert.Equal(t, NSVisualEffectMaterialAuto, NSVisualEffectMaterial(-1))
	assert.Equal(t, MacWindowLevelNormal, MacWindowLevel("normal"))
	assert.Equal(t, MacWindowCollectionBehaviorCanJoinAllSpaces, MacWindowCollectionBehavior(1))
	assert.Equal(t, MacWindowCollectionBehaviorFullScreenAuxiliary, MacWindowCollectionBehavior(1<<8))
}

func TestWebviewWindowOptions_Constants_Ugly(t *testing.T) {
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

	require.Equal(t, "main", options.Name)
	assert.Equal(t, "Main", options.Title)
	assert.Equal(t, 800, options.Width)
	assert.Equal(t, 600, options.Height)
	assert.True(t, options.Frameless)
	assert.True(t, options.Hidden)
	assert.True(t, options.AlwaysOnTop)
	assert.True(t, options.DisableResize)
	assert.True(t, options.EnableFileDrop)
	assert.Equal(t, RGBA{Red: 1, Green: 2, Blue: 3, Alpha: 4}, options.BackgroundColour)
}
