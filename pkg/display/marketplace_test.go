package display

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/marketplace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMarketplace_marketplaceRegistryURL_Good(t *testing.T) {
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", "")

	opts := core.NewOptions(
		core.Option{Key: "url", Value: "  https://override.example/registry  "},
	)

	require.Equal(t, "https://override.example/registry", marketplaceRegistryURL(opts))
}

func TestMarketplace_marketplaceRegistryURL_Bad(t *testing.T) {
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", "")

	require.Empty(t, marketplaceRegistryURL(core.NewOptions()))
}

func TestMarketplace_marketplaceRegistryURL_Ugly(t *testing.T) {
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", "  https://env.example/registry  ")

	require.Equal(t, "https://env.example/registry", marketplaceRegistryURL(core.NewOptions()))
}

func TestMarketplace_marketplaceInstallRoot_Good(t *testing.T) {
	root := marketplaceInstallRoot("  /tmp/custom/apps  ")

	require.Equal(t, "/tmp/custom/apps", root)
}

func TestMarketplace_marketplaceInstallRoot_Bad(t *testing.T) {
	t.Setenv("DIR_HOME", "")

	require.True(t, strings.HasSuffix(marketplaceInstallRoot(""), filepath.Join(".core", "apps")))
}

func TestMarketplace_marketplaceInstallRoot_Ugly(t *testing.T) {
	t.Setenv("DIR_HOME", "  /Users/tester  ")

	require.True(t, strings.HasSuffix(marketplaceInstallRoot(""), filepath.Join(".core", "apps")))
}

func TestMarketplace_registerMarketplaceActions_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	registry := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("manifests:\n  - name: core-ui\n    version: 1.2.3\n"))
	}))
	t.Cleanup(registry.Close)
	t.Setenv("CORE_MARKETPLACE_REGISTRY_URL", registry.URL)

	listResult := c.Action("display.marketplace.list").Run(context.Background(), core.NewOptions())
	require.True(t, listResult.OK)

	payload, ok := listResult.Value.(map[string]any)
	require.True(t, ok)
	require.Equal(t, registry.URL, payload["registry_url"])
	manifests, ok := payload["manifests"].([]marketplace.Manifest)
	require.True(t, ok)
	require.Len(t, manifests, 1)
	assert.Equal(t, "core-ui", manifests[0].Name)

	manifest := signedMarketplaceManifest(t, marketplace.Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	manifestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, err := yaml.Marshal(manifest)
		require.NoError(t, err)
		_, _ = w.Write(data)
	}))
	t.Cleanup(manifestServer.Close)

	fetchResult := c.Action("display.marketplace.fetch").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: manifestServer.URL},
	))
	require.True(t, fetchResult.OK)
	fetched, ok := fetchResult.Value.(marketplace.Manifest)
	require.True(t, ok)
	assert.Equal(t, manifest.Name, fetched.Name)
	assert.Equal(t, manifest.Ref, fetched.Ref)

	verifyResult := c.Action("display.marketplace.verify").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: manifestServer.URL},
	))
	require.True(t, verifyResult.OK)
	verified, ok := verifyResult.Value.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, marketplace.DigestManifest(manifest), verified["digest"])

	installDir := t.TempDir()
	gitLogDir := t.TempDir()
	gitLog := filepath.Join(gitLogDir, "git.log")
	gitBinary := filepath.Join(gitLogDir, "git")
	require.NoError(t, os.WriteFile(gitBinary, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > "+shellQuote(gitLog)+"\nlast=''\nfor arg in \"$@\"; do last=\"$arg\"; done\nmkdir -p \"$last\"\nexit 0\n"), 0o755))

	installResult := c.Action("display.marketplace.install").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: manifestServer.URL},
		core.Option{Key: "install_dir", Value: installDir},
		core.Option{Key: "git_binary", Value: gitBinary},
	))
	require.True(t, installResult.OK)
	installed, ok := installResult.Value.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, installDir, installed["install_dir"])
	resolvedInstallDir, err := filepath.EvalSymlinks(installDir)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(resolvedInstallDir, "core-ui"), installed["target_dir"])

	contents, err := os.ReadFile(gitLog)
	require.NoError(t, err)
	assert.Contains(t, string(contents), "clone")
	assert.Contains(t, string(contents), "--branch")
	assert.Contains(t, string(contents), "--")
}

func TestMarketplace_registerMarketplaceActions_Bad(t *testing.T) {
	_, c := newTestDisplayService(t)

	result := c.Action("display.marketplace.fetch").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "manifest url is required")
}

func TestMarketplace_registerMarketplaceActions_Ugly(t *testing.T) {
	_, c := newTestDisplayService(t)

	result := c.Action("display.marketplace.install").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "manifest url is required")
}

func signedMarketplaceManifest(t *testing.T, manifest marketplace.Manifest) marketplace.Manifest {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	payload := manifest.Name + "\n" + manifest.Version + "\n" + manifest.Repository + "\n" + manifest.Ref
	signature := ed25519.Sign(priv, []byte(payload))
	manifest.Signature = marketplace.Signature{
		Algorithm: "ed25519",
		PublicKey: base64.StdEncoding.EncodeToString(pub),
		Value:     base64.StdEncoding.EncodeToString(signature),
	}
	return manifest
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
