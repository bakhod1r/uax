package uax

import (
	"net/http"
	"testing"
)

// Chrome's User-Agent Reduction froze the platform values it sends. The parser
// must not repeat a frozen value as if it were observed, because the number is
// not merely imprecise — it is wrong. A phone running Android 15 says 10.
func TestAndroidReductionIsNotReportedAsAVersion(t *testing.T) {
	t.Parallel()

	c := Parse("Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36")

	if !c.OS.Reduced {
		t.Error("OS.Reduced = false for a reduced Chrome Android User-Agent")
	}
	if c.OS.Version.Known() {
		t.Errorf("os version: got %q, want none — Android 10 here is a placeholder", c.OS.Version.Raw)
	}
	if c.OS.Name != "Android" {
		t.Errorf("os name: got %q, want Android — the platform is still known", c.OS.Name)
	}
	if !c.Device.Reduced {
		t.Error("Device.Reduced = false although the model was elided to \"K\"")
	}
	if c.Device.Model != "" {
		t.Errorf("device model: got %q, want empty — \"K\" is a placeholder, not a code", c.Device.Model)
	}
	if c.Device.Type != DeviceMobile {
		t.Errorf("device type: got %q, want Mobile — form factor survives reduction", c.Device.Type)
	}
}

// A real Android device still reports its real version and model, and must not
// be caught by the reduction check.
func TestGenuineAndroidIsUntouched(t *testing.T) {
	t.Parallel()

	for _, ua := range []string{
		"Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36",
		// Android 10 with a real model: the version is genuine, only the
		// combination with the "K" placeholder indicates reduction.
		"Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36",
	} {
		c := Parse(ua)
		if c.OS.Reduced {
			t.Errorf("OS.Reduced = true for a genuine device\n  ua: %s", ua)
		}
		if !c.OS.Version.Known() {
			t.Errorf("os version lost for a genuine device\n  ua: %s", ua)
		}
		if c.Device.Model == "" {
			t.Errorf("device model lost for a genuine device\n  ua: %s", ua)
		}
	}
}

// macOS and Windows freeze their version too, but differently: the value they
// send is a floor rather than a placeholder, so it is kept and flagged.
func TestFrozenDesktopVersionsAreFlaggedButKept(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		ua   string
		want string
	}{
		{
			"macOS caps at 10.15.7",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
			"10_15_7",
		},
		{
			"Windows 11 still reports NT 10.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
			"10.0",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := Parse(tc.ua)
			if !c.OS.Reduced {
				t.Error("OS.Reduced = false for a frozen platform version")
			}
			if got := c.OS.Version.String(); got != tc.want {
				t.Errorf("os version: got %q, want %q — a floor is still information", got, tc.want)
			}
		})
	}
}

// Reduction is a statement about the User-Agent alone. Client Hints carry the
// real value, so a request that supplies them must not stay flagged.
func TestClientHintsClearReduction(t *testing.T) {
	t.Parallel()

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36")
	r.Header.Set(HeaderCHUAPlatform, `"Android"`)
	r.Header.Set(HeaderCHUAPlatformVersion, `"15"`)
	r.Header.Set(HeaderCHUAModel, `"SM-S928B"`)

	c := ParseRequest(r)
	if c.OS.Reduced {
		t.Error("OS.Reduced survived a Sec-CH-UA-Platform-Version that supplied the real value")
	}
	if c.OS.Version.Major != 15 {
		t.Errorf("os version: got %q, want 15 from the hint", c.OS.Version.String())
	}
	if c.Device.Reduced {
		t.Error("Device.Reduced survived a Sec-CH-UA-Model that supplied the real code")
	}
	if c.Device.Model != "SM-S928B" {
		t.Errorf("device model: got %q, want SM-S928B from the hint", c.Device.Model)
	}
}
