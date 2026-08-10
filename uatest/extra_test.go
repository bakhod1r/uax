package uatest

import (
	"testing"
)

func TestExtra(t *testing.T) {
	_ = Windows(10)
	_ = Model("Apple", "iPhone")
	_ = WithHints("Chrome", 100, "Windows", false)
}
