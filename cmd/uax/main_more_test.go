package main

import (
	"bytes"
	"testing"
)

func TestMainMore(t *testing.T) {
	var out bytes.Buffer

	_ = run([]string{"analyze", "main.go"}, &out)
	_ = run([]string{"json", "Mozilla/5.0"}, &out)
	_ = run([]string{"explain", "Mozilla/5.0"}, &out)
	_ = run([]string{"device", "unknown"}, &out)
	_ = run([]string{"robots", "--sitemap"}, &out)
	_ = run([]string{"robots", "--deny-provider", "OpenAI", "--allow", "search"}, &out)
	_ = run([]string{"robots", "--deny", "seo", "--sitemap", "https://e.com"}, &out)
}
