package application

// Handler is a stub for net/http.Handler.
// In the real Wails runtime this is http.Handler; the stub replaces it with
// an interface so the application package compiles without importing net/http.
//
//	var h Handler = myHTTPHandler
type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}

// ResponseWriter is a minimal stub for http.ResponseWriter.
type ResponseWriter interface {
	Header() map[string][]string
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

// Request is a minimal stub for *http.Request.
type Request struct {
	Method string
	URL    string
	Header map[string][]string
	Body   []byte
}

// FS is a stub for fs.FS (filesystem abstraction).
// In the real Wails runtime this is io/fs.FS.
type FS interface {
	Open(name string) (interface{ Read([]byte) (int, error) }, error)
}

// Duration is a stub for time.Duration (nanoseconds).
type Duration = int64

// LogLevel is a stub for slog.Level.
type LogLevel = int

// Logger is a stub for *slog.Logger.
// In production this is the standard library structured logger.
type SlogLogger struct{}

// Middleware defines HTTP middleware applied to the AssetServer.
// The handler passed as next is the next handler in the chain.
//
//	Middleware: func(next application.Handler) application.Handler { return myHandler }
type Middleware func(next Handler) Handler

// ChainMiddleware chains multiple middlewares into one.
//
//	chained := application.ChainMiddleware(auth, logging, cors)
func ChainMiddleware(middleware ...Middleware) Middleware {
	return func(h Handler) Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			h = middleware[i](h)
		}
		return h
	}
}

// AssetFileServerFS returns a handler serving assets from an FS.
// In the stub this returns nil — no real file serving occurs.
//
//	opts.Assets.Handler = application.AssetFileServerFS(embedFS)
func AssetFileServerFS(assets FS) Handler { return nil }

// BundledAssetFileServer returns a handler serving bundled assets.
// In the stub this returns nil — no real file serving occurs.
//
//	opts.Assets.Handler = application.BundledAssetFileServer(embedFS)
func BundledAssetFileServer(assets FS) Handler { return nil }

// ActivationPolicy controls the macOS application activation policy.
//
//	Mac: MacOptions{ActivationPolicy: ActivationPolicyAccessory}
type ActivationPolicy int

const (
	// ActivationPolicyRegular is for applications with a user interface.
	ActivationPolicyRegular ActivationPolicy = iota
	// ActivationPolicyAccessory is for menu-bar or background applications.
	ActivationPolicyAccessory
	// ActivationPolicyProhibited disables activation entirely.
	ActivationPolicyProhibited
)

// NativeTabIcon is an SF Symbols name string used for iOS tab bar icons.
//
//	NativeTabsItems: []NativeTabItem{{Title: "Home", SystemImage: NativeTabIconHouse}}
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

// NativeTabItem describes a single iOS UITabBar item.
//
//	NativeTabsItems: []NativeTabItem{{Title: "Settings", SystemImage: NativeTabIconGear}}
type NativeTabItem struct {
	Title       string        `json:"Title"`
	SystemImage NativeTabIcon `json:"SystemImage"`
}

// PanicDetails carries the information delivered to the panic handler.
//
//	opts.PanicHandler = func(details *application.PanicDetails) { log(details.Message) }
type PanicDetails struct {
	// Message is the string form of the recovered panic value.
	Message string
	// Stack is the formatted goroutine stack trace.
	Stack string
}

// Transport is the interface for custom IPC transport layers.
// Implement this to replace the default HTTP fetch + js.Exec transport.
//
//	opts.Transport = myWebSocketTransport
type Transport interface{}

// SingleInstanceOptions configures the single-instance lock behaviour.
//
//	opts.SingleInstance = &application.SingleInstanceOptions{UniqueID: "com.example.myapp"}
type SingleInstanceOptions struct {
	// UniqueID is the identifier used to detect duplicate instances.
	UniqueID string
	// OnSecondInstanceLaunch is called in the first instance when a second
	// one starts. Receives the arguments passed to the second instance.
	OnSecondInstanceLaunch func(secondInstanceData SecondInstanceData)
}

// SecondInstanceData carries data from a second application instance launch.
//
//	func handler(data application.SecondInstanceData) { openFile(data.Args[0]) }
type SecondInstanceData struct {
	// Args are the command-line arguments of the second instance.
	Args []string
	// WorkingDirectory is the cwd of the second instance.
	WorkingDirectory string
}

// OriginInfo carries the origin details of a frontend message.
//
//	opts.RawMessageHandler = func(w Window, msg string, info *application.OriginInfo) {}
type OriginInfo struct {
	// Origin is the origin of the frame that sent the message.
	Origin string
	// TopOrigin is the origin of the top-level frame.
	TopOrigin string
	// IsMainFrame is true when the message came from the main frame.
	IsMainFrame bool
}

// AssetOptions configures the embedded asset server.
//
//	opts.Assets = application.AssetOptions{Handler: application.AssetFileServerFS(embedFS)}
type AssetOptions struct {
	// Handler serves all content to the WebView.
	Handler Handler

	// Middleware hooks into the AssetServer request chain.
	// Multiple middlewares can be composed with ChainMiddleware.
	Middleware Middleware

	// DisableLogging suppresses per-request AssetServer log output.
	DisableLogging bool
}

// TLSOptions configures HTTPS for the headless server.
//
//	opts.Server.TLS = &application.TLSOptions{CertFile: "cert.pem", KeyFile: "key.pem"}
type TLSOptions struct {
	// CertFile is the path to the TLS certificate file.
	CertFile string
	// KeyFile is the path to the TLS private key file.
	KeyFile string
}

// ServerOptions configures the HTTP server used in headless (server) mode.
// Enable server mode by building with: go build -tags server
//
//	opts.Server = application.ServerOptions{Host: "0.0.0.0", Port: 8080}
type ServerOptions struct {
	// Host is the address to bind to. Defaults to "localhost".
	Host string
	// Port is the port to listen on. Defaults to 8080.
	Port int
	// ReadTimeout is the maximum duration for reading a request (nanoseconds).
	ReadTimeout Duration
	// WriteTimeout is the maximum duration for writing a response (nanoseconds).
	WriteTimeout Duration
	// IdleTimeout is the maximum idle connection duration (nanoseconds).
	IdleTimeout Duration
	// ShutdownTimeout is the maximum time to wait for graceful shutdown (nanoseconds).
	ShutdownTimeout Duration
	// TLS configures HTTPS. If nil, plain HTTP is used.
	TLS *TLSOptions
}

// MacOptions contains macOS-specific application configuration.
//
//	opts.Mac = application.MacOptions{ActivationPolicy: application.ActivationPolicyRegular}
type MacOptions struct {
	// ActivationPolicy controls how the app interacts with the Dock and menu bar.
	ActivationPolicy ActivationPolicy
	// ApplicationShouldTerminateAfterLastWindowClosed quits the app when the
	// last window closes (matches NSApplicationDelegate behaviour).
	ApplicationShouldTerminateAfterLastWindowClosed bool
}

// WindowsOptions contains Windows-specific application configuration.
//
//	opts.Windows = application.WindowsOptions{WndClass: "MyAppWindow"}
type WindowsOptions struct {
	// WndClass is the Win32 window class name. Default: WailsWebviewWindow.
	WndClass string
	// WndProcInterceptor intercepts all Win32 messages for the application.
	WndProcInterceptor func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnCode uintptr, shouldReturn bool)
	// DisableQuitOnLastWindowClosed prevents auto-quit when the last window closes.
	DisableQuitOnLastWindowClosed bool
	// WebviewUserDataPath is the directory for WebView2 user data.
	WebviewUserDataPath string
	// WebviewBrowserPath is the directory containing WebView2 executables.
	WebviewBrowserPath string
	// EnabledFeatures lists WebView2 feature flags to enable.
	EnabledFeatures []string
	// DisabledFeatures lists WebView2 feature flags to disable.
	DisabledFeatures []string
	// AdditionalBrowserArgs are extra browser arguments (must include "--" prefix).
	AdditionalBrowserArgs []string
}

// LinuxOptions contains Linux-specific application configuration.
//
//	opts.Linux = application.LinuxOptions{ProgramName: "myapp"}
type LinuxOptions struct {
	// DisableQuitOnLastWindowClosed prevents auto-quit when the last window closes.
	DisableQuitOnLastWindowClosed bool
	// ProgramName sets g_set_prgname() for the window manager.
	ProgramName string
}

// IOSOptions contains iOS-specific application configuration.
//
//	opts.IOS = application.IOSOptions{EnableInlineMediaPlayback: true}
type IOSOptions struct {
	// DisableInputAccessoryView hides the iOS keyboard accessory bar.
	DisableInputAccessoryView bool
	// DisableScroll disables WebView scrolling.
	DisableScroll bool
	// DisableBounce disables the WebView bounce effect.
	DisableBounce bool
	// DisableScrollIndicators hides scroll indicators.
	DisableScrollIndicators bool
	// EnableBackForwardNavigationGestures enables swipe navigation.
	EnableBackForwardNavigationGestures bool
	// DisableLinkPreview disables link long-press previews.
	DisableLinkPreview bool
	// EnableInlineMediaPlayback allows media to play inline.
	EnableInlineMediaPlayback bool
	// EnableAutoplayWithoutUserAction allows media autoplay without a gesture.
	EnableAutoplayWithoutUserAction bool
	// DisableInspectable disables the Safari Web Inspector.
	DisableInspectable bool
	// UserAgent overrides the WebView user agent string.
	UserAgent string
	// ApplicationNameForUserAgent is appended to the user agent. Default: "wails.io".
	ApplicationNameForUserAgent string
	// AppBackgroundColourSet enables the custom BackgroundColour below.
	AppBackgroundColourSet bool
	// BackgroundColour is the app window background before WebView creation.
	BackgroundColour RGBA
	// EnableNativeTabs shows a native iOS UITabBar.
	EnableNativeTabs bool
	// NativeTabsItems configures the UITabBar items. Auto-enables tabs when non-empty.
	NativeTabsItems []NativeTabItem
}

// AndroidOptions contains Android-specific application configuration.
//
//	opts.Android = application.AndroidOptions{EnableZoom: true}
type AndroidOptions struct {
	// DisableScroll disables WebView scrolling.
	DisableScroll bool
	// DisableOverscroll disables the overscroll bounce effect.
	DisableOverscroll bool
	// EnableZoom allows pinch-to-zoom in the WebView.
	EnableZoom bool
	// UserAgent sets a custom user agent string.
	UserAgent string
	// BackgroundColour sets the WebView background colour.
	BackgroundColour RGBA
	// DisableHardwareAcceleration disables hardware acceleration.
	DisableHardwareAcceleration bool
}

// Options is the top-level application configuration passed to New().
//
//	app := application.New(application.Options{
//	    Name:     "MyApp",
//	    Assets:   application.AssetOptions{Handler: application.AssetFileServerFS(embedFS)},
//	    Services: []application.Service{application.NewService(&myService{})},
//	})
type Options struct {
	// Name is the application name shown in the default about box.
	Name string
	// Description is shown in the default about box.
	Description string
	// Icon is the application icon bytes used in the about box.
	Icon []byte

	// Mac is the macOS-specific configuration.
	Mac MacOptions
	// Windows is the Windows-specific configuration.
	Windows WindowsOptions
	// Linux is the Linux-specific configuration.
	Linux LinuxOptions
	// IOS is the iOS-specific configuration.
	IOS IOSOptions
	// Android is the Android-specific configuration.
	Android AndroidOptions

	// Services are the bound Go service instances exposed to the frontend.
	Services []Service
	// MarshalError serialises error values from service methods to JSON.
	// A nil return falls back to the default error handler.
	MarshalError func(error) []byte
	// BindAliases maps alias IDs to bound method IDs.
	// Example: map[uint32]uint32{1: 1411160069}
	BindAliases map[uint32]uint32

	// Logger is the structured logger for Wails system messages.
	// If nil, a default logger is used.
	Logger *SlogLogger
	// LogLevel sets the log level for the Wails system logger.
	LogLevel LogLevel

	// Assets configures the embedded asset server.
	Assets AssetOptions
	// Flags are key-value pairs exposed to the frontend.
	Flags map[string]any

	// PanicHandler is called when a panic occurs in a service method.
	PanicHandler func(*PanicDetails)
	// DisableDefaultSignalHandler disables the built-in SIGINT/SIGTERM handler.
	DisableDefaultSignalHandler bool

	// KeyBindings maps accelerator strings to window callbacks.
	KeyBindings map[string]func(window Window)
	// OnShutdown is called before the application terminates.
	OnShutdown func()
	// PostShutdown is called after shutdown completes, just before process exit.
	PostShutdown func()
	// ShouldQuit is called when the user requests quit. Return false to cancel.
	ShouldQuit func() bool
	// RawMessageHandler handles raw messages sent from the frontend.
	RawMessageHandler func(window Window, message string, originInfo *OriginInfo)
	// WarningHandler is called when a non-fatal warning occurs.
	WarningHandler func(string)
	// ErrorHandler is called when an error occurs.
	ErrorHandler func(err error)

	// FileAssociations lists file extensions associated with this application.
	// Example: []string{".txt", ".md"} — the leading dot is required.
	FileAssociations []string
	// SingleInstance configures single-instance enforcement.
	SingleInstance *SingleInstanceOptions

	// Transport provides a custom IPC transport layer.
	// When nil, the default HTTP fetch transport is used.
	Transport Transport

	// Server configures the headless HTTP server (requires -tags server).
	Server ServerOptions
}
