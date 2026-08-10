package main

import (
	"bytes"
	"testing"
)

func TestMainExtra(t *testing.T) {
	var out bytes.Buffer

	_ = summarize("", &out)
	_ = explain("", &out)
	_ = describeModel("", &out)
	_ = describeDB(&out)
	_ = analyze("nonexistent.txt", &out)

	_ = run([]string{"robots", "--deny", "seo"}, &out)
	_ = run([]string{"robots", "--deny", "seo", "--sitemap", "xml", "--path", "/api", "--allow", "search", "--deny-provider", "OpenAI"}, &out)

	_ = run([]string{"robots", "--deny", "seo", "extra"}, &out) // error case

	usage(&out)
}
