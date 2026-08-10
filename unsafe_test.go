package uax

import (
	"bytes"
	"strings"
	"testing"
	"unsafe"
)

// The parse path hands the tokeniser a []byte aliasing the caller's string,
// without copying. That is what makes a parse allocate nothing, and it is sound
// only while nothing writes through the alias: a Go string is immutable, and
// code elsewhere may be holding the same backing array, interning it, or using
// it as a map key.
//
// Nothing writes to it today. These tests exist because that is an invariant
// rather than an observation — neither the race detector nor the fuzzer can see
// it break, so a future edit to the scanner would corrupt callers' strings in
// silence. They make the aliasing observable by building the string from bytes
// the test still owns, and checking those bytes afterwards.

// stringAliasing returns a string backed by b, so a write through the parser's
// view of it is visible in b.
func stringAliasing(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func TestParseDoesNotWriteThroughItsInput(t *testing.T) {
	inputs := []string{
		"Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"curl/8.4.0",
		"",
		"(((((",
		"a\x00b",
		"Mozilla/5.0 (" + strings.Repeat("x", 200) + ")",
	}

	for _, in := range inputs {
		name := in
		if len(name) > 40 {
			name = name[:40]
		}
		t.Run(name, func(t *testing.T) {
			buf := []byte(in)
			pristine := bytes.Clone(buf)

			c := Parse(stringAliasing(buf))
			// Touch the result so nothing above can be optimised away, and so
			// the resolved views into the input are actually materialised.
			_ = c.Browser.Name + c.OS.Name + c.Device.Model

			if !bytes.Equal(buf, pristine) {
				t.Fatalf("the parser wrote through its input\n got: %q\nwant: %q", buf, pristine)
			}
		})
	}
}

func TestTruncationDoesNotWriteThroughItsInput(t *testing.T) {
	// The over-long path slices the input rather than copying it, so it is the
	// place a bounds fix is most likely to reach for a mutable buffer.
	buf := []byte("Mozilla/5.0 " + strings.Repeat("(a)", 5000))
	pristine := bytes.Clone(buf)

	c := Parse(stringAliasing(buf))
	if !c.Truncated {
		t.Fatal("input was not truncated; the test is no longer exercising the path it targets")
	}
	if !bytes.Equal(buf, pristine) {
		t.Fatal("the parser wrote through its input while truncating")
	}
}

func TestResolvedFieldsDoNotOutliveTheirBackingArrayIncorrectly(t *testing.T) {
	// Fields resolved from the input are views into it. Callers keep results
	// beyond the request, so a value that changes when the caller's buffer is
	// reused would be a fact that quietly rewrites itself.
	buf := []byte("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36")
	c := Parse(stringAliasing(buf))

	model, browser := c.Device.Model, c.Browser.Name

	// Scribble over the caller's buffer, as a pooled or reused buffer would.
	for i := range buf {
		buf[i] = 'z'
	}

	if c.Device.Model != model {
		t.Errorf("Device.Model changed when the caller's buffer was reused: %q became %q", model, c.Device.Model)
	}
	if c.Browser.Name != browser {
		t.Errorf("Browser.Name changed when the caller's buffer was reused: %q became %q", browser, c.Browser.Name)
	}
}
