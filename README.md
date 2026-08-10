---
title: uax
# The site has no index.md on purpose: a landing page separate from the README
# is a second copy of the same text, and the two drift. This makes the README
# the root of the site instead.
permalink: /
---

# uax

Client intelligence for Go. Not another User-Agent parser: `uax` reads the
User-Agent **and** Client Hints **and** the fetch metadata around them, and
tells you what is talking to your server — including when it is lying.

```go
c := uax.ParseRequest(r)

c.Browser.Name    // "Chrome"
c.OS.Name         // "Android"
c.Device.Model    // "SM-S928B"   — verbatim from the request
c.IsAICrawler()   // false
c.Verdict.Risk    // 0.00
```

Zero dependencies. Zero allocations on the parse path.

Naming the hardware is a separate question, and a separate module. `uax`
extracts the model code; a catalogue turns it into a handset:

```go
import "github.com/bakhod1r/devicex"

p := uax.New(uax.Config{Device: devicex.Describe})
c := p.ParseRequest(r)

c.Device.Brand    // "Samsung"
c.Device.Name     // "Galaxy S24 Ultra"
c.Device.Type     // "Mobile"
```

Without the hook those stay empty, and `Device.Model` still reports the code —
which is the honest answer when nothing can name it. Twenty-five thousand
device names are a large thing to link into a binary that does not want them,
and they go stale on a schedule a parser does not share.

---

## Install

```bash
go get github.com/bakhod1r/uax
```

## Use

```go
// Pure function: no cache, no shared state, no retention.
c := uax.Parse(ua)

// Preferred: Client Hints are structured, and they survive the User-Agent
// reduction that is steadily emptying the legacy string.
c := uax.ParseRequest(r)

// Configured parser for a server under load.
p := uax.New(uax.Config{
    Cache:   uax.CacheConfig{Enabled: true, Size: 8192},
    Metrics: reporter,
})
```

HTTP middleware, parsing lazily so routes that ignore the client pay nothing:

```go
mux.Handle("/", middleware.New(p)(handler))

func handler(w http.ResponseWriter, r *http.Request) {
    c := middleware.FromContext(r.Context()) // never nil
}
```

---

## What it finds

| | |
|---|---|
| **Browser** | 55 rules — Chrome, Firefox, Safari, Edge, Opera, Brave, Vivaldi, Samsung Internet, Yandex, UC, QQ, Silk, IE; the regional majors (Whale, Coc Coc, Quark, Sogou, 360, Baidu, MIUI, Huawei, OPPO, vivo); the Gecko and Goanna forks (LibreWolf, Pale Moon, Basilisk); the independents (Konqueror, GNOME Web, Falkon, Midori); plus family, variant and mode (headless, WebView) |
| **OS** | 44 rules — Windows (NT mapped to its release), macOS, iOS, iPadOS, tvOS, watchOS, visionOS, Android, ChromeOS, HarmonyOS, KaiOS, Tizen, webOS, Sailfish, Symbian, ten Linux distributions, the BSDs, Solaris, Haiku, SteamOS, Roku, Chromecast, consoles, and architecture where stated |
| **Device** | the model code, verbatim from the request; brand, name and form factor come from a catalogue via `Config.Device` |
| **Bot** | 41 rules — categorised: SearchEngine, AICrawler, AIAgent, SocialCrawler, SEO, Monitoring, Security |
| **AI** | GPTBot, ClaudeBot, PerplexityBot, Google-Extended, Bytespider, Amazonbot, CCBot, and the live agents (ChatGPT-User, Claude-User) |
| **In-app** | 12 rules — Instagram, Facebook, TikTok, Telegram, WeChat, LINE, Snapchat, Slack |
| **Automation** | 9 rules — Headless Chrome/Firefox, Playwright, Puppeteer, Selenium, Cypress, WebDriver |
| **Form factor** | `Sec-CH-UA-Form-Factors` when sent — Desktop, Mobile, Tablet, XR, EInk, Watch, Automotive — stated by the browser rather than inferred |
| **HTTP clients** | 19 rules — curl, Wget, Go-http-client, python-requests, httpx, okhttp, Java, Postman, axios, Guzzle |

187 rules, each shipping with its own test. A rule without one fails the lint
suite, so the database cannot rot quietly.

### Frozen values are reported as frozen

Chrome's User-Agent Reduction replaced the platform section with constants.
Every reduced Android build says `Android 10; K` whatever it is really running,
macOS caps at `10_15_7`, and Windows has said `NT 10.0` since Windows 10.

Repeating those as observed facts is the most common way a parser lies, so uax
does not:

```go
c := uax.Parse("Mozilla/5.0 (Linux; Android 10; K) … Chrome/138.0.0.0 Mobile Safari/537.36")

c.OS.Name          // "Android"  — the platform is still known
c.OS.Version       // zero       — "10" here is a placeholder, not a version
c.OS.Reduced       // true
c.Device.Type      // Mobile     — the form factor survives reduction
c.Device.Model     // ""         — "K" is not a model code
c.Device.Reduced   // true
```

macOS and Windows are treated differently, because their frozen value is a
floor rather than a contradiction: the version is kept and `Reduced` is set
beside it. `Sec-CH-UA-Platform-Version` and `Sec-CH-UA-Model` carry the real
values, and merging them clears the flag.

---

## The device-name question

`SM-S928B` contains no hint that it means "Galaxy S24 Ultra". Samsung assigned
it. No regex derives it, and no amount of confidence substitutes for looking it
up — so this module does not try. It reports what the request carried:

```go
c.Device.Model  // "SM-S928B" — verbatim evidence, never modified
```

Everything past that is a catalogue's job, and catalogues have their own release
schedule: phones ship weekly, parsers do not. `Config.Device` is the seam.
[devicex](https://github.com/bakhod1r/devicex) fills it:

```go
p := uax.New(uax.Config{Device: devicex.Describe})

c.Device.Brand  // "Samsung"          — from the code's shape, so a phone released tomorrow still resolves
c.Device.Name   // "Galaxy S24 Ultra" — catalogue only; EMPTY when absent, never guessed
c.Device.Type   // "Tablet"           — SM-X710 is a Galaxy Tab, and the shape says so
```

The two tiers fail differently on purpose. The brand tier is rules and degrades
gracefully; the name tier is a fact about the world, and its failure is silence.

Without the hook, `uax device SM-S928B` says so plainly rather than guessing:

```
model  SM-S928B
brand  Samsung
name   (not in the catalogue)
```

No other Go parser resolves the marketing name at all; `uap-go` returns
`Samsung SM-S928B` and stops. The difference here is a mechanism with
provenance, not a bigger pile of remembered strings.

---

## Spoof detection

Client Hints and the User-Agent are two independent statements. When they
disagree, the disagreement *is* the signal — so merging records it instead of
destroying it.

```go
c := uax.ParseRequest(r)

c.Consistency().Findings
// [{Rule: "hints.platform_mismatch", Severity: "critical",
//   Left: "Windows", Right: "macOS"}]

c.Verdict.Risk       // 1.0
c.Verdict.Recommend  // "block" — advice, never a decision
```

Impossibility rules — Safari on Windows, a non-WebKit engine on iOS, a client
claiming two platforms at once — are facts, not heuristics. Weaker signals
(no `Accept-Language`, Chromium without `Sec-CH-UA`) move risk without ever
deciding anything alone.

Trust and risk answer different questions. A well-behaved crawler that names
itself is highly trusted and low risk: being a bot is not being a liar.

---

## Performance

Measured on an Apple M4 Pro, Go 1.26, against the full 187-rule database:

```
BenchmarkParseTo/chrome_windows      972 ns/op     0 B/op   0 allocs/op
BenchmarkParseTo/chrome_android     1093 ns/op     0 B/op   0 allocs/op
BenchmarkParseTo/firefox_macos       772 ns/op     0 B/op   0 allocs/op
BenchmarkParseTo/googlebot           532 ns/op     0 B/op   0 allocs/op
BenchmarkParseTo/curl                251 ns/op     0 B/op   0 allocs/op
BenchmarkParseTo/empty                21 ns/op     0 B/op   0 allocs/op
BenchmarkCacheHit                     23 ns/op     0 B/op   0 allocs/op
```

Against the ecosystem, same input, same machine:

| | ns/op | allocs/op | what it resolved from that string |
|---|--:|--:|---|
| **uax** | **1093** | **0** | Chrome 138.0.0.0 · Android 15 · SM-S928B · Mobile (plus **Samsung Galaxy S24 Ultra** with a catalogue) |
| mssola/user_agent | 786 | 20 | Chrome 138.0.0.0 · Android 15 · platform Linux |
| ua-parser/uap-go | 1412 | 8 | Chrome Mobile 138 · Android 15 · Samsung · SM-S928B |
| mileusna/useragent | 1426 | 16 | Chrome 138.0.0.0 · Android 15 · SM-S928B |
| avct/uasurfer | 1612 | 1 | Chrome 138 · Android 15 · phone |

`mssola` is faster because it resolves less — no device at all. Compared like
for like it wins on the work it skips. `mileusna` returns the model code but no
brand; `uap-go` returns both but no marketing name, which no published Go
parser resolves. Reproduce the table with `make compare`.

Zero allocations is enforced, not observed: `TestAllocationBudget` fails the
build if `ParseTo` ever allocates.

**Why it stays flat.** The input is tokenised once into a fixed stack buffer,
and each token is then looked up once in a combined index that answers for all
seven categories simultaneously. Cost tracks input length, not rule count:
`curl/8.4.0` costs 251 ns and a full Chrome-on-Android string 1093 ns, against
the same 187 rules. Adding a thousand more would not move either number. A
conventional regex-list parser is O(rules × length) and gets slower with every
rule anyone contributes.

---

## Design commitments

**Never invent.** Unknown model → empty `Device.Name`. Unknown client → `Kind`
is `Unknown`, not "Desktop by elimination". Everything derivable only from a
lookup is empty when the lookup misses. Facts that cannot be derived from the
input are not this module's to assert: they live in a catalogue that carries
their provenance.

**Never lie about certainty.** Every resolved field carries a `Confidence` and
a `Source` (`ua`, `hints`, `both`, `inferred`). A bot's identity is a claim, and
the verdict says so explicitly until `uax/verify` checks it against DNS and
published address ranges.

**No hidden state.** `uax.Parse` is pure. Caching, recording and metrics are all
opt-in, because a cache is a global side effect on what is otherwise a function,
it retains attacker-supplied strings, and it makes benchmarks measure the cache
instead of the parser.

**No network, ever, in the core.** No request, no file, no environment read.
The unknown-token recorder writes to memory the caller owns, after redacting
versions and opaque identifiers. See [SECURITY.md](SECURITY.md).

**Every rule is proven.** Classification lives in `internal/table`, in Go, and
a rule without a test fails the build — along with a rule that has no priority,
no confidence, two match clauses, or an ID its category already uses.

---

## CLI

```bash
go install github.com/bakhod1r/uax/cmd/uax@latest

uax parse   "Mozilla/5.0 (…)"    # one-line summary
uax explain "Mozilla/5.0 (…)"    # fields, confidence, verdict, contradictions
uax analyze access.log           # what actually reaches your service
uax device  SM-S928B             # resolve a bare model code
uax db                           # rule counts
uax robots --deny ai             # robots.txt from the bot database
```

`uax robots` is the bot database written in a different syntax. The token a rule
matches on is the token a crawler answers to in robots.txt, so the file cannot
drift away from the parser — a test asserts every emitted name parses back to
the bot it claims to block.

```
$ uax robots --deny ai --sitemap https://example.com/sitemap.xml
# Generated by uax from its bot database. Advisory only: robots.txt is
# obeyed voluntarily, and a crawler that ignores it is unaffected.

# AIAgent
User-agent: ChatGPT-User
Disallow: /
User-agent: Claude-User
Disallow: /
…
# AICrawler
User-agent: CCBot
Disallow: /
User-agent: ClaudeBot
Disallow: /
User-agent: GPTBot
Disallow: /
…
Sitemap: https://example.com/sitemap.xml
```

The two categories are separated because they are different questions. An
AICrawler collects training data; an AIAgent is fetching a page a person asked
for, right now. Blocking the second blocks a user.

Same selection is available to callers as `robots.Agents`, so a WAF rule set
or an nginx map can be built from it instead.

```
$ uax explain "Mozilla/5.0 (Linux; Android 15; SM-S928B) … Chrome/138.0.0.0 Mobile Safari/537.36"
FIELD     VALUE             CONFIDENCE
kind      Human
browser   Chrome 138.0.0.0  0.95
  family  Chrome
  source  ua
engine    Blink 537.36      0.95
os        Android 15        0.96
device    Mobile            0.99
  brand   Samsung
  model   SM-S928B
  name    Galaxy S24 Ultra

verdict: allow  trust=0.70  risk=0.00
```

The CLI wires in [devicex](https://github.com/bakhod1r/devicex), which is where
`Galaxy S24 Ultra` comes from. The library on its own reports the model code
and leaves the name empty.

---

## MCP server

`cmd/uax-mcp` hands the same answers to a model, over the Model Context
Protocol. It exists because a User-Agent is exactly the kind of string a model
will read and guess at, and the strings are built to defeat guessing: every
browser claims `Mozilla/5.0` for historical reasons, and nothing in `SM-S928B`
suggests Galaxy S24 Ultra.

```bash
go install github.com/bakhod1r/uax/cmd/uax-mcp@latest
```

```json
{
  "mcpServers": {
    "uax": { "command": "uax-mcp" }
  }
}
```

| Tool | Answers |
| --- | --- |
| `parse_user_agent` | Which browser, engine, OS, device, bot or library sent this — each field with its confidence and source |
| `check_consistency` | Whether the client's claims contradict each other, with a trust/risk verdict |
| `lookup_device` | What handset a bare model code names |
| `database_info` | What the rule database covers, its fingerprint, and the catalogue's snapshot date |

An unrecognised string and an unknown model code come back as tool errors, not
as empty fields — the difference between "no" and "no answer" is one a model
acts on. Nothing here reaches the network, reads a file or retains a request.

---

## Testing your own code

`uatest` builds fixtures that are internally consistent — right engine token,
plausible version, model code matching its brand — and every builder is
round-trip verified against the parser.

```go
ua := uatest.Chrome(138, uatest.Android(15, uatest.GalaxyS24Ultra))
r  := uatest.Request(ua, uatest.WithHints("Google Chrome", 138, "Android", true)...)
```

---

## Contributing a rule

Add the rule to its category in `internal/table/rules.go`, keeping the slice
ordered by descending priority — the resolver stops at the first match, so
order decides which rule wins:

```go
{ID: "my_browser", Name: "My Browser", Family: "Chrome", Engine: "Blink",
    Token: "MyBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.95},
```

Then add the test that proves it, in `internal/table/tests.go`:

```go
{RuleID: "my_browser", UA: "Mozilla/5.0 (…) MyBrowser/3.2 Safari/537.36",
    Expect: map[string]string{"browser": "My Browser"}},
```

Then `go test ./...`. A rule with no test does not compile past the lint suite.

Device *names* are not contributed here at all. They are a catalogue, they
carry a provenance, and they live in
[devicex](https://github.com/bakhod1r/devicex).

## License

MIT — see [LICENSE](LICENSE).

This module ships no device names. They are imported data rather than original
work, so they live in [devicex](https://github.com/bakhod1r/devicex) with their
provenance attached. [NOTICE.md](NOTICE.md) records the specifications these
rules implement.
