package uax

import (
	"os"
	"strings"
	"testing"
)

// The fingerprint's whole purpose is that the same database yields the same
// value. A fingerprint that varied per process would be worse than none: it
// would look like the rules had changed when nothing had.
func TestDatabaseVersionIsStable(t *testing.T) {
	t.Parallel()

	first := DatabaseVersion()
	if first == "" {
		t.Fatal("DatabaseVersion is empty")
	}
	for i := 0; i < 100; i++ {
		if got := DatabaseVersion(); got != first {
			t.Fatalf("DatabaseVersion changed between calls: %q then %q", first, got)
		}
	}
}

// The tag in the changelog and the constant in the code are read by different
// people and drift apart silently. Assert they agree.
func TestReleaseMatchesTheChangelog(t *testing.T) {
	t.Parallel()

	b, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Skipf("no changelog to check: %v", err)
	}
	want := "## v" + Release
	if !strings.Contains(string(b), want) {
		t.Errorf("CHANGELOG.md has no %q section for Release = %q", want, Release)
	}
}

// DatabaseRules counts what is actually compiled in. It fed a documented
// number in the README once and reported a fraction of the truth, because a
// category had been added to the tables and not to the summary.
func TestDatabaseSummaryCoversEveryCategory(t *testing.T) {
	t.Parallel()

	groups := builtinGroups()

	var counted int
	for _, g := range DatabaseSummary() {
		counted += g.Rules
	}

	var compiled int
	for _, g := range groups {
		compiled += len(g)
	}

	if counted != compiled {
		t.Errorf("DatabaseSummary reports %d rules but %d are compiled — a category is missing from the summary",
			counted, compiled)
	}
}
