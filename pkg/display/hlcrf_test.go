package display

import (
	"path/filepath"
	"strings"
	"testing"

	coreio "dappco.re/go/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHLCRF_DefaultHLCRFTag_Good(t *testing.T) {
	assert.Equal(t, "core-widget", defaultHLCRFTag("Widget.ts"))
}

func TestHLCRF_DefaultHLCRFTag_Bad(t *testing.T) {
	assert.Equal(t, "feature-card", defaultHLCRFTag("feature_card.html"))
}

func TestHLCRF_DefaultHLCRFTag_Ugly(t *testing.T) {
	assert.Equal(t, "core-", defaultHLCRFTag(""))
}

func TestHLCRF_BuildHLCRFComponents_Good(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "card.html"), "<article>Card</article>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"hlcrf:",
		"  - name: card.html",
		"  - tag: core-inline",
		"    template: <section>Inline</section>",
	}, "\n"), 0o644))

	svc := &Service{}

	script, err := svc.buildHLCRFComponents(filepath.Join(root, "index.html"))

	require.NoError(t, err)
	require.NotEmpty(t, script)
	assert.Contains(t, script, "customElements.define")
	assert.Contains(t, script, "article>Card</article>")
	assert.Contains(t, script, "<section>Inline</section>")
	assert.Contains(t, script, "core-card")
	assert.Contains(t, script, "core-inline")
}

func TestHLCRF_CompileHLCRFTemplate_Good(t *testing.T) {
	compiled := compileHLCRFTemplate(`<section data-slot="H">{{slot "H"}}</section><main>{{ slot "L-C" }}</main><footer>{{ slot "" }}{{ slot "default" }}</footer>`)

	assert.Contains(t, compiled, `<slot name="H"></slot>`)
	assert.Contains(t, compiled, `<slot name="L-C"></slot>`)
	assert.Contains(t, compiled, `<footer><slot></slot><slot></slot></footer>`)
}

func TestHLCRF_BuildHLCRFComponents_Bad(t *testing.T) {
	svc := &Service{}

	script, err := svc.buildHLCRFComponents(filepath.Join(t.TempDir(), "missing.html"))

	require.NoError(t, err)
	assert.Empty(t, script)
}

func TestHLCRF_BuildHLCRFComponents_Ugly(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"hlcrf:",
		"  - name: missing.html",
		"  - template: <span>Fallback</span>",
	}, "\n"), 0o644))

	svc := &Service{}

	script, err := svc.buildHLCRFComponents(filepath.Join(root, "index.html"))

	require.NoError(t, err)
	require.NotEmpty(t, script)
	assert.Contains(t, script, "<span>Fallback</span>")
	assert.NotContains(t, script, "missing.html")
}

func TestHLCRF_BuildHLCRFComponents_RejectsTraversal(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, coreio.Local.EnsureDir(filepath.Join(root, ".core")))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "index.html"), "<html></html>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, "outside.html"), "<span>Outside</span>", 0o644))
	require.NoError(t, coreio.Local.WriteMode(filepath.Join(root, ".core", "view.yaml"), strings.Join([]string{
		"hlcrf:",
		"  - name: ../outside.html",
	}, "\n"), 0o644))

	svc := &Service{}

	script, err := svc.buildHLCRFComponents(filepath.Join(root, "index.html"))

	require.Error(t, err)
	assert.Empty(t, script)
}
