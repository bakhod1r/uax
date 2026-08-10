---
title: Changelog
---

# Changelog

Notable changes per release. Entries say what changed for a caller, not what
was refactored.

The rule database has its own identity, reported by `uax.DatabaseVersion()`.
Rules change more often than the code does, and a classification is only
reproducible if you know which database produced it.

---

## v0.1.0 — first release

Initial public version.

### Parsing

- `Parse` reads a User-Agent and resolves browser, engine, operating system,
  architecture, device type, brand and model code.
- `ParseTo` writes into caller-owned storage and **allocates nothing**.
  `TestAllocationBudget` fails the build if that ever stops being true.
- `ParseRequest` additionally reads Client Hints, `Sec-Fetch-*` metadata and the
  `Accept` profile.
- Cost tracks input length, not database size: the engine makes one pass over
  the string and one map probe per token, so adding rules does not slow
  parsing. `BenchmarkScaling` shows the shape — 254 ns for `curl/8.4.0`,
  1097 ns for a full Chrome-on-Android string, both allocating nothing.
- `Sec-CH-UA-WoW64` and `Sec-CH-UA-Bitness` are read, and disagreements between
  the string and the hints are reported as `hints.arch_mismatch` and
  `hints.bitness_mismatch`. WoW64 does not excuse a bitness conflict: it
  describes a 32-bit process, while bitness describes the machine.

### Classification

- Bots are categorised — SearchEngine, AICrawler, AIAgent, SocialCrawler, SEO,
  Monitoring, Security — with the operating provider recorded.
- In-app WebViews, automation frameworks and HTTP client libraries are
  distinguished from browsers rather than lumped in with them.
- `Consistency` reports internal contradictions as findings with stable rule
  identifiers, and `Verdict` separates trust from risk with an advisory action.
- `VerdictFrom` takes a coherence check the caller already ran. A verdict is
  built on the findings, so `ComputeVerdict` has to produce them; a caller that
  displays the findings too would otherwise run every spoof rule twice.
  `ParseRequest` uses it, which is worth about 20% of its cost and one
  allocation per request.

### Client Hints

- Full `Sec-CH-UA*` support with GREASE brands rejected.
- `Sec-CH-UA-Form-Factors` resolves the form factor outright — including `XR`,
  `EInk`, `Watch` and `Automotive`, none of which the User-Agent can express.
  Values outside the specification are recorded and never interpreted.
- Hints win over the User-Agent, but disagreements are kept as findings rather
  than discarded.

### User-Agent Reduction

- Frozen platform values are reported as frozen, not repeated as facts.
  A reduced `Android 10; K` yields `OS.Reduced` with **no version** and no model,
  because the values are placeholders that contradict reality. macOS `10_15_7`
  and Windows `NT 10.0` are floors, so those versions are kept and flagged.
  Client Hints clear the flag.

### Devices

This module names no hardware. `SM-S928B` does not contain "Galaxy S24 Ultra" —
a manufacturer assigned that mapping — so answering it needs a catalogue, and
catalogues move on a schedule a parser does not share: phones ship weekly.

- `Device.Model` is the code the request carried, read positionally out of the
  Android platform comment. No table of model codes is involved, because the
  field *is* the code, whatever it happens to be.
- Everything past that — name, brand, family, form factor — comes from
  `Config.Device`, a hook whose signature carries no types from either side, so
  neither package imports the other:

  ```go
  import "github.com/bakhod1r/devicex"

  p := uax.New(uax.Config{Device: devicex.Describe})
  ```

- Without the hook, the model still resolves and the rest stays empty. That is
  the honest result of having nothing to answer it with, and it means a binary
  that does not want twenty-five thousand device names does not carry them.

### Integration

- Middleware for `net/http`, Gin, Echo and Fiber. It is lazy: a route that never
  calls `FromContext` pays nothing.
- Optional cache, metrics hook, `slog` attributes, unknown-token recorder and a
  runtime rule registry with atomic swap.
- `uatest` builds internally consistent fixtures, round-trip verified against
  the parser.

### CLI

- `uax parse`, `explain`, `json`, `analyze`, `device`, `db`.
- `uax robots` generates a robots.txt from the bot database. A test asserts
  every emitted token parses back to the bot it claims to block, so the file
  cannot drift away from the parser.

### MCP server

- `cmd/uax-mcp` exposes the parser over the Model Context Protocol, speaking
  JSON-RPC 2.0 on stdin and stdout. Four tools: `parse_user_agent`,
  `check_consistency`, `lookup_device`, `database_info`.
- It exists because a model asked to read a User-Agent will guess, and the
  strings are built to defeat guessing — every browser claims `Mozilla/5.0`,
  and `SM-S928B` contains no hint that it means Galaxy S24 Ultra. Answers carry
  the confidence and the evidence behind them, and a field that did not resolve
  comes back empty instead of filled in.
- A string that identifies nothing, and a model code the catalogue does not
  hold, are returned as tool errors — distinct from JSON-RPC errors, which mean
  the call itself was malformed. Collapsing the two would tell a client its
  request was broken when the answer was simply "no".
- The protocol is implemented in the command, not pulled in: a server parsing a
  wire format is the last place an unaudited dependency belongs.

### Rules

187 rules, each shipping with the test that proves it — 55 browsers, 44
operating systems, 7 engines, 41 bots, 19 HTTP clients, 12 in-app WebViews,
9 automation frameworks.

They live in `internal/table`, in Go, and are reviewed as code. A lint suite
holds what a generator used to: a rule cannot exist without a test, a priority,
a confidence in (0,1], exactly one match clause, or an ID its category already
uses, and each category stays ordered by descending priority because the
resolver's early exit depends on it.

### Known limitations
- Device marketing names depend on whatever catalogue is wired into
  `Config.Device`. `devicex` covers 25,475 codes as of its 2026-08-08 snapshot;
  a handset newer than that resolves a brand from the code's shape and no name.
- Verified bots (reverse DNS and published address ranges), version freshness
  and rule tracing are designed but not built. Nothing in this release claims
  them.
