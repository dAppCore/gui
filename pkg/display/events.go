package display

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/events"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/gorilla/websocket"
)

// EventType represents the type of event.
type EventType string

const (
	EventWindowFocus         EventType = "window.focus"
	EventWindowBlur          EventType = "window.blur"
	EventWindowMove          EventType = "window.move"
	EventWindowResize        EventType = "window.resize"
	EventWindowClose         EventType = "window.close"
	EventWindowCreate        EventType = "window.create"
	EventThemeChange         EventType = "theme.change"
	EventScreenChange        EventType = "screen.change"
	EventNotificationClick   EventType = "notification.click"
	EventTrayClick           EventType = "tray.click"
	EventTrayMenuItemClick   EventType = "tray.menuitem.click"
	EventKeybindingTriggered EventType = "keybinding.triggered"
	EventWindowFileDrop      EventType = "window.filedrop"
	EventDockVisibility      EventType = "dock.visibility-changed"
	EventAppStarted          EventType = "app.started"
	EventAppOpenedWithFile   EventType = "app.opened-with-file"
	EventAppWillTerminate    EventType = "app.will-terminate"
	EventAppActive           EventType = "app.active"
	EventAppInactive         EventType = "app.inactive"
	EventSystemPowerChange   EventType = "system.power-change"
	EventSystemSuspend       EventType = "system.suspend"
	EventSystemResume        EventType = "system.resume"
	EventContextMenuClick    EventType = "contextmenu.item-clicked"
	EventWebviewConsole      EventType = "webview.console"
	EventWebviewException    EventType = "webview.exception"
	EventCustomEvent         EventType = "custom.event"
	EventDockProgress        EventType = "dock.progress"
	EventDockBounce          EventType = "dock.bounce"
	EventNotificationAction  EventType = "notification.action"
	EventNotificationDismiss EventType = "notification.dismissed"
	EventChatConversation    EventType = "chat.conversation"
	EventChatMessage         EventType = "chat.message"
	EventChatToken           EventType = "chat.token"
	EventChatThinkingStart   EventType = "chat.thinking.start"
	EventChatThinkingAppend  EventType = "chat.thinking.append"
	EventChatThinkingEnd     EventType = "chat.thinking.end"
	EventChatToolCall        EventType = "chat.tool.call"
	EventChatToolResult      EventType = "chat.tool.result"
	EventChatImageQueued     EventType = "chat.image.queued"
)

const websocketReadTimeout = 30 * time.Second

// Event represents a display event sent to subscribers.
type Event struct {
	Type      EventType      `json:"type"`
	Timestamp int64          `json:"timestamp"`
	Window    string         `json:"window,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
}

// Subscription represents a client subscription to events.
type Subscription struct {
	ID         string      `json:"id"`
	EventTypes []EventType `json:"eventTypes"`
}

// WSEventManager manages WebSocket connections and event subscriptions.
type WSEventManager struct {
	upgrader    websocket.Upgrader
	clients     map[*websocket.Conn]*clientState
	mu          sync.RWMutex
	nextSubID   int
	eventBuffer chan Event
	closed      bool
	readTimeout time.Duration
}

// clientState tracks a client's subscriptions.
type clientState struct {
	subscriptions map[string]*Subscription
	writeMu       sync.Mutex
	mu            sync.RWMutex
}

// NewWSEventManager creates a new event manager.
func NewWSEventManager() *WSEventManager {
	em := &WSEventManager{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return trustedWebSocketOrigin(r)
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:     make(map[*websocket.Conn]*clientState),
		eventBuffer: make(chan Event, 100),
		readTimeout: websocketReadTimeout,
	}

	// Start event broadcaster
	go em.broadcaster()

	return em
}

func trustedWebSocketOrigin(r *http.Request) bool {
	if r == nil {
		return false
	}

	if r.URL == nil {
		return false
	}
	if path := strings.TrimSpace(r.URL.Path); path != "" && path != "/" && path != "/events" {
		return false
	}

	if !trustedWebSocketHost(r.Host) {
		return false
	}
	if !trustedWSRequestOrigin(r.RemoteAddr) {
		return false
	}

	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" || strings.EqualFold(origin, "null") {
		return true
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
		return trustedWebSocketHost(parsed.Host)
	case "wails", "core", "app":
		return true
	default:
		return false
	}
}

func trustedWSRequestOrigin(raw string) bool {
	if raw == "" {
		return false
	}
	host := raw
	if parsed, _, err := net.SplitHostPort(raw); err == nil {
		host = parsed
	}
	host = strings.Trim(host, "[]")
	return isLoopbackHost(host)
}

func isLoopbackHost(host string) bool {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return false
	}
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func trustedWebSocketHost(host string) bool {
	host = strings.TrimSpace(host)
	if host == "" {
		return false
	}

	name := host
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		name = parsedHost
	}
	name = strings.Trim(name, "[]")
	switch strings.ToLower(name) {
	case "localhost", "127.0.0.1", "::1":
		return true
	case "localhost:80", "localhost:443", "127.0.0.1:80", "127.0.0.1:443", "[::1]:80", "[::1]:443":
		return true
	default:
		return false
	}
}

// broadcaster sends events to all subscribed clients.
func (em *WSEventManager) broadcaster() {
	if em == nil || em.eventBuffer == nil {
		return
	}
	for event := range em.eventBuffer {
		em.mu.RLock()
		for conn, state := range em.clients {
			if state != nil && em.clientSubscribed(state, event.Type) {
				go em.sendEvent(conn, event)
			}
		}
		em.mu.RUnlock()
	}
}

// clientSubscribed checks if a client is subscribed to an event type.
func (em *WSEventManager) clientSubscribed(state *clientState, eventType EventType) bool {
	if state == nil {
		return false
	}
	state.mu.RLock()
	defer state.mu.RUnlock()

	for _, sub := range state.subscriptions {
		for _, et := range sub.EventTypes {
			if et == eventType || et == "*" {
				return true
			}
		}
	}
	return false
}

// sendEvent sends an event to a specific client.
func (em *WSEventManager) sendEvent(conn *websocket.Conn, event Event) {
	em.mu.RLock()
	state, exists := em.clients[conn]
	em.mu.RUnlock()

	if !exists || state == nil {
		return
	}

	marshalResult := core.JSONMarshal(event)
	if !marshalResult.OK {
		return
	}
	data, _ := marshalResult.Value.([]byte)

	state.writeMu.Lock()
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err := conn.WriteMessage(websocket.TextMessage, data)
	state.writeMu.Unlock()
	if err != nil {
		em.removeClient(conn)
	}
}

// HandleWebSocket handles WebSocket upgrade and connection.
func (em *WSEventManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if em == nil {
		if w != nil {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
		return
	}
	em.mu.RLock()
	closed := em.closed
	em.mu.RUnlock()
	if closed {
		if w != nil {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
		return
	}
	if !trustedWebSocketOrigin(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	conn, err := em.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	em.mu.Lock()
	if em.closed {
		em.mu.Unlock()
		_ = conn.Close()
		return
	}
	em.clients[conn] = &clientState{
		subscriptions: make(map[string]*Subscription),
	}
	em.mu.Unlock()

	em.prepareConnection(conn)

	// Handle incoming messages
	go em.handleMessages(conn)
}

func (em *WSEventManager) prepareConnection(conn *websocket.Conn) {
	if conn == nil {
		return
	}
	conn.SetReadLimit(64 * 1024)
	if em.readTimeout > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(em.readTimeout))
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(em.readTimeout))
		})
	}
}

// handleMessages processes incoming WebSocket messages.
func (em *WSEventManager) handleMessages(conn *websocket.Conn) {
	defer em.removeClient(conn)

	for {
		if em.readTimeout > 0 {
			_ = conn.SetReadDeadline(time.Now().Add(em.readTimeout))
		}
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg struct {
			Action     string      `json:"action"`
			ID         string      `json:"id,omitempty"`
			EventTypes []EventType `json:"eventTypes,omitempty"`
		}

		if unmarshalResult := core.JSONUnmarshal(message, &msg); !unmarshalResult.OK {
			em.closeWithPolicyViolation(conn, "invalid websocket message")
			return
		}

		switch msg.Action {
		case "subscribe":
			em.subscribe(conn, msg.ID, msg.EventTypes)
		case "unsubscribe":
			em.unsubscribe(conn, msg.ID)
		case "list":
			em.listSubscriptions(conn)
		default:
			em.closeWithPolicyViolation(conn, "unknown websocket action")
			return
		}
	}
}

func (em *WSEventManager) closeWithPolicyViolation(conn *websocket.Conn, reason string) {
	em.mu.RLock()
	state, exists := em.clients[conn]
	em.mu.RUnlock()
	if !exists || state == nil {
		return
	}
	state.writeMu.Lock()
	defer state.writeMu.Unlock()
	_ = conn.WriteJSON(map[string]any{
		"error":  reason,
		"status": websocket.ClosePolicyViolation,
	})
	_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, reason), time.Now().Add(2*time.Second))
	_ = conn.Close()
}

// subscribe adds a subscription for a client.
func (em *WSEventManager) subscribe(conn *websocket.Conn, id string, eventTypes []EventType) {
	em.mu.RLock()
	state, exists := em.clients[conn]
	em.mu.RUnlock()

	if !exists {
		return
	}

	// Generate ID if not provided
	if id == "" {
		em.mu.Lock()
		em.nextSubID++
		id = "sub-" + strconv.Itoa(em.nextSubID)
		em.mu.Unlock()
	}

	state.mu.Lock()
	state.subscriptions[id] = &Subscription{
		ID:         id,
		EventTypes: eventTypes,
	}
	state.mu.Unlock()

	// Send confirmation
	response := map[string]any{
		"type":       "subscribed",
		"id":         id,
		"eventTypes": eventTypes,
	}
	if marshalResult := core.JSONMarshal(response); marshalResult.OK {
		responseData, _ := marshalResult.Value.([]byte)
		em.writeClientMessage(state, conn, responseData)
	}
}

// unsubscribe removes a subscription for a client.
func (em *WSEventManager) unsubscribe(conn *websocket.Conn, id string) {
	em.mu.RLock()
	state, exists := em.clients[conn]
	em.mu.RUnlock()

	if !exists {
		return
	}

	state.mu.Lock()
	delete(state.subscriptions, id)
	state.mu.Unlock()

	// Send confirmation
	response := map[string]any{
		"type": "unsubscribed",
		"id":   id,
	}
	if marshalResult := core.JSONMarshal(response); marshalResult.OK {
		responseData, _ := marshalResult.Value.([]byte)
		em.writeClientMessage(state, conn, responseData)
	}
}

// listSubscriptions sends a list of active subscriptions to a client.
func (em *WSEventManager) listSubscriptions(conn *websocket.Conn) {
	em.mu.RLock()
	state, exists := em.clients[conn]
	em.mu.RUnlock()

	if !exists {
		return
	}

	state.mu.RLock()
	subs := make([]*Subscription, 0, len(state.subscriptions))
	for _, sub := range state.subscriptions {
		subs = append(subs, sub)
	}
	state.mu.RUnlock()

	response := map[string]any{
		"type":          "subscriptions",
		"subscriptions": subs,
	}
	if marshalResult := core.JSONMarshal(response); marshalResult.OK {
		responseData, _ := marshalResult.Value.([]byte)
		em.writeClientMessage(state, conn, responseData)
	}
}

func (em *WSEventManager) writeClientMessage(state *clientState, conn *websocket.Conn, data []byte) {
	state.writeMu.Lock()
	defer state.writeMu.Unlock()
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_ = conn.WriteMessage(websocket.TextMessage, data)
}

// removeClient removes a client and its subscriptions.
func (em *WSEventManager) removeClient(conn *websocket.Conn) {
	em.mu.Lock()
	delete(em.clients, conn)
	em.mu.Unlock()
	conn.Close()
}

// Emit sends an event to all subscribed clients.
func (em *WSEventManager) Emit(event Event) {
	if em == nil {
		return
	}
	event.Timestamp = time.Now().UnixMilli()
	em.mu.RLock()
	if em.closed || em.eventBuffer == nil {
		em.mu.RUnlock()
		return
	}
	select {
	case em.eventBuffer <- event:
	default:
		// Buffer full, drop event
	}
	em.mu.RUnlock()
}

// EmitWindowEvent is a helper to emit window-related events.
func (em *WSEventManager) EmitWindowEvent(eventType EventType, windowName string, data map[string]any) {
	em.Emit(Event{
		Type:   eventType,
		Window: windowName,
		Data:   data,
	})
}

// ConnectedClients returns the number of connected WebSocket clients.
func (em *WSEventManager) ConnectedClients() int {
	if em == nil {
		return 0
	}
	em.mu.RLock()
	defer em.mu.RUnlock()
	return len(em.clients)
}

// Info returns a snapshot of the live WebSocket event server.
//
//	info := display.GetEventManager().Info()
func (em *WSEventManager) Info() events.ServerInfo {
	if em == nil {
		return events.ServerInfo{}
	}
	em.mu.RLock()
	defer em.mu.RUnlock()

	subscriptionCount := 0
	for _, state := range em.clients {
		if state == nil {
			continue
		}
		state.mu.RLock()
		subscriptionCount += len(state.subscriptions)
		state.mu.RUnlock()
	}

	return events.ServerInfo{
		ConnectedClients:  len(em.clients),
		SubscriptionCount: subscriptionCount,
		BufferLength:      len(em.eventBuffer),
		BufferCapacity:    cap(em.eventBuffer),
	}
}

// Close shuts down the event manager.
func (em *WSEventManager) Close() {
	if em == nil {
		return
	}
	em.mu.Lock()
	if em.closed {
		em.mu.Unlock()
		return
	}
	em.closed = true
	conns := make([]*websocket.Conn, 0, len(em.clients))
	for conn := range em.clients {
		conns = append(conns, conn)
	}
	em.clients = make(map[*websocket.Conn]*clientState)
	if em.eventBuffer != nil {
		close(em.eventBuffer)
		em.eventBuffer = nil
	}
	em.mu.Unlock()

	for _, conn := range conns {
		_ = conn.Close()
	}
}

type windowEventSource interface {
	OnWindowEvent(func(event window.WindowEvent))
}

// AttachWindowListeners attaches event listeners to a specific window.
// Use: em.AttachWindowListeners(windowHandle)
func (em *WSEventManager) AttachWindowListeners(pw windowEventSource) {
	if pw == nil {
		return
	}

	pw.OnWindowEvent(func(e window.WindowEvent) {
		em.EmitWindowEvent(EventType("window."+e.Type), e.Name, e.Data)
	})
}
