// Package uax turns a request into a description of the client behind it.
//
// It reads the User-Agent string, the Client Hints headers, and the fetch and
// content-negotiation metadata that travel with them, and reports the browser,
// operating system, device, bot, in-app WebView, automation framework or HTTP
// library responsible — along with how confident it is, whether the client's
// claims contradict each other, and what it suggests doing about that.
//
// # Getting started
//
//	c := uax.Parse(r.UserAgent())
//	fmt.Println(c.Browser.Name, c.OS.Name, c.Device.Name)
//
// Prefer ParseRequest wherever the whole request is available. Client Hints are
// structured, are what browser vendors intend servers to read, and survive the
// User-Agent reduction that is steadily emptying the legacy string:
//
//	c := uax.ParseRequest(r)
//	if c.IsAICrawler() { … }
//
// # Configuration
//
// The package-level functions are pure: no cache, no recorder, no shared state.
// A server that wants caching, metrics or custom rules builds its own parser:
//
//	p := uax.New(uax.Config{
//	    Cache:   uax.CacheConfig{Enabled: true},
//	    Metrics: myReporter,
//	})
//
// # What this package will not do
//
// It never invents facts. A device model that is not in the database yields an
// empty Device.Name rather than a guess, because a wrong device name is worse
// than none. Predicates require positive evidence: an unrecognised client is
// unknown, not a desktop by elimination. And everything a client says about
// itself is a claim — Bot.Name is what the client typed, not proof of who sent
// the request.
//
// # Rules
//
// Classification rules live in internal/table, in Go. They are reviewed as
// code, and a set of lint tests holds the line that once belonged to a
// generator: every rule has a unique ID within its category, exactly one match
// clause, a priority, a confidence in (0,1], and at least one test proving it.
// A rule without a test fails the build.
//
// # Devices
//
// Naming hardware is a catalogue problem, not a parsing one — "SM-S928B" does
// not contain "Galaxy S24 Ultra", and no rule can derive it — and catalogues
// move on their own schedule, since phones ship weekly and parsers do not. So
// this package holds no device names. It extracts the model code a request
// carries and hands it to whatever the caller configured:
//
//	import "github.com/bakhod1r/devicex"
//
//	p := uax.New(uax.Config{Device: devicex.Describe})
//
// Without Config.Device, Device.Model still reports the code the request
// carried. Only the marketing name, brand and form factor go unanswered, which
// is the honest result of having nothing to answer them with.
package uax
