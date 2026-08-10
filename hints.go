package uax

import (
	"net/http"
	"strings"
)

// Hints is the parsed User-Agent Client Hints set.
//
// Client Hints matter because the legacy User-Agent string is being frozen and
// hollowed out: modern Chromium reports "Android 10; K" for every phone
// regardless of the real OS version or model. The hints carry the truth, but
// only the low-entropy ones are sent by default — the rest arrive only after
// the server requests them with Accept-CH.
type Hints struct {
	// Brands is the Sec-CH-UA list with GREASE entries removed. It contains
	// significant-version numbers only ("138"), unless FullVersionList was
	// requested.
	Brands []BrandHint `json:"brands,omitempty"`

	// FullVersionList is Sec-CH-UA-Full-Version-List: the same brands carrying
	// complete versions ("138.0.7204.92"). Requires Accept-CH.
	FullVersionList []BrandHint `json:"full_version_list,omitempty"`

	Mobile          bool   `json:"mobile"`
	Platform        string `json:"platform,omitempty"`         // "Windows", "macOS", "Android"
	PlatformVersion string `json:"platform_version,omitempty"` // Requires Accept-CH
	Model           string `json:"model,omitempty"`            // Requires Accept-CH
	Arch            string `json:"arch,omitempty"`             // Requires Accept-CH
	Bitness         string `json:"bitness,omitempty"`          // Requires Accept-CH

	// WoW64 is Sec-CH-UA-WoW64: a 32-bit browser running on 64-bit Windows.
	// It describes the process, not the machine — a WoW64 client still reports
	// Bitness "64", because that is what the platform is. Requires Accept-CH.
	WoW64 bool `json:"wow64,omitempty"`

	// FormFactors is Sec-CH-UA-Form-Factors: what shape the device physically
	// is, stated rather than inferred. Requires Accept-CH.
	//
	// Values are kept exactly as sent, including ones this package does not
	// recognise, so a hint added to the specification after this release is
	// still visible to callers.
	FormFactors []string `json:"form_factors,omitempty"`
}

// BrandHint is one entry of a Sec-CH-UA brand list.
type BrandHint struct {
	Brand   string  `json:"brand"`
	Version Version `json:"version"`
}

// Header names for every hint this package understands.
const (
	HeaderCHUA                = "Sec-CH-UA"
	HeaderCHUAMobile          = "Sec-CH-UA-Mobile"
	HeaderCHUAPlatform        = "Sec-CH-UA-Platform"
	HeaderCHUAPlatformVersion = "Sec-CH-UA-Platform-Version"
	HeaderCHUAModel           = "Sec-CH-UA-Model"
	HeaderCHUAArch            = "Sec-CH-UA-Arch"
	HeaderCHUABitness         = "Sec-CH-UA-Bitness"
	HeaderCHUAWoW64           = "Sec-CH-UA-WoW64"
	HeaderCHUAFullVersionList = "Sec-CH-UA-Full-Version-List"
	HeaderCHUAFormFactors     = "Sec-CH-UA-Form-Factors"
)

// AcceptCH is the header value asking a browser to send the high-entropy hints
// this package can use. Serve it as `Accept-CH: <AcceptCH>` to receive the OS
// version and device model that a modern User-Agent no longer contains.
const AcceptCH = "Sec-CH-UA-Platform-Version, Sec-CH-UA-Model, Sec-CH-UA-Arch, Sec-CH-UA-Bitness, Sec-CH-UA-Full-Version-List, Sec-CH-UA-Form-Factors"

// ParseHints reads every Client Hint from a header set. It returns nil when the
// request carries no hints at all, so "no hints" is distinguishable from "hints
// with empty values".
func ParseHints(h http.Header) *Hints {
	hints := &Hints{
		Brands:          parseBrandList(h.Get(HeaderCHUA)),
		FullVersionList: parseBrandList(h.Get(HeaderCHUAFullVersionList)),
		Mobile:          h.Get(HeaderCHUAMobile) == "?1",
		Platform:        normalizePlatform(unquote(h.Get(HeaderCHUAPlatform))),
		PlatformVersion: unquote(h.Get(HeaderCHUAPlatformVersion)),
		Model:           unquote(h.Get(HeaderCHUAModel)),
		Arch:            normalizeArch(unquote(h.Get(HeaderCHUAArch)), unquote(h.Get(HeaderCHUABitness))),
		Bitness:         unquote(h.Get(HeaderCHUABitness)),
		WoW64:           h.Get(HeaderCHUAWoW64) == "?1",
		FormFactors:     parseFormFactors(h.Get(HeaderCHUAFormFactors)),
	}

	if len(hints.Brands) == 0 && len(hints.FullVersionList) == 0 &&
		hints.Platform == "" && hints.PlatformVersion == "" &&
		hints.Model == "" && hints.Arch == "" && len(hints.FormFactors) == 0 &&
		h.Get(HeaderCHUAMobile) == "" {
		return nil
	}
	return hints
}

// ParseClientHints parses the three low-entropy hints a browser sends without
// being asked. Prefer ParseHints, which reads the complete set.
func ParseClientHints(secChUA, secChUAMobile, secChUAPlatform string) *Hints {
	if secChUA == "" && secChUAMobile == "" && secChUAPlatform == "" {
		return nil
	}
	return &Hints{
		Brands:   parseBrandList(secChUA),
		Mobile:   secChUAMobile == "?1",
		Platform: normalizePlatform(unquote(secChUAPlatform)),
	}
}

// parseBrandList parses the structured-header list browsers send, for example:
//
//	"Chromium";v="138", "Google Chrome";v="138", "Not/A)Brand";v="24"
//
// GREASE entries are deliberate garbage: browsers inject a randomly-shaped
// fake brand on every release specifically so that servers cannot hardcode the
// list. Treating one as a real browser is the classic mistake here, so they are
// dropped.
func parseBrandList(v string) []BrandHint {
	if v == "" {
		return nil
	}

	var out []BrandHint
	for len(v) > 0 {
		var item string
		if i := strings.IndexByte(v, ','); i >= 0 {
			item, v = v[:i], v[i+1:]
		} else {
			item, v = v, ""
		}

		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		brand, version := item, ""
		if i := strings.Index(item, ";v="); i >= 0 {
			brand, version = item[:i], item[i+3:]
		}

		brand = unquote(strings.TrimSpace(brand))
		if brand == "" || isGREASE(brand) {
			continue
		}

		out = append(out, BrandHint{Brand: brand, Version: ParseVersion(unquote(strings.TrimSpace(version)))})
	}
	return out
}

// isGREASE reports whether a brand is the intentional decoy browsers include.
// The specification does not fix its spelling — it has shipped as "Not;A Brand",
// "Not)A;Brand", "Not/A)Brand" and more — so the test is structural: a real
// brand name does not contain punctuation of this kind.
func isGREASE(brand string) bool {
	// Every GREASE entry shipped so far is built from the word "Brand" with
	// noise around it, and no real browser brand contains that word.
	if strings.Contains(brand, "Brand") {
		return true
	}
	// Belt and braces for future shapes: a real brand is not punctuation soup.
	return strings.HasPrefix(brand, "Not") && strings.ContainsAny(brand, ";)/(:=?_")
}

// Brand returns the significant brand: the entry that is not Chromium and not
// GREASE, which is the browser the user actually installed.
func (h *Hints) Brand() (BrandHint, bool) {
	if h == nil {
		return BrandHint{}, false
	}

	list := h.FullVersionList
	if len(list) == 0 {
		list = h.Brands
	}

	for _, b := range list {
		if b.Brand != "Chromium" {
			return b, true
		}
	}
	if len(list) > 0 {
		return list[0], true
	}
	return BrandHint{}, false
}

// unquote strips the quoting used by structured headers.
func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// normalizePlatform maps hint spellings onto the names used elsewhere in this
// package, so a caller can compare OS.Name with Hints.Platform directly.
func normalizePlatform(p string) string {
	switch p {
	case "macOS":
		return "macOS"
	case "Chrome OS", "Chromium OS":
		return "Chrome OS"
	case "Unknown":
		return ""
	default:
		return p
	}
}

// normalizeArch maps the hint's architecture and bitness onto Go's GOARCH
// vocabulary, which is what OS.Arch reports.
func normalizeArch(arch, bitness string) string {
	switch arch {
	case "x86":
		if bitness == "64" {
			return "amd64"
		}
		return "x86"
	case "arm":
		if bitness == "64" {
			return "arm64"
		}
		return "arm"
	case "arm64":
		return "arm64"
	default:
		return arch
	}
}
