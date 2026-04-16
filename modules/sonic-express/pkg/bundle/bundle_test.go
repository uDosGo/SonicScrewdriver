package bundle

import (
	"path/filepath"
	"testing"

	"github.com/sonic-family/sonic-express/pkg/manifest"
)

func TestBuildAndVerifyDraftBundle(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "test.she")

	m, err := manifest.BuildDryRun(manifest.BuildInput{
		SourcePath: tmp,
		Version:    "1.0.0",
		Channel:    "stable",
	})
	if err != nil {
		t.Fatalf("BuildDryRun error: %v", err)
	}

	if err := BuildDraftBundle(out, m); err != nil {
		t.Fatalf("BuildDraftBundle error: %v", err)
	}

	if err := VerifyDraftBundle(out); err != nil {
		t.Fatalf("VerifyDraftBundle error: %v", err)
	}
}
