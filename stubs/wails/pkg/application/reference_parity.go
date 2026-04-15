package application

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
)

type (
	runnable interface {
		Run()
	}

	clipboardImpl interface {
		SetText(string) bool
		Text() (string, bool)
	}

	menuImpl interface {
		Update()
		Destroy()
	}

	menuItemImpl interface {
		Destroy()
	}

	systemTrayImpl interface {
		Run()
		Destroy()
	}
)

type windowMessage struct {
	windowId   uint
	message    string
	originInfo *OriginInfo
}

type dragAndDropMessage struct {
	windowId   uint
	filenames  []string
	X          int
	Y          int
	DropTarget *DropTargetDetails
}

type webViewAssetRequest struct {
	Request    webview.Request
	windowId   uint
	windowName string
}

type windowKeyEvent struct {
	windowId          uint
	acceleratorString string
}

type windowEvent struct {
	WindowID uint
	EventID  uint
}

type eventHook struct {
	callback func(*CustomEvent)
}

const (
	webViewRequestHeaderWindowId   = "x-wails-window-id"
	webViewRequestHeaderWindowName = "x-wails-window-name"
)

var (
	windowMessageBuffer     = make(chan *windowMessage, 5)
	windowDragAndDropBuffer = make(chan *dragAndDropMessage, 5)
	windowKeyEvents         = make(chan *windowKeyEvent, 5)
	applicationEvents       = make(chan *ApplicationEvent, 5)
	windowEvents            = make(chan *windowEvent, 5)
	menuItemClicked         = make(chan uint, 5)

	menuItemRegistry sync.Map
	systemTrayIDs    atomic.Uint64
)

type noopResponseWriter struct {
	header http.Header
}

func (w *noopResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *noopResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (w *noopResponseWriter) WriteHeader(statusCode int) {}

func addDragAndDropMessage(windowId uint, filenames []string, dropTarget *DropTargetDetails) {
	message := &dragAndDropMessage{
		windowId:   windowId,
		filenames:  append([]string(nil), filenames...),
		DropTarget: dropTarget,
	}
	select {
	case windowDragAndDropBuffer <- message:
	default:
		select {
		case <-windowDragAndDropBuffer:
		default:
		}
		windowDragAndDropBuffer <- message
	}
}

func (r *webViewAssetRequest) URL() (string, error) {
	if r == nil || r.Request == nil {
		return "", nil
	}
	return r.Request.URL()
}

func (r *webViewAssetRequest) Method() (string, error) {
	if r == nil || r.Request == nil {
		return "", nil
	}
	return r.Request.Method()
}

func (r *webViewAssetRequest) Header() (http.Header, error) {
	if r == nil || r.Request == nil {
		return http.Header{}, nil
	}

	headers, err := r.Request.Header()
	if err != nil {
		return nil, err
	}
	if headers == nil {
		headers = make(http.Header)
	} else {
		headers = headers.Clone()
	}

	headers.Set(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(r.windowId), 10))
	if r.windowName != "" {
		headers.Set(webViewRequestHeaderWindowName, r.windowName)
	}
	return headers, nil
}

func (r *webViewAssetRequest) Body() (io.ReadCloser, error) {
	if r == nil || r.Request == nil {
		return io.NopCloser(strings.NewReader("")), nil
	}
	return r.Request.Body()
}

func (r *webViewAssetRequest) Response() webview.ResponseWriter {
	if r == nil || r.Request == nil {
		return &noopResponseWriter{}
	}
	if response := r.Request.Response(); response != nil {
		return response
	}
	return &noopResponseWriter{}
}

func (r *webViewAssetRequest) Close() error {
	if r == nil || r.Request == nil {
		return nil
	}
	return r.Request.Close()
}

func mergeApplicationDefaults(options *Options) {
	if options == nil {
		return
	}
	if options.Name == "" {
		options.Name = "My Wails Application"
	}
	if options.Description == "" {
		options.Description = "An application written using Wails"
	}
	if options.Windows.WndClass == "" {
		options.Windows.WndClass = "WailsWebviewWindow"
	}
}

func getServiceName(service Service) string {
	if name := strings.TrimSpace(service.Options().Name); name != "" {
		return name
	}
	if named, ok := service.Instance().(ServiceName); ok {
		if name := strings.TrimSpace(named.ServiceName()); name != "" {
			return name
		}
	}
	typ := reflect.TypeOf(service.Instance())
	if typ == nil {
		return ""
	}
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ.Name()
}

func addToMenuItemMap(menuItem *MenuItem) {
	if menuItem == nil {
		return
	}
	menuItemRegistry.Store(menuItem.id, menuItem)
}

func getMenuItemByID(id uint) *MenuItem {
	value, ok := menuItemRegistry.Load(id)
	if !ok {
		return nil
	}
	menuItem, _ := value.(*MenuItem)
	return menuItem
}

func removeMenuItemByID(id uint) {
	menuItemRegistry.Delete(id)
}

func newClipboard() *Clipboard {
	return &Clipboard{}
}

func newBrowserManager(app *App) *BrowserManager {
	initialiseAppManagers(app, Options{})
	return app.Browser
}

func newClipboardManager(app *App) *ClipboardManager {
	initialiseAppManagers(app, Options{})
	return app.Clipboard
}

func (cm *ClipboardManager) getClipboard() *Clipboard {
	return cm.instance()
}

func newContextMenuManager(app *App) *ContextMenuManager {
	initialiseAppManagers(app, Options{})
	return app.ContextMenu
}

func newDialogManager(app *App) *DialogManager {
	initialiseAppManagers(app, Options{})
	return app.Dialog
}

func newEnvironmentManager(app *App) *EnvironmentManager {
	initialiseAppManagers(app, Options{})
	return app.Environment
}

func newEventManager(app *App) *EventManager {
	initialiseAppManagers(app, Options{})
	return app.Event
}

func newKeyBindingManager(app *App) *KeyBindingManager {
	initialiseAppManagers(app, Options{})
	return app.KeyBinding
}

func (kbm *KeyBindingManager) handleWindowKeyEvent(event *windowKeyEvent) {
	if event == nil || globalApplication == nil {
		return
	}
	globalApplication.handleWindowKeyEvent(event)
}

func newMenuManager(app *App) *MenuManager {
	initialiseAppManagers(app, Options{})
	return app.Menu
}

func (mm *MenuManager) handleMenuItemClicked(menuItemID uint) {
	menuItem := getMenuItemByID(menuItemID)
	if menuItem == nil {
		if globalApplication != nil {
			globalApplication.warning("MenuItem #%d not found", menuItemID)
		}
		return
	}
	menuItem.handleClick()
}

func newScreenManager(app *App) *ScreenManager {
	initialiseAppManagers(app, Options{})
	return app.Screen
}

func newSystemTray(id uint) *SystemTray {
	tray := &SystemTray{}
	ensureTrayCompatState(tray)
	return tray
}

func newSystemTrayManager(app *App) *SystemTrayManager {
	initialiseAppManagers(app, Options{})
	return app.SystemTray
}

func (stm *SystemTrayManager) getNextID() uint {
	return uint(systemTrayIDs.Add(1))
}

func (stm *SystemTrayManager) destroy(tray *SystemTray) {
	if tray == nil {
		return
	}
	tray.Destroy()
}

func newWindowManager(app *App) *WindowManager {
	initialiseAppManagers(app, Options{})
	return app.Window
}

func (wm *WindowManager) add(window Window) {
	wm.Add(window)
}

func (wm *WindowManager) remove(windowID uint) {
	wm.Remove(windowID)
}

func (a *App) startupService(service Service) error {
	instance := service.Instance()
	startup, ok := instance.(ServiceStartup)
	if !ok {
		return nil
	}
	return startup.ServiceStartup(a.Context(), service.Options())
}

func (a *App) shutdownServices() {
	state := ensureAppCompatState(a)
	state.mu.RLock()
	services := append([]Service(nil), state.services...)
	state.mu.RUnlock()

	for i := len(services) - 1; i >= 0; i-- {
		shutdown, ok := services[i].Instance().(ServiceShutdown)
		if !ok {
			continue
		}
		if err := shutdown.ServiceShutdown(); err != nil {
			a.handleError(err)
		}
	}
}

func (a *App) cleanup() {
	a.shutdownServices()
}

func (a *App) dispatchOnMainThread(fn func()) {
	if fn != nil {
		fn()
	}
}

func (a *App) runOrDeferToAppRun(r runnable) {
	if r != nil {
		r.Run()
	}
}

func formatLogMessage(message string, args ...any) string {
	if len(args) == 0 {
		return message
	}
	return fmt.Sprintf(message, args...)
}

func (a *App) handleWarning(msg string) {
	if a == nil {
		return
	}
	if handler := a.Config().WarningHandler; handler != nil {
		handler(msg)
	}
	if a.Logger != nil {
		a.Logger.Warn(msg)
	}
}

func (a *App) handleError(err error) {
	if a == nil || err == nil {
		return
	}
	if handler := a.Config().ErrorHandler; handler != nil {
		handler(err)
	}
	if a.Logger != nil {
		a.Logger.Error(err.Error())
	}
}

func (a *App) handleFatalError(err error) {
	a.handleError(err)
}

func (a *App) info(message string, args ...any) {
	if a != nil && a.Logger != nil {
		a.Logger.Info(formatLogMessage(message, args...))
	}
}

func (a *App) debug(message string, args ...any) {
	if a != nil && a.Logger != nil {
		a.Logger.Debug(formatLogMessage(message, args...))
	}
}

func (a *App) fatal(message string, args ...any) {
	a.handleFatalError(fmt.Errorf("%s", formatLogMessage(message, args...)))
}

func (a *App) warning(message string, args ...any) {
	a.handleWarning(formatLogMessage(message, args...))
}

func (a *App) error(message string, args ...any) {
	a.handleError(fmt.Errorf("%s", formatLogMessage(message, args...)))
}

func (a *App) handleDragAndDropMessage(event *dragAndDropMessage) {
	if a == nil || event == nil || a.Window == nil {
		return
	}
	window, ok := a.Window.GetByID(event.windowId)
	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	window.handleDragAndDropMessage(event.filenames, event.DropTarget)
}

func (a *App) handleWindowMessage(event *windowMessage) {
	if a == nil || event == nil || a.Window == nil {
		return
	}
	window, ok := a.Window.GetByID(event.windowId)
	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	if handler := a.Config().RawMessageHandler; handler != nil {
		handler(window, event.message, event.originInfo)
		return
	}
	window.HandleMessage(event.message)
}

func (a *App) handleWebViewRequest(request *webViewAssetRequest) {
	if request == nil {
		return
	}
	_, _ = request.Header()
}

func (a *App) handleWindowEvent(event *windowEvent) {
	if a == nil || event == nil || a.Window == nil {
		return
	}
	window, ok := a.Window.GetByID(event.WindowID)
	if !ok {
		a.warning("WebviewWindow #%d not found", event.WindowID)
		return
	}
	window.HandleWindowEvent(event.EventID)
}

func (a *App) handleWindowKeyEvent(event *windowKeyEvent) {
	if a == nil || event == nil || a.Window == nil {
		return
	}
	window, ok := a.Window.GetByID(event.windowId)
	if !ok {
		a.warning("WebviewWindow #%d not found", event.windowId)
		return
	}
	window.HandleKeyEvent(event.acceleratorString)
	if a.KeyBinding != nil {
		a.KeyBinding.Process(event.acceleratorString, window)
	}
}

func (a *App) shouldQuit() bool {
	if a == nil {
		return true
	}
	if callback := a.Config().ShouldQuit; callback != nil {
		return callback()
	}
	return true
}

func (m *Menu) processRadioGroups() {
	var radioGroup []*MenuItem
	flush := func() {
		if len(radioGroup) == 0 {
			return
		}
		for _, item := range radioGroup {
			item.radioGroupMembers = radioGroup
		}
		radioGroup = nil
	}
	for _, item := range m.Items {
		if item == nil {
			continue
		}
		if item.itemType != menuItemTypeRadio {
			flush()
		}
		if item.itemType == menuItemTypeSubmenu && item.submenu != nil {
			item.submenu.processRadioGroups()
			continue
		}
		if item.itemType == menuItemTypeRadio {
			radioGroup = append(radioGroup, item)
		}
	}
	flush()
}

func (m *MenuItem) handleClick() {
	if m == nil || m.disabled {
		return
	}
	switch m.itemType {
	case menuItemTypeCheckbox:
		m.checked = !m.checked
	case menuItemTypeRadio:
		for _, item := range m.radioGroupMembers {
			item.checked = item == m
		}
	}

	if m.callback == nil {
		return
	}

	ctx := newContext().
		withClickedMenuItem(m).
		withChecked(m.checked)
	if m.contextMenuData != nil {
		ctx.withContextMenuData(m.contextMenuData.clone())
	}
	m.callback(ctx)
}

func (m *MenuItem) setContextData(data *ContextMenuData) {
	if m == nil {
		return
	}
	if data == nil {
		m.contextMenuData = nil
	} else {
		m.contextMenuData = data.clone()
	}
	if m.submenu != nil {
		m.submenu.setContextData(data)
	}
}
