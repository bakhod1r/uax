package uax

import "testing"

func TestToJSON(t *testing.T) {
	c := &Client{
		Browser: Browser{Name: "Chrome"},
	}
	b := c.ToJSON()
	if len(b) == 0 {
		t.Fatal("expected JSON bytes, got empty")
	}
}
