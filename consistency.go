package uax

// Severity is how much weight a finding carries.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

func (s Severity) weight() float32 {
	switch s {
	case SeverityCritical:
		return 1.00
	case SeverityHigh:
		return 0.60
	case SeverityMedium:
		return 0.30
	case SeverityLow:
		return 0.10
	}
	return 0
}

// Finding is one internal contradiction, with the evidence behind it.
//
// Rule is a stable identifier, safe to log, alert on and assert in tests. It
// never changes spelling once released, so a dashboard built on it keeps
// working.
type Finding struct {
	Rule     string   `json:"rule"`
	Severity Severity `json:"severity"`
	Detail   string   `json:"detail,omitempty"`
	Left     string   `json:"left,omitempty"`
	Right    string   `json:"right,omitempty"`
}

// Consistency is the result of checking a client's claims against each other.
//
// It answers "does this client's story hold together", which is a different
// and more answerable question than "is this client lying". A clean result is
// not proof of honesty; a dirty one is strong evidence of editing.
type Consistency struct {
	Findings []Finding `json:"findings,omitempty"`
	Score    float32   `json:"score"` // 0 = coherent, 1 = impossible
}

// OK reports whether nothing contradicted anything else.
func (c Consistency) OK() bool { return len(c.Findings) == 0 }

// Has reports whether a specific rule fired.
func (c Consistency) Has(rule string) bool {
	for _, f := range c.Findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

// spoofRule is one impossibility check over an already-parsed client.
type spoofRule struct {
	id       string
	severity Severity
	detail   string
	match    func(*Client) (left, right string, hit bool)
}

// spoofRules are combinations that cannot occur on real hardware.
//
// Each one is an impossibility, not a heuristic: no amount of unusual
// configuration produces Safari on Windows NT in 2026, or a Gecko engine on
// iOS where Apple permits only WebKit. Heuristics that merely look unusual
// belong in Verdict, where they influence risk without asserting a fact.
var spoofRules = []spoofRule{
	{
		id:       "impossible.safari_on_windows",
		severity: SeverityCritical,
		detail:   "Safari has not shipped for Windows since 2012",
		match: func(c *Client) (string, string, bool) {
			return c.Browser.Name, c.OS.Name, c.Browser.Name == "Safari" && c.OS.Name == "Windows"
		},
	},
	{
		id:       "impossible.safari_on_android",
		severity: SeverityCritical,
		detail:   "Safari has never shipped for Android",
		match: func(c *Client) (string, string, bool) {
			return c.Browser.Name, c.OS.Name, c.Browser.Name == "Safari" && c.OS.Family == "Android"
		},
	},
	{
		id:       "impossible.non_webkit_on_ios",
		severity: SeverityCritical,
		detail:   "iOS permits only WebKit, so every browser there reports it",
		match: func(c *Client) (string, string, bool) {
			hit := c.OS.Family == "iOS" && c.Engine.Name != "" && c.Engine.Name != "WebKit"
			return c.Engine.Name, c.OS.Name, hit
		},
	},
	{
		id:       "impossible.trident_off_windows",
		severity: SeverityCritical,
		detail:   "Trident is a Windows-only engine",
		match: func(c *Client) (string, string, bool) {
			hit := c.Engine.Name == "Trident" && c.OS.Name != "" && c.OS.Family != "Windows"
			return c.Engine.Name, c.OS.Name, hit
		},
	},
	{
		id:       "impossible.two_platforms",
		severity: SeverityCritical,
		detail:   "the User-Agent names two mutually exclusive platforms",
		match: func(c *Client) (string, string, bool) {
			raw := c.Raw
			hit := containsAll(raw, "iPhone", "Android") ||
				containsAll(raw, "Windows NT", "Android") ||
				containsAll(raw, "iPhone", "Windows NT")
			return c.OS.Name, "", hit
		},
	},
	{
		id:       "impossible.bot_with_device",
		severity: SeverityMedium,
		detail:   "a crawler declared a specific consumer device model",
		match: func(c *Client) (string, string, bool) {
			hit := c.Bot != nil && c.Bot.Category == BotSearchEngine && c.Device.Name != ""
			if c.Bot == nil {
				return "", "", false
			}
			return c.Bot.Name, c.Device.Name, hit
		},
	},
	{
		id:       "anachronism.future_browser",
		severity: SeverityHigh,
		detail:   "the browser version is far beyond anything released",
		match: func(c *Client) (string, string, bool) {
			// A version this far ahead of the release cadence is a typo or an
			// edit, not a browser. The bound is deliberately generous so it
			// does not fire on a real release the database has not seen yet.
			hit := c.Browser.Version.Major > maxPlausibleMajor(c.Browser.Family)
			return c.Browser.Name + " " + c.Browser.Version.Raw, "", hit
		},
	},
}

// maxPlausibleMajor bounds a browser family's major version. It is intentionally
// loose: its job is to catch "Chrome/9000", not to know today's release.
func maxPlausibleMajor(family string) int {
	switch family {
	case "Chrome", "Firefox", "Edge", "Opera":
		return 400 // roughly a decade of releases beyond the present
	case "Safari":
		return 40
	case "IE":
		return 12
	default:
		return 1000
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	// strings.Contains, kept local so the check reads the same as the rules.
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// CheckConsistency evaluates a client against every impossibility rule and
// combines the result with any disagreement already found between the
// User-Agent and Client Hints.
func CheckConsistency(c *Client) Consistency {
	var cons Consistency
	cons.Findings = append(cons.Findings, c.hintFindings...)

	for _, r := range spoofRules {
		left, right, hit := r.match(c)
		if !hit {
			continue
		}
		cons.Findings = append(cons.Findings, Finding{
			Rule: r.id, Severity: r.severity, Detail: r.detail, Left: left, Right: right,
		})
	}

	// The score is the strongest single finding rather than a sum: one
	// impossibility is already conclusive, and adding up weak signals would
	// let ordinary clients drift into suspicion.
	for _, f := range cons.Findings {
		if w := f.Severity.weight(); w > cons.Score {
			cons.Score = w
		}
	}
	return cons
}

// Consistency returns the coherence check for this client.
func (c *Client) Consistency() Consistency { return CheckConsistency(c) }
