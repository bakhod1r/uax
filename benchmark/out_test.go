package main

import (
	"testing"

	"github.com/avct/uasurfer"
	"github.com/bakhod1r/uax"
	"github.com/mileusna/useragent"
	"github.com/mssola/user_agent"
)

// TestWhatEachResolves prints what every library extracts from the same input.
// A latency comparison is meaningless without it: a parser that resolves less
// is expected to be faster, and the only fair reading is cost against work.
func TestWhatEachResolves(t *testing.T) {
	c := uax.Parse(ua)
	t.Logf("uax        browser=%s %s | os=%s %s | device=%s %s %q | type=%s",
		c.Browser.Name, c.Browser.Version.Raw, c.OS.Name, c.OS.Version.Raw,
		c.Device.Brand, c.Device.Model, c.Device.Name, c.Device.Type)

	m := useragent.Parse(ua)
	t.Logf("mileusna   browser=%s %s | os=%s %s | device=%q | mobile=%v",
		m.Name, m.Version, m.OS, m.OSVersion, m.Device, m.Mobile)

	s := user_agent.New(ua)
	bn, bv := s.Browser()
	t.Logf("mssola     browser=%s %s | os=%s | platform=%s | mobile=%v",
		bn, bv, s.OS(), s.Platform(), s.Mobile())

	u := uasurfer.Parse(ua)
	t.Logf("uasurfer   browser=%v %d | os=%v %d | device=%v",
		u.Browser.Name, u.Browser.Version.Major, u.OS.Name, u.OS.Version.Major, u.DeviceType)

	p := uapParser.Parse(ua)
	t.Logf("uap-go     browser=%s %s | os=%s %s | device=%s %s %s",
		p.UserAgent.Family, p.UserAgent.Major, p.Os.Family, p.Os.Major,
		p.Device.Brand, p.Device.Model, p.Device.Family)
}

// TestTheHardCases is the other half of the comparison.
//
// On an ordinary Chrome string every library above agrees, which makes the
// latency table look like the whole story. It is not. The strings below are the
// ones that decide whether a result can be acted on: a browser string that is
// really a crawler, a headless browser, an in-app WebView, a library that is
// not a browser at all, and a string Chrome froze on purpose.
//
// Nothing here is scored. Each library is asked what it sees and the answers
// are printed side by side, because "wrong" depends on what the caller needed —
// a library that never claimed to identify bots is not failing at it.
func TestTheHardCases(t *testing.T) {
	cases := []struct {
		name string
		ua   string
		want string
	}{
		{
			"googlebot",
			"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			"a crawler, not a browser — it says Mozilla like everything else",
		},
		{
			"chrome_headless",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/138.0.0.0 Safari/537.36",
			"automation — identical to Chrome but for one token",
		},
		{
			"instagram_webview",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 320.0.0.0",
			"an in-app WebView — the traffic is Instagram's, not Safari's",
		},
		{
			"curl",
			"curl/8.7.1",
			"an HTTP library — no browser, no OS, and saying otherwise is invention",
		},
		{
			"reduced_android",
			"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36",
			"frozen by User-Agent Reduction — 'Android 10' and 'K' are placeholders, not facts",
		},
		{
			"claudebot",
			"Mozilla/5.0 (compatible; ClaudeBot/1.0; +claudebot@anthropic.com)",
			"an AI crawler — a category that did not exist when most of these were written",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("want: %s", tc.want)

			c := uax.Parse(tc.ua)
			line := "uax        browser=" + or(c.Browser.Name) + " | os=" + or(c.OS.Name)
			if c.Bot != nil {
				line += " | BOT=" + c.Bot.Name + " (" + string(c.Bot.Category) + ")"
			}
			if c.Automation != nil {
				line += " | AUTOMATION=" + c.Automation.Name
			}
			if c.App != nil {
				line += " | WEBVIEW=" + c.App.Name
			}
			if c.HTTPClient != nil {
				line += " | HTTP_CLIENT=" + c.HTTPClient.Name
			}
			if c.OS.Reduced {
				line += " | REDUCED (version withheld: it is a placeholder)"
			}
			t.Log(line + " | kind=" + string(c.Kind))

			m := useragent.Parse(tc.ua)
			t.Logf("mileusna   browser=%s | os=%s | bot=%v", or(m.Name), or(m.OS), m.Bot)

			s := user_agent.New(tc.ua)
			bn, _ := s.Browser()
			t.Logf("mssola     browser=%s | os=%s | bot=%v", or(bn), or(s.OS()), s.Bot())

			u := uasurfer.Parse(tc.ua)
			t.Logf("uasurfer   browser=%v | os=%v | device=%v", u.Browser.Name, u.OS.Name, u.DeviceType)

			p := uapParser.Parse(tc.ua)
			t.Logf("uap-go     browser=%s | os=%s | device=%s",
				or(p.UserAgent.Family), or(p.Os.Family), or(p.Device.Family))
		})
	}
}

// or marks an empty answer as empty. Printing nothing there reads as a column
// that failed to align rather than a parser that resolved nothing, and the
// difference between those is the whole point of the table.
func or(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}
