package chat

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamRenderer_Good_ParsesThinkingContentAndToolCalls(t *testing.T) {
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
	require.NoError(t, renderer.Render(strings.NewReader(stream)))

	message := renderer.Message("msg-1", "lemer", testTime())
	require.NotNil(t, message.Thinking)
	assert.Equal(t, "Hello", message.Content)
	assert.Equal(t, "Let me think", message.Thinking.Content)
	require.Len(t, message.ToolCalls, 1)
	assert.Equal(t, "layout_suggest", message.ToolCalls[0].Name)
	assert.Equal(t, 2.0, message.ToolCalls[0].Arguments["window_count"])
	assert.Equal(t, "tool_calls", message.FinishReason)
}

func testTime() time.Time {
	return time.Unix(1_700_000_000, 0).UTC()
}
