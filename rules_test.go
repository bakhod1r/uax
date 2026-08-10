package uax

import (
	"testing"

	"github.com/bakhod1r/uax/internal/table"
)

// TestInlineRuleTests runs the expectations that travel with each rule in
// internal/table. A rule and its proof ship together — the lint suite refuses
// a rule that has no test — and this runs every one of them against the real
// parser rather than against the matcher alone.
func TestInlineRuleTests(t *testing.T) {
	t.Parallel()

	if len(table.Tests) == 0 {
		t.Fatal("internal/table declares no rule tests")
	}

	for _, tc := range table.Tests {
		t.Run(tc.RuleID, func(t *testing.T) {
			t.Parallel()
			c := Parse(tc.UA)

			for field, want := range tc.Expect {
				got := ""
				switch field {
				case "browser":
					got = c.Browser.Name
				case "engine":
					got = c.Engine.Name
				case "os":
					got = c.OS.Name
				case "brand":
					got = c.Device.Brand
				case "name":
					got = c.Device.Name
				case "model":
					got = c.Device.Model
				case "type":
					got = string(c.Device.Type)
				case "bot":
					if c.Bot != nil {
						got = c.Bot.Name
					}
				case "category":
					// "category" means whichever categorised thing the rule
					// under test produces.
					switch {
					case c.Bot != nil:
						got = string(c.Bot.Category)
					case c.HTTPClient != nil:
						got = c.HTTPClient.Category
					case c.Automation != nil:
						got = c.Automation.Category
					case c.App != nil:
						got = c.App.Category
					}
				case "provider":
					if c.Bot != nil {
						got = c.Bot.Provider
					}
				case "app":
					if c.App != nil {
						got = c.App.Name
					}
				case "automation":
					if c.Automation != nil {
						got = c.Automation.Name
					}
				case "client":
					if c.HTTPClient != nil {
						got = c.HTTPClient.Name
					}
				default:
					t.Fatalf("unknown expectation field %q — add it here or fix the rule", field)
				}

				if got != want {
					t.Errorf("%s: got %q, want %q\n  ua: %s", field, got, want, tc.UA)
				}
			}
		})
	}
}

// TestEveryRuleHasATest is the guard that keeps the database honest as it
// grows. The generator enforces it too; this fails loudly if a generated table
// is ever committed out of sync with its data.
func TestEveryRuleHasATest(t *testing.T) {
	t.Parallel()

	tested := make(map[string]bool, len(table.Tests))
	for _, tc := range table.Tests {
		tested[tc.RuleID] = true
	}

	groups := map[string][]table.Rule{
		"browser":    table.Browsers,
		"engine":     table.Engines,
		"os":         table.OSes,
		"bot":        table.Bots,
		"app":        table.Apps,
		"automation": table.Automations,
		"client":     table.Clients,
	}

	for group, rules := range groups {
		for _, r := range rules {
			if !tested[r.ID] {
				t.Errorf("%s rule %q has no test", group, r.ID)
			}
		}
	}
}

// TestRulesArePrioritySorted protects the early-exit logic in index.match,
// which is only correct while rules are ordered by descending priority.
func TestRulesArePrioritySorted(t *testing.T) {
	t.Parallel()

	groups := map[string][]table.Rule{
		"browser": table.Browsers, "engine": table.Engines, "os": table.OSes,
		"automation": table.Automations, "client": table.Clients,
	}

	for group, rules := range groups {
		for i := 1; i < len(rules); i++ {
			if rules[i].Priority > rules[i-1].Priority {
				t.Errorf("%s: rule %q (priority %d) sorts after %q (priority %d)",
					group, rules[i].ID, rules[i].Priority, rules[i-1].ID, rules[i-1].Priority)
			}
		}
	}
}
