---
title: Contributing
---

# Contributing to uax

Thank you for considering it. This document is short because most of what
matters is enforced by tests rather than by convention — if it builds and the
suite passes, it is close to mergeable.

## The one rule that matters

**Never invent a fact.** This package reports what a request said about itself.
Where it cannot tell, it reports nothing: an unrecognised client is `Unknown`,
not "Desktop by elimination"; an unknown model code yields an empty
`Device.Name`, not a plausible-looking guess. A wrong answer is worse than no
answer here, because a wrong answer is a fact-shaped value that a caller will
store, aggregate and report.

Everything else in this file is a consequence of that.

## Contributing a rule

Rules live in `internal/table`, in Go, and are reviewed as code. Add the rule
to its category in `rules.go`, keeping the slice ordered by descending
priority — the resolver stops at the first match, so order decides which rule
wins:

```go
{ID: "my_browser", Name: "My Browser", Family: "Chrome", Engine: "Blink",
    Token: "MyBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.95},
```

Then add the test that proves it, in `tests.go`, using a **real** User-Agent —
one you have observed, not one you composed to match the rule:

```go
{RuleID: "my_browser", UA: "Mozilla/5.0 (…) MyBrowser/3.2 Safari/537.36",
    Expect: map[string]string{"browser": "My Browser"}},
```

`go test ./...` then checks, and refuses to pass otherwise, that the rule has a
test, a priority, a confidence in (0,1], exactly one match clause and an ID its
category does not already use, and that the category is still ordered.

Choosing a priority: brand-specific tokens sit above generic ones. A Chromium
fork must outrank the bare `Chrome` token (140), a Gecko fork must outrank
`Firefox` (150), and `Safari` (100) is last because almost every WebKit browser
also says `Safari`.

Set `Confidence` to what the match really proves. An exact product token is
strong; a substring that could appear in prose is not.

## What does not belong here

**Device marketing names.** `SM-S928B` does not contain "Galaxy S24 Ultra" — a
manufacturer assigned it — so it is imported data with a provenance, not a
rule. It lives in [devicex](https://github.com/bakhod1r/devicex) and reaches
this package through `Config.Device`.

**A dependency in the library.** `uax` sits in front of every request, so a
transitive dependency is a supply-chain risk for every user. CI fails if the
library package links anything external. `cmd/uax` and the framework adapters
are separate and may depend on what they need.

**Anything that allocates on the parse path.** `TestAllocationBudget` fails the
build if `ParseTo` ever allocates. Resolved fields are views into the caller's
input, not copies.

## Before opening a pull request

```bash
make all      # lint and the full suite, including the framework adapters
make race     # the parser is used concurrently by definition
make bench    # if you touched the parse path
```

If you changed anything under `internal/scan` or touched `unsafeBytes`, read
[SECURITY.md](SECURITY.md) first. The tokeniser reads the caller's string in
place, and `unsafe_test.go` is what stands between that and memory corruption.

## Bug reports

Include the exact input. For this package that is a User-Agent string or a set
of request headers, which is enough to reproduce anything it does. A parse that
returns the wrong answer is more useful to us as one line of input than as a
description.

For a security issue, do not open a public issue —
see [SECURITY.md](SECURITY.md).

## Style

Match the surrounding code. Comments explain *why* a thing is the way it is,
because the *what* is already in the code; a comment that restates the line
below it will be asked about in review.
