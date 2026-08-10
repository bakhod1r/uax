package uax

// Action is what uax suggests doing with a client.
//
// It is advice, never a decision. This package sees one request's headers; it
// cannot see rate, reputation, session history or business context, and any of
// those may reasonably overrule it. The caller decides.
type Action string

const (
	ActionAllow     Action = "allow"
	ActionChallenge Action = "challenge"
	ActionRateLimit Action = "rate_limit"
	ActionBlock     Action = "block"
)

// Reason is one weighted input to a verdict, carrying a stable code so it can
// be logged, alerted on and asserted in tests.
type Reason struct {
	Code   string  `json:"code"`
	Weight float32 `json:"weight"`
	Detail string  `json:"detail,omitempty"`
}

// Verdict is the actionable summary of everything uax determined.
//
// The two scores answer different questions and should not be conflated:
//
//	Trust is how much the client's self-description can be believed. A verified
//	browser sending coherent hints scores high. An honest, well-behaved crawler
//	also scores high — being a bot is not being a liar.
//
//	Risk is how likely this request is to be hostile. A client whose claims
//	contradict each other scores high, whatever it claims to be.
type Verdict struct {
	Kind      Kind     `json:"kind"`
	Trust     float32  `json:"trust"`
	Risk      float32  `json:"risk"`
	Reasons   []Reason `json:"reasons,omitempty"`
	Recommend Action   `json:"recommend"`
}

// ComputeVerdict derives a verdict from an already-parsed client.
func ComputeVerdict(c *Client) *Verdict {
	return VerdictFrom(c, CheckConsistency(c))
}

// VerdictFrom derives a verdict from a client and a coherence check already
// performed on it.
//
// A verdict is built on top of the consistency findings, so ComputeVerdict has
// to run the spoof rules to produce one. Callers that display the findings as
// well — the CLI, the MCP server, ParseRequest — would otherwise run them
// twice per request. Pass a Consistency from any source but this client and the
// verdict describes a request that never arrived.
func VerdictFrom(c *Client, cons Consistency) *Verdict {
	v := &Verdict{
		Kind:      c.Kind,
		Trust:     0.5, // neither believed nor disbelieved until evidence arrives
		Risk:      0,
		Recommend: ActionAllow,
	}

	// ---- coherence ---------------------------------------------------------
	for _, f := range cons.Findings {
		v.Reasons = append(v.Reasons, Reason{
			Code: f.Rule, Weight: f.Severity.weight(), Detail: f.Detail,
		})
	}
	if cons.Score > 0 {
		v.Risk = cons.Score
		v.Trust = clamp01(v.Trust - cons.Score)
	}

	// ---- what the client says it is ---------------------------------------
	switch {
	case c.Bot != nil:
		// A crawler that names itself is being cooperative. Its trust reflects
		// how specific the identification was: a named, attributed bot is far
		// more credible than something that merely contains "bot".
		v.Trust = clamp01(v.Trust + (c.Bot.Confidence-0.5)*0.8)
		v.addReason("client.self_identified_bot", 0.1, "the client identified itself as "+c.Bot.Name)

		if c.Bot.Category == BotGeneric {
			v.addReason("bot.unattributed", 0.2, "self-declared bot with no known operator")
			v.Risk = maxF(v.Risk, 0.3)
		}
		// Identity is a claim until verified. The uax/verify module turns a
		// claimed Googlebot into a verified one using DNS and address ranges.
		if c.Bot.IsSearchEngine() || c.Bot.IsAICrawler() {
			v.addReason("bot.identity_unverified", 0, "identity is claimed, not verified; see uax/verify")
		}

	case c.Automation != nil && c.Automation.Category != "Embedded":
		v.addReason("client.automation_declared", 0.3, c.Automation.Name+" declared itself")
		v.Risk = maxF(v.Risk, 0.35)
		v.Trust = clamp01(v.Trust + 0.2) // declaring it at all is honest

	case c.HTTPClient != nil:
		v.addReason("client.http_library", 0.1, c.HTTPClient.Name+" is not a browser")
		v.Trust = clamp01(v.Trust + 0.3)

	case c.Browser.Name != "":
		v.Trust = clamp01(v.Trust + 0.2)
		if c.Browser.Source == SourceBoth {
			// The User-Agent and the hints independently agree. This is the
			// strongest ordinary evidence a request can carry.
			v.addReason("hints.corroborate_ua", 0, "Client Hints agree with the User-Agent")
			v.Trust = clamp01(v.Trust + 0.25)
		}

	default:
		v.addReason("client.unrecognised", 0.2, "nothing in the request identified the client")
		v.Risk = maxF(v.Risk, 0.25)
		v.Trust = clamp01(v.Trust - 0.2)
	}

	// ---- header shape ------------------------------------------------------
	// A browser is a large, chatty program: it sends a language preference and
	// fetch metadata. A script pretending to be one usually sends neither.
	// This is weak evidence on its own, which is why it moves risk a little and
	// never decides anything by itself.
	if c.fromRequest && c.Browser.Name != "" && c.Bot == nil {
		if c.AcceptProfile != nil && c.AcceptProfile.IsMissingLanguage() {
			v.addReason("headers.no_accept_language", 0.15, "a browser without Accept-Language is unusual")
			v.Risk = maxF(v.Risk, 0.2)
		}
		if c.Hints == nil && c.Browser.Family == "Chrome" {
			// Every modern Chromium sends Sec-CH-UA unprompted.
			v.addReason("headers.chromium_without_hints", 0.2, "Chromium sent no Sec-CH-UA")
			v.Risk = maxF(v.Risk, 0.25)
		}
	}

	if c.Truncated {
		v.addReason("input.truncated", 0.1, "the User-Agent exceeded the configured length limit")
		v.Risk = maxF(v.Risk, 0.2)
	}

	v.Recommend = recommend(v)
	return v
}

// recommend maps the scores onto advice.
//
// The thresholds are deliberately reluctant. A false block costs a real user
// their access, while a false allow costs one served request, so nothing short
// of an outright impossibility reaches ActionBlock.
func recommend(v *Verdict) Action {
	switch {
	case v.Risk >= 1.0:
		return ActionBlock
	case v.Risk >= 0.5:
		return ActionChallenge
	case v.Kind == KindBot && v.Trust < 0.5:
		return ActionRateLimit
	case v.Risk >= 0.25:
		return ActionRateLimit
	default:
		return ActionAllow
	}
}

func (v *Verdict) addReason(code string, weight float32, detail string) {
	v.Reasons = append(v.Reasons, Reason{Code: code, Weight: weight, Detail: detail})
}

// Has reports whether a reason code is present.
func (v *Verdict) Has(code string) bool {
	if v == nil {
		return false
	}
	for _, r := range v.Reasons {
		if r.Code == code {
			return true
		}
	}
	return false
}

// IsSuspicious reports whether anything about the client failed to add up.
func (c *Client) IsSuspicious() bool {
	return c.Verdict != nil && c.Verdict.Risk >= 0.25
}

func clamp01(f float32) float32 {
	switch {
	case f < 0:
		return 0
	case f > 1:
		return 1
	}
	return f
}

func maxF(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
