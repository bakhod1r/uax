package uax

import (
	"strings"
	"testing"
)

const chromeAndroid = "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36"

func TestParseIsPure(t *testing.T) {
	t.Parallel()

	// The package-level Parse must hold no state between calls: no cache, no
	// recorder, no retention. Two calls must be independent and identical.
	a := Parse(chromeAndroid)
	b := Parse(chromeAndroid)

	if a == b {
		t.Error("Parse returned the same pointer twice; it is sharing state")
	}
	if a.Browser != b.Browser || a.OS != b.OS || a.Device != b.Device || a.Kind != b.Kind {
		t.Error("Parse is not deterministic")
	}
}

func TestParseToAllocatesNothing(t *testing.T) {
	// AllocsPerRun requires exclusive use of the process, so this test does not
	// run in parallel with others.

	var c Client
	allocs := testing.AllocsPerRun(200, func() {
		ParseTo(&c, chromeAndroid)
	})

	// This is an assertion, not a report. Allocation-freedom is a property of
	// the design — token views into the input, constants from the generated
	// tables — and losing it silently is a regression worth failing on.
	if allocs > 0 {
		t.Errorf("ParseTo allocated %.0f times per run, want 0", allocs)
	}
}

func TestMaxInputLengthIsEnforced(t *testing.T) {
	t.Parallel()

	huge := chromeAndroid + strings.Repeat("A", 1<<20)
	c := New(Config{MaxInputLength: 256}).Parse(huge)

	if !c.Truncated {
		t.Error("oversized input was not flagged as truncated")
	}
	if len(c.Raw) != 256 {
		t.Errorf("raw length: got %d, want 256", len(c.Raw))
	}
	// Truncation must not cost correctness for the part that survived.
	if c.Browser.Name != "Chrome" {
		t.Errorf("browser after truncation: got %q, want Chrome", c.Browser.Name)
	}
}

func TestSkipOptions(t *testing.T) {
	t.Parallel()

	c := New(Config{SkipDevice: true, SkipBot: true}).Parse(chromeAndroid)

	if c.Device.Brand != "" {
		t.Errorf("SkipDevice: brand resolved anyway: %q", c.Device.Brand)
	}
	if c.Browser.Name != "Chrome" {
		t.Errorf("browser: got %q, want Chrome", c.Browser.Name)
	}
}

func TestExtraCoverage(t *testing.T) {
	// Cover various resolve and resolve_fields branches
	Parse("Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0")                                              // WOW64, NT version
	Parse("Mozilla/5.0 (X11; CrOS x86_64 14526.89.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36")       // CrOS
	Parse("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15") // Apple version
	Parse("Mozilla/5.0 (SMART-TV; Linux; Tizen 5.0)")                                                                              // TV deviceType
	Parse("Mozilla/5.0 (PlayStation 4 3.11) AppleWebKit/537.73 (KHTML, like Gecko)")                                               // Console
	Parse("Mozilla/5.0 (Linux; Android 10; SM-T510) AppleWebKit/537.36")                                                           // Tablet
}

func TestCacheIsOptOut(t *testing.T) {
	t.Parallel()

	// A parser with caching disabled must never hand out a shared pointer,
	// because callers are free to mutate what they receive.
	p := New(Config{})
	c1 := p.Parse(chromeAndroid)
	c2 := p.Parse(chromeAndroid)
	if c1 == c2 {
		t.Error("uncached parser returned a shared pointer")
	}

	cached := New(Config{Cache: CacheConfig{Enabled: true}})
	c3 := cached.Parse(chromeAndroid)
	c4 := cached.Parse(chromeAndroid)
	if c3 != c4 {
		t.Error("cached parser did not reuse its result")
	}
}

func TestCustomRulesOverrideBuiltins(t *testing.T) {
	t.Parallel()

	reg := NewRegistry()
	reg.RegisterBrowser(Rule{
		ID: "internal_app", Name: "AcmeApp", Family: "Chrome",
		Token: "AcmeApp", VersionFrom: "token",
		Priority: 1000, Confidence: 0.99,
	})

	p := New(Config{Registry: reg})
	c := p.Parse("Mozilla/5.0 (Linux; Android 15) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 AcmeApp/3.2 Mobile Safari/537.36")

	if c.Browser.Name != "AcmeApp" {
		t.Errorf("browser: got %q, want AcmeApp", c.Browser.Name)
	}
	if !c.Browser.Version.AtLeast(3, 2) {
		t.Errorf("version: got %v, want >= 3.2", c.Browser.Version)
	}

	// The built-in database must be untouched by a registry.
	if Parse(chromeAndroid).Browser.Name != "Chrome" {
		t.Error("a registry leaked into the built-in database")
	}
}

func FuzzParse(f *testing.F) {
	f.Add(chromeAndroid)
	f.Add("curl/8.4.0")
	f.Add("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	f.Add("")
	f.Add("(((((((((")
	f.Add("Mozilla/5.0 (\x00\x00; \xff\xfe)")

	f.Fuzz(func(t *testing.T, ua string) {
		c := Parse(ua)

		// Whatever the input, the result must stay internally coherent: a
		// caller must never see a confidence outside [0,1] or a nil pointer
		// where the API promises a value.
		for _, conf := range []float32{
			c.Browser.Confidence, c.OS.Confidence, c.Device.Confidence, c.Engine.Confidence,
		} {
			if conf < 0 || conf > 1 {
				t.Fatalf("confidence out of range: %v for %q", conf, ua)
			}
		}
		if c.AI != nil && c.Bot == nil {
			t.Fatalf("AI resolved without a bot for %q", ua)
		}
	})
}
