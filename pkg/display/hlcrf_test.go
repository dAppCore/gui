package display

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "card.html"), []byte("<article>Card</article>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte(strings.Join([]string{
		"hlcrf:",
		"  - name: card.html",
		"  - tag: core-inline",
		"    template: <section>Inline</section>",
	}, "\n")), 0o644))

	svc := &Service{}

	script := svc.buildHLCRFComponents(filepath.Join(root, "index.html"))

	require.NotEmpty(t, script)
	assert.Contains(t, script, "customElements.define")
	assert.Contains(t, script, "article>Card</article>")
	assert.Contains(t, script, "<section>Inline</section>")
	assert.Contains(t, script, "core-card")
	assert.Contains(t, script, "core-inline")
}

func TestHLCRF_BuildHLCRFComponents_Bad(t *testing.T) {
	svc := &Service{}

	script := svc.buildHLCRFComponents(filepath.Join(t.TempDir(), "missing.html"))

	assert.Empty(t, script)
}

func TestHLCRF_BuildHLCRFComponents_Ugly(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".core"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("<html></html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".core", "view.yaml"), []byte(strings.Join([]string{
		"hlcrf:",
		"  - name: missing.html",
		"  - template: <span>Fallback</span>",
	}, "\n")), 0o644))

	svc := &Service{}

	script := svc.buildHLCRFComponents(filepath.Join(root, "index.html"))

	require.NotEmpty(t, script)
	assert.Contains(t, script, "<span>Fallback</span>")
	assert.NotContains(t, script, "missing.html")
}
