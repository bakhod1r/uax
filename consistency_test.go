package uax

import "testing"

func TestImpossibleCombinationsAreDetected(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ua   string
		rule string
	}{
		{
			"safari on windows",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			"impossible.safari_on_windows",
		},
		{
			"iphone and android at once",
			"Mozilla/5.0 (iPhone; Android 13; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
			"impossible.two_platforms",
		},
		{
			"impossible browser version",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/9000.0.0.0 Safari/537.36",
			"anachronism.future_browser",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cons := Parse(tc.ua).Consistency()
			if !cons.Has(tc.rule) {
				t.Errorf("rule %q did not fire; findings: %+v", tc.rule, cons.Findings)
			}
			if cons.OK() {
				t.Error("OK() reported a coherent client")
			}
			if cons.Score <= 0 {
				t.Errorf("score: got %v, want > 0", cons.Score)
			}
		})
	}
}

func TestOrdinaryClientsAreCoherent(t *testing.T) {
	t.Parallel()

	for _, ua := range []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"curl/8.4.0",
	} {
		if cons := Parse(ua).Consistency(); !cons.OK() {
			t.Errorf("a real client was flagged: %+v\n  ua: %s", cons.Findings, ua)
		}
	}
}

// The score reflects the single strongest finding rather than a sum, so a
// client with several weak oddities never inherits a critical score.
func TestScoreIsTheStrongestFinding(t *testing.T) {
	t.Parallel()

	c := &Client{hintFindings: []Finding{
		{Rule: "a", Severity: SeverityLow},
		{Rule: "b", Severity: SeverityMedium},
		{Rule: "c", Severity: SeverityLow},
	}}

	if got := CheckConsistency(c).Score; got != SeverityMedium.weight() {
		t.Errorf("score: got %v, want %v", got, SeverityMedium.weight())
	}
}

func TestConsistencyEdgeCases(t *testing.T) {
	if got := Severity("unknown").weight(); got != 0 {
		t.Errorf("expected unknown severity to be 0, got %v", got)
	}

	c := Consistency{Findings: []Finding{{Rule: "some_rule"}}}
	if c.Has("other_rule") {
		t.Error("expected false for missing rule")
	}

	if got := maxPlausibleMajor("IE"); got != 12 {
		t.Errorf("expected 12 for IE, got %d", got)
	}
}
