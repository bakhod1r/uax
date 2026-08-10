package uax

import (
	"strconv"
	"strings"
)

// AcceptProfile represents the client's content negotiation headers.
type AcceptProfile struct {
	Accept         string `json:"accept,omitempty"`
	AcceptLanguage string `json:"accept_language,omitempty"`
	AcceptEncoding string `json:"accept_encoding,omitempty"`
}

// ParseAcceptProfile creates an AcceptProfile from the provided headers.
func ParseAcceptProfile(accept, acceptLanguage, acceptEncoding string) *AcceptProfile {
	if accept == "" && acceptLanguage == "" && acceptEncoding == "" {
		return nil
	}

	return &AcceptProfile{
		Accept:         accept,
		AcceptLanguage: acceptLanguage,
		AcceptEncoding: acceptEncoding,
	}
}

// IsMissingLanguage returns true if a language header was expected but not provided.
// Most legitimate browsers send an Accept-Language header.
func (a *AcceptProfile) IsMissingLanguage() bool {
	return a.AcceptLanguage == ""
}

// PreferredLanguage is the base language the client most wants, as a two- or
// three-letter subtag: "en-GB;q=0.9, fr;q=1.0" reports "fr".
//
// It reads the whole header rather than only its first entry, because the
// order a client writes its languages in is not the order it prefers them —
// that is what the q-value is for. Entries without one default to q=1.0, as
// RFC 9110 says, and ties keep the order the client wrote.
//
// The wildcard "*" is not a language and is never returned. A malformed
// q-value makes that one entry unusable rather than the whole header: a
// client that garbles its preferences is still a client with preferences.
func (a *AcceptProfile) PreferredLanguage() string {
	if a == nil || a.AcceptLanguage == "" {
		return ""
	}

	best := ""
	bestQ := -1.0

	for entry := range strings.SplitSeq(a.AcceptLanguage, ",") {
		tag, q := entry, 1.0
		if i := strings.IndexByte(entry, ';'); i >= 0 {
			tag = entry[:i]
			q = quality(entry[i+1:])
		}

		tag = strings.TrimSpace(tag)
		if tag == "" || tag == "*" {
			continue
		}

		// The base language, not the tag: "en-GB" and "en-US" are both "en",
		// and a caller asking which language to serve wants the language.
		if i := strings.IndexByte(tag, '-'); i > 0 {
			tag = tag[:i]
		}
		if len(tag) < 2 {
			continue
		}

		// Strictly greater, so an earlier entry wins a tie — the client's own
		// ordering is the tiebreak it expressed.
		if q > bestQ {
			best, bestQ = tag, q
		}
	}

	return best
}

// quality reads the q-value out of an Accept-Language parameter list. Anything
// it cannot read is q=0: unusable, but only for the entry that carried it.
func quality(params string) float64 {
	for param := range strings.SplitSeq(params, ";") {
		param = strings.TrimSpace(param)
		if !strings.HasPrefix(param, "q=") {
			continue
		}
		q, err := strconv.ParseFloat(strings.TrimSpace(param[2:]), 64)
		if err != nil || q < 0 || q > 1 {
			return 0
		}
		return q
	}
	return 1.0
}
