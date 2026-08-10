package scan

import (
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "empty",
			input: "",
			want:  []Token{},
		},
		{
			name:  "simple product",
			input: "Mozilla/5.0",
			want: []Token{
				{Kind: KindProduct, Value: "Mozilla", Offset: 0},
				{Kind: KindVersion, Value: "5.0", Offset: 8},
			},
		},
		{
			name:  "product with comment",
			input: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			want: []Token{
				{Kind: KindProduct, Value: "Mozilla", Offset: 0},
				{Kind: KindVersion, Value: "5.0", Offset: 8},
				{Kind: KindSpace, Value: " ", Offset: 11},
				{Kind: KindComment, Value: "Windows NT 10.0; Win64; x64", Offset: 13},
			},
		},
		{
			name:  "multiple products",
			input: "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			want: []Token{
				{Kind: KindProduct, Value: "AppleWebKit", Offset: 0},
				{Kind: KindVersion, Value: "537.36", Offset: 12},
				{Kind: KindSpace, Value: " ", Offset: 18},
				{Kind: KindComment, Value: "KHTML, like Gecko", Offset: 20},
				{Kind: KindSpace, Value: " ", Offset: 38},
				{Kind: KindProduct, Value: "Chrome", Offset: 39},
				{Kind: KindVersion, Value: "120.0.0.0", Offset: 46},
				{Kind: KindSpace, Value: " ", Offset: 55},
				{Kind: KindProduct, Value: "Safari", Offset: 56},
				{Kind: KindVersion, Value: "537.36", Offset: 63},
			},
		},
		{
			name:  "malformed no version",
			input: "ProductOnly",
			want: []Token{
				{Kind: KindProduct, Value: "ProductOnly", Offset: 0},
			},
		},
		{
			name:  "nested parens in comment",
			input: "Mozilla/5.0 (Linux; Android 10; SM-A205U) (KHTML, like Gecko)",
			want: []Token{
				{Kind: KindProduct, Value: "Mozilla", Offset: 0},
				{Kind: KindVersion, Value: "5.0", Offset: 8},
				{Kind: KindSpace, Value: " ", Offset: 11},
				{Kind: KindComment, Value: "Linux; Android 10; SM-A205U", Offset: 13},
				{Kind: KindSpace, Value: " ", Offset: 41},
				{Kind: KindComment, Value: "KHTML, like Gecko", Offset: 43},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner([]byte(tt.input))
			var got []Token
			for s.Next() {
				tok := s.Token()
				// We convert []byte to string for easier assertion, but keep in mind Token has []byte Value
				got = append(got, Token{
					Kind:   tok.Kind,
					Value:  string(tok.Raw),
					Offset: tok.Offset,
				})
			}

			// Convert expected strings to our assertable Token struct
			var wantTokens []Token
			for _, w := range tt.want {
				wantTokens = append(wantTokens, Token{
					Kind:   w.Kind,
					Value:  w.Value, // This is already string in our test struct
					Offset: w.Offset,
				})
			}

			if len(got) == 0 && len(wantTokens) == 0 {
				return
			}
			if !reflect.DeepEqual(got, wantTokens) {
				t.Errorf("Scan() got = %v, want %v", got, wantTokens)
			}
		})
	}
}

// Token is redefined here just for test assertions since we want to store string Value instead of []byte to make struct literal easier.
type Token struct {
	Kind   Kind
	Value  string
	Offset int
}

func TestScannerAllocs(t *testing.T) {
	input := []byte("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	allocs := testing.AllocsPerRun(1000, func() {
		s := NewScanner(input)
		for s.Next() {
			_ = s.Token()
		}
	})
	if allocs > 0 {
		t.Errorf("Scanner allocates %v times per run, want 0", allocs)
	}
}

// FuzzScan asserts the two properties everything above the tokeniser relies on:
// it never panics, and every token it emits is a real slice of the input at the
// offset it reported. A token that lies about its position would produce
// silently wrong versions and device models rather than an obvious failure.
func FuzzScan(f *testing.F) {
	f.Add([]byte("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 Chrome/138.0.0.0 Mobile Safari/537.36"))
	f.Add([]byte("curl/8.4.0"))
	f.Add([]byte(""))
	f.Add([]byte("((((((((("))
	f.Add([]byte(")))"))
	f.Add([]byte("a/\x00\xff/b"))

	f.Fuzz(func(t *testing.T, data []byte) {
		s := NewScanner(data)
		for n := 0; s.Next(); n++ {
			if n > len(data)+2 {
				t.Fatalf("scanner did not terminate on %q", data)
			}

			tok := s.Token()
			if tok.Offset < 0 || tok.Offset > len(data) {
				t.Fatalf("offset %d out of range for %q", tok.Offset, data)
			}
			if tok.Offset+len(tok.Raw) > len(data) {
				t.Fatalf("token overruns the input: offset %d len %d of %d", tok.Offset, len(tok.Raw), len(data))
			}
			if len(tok.Raw) > 0 && string(data[tok.Offset:tok.Offset+len(tok.Raw)]) != string(tok.Raw) {
				t.Fatalf("token at %d does not match the input it claims", tok.Offset)
			}
		}
	})
}
