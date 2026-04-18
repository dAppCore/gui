package display

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"forge.lthn.ai/core/gui/pkg/events"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWindowListener struct {
	handler func(window.WindowEvent)
}

func (w *testWindowListener) Name() string                    { return "listener" }
func (w *testWindowListener) Title() string                   { return "listener" }
func (w *testWindowListener) Position() (int, int)            { return 0, 0 }
func (w *testWindowListener) Size() (int, int)                { return 0, 0 }
func (w *testWindowListener) IsMaximised() bool               { return false }
func (w *testWindowListener) IsFocused() bool                 { return false }
func (w *testWindowListener) IsVisible() bool                 { return false }
func (w *testWindowListener) IsFullscreen() bool              { return false }
func (w *testWindowListener) IsMinimised() bool               { return false }
func (w *testWindowListener) GetBounds() (int, int, int, int) { return 0, 0, 0, 0 }
func (w *testWindowListener) GetZoom() float64                { return 1 }
func (w *testWindowListener) SetTitle(string)                 {}
func (w *testWindowListener) SetPosition(int, int)            {}
func (w *testWindowListener) SetSize(int, int)                {}
func (w *testWindowListener) SetBackgroundColour(uint8, uint8, uint8, uint8) {
}
func (w *testWindowListener) SetVisibility(bool)           {}
func (w *testWindowListener) SetAlwaysOnTop(bool)          {}
func (w *testWindowListener) SetBounds(int, int, int, int) {}
func (w *testWindowListener) SetURL(string)                {}
func (w *testWindowListener) SetHTML(string)               {}
func (w *testWindowListener) SetZoom(float64)              {}
func (w *testWindowListener) SetContentProtection(bool)    {}
func (w *testWindowListener) Maximise()                    {}
func (w *testWindowListener) Restore()                     {}
func (w *testWindowListener) Minimise()                    {}
func (w *testWindowListener) Focus()                       {}
func (w *testWindowListener) Close()                       {}
func (w *testWindowListener) Show()                        {}
func (w *testWindowListener) Hide()                        {}
func (w *testWindowListener) Fullscreen()                  {}
func (w *testWindowListener) UnFullscreen()                {}
func (w *testWindowListener) ToggleFullscreen()            {}
func (w *testWindowListener) ToggleMaximise()              {}
func (w *testWindowListener) ExecJS(string)                {}
func (w *testWindowListener) Flash(bool)                   {}
func (w *testWindowListener) Print() error                 { return nil }
func (w *testWindowListener) OnWindowEvent(handler func(window.WindowEvent)) {
	w.handler = handler
}
func (w *testWindowListener) OnFileDrop(func([]string, string)) {}

func (w *testWindowListener) emit(event window.WindowEvent) {
	if w.handler != nil {
		w.handler(event)
	}
}

func dialWSEventManager(t *testing.T, em *WSEventManager) (*websocket.Conn, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(em.HandleWebSocket))
	t.Cleanup(server.Close)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	return conn, func() {
		_ = conn.Close()
		em.Close()
	}
}

func readJSONMessage(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()

	_, data, err := conn.ReadMessage()
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))
	return payload
}

func writeJSONMessage(t *testing.T, conn *websocket.Conn, payload map[string]any) {
	t.Helper()
	require.NoError(t, conn.WriteJSON(payload))
}

func TestWSEventManager_HandleWebSocket_Good(t *testing.T) {
	em := NewWSEventManager()
	conn, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	writeJSONMessage(t, conn, map[string]any{
		"action":     "subscribe",
		"id":         "sub-1",
		"eventTypes": []string{"*"},
	})
	subscribeAck := readJSONMessage(t, conn)
	assert.Equal(t, "subscribed", subscribeAck["type"])
	assert.Equal(t, "sub-1", subscribeAck["id"])

	info := em.Info()
	assert.Equal(t, 1, info.ConnectedClients)
	assert.Equal(t, 1, info.SubscriptionCount)

	em.Emit(Event{Type: EventCustomEvent, Window: "main", Data: map[string]any{"hello": "world"}})
	eventPayload := readJSONMessage(t, conn)
	assert.Equal(t, string(EventCustomEvent), eventPayload["type"])
	assert.Equal(t, "main", eventPayload["window"])

	writeJSONMessage(t, conn, map[string]any{"action": "list"})
	listPayload := readJSONMessage(t, conn)
	assert.Equal(t, "subscriptions", listPayload["type"])
	assert.Len(t, listPayload["subscriptions"].([]any), 1)

	writeJSONMessage(t, conn, map[string]any{"action": "unsubscribe", "id": "sub-1"})
	unsubscribeAck := readJSONMessage(t, conn)
	assert.Equal(t, "unsubscribed", unsubscribeAck["type"])

	writeJSONMessage(t, conn, map[string]any{"action": "list"})
	emptyList := readJSONMessage(t, conn)
	assert.Equal(t, "subscriptions", emptyList["type"])
	assert.Len(t, emptyList["subscriptions"].([]any), 0)

	require.NoError(t, conn.Close())
	require.Eventually(t, func() bool {
		return em.ConnectedClients() == 0
	}, 2*time.Second, 20*time.Millisecond)
}

func TestWSEventManager_HandleWebSocket_RejectsRemoteOrigin(t *testing.T) {
	em := NewWSEventManager()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.Header.Set("Origin", "https://evil.example")
	recorder := httptest.NewRecorder()

	em.HandleWebSocket(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestWSEventManager_HandleWebSocket_RejectsLoopbackSpoofedOrigin(t *testing.T) {
	em := NewWSEventManager()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Origin", "file://malicious")
	recorder := httptest.NewRecorder()

	em.HandleWebSocket(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestWSEventManager_HandleWebSocket_NilReceiverFailsClosed(t *testing.T) {
	var em *WSEventManager

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	recorder := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		em.HandleWebSocket(recorder, req)
	})
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestWSEventManager_HandleWebSocket_NilWriterFailsClosed(t *testing.T) {
	em := NewWSEventManager()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	assert.NotPanics(t, func() {
		em.HandleWebSocket(nil, req)
	})
}

func TestWSEventManager_HandleWebSocket_RejectsAfterClose(t *testing.T) {
	em := NewWSEventManager()
	em.Close()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	recorder := httptest.NewRecorder()

	em.HandleWebSocket(recorder, req)

	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestEvents_trustedWebSocketOrigin_Good(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{
			name: "localhost without origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				return r
			}(),
			want: true,
		},
		{
			name: "local origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				r.Header.Set("Origin", "http://localhost:8080")
				return r
			}(),
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, trustedWebSocketOrigin(tc.req))
		})
	}
}

func TestEvents_trustedWebSocketOrigin_Bad(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
	}{
		{
			name: "nil request",
		},
		{
			name: "wrong path",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/other", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				return r
			}(),
		},
		{
			name: "remote client",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "203.0.113.10:2222"
				return r
			}(),
		},
		{
			name: "remote origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				r.Header.Set("Origin", "https://evil.example")
				return r
			}(),
		},
		{
			name: "file origin",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
				r.RemoteAddr = "127.0.0.1:12345"
				r.Header.Set("Origin", "file://local")
				return r
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.False(t, trustedWebSocketOrigin(tc.req))
		})
	}
}

func TestEvents_trustedWebSocketOrigin_Ugly(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/events", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Origin", "://bad")

	assert.False(t, trustedWebSocketOrigin(req))
	assert.False(t, trustedWebSocketOrigin(&http.Request{}))
}

func TestEvents_trustedWebSocketHost_Good(t *testing.T) {
	assert.True(t, trustedWebSocketHost("localhost"))
	assert.True(t, trustedWebSocketHost("127.0.0.1:443"))
	assert.True(t, trustedWebSocketHost("[::1]:80"))
}

func TestEvents_trustedWebSocketHost_Bad(t *testing.T) {
	assert.False(t, trustedWebSocketHost(""))
	assert.False(t, trustedWebSocketHost("example.com"))
}

func TestEvents_trustedWebSocketHost_Ugly(t *testing.T) {
	assert.False(t, trustedWebSocketHost("not a host"))
}

func TestEvents_isLoopbackHost_Good(t *testing.T) {
	assert.True(t, isLoopbackHost("localhost"))
	assert.True(t, isLoopbackHost("127.0.0.1"))
	assert.True(t, isLoopbackHost("::1"))
}

func TestEvents_isLoopbackHost_Bad(t *testing.T) {
	assert.False(t, isLoopbackHost(""))
	assert.False(t, isLoopbackHost("example.com"))
}

func TestEvents_isLoopbackHost_Ugly(t *testing.T) {
	assert.False(t, isLoopbackHost("203.0.113.10"))
}

func TestWSEventManager_HandleWebSocket_ClosesOnMalformedMessage(t *testing.T) {
	em := NewWSEventManager()
	conn, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"action":`)))

	payload := readJSONMessage(t, conn)
	assert.Equal(t, "invalid websocket message", payload["error"])
	assert.Equal(t, float64(websocket.ClosePolicyViolation), payload["status"])

	_, _, err := conn.ReadMessage()
	require.Error(t, err)
}

func TestWSEventManager_HandleWebSocket_ClosesOnUnknownAction(t *testing.T) {
	em := NewWSEventManager()
	conn, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"action":"bogus"}`)))

	payload := readJSONMessage(t, conn)
	assert.Equal(t, "unknown websocket action", payload["error"])
	assert.Equal(t, float64(websocket.ClosePolicyViolation), payload["status"])

	_, _, err := conn.ReadMessage()
	require.Error(t, err)
}

func TestWSEventManager_HandleWebSocket_ClosesOnReadTimeout(t *testing.T) {
	em := NewWSEventManager()
	em.readTimeout = 10 * time.Millisecond

	_, cleanup := dialWSEventManager(t, em)
	defer cleanup()

	require.Eventually(t, func() bool {
		return em.ConnectedClients() == 0
	}, 500*time.Millisecond, 10*time.Millisecond)
}

func TestWSEventManager_Emit_Ugly(t *testing.T) {
	em := &WSEventManager{
		clients:     map[*websocket.Conn]*clientState{},
		eventBuffer: make(chan Event, 1),
	}

	em.Emit(Event{Type: EventTrayClick})
	em.Emit(Event{Type: EventTrayClick})

	assert.Len(t, em.eventBuffer, 1)
	info := em.Info()
	assert.Equal(t, 0, info.ConnectedClients)
	assert.Equal(t, 1, info.BufferLength)
}

func TestWSEventManager_EmitWindowEvent_Good(t *testing.T) {
	em := &WSEventManager{
		clients:     map[*websocket.Conn]*clientState{},
		eventBuffer: make(chan Event, 2),
	}

	em.EmitWindowEvent(EventWindowMove, "editor", map[string]any{"x": 10, "y": 20})
	require.Len(t, em.eventBuffer, 1)

	event := <-em.eventBuffer
	assert.Equal(t, EventWindowMove, event.Type)
	assert.Equal(t, "editor", event.Window)
	assert.Equal(t, 10, event.Data["x"])
}

func TestWSEventManager_ClientSubscribed_Good(t *testing.T) {
	em := &WSEventManager{}
	state := &clientState{
		subscriptions: map[string]*Subscription{
			"sub-1": {ID: "sub-1", EventTypes: []EventType{EventWindowFocus}},
			"sub-2": {ID: "sub-2", EventTypes: []EventType{"*"}},
		},
	}

	assert.True(t, em.clientSubscribed(state, EventWindowFocus))
	assert.True(t, em.clientSubscribed(state, EventCustomEvent))
}

func TestWSEventManager_ClientSubscribed_Bad(t *testing.T) {
	em := &WSEventManager{}
	state := &clientState{subscriptions: map[string]*Subscription{}}

	assert.False(t, em.clientSubscribed(state, EventWindowClose))
}

func TestWSEventManager_ConnectedClients_Good(t *testing.T) {
	em := &WSEventManager{
		clients: map[*websocket.Conn]*clientState{},
	}
	em.clients[nil] = &clientState{subscriptions: map[string]*Subscription{}}

	assert.Equal(t, 1, em.ConnectedClients())
}

func TestWSEventManager_NilSafety(t *testing.T) {
	var em *WSEventManager

	assert.Equal(t, 0, em.ConnectedClients())
	assert.Equal(t, events.ServerInfo{}, em.Info())
}

func TestWSEventManager_AttachWindowListeners_Good(t *testing.T) {
	em := &WSEventManager{
		clients:     map[*websocket.Conn]*clientState{},
		eventBuffer: make(chan Event, 1),
	}

	listener := &testWindowListener{}
	em.AttachWindowListeners(listener)
	listener.emit(window.WindowEvent{Type: "focus", Name: "main", Data: map[string]any{"reason": "user"}})

	require.Len(t, em.eventBuffer, 1)
	event := <-em.eventBuffer
	assert.Equal(t, EventType("window.focus"), event.Type)
	assert.Equal(t, "main", event.Window)
	assert.Equal(t, "user", event.Data["reason"])
}

func TestWSEventManager_AttachWindowListeners_Bad(t *testing.T) {
	em := &WSEventManager{}
	em.AttachWindowListeners(nil)
	assert.Empty(t, em.eventBuffer)
}

func TestWSEventManager_CloseIsIdempotent(t *testing.T) {
	em := NewWSEventManager()

	assert.NotPanics(t, func() {
		em.Close()
		em.Close()
		em.Emit(Event{Type: EventCustomEvent})
	})
}
