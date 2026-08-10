package table

import (
	"testing"
)

// These are the invariants the YAML generator used to enforce before the rules
// moved into Go. They are the reason the database can be trusted, so they
// outlived the generator: a rule that breaks one of them fails the build.

// groups is every category, named so a failure says which file to open.
func groups() map[string][]Rule {
	return map[string][]Rule{
		"Browsers":    Browsers,
		"Engines":     Engines,
		"OSes":        OSes,
		"Bots":        Bots,
		"Apps":        Apps,
		"Automations": Automations,
		"Clients":     Clients,
	}
}

func TestRuleIDsAreUniqueWithinACategory(t *testing.T) {
	// Uniqueness is per category, which is the contract the database was
	// written to: one token can be evidence in two categories at once, and
	// "headless_chrome" is deliberately both a browser and an automation
	// signal. An ID identifies a rule within the category that defines it.
	for name, rules := range groups() {
		seen := make(map[string]bool, len(rules))
		for _, r := range rules {
			if r.ID == "" {
				t.Errorf("%s: rule with empty ID", name)
				continue
			}
			if seen[r.ID] {
				t.Errorf("%s: duplicate rule ID %q", name, r.ID)
				continue
			}
			seen[r.ID] = true
		}
	}
}

func TestEveryRuleHasExactlyOneMatchClause(t *testing.T) {
	for name, rules := range groups() {
		for _, r := range rules {
			n := 0
			for _, v := range []string{r.Token, r.Prefix, r.Contains} {
				if v != "" {
					n++
				}
			}
			if n != 1 {
				t.Errorf("%s: rule %q has %d match clauses, want exactly 1", name, r.ID, n)
			}
		}
	}
}

func TestEveryRuleHasAPriorityAndAConfidence(t *testing.T) {
	for name, rules := range groups() {
		for _, r := range rules {
			if r.Priority <= 0 {
				t.Errorf("%s: rule %q has priority %d, want > 0", name, r.ID, r.Priority)
			}
			if r.Confidence <= 0 || r.Confidence > 1 {
				t.Errorf("%s: rule %q has confidence %v, want (0,1]", name, r.ID, r.Confidence)
			}
		}
	}
}

func TestRulesAreSortedByDescendingPriority(t *testing.T) {
	// The resolver stops at the first match, so order is not cosmetic: an
	// out-of-order rule silently shadows a more specific one.
	for name, rules := range groups() {
		for i := 1; i < len(rules); i++ {
			if rules[i-1].Priority < rules[i].Priority {
				t.Errorf("%s: rule %q (priority %d) precedes %q (priority %d); highest priority must come first",
					name, rules[i-1].ID, rules[i-1].Priority, rules[i].ID, rules[i].Priority)
			}
		}
	}
}

func TestEveryRuleIsProvenByATest(t *testing.T) {
	// The rule that keeps the database honest as it grows: a rule with no test
	// is a claim nobody checked.
	proven := make(map[string]bool, len(Tests))
	for _, tc := range Tests {
		proven[tc.RuleID] = true
	}
	for name, rules := range groups() {
		for _, r := range rules {
			if !proven[r.ID] {
				t.Errorf("%s: rule %q has no test in tests.go", name, r.ID)
			}
		}
	}
}

func TestEveryTestNamesARealRule(t *testing.T) {
	// The other direction: a test whose rule was deleted stops proving
	// anything and would otherwise sit there looking like coverage.
	known := make(map[string]bool)
	for _, rules := range groups() {
		for _, r := range rules {
			known[r.ID] = true
		}
	}
	for _, tc := range Tests {
		if !known[tc.RuleID] {
			t.Errorf("test for rule %q, which no category defines", tc.RuleID)
		}
		if tc.UA == "" {
			t.Errorf("test for rule %q has an empty User-Agent", tc.RuleID)
		}
		if len(tc.Expect) == 0 {
			t.Errorf("test for rule %q expects nothing", tc.RuleID)
		}
	}
}
