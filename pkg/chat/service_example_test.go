package chat

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	core "dappco.re/go"
)

func ExampleRegister() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Hello from chat"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}))
	defer server.Close()

	storeDir, err := os.MkdirTemp("", "chat-example-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(storeDir)

	c := core.New(
		core.WithService(Register(
			func(o *Options) { o.APIURL = server.URL },
			func(o *Options) { o.StorePath = filepath.Join(storeDir, "chat.db") },
			func(o *Options) { o.ToolExecutor = &mockToolExecutor{} },
			func(o *Options) { o.Now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() } },
		)),
		core.WithServiceLock(),
	)
	if !c.ServiceStartup(context.Background(), nil).OK {
		panic("chat startup failed")
	}

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hello"},
	))
	if !send.OK {
		panic(send.Value)
	}

	conversations := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	history := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conversations.Value.([]Conversation)[0].ID},
	))

	fmt.Println(len(history.Value.([]Message)))
	fmt.Println(history.Value.([]Message)[1].Content)
	// Output:
	// 2
	// Hello from chat
}
