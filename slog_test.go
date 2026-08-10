package uax

import (
	"log/slog"
	"testing"
)

func TestClientAttrs(t *testing.T) {
	var nilClient *Client
	if attrs := nilClient.Attrs(); attrs != nil {
		t.Errorf("expected nil for nil client, got %v", attrs)
	}

	c := &Client{
		Kind: KindBot,
		Browser: Browser{
			Name:    "Chrome",
			Version: Version{Raw: "120.0.0.0"},
			Mode:    "Headless",
		},
		OS: OS{
			Name:    "Windows",
			Version: Version{Raw: "10"},
		},
		Device: Device{
			Type:  DeviceDesktop,
			Brand: "Dell",
			Model: "XPS",
			Name:  "XPS 13",
		},
		Bot: &Bot{
			Name:     "Googlebot",
			Category: BotSearchEngine,
			Provider: "Google",
		},
		HTTPClient: &HTTPClient{Name: "curl"},
		Automation: &Automation{Name: "Puppeteer"},
		Verdict: &Verdict{
			Trust:     0.1,
			Risk:      0.9,
			Recommend: ActionBlock,
		},
	}

	attrs := c.Attrs()
	if len(attrs) == 0 {
		t.Error("expected attributes")
	}

	val := c.LogValue()
	if val.Kind() != slog.KindGroup {
		t.Errorf("expected LogValue to return a group, got %v", val.Kind())
	}
}
