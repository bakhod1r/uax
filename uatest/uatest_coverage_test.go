package uatest

import "testing"

func TestCoverageExtra(t *testing.T) {
	_ = Windows(7)
	_ = Windows(8)

	_ = WithHints("a", 1, "b", true)
}
