package chat

import (
	"bufio"
	"io"
	"slices"
	"strings"
	"time"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/internal/textutil"
	coreerr "dappco.re/go/log"
)

type StreamCallbacks struct {
	OnStart          func(streamID string)
	OnToken          func(content string)
	OnThinkingStart  func(state ThinkingState)
	OnThinkingAppend func(content string)
	OnThinkingEnd    func(state ThinkingState)
	OnToolCall       func(call ToolCall)
	OnFinish         func(reason string)
}

type StreamRenderer struct {
	callbacks      StreamCallbacks
	now            func() time.Time
	content        strings.Builder
	thinking       strings.Builder
	thinkingState  ThinkingState
	toolCalls      map[int]*streamToolCall
	parsedToolArgs map[int]map[string]any
	toolOrder      []int
	streamID       string
	started        bool
	finishReason   string
}

type streamToolCall struct {
	ID        string
	Name      string
	Arguments strings.Builder
}

type streamChunk struct {
	ID      string `json:"id"`
	Choices []struct {
		Delta struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			Reasoning string `json:"reasoning"`
			Thinking  string `json:"thinking"`
			Thought   string `json:"thought"`
			ToolCalls []struct {
				Index    int    `json:"index"`
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func NewStreamRenderer(callbacks StreamCallbacks) *StreamRenderer {
	return &StreamRenderer{
		callbacks:      callbacks,
		now:            time.Now,
		toolCalls:      make(map[int]*streamToolCall),
		parsedToolArgs: make(map[int]map[string]any),
	}
}

func (r *StreamRenderer) Render(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 4096), 1024*1024)

	var dataLines []string
	flush := func() error {
		if len(dataLines) == 0 {
			return nil
		}
		payload := strings.Join(dataLines, "\n")
		dataLines = nil
		return r.handleData(payload)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if err := flush(); err != nil {
				return err
			}
			continue
		}
		if strings.HasPrefix(line, "data:") {
			dataLines = append(dataLines, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := flush(); err != nil {
		return err
	}

	r.completeThinking()
	for _, call := range r.ToolCalls() {
		if r.callbacks.OnToolCall != nil {
			r.callbacks.OnToolCall(call)
		}
	}
	if r.callbacks.OnFinish != nil {
		r.callbacks.OnFinish(r.finishReason)
	}
	return nil
}

func (r *StreamRenderer) handleData(payload string) error {
	if payload == "" {
		return nil
	}
	if payload == "[DONE]" {
		return nil
	}

	var chunk streamChunk
	result := core.JSONUnmarshalString(payload, &chunk)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return err
		}
		return coreerr.E("chat.StreamRenderer.handleData", "failed to decode stream chunk", nil)
	}
	if !r.started {
		r.started = true
		r.streamID = chunk.ID
		if r.callbacks.OnStart != nil {
			r.callbacks.OnStart(chunk.ID)
		}
	}

	for _, choice := range chunk.Choices {
		delta := choice.Delta
		thought := textutil.FirstNonEmpty(delta.Thinking, delta.Reasoning, delta.Thought)
		if thought != "" {
			r.appendThinking(thought)
		}
		if delta.Content != "" {
			r.completeThinking()
			r.content.WriteString(delta.Content)
			if r.callbacks.OnToken != nil {
				r.callbacks.OnToken(delta.Content)
			}
		}
		for _, call := range delta.ToolCalls {
			r.appendToolCall(call.Index, call.ID, call.Function.Name, call.Function.Arguments)
		}
		if choice.FinishReason != "" {
			r.finishReason = choice.FinishReason
		}
	}
	return nil
}

func (r *StreamRenderer) appendThinking(content string) {
	if content == "" {
		return
	}
	if !r.thinkingState.Active {
		r.thinkingState = ThinkingState{
			Active:    true,
			StartedAt: r.now(),
		}
		if r.callbacks.OnThinkingStart != nil {
			r.callbacks.OnThinkingStart(r.thinkingState)
		}
	}
	r.thinking.WriteString(content)
	if r.callbacks.OnThinkingAppend != nil {
		r.callbacks.OnThinkingAppend(content)
	}
}

func (r *StreamRenderer) completeThinking() {
	if !r.thinkingState.Active {
		return
	}
	r.thinkingState.Active = false
	r.thinkingState.Content = r.thinking.String()
	r.thinkingState.EndedAt = r.now()
	r.thinkingState.DurationMS = r.thinkingState.EndedAt.Sub(r.thinkingState.StartedAt).Milliseconds()
	if r.callbacks.OnThinkingEnd != nil {
		r.callbacks.OnThinkingEnd(r.thinkingState)
	}
}

func (r *StreamRenderer) appendToolCall(index int, id, name, arguments string) {
	call := r.toolCalls[index]
	if call == nil {
		call = &streamToolCall{}
		r.toolCalls[index] = call
		r.toolOrder = append(r.toolOrder, index)
	}
	if id != "" {
		call.ID = id
	}
	if name != "" {
		call.Name = name
	}
	if arguments != "" {
		call.Arguments.WriteString(arguments)
		delete(r.parsedToolArgs, index)
	}
}

func (r *StreamRenderer) ToolCalls() []ToolCall {
	order := slices.Clone(r.toolOrder)
	slices.Sort(order)
	result := make([]ToolCall, 0, len(order))
	for _, index := range order {
		call := r.toolCalls[index]
		if call == nil {
			continue
		}
		arguments, ok := r.parsedToolArgs[index]
		if !ok {
			arguments = map[string]any{}
			raw := strings.TrimSpace(call.Arguments.String())
			if raw != "" {
				if decode := core.JSONUnmarshalString(raw, &arguments); !decode.OK {
					arguments = map[string]any{"raw": raw}
				}
			}
			r.parsedToolArgs[index] = cloneArguments(arguments)
		}
		result = append(result, ToolCall{
			ID:        call.ID,
			Name:      call.Name,
			Arguments: cloneArguments(arguments),
		})
	}
	return result
}

func (r *StreamRenderer) Thinking() *ThinkingState {
	if r.thinking.Len() == 0 && !r.thinkingState.Active {
		return nil
	}
	state := r.thinkingState
	state.Content = r.thinking.String()
	return &state
}

func (r *StreamRenderer) Message(messageID, model string, createdAt time.Time) ChatMessage {
	return ChatMessage{
		ID:           messageID,
		Role:         "assistant",
		Content:      r.content.String(),
		CreatedAt:    createdAt,
		Model:        model,
		Thinking:     r.Thinking(),
		ToolCalls:    r.ToolCalls(),
		FinishReason: r.finishReason,
	}
}
