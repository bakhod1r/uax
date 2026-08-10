package uax

import "net/http"

// Client is everything uax could determine about one caller.
//
// Fields that describe an identity always present in a request — browser, OS,
// device, engine — are values. Fields that describe something that may simply
// not be there are pointers, so "no bot was detected" and "a bot with empty
// fields" can never be confused.
type Client struct {
	Browser Browser `json:"browser,omitempty"`
	OS      OS      `json:"os,omitempty"`
	Device  Device  `json:"device,omitempty"`
	Engine  Engine  `json:"engine,omitempty"`

	Bot        *Bot        `json:"bot,omitempty"`
	AI         *AI         `json:"ai,omitempty"`
	App        *App        `json:"app,omitempty"`
	Automation *Automation `json:"automation,omitempty"`
	HTTPClient *HTTPClient `json:"http_client,omitempty"`

	Hints         *Hints         `json:"client_hints,omitempty"`
	Fetch         *FetchContext  `json:"fetch,omitempty"`
	AcceptProfile *AcceptProfile `json:"accept_profile,omitempty"`

	Kind    Kind     `json:"kind,omitempty"`
	Verdict *Verdict `json:"verdict,omitempty"`

	// SpoofScore is 0 when nothing contradicts itself and rises towards 1 as
	// the client's claims stop fitting together. See Consistency for the
	// individual findings behind the number.
	SpoofScore float32 `json:"spoof_score,omitempty"`

	// hintFindings records disagreements discovered while folding Client Hints
	// into the User-Agent result. They are kept because merging would otherwise
	// destroy the evidence: once the hints have overwritten the string, the
	// contradiction is unrecoverable.
	hintFindings []Finding

	// Storage backing the optional pointers above. Keeping it inside the
	// Client lets ParseTo fill them without touching the heap: the pointer API
	// is what callers want (nil says "absent" unambiguously), and the storage
	// is what keeps it free.
	// fromRequest records whether this client was parsed from a full request.
	// Header-shape signals are only meaningful then: a bare User-Agent string
	// has no headers to be missing, and penalising it for their absence would
	// make every offline parse look suspicious.
	fromRequest bool

	botStore        Bot
	aiStore         AI
	appStore        App
	automationStore Automation
	httpClientStore HTTPClient

	// Raw is the input as parsed. It is the truncated form when Truncated is
	// set, so it always matches what the fields above were derived from.
	Raw       string `json:"raw,omitempty"`
	Truncated bool   `json:"truncated,omitempty"`
}

// ParseRequest parses a request's User-Agent together with its Client Hints,
// fetch metadata and Accept headers.
//
// Client Hints are authoritative where present: they are structured, and they
// survive the User-Agent reduction that is steadily hollowing out the legacy
// string. Where hints and the User-Agent disagree, both are kept and the
// disagreement itself is recorded as a spoofing signal.
func ParseRequest(r *http.Request) *Client {
	return defaultParser.ParseRequest(r)
}

// ParseRequest parses an http.Request under the parser's configuration.
func (p *Parser) ParseRequest(r *http.Request) *Client {
	c := p.Parse(r.UserAgent())

	// A cached Client is shared with other requests; per-request headers must
	// never be written into it.
	if p.cache != nil {
		clone := *c
		c = &clone
	}

	c.fromRequest = true
	c.Hints = ParseHints(r.Header)
	c.hintFindings = mergeHints(c)

	c.Fetch = ParseFetchContext(
		r.Header.Get("Sec-Fetch-Site"),
		r.Header.Get("Sec-Fetch-Mode"),
		r.Header.Get("Sec-Fetch-Dest"),
		r.Header.Get("Sec-Fetch-User"),
	)

	c.AcceptProfile = ParseAcceptProfile(
		r.Header.Get("Accept"),
		r.Header.Get("Accept-Language"),
		r.Header.Get("Accept-Encoding"),
	)

	// One coherence check feeds both. The verdict is derived from the findings,
	// so computing it separately would run every spoof rule a second time on
	// every request.
	cons := CheckConsistency(c)
	c.SpoofScore = cons.Score
	c.Verdict = VerdictFrom(c, cons)

	if p.cfg.Metrics != nil && c.SpoofScore >= SeverityCritical.weight() {
		p.cfg.Metrics.IncSpoofDetected()
	}
	return c
}

// ------------------------------------------------------------------ helpers
//
// Every predicate below requires positive evidence. None of them treats the
// absence of a signal as proof of its opposite: IsDesktop is false for an
// unrecognised client rather than true by elimination.

// IsMobile reports whether the client is a phone.
func (c *Client) IsMobile() bool {
	if c.Hints != nil && c.Hints.Mobile {
		return true
	}
	return c.Device.Type == DeviceMobile
}

// IsTablet reports whether the client is a tablet.
func (c *Client) IsTablet() bool { return c.Device.Type == DeviceTablet }

// IsDesktop reports whether the client is a desktop or laptop.
func (c *Client) IsDesktop() bool { return c.Device.Type == DeviceDesktop }

// IsTV reports whether the client is a television.
func (c *Client) IsTV() bool { return c.Device.Type == DeviceTV }

// IsConsole reports whether the client is a games console.
func (c *Client) IsConsole() bool { return c.Device.Type == DeviceConsole }

// IsBot reports whether the client identified itself as a crawler, agent or
// scraper. This is a claim by the client, not a verified fact.
func (c *Client) IsBot() bool { return c.Bot != nil }

// IsAICrawler reports whether an AI system is fetching, either to collect data
// or to answer a user's live request.
func (c *Client) IsAICrawler() bool { return c.AI != nil }

// IsWebView reports whether the page is rendered inside an embedded browser.
func (c *Client) IsWebView() bool { return c.Browser.Mode == "WebView" }

// IsHeadless reports whether the client declared itself headless.
func (c *Client) IsHeadless() bool {
	return c.Browser.Mode == "Headless" ||
		(c.Automation != nil && c.Automation.Category == "Headless")
}

// IsAutomated reports whether a browser-automation framework was declared.
func (c *Client) IsAutomated() bool {
	return c.Automation != nil && c.Automation.Category != "Embedded"
}

// IsCLI reports whether the client is an interactive command-line tool.
func (c *Client) IsCLI() bool {
	return c.HTTPClient != nil && c.HTTPClient.Category == "CLI"
}

// IsTool reports whether the client is any non-browser HTTP client.
func (c *Client) IsTool() bool { return c.HTTPClient != nil }

// IsWindows reports whether the client runs Windows.
func (c *Client) IsWindows() bool { return c.OS.Family == "Windows" }

// IsMacOS reports whether the client runs macOS.
func (c *Client) IsMacOS() bool { return c.OS.Family == "macOS" }

// IsLinux reports whether the client runs a Linux distribution. Android is a
// separate platform here, as callers almost always mean desktop Linux.
func (c *Client) IsLinux() bool { return c.OS.Family == "Linux" }

// IsAndroid reports whether the client runs Android, including Android TV.
func (c *Client) IsAndroid() bool { return c.OS.Family == "Android" }

// IsIOS reports whether the client runs iOS or iPadOS.
func (c *Client) IsIOS() bool { return c.OS.Family == "iOS" }
