package uax

import "testing"

func TestParseAcceptProfile(t *testing.T) {
	ap := ParseAcceptProfile("text/html", "en-US", "gzip")
	if ap == nil {
		t.Fatal("expected AcceptProfile, got nil")
	}
	if ap.IsMissingLanguage() {
		t.Error("expected IsMissingLanguage() to be false")
	}

	apMissing := ParseAcceptProfile("text/html", "", "gzip")
	if !apMissing.IsMissingLanguage() {
		t.Error("expected IsMissingLanguage() to be true")
	}

	apNil := ParseAcceptProfile("", "", "")
	if apNil != nil {
		t.Error("expected nil AcceptProfile for empty strings")
	}
}

func TestPreferredLanguage(t *testing.T) {
	tests := []struct {
		name           string
		acceptLanguage string
		expected       string
	}{
		{
			name:           "basic sorting",
			acceptLanguage: "uz-UZ,uz;q=0.9,en-US;q=0.8,en;q=0.7,ru;q=0.6",
			expected:       "uz",
		},
		{
			name:           "mixed q values",
			acceptLanguage: "en;q=0.8, ru;q=0.9",
			expected:       "ru",
		},
		{
			name:           "single language",
			acceptLanguage: "en-US",
			expected:       "en",
		},
		{
			name:           "empty header",
			acceptLanguage: "",
			expected:       "",
		},
		{
			name:           "malformed parts",
			acceptLanguage: " ;q=, en;q=0.5",
			expected:       "en",
		},
		{
			name:           "only one letter language",
			acceptLanguage: "e, en;q=0.5",
			expected:       "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := &AcceptProfile{AcceptLanguage: tt.acceptLanguage}
			if got := ap.PreferredLanguage(); got != tt.expected {
				t.Errorf("PreferredLanguage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPreferredLanguageReadsQualityNotOrder(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"single tag", "en-US", "en"},
		{"the base language, not the region", "pt-BR", "pt"},
		{"q decides, not position", "en-GB;q=0.9, fr;q=1.0", "fr"},
		{"a missing q means 1.0", "de, en;q=0.9", "de"},
		{"whitespace is not significant", "  en-US ; q=0.8 ,  ja ; q=0.9  ", "ja"},
		{"an earlier entry wins a tie", "es;q=0.5, it;q=0.5", "es"},
		{"the wildcard is not a language", "*", ""},
		{"the wildcard is skipped, not fatal", "*;q=1.0, nl;q=0.5", "nl"},
		{"three-letter subtags are languages too", "haw;q=1.0", "haw"},
		{"a garbled q disqualifies its own entry only", "en;q=banana, sv;q=0.1", "sv"},
		{"q above the permitted range is garbled", "en;q=9, sv;q=0.1", "sv"},
		{"empty entries are ignored", ",, en ,,", "en"},
		{"a one-letter subtag is not a language", "x, en", "en"},
		{"nothing usable", ", ;q=0.5, *", ""},
		{"empty header", "", ""},
		{"param not starting with q= is ignored", "en;v=1;q=0.9, sv;q=0.1", "en"},
		{"no q param means 1.0", "en;v=1, sv;q=0.9", "en"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := &AcceptProfile{AcceptLanguage: tc.header}
			if got := a.PreferredLanguage(); got != tc.want {
				t.Errorf("PreferredLanguage(%q) = %q, want %q", tc.header, got, tc.want)
			}
		})
	}
}

func TestPreferredLanguageOnNilProfile(t *testing.T) {
	var a *AcceptProfile
	if got := a.PreferredLanguage(); got != "" {
		t.Errorf("PreferredLanguage on nil = %q, want empty", got)
	}
}
