package chat

import (
	core "dappco.re/go"
	"strings"
	"time"
)

func TestStreamRenderer_Good_ParsesThinkingContentAndToolCalls(t *core.T) {
	stream := strings.Join([]string{
		`data: {"id":"chatcmpl-1","choices":[{"delta":{"thinking":"Let me think"}}]}`,
		"",
		`data: {"id":"chatcmpl-1","choices":[{"delta":{"content":"Hello"}}]}`,
		"",
		`data: {"id":"chatcmpl-1","choices":[{"delta":{"tool_calls":[{"index":0,"id":"call-1","function":{"name":"layout_suggest","arguments":"{\"window_count\":2}"}}]}}]}`,
		"",
		`data: {"id":"chatcmpl-1","choices":[{"finish_reason":"tool_calls"}]}`,
		"",
		`data: [DONE]`,
		"",
	}, "\n")

	renderer := NewStreamRenderer(StreamCallbacks{})
	core.RequireNoError(t, renderer.Render(strings.NewReader(stream)))

	message := renderer.Message("msg-1", "lemer", testTime())
	core.AssertNotNil(t, message.Thinking)
	core.AssertEqual(t, "Hello", message.Content)
	core.AssertEqual(t, "Let me think", message.Thinking.Content)
	core.AssertLen(t, message.ToolCalls, 1)
	core.AssertEqual(t, "layout_suggest", message.ToolCalls[0].Name)
	core.AssertEqual(t, 2.0, message.ToolCalls[0].Arguments["window_count"])
	core.AssertEqual(t, "tool_calls", message.FinishReason)
}

func testTime() time.Time {
	return time.Unix(1_700_000_000, 0).UTC()
}
