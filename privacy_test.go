package uax

import "testing"

func TestPrivacyRatesDisclosure(t *testing.T) {
	t.Parallel()

	// A phone that names its exact model, its full browser build and its OS
	// patch level has volunteered a great deal.
	phone := Parse("Mozilla/5.0 (Linux; Android 15.1.2; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.7204.92 Mobile Safari/537.36")
	if got := phone.Privacy().Risk; got != PrivacyHigh {
		t.Errorf("risk for an exact model: got %q, want %q (%+v)", got, PrivacyHigh, phone.Privacy().Identifiers)
	}

	// A reduced modern User-Agent discloses very little, which is the point of
	// the reduction.
	reduced := Parse("Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36")
	if got := reduced.Privacy().Risk; got == PrivacyHigh {
		t.Errorf("risk for a reduced User-Agent: got %q, want low or medium (%+v)", got, reduced.Privacy().Identifiers)
	}
}

func TestPrivacyNamesItsReasons(t *testing.T) {
	t.Parallel()

	p := Parse("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36").Privacy()
	if len(p.Identifiers) == 0 {
		t.Error("a risk rating was reported with no reasons behind it")
	}
}

func TestCoarseStripsIdentifiers(t *testing.T) {
	t.Parallel()

	full := parseWithDevices("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.7204.92 Mobile Safari/537.36")
	c := full.Coarse()

	if c.Device.Model != "" || c.Device.Name != "" {
		t.Errorf("device identity survived Coarse: %q / %q", c.Device.Model, c.Device.Name)
	}
	if c.Browser.Version.Patch != 0 || c.OS.Arch != "" {
		t.Error("narrowing version or architecture detail survived Coarse")
	}
	if c.Raw != "" {
		t.Error("the raw User-Agent survived Coarse")
	}

	// What survives must still answer the question analytics actually asks.
	if c.Browser.Name != "Chrome" || c.Browser.Version.Major != 138 || c.OS.Name != "Android" {
		t.Errorf("Coarse destroyed the useful part: %+v / %+v", c.Browser, c.OS)
	}
	if c.Device.Brand != "Samsung" || c.Device.Type != DeviceMobile {
		t.Errorf("Coarse destroyed the coarse device facts: %+v", c.Device)
	}

	// The original must be untouched.
	if full.Device.Model != "SM-S928B" {
		t.Error("Coarse mutated its receiver")
	}
}

func TestCoarseDoesNotAliasTheOriginal(t *testing.T) {
	t.Parallel()

	full := Parse("Mozilla/5.0 (compatible; GPTBot/1.2; +https://openai.com/gptbot)")
	c := full.Coarse()

	if c.Bot == full.Bot {
		t.Fatal("the coarse copy shares the original's bot")
	}
	c.Bot.Name = "mutated"
	if full.Bot.Name != "GPTBot" {
		t.Error("mutating the coarse copy changed the original")
	}
}

func TestPrivacyEdgeCases(t *testing.T) {
	// 4+ identifiers triggers High risk
	req := request("Mozilla/5.0 (Linux; Android 15.1; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.7204.92 Mobile Safari/537.36 Instagram 1.2.3",
		HeaderCHUAModel, `"Model123"`,
	)
	c := ParseRequest(req)
	c.Device.Name = "Galaxy S24" // manually inject to increase identifiers

	p := c.Privacy()
	if p.Risk != PrivacyHigh {
		t.Error("expected PrivacyHigh for >= 4 identifiers")
	}
}

func TestCoarseEdgeCases(t *testing.T) {
	var c *Client
	if c.Coarse() != nil {
		t.Error("expected Coarse() on nil client to return nil")
	}

	c = &Client{
		AI:         &AI{Name: "AI"},
		Automation: &Automation{Name: "Auto"},
		HTTPClient: &HTTPClient{Name: "HTTP"},
		Hints:      &Hints{Model: "abc"},
	}
	c2 := c.Coarse()
	if c2.AI == c.AI || c2.Automation == c.Automation || c2.HTTPClient == c.HTTPClient || c2.Hints == c.Hints {
		t.Error("expected deep copy for AI, Automation, HTTPClient, Hints")
	}
}

func TestBotCategories(t *testing.T) {
	b1 := &Bot{Category: BotSocialCrawler}
	if !b1.IsSocialCrawler() {
		t.Error("expected IsSocialCrawler")
	}
	b2 := &Bot{Category: BotSEO}
	if !b2.IsSEO() {
		t.Error("expected IsSEO")
	}
}
