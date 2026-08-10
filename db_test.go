package uax

import "testing"

func TestDatabaseRules(t *testing.T) {
	if DatabaseRules() <= 0 {
		t.Error("expected DatabaseRules to be > 0")
	}
	summary := DatabaseSummary()
	if len(summary) == 0 {
		t.Error("expected DatabaseSummary to have entries")
	}
}
