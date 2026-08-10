package scan

import "testing"

func TestCoverageExtra(t *testing.T) {
	s := NewScanner([]byte("    "))
	for s.Next() {
	}

	s = NewScanner(nil)
	for s.Next() {
	}
}
