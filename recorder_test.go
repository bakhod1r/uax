package uax

import "testing"

func TestRedactKeepsShapeDropsIdentity(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"MyApp/3.14.159":      "MyApp/N.N.N",
		"Build/a91f0c2e77b14": "Build/X",
		"SomeBrowser":         "SomeBrowser",
		"Foo/12":              "Foo/N",
	}
	for in, want := range cases {
		if got := Redact(in); got != want {
			t.Errorf("Redact(%q) = %q, want %q", in, got, want)
		}
	}
}

// Two installs of the same app must collapse to one shape; that is the whole
// point of redacting before storing.
func TestRecorderCollapsesVersions(t *testing.T) {
	t.Parallel()

	r := NewRecorder()
	r.Record("MyApp/1.2.3")
	r.Record("MyApp/9.9.9")

	top := r.Top(10)
	if len(top) != 1 {
		t.Fatalf("distinct shapes: got %d, want 1 (%+v)", len(top), top)
	}
	if top[0].Count != 2 {
		t.Errorf("count: got %d, want 2", top[0].Count)
	}
}

// A hostile client must not be able to grow the map without limit.
func TestRecorderIsBounded(t *testing.T) {
	t.Parallel()

	r := NewRecorderWithLimit(10)
	for i := 0; i < 1000; i++ {
		r.Record(string(rune('a'+i%26)) + string(rune('a'+i/26%26)) + string(rune('a'+i/676%26)))
	}
	if r.Len() > 10 {
		t.Errorf("recorder holds %d shapes, want at most 10", r.Len())
	}
}

// Recording is opt-in, and a parser without one must never accumulate anything.
func TestRecordingIsOptIn(t *testing.T) {
	t.Parallel()

	if Default().cfg.Recorder != nil {
		t.Error("the default parser has a recorder attached")
	}

	rec := NewRecorder()
	p := New(Config{Recorder: rec})
	p.Parse("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 WeirdThing/4.2")

	found := false
	for _, u := range rec.Top(10) {
		if u.Shape == "WeirdThing" {
			found = true
		}
	}
	if !found {
		t.Errorf("the unrecognised token was not recorded: %+v", rec.Top(10))
	}
}

func TestRecorderEdgeCases(t *testing.T) {
	t.Parallel()

	// Limit <= 0 should default
	r1 := NewRecorderWithLimit(0)
	if r1.limit != DefaultRecorderLimit {
		t.Errorf("expected default limit for limit=0, got %d", r1.limit)
	}

	// Record empty string
	r1.Record("")
	if r1.Len() != 0 {
		t.Errorf("expected 0 for empty token, got %d", r1.Len())
	}

	// Reset
	r1.Record("Token")
	r1.Reset()
	if r1.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", r1.Len())
	}

	// Top with n > length
	r1.Record("A")
	r1.Record("B")
	r1.Record("C")
	top := r1.Top(10)
	if len(top) != 3 {
		t.Errorf("expected 3 items from Top(10), got %d", len(top))
	}

	// Top with n < length
	top = r1.Top(2)
	if len(top) != 2 {
		t.Errorf("expected 2 items from Top(2), got %d", len(top))
	}
}

// TestRecordSkipsKnown entities to cover ua.go record logic
func TestRecordSkipsKnownEntities(t *testing.T) {
	t.Parallel()

	rec := NewRecorder()
	p := New(Config{Recorder: rec})

	// Bot case
	p.Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	// HTTPClient case
	p.Parse("curl/7.64.1")

	for _, u := range rec.Top(10) {
		if u.Shape == "Googlebot" || u.Shape == "curl" {
			t.Errorf("expected known entities to be skipped, but recorded: %s", u.Shape)
		}
	}
}
