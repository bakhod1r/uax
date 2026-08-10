package uax

import (
	"testing"
)

func TestVersion_AtLeast(t *testing.T) {
	v := Version{Major: 15, Minor: 4, Patch: 1}

	if !v.AtLeast(15) {
		t.Error("expected true for 15")
	}
	if !v.AtLeast(15, 4) {
		t.Error("expected true for 15.4")
	}
	if !v.AtLeast(15, 4, 1) {
		t.Error("expected true for 15.4.1")
	}
	if !v.AtLeast(15, 4, 0) {
		t.Error("expected true for 15.4.0")
	}
	if v.AtLeast(16) {
		t.Error("expected false for 16")
	}
	if v.AtLeast(15, 5) {
		t.Error("expected false for 15.5")
	}
	if v.AtLeast(15, 4, 2) {
		t.Error("expected false for 15.4.2")
	}
}

func TestVersion_Compare(t *testing.T) {
	v1 := Version{Major: 1, Minor: 2, Patch: 3}

	if got := v1.Compare(Version{Major: 1, Minor: 2, Patch: 3}); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
	if got := v1.Compare(Version{Major: 1, Minor: 2, Patch: 2}); got <= 0 {
		t.Errorf("expected >0, got %d", got)
	}
	if got := v1.Compare(Version{Major: 1, Minor: 3, Patch: 0}); got >= 0 {
		t.Errorf("expected <0, got %d", got)
	}
}

func TestVersion_String(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3, Raw: "1.2.3-beta"}
	if got := v.String(); got != "1.2.3-beta" {
		t.Errorf("expected '1.2.3-beta', got %q", got)
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input string
		want  Version
	}{
		{"", Version{KnownVersion: false}},
		{"1", Version{Major: 1, KnownVersion: true, Raw: "1"}},
		{"1.2", Version{Major: 1, Minor: 2, KnownVersion: true, Raw: "1.2"}},
		{"1.2.3", Version{Major: 1, Minor: 2, Patch: 3, KnownVersion: true, Raw: "1.2.3"}},
		{"1.2.3.4", Version{Major: 1, Minor: 2, Patch: 3, KnownVersion: true, Raw: "1.2.3.4"}},
		{"1.2.3-beta", Version{Major: 1, Minor: 2, Patch: 3, KnownVersion: true, Raw: "1.2.3-beta"}},
		{"a.b.c", Version{KnownVersion: false, Raw: "a.b.c"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseVersion(tt.input)
			if got.Major != tt.want.Major || got.Minor != tt.want.Minor || got.Patch != tt.want.Patch || got.KnownVersion != tt.want.KnownVersion {
				t.Errorf("ParseVersion(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}
