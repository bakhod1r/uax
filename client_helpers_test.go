package uax

import "testing"

// The predicates answer questions about a parsed client, so they are tested
// through real User-Agents rather than hand-built structs. A hand-built struct
// can assert a shape the parser would never produce.
func TestClientPredicates(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ua   string
		want func(*Client) bool
		desc string
	}{
		{"windows desktop", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
			(*Client).IsDesktop, "IsDesktop"},
		{"android phone", chromeAndroid, (*Client).IsMobile, "IsMobile"},
		{"ipad", "Mozilla/5.0 (iPad; CPU OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
			(*Client).IsTablet, "IsTablet"},
		{"googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			(*Client).IsBot, "IsBot"},
		{"gptbot", "Mozilla/5.0 (compatible; GPTBot/1.2; +https://openai.com/gptbot)",
			(*Client).IsAICrawler, "IsAICrawler"},
		{"curl", "curl/8.4.0", (*Client).IsCLI, "IsCLI"},
		{"headless", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/120.0.0.0 Safari/537.36",
			(*Client).IsHeadless, "IsHeadless"},
		{"android webview", "Mozilla/5.0 (Linux; Android 14; SM-A546E Build/UP1A.231005.007; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.0.0 Mobile Safari/537.36",
			(*Client).IsWebView, "IsWebView"},
		{"smart tv", "Mozilla/5.0 (SMART-TV; LINUX; Tizen 6.0) AppleWebKit/537.36 (KHTML, like Gecko) Version/6.0 TV Safari/537.36",
			(*Client).IsTV, "IsTV"},
		{"playstation", "Mozilla/5.0 (PlayStation; PlayStation 5/2.26) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Safari/605.1.15",
			(*Client).IsConsole, "IsConsole"},
		{"playwright", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Playwright/1.40",
			(*Client).IsAutomated, "IsAutomated"},
		{"tool", "Go-http-client/1.1", (*Client).IsTool, "IsTool"},
		{"windows", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", (*Client).IsWindows, "IsWindows"},
		{"macos", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", (*Client).IsMacOS, "IsMacOS"},
		{"linux", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", (*Client).IsLinux, "IsLinux"},
		{"android", chromeAndroid, (*Client).IsAndroid, "IsAndroid"},
		{"ios", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1", (*Client).IsIOS, "IsIOS"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if c := Parse(tc.ua); !tc.want(c) {
				t.Errorf("%s() = false, want true (device type %q, os %q)", tc.desc, c.Device.Type, c.OS.Name)
			}
		})
	}
}

// Client Hints outrank the User-Agent, because they are the signal that
// survives User-Agent reduction.
func TestHintsOverrideUserAgentForMobile(t *testing.T) {
	t.Parallel()

	c := &Client{Hints: &Hints{Mobile: true}}
	if !c.IsMobile() {
		t.Error("Sec-CH-UA-Mobile: ?1 was ignored")
	}
}

// The predicates must never answer by elimination: an unrecognised client is
// unknown, not a desktop.
func TestPredicatesRequirePositiveEvidence(t *testing.T) {
	t.Parallel()

	c := Parse("something-nobody-has-ever-seen/1.0")
	if c.IsDesktop() || c.IsMobile() || c.IsTablet() {
		t.Errorf("an unrecognised client claimed a form factor: %q", c.Device.Type)
	}
}
