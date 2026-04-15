package marketplace

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func signedManifest(t *testing.T, manifest Manifest) Manifest {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	payload := manifest.Name + "\n" + manifest.Version + "\n" + manifest.Repository + "\n" + manifest.Ref
	signature := ed25519.Sign(priv, []byte(payload))
	manifest.Signature = Signature{
		Algorithm: "ed25519",
		PublicKey: base64.StdEncoding.EncodeToString(pub),
		Value:     base64.StdEncoding.EncodeToString(signature),
	}
	return manifest
}

func TestMarketplace_FetchManifest_Good(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("name: core-ui\nversion: 1.2.3\nref: main\n"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	manifest, err := installer.FetchManifest(context.Background(), server.URL)
	require.NoError(t, err)
	assert.Equal(t, "core-ui", manifest.Name)
	assert.Equal(t, "1.2.3", manifest.Version)
	assert.Equal(t, server.URL, manifest.Repository)
	assert.Equal(t, "main", manifest.Ref)
}

func TestMarketplace_FetchManifest_Bad(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	_, err := installer.FetchManifest(context.Background(), server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "manifest fetch failed")
}

func TestMarketplace_FetchManifest_Ugly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(": not-yaml"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	_, err := installer.FetchManifest(context.Background(), server.URL)
	require.Error(t, err)
}

func TestMarketplace_VerifyManifest_Good(t *testing.T) {
	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})

	require.NoError(t, VerifyManifest(manifest))
}

func TestMarketplace_VerifyManifest_Bad(t *testing.T) {
	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	manifest.Ref = "dev"

	err := VerifyManifest(manifest)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "manifest signature verification failed")
}

func TestMarketplace_VerifyManifest_Ugly(t *testing.T) {
	manifest := Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
		Signature: Signature{
			Algorithm: "ed25519",
			PublicKey: "not-base64",
			Value:     "also-not-base64",
		},
	}

	require.Error(t, VerifyManifest(manifest))
}

func TestMarketplace_VerifyManifest_RequiresSignature(t *testing.T) {
	manifest := Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}

	require.Error(t, VerifyManifest(manifest))
}

func TestMarketplace_Install_Good(t *testing.T) {
	scriptDir := t.TempDir()
	logFile := filepath.Join(scriptDir, "git.log")
	targetRoot := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "git")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + shellQuote(logFile) + "\nlast=''\nfor arg in \"$@\"; do last=\"$arg\"; done\nmkdir -p \"$last\"\nexit 0\n"
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0o755))

	installer := Installer{
		GitBinary:  scriptPath,
		InstallDir: targetRoot,
	}

	targetDir, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "Core UI",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}))
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(targetRoot, "core-ui"), targetDir)
	_, err = os.Stat(targetDir)
	require.NoError(t, err)

	contents, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.Contains(t, string(contents), "clone")
	assert.Contains(t, string(contents), "--branch")
	assert.Contains(t, string(contents), "main")
}

func TestMarketplace_Install_Bad(t *testing.T) {
	installer := Installer{InstallDir: ""}
	_, err := installer.Install(context.Background(), Manifest{Name: "core-ui"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "install dir is required")
}

func TestMarketplace_Install_RejectsTraversalName(t *testing.T) {
	installer := Installer{InstallDir: t.TempDir()}
	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "../../escape",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path separators")
}

func TestMarketplace_Install_Ugly(t *testing.T) {
	scriptDir := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "git")
	require.NoError(t, os.WriteFile(scriptPath, []byte("#!/bin/sh\nexit 1\n"), 0o755))

	installer := Installer{
		GitBinary:  scriptPath,
		InstallDir: t.TempDir(),
	}

	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "core-ui",
		Repository: "https://example.com/core-ui.git",
	}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "git clone failed")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
