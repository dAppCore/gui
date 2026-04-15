package marketplace

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name        string            `yaml:"name" json:"name"`
	Version     string            `yaml:"version" json:"version"`
	Repository  string            `yaml:"repository" json:"repository"`
	Ref         string            `yaml:"ref" json:"ref"`
	Description string            `yaml:"description" json:"description"`
	Files       map[string]string `yaml:"files" json:"files"`
	Signature   Signature         `yaml:"signature" json:"signature"`
}

type Signature struct {
	Algorithm string `yaml:"algorithm" json:"algorithm"`
	PublicKey string `yaml:"public_key" json:"public_key"`
	Value     string `yaml:"value" json:"value"`
}

type Installer struct {
	HTTPClient *http.Client
	GitBinary  string
	InstallDir string
}

const maxManifestBytes = 1 << 20

func (i Installer) FetchManifest(ctx context.Context, manifestURL string) (Manifest, error) {
	client := i.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, manifestURL, nil)
	if err != nil {
		return Manifest{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return Manifest{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return Manifest{}, fmt.Errorf("manifest fetch failed: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxManifestBytes+1))
	if err != nil {
		return Manifest{}, err
	}
	if len(body) > maxManifestBytes {
		return Manifest{}, fmt.Errorf("manifest fetch failed: manifest exceeds %d bytes", maxManifestBytes)
	}
	var manifest Manifest
	if err := yaml.Unmarshal(body, &manifest); err != nil {
		return Manifest{}, err
	}
	if manifest.Repository == "" {
		manifest.Repository = manifestURL
	}
	return manifest, nil
}

func VerifyManifest(manifest Manifest) error {
	if strings.ToLower(strings.TrimSpace(manifest.Signature.Algorithm)) != "ed25519" {
		return errors.New("manifest signature algorithm must be ed25519")
	}
	if manifest.Signature.Value == "" || manifest.Signature.PublicKey == "" {
		return errors.New("manifest signature is required")
	}
	payload := manifest.Name + "\n" + manifest.Version + "\n" + manifest.Repository + "\n" + manifest.Ref
	signature, err := base64.StdEncoding.DecodeString(manifest.Signature.Value)
	if err != nil {
		return err
	}
	publicKey, err := base64.StdEncoding.DecodeString(manifest.Signature.PublicKey)
	if err != nil {
		return err
	}
	if len(signature) != ed25519.SignatureSize {
		return errors.New("manifest signature has invalid size")
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return errors.New("manifest public key has invalid size")
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), []byte(payload), signature) {
		return errors.New("manifest signature verification failed")
	}
	return nil
}

func (i Installer) Install(ctx context.Context, manifest Manifest) (string, error) {
	if strings.TrimSpace(i.InstallDir) == "" {
		return "", errors.New("install dir is required")
	}
	if err := VerifyManifest(manifest); err != nil {
		return "", err
	}
	if err := validateManifestName(manifest.Name); err != nil {
		return "", err
	}
	if err := os.MkdirAll(i.InstallDir, 0o755); err != nil {
		return "", err
	}
	targetDir := filepath.Join(i.InstallDir, safeName(manifest.Name))
	rootAbs, err := filepath.Abs(i.InstallDir)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(targetDir)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("install path escapes install dir")
	}
	_ = os.RemoveAll(targetDir)
	args := []string{"clone", "--depth", "1"}
	if manifest.Ref != "" {
		args = append(args, "--branch", manifest.Ref)
	}
	args = append(args, manifest.Repository, targetDir)
	binary := i.GitBinary
	if strings.TrimSpace(binary) == "" {
		binary = "git"
	}
	cmd := exec.CommandContext(ctx, binary, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git clone failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return targetDir, nil
}

func validateManifestName(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errors.New("manifest name is required")
	}
	if strings.ContainsAny(trimmed, `/\`) {
		return errors.New("manifest name must not contain path separators")
	}
	if strings.Contains(trimmed, "..") {
		return errors.New("manifest name must not contain path traversal segments")
	}
	return nil
}

func DigestManifest(manifest Manifest) string {
	hash := sha256.Sum256([]byte(manifest.Name + ":" + manifest.Version + ":" + manifest.Repository + ":" + manifest.Ref))
	return hex.EncodeToString(hash[:])
}

func safeName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return "module"
	}
	var builder strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			builder.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || r == '.':
			builder.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				builder.WriteRune('-')
				lastDash = true
			}
		}
	}
	cleaned := strings.Trim(builder.String(), "-._")
	if cleaned == "" {
		return "module"
	}
	return cleaned
}
