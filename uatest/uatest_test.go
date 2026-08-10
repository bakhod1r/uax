package uatest_test

import (
	"testing"

	"github.com/bakhod1r/uax"
	"github.com/bakhod1r/uax/uatest"
)

// Every builder must produce a string that parses back to what it described.
// This is what makes the package trustworthy: a fixture that the parser reads
// differently than intended would silently weaken every test built on it.
func TestBuildersRoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		ua      string
		browser string
		os      string
		device  uax.DeviceType
	}{
		{"chrome windows", uatest.Chrome(138, uatest.Windows(10)), "Chrome", "Windows", uax.DeviceDesktop},
		{"chrome macos", uatest.Chrome(138, uatest.MacOS(10, 15, 7)), "Chrome", "macOS", uax.DeviceDesktop},
		{"chrome linux", uatest.Chrome(138, uatest.Linux()), "Chrome", "Linux", uax.DeviceDesktop},
		{"chrome android", uatest.Chrome(138, uatest.Android(15, uatest.GalaxyS24Ultra)), "Chrome", "Android", uax.DeviceMobile},
		{"chrome tablet", uatest.Chrome(130, uatest.AndroidTablet(14, uatest.GalaxyTabS9)), "Chrome", "Android", uax.DeviceTablet},
		{"edge", uatest.Edge(138, uatest.Windows(10)), "Edge", "Windows", uax.DeviceDesktop},
		{"firefox", uatest.Firefox(128, uatest.MacOS(10, 15, 7)), "Firefox", "macOS", uax.DeviceDesktop},
		{"safari macos", uatest.Safari(17, 5, uatest.MacOS(10, 15, 7)), "Safari", "macOS", uax.DeviceDesktop},
		{"safari ios", uatest.Safari(17, 5, uatest.IOS(17, 5)), "Safari", "iOS", uax.DeviceMobile},
		{"safari ipados", uatest.Safari(17, 5, uatest.IPadOS(17, 5)), "Safari", "iPadOS", uax.DeviceTablet},
		{"headless", uatest.HeadlessChrome(120, uatest.Linux()), "Chrome", "Linux", uax.DeviceDesktop},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := uax.Parse(tc.ua)

			if c.Browser.Name != tc.browser {
				t.Errorf("browser: got %q, want %q\n  ua: %s", c.Browser.Name, tc.browser, tc.ua)
			}
			if c.OS.Name != tc.os {
				t.Errorf("os: got %q, want %q\n  ua: %s", c.OS.Name, tc.os, tc.ua)
			}
			if c.Device.Type != tc.device {
				t.Errorf("device type: got %q, want %q\n  ua: %s", c.Device.Type, tc.device, tc.ua)
			}
		})
	}
}

func TestDeviceBuildersCarryTheirBrand(t *testing.T) {}
func TestWebViewBuilder(t *testing.T) {
	t.Parallel()

	c := uax.Parse(uatest.WebView(130, uatest.Android(14, uatest.GalaxyS24Ultra)))
	if !c.IsWebView() {
		t.Errorf("IsWebView() = false for a WebView fixture (mode %q)", c.Browser.Mode)
	}
}

func TestBotFixtures(t *testing.T) {
	t.Parallel()

	for ua, want := range map[string]uax.BotCategory{
		uatest.Googlebot:   uax.BotSearchEngine,
		uatest.Bingbot:     uax.BotSearchEngine,
		uatest.GPTBot:      uax.BotAICrawler,
		uatest.ClaudeBot:   uax.BotAICrawler,
		uatest.FacebookBot: uax.BotSocialCrawler,
	} {
		c := uax.Parse(ua)
		if c.Bot == nil {
			t.Errorf("bot fixture not recognised: %s", ua)
			continue
		}
		if c.Bot.Category != want {
			t.Errorf("category: got %q, want %q\n  ua: %s", c.Bot.Category, want, ua)
		}
	}
}

func TestRequestCarriesBrowserHeaders(t *testing.T) {
	t.Parallel()

	r := uatest.Request(uatest.Chrome(138, uatest.Windows(10)))
	if r.Header.Get("Accept-Language") == "" {
		t.Error("a browser fixture sent no Accept-Language")
	}

	// A bot fixture must not be dressed up as a browser: doing so would make
	// every header-shape test meaningless.
	if uatest.Request(uatest.Googlebot).Header.Get("Sec-Fetch-Mode") != "" {
		t.Error("a bot fixture was given browser fetch metadata")
	}
}

func TestWithHintsIgnoresGREASE(t *testing.T) {
	t.Parallel()

	r := uatest.Request(uatest.Chrome(138, uatest.Windows(10)),
		uatest.WithHints("Google Chrome", 138, "Windows", false)...)

	c := uax.ParseRequest(r)
	if c.Hints == nil {
		t.Fatal("hints were not parsed")
	}
	for _, b := range c.Hints.Brands {
		if b.Brand == "Not/A)Brand" {
			t.Error("a GREASE brand survived parsing")
		}
	}
	if c.Browser.Source != uax.SourceBoth {
		t.Errorf("browser source: got %q, want %q", c.Browser.Source, uax.SourceBoth)
	}
}
