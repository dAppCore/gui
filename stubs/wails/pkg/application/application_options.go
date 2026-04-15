package application

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Options is the top-level application configuration passed to New().
// app := application.New(application.Options{Name: "MyApp", Services: []Service{...}})
type Options struct {
	// Name is displayed in the default about box.
	Name string

	// Description is displayed in the default about box.
	Description string

	// Icon is the application icon shown in the about box.
	Icon []byte

	// Mac contains macOS-specific options.
	Mac MacOptions

	// Windows contains Windows-specific options.
	Windows WindowsOptions

	// Linux contains Linux-specific options.
	Linux LinuxOptions

	// IOS contains iOS-specific options.
	IOS IOSOptions

	// Android contains Android-specific options.
	Android AndroidOptions

	// Services lists bound Go services exposed to the frontend.
	Services []Service

	// MarshalError serialises service method errors to JSON.
	// Return nil to fall back to the default error handler.
	MarshalError func(error) []byte

	// BindAliases maps alias IDs to bound method IDs.
	// Example: map[uint32]uint32{1: 1411160069}
	BindAliases map[uint32]uint32

	// Logger is the Wails system logger used by the runtime.
	Logger *slog.Logger

	// LogLevel defines the desired system log level when Logger is nil.
	LogLevel slog.Level

	// Assets configures the embedded asset server.
	Assets AssetOptions

	// Flags are key/value pairs available to the frontend at startup.
	Flags map[string]any

	// PanicHandler is called when an unhandled panic occurs.
	PanicHandler func(*PanicDetails)

	// DisableDefaultSignalHandler prevents Wails from handling OS signals.
	DisableDefaultSignalHandler bool

	// KeyBindings maps accelerator strings to window callbacks.
	// Example: map[string]func(Window){"Ctrl+Q": func(w Window) { w.Close() }}
	KeyBindings map[string]func(window Window)

	// OnShutdown is called before the application terminates.
	OnShutdown func()

	// PostShutdown is called after shutdown, just before process exit.
	PostShutdown func()

	// ShouldQuit is called when the user attempts to quit.
	// Return false to prevent the application from quitting.
	ShouldQuit func() bool

	// RawMessageHandler handles raw messages from the frontend.
	RawMessageHandler func(window Window, message string, originInfo *OriginInfo)

	// WarningHandler is called when a non-fatal warning occurs.
	WarningHandler func(string)

	// ErrorHandler is called when a non-fatal error occurs.
	ErrorHandler func(err error)

	// FileAssociations lists file extensions associated with this application.
	// Each extension must include the leading dot, e.g. ".txt".
	FileAssociations []string

	// SingleInstance configures single-instance enforcement.
	SingleInstance *SingleInstanceOptions

	// Transport configures the IPC transport implementation.
	Transport Transport

	// Server configures the headless HTTP server (enabled via the "server" build tag).
	Server ServerOptions
}

// ServerOptions configures the headless HTTP server started in server mode.
// opts.Server = application.ServerOptions{Host: "0.0.0.0", Port: 8080}
type ServerOptions struct {
	// Host is the bind address. Defaults to "localhost".
	Host string

	// Port is the TCP port. Defaults to 8080.
	Port int

	// ReadTimeout is the maximum duration for reading a complete request.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out a response write.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum duration to wait for the next request.
	IdleTimeout time.Duration

	// ShutdownTimeout is the maximum duration to wait for active connections on shutdown.
	ShutdownTimeout time.Duration

	// TLS configures HTTPS. If nil, plain HTTP is used.
	TLS *TLSOptions
}

// TLSOptions configures HTTPS for the headless server.
// opts.Server.TLS = &application.TLSOptions{CertFile: "cert.pem", KeyFile: "key.pem"}
type TLSOptions struct {
	// CertFile is the path to the TLS certificate file.
	CertFile string

	// KeyFile is the path to the TLS private key file.
	KeyFile string
}

// AssetOptions configures the embedded asset server.
// opts.Assets = application.AssetOptions{Handler: http.FileServer(http.FS(assets))}
type AssetOptions struct {
	// Handler serves all content to the WebView.
	Handler http.Handler

	// Middleware injects into the asset server request chain before Wails internals.
	Middleware Middleware

	// DisableLogging suppresses per-request asset server log output.
	DisableLogging bool
}

// Middleware is an HTTP middleware function for the asset server.
// type Middleware func(next http.Handler) http.Handler
type Middleware func(next http.Handler) http.Handler

// MessageProcessor is the placeholder type passed to Transport.Start.
// The stub does not perform runtime IPC processing, but the type exists so
// callers can satisfy the same API surface as real Wails.
type MessageProcessor struct{}

// Transport is the custom IPC transport contract supported by Wails.
//
//	type MyTransport struct{}
//	func (t *MyTransport) Start(ctx context.Context, processor *application.MessageProcessor) error { return nil }
//	func (t *MyTransport) Stop() error { return nil }
type Transport interface {
	Start(ctx context.Context, processor *MessageProcessor) error
	Stop() error
}

// AssetServerTransport optionally allows a custom transport to serve app assets.
type AssetServerTransport interface {
	Transport
	ServeAssets(assetHandler http.Handler) error
}

// ChainMiddleware composes multiple Middleware values into a single Middleware.
// chain := application.ChainMiddleware(authMiddleware, loggingMiddleware)
func ChainMiddleware(middleware ...Middleware) Middleware {
	return func(handler http.Handler) http.Handler {
		for index := len(middleware) - 1; index >= 0; index-- {
			handler = middleware[index](handler)
		}
		return handler
	}
}

// PanicDetails carries information about an unhandled panic.
// opts.PanicHandler = func(d *application.PanicDetails) { log(d.StackTrace) }
type PanicDetails struct {
	StackTrace     string
	Error          error
	FullStackTrace string
}

// OriginInfo carries frame-origin metadata for raw message handling.
// opts.RawMessageHandler = func(w Window, msg string, o *application.OriginInfo) { ... }
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}

// SingleInstanceOptions configures single-instance enforcement.
// opts.SingleInstance = &application.SingleInstanceOptions{UniqueID: "com.example.myapp"}
type SingleInstanceOptions struct {
	// UniqueID identifies the application instance, e.g. "com.example.myapp".
	UniqueID string

	// OnSecondInstanceLaunch is called when a second instance attempts to start.
	OnSecondInstanceLaunch func(data SecondInstanceData)

	// AdditionalData passes custom data from the second instance to the first.
	AdditionalData map[string]string
}

// SecondInstanceData describes a second-instance launch event.
type SecondInstanceData struct {
	Args           []string          `json:"args"`
	WorkingDir     string            `json:"workingDir"`
	AdditionalData map[string]string `json:"additionalData,omitempty"`
}

// ActivationPolicy controls when a macOS application activates.
// opts.Mac.ActivationPolicy = application.ActivationPolicyAccessory
type ActivationPolicy int

const (
	// ActivationPolicyRegular is used for applications with a main window.
	ActivationPolicyRegular ActivationPolicy = iota
	// ActivationPolicyAccessory is used for menu-bar or background applications.
	ActivationPolicyAccessory
	// ActivationPolicyProhibited prevents the application from activating.
	ActivationPolicyProhibited
)

// MacOptions contains macOS-specific application options.
// opts.Mac = application.MacOptions{ActivationPolicy: application.ActivationPolicyRegular}
type MacOptions struct {
	// ActivationPolicy controls how and when the application becomes active.
	ActivationPolicy ActivationPolicy

	// ApplicationShouldTerminateAfterLastWindowClosed quits the app when
	// the last window closes, matching standard macOS behaviour.
	ApplicationShouldTerminateAfterLastWindowClosed bool
}

// WindowsOptions contains Windows-specific application options.
// opts.Windows = application.WindowsOptions{WndClass: "MyAppWindow"}
type WindowsOptions struct {
	// WndClass is the Windows window class name. Defaults to "WailsWebviewWindow".
	WndClass string

	// WndProcInterceptor intercepts all Win32 messages in the main message loop.
	WndProcInterceptor func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnCode uintptr, shouldReturn bool)

	// DisableQuitOnLastWindowClosed prevents auto-quit when the last window closes.
	DisableQuitOnLastWindowClosed bool

	// WebviewUserDataPath sets the WebView2 user-data directory.
	// Defaults to %APPDATA%\[BinaryName.exe].
	WebviewUserDataPath string

	// WebviewBrowserPath sets the directory containing WebView2 executables.
	// Defaults to the system-installed WebView2.
	WebviewBrowserPath string

	// EnabledFeatures lists WebView2 feature flags to enable.
	EnabledFeatures []string

	// DisabledFeatures lists WebView2 feature flags to disable.
	DisabledFeatures []string

	// AdditionalBrowserArgs are extra Chromium command-line arguments.
	// Each argument must include the "--" prefix, e.g. "--remote-debugging-port=9222".
	AdditionalBrowserArgs []string
}

// LinuxOptions contains Linux-specific application options.
// opts.Linux = application.LinuxOptions{ProgramName: "myapp"}
type LinuxOptions struct {
	// DisableQuitOnLastWindowClosed prevents auto-quit when the last window closes.
	DisableQuitOnLastWindowClosed bool

	// ProgramName is passed to g_set_prgname() for window grouping in GTK.
	ProgramName string
}

// IOSOptions contains iOS-specific application options.
// opts.IOS = application.IOSOptions{EnableInlineMediaPlayback: true}
type IOSOptions struct {
	// DisableInputAccessoryView hides the Next/Previous/Done toolbar above the keyboard.
	DisableInputAccessoryView bool

	// DisableScroll disables WebView scrolling.
	DisableScroll bool

	// DisableBounce disables the iOS overscroll bounce effect.
	DisableBounce bool

	// DisableScrollIndicators hides scroll indicator bars.
	DisableScrollIndicators bool

	// EnableBackForwardNavigationGestures enables swipe navigation gestures.
	EnableBackForwardNavigationGestures bool

	// DisableLinkPreview disables long-press link preview (peek and pop).
	DisableLinkPreview bool

	// EnableInlineMediaPlayback allows video to play inline rather than full-screen.
	EnableInlineMediaPlayback bool

	// EnableAutoplayWithoutUserAction allows media to autoplay without a user gesture.
	EnableAutoplayWithoutUserAction bool

	// DisableInspectable disables the Safari remote web inspector.
	DisableInspectable bool

	// UserAgent overrides the WKWebView user agent string.
	UserAgent string

	// ApplicationNameForUserAgent appends a name to the user agent. Defaults to "wails.io".
	ApplicationNameForUserAgent string

	// AppBackgroundColourSet enables BackgroundColour for the main iOS window.
	AppBackgroundColourSet bool

	// BackgroundColour is applied to the iOS app window before any WebView is created.
	BackgroundColour RGBA

	// EnableNativeTabs enables a native UITabBar at the bottom of the screen.
	EnableNativeTabs bool

	// NativeTabsItems configures the labels and SF Symbol icons for the native UITabBar.
	NativeTabsItems []NativeTabItem
}

// NativeTabItem describes a single tab in the iOS native UITabBar.
// item := application.NativeTabItem{Title: "Home", SystemImage: application.NativeTabIconHouse}
type NativeTabItem struct {
	Title       string        `json:"Title"`
	SystemImage NativeTabIcon `json:"SystemImage"`
}

// NativeTabIcon is an SF Symbols name used for a UITabBar icon.
// tab := application.NativeTabItem{SystemImage: application.NativeTabIconGear}
type NativeTabIcon string

const (
	NativeTabIconNone    NativeTabIcon = ""
	NativeTabIconHouse   NativeTabIcon = "house"
	NativeTabIconGear    NativeTabIcon = "gear"
	NativeTabIconStar    NativeTabIcon = "star"
	NativeTabIconPerson  NativeTabIcon = "person"
	NativeTabIconBell    NativeTabIcon = "bell"
	NativeTabIconMagnify NativeTabIcon = "magnifyingglass"
	NativeTabIconList    NativeTabIcon = "list.bullet"
	NativeTabIconFolder  NativeTabIcon = "folder"
)

// AndroidOptions contains Android-specific application options.
// opts.Android = application.AndroidOptions{EnableZoom: true}
type AndroidOptions struct {
	// DisableScroll disables WebView scrolling.
	DisableScroll bool

	// DisableOverscroll disables the overscroll bounce effect.
	DisableOverscroll bool

	// EnableZoom enables pinch-to-zoom in the WebView.
	EnableZoom bool

	// UserAgent overrides the WebView user agent string.
	UserAgent string

	// BackgroundColour sets the WebView background colour.
	BackgroundColour RGBA

	// DisableHardwareAcceleration disables GPU acceleration for the WebView.
	DisableHardwareAcceleration bool
}
