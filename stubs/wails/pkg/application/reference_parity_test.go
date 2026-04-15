package application

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

type lifecycleService struct {
	started   int
	shutdowns int
	ctx       context.Context
	options   ServiceOptions
}

func (s *lifecycleService) ServiceStartup(ctx context.Context, options ServiceOptions) error {
	s.started++
	s.ctx = ctx
	s.options = options
	return nil
}

func (s *lifecycleService) ServiceShutdown() error {
	s.shutdowns++
	return nil
}

type stubRequest struct {
	url    string
	method string
	header http.Header
	body   string
	resp   *noopResponseWriter
}

func (r *stubRequest) URL() (string, error)         { return r.url, nil }
func (r *stubRequest) Method() (string, error)      { return r.method, nil }
func (r *stubRequest) Header() (http.Header, error) { return r.header, nil }
func (r *stubRequest) Body() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(r.body)), nil
}
func (r *stubRequest) Response() webview.ResponseWriter { return r.resp }
func (r *stubRequest) Close() error                     { return nil }

func TestReferenceParity_NewRunQuitLifecycle(t *testing.T) {
	globalApplication = nil

	service := &lifecycleService{}
	shutdownCalled := 0
	postShutdownCalled := 0

	app := New(Options{
		Services: []Service{NewServiceWithOptions(service, ServiceOptions{Name: "lifecycle"})},
		OnShutdown: func() {
			shutdownCalled++
		},
		PostShutdown: func() {
			postShutdownCalled++
		},
	})

	if app.Config().Name != "My Wails Application" {
		t.Fatalf("expected default app name, got %q", app.Config().Name)
	}
	if app.Config().Description != "An application written using Wails" {
		t.Fatalf("expected default description, got %q", app.Config().Description)
	}
	if app.Config().Windows.WndClass != "WailsWebviewWindow" {
		t.Fatalf("expected default window class, got %q", app.Config().Windows.WndClass)
	}

	if err := app.Run(); err != nil {
		t.Fatalf("Run() failed: %v", err)
	}
	if service.started != 1 {
		t.Fatalf("expected service startup once, got %d", service.started)
	}
	if service.ctx == nil {
		t.Fatalf("expected service startup context")
	}
	if got := getServiceName(NewServiceWithOptions(service, ServiceOptions{Name: "lifecycle"})); got != "lifecycle" {
		t.Fatalf("unexpected service name: %q", got)
	}

	app.Quit()
	if service.shutdowns != 1 {
		t.Fatalf("expected service shutdown once, got %d", service.shutdowns)
	}
	if shutdownCalled != 1 {
		t.Fatalf("expected OnShutdown callback once, got %d", shutdownCalled)
	}
	if postShutdownCalled != 1 {
		t.Fatalf("expected PostShutdown callback once, got %d", postShutdownCalled)
	}
}

func TestReferenceParity_WebViewAssetRequestHeaders(t *testing.T) {
	request := &webViewAssetRequest{
		Request: &stubRequest{
			url:    "https://example.com/app",
			method: http.MethodPost,
			header: http.Header{"X-Test": []string{"ok"}},
			body:   "payload",
			resp:   &noopResponseWriter{},
		},
		windowId:   42,
		windowName: "main",
	}

	url, err := request.URL()
	if err != nil || url != "https://example.com/app" {
		t.Fatalf("unexpected URL result: %q %v", url, err)
	}
	method, err := request.Method()
	if err != nil || method != http.MethodPost {
		t.Fatalf("unexpected Method result: %q %v", method, err)
	}
	headers, err := request.Header()
	if err != nil {
		t.Fatalf("Header() failed: %v", err)
	}
	if got := headers.Get(webViewRequestHeaderWindowId); got != "42" {
		t.Fatalf("unexpected window ID header: %q", got)
	}
	if got := headers.Get(webViewRequestHeaderWindowName); got != "main" {
		t.Fatalf("unexpected window name header: %q", got)
	}
	body, err := request.Body()
	if err != nil {
		t.Fatalf("Body() failed: %v", err)
	}
	defer body.Close()
	payload, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("ReadAll() failed: %v", err)
	}
	if string(payload) != "payload" {
		t.Fatalf("unexpected body payload: %q", payload)
	}
	if request.Response() == nil {
		t.Fatalf("expected Response() to return a writer")
	}
}

func TestReferenceParity_WindowKeyAndMenuClickHandlers(t *testing.T) {
	globalApplication = nil

	app := New(Options{})
	window := app.Window.NewWithOptions(WebviewWindowOptions{Name: "main"})

	windowTriggered := false
	globalTriggered := false
	window.RegisterKeyBinding("CmdOrCtrl+K", func(Window) { windowTriggered = true })
	app.KeyBinding.Add("CmdOrCtrl+K", func(Window) { globalTriggered = true })

	app.handleWindowKeyEvent(&windowKeyEvent{
		windowId:          window.ID(),
		acceleratorString: "CmdOrCtrl+K",
	})

	if !windowTriggered {
		t.Fatalf("expected window key binding callback to fire")
	}
	if !globalTriggered {
		t.Fatalf("expected global key binding callback to fire")
	}

	clicked := false
	menuItem := NewMenuItem("Open")
	menuItem.OnClick(func(ctx *Context) {
		clicked = ctx != nil && ctx.clickedMenuItem == menuItem
	})
	app.Menu.handleMenuItemClicked(menuItem.id)
	if !clicked {
		t.Fatalf("expected menu item click handler to fire")
	}
}
