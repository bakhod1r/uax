package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// What these tests hold the tools to.
//
// The output is prose a model reads, so asserting on its exact wording would
// break on every rephrasing while proving nothing. What matters is that the
// facts are present, that a fact the server does not have is absent rather than
// invented, and that "I could not tell" is reported as a tool error the model
// can act on rather than as a confident empty answer.

const (
	chromeWindows = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"
	chromeAndroid = "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36"
	safariIOS     = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1"
	googlebot     = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	curlAgent     = "curl/8.7.1"
	// Safari does not run on Windows and has not since 2012, so this string
	// cannot describe a real client.
	impossible = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15"
)

// contains fails unless every wanted fragment appears in the output.
func contains(t *testing.T, got string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("output is missing %q\n--- output ---\n%s", w, got)
		}
	}
}

// jsonArg builds the arguments object for a User-Agent tool. It encodes through
// encoding/json rather than quoting by hand: several of the inputs below are
// deliberately hostile — NUL bytes, invalid UTF-8 — and hand-quoting them would
// produce a malformed request, testing the JSON parser instead of the server.
func jsonArg(ua string) string {
	b, err := json.Marshal(map[string]string{"user_agent": ua})
	if err != nil {
		panic(err)
	}
	return string(b)
}

// --------------------------------------------------------- parse_user_agent

func TestParseUserAgentReportsWhatTheStringEstablishes(t *testing.T) {
	for _, tc := range []struct {
		name string
		ua   string
		want []string
	}{
		{"desktop browser", chromeWindows, []string{"browser: Chrome 138", "os: Windows", "confidence"}},
		{"android handset", chromeAndroid, []string{"browser: Chrome", "os: Android 15", "device:", "SM-S928B"}},
		{"ios browser", safariIOS, []string{"browser: Safari", "os: iOS"}},
		{"named crawler", googlebot, []string{"bot: Googlebot"}},
		{"http library", curlAgent, []string{"http client: curl"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := callTool(t, "parse_user_agent", jsonArg(tc.ua))
			if res.IsError {
				t.Fatalf("unexpected tool error: %s", text(res))
			}
			contains(t, text(res), tc.want...)
		})
	}
}

func TestParseUserAgentNamesTheHandsetBehindAModelCode(t *testing.T) {
	// The point of wiring the catalogue in: SM-S928B contains no hint that it
	// means Galaxy S24 Ultra, and a model cannot derive it.
	res := callTool(t, "parse_user_agent", jsonArg(chromeAndroid))
	contains(t, text(res), "brand: Samsung", "model: SM-S928B", "name: Galaxy S24 Ultra")
}

func TestParseUserAgentSaysSoWhenNothingResolves(t *testing.T) {
	// An unrecognised string must come back as a tool error, not as a page of
	// empty fields a model would read as a confident negative.
	res := callTool(t, "parse_user_agent", jsonArg("qqqq zzzz"))
	if !res.IsError {
		t.Fatalf("nonsense resolved to a result:\n%s", text(res))
	}
	// The message has to point somewhere: "no" and "no, and here is why" differ
	// in whether the model retries usefully.
	contains(t, text(res), "database_info")
}

func TestParseUserAgentSurvivesHostileInput(t *testing.T) {
	// The server sits in front of a parser whose whole premise is that its
	// input is chosen by an unauthenticated remote client.
	for _, ua := range []string{
		strings.Repeat("(", 5000),
		strings.Repeat("Mozilla/5.0 ", 500),
		"\x00\x01\x02 binary",
		"Mozilla/5.0 (\xff\xfe invalid utf8)",
		strings.Repeat("A", 100000),
		"Mozilla/5.0 (Linux; Android 15; " + strings.Repeat("X", 4000) + ")",
	} {
		// No assertion on the answer: the requirement is that there is one, and
		// that the process is still alive to give it.
		callTool(t, "parse_user_agent", jsonArg(ua))
	}
}

// -------------------------------------------------------- check_consistency

func TestCheckConsistencyReportsAContradiction(t *testing.T) {
	res := callTool(t, "check_consistency", jsonArg(impossible))
	got := text(res)
	if strings.Contains(got, "No contradictions") {
		t.Fatalf("Safari on Windows passed as consistent:\n%s", got)
	}
	contains(t, got, "contradiction", "trust:", "risk:", "recommend:")
}

func TestCheckConsistencyPassesAnOrdinaryClient(t *testing.T) {
	got := text(callTool(t, "check_consistency", jsonArg(chromeWindows)))
	contains(t, got, "No contradictions")
	// Consistency is not honesty, and a model that reads it as honesty will
	// authorise on it. The output has to say so.
	contains(t, got, "not proof")
}

func TestCheckConsistencyDoesNotTreatABotAsAThreat(t *testing.T) {
	// A crawler that names itself is the honest case, not the suspicious one.
	// A verdict that punishes self-identification teaches clients to lie.
	got := text(callTool(t, "check_consistency", jsonArg(googlebot)))
	contains(t, got, "trust:", "risk:")
	if strings.Contains(got, "contradiction(s) found") {
		t.Errorf("Googlebot flagged as contradictory:\n%s", got)
	}
}

func TestCheckConsistencyCallsTheVerdictAdvice(t *testing.T) {
	// Verdict.Recommend is advice; a model that reads it as an authorisation
	// decision builds an access control system out of a parsed string.
	got := text(callTool(t, "check_consistency", jsonArg(impossible)))
	contains(t, got, "advice")
}

// ------------------------------------------------------------ lookup_device

func TestLookupDeviceNamesAKnownHandset(t *testing.T) {
	got := text(callTool(t, "lookup_device", `{"code":"SM-S928B"}`))
	contains(t, got, "Samsung", "Galaxy S24 Ultra")
}

func TestLookupDeviceStripsTheBuildSuffix(t *testing.T) {
	// Android appends "Build/..." inside the platform comment, so the code a
	// caller copies out of a User-Agent usually arrives with it attached.
	withBuild := text(callTool(t, "lookup_device", `{"code":"SM-S928B Build/UP1A.231005.007"}`))
	contains(t, withBuild, "Samsung", "Galaxy S24 Ultra")
}

func TestLookupDeviceGivesTheBrandWithoutInventingAName(t *testing.T) {
	// A code whose shape identifies its manufacturer but which the catalogue
	// has never seen. The brand is derivable; the marketing name is not, and a
	// plausible-looking wrong name is worse than none.
	res := callTool(t, "lookup_device", `{"code":"SM-ZZ999Q"}`)
	got := text(res)
	if res.IsError {
		t.Fatalf("an unseen Samsung code should still resolve its brand:\n%s", got)
	}
	contains(t, got, "Samsung", "not in the catalogue")
}

func TestLookupDeviceRefusesACodeItCannotPlace(t *testing.T) {
	res := callTool(t, "lookup_device", `{"code":"ZZ-NONSENSE-9999"}`)
	if !res.IsError {
		t.Fatalf("nonsense code produced a confident answer:\n%s", text(res))
	}
}

// ------------------------------------------------------------ database_info

func TestDatabaseInfoReportsWhatProducedAClassification(t *testing.T) {
	// The fingerprint is what makes a stored classification reproducible: rules
	// change between releases, and a result without one cannot be re-checked.
	got := text(callTool(t, "database_info", `{}`))
	contains(t, got, "total", "uax release:", "rule fingerprint:", "device catalogue:", "catalogue source:")
}

func TestDatabaseInfoTakesNoArguments(t *testing.T) {
	// A model will send an empty object, omit arguments entirely, or send
	// something irrelevant. None of those is a reason to fail.
	for _, args := range []string{`{}`, `null`, `{"unexpected":"value"}`} {
		if res := callTool(t, "database_info", args); res.IsError {
			t.Errorf("arguments %s produced an error: %s", args, text(res))
		}
	}
}
