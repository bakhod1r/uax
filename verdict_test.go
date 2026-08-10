package uax

import "testing"

func TestVerdictSeparatesTrustFromRisk(t *testing.T) {
	t.Parallel()

	// An honest, well-known crawler is highly trustworthy and low risk.
	// "Bot" is not a synonym for "threat", and the two scores must show that.
	bot := ParseRequest(request("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"))
	if bot.Verdict.Kind != KindBot {
		t.Errorf("kind: got %q, want %q", bot.Verdict.Kind, KindBot)
	}
	if bot.Verdict.Trust < 0.7 {
		t.Errorf("trust for a self-identified Googlebot: got %v, want >= 0.7", bot.Verdict.Trust)
	}
	if bot.Verdict.Risk > 0.25 {
		t.Errorf("risk for a self-identified Googlebot: got %v, want <= 0.25", bot.Verdict.Risk)
	}

	// A contradictory client is the opposite: whatever it claims, it is risky.
	spoofed := ParseRequest(request(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"))
	if spoofed.Verdict.Risk < 1.0 {
		t.Errorf("risk for Safari on Windows: got %v, want 1.0", spoofed.Verdict.Risk)
	}
	if spoofed.Verdict.Recommend != ActionBlock {
		t.Errorf("recommendation: got %q, want %q", spoofed.Verdict.Recommend, ActionBlock)
	}
}

func TestVerdictReasonsAreStableCodes(t *testing.T) {
	t.Parallel()

	v := ParseRequest(request(
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")).Verdict

	if !v.Has("impossible.safari_on_windows") {
		t.Errorf("expected a stable reason code; got %+v", v.Reasons)
	}
	for _, r := range v.Reasons {
		if r.Code == "" {
			t.Error("a reason was emitted without a code")
		}
	}
}

// Ordinary browsing must not be treated as an attack. False positives here cost
// real users their access, so the bar for anything above Allow is high.
func TestOrdinaryBrowsersAreAllowed(t *testing.T) {
	t.Parallel()

	r := request(chromeWindows,
		HeaderCHUA, `"Not/A)Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
		HeaderCHUAPlatform, `"Windows"`,
		HeaderCHUAMobile, "?0",
		"Accept-Language", "en-US,en;q=0.9",
		"Accept", "text/html,application/xhtml+xml",
	)

	v := ParseRequest(r).Verdict
	if v.Recommend != ActionAllow {
		t.Errorf("recommendation: got %q, want %q (reasons: %+v)", v.Recommend, ActionAllow, v.Reasons)
	}
	if v.Trust < 0.9 {
		t.Errorf("trust for a corroborated browser: got %v, want >= 0.9", v.Trust)
	}
}

func TestUnattributedBotIsRateLimited(t *testing.T) {
	t.Parallel()

	v := ParseRequest(request("SomeRandomBot/1.0 (+http://example.com/bot)")).Verdict
	if !v.Has("bot.unattributed") {
		t.Errorf("expected bot.unattributed; got %+v", v.Reasons)
	}
	if v.Recommend == ActionAllow {
		t.Error("an unattributed self-declared bot was allowed unconditionally")
	}
}

// Claimed identity is never reported as verified identity.
func TestBotIdentityIsMarkedUnverified(t *testing.T) {
	t.Parallel()

	v := ParseRequest(request("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")).Verdict
	if !v.Has("bot.identity_unverified") {
		t.Errorf("a claimed Googlebot was not marked unverified: %+v", v.Reasons)
	}
}
