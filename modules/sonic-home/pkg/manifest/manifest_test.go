package manifest

import "testing"

func TestBuildDryRunDefaults(t *testing.T) {
	got, err := BuildDryRun(BuildInput{})
	if err != nil {
		t.Fatalf("BuildDryRun returned error: %v", err)
	}

	if got.SchemaVersion != "1.0.0" {
		t.Fatalf("unexpected schema version: %s", got.SchemaVersion)
	}

	if got.BundleID != "com.uhome.nest" {
		t.Fatalf("unexpected bundle id: %s", got.BundleID)
	}

	if len(got.Architecture) == 0 {
		t.Fatal("expected default architecture")
	}
}
