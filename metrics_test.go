package uax

import "testing"

func TestNoopReporter(t *testing.T) {
	n := NoopReporter{}
	n.IncParse()
	n.IncCacheHit()
	n.IncBotDetected()
	n.IncSpoofDetected()
}

type mockReporter struct {
	NoopReporter
	parsed int
}

func (m *mockReporter) IncParse() { m.parsed++ }

func TestParserMetrics(t *testing.T) {
	m := &mockReporter{}
	p := New(Config{Metrics: m})

	p.Parse("Mozilla/5.0")
	if m.parsed == 0 {
		t.Error("expected metrics to be invoked")
	}
}

type mockSpoofReporter struct {
	NoopReporter
	spoofed int
}

func (m *mockSpoofReporter) IncSpoofDetected() { m.spoofed++ }

func TestParseRequestMetricsSpoof(t *testing.T) {
	m := &mockSpoofReporter{}
	p := New(Config{Metrics: m})

	// Create a spoofed request: Safari on Windows
	req := request("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15")
	p.ParseRequest(req)

	if m.spoofed == 0 {
		t.Error("expected IncSpoofDetected to be called")
	}
}

func TestParseRequestWithCache(t *testing.T) {
	p := New(Config{Cache: CacheConfig{Enabled: true, Size: 100}})
	req := request("Mozilla/5.0")

	c1 := p.ParseRequest(req)
	c2 := p.ParseRequest(req)

	if c1 == c2 {
		t.Error("expected ParseRequest to return cloned clients when cache is enabled")
	}
}
