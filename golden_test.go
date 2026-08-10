package uax

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// goldenCase mirrors one line of testdata/corpus/*.jsonl. Every field except
// "ua" is optional: an empty expectation is not checked, so a corpus line
// asserts only what it explicitly states.
type goldenCase struct {
	UA string `json:"ua"`

	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	Engine         string `json:"engine"`

	OS        string `json:"os"`
	OSVersion string `json:"os_version"`
	Arch      string `json:"arch"`

	DeviceType  string `json:"device_type"`
	DeviceBrand string `json:"device_brand"`
	DeviceModel string `json:"device_model"`
	DeviceName  string `json:"device_name"`
	// A sub-brand is reported as its product line, not as the manufacturer: a
	// Redmi Note is made by Xiaomi. Brand and family together keep both facts.
	DeviceFamily string `json:"device_family"`

	Bot         string `json:"bot"`
	BotCategory string `json:"bot_category"`
	BotProvider string `json:"bot_provider"`
	AI          string `json:"ai"`

	App        string `json:"app"`
	Automation string `json:"automation"`

	HTTPClient         string `json:"http_client"`
	HTTPClientVersion  string `json:"http_client_version"`
	HTTPClientCategory string `json:"http_client_category"`

	Kind     string `json:"kind"`
	WebView  bool   `json:"webview"`
	Headless bool   `json:"headless"`
	CLI      bool   `json:"cli"`
}

func loadGolden(t *testing.T) []goldenCase {
	t.Helper()

	paths, err := filepath.Glob(filepath.Join("testdata", "corpus", "*.jsonl"))
	if err != nil {
		t.Fatalf("glob corpus: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("no corpus files found under testdata/corpus")
	}

	var cases []goldenCase
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			t.Fatalf("open %s: %v", p, err)
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
		for line := 1; sc.Scan(); line++ {
			text := strings.TrimSpace(sc.Text())
			if text == "" || strings.HasPrefix(text, "//") {
				continue
			}
			var c goldenCase
			if err := json.Unmarshal([]byte(text), &c); err != nil {
				f.Close()
				t.Fatalf("%s:%d: %v", p, line, err)
			}
			cases = append(cases, c)
		}
		if err := sc.Err(); err != nil {
			f.Close()
			t.Fatalf("read %s: %v", p, err)
		}
		f.Close()
	}
	return cases
}

func TestGoldenCorpus(t *testing.T) {
	t.Parallel()

	for _, tc := range loadGolden(t) {
		name := tc.UA
		if name == "" {
			name = "(empty)"
		}
		if len(name) > 60 {
			name = name[:60]
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := parseWithDevices(tc.UA)

			eq := func(field, want, got string) {
				t.Helper()
				if want != "" && want != got {
					t.Errorf("%s: want %q, got %q", field, want, got)
				}
			}
			// Versions are compared by value, not by spelling: a User-Agent may
			// write the same version as "17_5", "17.5" or "10.0", and the
			// corpus asserts what the version is, not how it was written.
			eqVersion := func(field, want string, got Version) {
				t.Helper()
				if want == "" {
					return
				}
				if w := ParseVersion(want); w.Compare(got) != 0 {
					t.Errorf("%s: want %q, got %q", field, want, got.Raw)
				}
			}
			isTrue := func(field string, want, got bool) {
				t.Helper()
				if want && !got {
					t.Errorf("%s: want true, got false", field)
				}
			}

			eq("browser", tc.Browser, c.Browser.Name)
			eqVersion("browser.version", tc.BrowserVersion, c.Browser.Version)
			eq("engine", tc.Engine, c.Engine.Name)

			eq("os", tc.OS, c.OS.Name)
			eqVersion("os.version", tc.OSVersion, c.OS.Version)
			eq("os.arch", tc.Arch, c.OS.Arch)

			eq("device.type", tc.DeviceType, string(c.Device.Type))
			eq("device.brand", tc.DeviceBrand, c.Device.Brand)
			eq("device.model", tc.DeviceModel, c.Device.Model)
			eq("device.name", tc.DeviceName, c.Device.Name)
			eq("device.family", tc.DeviceFamily, c.Device.Family)

			if tc.Bot != "" {
				if c.Bot == nil {
					t.Errorf("bot: want %q, got nil", tc.Bot)
				} else {
					eq("bot", tc.Bot, c.Bot.Name)
					eq("bot.category", tc.BotCategory, string(c.Bot.Category))
					eq("bot.provider", tc.BotProvider, c.Bot.Provider)
				}
			}
			if tc.AI != "" {
				if c.AI == nil {
					t.Errorf("ai: want %q, got nil", tc.AI)
				} else {
					eq("ai", tc.AI, c.AI.Name)
				}
			}
			if tc.App != "" {
				if c.App == nil {
					t.Errorf("app: want %q, got nil", tc.App)
				} else {
					eq("app", tc.App, c.App.Name)
				}
			}
			if tc.Automation != "" {
				if c.Automation == nil {
					t.Errorf("automation: want %q, got nil", tc.Automation)
				} else {
					eq("automation", tc.Automation, c.Automation.Name)
				}
			}
			if tc.HTTPClient != "" {
				if c.HTTPClient == nil {
					t.Errorf("http_client: want %q, got nil", tc.HTTPClient)
				} else {
					eq("http_client", tc.HTTPClient, c.HTTPClient.Name)
					eqVersion("http_client.version", tc.HTTPClientVersion, c.HTTPClient.Version)
					eq("http_client.category", tc.HTTPClientCategory, c.HTTPClient.Category)
				}
			}

			eq("kind", tc.Kind, string(c.Kind))
			isTrue("webview", tc.WebView, c.IsWebView())
			isTrue("headless", tc.Headless, c.IsHeadless())
			isTrue("cli", tc.CLI, c.IsCLI())
		})
	}
}

// TestGoldenNoPanic guarantees the corpus never panics through the request path
// either, not just the plain string path.
func TestGoldenNoPanic(t *testing.T) {
	t.Parallel()
	for _, tc := range loadGolden(t) {
		_ = parseWithDevices(tc.UA)
	}
}
