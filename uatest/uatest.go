// Package uatest builds realistic User-Agent strings and requests for tests.
//
// It exists because testing client-aware code is otherwise miserable: people
// paste a User-Agent from a bug report, it drifts out of date, and nobody
// remembers which of its forty characters mattered. The builders here produce
// strings that are internally consistent — the right engine token, a plausible
// version, a model code that matches its brand — and every one of them is
// verified to parse back to what it described.
package uatest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bakhod1r/uax"
)

// Platform describes the operating system a generated User-Agent runs on.
type Platform struct {
	comment string // the parenthesised section
	mobile  bool
	tablet  bool
}

// Windows returns a Windows platform. Version is the marketing release (10, 11).
func Windows(version int) Platform {
	nt := "10.0"
	switch version {
	case 7:
		nt = "6.1"
	case 8:
		nt = "6.2"
	}
	// Windows 11 reports NT 10.0 like Windows 10; only Client Hints can tell
	// them apart, which is exactly the behaviour tests should exercise.
	return Platform{comment: "Windows NT " + nt + "; Win64; x64"}
}

// MacOS returns a macOS platform, using Apple's underscore form.
func MacOS(major, minor, patch int) Platform {
	return Platform{comment: fmt.Sprintf("Macintosh; Intel Mac OS X %d_%d_%d", major, minor, patch)}
}

// Linux returns a desktop Linux platform.
func Linux() Platform { return Platform{comment: "X11; Linux x86_64"} }

// Android returns an Android platform. Pass a Device to name the handset.
func Android(version int, device ...Device) Platform {
	model := "K" // what modern Chromium reports when the model is reduced
	if len(device) > 0 {
		model = device[0].model
	}
	return Platform{comment: fmt.Sprintf("Linux; Android %d; %s", version, model), mobile: true}
}

// AndroidTablet is Android without the "Mobile" token, which is how Chrome
// signals the larger form factor.
func AndroidTablet(version int, device ...Device) Platform {
	p := Android(version, device...)
	p.mobile, p.tablet = false, true
	return p
}

// IOS returns an iPhone platform.
func IOS(major, minor int) Platform {
	return Platform{
		comment: fmt.Sprintf("iPhone; CPU iPhone OS %d_%d like Mac OS X", major, minor),
		mobile:  true,
	}
}

// IPadOS returns an iPad platform.
func IPadOS(major, minor int) Platform {
	return Platform{
		comment: fmt.Sprintf("iPad; CPU OS %d_%d like Mac OS X", major, minor),
		mobile:  true,
		tablet:  true,
	}
}

// Device names a handset by its model code.
type Device struct {
	brand string
	model string
}

// Model returns a device with the given brand and model code. The code is what
// appears in the User-Agent; the brand is what a parser should recover from it.
func Model(brand, model string) Device { return Device{brand: brand, model: model} }

// Named devices, kept here so tests read as intent rather than as model codes.
var (
	GalaxyS24Ultra = Device{"Samsung", "SM-S928B"}
	GalaxyTabS9    = Device{"Samsung", "SM-X710"}
	Pixel8         = Device{"Google", "Pixel 8"}
	// CPH is OPPO's code space, shared with its OnePlus sub-brand; without the
	// imported catalogue a parser can only prove the parent brand.
	OppoCPH9999 = Device{"OPPO", "CPH9999"}
)

// Chrome returns a Chrome User-Agent for the platform.
func Chrome(major int, p Platform) string {
	return blink(p, fmt.Sprintf("Chrome/%d.0.0.0", major))
}

// Edge returns an Edge User-Agent, which carries both the Chrome and Edg tokens
// exactly as the real browser does.
func Edge(major int, p Platform) string {
	return blink(p, fmt.Sprintf("Chrome/%d.0.0.0", major)) + fmt.Sprintf(" Edg/%d.0.0.0", major)
}

// Firefox returns a Firefox User-Agent.
func Firefox(major int, p Platform) string {
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%d.0) Gecko/20100101 Firefox/%d.0", p.comment, major, major)
}

// Safari returns a Safari User-Agent, placing the real version in the separate
// "Version/" product as Apple does.
func Safari(major, minor int, p Platform) string {
	ua := fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%d.%d", p.comment, major, minor)
	if p.mobile {
		return ua + " Mobile/15E148 Safari/604.1"
	}
	return ua + " Safari/605.1.15"
}

// HeadlessChrome returns the User-Agent a headless Chrome sends.
func HeadlessChrome(major int, p Platform) string {
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/%d.0.0.0 Safari/537.36",
		p.comment, major)
}

// WebView returns an Android WebView User-Agent, including the "wv" marker.
func WebView(major int, p Platform) string {
	comment := p.comment + "; wv"
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/%d.0.0.0 Mobile Safari/537.36",
		comment, major)
}

func blink(p Platform, product string) string {
	ua := fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) %s", p.comment, product)
	if p.mobile {
		return ua + " Mobile Safari/537.36"
	}
	return ua + " Safari/537.36"
}

// Bots, as their operators publish them.
const (
	Googlebot   = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	Bingbot     = "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"
	GPTBot      = "Mozilla/5.0 (compatible; GPTBot/1.2; +https://openai.com/gptbot)"
	ClaudeBot   = "Mozilla/5.0 (compatible; ClaudeBot/1.0; +claudebot@anthropic.com)"
	Curl        = "curl/8.4.0"
	GoHTTP      = "Go-http-client/2.0"
	PythonReqs  = "python-requests/2.31.0"
	FacebookBot = "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)"
)

// Request builds a GET request carrying the given User-Agent plus the headers a
// real browser sends alongside it. Extra headers may be supplied as key/value
// pairs, and override the defaults.
func Request(ua string, headers ...string) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("User-Agent", ua)

	// A browser is chatty. A bot or a script is not, and generating those
	// headers for one would make the fixture lie.
	if looksLikeBrowser(ua) {
		r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Header.Set("Accept-Language", "en-US,en;q=0.9")
		r.Header.Set("Accept-Encoding", "gzip, deflate, br")
		r.Header.Set("Sec-Fetch-Site", "none")
		r.Header.Set("Sec-Fetch-Mode", "navigate")
		r.Header.Set("Sec-Fetch-Dest", "document")
		r.Header.Set("Sec-Fetch-User", "?1")
	}

	for i := 0; i+1 < len(headers); i += 2 {
		r.Header.Set(headers[i], headers[i+1])
	}
	return r
}

// WithHints returns the header pairs a Chromium browser sends for the given
// brand, version and platform, including a GREASE entry so callers exercise
// the code path that must ignore it.
func WithHints(brand string, major int, platform string, mobile bool) []string {
	m := "?0"
	if mobile {
		m = "?1"
	}
	return []string{
		uax.HeaderCHUA, fmt.Sprintf(`"Not/A)Brand";v="8", "Chromium";v="%d", %q;v="%d"`, major, brand, major),
		uax.HeaderCHUAMobile, m,
		uax.HeaderCHUAPlatform, fmt.Sprintf("%q", platform),
	}
}

func looksLikeBrowser(ua string) bool {
	return strings.HasPrefix(ua, "Mozilla/") && !strings.Contains(ua, "compatible;")
}
