package display

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				OK: true,
			}
		case chat.QueryHistory:
			return core.Result{
				Value: chat.Conversation{ID: "conv-1", Title: "Chat Route"},
				OK: true,
			}
		case chat.QueryConversationSearch:
			return core.Result{
				Value: []any{"chat-match"},
				OK: true,
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

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "core://store?q=theme", nil)

	svc.AssetMiddleware()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not be called for core routes")
	})).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "text/html; charset=utf-8", recorder.Header().Get("Content-Type"))
	assert.Contains(t, recorder.Body.String(), "core://store")
}

func TestScheme_AssetMiddleware_Bad(t *testing.T) {
	svc, _ := newTestDisplayService(t)
	svc.registerDefaultSchemes()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://example.com/app", nil)
	nextCalled := false

	svc.AssetMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		_, _ = w.Write([]byte("next"))
	})).ServeHTTP(recorder, request)

	require.True(t, nextCalled)
	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "next", recorder.Body.String())
}

func TestScheme_AssetMiddleware_Ugly(t *testing.T) {
	svc, _ := New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "core://missing", nil)

	svc.AssetMiddleware()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not be called for core routes")
	})).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "core route not found")
}
