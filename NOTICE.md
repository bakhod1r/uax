---
title: Third-party data
---

# Third-party data

The Go code in this repository is original, and so are the classification rules
in `internal/table`. This file records the outside material the rules read, and
the one body of data this module deliberately does not carry.

---

## Device marketing names — not in this repository

The mapping from a model code to a marketing name (`SM-G973F` → `Galaxy S10`)
is a fact assigned by the manufacturer. No rule derives it, so it can only be
imported — and imported data belongs with its provenance, on its own release
schedule, rather than inside a parser.

It therefore lives in [devicex](https://github.com/bakhod1r/devicex), which
records its own sources and reaches this module through `uax.Config.Device`.
Nothing in this repository names hardware; `Device.Model` reports the code the
request carried and nothing more.

---

## Specifications referenced by the rules

These are read, not copied. They are listed because the behaviour they define
is load-bearing and a reader should be able to check it.

| | |
|---|---|
| User-Agent Client Hints, including `Sec-CH-UA-Form-Factors` | [wicg.github.io/ua-client-hints](https://wicg.github.io/ua-client-hints/) |
| User-Agent Reduction (frozen platform values) | [chromium.org/updates/ua-reduction](https://www.chromium.org/updates/ua-reduction/) |
| Fetch Metadata request headers (`Sec-Fetch-*`) | [w3c.github.io/webappsec-fetch-metadata](https://w3c.github.io/webappsec-fetch-metadata/) |

The bot product tokens in `internal/table` are the identifiers their operators
publish for exactly this purpose, each recorded with its own test case.
