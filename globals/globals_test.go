package globals

import (
	"strings"
	"testing"
)

func TestVersionNumberExistence(t *testing.T) {
	if len(VersionHistory) == 0 {
		t.Fatalf("VersionHistory is empty")
	}
}

func TestVersionFormats(t *testing.T) {
	for i, version := range VersionHistory {
		if !strings.HasPrefix(version.SemVer, "v") {
			t.Fatalf("%d version does not start with 'v ': got %s", i, version.SemVer)
		}

		noPrefixSemVer := strings.TrimPrefix(version.SemVer, "v")
		parts := strings.Split(noPrefixSemVer, ".")
		if len(parts) != 3 {
			t.Fatalf("%d Version number does not follow major.minor.patch format: %s", i, noPrefixSemVer)
		}
	}
}

func TestDuplicateVersions(t *testing.T) {
	versionMap := make(map[string]bool)
	for i, vh := range VersionHistory {
		if _, exists := versionMap[vh.SemVer]; exists {
			t.Fatalf("Duplicate version found: %s", vh.SemVer)
		}
		versionMap[vh.SemVer] = true

		if i == 0 && vh.SemVer != "v"+Version {
			t.Fatalf("Global Version (%s) doesn't match first VersionHistory entry (%s)", Version, vh.SemVer)
		}
	}
}

func TestNotLocal(t *testing.T) {
	if Env != "prod" {
		t.Fatalf("Expected Env to be prod, got %s", Env)
	}
}
