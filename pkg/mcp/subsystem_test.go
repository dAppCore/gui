package mcp

import (
	"context"
	"reflect"
	"testing"

	core "dappco.re/go/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubsystem_renderCallToolResult_Good(t *testing.T) {
	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "alpha"},
			&mcp.ImageContent{Data: []byte("png"), MIMEType: "image/png"},
		},
	}

	rendered := renderCallToolResult(result)

	assert.Contains(t, rendered, "alpha")
	assert.Contains(t, rendered, "\"mimeType\":\"image/png\"")
	assert.Contains(t, rendered, "\n")
}

func TestSubsystem_renderCallToolResult_Bad(t *testing.T) {
	rendered := renderCallToolResult(&mcp.CallToolResult{})

	assert.Contains(t, rendered, "\"content\":null")
}

func TestSubsystem_renderCallToolResult_Ugly(t *testing.T) {
	assert.Equal(t, "", renderCallToolResult(nil))
}

func TestSubsystem_normalizeSchema_Good(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}

	assert.Equal(t, schema, normalizeSchema(schema))
}

func TestSubsystem_normalizeSchema_Bad(t *testing.T) {
	assert.Nil(t, normalizeSchema(nil))
}

func TestSubsystem_normalizeSchema_Ugly(t *testing.T) {
	type payload struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	assert.Equal(t, map[string]any{"name": "core", "count": float64(2)}, normalizeSchema(payload{Name: "core", Count: 2}))
}

func TestSubsystem_schemaForType_Good(t *testing.T) {
	type sample struct {
		Name    string `json:"name,omitempty"`
		Alias   string `json:",omitempty"`
		Count   int
		skip    string
		Ignored string `json:"-"`
	}

	schema := schemaForType(reflect.TypeOf(sample{}))

	assert.Equal(t, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":  map[string]any{"type": "string"},
			"Alias": map[string]any{"type": "string"},
			"Count": map[string]any{"type": "integer"},
		},
		"required": []string{"Count"},
	}, schema)
}

func TestSubsystem_schemaForType_Bad(t *testing.T) {
	assert.Nil(t, schemaForType(nil))
}

func TestSubsystem_schemaForType_Ugly(t *testing.T) {
	assert.Equal(t, map[string]any{"type": "string"}, schemaForType(reflect.TypeOf(make(chan int))))
}

func TestSubsystem_CallTool_Bad_UnknownTool(t *testing.T) {
	sub := New(core.New(core.WithServiceLock()))

	_, err := sub.CallTool(context.Background(), "missing_tool", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tool not found")
}

func TestSubsystem_CallTool_Ugly_InvalidArguments(t *testing.T) {
	sub := New(core.New(core.WithServiceLock()))
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	_, err := sub.CallTool(context.Background(), "layout_suggest", map[string]any{
		"window_count": map[string]any{"unexpected": true},
	})
	require.Error(t, err)
}
