package uax

// Parse parses a User-Agent string.
//
// It is a pure function: no cache, no shared state, no retention of the input.
// For a configured parser — caching, recording, metrics, or skipping work you
// do not need — use New.
func Parse(raw string) *Client {
	return defaultParser.Parse(raw)
}

// ParseTo parses into a caller-supplied Client, allocating nothing. Use it on
// hot paths that can reuse a Client between requests; Parse is the ordinary
// entry point.
func ParseTo(dst *Client, raw string) {
	defaultParser.ParseTo(dst, raw)
}

// Parse parses a User-Agent string under the parser's configuration.
func (p *Parser) Parse(raw string) *Client {
	if p.cache != nil {
		if c, ok := p.cache.Get(raw); ok {
			p.metric(MetricsReporter.IncCacheHit)
			return c
		}
	}

	c := new(Client)
	p.ParseTo(c, raw)

	if p.cache != nil {
		p.cache.Set(raw, c)
	}
	return c
}

// ParseTo parses into dst, reusing its storage. dst is fully overwritten, so a
// stale Client can be handed back without clearing it first.
func (p *Parser) ParseTo(dst *Client, raw string) {
	*dst = Client{}
	p.metric(MetricsReporter.IncParse)

	if len(raw) > p.maxInput {
		raw = raw[:p.maxInput]
		dst.Truncated = true
	}
	dst.Raw = raw
	if raw == "" {
		dst.Kind = KindUnknown
		return
	}

	// One pass over the input. Everything below reads tokens, never bytes.
	var t tokens
	collect(&t, raw)

	// Every category is resolved in one sweep over the collected tokens.
	m := p.indexes().matchAll(&t)

	dst.Browser = resolveBrowser(&t, m[catBrowserIdx])
	dst.Engine = resolveEngine(&t, m[catEngineIdx], m[catBrowserIdx])
	dst.OS = resolveOS(&t, m[catOSIdx])

	if !p.cfg.SkipDevice {
		dst.Device = p.resolveDevice(&t)
	}

	if !p.cfg.SkipBot {
		dst.Bot = resolveBot(&t, m[catBotIdx], &dst.botStore)
		dst.AI = aiFor(dst.Bot, &dst.aiStore)
		dst.App = resolveApp(&t, m[catAppIdx], &dst.appStore)
		dst.Automation = resolveAutomation(&t, m[catAutomationIdx], &dst.automationStore)

		// A browser and an HTTP client are mutually exclusive claims. When a
		// real browser identified itself, a stray library-shaped token is noise.
		if dst.Browser.Name == "" {
			dst.HTTPClient = resolveHTTPClient(&t, m[catClientIdx], &dst.httpClientStore)
		}
	}

	applyBrowserMode(&t, dst)
	dst.Device.Type = deviceType(&t, dst, m[catOSIdx].rule)

	// After the form factor, because reduction removes the model but not the
	// evidence that the client is a phone.
	applyReduction(&t, dst)

	dst.Kind = kind(dst)

	if dst.Kind == KindBot || dst.Kind == KindAICrawler || dst.Kind == KindAIAgent {
		p.metric(MetricsReporter.IncBotDetected)
	}

	if p.cfg.Recorder != nil {
		p.record(&t, dst)
	}
}

// applyBrowserMode records how a browser is running rather than which browser
// it is. An Android WebView is Chrome by engine and by rules, but a caller
// asking what rendered the page deserves the more specific answer.
func applyBrowserMode(t *tokens, c *Client) {
	switch {
	case c.Automation != nil && c.Automation.Category == "Headless":
		c.Browser.Mode = "Headless"

	case t.hasWebViewFlag:
		c.Browser.Mode = "WebView"
		if c.Browser.Family == "Chrome" && c.OS.Family == "Android" {
			c.Browser.Name = "Android WebView"
		}

	case c.App != nil:
		// An in-app browser is a WebView even without the "wv" marker, which
		// iOS never sends.
		c.Browser.Mode = "WebView"
	}
}

// record reports product tokens that matched no rule, so a deployment can see
// what its own traffic contains that the database does not cover yet.
func (p *Parser) record(t *tokens, c *Client) {
	for i := 0; i < t.nprod; i++ {
		name := t.prodName(i)
		switch name {
		case "", "Mozilla", "AppleWebKit", "KHTML", "Gecko", "like", "Safari", "Version", "Mobile":
			continue
		case c.Browser.Name, c.OS.Name, c.Engine.Name:
			continue
		}
		if c.Bot != nil && name == c.Bot.Name {
			continue
		}
		if c.HTTPClient != nil && name == c.HTTPClient.Name {
			continue
		}
		p.cfg.Recorder.Record(name)
	}
}

// metric invokes a counter only when metrics are configured, keeping the
// no-metrics path free of interface calls.
func (p *Parser) metric(f func(MetricsReporter)) {
	if p.cfg.Metrics != nil {
		f(p.cfg.Metrics)
	}
}
