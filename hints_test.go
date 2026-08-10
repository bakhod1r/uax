package uax

import "testing"

func TestParseBrandListDropsGREASE(t *testing.T) {
	t.Parallel()

	// Browsers inject a randomly-shaped fake brand on every release precisely
	// so servers cannot hardcode the list. Every spelling below has shipped.
	for _, header := range []string{
		`"Not/A)Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
		`"Not_A Brand";v="24", "Chromium";v="128", "Google Chrome";v="128"`,
		`" Not;A Brand";v="99", "Chromium";v="120", "Google Chrome";v="120"`,
		`"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`,
	} {
		brands := parseBrandList(header)
		if len(brands) != 2 {
			t.Errorf("got %d real brands, want 2, from %s -> %+v", len(brands), header, brands)
		}
		for _, b := range brands {
			if b.Brand != "Chromium" && b.Brand != "Google Chrome" {
				t.Errorf("GREASE brand %q survived: %s", b.Brand, header)
			}
		}
	}
}

func TestBrandPrefersTheInstalledBrowser(t *testing.T) {
	t.Parallel()

	h := &Hints{Brands: parseBrandList(`"Chromium";v="138", "Microsoft Edge";v="138"`)}
	b, ok := h.Brand()
	if !ok || b.Brand != "Microsoft Edge" {
		t.Errorf("Brand() = %+v, want Microsoft Edge", b)
	}
	if b.Version.Major != 138 {
		t.Errorf("version: got %d, want 138", b.Version.Major)
	}
}

func TestFullVersionListWins(t *testing.T) {
	t.Parallel()

	h := &Hints{
		Brands:          parseBrandList(`"Google Chrome";v="138"`),
		FullVersionList: parseBrandList(`"Google Chrome";v="138.0.7204.92"`),
	}
	b, _ := h.Brand()
	if b.Version.Patch != 7204 {
		t.Errorf("expected the full version, got %q", b.Version.Raw)
	}
}

func TestNoHintHeadersYieldsNil(t *testing.T) {
	t.Parallel()

	if h := ParseHints(request("").Header); h != nil {
		t.Errorf("ParseHints on a request with no hints = %+v, want nil", h)
	}
}

func TestArchNormalisation(t *testing.T) {
	t.Parallel()

	cases := []struct{ arch, bitness, want string }{
		{"x86", "64", "amd64"},
		{"x86", "32", "x86"},
		{"arm", "64", "arm64"},
		{"arm", "32", "arm"},
		{"arm64", "", "arm64"},
		{"unknown", "", "unknown"},
	}
	for _, tc := range cases {
		if got := normalizeArch(tc.arch, tc.bitness); got != tc.want {
			t.Errorf("normalizeArch(%q, %q) = %q, want %q", tc.arch, tc.bitness, got, tc.want)
		}
	}
}

func TestHintsEdgeCases(t *testing.T) {
	// ParseClientHints empty check
	if got := ParseClientHints("", "", ""); got != nil {
		t.Error("expected nil for empty hints")
	}
	if got := ParseClientHints(`"Chrome";v="100"`, "?1", `"Windows"`); got == nil {
		t.Error("expected valid hints")
	}

	// parseBrandList empty item
	brands := parseBrandList(`"Chromium";v="138", , "Google Chrome";v="138"`)
	if len(brands) != 2 {
		t.Errorf("expected 2 brands, got %d", len(brands))
	}

	// Brand on nil hints
	var h *Hints
	if _, ok := h.Brand(); ok {
		t.Error("expected false for nil hints")
	}

	// Brand when only Chromium is present
	h = &Hints{Brands: parseBrandList(`"Chromium";v="138"`)}
	if b, ok := h.Brand(); !ok || b.Brand != "Chromium" {
		t.Error("expected Chromium as fallback")
	}

	// normalizePlatform
	if got := normalizePlatform("Unknown"); got != "" {
		t.Error("expected empty string for Unknown platform")
	}
	if got := normalizePlatform("Chrome OS"); got != "Chrome OS" {
		t.Error("expected Chrome OS")
	}
}
