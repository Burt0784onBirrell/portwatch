package version_test

import (
	"strings"
	"testing"

	"github.com/danvolchek/portwatch/internal/version"
)

func TestGet_ReturnsInfo(t *testing.T) {
	info := version.Get()

	if info.Version == "" {
		t.Error("expected Version to be non-empty")
	}
	if info.Commit == "" {
		t.Error("expected Commit to be non-empty")
	}
	if info.Date == "" {
		t.Error("expected Date to be non-empty")
	}
}

func TestGet_DefaultsAreSet(t *testing.T) {
	// Without ldflags the package-level vars fall back to their default values.
	info := version.Get()

	if info.Version != version.Version {
		t.Errorf("Get().Version = %q, want %q", info.Version, version.Version)
	}
	if info.Commit != version.Commit {
		t.Errorf("Get().Commit = %q, want %q", info.Commit, version.Commit)
	}
	if info.Date != version.Date {
		t.Errorf("Get().Date = %q, want %q", info.Date, version.Date)
	}
}

func TestInfo_String_ContainsVersion(t *testing.T) {
	info := version.Info{
		Version: "1.2.3",
		Commit:  "abc1234",
		Date:    "2024-01-15",
	}

	s := info.String()

	for _, want := range []string{"portwatch", "1.2.3", "abc1234", "2024-01-15"} {
		if !strings.Contains(s, want) {
			t.Errorf("String() = %q, expected to contain %q", s, want)
		}
	}
}

func TestInfo_String_Format(t *testing.T) {
	info := version.Info{
		Version: "0.1.0",
		Commit:  "deadbeef",
		Date:    "2024-06-01",
	}

	got := info.String()
	want := "portwatch 0.1.0 (commit deadbeef, built 2024-06-01)"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
