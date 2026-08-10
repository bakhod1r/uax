package uax

import "testing"

func TestParseFetchContext(t *testing.T) {
	fc := ParseFetchContext("same-origin", "navigate", "document", "?1")
	if fc == nil {
		t.Fatal("expected FetchContext, got nil")
	}
	if !fc.IsNavigation() {
		t.Error("expected IsNavigation() to be true")
	}
	if fc.Site != "same-origin" {
		t.Errorf("expected Site=same-origin, got %q", fc.Site)
	}

	fcNil := ParseFetchContext("", "", "", "")
	if fcNil != nil {
		t.Error("expected nil FetchContext for empty strings")
	}

	fcImage := ParseFetchContext("cross-site", "no-cors", "image", "")
	if fcImage.IsNavigation() {
		t.Error("expected Image fetch to not be Navigation")
	}
}
