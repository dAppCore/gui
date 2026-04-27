package webview

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnostics_ComputedStyleScript_Good(t *testing.T) {
	script := ComputedStyleScript("#app")

	require.Contains(t, script, "document.querySelector")
	assert.Contains(t, script, "window.getComputedStyle(el)")
	assert.Contains(t, script, mustJSON("#app"))
}

func TestDiagnostics_ComputedStyleScript_Bad(t *testing.T) {
	script := ComputedStyleScript("")

	assert.Contains(t, script, `document.querySelector("")`)
	assert.Contains(t, script, "return null;")
}

func TestDiagnostics_ComputedStyleScript_Ugly(t *testing.T) {
	selector := "#app\"\n\t"
	script := ComputedStyleScript(selector)

	assert.Contains(t, script, mustJSON(selector))
}

func TestDiagnostics_HighlightScript_Good(t *testing.T) {
	script := HighlightScript(".card", "#00ff00")

	assert.Contains(t, script, mustJSON(".card"))
	assert.Contains(t, script, mustJSON("#00ff00"))
	assert.Contains(t, script, `outline = "3px solid " +`)
}

func TestDiagnostics_HighlightScript_Bad(t *testing.T) {
	script := HighlightScript(".card", "")

	assert.Contains(t, script, `3px solid `)
	assert.Contains(t, script, mustJSON("#ff9800"))
}

func TestDiagnostics_HighlightScript_Ugly(t *testing.T) {
	selector := `#card" + alert(1) + "`
	script := HighlightScript(selector, "#123456")

	assert.Contains(t, script, mustJSON(selector))
	assert.Contains(t, script, mustJSON("#123456"))
}

func TestDiagnostics_NetworkLogScript_Good(t *testing.T) {
	script := NetworkLogScript(5)

	assert.Contains(t, script, "slice(-5)")
	assert.Contains(t, script, `window.__coreNetworkLog`)
}

func TestDiagnostics_NetworkLogScript_Bad(t *testing.T) {
	script := NetworkLogScript(0)

	assert.Contains(t, script, `performance.getEntriesByType("resource")`)
	assert.NotContains(t, script, "slice(-")
}

func TestDiagnostics_NetworkLogScript_Ugly(t *testing.T) {
	script := NetworkLogScript(-7)

	assert.Contains(t, script, `performance.getEntriesByType("resource")`)
	assert.NotContains(t, script, "slice(-")
}

func TestDiagnostics_normalizeWhitespace_Good(t *testing.T) {
	assert.Equal(t, "hello world", normalizeWhitespace("  hello world  "))
}

func TestDiagnostics_normalizeWhitespace_Bad(t *testing.T) {
	assert.Empty(t, normalizeWhitespace(""))
}

func TestDiagnostics_normalizeWhitespace_Ugly(t *testing.T) {
	assert.Empty(t, normalizeWhitespace("\n\t   "))
}

func mustJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
