package uax

import (
	"net/http"
	"testing"
)

func request(ua string, headers ...string) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("User-Agent", ua)
	for i := 0; i+1 < len(headers); i += 2 {
		r.Header.Set(headers[i], headers[i+1])
	}
	return r
}

const chromeWindows = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"

func TestHintsAndUserAgentAgreeing(t *testing.T) {
	t.Parallel()

	c := ParseRequest(request(chromeWindows,
		HeaderCHUA, `"Not/A)Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
		HeaderCHUAPlatform, `"Windows"`,
		HeaderCHUAMobile, "?0",
	))

	if c.Browser.Name != "Chrome" {
		t.Errorf("browser: got %q, want Chrome", c.Browser.Name)
	}
	if c.Browser.Source != SourceBoth {
		t.Errorf("source: got %q, want %q", c.Browser.Source, SourceBoth)
	}
	if !c.Consistency().OK() {
		t.Errorf("agreeing sources produced findings: %+v", c.Consistency().Findings)
	}
}

// Merging must not destroy the evidence of a disagreement. This is the whole
// point of recording findings during the merge: once the hints have overwritten
// the User-Agent's values, the contradiction is otherwise unrecoverable.
func TestPlatformMismatchSurvivesTheMerge(t *testing.T) {
	t.Parallel()

	c := ParseRequest(request(chromeWindows, HeaderCHUAPlatform, `"macOS"`))

	if c.OS.Name != "macOS" {
		t.Errorf("os: got %q, want macOS (hints are authoritative)", c.OS.Name)
	}
	if !c.Consistency().Has("hints.platform_mismatch") {
		t.Errorf("the disagreement was lost: %+v", c.Consistency().Findings)
	}
	if !c.IsSuspicious() {
		t.Error("a platform mismatch did not raise suspicion")
	}
}

func TestBrowserMismatchIsRecorded(t *testing.T) {
	t.Parallel()

	c := ParseRequest(request(chromeWindows,
		HeaderCHUA, `"Firefox";v="128"`,
		HeaderCHUAPlatform, `"Windows"`,
	))

	if !c.Consistency().Has("hints.browser_mismatch") {
		t.Errorf("browser mismatch not recorded: %+v", c.Consistency().Findings)
	}
	if c.Browser.Name != "Firefox" {
		t.Errorf("browser: got %q, want Firefox (hints win)", c.Browser.Name)
	}
}

func TestVersionMismatchIsRecorded(t *testing.T) {
	t.Parallel()

	// The User-Agent says Chrome 138; the hints say 120. Real browsers cannot
	// disagree with themselves, so the string was edited.
	c := ParseRequest(request(chromeWindows,
		HeaderCHUA, `"Chromium";v="120", "Google Chrome";v="120"`,
		HeaderCHUAPlatform, `"Windows"`,
	))

	if !c.Consistency().Has("hints.browser_version_mismatch") {
		t.Errorf("version mismatch not recorded: %+v", c.Consistency().Findings)
	}
}

// The high-entropy hints are the only way to recover what a reduced
// User-Agent no longer carries.
func TestHighEntropyHintsFillInWhatTheUserAgentLost(t *testing.T) {
	t.Parallel()

	reduced := "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36"

	c := ParseRequest(request(reduced,
		HeaderCHUAPlatform, `"Android"`,
		HeaderCHUAPlatformVersion, `"15"`,
		HeaderCHUAModel, `"SM-S928B"`,
		HeaderCHUAMobile, "?1",
	))

	if c.OS.Version.Major != 15 {
		t.Errorf("os version: got %q, want 15", c.OS.Version.Raw)
	}
	if c.Device.Model != "SM-S928B" {
		t.Errorf("model: got %q, want SM-S928B", c.Device.Model)
	}
	if c.Device.Source != SourceHints {
		t.Errorf("device source: got %q, want %q", c.Device.Source, SourceHints)
	}
}

// Windows 11 reports NT 10.0 like Windows 10; only the platform version hint
// separates them.
func TestWindows11IsRecoveredFromHints(t *testing.T) {
	t.Parallel()

	c := ParseRequest(request(chromeWindows,
		HeaderCHUAPlatform, `"Windows"`,
		HeaderCHUAPlatformVersion, `"15.0.0"`,
	))

	if c.OS.Version.Major != 11 {
		t.Errorf("windows version: got %d, want 11", c.OS.Version.Major)
	}
}

func TestNoHintsLeavesTheUserAgentAlone(t *testing.T) {
	t.Parallel()

	c := ParseRequest(request(chromeWindows))
	if c.Hints != nil {
		t.Errorf("hints: got %+v, want nil", c.Hints)
	}
	if c.OS.Name != "Windows" || c.Browser.Name != "Chrome" {
		t.Errorf("parsed %q / %q, want Chrome / Windows", c.Browser.Name, c.OS.Name)
	}
}

func TestMergeMissingUserAgentData(t *testing.T) {
	c := ParseRequest(request("Unknown UA",
		HeaderCHUA, `"Google Chrome";v="138"`,
		HeaderCHUAPlatform, `"Windows"`,
		HeaderCHUAMobile, "?0",
		HeaderCHUAArch, `"x86"`,
		HeaderCHUABitness, `"64"`,
		HeaderCHUAFormFactors, `"Desktop"`,
	))
	if c.Browser.Name != "Chrome" {
		t.Errorf("expected Chrome, got %q", c.Browser.Name)
	}
	if c.OS.Name != "Windows" {
		t.Errorf("expected Windows, got %q", c.OS.Name)
	}
	if c.OS.Arch != "amd64" {
		t.Errorf("expected amd64, got %q", c.OS.Arch)
	}
	if c.Hints.Bitness != "64" {
		t.Errorf("expected 64, got %q", c.Hints.Bitness)
	}
	if c.Device.Type != DeviceDesktop {
		t.Errorf("expected Desktop, got %q", c.Device.Type)
	}
}

func TestMergeMismatches(t *testing.T) {
	// For mobile_mismatch, UA is desktop but hints say mobile. No form factor hint so it remains desktop.
	c := ParseRequest(request(chromeWindows,
		HeaderCHUAMobile, "?1",
	))
	if !c.Consistency().Has("hints.mobile_mismatch") {
		t.Error("expected mobile mismatch")
	}

	// For form_factor_mismatch, hints contradict each other.
	c2 := ParseRequest(request(chromeWindows,
		HeaderCHUAMobile, "?1",
		HeaderCHUAFormFactors, `"Desktop"`,
	))
	if !c2.Consistency().Has("hints.form_factor_mismatch") {
		t.Error("expected form factor mismatch")
	}

	// For arch_mismatch and bitness_mismatch
	c3 := ParseRequest(request(chromeWindows,
		HeaderCHUAArch, `"arm"`,
		HeaderCHUABitness, `"32"`,
	))
	if !c3.Consistency().Has("hints.arch_mismatch") {
		t.Error("expected arch mismatch")
	}
	if !c3.Consistency().Has("hints.bitness_mismatch") {
		t.Error("expected bitness mismatch")
	}

	// For bitness mismatch with a 32-bit architecture
	c4 := ParseRequest(request("Mozilla/5.0 (Windows NT 10.0; Win32; x86) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36",
		HeaderCHUAArch, `"x86"`,
		HeaderCHUABitness, `"64"`, // x86 implies 32-bit
	))
	if !c4.Consistency().Has("hints.bitness_mismatch") {
		t.Error("expected bitness mismatch for 32-bit arch")
	}
}

func TestModelMismatch(t *testing.T) {
	// Android UA with SM-S928B
	androidUA := "Mozilla/5.0 (Linux; Android 14; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36"
	c := ParseRequest(request(androidUA, HeaderCHUAModel, `"UnknownModel"`))
	if !c.Consistency().Has("hints.model_mismatch") {
		t.Error("expected model mismatch")
	}
}

func TestNormalizeBrandEdgeCases(t *testing.T) {
	c1 := ParseRequest(request("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 OPR/105.0.0.0",
		HeaderCHUA, `"Opera GX";v="105", "Chromium";v="138"`,
	))
	if c1.Browser.Name != "Opera GX" {
		t.Errorf("expected Opera GX, got %q", c1.Browser.Name)
	}

	c2 := ParseRequest(request("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 YaBrowser/24.1.0.0 Safari/537.36",
		HeaderCHUA, `"Yandex Browser";v="24", "Chromium";v="138"`,
	))
	if c2.Browser.Name != "Yandex Browser" {
		t.Errorf("expected Yandex Browser, got %q", c2.Browser.Name)
	}

	c3 := ParseRequest(request("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0",
		HeaderCHUA, `"Microsoft Edge";v="138", "Chromium";v="138"`,
	))
	if c3.Browser.Name != "Edge" {
		t.Errorf("expected Edge, got %q", c3.Browser.Name)
	}
}

func TestSameOS(t *testing.T) {
	if !sameOS("Ubuntu", "Linux") {
		t.Error("Ubuntu should match Linux")
	}
	if !sameOS("Android TV", "Android") {
		t.Error("Android TV should match Android")
	}
	if sameOS("iOS", "macOS") {
		t.Error("iOS should not match macOS")
	}
	if sameOS("Linux", "Windows") {
		t.Error("Linux should not match Windows")
	}
}
