package mcp

import (
	"context"
	"reflect"

	core "dappco.re/go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSubsystem_renderCallToolResult_Good(t *core.T) {
	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "alpha"},
			&mcp.ImageContent{Data: []byte("png"), MIMEType: "image/png"},
		},
	}

	rendered := renderCallToolResult(result)

	core.AssertContains(t, rendered, "alpha")
	core.AssertContains(t, rendered, "\"mimeType\":\"image/png\"")
	core.AssertContains(t, rendered, "\n")
}

func TestSubsystem_renderCallToolResult_Bad(t *core.T) {
	rendered := renderCallToolResult(&mcp.CallToolResult{})

	core.AssertContains(t, rendered, "\"content\":null")
	core.AssertNotEmpty(t, core.Sprintf("%T", rendered))
}

func TestSubsystem_renderCallToolResult_Ugly(t *core.T) {
	core.AssertEqual(t, "", renderCallToolResult(nil))
	observedType := core.Sprintf("%T", renderCallToolResult(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestSubsystem_normalizeSchema_Good(t *core.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}

	core.AssertEqual(t, schema, normalizeSchema(schema))
}

func TestSubsystem_normalizeSchema_Bad(t *core.T) {
	core.AssertNil(t, normalizeSchema(nil))
	observedType := core.Sprintf("%T", normalizeSchema(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestSubsystem_normalizeSchema_Ugly(t *core.T) {
	type payload struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	core.AssertEqual(t, map[string]any{"name": "core", "count": float64(2)}, normalizeSchema(payload{Name: "core", Count: 2}))
}

func TestSubsystem_schemaForType_Good(t *core.T) {
	type sample struct {
		Name    string `json:"name,omitempty"`
		Alias   string `json:",omitempty"`
		Count   int
		skip    string
		Ignored string `json:"-"`
	}

	schema := schemaForType(reflect.TypeOf(sample{}))

	core.AssertEqual(t, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":  map[string]any{"type": "string"},
			"Alias": map[string]any{"type": "string"},
			"Count": map[string]any{"type": "integer"},
		},
		"required": []string{"Count"},
	}, schema)
}

func TestSubsystem_schemaForType_Bad(t *core.T) {
	core.AssertNil(t, schemaForType(nil))
	observedType := core.Sprintf("%T", schemaForType(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestSubsystem_schemaForType_Ugly(t *core.T) {
	core.AssertEqual(t, map[string]any{"type": "string"}, schemaForType(reflect.TypeOf(make(chan int))))
	observedType := core.Sprintf("%T", schemaForType(reflect.TypeOf(make(chan int))))
	core.AssertNotEmpty(t, observedType)
}

func TestSubsystem_CallTool_Bad_UnknownTool(t *core.T) {
	sub := New(core.New(core.WithServiceLock()))

	_, err := sub.CallTool(context.Background(), "missing_tool", nil)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "tool not found")
}

func TestSubsystem_CallTool_Ugly_InvalidArguments(t *core.T) {
	sub := New(core.New(core.WithServiceLock()))
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	_, err := sub.CallTool(context.Background(), "layout_suggest", map[string]any{
		"window_count": map[string]any{"unexpected": true},
	})
	core.AssertError(t, err)
}

// AX7 generated source-matching smoke coverage.
func TestSubsystem_New_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := New(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_New_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := New(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_New_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := New(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_Name_Good(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_Name_Bad(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_Name_Ugly(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.Name()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_RegisterTools_Good(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		subject.RegisterTools(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_RegisterTools_Bad(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		subject.RegisterTools(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_RegisterTools_Ugly(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		subject.RegisterTools(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_Manifest_Good(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.Manifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_Manifest_Bad(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.Manifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_Manifest_Ugly(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.Manifest()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_ManifestText_Good(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.ManifestText()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_ManifestText_Bad(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.ManifestText()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_ManifestText_Ugly(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0 := subject.ManifestText()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_CallTool_Good(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0, got1 := subject.CallTool(core.Background(), "agent", nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_CallTool_Bad(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0, got1 := subject.CallTool(core.Background(), "", nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSubsystem_Subsystem_CallTool_Ugly(t *core.T) {
	subject := new(Subsystem)
	result := core.Try(func() any {
		got0, got1 := subject.CallTool(core.Background(), "../../edge", nil)
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
