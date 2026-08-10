package uax

import (
	"sort"
	"strings"
	"sync"
)

// Recorder collects product tokens that matched no rule.
//
// It answers the question every deployment eventually asks: what is my traffic
// sending that the database does not know about? Feeding those tokens back into
// internal/table is how the rules improve.
//
// Two deliberate choices make it safe to run in production:
//
//	Nothing is recorded unless Config.Recorder is set. There is no default
//	recorder, no default-on telemetry and no network access anywhere in this
//	package — the tokens go into memory the caller owns, and nowhere else.
//
//	Versions and long opaque identifiers are redacted before storage, so
//	"MyApp/3.14.159" and a per-install build ID collapse into one shape. What is
//	kept describes a class of client, not an individual one.
type Recorder struct {
	mu     sync.RWMutex
	tokens map[string]int
	limit  int
}

// DefaultRecorderLimit bounds how many distinct shapes a recorder retains.
// Without a bound, a hostile client could grow the map without limit simply by
// sending a novel token on every request.
const DefaultRecorderLimit = 4096

// NewRecorder returns a recorder holding at most DefaultRecorderLimit shapes.
func NewRecorder() *Recorder { return NewRecorderWithLimit(DefaultRecorderLimit) }

// NewRecorderWithLimit returns a recorder holding at most limit distinct shapes.
// Once full it stops adding new ones and keeps counting the ones it has, so
// memory is bounded regardless of what arrives.
func NewRecorderWithLimit(limit int) *Recorder {
	if limit <= 0 {
		limit = DefaultRecorderLimit
	}
	return &Recorder{tokens: make(map[string]int), limit: limit}
}

// Record stores the redacted shape of a token.
func (r *Recorder) Record(token string) {
	if token == "" {
		return
	}

	shape := Redact(token)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, seen := r.tokens[shape]; !seen && len(r.tokens) >= r.limit {
		return
	}
	r.tokens[shape]++
}

// Unknown is one recorded shape and how often it was seen.
type Unknown struct {
	Shape string `json:"shape"`
	Count int    `json:"count"`
}

// Top returns the n most frequent shapes, most frequent first. Ties break by
// shape so the output is stable and diffable.
func (r *Recorder) Top(n int) []Unknown {
	r.mu.RLock()
	out := make([]Unknown, 0, len(r.tokens))
	for shape, count := range r.tokens {
		out = append(out, Unknown{Shape: shape, Count: count})
	}
	r.mu.RUnlock()

	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Shape < out[j].Shape
	})

	if n > 0 && len(out) > n {
		out = out[:n]
	}
	return out
}

// Len returns how many distinct shapes are held.
func (r *Recorder) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tokens)
}

// Reset discards everything recorded so far.
func (r *Recorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens = make(map[string]int)
}

// Redact reduces a token to its shape: digit runs become "N", and long opaque
// identifiers become "X".
//
//	"MyApp/3.14.159"      -> "MyApp/N.N.N"
//	"Build/a91f0c2e77b14" -> "Build/X"
//
// The point is to keep what identifies a kind of client while dropping what
// could identify an installation, so the result is safe to log, aggregate and
// share upstream as a rule request.
func Redact(token string) string {
	var b strings.Builder
	b.Grow(len(token))

	for i := 0; i < len(token); {
		c := token[i]

		if c >= '0' && c <= '9' {
			for i < len(token) && token[i] >= '0' && token[i] <= '9' {
				i++
			}
			b.WriteByte('N')
			continue
		}

		// A long run of letters and digits mixed together is an identifier,
		// not a product name.
		if isAlnum(c) {
			start := i
			digits := 0
			for i < len(token) && isAlnum(token[i]) {
				if token[i] >= '0' && token[i] <= '9' {
					digits++
				}
				i++
			}
			word := token[start:i]
			if len(word) >= 8 && digits >= 2 {
				b.WriteByte('X')
			} else {
				b.WriteString(word)
			}
			continue
		}

		b.WriteByte(c)
		i++
	}
	return b.String()
}

func isAlnum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
