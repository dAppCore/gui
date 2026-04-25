package marketplace

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
	t.Run("invalid yaml", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(": not-yaml"))
		}))
		t.Cleanup(server.Close)

		installer := Installer{HTTPClient: server.Client()}
		_, err := installer.FetchManifest(context.Background(), server.URL)
		require.Error(t, err)
	})

	t.Run("size limit", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("name: " + strings.Repeat("x", maxManifestBytes)))
		}))
		t.Cleanup(server.Close)

		installer := Installer{HTTPClient: server.Client()}
		_, err := installer.FetchManifest(context.Background(), server.URL)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds")
	})
}

func TestMarketplace_List_Bad(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("boom"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	_, err := installer.List(context.Background(), server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "marketplace list failed")
}

func TestMarketplace_List_Ugly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("a", maxManifestBytes+1)))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	_, err := installer.List(context.Background(), server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds")
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
	resolvedRoot, err := filepath.EvalSymlinks(targetRoot)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(resolvedRoot, "core-ui"), targetDir)
	_, err = os.Stat(targetDir)
	require.NoError(t, err)

	contents, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.Contains(t, string(contents), "clone")
	assert.Contains(t, string(contents), "--branch")
	assert.Contains(t, string(contents), "main")
	assert.Contains(t, string(contents), "--")

	installedManifest, err := os.ReadFile(filepath.Join(targetDir, ".core", "marketplace.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(installedManifest), "name: Core UI")
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

func TestMarketplace_Install_RejectsDashPrefixedRepository(t *testing.T) {
	installer := Installer{InstallDir: t.TempDir()}
	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "--upload-pack=sh",
		Ref:        "main",
	}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "repository must not begin with a dash")
}

func TestMarketplace_Install_RejectsDashPrefixedRef(t *testing.T) {
	installer := Installer{InstallDir: t.TempDir()}
	_, err := installer.Install(context.Background(), signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "--upload-pack=sh",
	}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ref must not begin with a dash")
}

func TestMarketplace_Install_Ugly(t *testing.T) {
	scriptDir := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "git")
	require.NoError(t, os.WriteFile(scriptPath, []byte("#!/bin/sh\nprintf '%s\\n' 'fatal: https://token:secret@example.com/repo.git' >&2\nexit 1\n"), 0o755))

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
	assert.NotContains(t, err.Error(), "secret")
	assert.NotContains(t, err.Error(), "token:")
}

func TestMarketplace_Install_CleansUpOnCloneFailure(t *testing.T) {
	scriptDir := t.TempDir()
	scriptPath := filepath.Join(scriptDir, "git")
	targetRoot := t.TempDir()
	script := "#!/bin/sh\nlast=''\nfor arg in \"$@\"; do last=\"$arg\"; done\nmkdir -p \"$last\"\ntouch \"$last/partial\"\nexit 1\n"
	require.NoError(t, os.WriteFile(scriptPath, []byte(script), 0o755))

	installer := Installer{
		GitBinary:  scriptPath,
		InstallDir: targetRoot,
	}

	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	_, err := installer.Install(context.Background(), manifest)
	require.Error(t, err)

	targetDir := filepath.Join(targetRoot, "core-ui")
	_, statErr := os.Stat(targetDir)
	assert.Error(t, statErr)
	assert.True(t, os.IsNotExist(statErr))
}

func TestMarketplace_Verify_Good(t *testing.T) {
	manifest := signedManifest(t, Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		data, err := yaml.Marshal(manifest)
		require.NoError(t, err)
		_, _ = w.Write(data)
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	verified, err := installer.Verify(context.Background(), server.URL)
	require.NoError(t, err)
	assert.Equal(t, manifest.Name, verified.Name)
	assert.Equal(t, manifest.Ref, verified.Ref)
}

func TestMarketplace_List_Good(t *testing.T) {
	manifests := []Manifest{
		{Name: "core-ui", Version: "1.2.3"},
		{Name: "core-chat", Version: "0.9.0"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("manifests:\n  - name: core-ui\n    version: 1.2.3\n  - name: core-chat\n    version: 0.9.0\n"))
	}))
	t.Cleanup(server.Close)

	installer := Installer{HTTPClient: server.Client()}
	listed, err := installer.List(context.Background(), server.URL)
	require.NoError(t, err)
	assert.Equal(t, manifests, listed)
}

func TestMarketplace_DigestManifest_Good(t *testing.T) {
	manifest := Manifest{
		Name:       "core-ui",
		Version:    "1.2.3",
		Repository: "https://example.com/core-ui.git",
		Ref:        "main",
	}

	got := DigestManifest(manifest)
	expected := sha256.Sum256([]byte(manifest.Name + ":" + manifest.Version + ":" + manifest.Repository + ":" + manifest.Ref))

	assert.Equal(t, hex.EncodeToString(expected[:]), got)
}

func TestMarketplace_DigestManifest_Bad(t *testing.T) {
	base := Manifest{Name: "core-ui", Version: "1.2.3", Repository: "https://example.com/core-ui.git", Ref: "main"}
	changed := base
	changed.Ref = "dev"

	assert.NotEqual(t, DigestManifest(base), DigestManifest(changed))
}

func TestMarketplace_DigestManifest_Ugly(t *testing.T) {
	got := DigestManifest(Manifest{})

	assert.Len(t, got, 64)
	assert.NotEmpty(t, got)
}

func TestMarketplace_safeName_EmptyFallbackUsesInputDigest(t *testing.T) {
	slashes := safeName("////")
	ats := safeName("@@@")

	assert.Equal(t, "module-0ea28b45", slashes)
	assert.Equal(t, "module-2ec847d8", ats)
	assert.NotEqual(t, slashes, ats)
	assertSafeModuleName(t, slashes)
	assertSafeModuleName(t, ats)
	assert.Equal(t, "valid-name", safeName("valid-name"))
}

func TestMarketplace_validateManifestName_Good(t *testing.T) {
	require.NoError(t, validateManifestName("core-ui"))
}

func TestMarketplace_validateManifestName_Bad(t *testing.T) {
	err := validateManifestName("")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestMarketplace_validateManifestName_Ugly(t *testing.T) {
	err := validateManifestName("..")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

func TestMarketplace_validateCloneArg_Good(t *testing.T) {
	require.NoError(t, validateCloneArg("repository", "https://example.com/core-ui.git"))
}

func TestMarketplace_validateCloneArg_Bad(t *testing.T) {
	err := validateCloneArg("repository", "--upload-pack=sh")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dash")
}

func TestMarketplace_validateCloneArg_Ugly(t *testing.T) {
	err := validateCloneArg("repository", "https://example.com/core-ui.git\n--depth 1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid control characters")
}

func TestMarketplace_validateRepositorySource_Good(t *testing.T) {
	require.NoError(t, validateRepositorySource("https://example.com/core-ui.git"))
	require.NoError(t, validateRepositorySource("git@example.com:core-ui.git"))
}

func TestMarketplace_validateRepositorySource_Bad(t *testing.T) {
	err := validateRepositorySource("ext::sh -c id")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "remote helper")
}

func TestMarketplace_validateRepositorySource_Ugly(t *testing.T) {
	err := validateRepositorySource("file:///tmp/core-ui.git")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not allowed")
}

func TestMarketplace_decodeManifestList_Good(t *testing.T) {
	manifests, err := decodeManifestList([]byte(`[{"name":"core-ui","version":"1.2.3"},{"name":"core-chat","version":"0.9.0"}]`))

	require.NoError(t, err)
	require.Len(t, manifests, 2)
	assert.Equal(t, "core-ui", manifests[0].Name)
	assert.Equal(t, "core-chat", manifests[1].Name)
}

func TestMarketplace_decodeManifestList_Bad(t *testing.T) {
	manifests, err := decodeManifestList([]byte("   "))

	require.NoError(t, err)
	assert.Nil(t, manifests)
}

func TestMarketplace_decodeManifestList_Ugly(t *testing.T) {
	_, err := decodeManifestList([]byte(": not-yaml"))

	require.Error(t, err)
}

func TestMarketplace_sanitizeCommandOutput_Good(t *testing.T) {
	got := sanitizeCommandOutput([]byte("fatal: https://token:secret@example.com/repo.git"))

	assert.Contains(t, got, "[redacted]@")
	assert.NotContains(t, got, "secret")
	assert.NotContains(t, got, "token:")
}

func TestMarketplace_sanitizeCommandOutput_Bad(t *testing.T) {
	assert.Equal(t, "command produced no output", sanitizeCommandOutput(nil))
	assert.Equal(t, "command produced no output", sanitizeCommandOutput([]byte("   \n")))
}

func TestMarketplace_sanitizeCommandOutput_Ugly(t *testing.T) {
	got := sanitizeCommandOutput([]byte(strings.Repeat("a", 1024)))

	assert.Len(t, got, 515)
	assert.True(t, strings.HasSuffix(got, "..."))
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func assertSafeModuleName(t *testing.T, value string) {
	t.Helper()

	require.NotEmpty(t, value)
	assert.LessOrEqual(t, len(value), 32)
	for _, r := range value {
		assert.Truef(t, (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-', "unsafe character %q in %q", r, value)
	}
}
