package application

import "github.com/wailsapp/wails/v3/pkg/events"

// WindowState represents the starting visual state of a window.
type WindowState int

const (
	WindowStateNormal     WindowState = iota
	WindowStateMinimised  WindowState = iota
	WindowStateMaximised  WindowState = iota
	WindowStateFullscreen WindowState = iota
)

// WindowStartPosition determines how the initial X/Y coordinates are interpreted.
type WindowStartPosition int

const (
	// WindowCentered places the window in the centre of the screen on first show.
	WindowCentered WindowStartPosition = 0
	// WindowXY places the window at the explicit X/Y coordinates.
	WindowXY WindowStartPosition = 1
)

// BackgroundType controls window transparency.
type BackgroundType int

const (
	BackgroundTypeSolid       BackgroundType = iota
	BackgroundTypeTransparent BackgroundType = iota
	BackgroundTypeTranslucent BackgroundType = iota
)

// WebviewWindowOptions holds all configuration for a webview window.
//
//	opts := WebviewWindowOptions{
//	    Name:   "main",
//	    Title:  "My App",
//	    Width:  1280,
//	    Height: 800,
//	    Mac:    MacWindow{TitleBar: MacTitleBarHiddenInset},
//	}
type WebviewWindowOptions struct {
	// Name is a unique identifier for the window.
	Name string
	// Title is shown in the title bar.
	Title string
	// Width is the initial width in logical pixels.
	Width int
	// Height is the initial height in logical pixels.
	Height int
	// AlwaysOnTop makes the window float above others.
	AlwaysOnTop bool
	// URL is the URL to load on creation.
	URL string
	// HTML is inline HTML to load (alternative to URL).
	HTML string
	// JS is inline JavaScript to inject.
	JS string
	// CSS is inline CSS to inject.
	CSS string
	// DisableResize prevents the user from resizing the window.
	DisableResize bool
	// Frameless removes the OS window frame.
	Frameless bool
	// MinWidth is the minimum allowed width.
	MinWidth int
	// MinHeight is the minimum allowed height.
	MinHeight int
	// MaxWidth is the maximum allowed width (0 = unlimited).
	MaxWidth int
	// MaxHeight is the maximum allowed height (0 = unlimited).
	MaxHeight int
	// StartState sets the visual state when first shown.
	StartState WindowState
	// BackgroundType controls transparency.
	BackgroundType BackgroundType
	// BackgroundColour is the solid background fill colour.
	BackgroundColour RGBA
	// InitialPosition controls how X/Y are interpreted.
	InitialPosition WindowStartPosition
	// X is the initial horizontal position.
	X int
	// Y is the initial vertical position.
	Y int
	// Hidden creates the window without showing it.
	Hidden bool
	// Zoom sets the initial zoom magnification (0 = default = 1.0).
	Zoom float64
	// ZoomControlEnabled allows the user to change zoom.
	ZoomControlEnabled bool
	// EnableFileDrop enables drag-and-drop of files onto the window.
	EnableFileDrop bool
	// OpenInspectorOnStartup opens the web inspector when first shown.
	OpenInspectorOnStartup bool
	// DevToolsEnabled exposes the developer tools (default true in debug builds).
	DevToolsEnabled bool
	// DefaultContextMenuDisabled disables the built-in right-click menu.
	DefaultContextMenuDisabled bool
	// KeyBindings is a map of accelerator strings to window callbacks.
	KeyBindings map[string]func(window Window)
	// IgnoreMouseEvents passes all mouse events through (Windows + macOS only).
	IgnoreMouseEvents bool
	// ContentProtectionEnabled prevents screen capture (Windows + macOS only).
	ContentProtectionEnabled bool
	// HideOnFocusLost hides the window when it loses focus.
	HideOnFocusLost bool
	// HideOnEscape hides the window when Escape is pressed.
	HideOnEscape bool
	// UseApplicationMenu uses the application-level menu for this window.
	UseApplicationMenu bool
	// MinimiseButtonState controls the minimise button state.
	MinimiseButtonState ButtonState
	// MaximiseButtonState controls the maximise/zoom button state.
	MaximiseButtonState ButtonState
	// CloseButtonState controls the close button state.
	CloseButtonState ButtonState
	// Mac contains macOS-specific window options.
	Mac MacWindow
	// Windows contains Windows-specific window options.
	Windows WindowsWindow
	// Linux contains Linux-specific window options.
	Linux LinuxWindow
}

// -------------------------
// macOS-specific types
// -------------------------

// MacBackdrop is the backdrop material for a macOS window.
type MacBackdrop int

const (
	// MacBackdropNormal uses an opaque background.
	MacBackdropNormal MacBackdrop = iota
	// MacBackdropTransparent shows the desktop behind the window.
	MacBackdropTransparent MacBackdrop = iota
	// MacBackdropTranslucent applies a frosted-glass vibrancy effect.
	MacBackdropTranslucent MacBackdrop = iota
	// MacBackdropLiquidGlass uses the Apple Liquid Glass effect (macOS 15+).
	MacBackdropLiquidGlass MacBackdrop = iota
)

// MacToolbarStyle controls toolbar layout relative to the title bar.
type MacToolbarStyle int

const (
	MacToolbarStyleAutomatic      MacToolbarStyle = iota
	MacToolbarStyleExpanded       MacToolbarStyle = iota
	MacToolbarStylePreference     MacToolbarStyle = iota
	MacToolbarStyleUnified        MacToolbarStyle = iota
	MacToolbarStyleUnifiedCompact MacToolbarStyle = iota
)

// MacLiquidGlassStyle defines the tint of the Liquid Glass effect.
type MacLiquidGlassStyle int

const (
	LiquidGlassStyleAutomatic MacLiquidGlassStyle = iota
	LiquidGlassStyleLight     MacLiquidGlassStyle = iota
	LiquidGlassStyleDark      MacLiquidGlassStyle = iota
	LiquidGlassStyleVibrant   MacLiquidGlassStyle = iota
)

// NSVisualEffectMaterial maps to the NSVisualEffectMaterial macOS enum.
type NSVisualEffectMaterial int

const (
	NSVisualEffectMaterialAppearanceBased       NSVisualEffectMaterial = 0
	NSVisualEffectMaterialLight                 NSVisualEffectMaterial = 1
	NSVisualEffectMaterialDark                  NSVisualEffectMaterial = 2
	NSVisualEffectMaterialTitlebar              NSVisualEffectMaterial = 3
	NSVisualEffectMaterialSelection             NSVisualEffectMaterial = 4
	NSVisualEffectMaterialMenu                  NSVisualEffectMaterial = 5
	NSVisualEffectMaterialPopover               NSVisualEffectMaterial = 6
	NSVisualEffectMaterialSidebar               NSVisualEffectMaterial = 7
	NSVisualEffectMaterialHeaderView            NSVisualEffectMaterial = 10
	NSVisualEffectMaterialSheet                 NSVisualEffectMaterial = 11
	NSVisualEffectMaterialWindowBackground      NSVisualEffectMaterial = 12
	NSVisualEffectMaterialHUDWindow             NSVisualEffectMaterial = 13
	NSVisualEffectMaterialFullScreenUI          NSVisualEffectMaterial = 15
	NSVisualEffectMaterialToolTip               NSVisualEffectMaterial = 17
	NSVisualEffectMaterialContentBackground     NSVisualEffectMaterial = 18
	NSVisualEffectMaterialUnderWindowBackground NSVisualEffectMaterial = 21
	NSVisualEffectMaterialUnderPageBackground   NSVisualEffectMaterial = 22
	NSVisualEffectMaterialAuto                  NSVisualEffectMaterial = -1
)

// MacLiquidGlass configures the Liquid Glass compositor effect.
type MacLiquidGlass struct {
	Style        MacLiquidGlassStyle
	Material     NSVisualEffectMaterial
	CornerRadius float64
	TintColor    *RGBA
	GroupID      string
	GroupSpacing float64
}

// MacAppearanceType is the NSAppearance name string for a macOS window.
type MacAppearanceType string

const (
	DefaultAppearance                                    MacAppearanceType = ""
	NSAppearanceNameAqua                                 MacAppearanceType = "NSAppearanceNameAqua"
	NSAppearanceNameDarkAqua                             MacAppearanceType = "NSAppearanceNameDarkAqua"
	NSAppearanceNameVibrantLight                         MacAppearanceType = "NSAppearanceNameVibrantLight"
	NSAppearanceNameAccessibilityHighContrastAqua        MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastAqua"
	NSAppearanceNameAccessibilityHighContrastDarkAqua    MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastDarkAqua"
	NSAppearanceNameAccessibilityHighContrastVibrantLight MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantLight"
	NSAppearanceNameAccessibilityHighContrastVibrantDark  MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantDark"
)

// MacWindowLevel controls Z-ordering relative to other windows.
type MacWindowLevel string

const (
	MacWindowLevelNormal      MacWindowLevel = "normal"
	MacWindowLevelFloating    MacWindowLevel = "floating"
	MacWindowLevelTornOffMenu MacWindowLevel = "tornOffMenu"
	MacWindowLevelModalPanel  MacWindowLevel = "modalPanel"
	MacWindowLevelMainMenu    MacWindowLevel = "mainMenu"
	MacWindowLevelStatus      MacWindowLevel = "status"
	MacWindowLevelPopUpMenu   MacWindowLevel = "popUpMenu"
	MacWindowLevelScreenSaver MacWindowLevel = "screenSaver"
)

// MacWindowCollectionBehavior is a bitmask controlling Spaces and fullscreen behaviour.
type MacWindowCollectionBehavior int

const (
	MacWindowCollectionBehaviorDefault                   MacWindowCollectionBehavior = 0
	MacWindowCollectionBehaviorCanJoinAllSpaces          MacWindowCollectionBehavior = 1 << 0
	MacWindowCollectionBehaviorMoveToActiveSpace         MacWindowCollectionBehavior = 1 << 1
	MacWindowCollectionBehaviorManaged                   MacWindowCollectionBehavior = 1 << 2
	MacWindowCollectionBehaviorTransient                 MacWindowCollectionBehavior = 1 << 3
	MacWindowCollectionBehaviorStationary                MacWindowCollectionBehavior = 1 << 4
	MacWindowCollectionBehaviorParticipatesInCycle       MacWindowCollectionBehavior = 1 << 5
	MacWindowCollectionBehaviorIgnoresCycle              MacWindowCollectionBehavior = 1 << 6
	MacWindowCollectionBehaviorFullScreenPrimary         MacWindowCollectionBehavior = 1 << 7
	MacWindowCollectionBehaviorFullScreenAuxiliary       MacWindowCollectionBehavior = 1 << 8
	MacWindowCollectionBehaviorFullScreenNone            MacWindowCollectionBehavior = 1 << 9
	MacWindowCollectionBehaviorFullScreenAllowsTiling    MacWindowCollectionBehavior = 1 << 11
	MacWindowCollectionBehaviorFullScreenDisallowsTiling MacWindowCollectionBehavior = 1 << 12
)

// MacWebviewPreferences holds WKWebView preference flags for macOS.
// Use integer tristate: 0 = unset, 1 = true, 2 = false.
type MacWebviewPreferences struct {
	TabFocusesLinks                     int
	TextInteractionEnabled              int
	FullscreenEnabled                   int
	AllowsBackForwardNavigationGestures int
}

// MacTitleBar configures the macOS title bar appearance.
type MacTitleBar struct {
	// AppearsTransparent removes the title bar background.
	AppearsTransparent bool
	// Hide removes the title bar entirely.
	Hide bool
	// HideTitle hides only the text title.
	HideTitle bool
	// FullSizeContent extends window content into the title bar area.
	FullSizeContent bool
	// UseToolbar replaces the title bar with an NSToolbar.
	UseToolbar bool
	// HideToolbarSeparator removes the line between toolbar and content.
	HideToolbarSeparator bool
	// ShowToolbarWhenFullscreen keeps the toolbar visible in fullscreen.
	ShowToolbarWhenFullscreen bool
	// ToolbarStyle selects the toolbar layout style.
	ToolbarStyle MacToolbarStyle
}

// MacWindow contains macOS-specific window options.
type MacWindow struct {
	Backdrop                        MacBackdrop
	DisableShadow                   bool
	TitleBar                        MacTitleBar
	Appearance                      MacAppearanceType
	InvisibleTitleBarHeight         int
	EventMapping                    map[events.WindowEventType]events.WindowEventType
	EnableFraudulentWebsiteWarnings bool
	WebviewPreferences              MacWebviewPreferences
	WindowLevel                     MacWindowLevel
	CollectionBehavior              MacWindowCollectionBehavior
	LiquidGlass                     MacLiquidGlass
}

// Pre-built MacTitleBar configurations — use directly in MacWindow.TitleBar.

// MacTitleBarDefault produces the standard macOS title bar.
//
//	Mac: MacWindow{TitleBar: MacTitleBarDefault}
var MacTitleBarDefault = MacTitleBar{}

// MacTitleBarHidden hides the title text while keeping the traffic-light buttons.
//
//	Mac: MacWindow{TitleBar: MacTitleBarHidden}
var MacTitleBarHidden = MacTitleBar{
	AppearsTransparent: true,
	HideTitle:          true,
	FullSizeContent:    true,
}

// MacTitleBarHiddenInset keeps traffic lights slightly inset from the window edge.
//
//	Mac: MacWindow{TitleBar: MacTitleBarHiddenInset}
var MacTitleBarHiddenInset = MacTitleBar{
	AppearsTransparent:   true,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           true,
	HideToolbarSeparator: true,
}

// MacTitleBarHiddenInsetUnified uses the unified toolbar style for a more compact look.
//
//	Mac: MacWindow{TitleBar: MacTitleBarHiddenInsetUnified}
var MacTitleBarHiddenInsetUnified = MacTitleBar{
	AppearsTransparent:   true,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           true,
	HideToolbarSeparator: true,
	ToolbarStyle:         MacToolbarStyleUnified,
}

// -------------------------
// Windows-specific types
// -------------------------

// BackdropType selects the Windows 11 compositor effect.
type BackdropType int32

const (
	Auto    BackdropType = 0
	None    BackdropType = 1
	Mica    BackdropType = 2
	Acrylic BackdropType = 3
	Tabbed  BackdropType = 4
)

// Theme selects dark or light mode on Windows.
type Theme int

const (
	SystemDefault Theme = 0
	Dark          Theme = 1
	Light         Theme = 2
)

// WindowTheme holds custom title-bar colours for a Windows window.
type WindowTheme struct {
	BorderColour    *uint32
	TitleBarColour  *uint32
	TitleTextColour *uint32
}

// TextTheme holds foreground/background colour pair for menu text.
type TextTheme struct {
	Text       *uint32
	Background *uint32
}

// MenuBarTheme holds per-state text themes for the Windows menu bar.
type MenuBarTheme struct {
	Default  *TextTheme
	Hover    *TextTheme
	Selected *TextTheme
}

// ThemeSettings holds custom colours for dark/light mode on Windows.
// Colours use 0x00BBGGRR encoding.
type ThemeSettings struct {
	DarkModeActive    *WindowTheme
	DarkModeInactive  *WindowTheme
	LightModeActive   *WindowTheme
	LightModeInactive *WindowTheme
	DarkModeMenuBar   *MenuBarTheme
	LightModeMenuBar  *MenuBarTheme
}

// CoreWebView2PermissionKind identifies a WebView2 permission category.
type CoreWebView2PermissionKind uint32

const (
	CoreWebView2PermissionKindUnknownPermission CoreWebView2PermissionKind = iota
	CoreWebView2PermissionKindMicrophone
	CoreWebView2PermissionKindCamera
	CoreWebView2PermissionKindGeolocation
	CoreWebView2PermissionKindNotifications
	CoreWebView2PermissionKindOtherSensors
	CoreWebView2PermissionKindClipboardRead
)

// CoreWebView2PermissionState sets whether a permission is granted.
type CoreWebView2PermissionState uint32

const (
	CoreWebView2PermissionStateDefault CoreWebView2PermissionState = iota
	CoreWebView2PermissionStateAllow
	CoreWebView2PermissionStateDeny
)

// WindowsWindow contains Windows-specific window options.
type WindowsWindow struct {
	BackdropType                      BackdropType
	DisableIcon                       bool
	Theme                             Theme
	CustomTheme                       ThemeSettings
	DisableFramelessWindowDecorations bool
	WindowMask                        []byte
	WindowMaskDraggable               bool
	ResizeDebounceMS                  uint16
	WindowDidMoveDebounceMS           uint16
	EventMapping                      map[events.WindowEventType]events.WindowEventType
	HiddenOnTaskbar                   bool
	EnableSwipeGestures               bool
	Menu                              *Menu
	Permissions                       map[CoreWebView2PermissionKind]CoreWebView2PermissionState
	ExStyle                           int
	GeneralAutofillEnabled            bool
	PasswordAutosaveEnabled           bool
}

// -------------------------
// Linux-specific types
// -------------------------

// WebviewGpuPolicy controls GPU acceleration for the Linux webview.
type WebviewGpuPolicy int

const (
	WebviewGpuPolicyAlways   WebviewGpuPolicy = iota
	WebviewGpuPolicyOnDemand WebviewGpuPolicy = iota
	WebviewGpuPolicyNever    WebviewGpuPolicy = iota
)

// LinuxMenuStyle selects how the application menu is rendered on Linux.
type LinuxMenuStyle int

const (
	LinuxMenuStyleMenuBar     LinuxMenuStyle = iota
	LinuxMenuStylePrimaryMenu LinuxMenuStyle = iota
)

// LinuxWindow contains Linux-specific window options.
type LinuxWindow struct {
	Icon                    []byte
	WindowIsTranslucent     bool
	WebviewGpuPolicy        WebviewGpuPolicy
	WindowDidMoveDebounceMS uint16
	Menu                    *Menu
	MenuStyle               LinuxMenuStyle
}

// NewRGB constructs an RGBA value with full opacity from RGB components.
//
//	colour := NewRGB(255, 128, 0)
func NewRGB(red, green, blue uint8) RGBA {
	return RGBA{Red: red, Green: green, Blue: blue, Alpha: 255}
}

// NewRGBPtr encodes RGB as a packed uint32 pointer (0x00BBGGRR) for ThemeSettings.
//
//	theme.BorderColour = NewRGBPtr(255, 0, 0)
func NewRGBPtr(red, green, blue uint8) *uint32 {
	result := uint32(red) | uint32(green)<<8 | uint32(blue)<<16
	return &result
}
