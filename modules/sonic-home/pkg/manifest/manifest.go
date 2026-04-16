package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Manifest struct {
	SchemaVersion string       `json:"schema_version"`
	BundleID      string       `json:"bundle_id"`
	Version       string       `json:"version"`
	ReleaseChannel string      `json:"release_channel"`
	BuildDate     string       `json:"build_date"`
	Architecture  []string     `json:"architecture"`
	SourcePath    string       `json:"source_path"`
	SourceDigest  string       `json:"source_digest"`
	Components    []Component  `json:"components"`
	Scripts       ScriptPaths  `json:"scripts"`
	Signatures    []Signature  `json:"signatures"`
	UpdateInfo    UpdateInfo   `json:"update_info"`
}

type Component struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Checksum  string `json:"checksum"`
	SizeBytes int64  `json:"size_bytes"`
}

type ScriptPaths struct {
	PreInstall  string `json:"pre_install"`
	PostInstall string `json:"post_install"`
	Healthcheck string `json:"healthcheck"`
}

type Signature struct {
	KeyID     string `json:"key_id"`
	Algorithm string `json:"algorithm"`
	Signature string `json:"signature"`
}

type UpdateInfo struct {
	PreviousVersion string `json:"previous_version"`
	DeltaAvailable  bool   `json:"delta_available"`
	ChangelogURL    string `json:"changelog_url"`
}

type BuildInput struct {
	SourcePath   string
	Version      string
	Channel      string
	Architecture []string
}

func BuildDryRun(input BuildInput) (Manifest, error) {
	absSource := input.SourcePath
	if absSource == "" {
		absSource = "."
	}

	resolvedSource, err := filepath.Abs(absSource)
	if err != nil {
		return Manifest{}, err
	}

	// A deterministic dry-run digest helps spot source drift.
	sourceDigest := sha256.Sum256([]byte(resolvedSource))

	return Manifest{
		SchemaVersion: "1.0.0",
		BundleID:      "com.uhome.nest",
		Version:       defaultIfEmpty(input.Version, "0.0.0-dev"),
		ReleaseChannel: defaultIfEmpty(input.Channel, "edge"),
		BuildDate:     time.Now().UTC().Format(time.RFC3339),
		Architecture:  defaultArchitectures(input.Architecture),
		SourcePath:    resolvedSource,
		SourceDigest:  "sha256:" + hex.EncodeToString(sourceDigest[:]),
		Components: []Component{
			{
				Name:      "uhome-api",
				Type:      "binary",
				Source:    "payload/base/usr/local/bin/uhome-api",
				Checksum:  "sha256:pending",
				SizeBytes: 0,
			},
		},
		Scripts: ScriptPaths{
			PreInstall:  "payload/scripts/pre-install.sh",
			PostInstall: "payload/scripts/post-install.sh",
			Healthcheck: "payload/scripts/healthcheck.sh",
		},
		Signatures: []Signature{},
		UpdateInfo: UpdateInfo{
			PreviousVersion: "",
			DeltaAvailable:  false,
			ChangelogURL:    "",
		},
	}, nil
}

func Write(path string, manifest Manifest) error {
	body, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o644)
}

func defaultArchitectures(arch []string) []string {
	if len(arch) == 0 {
		return []string{"amd64"}
	}
	return arch
}

func defaultIfEmpty(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
