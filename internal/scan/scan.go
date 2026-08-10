package scan

type Kind int

const (
	KindProduct Kind = iota
	KindVersion
	KindComment
	KindSpace
)

// RawToken represents a zero-allocation view into the input string.
type RawToken struct {
	Kind   Kind
	Raw    []byte
	Offset int
}

// Scanner is an allocation-free tokenizer for User-Agent strings.
type Scanner struct {
	data []byte
	pos  int
	tok  RawToken
}

// NewScanner creates a new scanner over the given byte slice.
func NewScanner(data []byte) *Scanner {
	return &Scanner{data: data, pos: 0}
}

// Next advances the scanner to the next token. It returns true if a token is found.
func (s *Scanner) Next() bool {
	if s.pos >= len(s.data) {
		return false
	}

	start := s.pos
	b := s.data[s.pos]

	if b == ' ' {
		s.pos++
		for s.pos < len(s.data) && s.data[s.pos] == ' ' {
			s.pos++
		}
		s.tok = RawToken{
			Kind:   KindSpace,
			Raw:    s.data[start:s.pos],
			Offset: start,
		}
		return true
	}

	if b == '(' {
		s.pos++
		contentStart := s.pos
		depth := 1
		for s.pos < len(s.data) {
			if s.data[s.pos] == '(' {
				depth++
			} else if s.data[s.pos] == ')' {
				depth--
				if depth == 0 {
					break
				}
			}
			s.pos++
		}
		end := s.pos
		if s.pos < len(s.data) {
			s.pos++ // consume ')'
		}
		s.tok = RawToken{
			Kind:   KindComment,
			Raw:    s.data[contentStart:end],
			Offset: contentStart,
		}
		return true
	}

	// A token is a product name or the version that follows one, and the
	// separator is what distinguishes them: "Chrome/138.0" is two tokens, and
	// only the second is a version. Read to the next '/', ' ' or '('.
	isVersion := start > 0 && s.data[start-1] == '/'

	for s.pos < len(s.data) {
		c := s.data[s.pos]
		if c == ' ' || c == '(' {
			break
		}
		if !isVersion && c == '/' {
			break
		}
		s.pos++
	}

	if isVersion {
		s.tok = RawToken{
			Kind:   KindVersion,
			Raw:    s.data[start:s.pos],
			Offset: start,
		}
		return true
	}

	s.tok = RawToken{
		Kind:   KindProduct,
		Raw:    s.data[start:s.pos],
		Offset: start,
	}

	if s.pos < len(s.data) && s.data[s.pos] == '/' {
		s.pos++ // consume '/'
	}

	return true
}

// Token returns the current token.
func (s *Scanner) Token() RawToken {
	return s.tok
}
