---
title: uax
permalink: /
---

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="assets/logo-dark.png">
    <img src="assets/logo.png" alt="uax" width="340">
  </picture>
</div>

<div align="center">

[![CI](https://github.com/bakhod1r/uax/actions/workflows/ci.yml/badge.svg)](https://github.com/bakhod1r/uax/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/bakhod1r/uax.svg)](https://pkg.go.dev/github.com/bakhod1r/uax)
[![Release](https://img.shields.io/github/v/release/bakhod1r/uax?sort=semver)](https://github.com/bakhod1r/uax/releases)
[![Zero dependencies](https://img.shields.io/badge/dependencies-none-success)](https://github.com/bakhod1r/uax/blob/main/go.mod)

</div>

Client intelligence for Go. Not another User-Agent parser: `uax` reads the
User-Agent **and** Client Hints **and** the fetch metadata around them, and
tells you what is talking to your server — including when it is lying.

```bash
go get github.com/bakhod1r/uax
```

```go
c := uax.ParseRequest(r)

c.Browser.Name    // "Chrome"
c.OS.Name         // "Android"
c.Device.Model    // "SM-S928B"   — verbatim from the request
c.IsAICrawler()   // false
c.Verdict.Risk    // 0.00
```

Zero dependencies. Zero allocations on the parse path. 187 rules, each shipping
with the test that proves it.

---

## Start here

| | |
|---|---|
| **[Tutorial](TUTORIAL.md)** | Nine steps, in the order the questions come up. Every snippet is complete and runs. |
| [Reference](https://github.com/bakhod1r/uax#readme) | Every field, every rule category, the design commitments and why they are commitments |
| [API docs](https://pkg.go.dev/github.com/bakhod1r/uax) | Generated from the source |
| [Changelog](CHANGELOG.md) | What changed, and what it means for a caller |

---

## The one rule

It never invents a fact.

An unrecognised client is `Unknown`, not "Desktop by elimination". An unknown
model code leaves `Device.Name` empty rather than guessing. A frozen platform
version is reported as frozen instead of repeated as an observation.

A wrong answer is worse than no answer here, because a wrong answer is a
fact-shaped value that a caller stores, aggregates and puts in a report.

```go
c := uax.Parse("Mozilla/5.0 (Linux; Android 10; K) … Chrome/138.0.0.0 Mobile Safari/537.36")

c.OS.Name        // "Android"  — the platform is still known
c.OS.Version     // zero       — "10" is a placeholder Chrome froze, not a version
c.OS.Reduced     // true
c.Device.Model   // ""         — "K" is not a model code
```

Most parsers report `Android 10` there and are wrong on nearly every modern
Android request.

---

## What it finds

- **Browsers** — 55 rules, with family, variant and mode (headless, WebView)
- **Operating systems** — 44 rules, architecture where stated
- **Devices** — the model code verbatim; brand and marketing name through a
  catalogue hook, empty when the catalogue misses
- **Bots** — 41 rules, categorised: SearchEngine, AICrawler, AIAgent,
  SocialCrawler, SEO, Monitoring, Security
- **Not-browsers** — 12 in-app WebViews, 9 automation frameworks, 19 HTTP
  client libraries, each kept distinct from a browser
- **Client Hints** — full `Sec-CH-UA*`, GREASE rejected, disagreements with the
  User-Agent kept as findings rather than discarded
- **Consistency and verdict** — contradictions with stable rule IDs, trust
  separated from risk

---

## What it will not do

It does not verify anything a client says. Everything reported is a claim the
client made about itself: `Bot.Name` is the name in the string, not proof of who
sent the request, and `Verdict.Recommend` is advice, never an authorisation
decision.

Absence of evidence is not evidence either — automation configured to hide is
invisible here by design, so a nil `Automation` never means "human".

See [SECURITY.md](SECURITY.md) for the full statement of what is promised and
what is not.

---

<div align="center">
  <sub>MIT · <a href="https://github.com/bakhod1r/uax">github.com/bakhod1r/uax</a></sub>
</div>
