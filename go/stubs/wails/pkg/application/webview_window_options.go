package application

import "github.com/wailsapp/wails/v3/pkg/events"

// WindowState represents the visible state of a window.
//
//	opts := application.WebviewWindowOptions{StartState: application.WindowStateMaximised}
type WindowState int

const (
	// WindowStateNormal is the default windowed state.
	WindowStateNormal WindowState = iota
	// WindowStateMinimised is the minimised (iconified) state.
	WindowStateMinimised
	// WindowStateMaximised fills the available work area.
	WindowStateMaximised
	// WindowStateFullscreen occupies the full screen.
	WindowStateFullscreen
)

// WindowStartPosition controls where a window first appears.
//
//	opts := application.WebviewWindowOptions{InitialPosition: application.WindowCentered}
type WindowStartPosition int

const (
	// WindowCentered places the window at the centre of the screen.
	WindowCentered WindowStartPosition = 0
	// WindowXY places the window at the coordinates given by X and Y.
	WindowXY WindowStartPosition = 1
)

// BackgroundType determines how the window background is rendered.
//
//	opts := application.WebviewWindowOptions{BackgroundType: application.BackgroundTypeTranslucent}
type BackgroundType int

const (
	// BackgroundTypeSolid renders a solid opaque background.
	BackgroundTypeSolid BackgroundType = iota
	// BackgroundTypeTransparent renders a fully transparent background.
	BackgroundTypeTransparent
	// BackgroundTypeTranslucent renders a frosted/blur translucent background.
	BackgroundTypeTranslucent
)

// NewRGB constructs an opaque RGBA value from red, green, blue components.
//
//	colour := application.NewRGB(0xff, 0x00, 0x00) // red
func NewRGB(red, green, blue uint8) RGBA {
	return RGBA{Red: red, Green: green, Blue: blue, Alpha: 255}
}

// NewRGBPtr encodes red, green, blue as a packed *uint32 in 0x00BBGGRR order.
//
//	ptr := application.NewRGBPtr(0xff, 0x80, 0x00)
func NewRGBPtr(red, green, blue uint8) *uint32 {
	value := uint32(red) | uint32(green)<<8 | uint32(blue)<<16
	return &value
}

/******* Windows Options *******/

// BackdropType selects the translucent backdrop style on Windows 11.
//
//	opts.Windows = application.WindowsWindow{BackdropType: application.Mica}
type BackdropType int32

const (
	// Auto lets the system choose the best backdrop.
	Auto BackdropType = 0
	// None disables the translucent backdrop.
	None BackdropType = 1
	// Mica applies the Mica material (Windows 11 22H2+).
	Mica BackdropType = 2
	// Acrylic applies the Acrylic blur-behind material.
	Acrylic BackdropType = 3
	// Tabbed applies the Tabbed/MICA-Alt material.
	Tabbed BackdropType = 4
)

// CoreWebView2PermissionKind enumerates the types of WebView2 permissions.
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

// CoreWebView2PermissionState enumerates the allowed states for a WebView2 permission.
type CoreWebView2PermissionState uint32

const (
	CoreWebView2PermissionStateDefault CoreWebView2PermissionState = iota
	CoreWebView2PermissionStateAllow
	CoreWebView2PermissionStateDeny
)

// Theme selects between the system default, dark, and light UI themes on Windows.
//
//	opts.Windows = application.WindowsWindow{Theme: application.Dark}
type Theme int

const (
	// SystemDefault follows the OS theme and reacts to changes.
	SystemDefault Theme = 0
	// Dark forces the dark theme.
	Dark Theme = 1
	// Light forces the light theme.
	Light Theme = 2
)

// WindowTheme defines colour overrides for a single window activity state.
//
//	wt := &application.WindowTheme{TitleBarColour: application.NewRGBPtr(0x1e, 0x1e, 0x1e)}
type WindowTheme struct {
	// BorderColour is the colour of the window border (0x00BBGGRR).
	BorderColour *uint32
	// TitleBarColour is the colour of the title bar (0x00BBGGRR).
	TitleBarColour *uint32
	// TitleTextColour is the colour of the title text (0x00BBGGRR).
	TitleTextColour *uint32
}

// TextTheme defines foreground and background colours for a text element.
type TextTheme struct {
	// Text is the foreground colour.
	Text *uint32
	// Background is the background colour.
	Background *uint32
}

// MenuBarTheme defines colours for a menu bar in default, hovered, and selected states.
type MenuBarTheme struct {
	// Default is the theme used when the item is neither hovered nor selected.
	Default *TextTheme
	// Hover is the theme used when the pointer is over the item.
	Hover *TextTheme
	// Selected is the theme used when the item is selected.
	Selected *TextTheme
}

// ThemeSettings defines custom colours used in dark or light mode.
// Colour values use packed 0x00BBGGRR encoding — use NewRGBPtr to construct them.
//
//	ts := application.ThemeSettings{
//	    DarkModeActive: &application.WindowTheme{TitleBarColour: application.NewRGBPtr(0x1e, 0x1e, 0x2e)},
//	}
type ThemeSettings struct {
	// DarkModeActive applies when the window is active in dark mode.
	DarkModeActive *WindowTheme
	// DarkModeInactive applies when the window is inactive in dark mode.
	DarkModeInactive *WindowTheme
	// LightModeActive applies when the window is active in light mode.
	LightModeActive *WindowTheme
	// LightModeInactive applies when the window is inactive in light mode.
	LightModeInactive *WindowTheme
	// DarkModeMenuBar applies to the menu bar in dark mode.
	DarkModeMenuBar *MenuBarTheme
	// LightModeMenuBar applies to the menu bar in light mode.
	LightModeMenuBar *MenuBarTheme
}

// WindowsWindow contains Windows-specific window configuration.
//
//	opts.Windows = application.WindowsWindow{BackdropType: application.Mica, Theme: application.Dark}
type WindowsWindow struct {
	// BackdropType selects the translucent material. Requires Windows 11 22621+.
	// Only used when BackgroundType is BackgroundTypeTranslucent.
	// Default: Auto
	BackdropType BackdropType

	// DisableIcon removes the application icon from the title bar.
	// Default: false
	DisableIcon bool

	// Theme selects between dark, light, or system-default title bar styling.
	// Default: SystemDefault
	Theme Theme

	// CustomTheme overrides colours for dark/light active/inactive states.
	// Default: zero value (no override)
	CustomTheme ThemeSettings

	// DisableFramelessWindowDecorations suppresses Aero shadow and rounded corners
	// when the window is frameless. Rounded corners require Windows 11.
	// Default: false
	DisableFramelessWindowDecorations bool

	// WindowMask sets the window shape via a PNG with an alpha channel.
	// Default: nil
	WindowMask []byte

	// WindowMaskDraggable allows the window to be dragged via the mask area.
	// Default: false
	WindowMaskDraggable bool

	// ResizeDebounceMS debounces WebView2 redraws during resize.
	// Default: 0
	ResizeDebounceMS uint16

	// WindowDidMoveDebounceMS debounces the WindowDidMove event.
	// Default: 0
	WindowDidMoveDebounceMS uint16

	// EventMapping translates platform window events to common event types.
	// Default: nil
	EventMapping map[events.WindowEventType]events.WindowEventType

	// HiddenOnTaskbar excludes the window from the taskbar.
	// Default: false
	HiddenOnTaskbar bool

	// EnableSwipeGestures enables horizontal swipe gestures.
	// Default: false
	EnableSwipeGestures bool

	// Menu is the window-level menu.
	Menu *Menu

	// Permissions configures WebView2 permission grants.
	// Default: nil (system defaults apply)
	Permissions map[CoreWebView2PermissionKind]CoreWebView2PermissionState

	// ExStyle is the extended window style flags (WS_EX_*).
	ExStyle int

	// GeneralAutofillEnabled enables general autofill in WebView2.
	GeneralAutofillEnabled bool

	// PasswordAutosaveEnabled enables password autosave in WebView2.
	PasswordAutosaveEnabled bool
}

/****** Mac Options *******/

// MacBackdrop controls the translucency of the macOS window background.
//
//	opts.Mac = application.MacWindow{Backdrop: application.MacBackdropTranslucent}
type MacBackdrop int

const (
	// MacBackdropNormal renders a standard opaque background.
	MacBackdropNormal MacBackdrop = iota
	// MacBackdropTransparent renders a fully transparent background.
	MacBackdropTransparent
	// MacBackdropTranslucent renders a frosted vibrancy background.
	MacBackdropTranslucent
	// MacBackdropLiquidGlass applies Apple's Liquid Glass effect (macOS 15+,
	// falls back to translucent on earlier releases).
	MacBackdropLiquidGlass
)

// MacToolbarStyle controls the toolbar layout relative to the title bar.
//
//	opts.Mac.TitleBar.ToolbarStyle = application.MacToolbarStyleUnified
type MacToolbarStyle int

const (
	// MacToolbarStyleAutomatic lets the system decide based on configuration.
	MacToolbarStyleAutomatic MacToolbarStyle = iota
	// MacToolbarStyleExpanded shows the toolbar below the title bar.
	MacToolbarStyleExpanded
	// MacToolbarStylePreference shows the toolbar below the title bar with
	// equal-width items where possible.
	MacToolbarStylePreference
	// MacToolbarStyleUnified merges the title bar and toolbar into one row.
	MacToolbarStyleUnified
	// MacToolbarStyleUnifiedCompact is like Unified but with reduced margins.
	MacToolbarStyleUnifiedCompact
)

// MacLiquidGlassStyle defines the appearance of the Liquid Glass effect.
type MacLiquidGlassStyle int

const (
	// LiquidGlassStyleAutomatic lets the system choose the best style.
	LiquidGlassStyleAutomatic MacLiquidGlassStyle = iota
	// LiquidGlassStyleLight uses a light glass appearance.
	LiquidGlassStyleLight
	// LiquidGlassStyleDark uses a dark glass appearance.
	LiquidGlassStyleDark
	// LiquidGlassStyleVibrant uses an enhanced vibrant glass appearance.
	LiquidGlassStyleVibrant
)

// NSVisualEffectMaterial mirrors NSVisualEffectMaterial from the macOS SDK.
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
	// NSVisualEffectMaterialAuto selects the material automatically based on Style.
	NSVisualEffectMaterialAuto NSVisualEffectMaterial = -1
)

// MacLiquidGlass configures the Liquid Glass visual effect (macOS 15+).
//
//	opts.Mac.LiquidGlass = application.MacLiquidGlass{Style: application.LiquidGlassStyleDark}
type MacLiquidGlass struct {
	// Style of the glass effect.
	Style MacLiquidGlassStyle

	// Material for the NSVisualEffectView fallback.
	// Use NSVisualEffectMaterialAuto for automatic selection based on Style.
	Material NSVisualEffectMaterial

	// CornerRadius specifies the corner radius in points (0 for square corners).
	CornerRadius float64

	// TintColor adds an optional colour tint to the glass (nil for no tint).
	TintColor *RGBA

	// GroupID merges multiple glass windows into a single visual group.
	GroupID string

	// GroupSpacing is the spacing between grouped glass elements in points.
	GroupSpacing float64
}

// MacAppearanceType selects a Cocoa NSAppearance for the window.
//
//	opts.Mac = application.MacWindow{Appearance: application.NSAppearanceNameDarkAqua}
type MacAppearanceType string

const (
	// DefaultAppearance follows the system setting.
	DefaultAppearance MacAppearanceType = ""
	// NSAppearanceNameAqua is the standard light system appearance.
	NSAppearanceNameAqua MacAppearanceType = "NSAppearanceNameAqua"
	// NSAppearanceNameDarkAqua is the standard dark system appearance.
	NSAppearanceNameDarkAqua MacAppearanceType = "NSAppearanceNameDarkAqua"
	// NSAppearanceNameVibrantLight is the light vibrant appearance.
	NSAppearanceNameVibrantLight MacAppearanceType = "NSAppearanceNameVibrantLight"
	// NSAppearanceNameAccessibilityHighContrastAqua is high-contrast light.
	NSAppearanceNameAccessibilityHighContrastAqua MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastAqua"
	// NSAppearanceNameAccessibilityHighContrastDarkAqua is high-contrast dark.
	NSAppearanceNameAccessibilityHighContrastDarkAqua MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastDarkAqua"
	// NSAppearanceNameAccessibilityHighContrastVibrantLight is high-contrast light vibrant.
	NSAppearanceNameAccessibilityHighContrastVibrantLight MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantLight"
	// NSAppearanceNameAccessibilityHighContrastVibrantDark is high-contrast dark vibrant.
	NSAppearanceNameAccessibilityHighContrastVibrantDark MacAppearanceType = "NSAppearanceNameAccessibilityHighContrastVibrantDark"
)

// MacWindowLevel controls the z-order stacking group of the window.
//
//	opts.Mac = application.MacWindow{WindowLevel: application.MacWindowLevelFloating}
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

// MacWindowCollectionBehavior controls how the window participates in macOS
// Spaces and fullscreen. Values correspond to NSWindowCollectionBehavior bits
// and may be combined with bitwise OR.
//
//	opts.Mac.CollectionBehavior = application.MacWindowCollectionBehaviorCanJoinAllSpaces |
//	    application.MacWindowCollectionBehaviorFullScreenAuxiliary
type MacWindowCollectionBehavior int

const (
	// MacWindowCollectionBehaviorDefault uses FullScreenPrimary for backwards compatibility.
	MacWindowCollectionBehaviorDefault MacWindowCollectionBehavior = 0
	// MacWindowCollectionBehaviorCanJoinAllSpaces shows the window on all Spaces.
	MacWindowCollectionBehaviorCanJoinAllSpaces MacWindowCollectionBehavior = 1 << 0
	// MacWindowCollectionBehaviorMoveToActiveSpace moves the window to the active Space when shown.
	MacWindowCollectionBehaviorMoveToActiveSpace MacWindowCollectionBehavior = 1 << 1
	// MacWindowCollectionBehaviorManaged is the default managed window behaviour.
	MacWindowCollectionBehaviorManaged MacWindowCollectionBehavior = 1 << 2
	// MacWindowCollectionBehaviorTransient marks the window as temporary.
	MacWindowCollectionBehaviorTransient MacWindowCollectionBehavior = 1 << 3
	// MacWindowCollectionBehaviorStationary keeps the window stationary during Space switches.
	MacWindowCollectionBehaviorStationary MacWindowCollectionBehavior = 1 << 4
	// MacWindowCollectionBehaviorParticipatesInCycle includes the window in Cmd+` cycling.
	MacWindowCollectionBehaviorParticipatesInCycle MacWindowCollectionBehavior = 1 << 5
	// MacWindowCollectionBehaviorIgnoresCycle excludes the window from Cmd+` cycling.
	MacWindowCollectionBehaviorIgnoresCycle MacWindowCollectionBehavior = 1 << 6
	// MacWindowCollectionBehaviorFullScreenPrimary allows the window to enter fullscreen.
	MacWindowCollectionBehaviorFullScreenPrimary MacWindowCollectionBehavior = 1 << 7
	// MacWindowCollectionBehaviorFullScreenAuxiliary allows the window to overlay fullscreen apps.
	MacWindowCollectionBehaviorFullScreenAuxiliary MacWindowCollectionBehavior = 1 << 8
	// MacWindowCollectionBehaviorFullScreenNone prevents the window from entering fullscreen (10.7+).
	MacWindowCollectionBehaviorFullScreenNone MacWindowCollectionBehavior = 1 << 9
	// MacWindowCollectionBehaviorFullScreenAllowsTiling allows side-by-side tiling (10.11+).
	MacWindowCollectionBehaviorFullScreenAllowsTiling MacWindowCollectionBehavior = 1 << 11
	// MacWindowCollectionBehaviorFullScreenDisallowsTiling prevents tiling in fullscreen (10.11+).
	MacWindowCollectionBehaviorFullScreenDisallowsTiling MacWindowCollectionBehavior = 1 << 12
)

// MacWebviewPreferences configures Safari-level webview behaviour on macOS.
type MacWebviewPreferences struct {
	// TabFocusesLinks enables keyboard navigation to links via Tab.
	TabFocusesLinks bool
	// TextInteractionEnabled allows the user to select and interact with text.
	TextInteractionEnabled bool
	// FullscreenEnabled allows the webview to enter fullscreen.
	FullscreenEnabled bool
	// AllowsBackForwardNavigationGestures enables horizontal swipe for navigation.
	AllowsBackForwardNavigationGestures bool
}

// MacTitleBar configures the macOS title bar appearance.
//
//	opts.Mac = application.MacWindow{TitleBar: application.MacTitleBarHiddenInset}
type MacTitleBar struct {
	// AppearsTransparent makes the title bar background transparent.
	AppearsTransparent bool
	// Hide removes the title bar entirely.
	Hide bool
	// HideTitle omits the window title text.
	HideTitle bool
	// FullSizeContent extends the content area behind the title bar.
	FullSizeContent bool
	// UseToolbar renders a toolbar in place of the standard title bar.
	UseToolbar bool
	// HideToolbarSeparator removes the separator line below the toolbar.
	HideToolbarSeparator bool
	// ShowToolbarWhenFullscreen keeps the toolbar visible in fullscreen mode.
	ShowToolbarWhenFullscreen bool
	// ToolbarStyle controls how the toolbar relates to the title bar.
	ToolbarStyle MacToolbarStyle
}

// MacTitleBarDefault is the stock macOS title bar with all decorations visible.
var MacTitleBarDefault = MacTitleBar{}

// MacTitleBarHidden hides the title text and extends content behind the title bar,
// while keeping the traffic-light window controls visible.
var MacTitleBarHidden = MacTitleBar{
	AppearsTransparent: true,
	HideTitle:          true,
	FullSizeContent:    true,
}

// MacTitleBarHiddenInset is like MacTitleBarHidden but uses an inset toolbar so the
// traffic lights sit slightly further from the window edge.
var MacTitleBarHiddenInset = MacTitleBar{
	AppearsTransparent:   true,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           true,
	HideToolbarSeparator: true,
}

// MacTitleBarHiddenInsetUnified is like MacTitleBarHiddenInset but merges the toolbar
// and title bar into a single unified row.
var MacTitleBarHiddenInsetUnified = MacTitleBar{
	AppearsTransparent:   true,
	HideTitle:            true,
	FullSizeContent:      true,
	UseToolbar:           true,
	HideToolbarSeparator: true,
	ToolbarStyle:         MacToolbarStyleUnified,
}

// MacWindow contains macOS-specific window configuration.
//
//	opts.Mac = application.MacWindow{
//	    Backdrop:  application.MacBackdropTranslucent,
//	    TitleBar:  application.MacTitleBarHiddenInset,
//	    Appearance: application.NSAppearanceNameDarkAqua,
//	}
type MacWindow struct {
	// Backdrop controls the translucency of the window background.
	Backdrop MacBackdrop
	// DisableShadow removes the drop shadow cast by the window.
	DisableShadow bool
	// TitleBar configures the title bar appearance.
	TitleBar MacTitleBar
	// Appearance sets a specific NSAppearance for the window.
	Appearance MacAppearanceType
	// InvisibleTitleBarHeight sets the height (in points) of a draggable but
	// invisible title bar region at the top of the content area.
	InvisibleTitleBarHeight int
	// EventMapping translates platform window events to common event types.
	EventMapping map[events.WindowEventType]events.WindowEventType
	// EnableFraudulentWebsiteWarnings shows browser-level phishing warnings.
	// Default: false
	EnableFraudulentWebsiteWarnings bool
	// WebviewPreferences configures Safari webview-level preferences.
	WebviewPreferences MacWebviewPreferences
	// WindowLevel controls the z-order stacking group of the window.
	WindowLevel MacWindowLevel
	// CollectionBehavior controls how the window interacts with Spaces and fullscreen.
	CollectionBehavior MacWindowCollectionBehavior
	// LiquidGlass configures the Liquid Glass visual effect (macOS 15+).
	LiquidGlass MacLiquidGlass
}

/******** Linux Options ********/

// WebviewGpuPolicy controls hardware acceleration for the Linux webview.
//
//	opts.Linux = application.LinuxWindow{WebviewGpuPolicy: application.WebviewGpuPolicyAlways}
type WebviewGpuPolicy int

const (
	// WebviewGpuPolicyAlways always enables hardware acceleration.
	WebviewGpuPolicyAlways WebviewGpuPolicy = iota
	// WebviewGpuPolicyOnDemand enables acceleration as requested by web content.
	WebviewGpuPolicyOnDemand
	// WebviewGpuPolicyNever always disables hardware acceleration.
	WebviewGpuPolicyNever
)

// LinuxMenuStyle controls how the application menu is displayed on Linux (GTK4 only).
// On GTK3 builds this option is ignored and MenuBar style is always used.
//
//	opts.Linux = application.LinuxWindow{MenuStyle: application.LinuxMenuStylePrimaryMenu}
type LinuxMenuStyle int

const (
	// LinuxMenuStyleMenuBar shows a traditional menu bar below the title bar.
	LinuxMenuStyleMenuBar LinuxMenuStyle = iota
	// LinuxMenuStylePrimaryMenu shows a primary menu button in the header bar (GNOME style).
	LinuxMenuStylePrimaryMenu
)

// LinuxWindow contains Linux-specific window configuration.
//
//	opts.Linux = application.LinuxWindow{
//	    WindowIsTranslucent: true,
//	    WebviewGpuPolicy:    application.WebviewGpuPolicyAlways,
//	}
type LinuxWindow struct {
	// Icon is the window icon shown when the window is minimised.
	// Provide PNG-encoded image data.
	Icon []byte

	// WindowIsTranslucent makes the window background transparent.
	WindowIsTranslucent bool

	// WebviewGpuPolicy sets the hardware acceleration policy for the webview.
	// Defaults to WebviewGpuPolicyNever when LinuxWindow is nil in options.
	WebviewGpuPolicy WebviewGpuPolicy

	// WindowDidMoveDebounceMS is the debounce time in milliseconds for the
	// WindowDidMove event.
	WindowDidMoveDebounceMS uint16

	// Menu is the window-level menu.
	Menu *Menu

	// MenuStyle controls how the menu is displayed (GTK4 only; ignored on GTK3).
	MenuStyle LinuxMenuStyle
}
