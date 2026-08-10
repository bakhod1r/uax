---
title: Security
---

# Security

## Reporting a vulnerability

Report privately through [GitHub Security
Advisories](https://github.com/bakhod1r/uax/security/advisories/new). Please do
not open a public issue for a vulnerability.

Include the input that triggers it. For this package that usually means a
User-Agent string or a set of request headers, which is enough to reproduce
anything it does.

## What this package promises

`uax` parses input that is hostile by definition: a User-Agent is a string an
unauthenticated remote client chooses freely. Three properties are treated as
security properties rather than niceties, and a failure of any of them is a
vulnerability worth reporting.

**It never panics.** Whatever it is handed — truncated UTF-8, unbalanced
parentheses, megabytes of parentheses, embedded NUL bytes, nothing at all — it
returns a result. `FuzzParse` and the tokeniser's fuzzer run in CI for this
reason. A panic reachable from `Parse`, `ParseTo` or `ParseRequest` is a denial
of service against every caller, because a parser sits in front of the
handler.

**Work is bounded by input length, not by content.** The engine makes one pass
over the string and one map probe per token, so no input makes it quadratic.
`Config.MaxInputLength` caps the input at 4096 bytes by default; longer input
is truncated and flagged as `Client.Truncated` rather than rejected, because a
parser that errors on a hostile string hands the caller a decision it cannot
make. An input that costs disproportionately more than its length suggests is a
bug.

**It does no I/O.** The library issues no network request, opens no file and
reads no environment variable. It imports `net/http` for its header and request
*types* — `ParseRequest` reads `r.Header` and nothing else — and never
constructs a client or a connection. The optional recorder writes to memory the
caller owns. If you find a code path that reaches outside the process, that is
a bug regardless of what it does.

## What this package does not promise

**It does not verify anything a client says.** Everything reported is a claim
the client made about itself. `Bot.Name` is the name in the string, not proof
of who sent the request: any client can send `Googlebot/2.1`. `Consistency`
and `Verdict` report that claims contradict each other, and `Verdict.Recommend`
is advice — never an authorisation decision. Verifying a crawler needs reverse
DNS or the operator's published address ranges, which this package does not do.

Treating a parse result as authentication is a vulnerability in the calling
application, not in this package.

**Absence of evidence is not evidence.** An automation framework configured to
hide is invisible here by design, so a nil `Automation` never means "human".

## Unsafe code

`resolve_fields.go` converts between `string` and `[]byte` without copying, so
the tokeniser reads the caller's input in place and the parse allocates
nothing. This is sound only while nothing writes through the result:
`internal/scan` reads `s.data` and never assigns to it. A change that writes to
the scanner's buffer would corrupt the caller's immutable string, and is a
memory-safety bug rather than a style question.

## Supported versions

The latest minor release. `uax.Release` reports the version, and
`uax.DatabaseVersion()` fingerprints the rule set that produced a given
classification.
