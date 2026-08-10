package uax

import (
	"strings"

	"github.com/bakhod1r/uax/internal/scan"
	"github.com/bakhod1r/uax/internal/table"
)

// The resolver reads a User-Agent exactly once and resolves every category in
// a single sweep.
//
// A conventional parser loops over N rules and runs each against the string,
// costing O(rules × len) and getting slower with every rule added. Here the
// input is tokenised once into a fixed stack buffer, and each resulting token
// is looked up once in a combined index that answers for all eight categories
// at once. Cost tracks input length, not rule count.
//
// Rules that cannot be indexed — prefix and substring forms — are evaluated
// afterwards, and only while they could still outrank what already matched.

// maxTokens caps how much of a single User-Agent the resolver will consider.
// Real strings hold well under 40 tokens; the cap keeps the working set on the
// stack and bounds the work an adversarial input can cause.
const maxTokens = 32

// Rule categories. A dense enum lets a match result be a small array indexed by
// category rather than eight variables threaded through by hand.
type category uint8

const (
	catBrowserIdx category = iota
	catEngineIdx
	catOSIdx
	catBotIdx
	catAppIdx
	catAutomationIdx
	catClientIdx
	numCategories
)

// span is a range of the input. Storing offsets rather than string headers
// keeps the per-parse buffer small — which matters because Go zeroes a local
// struct on declaration, and an array of string headers would cost more to
// zero than the parse itself costs to run.
type span struct {
	off, length uint16
}

func (s span) of(raw string) string {
	if s.length == 0 {
		return ""
	}
	return raw[s.off : s.off+s.length]
}

// product is one "Name/1.2.3" pair, as offsets into the input.
type product struct {
	name    span
	version span
}

// tokens is the once-per-parse view of a User-Agent: the products outside the
// parentheses, and the semicolon-separated fields inside them.
type tokens struct {
	raw string

	prod  [maxTokens]product
	nprod int

	field  [maxTokens]span
	nfield int
	// head[i] is field[i] cut at its first '/' or space, precomputed because
	// every category would otherwise recompute it.
	head [maxTokens]span

	// bigrams is a 256-bit presence filter over adjacent byte pairs. It costs
	// one pass over the input and lets a substring rule be rejected with a
	// single bit test instead of a scan.
	//
	// Pairs rather than single bytes because single bytes are far too common to
	// discriminate: nearly every User-Agent contains a lowercase "b", so a
	// byte filter would let "bot" through every time, while the pair "bo" is
	// genuinely rare. Collisions only cost a scan that would have happened
	// anyway, so the filter may be approximate but is never wrong.
	bigrams [4]uint64

	// Signals worth knowing before any rule runs.
	hasMobileToken bool // the bare "Mobile" product, as in "Mobile Safari"
	hasTabletToken bool
	hasTVToken     bool
	hasWebViewFlag bool // the "wv" field Android WebView adds
}

// Accessors. Each returns a view into the input; none copies.
func (t *tokens) prodName(i int) string    { return t.prod[i].name.of(t.raw) }
func (t *tokens) prodVersion(i int) string { return t.prod[i].version.of(t.raw) }
func (t *tokens) fieldAt(i int) string     { return t.field[i].of(t.raw) }
func (t *tokens) headAt(i int) string      { return t.head[i].of(t.raw) }

// collect tokenises the input. This is the only pass over the raw bytes.
func collect(t *tokens, raw string) {
	t.raw = raw
	t.nprod = 0
	t.nfield = 0
	t.hasMobileToken = false
	t.hasTabletToken = false
	t.hasTVToken = false
	t.hasWebViewFlag = false
	t.bigrams = [4]uint64{}

	for i := 0; i+1 < len(raw); i++ {
		h := bigramHash(raw[i], raw[i+1])
		t.bigrams[h>>6] |= 1 << (h & 63)
	}

	s := scan.NewScanner(unsafeBytes(raw))
	for s.Next() {
		tok := s.Token()
		switch tok.Kind {
		case scan.KindProduct:
			if t.nprod >= maxTokens {
				continue
			}
			name := unsafeString(tok.Raw)
			switch name {
			case "Mobile":
				t.hasMobileToken = true
			case "Tablet":
				t.hasTabletToken = true
			case "TV", "SMART-TV":
				t.hasTVToken = true
			}
			t.prod[t.nprod] = product{name: span{uint16(tok.Offset), uint16(len(tok.Raw))}}
			t.nprod++

		case scan.KindVersion:
			// A version always belongs to the product just emitted.
			if t.nprod > 0 {
				t.prod[t.nprod-1].version = span{uint16(tok.Offset), uint16(len(tok.Raw))}
			}

		case scan.KindComment:
			addFields(t, unsafeString(tok.Raw), tok.Offset)
		}
	}
}

// addFields splits a comment's contents on ';' and records each trimmed field.
// Splitting uses index arithmetic, so every field stays a view into the
// original string.
func addFields(t *tokens, comment string, base int) {
	for len(comment) > 0 {
		var f string
		start := base
		if i := strings.IndexByte(comment, ';'); i >= 0 {
			f, comment, base = comment[:i], comment[i+1:], base+i+1
		} else {
			f, comment, base = comment, "", base+len(comment)
		}

		// Trim without losing track of where the field sits in the input.
		for len(f) > 0 && f[0] == ' ' {
			f, start = f[1:], start+1
		}
		for len(f) > 0 && f[len(f)-1] == ' ' {
			f = f[:len(f)-1]
		}
		if f == "" {
			continue
		}

		switch f {
		case "wv":
			t.hasWebViewFlag = true
		case "SMART-TV", "SmartTV":
			t.hasTVToken = true
		case "Tablet":
			t.hasTabletToken = true
		case "Mobile":
			t.hasMobileToken = true
		}

		if t.nfield < maxTokens {
			t.field[t.nfield] = span{uint16(start), uint16(len(f))}
			t.head[t.nfield] = span{uint16(start), uint16(len(head(f)))}
			t.nfield++
		}
	}
}

// head returns the field up to its first '/' or space: "Opera Mini/7.5" yields
// "Opera Mini", "Tizen 6.0" yields "Tizen".
func head(f string) string {
	if i := strings.IndexByte(f, '/'); i > 0 {
		return f[:i]
	}
	if i := strings.IndexByte(f, ' '); i > 0 {
		return f[:i]
	}
	return f
}

// bigramHash maps a byte pair into the 256-slot filter.
func bigramHash(a, b byte) uint32 {
	return (uint32(a)*31 + uint32(b)) & 255
}

// mayContain reports whether the input could contain a needle with this
// bigram. A false answer is conclusive; a true answer means "scan to be sure".
func (t *tokens) mayContain(h uint32) bool {
	return t.bigrams[h>>6]&(1<<(h&63)) != 0
}

// rarestBigram picks the least likely adjacent pair of a needle, so its absence
// rejects the rule without touching the input.
func rarestBigram(s string) uint32 {
	if len(s) < 2 {
		// A one-byte needle cannot be filtered; treat it as always possible.
		return alwaysPossible
	}
	best, bestRank := bigramHash(s[0], s[1]), 1<<30
	for i := 0; i+1 < len(s); i++ {
		if r := byteRarity(s[i]) + byteRarity(s[i+1]); r < bestRank {
			best, bestRank = bigramHash(s[i], s[i+1]), r
		}
	}
	return best
}

// alwaysPossible is a sentinel bigram that never filters anything out.
const alwaysPossible = ^uint32(0)

// ------------------------------------------------------------------- index

// slot holds the highest-priority rule claiming a token, per category, so a
// single map probe answers for every category at once.
type slot [numCategories]*table.Rule

// unindexed is a prefix or substring rule, tagged with its category and with
// its rarest byte precomputed for the substring prefilter.
type unindexed struct {
	rule *table.Rule
	cat  category
	// rare is the rarest bigram of a substring rule's needle, or
	// alwaysPossible when the needle is too short to filter.
	rare uint32
}

// indexSet is one compiled generation of the whole rule database.
type indexSet struct {
	exact map[string]*slot

	// prefixByByte buckets prefix rules by their first byte, so a candidate
	// string is only tested against rules that could possibly match it. Without
	// it, every prefix rule would be tried against every token.
	prefixByByte [256][]unindexed

	// contains is sorted by descending priority, so a category can stop as soon
	// as nothing left could outrank what it already matched.
	contains []unindexed
}

// match is one category's resolved rule and the candidate string that matched.
type match struct {
	rule *table.Rule
	str  string
}

// buildIndexes compiles the rule categories into one index.
func buildIndexes(groups [numCategories][]table.Rule) *indexSet {
	ix := &indexSet{exact: make(map[string]*slot, 256)}

	for cat := category(0); cat < numCategories; cat++ {
		for i := range groups[cat] {
			r := &groups[cat][i]
			switch {
			case r.Token != "":
				s := ix.exact[r.Token]
				if s == nil {
					s = new(slot)
					ix.exact[r.Token] = s
				}
				// Rules arrive priority-sorted, so the first claim wins.
				if s[cat] == nil {
					s[cat] = r
				}

			case r.Prefix != "":
				b := r.Prefix[0]
				ix.prefixByByte[b] = append(ix.prefixByByte[b], unindexed{rule: r, cat: cat})

			case r.Contains != "":
				ix.contains = append(ix.contains, unindexed{rule: r, cat: cat, rare: rarestBigram(r.Contains)})
			}
		}
	}

	for b := range ix.prefixByByte {
		sortUnindexed(ix.prefixByByte[b])
	}
	sortUnindexed(ix.contains)
	return ix
}

// sortUnindexed orders rules by descending priority. Insertion sort suits a
// nearly-sorted slice and runs once, at build time.
func sortUnindexed(u []unindexed) {
	for i := 1; i < len(u); i++ {
		for j := i; j > 0 && u[j].rule.Priority > u[j-1].rule.Priority; j-- {
			u[j], u[j-1] = u[j-1], u[j]
		}
	}
}

// matchAll resolves every category in one sweep.
func (ix *indexSet) matchAll(t *tokens) [numCategories]match {
	var best [numCategories]match

	consider := func(s *slot, str string) {
		for cat := category(0); cat < numCategories; cat++ {
			r := s[cat]
			if r == nil {
				continue
			}
			if b := best[cat].rule; b == nil || r.Priority > b.Priority {
				best[cat] = match{rule: r, str: str}
			}
		}
	}

	// Tiers 1 and 2 share one loop over the candidates: each string is probed
	// once in the exact index, then tested against the prefix rules that start
	// with its first byte. Both are driven by the input rather than by the rule
	// list, which is what keeps cost independent of database size.
	tryPrefixes := func(cand string) {
		for _, u := range ix.prefixByByte[cand[0]] {
			if b := best[u.cat].rule; b != nil && u.rule.Priority <= b.Priority {
				continue
			}
			if len(cand) >= len(u.rule.Prefix) && strings.HasPrefix(cand, u.rule.Prefix) {
				best[u.cat] = match{rule: u.rule, str: cand}
			}
		}
	}

	for i := 0; i < t.nprod; i++ {
		name := t.prodName(i)
		if name == "" {
			continue
		}
		if s := ix.exact[name]; s != nil {
			consider(s, name)
		}
		tryPrefixes(name)
	}
	for i := 0; i < t.nfield; i++ {
		f := t.fieldAt(i)
		if f == "" {
			continue
		}
		if s := ix.exact[f]; s != nil {
			consider(s, f)
		} else if h := t.headAt(i); h != f {
			if s := ix.exact[h]; s != nil {
				consider(s, f)
			}
		}
		tryPrefixes(f)
	}

	// Tier 3: substring rules over the raw input — the slowest form, so it runs
	// last, behind a bit test that rejects most rules without touching a byte.
	for _, u := range ix.contains {
		if b := best[u.cat].rule; b != nil && u.rule.Priority <= b.Priority {
			continue
		}
		if u.rare != alwaysPossible && !t.mayContain(u.rare) {
			continue
		}
		if strings.Contains(t.raw, u.rule.Contains) {
			best[u.cat] = match{rule: u.rule, str: u.rule.Contains}
		}
	}

	return best
}

// hasPrefixAny reports the first comment field or product token carrying the
// given prefix. The first-byte test keeps the common miss to a few
// instructions.
func hasPrefixAny(t *tokens, prefix string) (string, bool) {
	c := prefix[0]
	for i := 0; i < t.nfield; i++ {
		if f := t.fieldAt(i); len(f) >= len(prefix) && f[0] == c && strings.HasPrefix(f, prefix) {
			return f, true
		}
	}
	for i := 0; i < t.nprod; i++ {
		if p := t.prodName(i); len(p) >= len(prefix) && p[0] == c && strings.HasPrefix(p, prefix) {
			return p, true
		}
	}
	return "", false
}

// byteRarity ranks how likely a byte is to appear in a User-Agent. Lower is
// rarer. The ranking only has to be roughly right: it picks which byte of a
// substring rule to test first, and a poor pick costs a scan that would have
// happened anyway.
func byteRarity(c byte) int {
	switch {
	case c == ' ' || c == '/' || c == '.' || c == ';' || c == ',' || c == '(' || c == ')':
		return 120
	case c >= 'a' && c <= 'z':
		return 100
	case c >= '0' && c <= '9':
		return 90
	case c >= 'A' && c <= 'Z':
		return 40
	default:
		return 10 // punctuation and anything unusual
	}
}
