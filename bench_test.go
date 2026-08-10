package uax

import (
	"fmt"
	"testing"
)

// Representative inputs. Real traffic is dominated by a handful of shapes, and
// the awkward ones (an empty string, a bot, a hostile blob) are exactly where a
// parser tends to fall over.
var benchInputs = []struct {
	name string
	ua   string
}{
	{"chrome_windows", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
	{"chrome_android", "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36"},
	{"safari_ios", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1"},
	{"firefox_macos", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:128.0) Gecko/20100101 Firefox/128.0"},
	{"googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
	{"curl", "curl/8.4.0"},
	{"empty", ""},
	{"junk", "\x00\xff not a user agent at all ((((("},
}

func BenchmarkParseTo(b *testing.B) {
	for _, in := range benchInputs {
		b.Run(in.name, func(b *testing.B) {
			var c Client
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				ParseTo(&c, in.ua)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	for _, in := range benchInputs[:4] {
		b.Run(in.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				_ = Parse(in.ua)
			}
		})
	}
}

// BenchmarkParseUncached measures the parser rather than the cache. Benchmarking
// a cached parser on one repeated string measures a map lookup and tells you
// nothing about parsing, which is how libraries end up publishing numbers they
// cannot defend.
func BenchmarkParseUncached(b *testing.B) {
	uas := make([]string, 1024)
	for i := range uas {
		uas[i] = fmt.Sprintf(
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.%d.0 Safari/537.36", i)
	}

	var c Client
	b.ReportAllocs()
	b.ResetTimer()
	i := 0
	for b.Loop() {
		ParseTo(&c, uas[i%len(uas)])
		i++
	}
}

func BenchmarkCacheHit(b *testing.B) {
	p := New(Config{Cache: CacheConfig{Enabled: true}})
	ua := benchInputs[0].ua
	p.Parse(ua)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = p.Parse(ua)
	}
}

func BenchmarkParseRequest(b *testing.B) {
	r := request(benchInputs[0].ua,
		HeaderCHUA, `"Not/A)Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
		HeaderCHUAPlatform, `"Windows"`,
		HeaderCHUAMobile, "?0",
	)

	b.ReportAllocs()

	for b.Loop() {
		_ = ParseRequest(r)
	}
}

// BenchmarkParseParallel is the shape the middleware actually runs in: one
// parser, every request goroutine at once. A parser that shared mutable state
// would show up here as contention long before it showed up as a race.
func BenchmarkParseParallel(b *testing.B) {
	p := New(Config{})

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = p.Parse(benchInputs[i%len(benchInputs)].ua)
			i++
		}
	})
}

// BenchmarkConsistency measures the coherence check on its own. It is the unit
// the verdict is built from, so its cost is paid at least once by every caller
// that asks for either.
func BenchmarkConsistency(b *testing.B) {
	for _, in := range benchInputs[:4] {
		c := Parse(in.ua)
		b.Run(in.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = CheckConsistency(c)
			}
		})
	}
}

// BenchmarkVerdict is why VerdictFrom exists.
//
// A verdict is derived from the consistency findings, so ComputeVerdict has to
// run every spoof rule to produce one. A caller that shows the findings as well
// — ParseRequest, the CLI, the MCP server — already has them, and the two
// sub-benchmarks below are the difference between passing them in and making
// the verdict recompute them.
func BenchmarkVerdict(b *testing.B) {
	c := Parse(benchInputs[0].ua)
	cons := CheckConsistency(c)

	b.Run("recomputed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = ComputeVerdict(c)
		}
	})
	b.Run("reused", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = VerdictFrom(c, cons)
		}
	})
}

// BenchmarkScaling documents the property that separates this parser from a
// regex list: cost tracks input length, not the number of rules. Both cases
// below run against the same full database.
func BenchmarkScaling(b *testing.B) {
	short := "curl/8.4.0"
	long := benchInputs[1].ua

	var c Client
	b.Run("short_input", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			ParseTo(&c, short)
		}
	})
	b.Run("long_input", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			ParseTo(&c, long)
		}
	})
}

// TestAllocationBudget is the gate. Allocation-freedom on the hot path is a
// design property — token views into the input, constants from the rule
// tables — and it should fail the build the day it is lost, not appear as a
// number nobody reads.
func TestAllocationBudget(t *testing.T) {
	var c Client

	for _, in := range benchInputs {
		allocs := testing.AllocsPerRun(100, func() { ParseTo(&c, in.ua) })
		if allocs > 0 {
			t.Errorf("ParseTo(%s): %.0f allocs/op, want 0", in.name, allocs)
		}
	}
}
