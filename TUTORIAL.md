---
title: Tutorial
---

# Tutorial

The README is a reference: it tells you what exists. This walks you through
using it, in the order the questions actually come up.

Every snippet below is complete and runnable. Nothing is elided behind a `…`.

- [1. The first parse](#1-the-first-parse)
- [2. Read the whole request, not the string](#2-read-the-whole-request-not-the-string)
- [3. Put it in a server](#3-put-it-in-a-server)
- [4. Name the handset](#4-name-the-handset)
- [5. Tell bots apart from people](#5-tell-bots-apart-from-people)
- [6. Catch a client that is lying](#6-catch-a-client-that-is-lying)
- [7. Make it fast under load](#7-make-it-fast-under-load)
- [8. Test your own handlers](#8-test-your-own-handlers)
- [9. Find out what you do not recognise](#9-find-out-what-you-do-not-recognise)
- [What to avoid](#what-to-avoid)

---

## 1. The first parse

```bash
mkdir uax-tutorial && cd uax-tutorial
go mod init example.com/tutorial
go get github.com/bakhod1r/uax
```

```go
package main

import (
	"fmt"

	"github.com/bakhod1r/uax"
)

func main() {
	c := uax.Parse("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36")

	fmt.Println(c.Browser.Name)    // Chrome
	fmt.Println(c.OS.Name)         // Android
	fmt.Println(c.Device.Model)    // SM-S928B
	fmt.Println(c.Device.Type)     // Mobile
	fmt.Println(c.Device.Name)     // (empty — see step 4)
}
```

`uax.Parse` is a pure function. No cache, no shared state, nothing retained. You
can call it from anywhere, concurrently, without thinking about it.

**`Device.Name` is empty and that is correct.** `SM-S928B` does not contain the
words "Galaxy S24 Ultra" — Samsung assigned that mapping, and no parsing
derives it. Step 4 fills it in.

### Confidence and source

Every resolved field carries how sure the answer is and where it came from:

```go
fmt.Println(c.Browser.Confidence)  // 0.95
fmt.Println(c.Browser.Source)      // ua
```

`Source` is one of `ua`, `hints`, `both` or `inferred`. It matters more than it
looks: `ua` means a client typed it and could have typed anything, while
`hints` means the browser sent it as structured data. When you store a
classification, store the source with it.

---

## 2. Read the whole request, not the string

The User-Agent is being deliberately emptied. Chrome's User-Agent Reduction
froze the platform section: every reduced Android build now says `Android 10; K`
regardless of what it runs. The real values moved to Client Hints.

So prefer `ParseRequest`:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	c := uax.ParseRequest(r)   // reads User-Agent, Sec-CH-UA*, Sec-Fetch-*, Accept
	fmt.Fprintln(w, c.Browser.Name, c.OS.Name, c.OS.Version)
}
```

See what reduction does, and what `uax` refuses to do about it:

```go
c := uax.Parse("Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 " +
	"(KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36")

c.OS.Name        // "Android"  — still known
c.OS.Version     // zero       — "10" is a placeholder, not a version
c.OS.Reduced     // true
c.Device.Model   // ""         — "K" is not a model code
```

Most parsers report `Android 10` here and are wrong on nearly every modern
Android request. `uax` reports that it cannot tell. **Check `OS.Reduced` before
you put an OS version in a chart** — otherwise your chart says the whole world
runs Android 10.

If the client sent `Sec-CH-UA-Platform-Version`, `ParseRequest` uses it and
clears the flag. That is the whole reason to prefer it.

---

## 3. Put it in a server

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/bakhod1r/uax"
	"github.com/bakhod1r/uax/middleware"
)

func main() {
	p := uax.New(uax.Config{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := middleware.FromContext(r.Context())   // never nil
		fmt.Fprintf(w, "%s on %s\n", c.Browser.Name, c.OS.Name)
	})

	http.ListenAndServe(":8080", middleware.New(p)(mux))
}
```

```bash
curl -s localhost:8080 -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15'
# Safari on macOS
```

Parsing is **lazy**: a route that never calls `FromContext` costs nothing. You
can mount it globally without auditing which handlers need it.

`FromContext` never returns nil. On a route where the middleware was not
mounted you get an empty `*uax.Client`, so a handler cannot panic in production
over a routing mistake.

Gin, Echo and Fiber have their own adapters, in separate modules so the library
never pulls a web framework into your build:

```bash
go get github.com/bakhod1r/uax/middleware/gin@latest
```

```go
r.Use(uaxgin.Middleware(nil))   // nil = the package default parser

r.GET("/", func(c *gin.Context) {
	client := uaxgin.FromContext(c)
	...
})
```

---

## 4. Name the handset

`Device.Model` is evidence — the code the request carried, byte for byte.
Turning it into "Galaxy S24 Ultra" needs a catalogue, and catalogues have a
release schedule a parser does not share: phones ship weekly.

So it is a separate module, wired through a hook:

```bash
go get github.com/bakhod1r/devicex
```

```go
p := uax.New(uax.Config{Device: devicex.Describe})
c := p.Parse(androidUA)

c.Device.Model   // "SM-S928B"
c.Device.Brand   // "Samsung"
c.Device.Name    // "Galaxy S24 Ultra"
c.Device.Type    // "Mobile"
```

Two tiers that fail differently, on purpose:

| tier | source | when it misses |
|---|---|---|
| brand | the shape of the code | a phone released tomorrow still resolves `Samsung` |
| name | the catalogue | **empty** — never approximated |

A wrong device name is worse than no device name, because it is a fact-shaped
value that you will store, aggregate and put in a report. Handle the empty case
rather than expecting a name:

```go
name := c.Device.Name
if name == "" {
	name = c.Device.Model   // the code is still true
}
```

Without the hook everything past `Model` stays empty, and your binary does not
carry twenty-five thousand device names it never asked for.

---

## 5. Tell bots apart from people

```go
c := uax.Parse("Mozilla/5.0 (compatible; ClaudeBot/1.0; +claudebot@anthropic.com)")

c.Kind            // AICrawler
c.Bot.Name        // ClaudeBot
c.Bot.Category    // AICrawler
c.Bot.Provider    // Anthropic
c.Browser.Name    // "" — it is not a browser, and saying "Mozilla" would be a lie
```

Bots are categorised because the categories call for different responses:

| category | what it is |
|---|---|
| `SearchEngine` | Googlebot, Bingbot — you probably want these |
| `AICrawler` | GPTBot, ClaudeBot, CCBot — collecting training data |
| `AIAgent` | ChatGPT-User, Claude-User — **fetching for a person, right now** |
| `SocialCrawler` | link previews |
| `SEO`, `Monitoring`, `Security` | tools |

**`AICrawler` and `AIAgent` are not the same decision.** A crawler collects
training data on its own schedule. An agent is fetching a page because a human
asked, this second. Blocking the second blocks a user.

Other things a "browser" can turn out to be:

```go
c.App          // in-app WebView (Instagram, TikTok, WeChat…)
c.Automation   // Playwright, Puppeteer, Selenium, Headless Chrome
c.HTTPClient   // curl, python-requests, okhttp, Go-http-client
```

Convenience predicates exist — `IsBot`, `IsAICrawler`, `IsWebView`,
`IsHeadless`, `IsAutomated`, `IsMobile`, `IsAndroid` and so on:

```go
if c.IsAICrawler() {
	w.Header().Set("X-Robots-Tag", "noai")
}
```

Each requires positive evidence. `IsDesktop()` is **false** for an unrecognised
client, not true by elimination.

### The limit worth knowing

`Bot.Name` is the name in the string, not proof of who sent the request.
Anything can send `Googlebot/2.1`. Verifying a crawler needs reverse DNS or the
operator's published address ranges, which this module does not do. And
automation configured to hide is invisible here by design — **a nil
`Automation` never means "human"**.

---

## 6. Catch a client that is lying

The User-Agent and Client Hints are two independent statements. When they
disagree, the disagreement is the signal, so merging records it instead of
discarding it.

```go
c := uax.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/605.1.15 " +
	"(KHTML, like Gecko) Version/17.5 Safari/605.1.15")

cons := c.Consistency()
for _, f := range cons.Findings {
	fmt.Printf("[%s] %s: %s vs %s\n", f.Severity, f.Rule, f.Left, f.Right)
}
// [critical] impossible.safari_on_windows: Safari has not shipped for Windows since 2012
```

The verdict separates two questions that get confused constantly:

```go
v := uax.VerdictFrom(c, cons)   // pass the check you already ran

v.Trust      // how much of what it says is believable
v.Risk       // how likely it is to be hostile
v.Recommend  // "allow" | "challenge" | "block"
for _, r := range v.Reasons {
	fmt.Println(r.Code, r.Detail)
}
```

> `ParseRequest` fills in `c.Verdict` for you. `Parse` does not — a verdict is a
> pointer, and the parse path allocates nothing. Use `VerdictFrom` when you
> already have the findings, and `ComputeVerdict(c)` when you do not; passing in
> the check you ran avoids running every rule twice.

**Trust and risk are not opposites.** A well-behaved crawler that names itself
is highly trusted *and* low risk — being a bot is not being a liar. A string
that contradicts itself is low trust and high risk whatever it claims to be.

Three things to be clear about before you act on this:

1. `Recommend` is **advice**, never an authorisation decision. Nothing here is
   authentication; everything reported is a claim the client made about itself.
2. A finding is evidence of an edited string, not proof of intent. Privacy
   tools, corporate proxies and old middleboxes all rewrite User-Agents.
3. Consistency is not honesty. A carefully edited string is perfectly
   consistent — which is why "no findings" is reported as no contradictions
   found, not as "this client is genuine".

---

## 7. Make it fast under load

The default is already fast and allocates nothing:

```
BenchmarkParseTo/chrome_android   1093 ns/op   0 B/op   0 allocs/op
BenchmarkParseTo/curl              251 ns/op   0 B/op   0 allocs/op
```

Cost tracks input length, not rule count — one pass over the string, one map
probe per token — so more rules never make it slower.

If the same User-Agents repeat, which on real traffic they do:

```go
p := uax.New(uax.Config{
	Cache: uax.CacheConfig{Enabled: true, Size: 8192},
})
```

A cache hit is about 23 ns. It is **opt-in on purpose**: a cache is a global
side effect on what is otherwise a pure function, it retains attacker-supplied
strings, and it makes benchmarks measure the cache instead of the parser.

For a hot path that reuses storage between requests:

```go
var c uax.Client
for _, ua := range lines {
	p.ParseTo(&c, ua)   // zero allocations, c fully overwritten each time
	...
}
```

`ParseTo` overwrites `c` completely, so a stale client can be handed back
without clearing it. But the resolved strings are views into the input you
passed — **do not keep a `Client` past the lifetime of the string it came
from**, and do not hold one across goroutines.

Skip work you do not need:

```go
p := uax.New(uax.Config{
	SkipDevice: true,   // no device resolution
	SkipBot:    true,   // no bot, AI, WebView or automation resolution
})
```

---

## 8. Test your own handlers

Hand-written User-Agent fixtures drift into strings no browser ever sent, and
then you are testing your parser's tolerance for nonsense. `uatest` builds
internally consistent ones — right engine token, plausible version, model code
matching its brand — each round-trip verified against the parser:

```go
import "github.com/bakhod1r/uax/uatest"

ua := uatest.Chrome(138, uatest.Android(15, uatest.Model("Samsung", "SM-S928B")))
r := uatest.Request(ua, uatest.WithHints("Google Chrome", 138, "Android", true)...)

rec := httptest.NewRecorder()
myHandler(rec, r)
```

Other builders: `Firefox`, `Safari`, `Edge`, `HeadlessChrome`, `WebView`, and
platforms `Windows`, `MacOS`, `Linux`, `IOS`, `IPadOS`, `Android`,
`AndroidTablet`.

---

## 9. Find out what you do not recognise

No rule database is complete. The honest question is *what is it missing on
your traffic*, and the CLI answers it against a real log:

```bash
go install github.com/bakhod1r/uax/cmd/uax@latest

uax analyze access.log     # what actually reaches your service, and how much is unrecognised
uax explain "Mozilla/5.0 (…)"   # every field, its confidence, the verdict and any contradictions
uax device SM-S928B        # resolve a bare model code
uax db                     # rule counts and the database fingerprint
```

Record the fingerprint next to anything you store:

```go
uax.Release            // "0.1.0"
uax.DatabaseVersion()  // fingerprint of the rule set
```

Rules change between releases. A stored classification is only reproducible
against the database that produced it.

In a long-running service, collect what it could not classify:

```go
rec := uax.NewRecorder()
p := uax.New(uax.Config{Recorder: rec})

// …serve traffic…

for _, u := range rec.Top(20) {
	log.Printf("unrecognised %q seen %d times", u.Shape, u.Count)
}
```

The recorder writes to memory you own, after redacting versions and opaque
identifiers — so what you get back is the *shape* of an unknown token, not the
raw string a client sent. Use `NewRecorderWithLimit(n)` to bound it yourself.

---

## What to avoid

**Do not authenticate with it.** Everything here is a claim a client made about
itself. Treating a parse result as authentication is a vulnerability in your
application.

**Do not read absence as evidence.** No automation detected does not mean human.
No bot detected does not mean person. Both mean *nothing matched*.

**Do not store `OS.Version` without checking `OS.Reduced`.** Frozen values look
exactly like real ones and will quietly poison a dashboard.

**Do not invent a device name.** If `Device.Name` is empty, the catalogue does
not have it. Fall back to `Device.Model`, which is true.

**Do not hold a `Client` longer than its input.** Fields are views into the
string you passed in.

---

## Where to go next

- [README](./) — the full reference: every field, every rule category, the
  design commitments and why they are commitments
- [CHANGELOG](CHANGELOG.md) — what changed and what it means for a caller
- [SECURITY.md](SECURITY.md) — what this package promises, and what it
  explicitly does not
- [CONTRIBUTING.md](CONTRIBUTING.md) — adding a rule, and the one rule that
  governs all of them: never invent a fact
- [pkg.go.dev](https://pkg.go.dev/github.com/bakhod1r/uax) — generated API docs
