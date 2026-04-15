package display

import (
	"context"
	"net/url"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type testApplicationResponseWriter struct {
	header map[string][]string
	body   []byte
	status int
}

func newTestApplicationResponseWriter() *testApplicationResponseWriter {
	return &testApplicationResponseWriter{header: map[string][]string{}}
}

func (w *testApplicationResponseWriter) Header() map[string][]string {
	if w.header == nil {
		w.header = map[string][]string{}
	}
	return w.header
}

func (w *testApplicationResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return len(data), nil
}

func (w *testApplicationResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

type testApplicationHandler struct {
	called bool
}

func (h *testApplicationHandler) ServeHTTP(_ application.ResponseWriter, _ *application.Request) {
	h.called = true
}

func TestScheme_ResolveScheme_Good(t *testing.T) {
	svc, c := newTestDisplayService(t)
	svc.registerDefaultSchemes()
	svc.configFile = nil
	svc.storage.Set("origin-a", "localStorage", "theme", "dark")
	svc.configData["window"] = map[string]any{"default_width": 1024, "default_height": 768}

	c.Action("gui.chat.models", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{
			Value: []chat.ModelEntry{
				{Name: "Alpha", Architecture: "gemma", SizeBytes: 2048, Loaded: true, Backend: "local"},
				{Name: "Beta", Architecture: "phi", SizeBytes: 4096, Loaded: false},
			},
			OK: true,
		}
	})
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case chat.QueryConversationList:
			return core.Result{
				Value: []chat.ConversationSummary{{ID: "conv-1", Title: "Chat Route", MessageCount: 3}},
				OK:    true,
			}
		case chat.QueryHistory:
			return core.Result{
				Value: chat.Conversation{ID: "conv-1", Title: "Chat Route"},
				OK:    true,
			}
		case chat.QueryConversationSearch:
			return core.Result{
				Value: []any{"chat-match"},
				OK:    true,
			}
		default:
			return core.Result{}
		}
	})

	storeResult := svc.ResolveScheme(context.Background(), "core://store?q=theme")
	require.True(t, storeResult.OK)
	storePayload := storeResult.Value.(map[string]any)
	assert.Equal(t, "text/html", storePayload["content_type"])
	assert.Contains(t, storePayload["body"].(string), "origin-a")
	assert.Contains(t, storePayload["body"].(string), "dark")

	entryResult := svc.ResolveScheme(context.Background(), "core://store/localStorage/theme")
	require.True(t, entryResult.OK)
	entryPayload := entryResult.Value.(map[string]any)
	assert.Equal(t, "store", entryPayload["route"])
	assert.Contains(t, entryPayload["body"].(string), "localStorage")
	assert.Contains(t, entryPayload["body"].(string), "theme")

	settingsResult := svc.ResolveScheme(context.Background(), "core://settings/window")
	require.True(t, settingsResult.OK)
	settingsPayload := settingsResult.Value.(map[string]any)
	assert.Equal(t, "settings", settingsPayload["route"])
	assert.Contains(t, settingsPayload["body"].(string), "default_width")
	assert.Contains(t, settingsPayload["body"].(string), "1024")

	modelResult := svc.ResolveScheme(context.Background(), "core://models/alpha")
	require.True(t, modelResult.OK)
	modelPayload := modelResult.Value.(map[string]any)
	assert.Equal(t, "models", modelPayload["route"])
	assert.Contains(t, modelPayload["body"].(string), "Alpha")
	assert.Contains(t, modelPayload["body"].(string), "2048")

	chatListResult := svc.ResolveScheme(context.Background(), "core://chat")
	require.True(t, chatListResult.OK)
	chatListPayload := chatListResult.Value.(map[string]any)
	assert.Equal(t, "chat", chatListPayload["route"])
	assert.Contains(t, chatListPayload["body"].(string), "Chat Route")

	chatHistoryResult := svc.ResolveScheme(context.Background(), "core://chat?conversation_id=conv-1")
	require.True(t, chatHistoryResult.OK)
	chatHistoryPayload := chatHistoryResult.Value.(map[string]any)
	assert.Equal(t, "chat", chatHistoryPayload["route"])
	assert.Contains(t, chatHistoryPayload["body"].(string), "conv-1")

	chatSearchResult := svc.handleStoreSearch(context.Background(), url.Values{"q": []string{"chat"}})
	require.True(t, chatSearchResult.OK)
	chatSearchPayload := chatSearchResult.Value.(map[string]any)
	assert.Contains(t, chatSearchPayload["body"].(string), "core://chat")
}

func TestScheme_ResolveScheme_Bad(t *testing.T) {
	svc := &Service{}

	emptyResult := svc.ResolveScheme(context.Background(), "")
	require.False(t, emptyResult.OK)

	malformedResult := svc.ResolveScheme(context.Background(), "://bad-url")
	require.False(t, malformedResult.OK)

	noHandlerResult := svc.ResolveScheme(context.Background(), "core://store")
	require.False(t, noHandlerResult.OK)
}

func TestScheme_ResolveScheme_Ugly(t *testing.T) {
	svc, _ := newTestDisplayService(t)
	svc.registerDefaultSchemes()

	result := svc.ResolveScheme(context.Background(), "core://wallet/treasury?amount=1")
	require.True(t, result.OK)
	payload := result.Value.(map[string]any)
	assert.Equal(t, "wallet", payload["route"])
	assert.Equal(t, false, payload["available"])
	assert.Contains(t, payload["body"].(string), "no backend is registered for this route")

	searchResult := svc.handleStoreSearch(context.Background(), url.Values{"q": []string{"missing"}})
	require.True(t, searchResult.OK)
	searchPayload := searchResult.Value.(map[string]any)
	assert.Contains(t, searchPayload["body"].(string), "No matches found in Core storage.")
}

func TestScheme_AssetMiddleware_Good(t *testing.T) {
	svc, _ := newTestDisplayService(t)
	svc.registerDefaultSchemes()

	recorder := newTestApplicationResponseWriter()
	request := &application.Request{Method: "GET", URL: "core://store?q=theme"}

	svc.AssetMiddleware()(&testApplicationHandler{}).ServeHTTP(recorder, request)

	require.Equal(t, 200, recorder.status)
	assert.Equal(t, "text/html; charset=utf-8", recorder.Header()["Content-Type"][0])
	assert.Contains(t, string(recorder.body), "core://store")
}

func TestScheme_AssetMiddleware_Bad(t *testing.T) {
	svc, _ := newTestDisplayService(t)
	svc.registerDefaultSchemes()

	recorder := newTestApplicationResponseWriter()
	request := &application.Request{Method: "GET", URL: "https://example.com/app"}
	next := &testApplicationHandler{}

	svc.AssetMiddleware()(next).ServeHTTP(recorder, request)

	require.True(t, next.called)
	require.Equal(t, 0, recorder.status)
	assert.Empty(t, recorder.body)
}

func TestScheme_AssetMiddleware_Ugly(t *testing.T) {
	svc, _ := New()

	recorder := newTestApplicationResponseWriter()
	request := &application.Request{Method: "GET", URL: "core://missing"}

	svc.AssetMiddleware()(&testApplicationHandler{}).ServeHTTP(recorder, request)

	require.Equal(t, 404, recorder.status)
	assert.Contains(t, string(recorder.body), "core route not found")
}
