package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{"empty", []string{}, false, "uax — client intelligence"},
		{"bare ua", []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"}, false, "Chrome"},
		{"parse", []string{"parse", "curl/8.4.0"}, false, "curl"},
		{"parse no arg", []string{"parse"}, true, ""},
		{"explain", []string{"explain", "curl/8.4.0"}, false, "Tool"},
		{"explain no arg", []string{"explain"}, true, ""},
		{"json", []string{"json", "curl/8.4.0"}, false, "curl"},
		{"json no arg", []string{"json"}, true, ""},
		{"db", []string{"db"}, false, "model codes"},
		{"device", []string{"device", "CPH2451"}, false, "OnePlus"},
		{"device no arg", []string{"device"}, true, ""},
		{"robots default", []string{"robots"}, false, "User-agent: GPTBot"},
		{"robots by category", []string{"robots", "--deny", "seo"}, false, "AhrefsBot"},
		{"robots by provider", []string{"robots", "--deny-provider", "OpenAI"}, false, "OAI-SearchBot"},
		{"robots allow", []string{"robots", "--allow", "search"}, false, "User-agent: *"},
		{"robots sitemap", []string{"robots", "--sitemap", "https://e.com/s.xml"}, false, "Sitemap: https://e.com/s.xml"},
		{"robots unknown category", []string{"robots", "--deny", "nonsense"}, true, ""},
		{"robots missing deny", []string{"robots", "--deny"}, true, ""},
		{"robots missing allow", []string{"robots", "--allow"}, true, ""},
		{"robots missing deny provider", []string{"robots", "--deny-provider"}, true, ""},
		{"robots missing path", []string{"robots", "--path"}, true, ""},
		{"robots missing sitemap", []string{"robots", "--sitemap"}, true, ""},
		{"version", []string{"version"}, false, "uax 0."},
		{"version database fingerprint", []string{"version"}, false, "database "},
		{"invalid command", []string{"unknown_cmd"}, true, ""},
		{"parse full", []string{"parse", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"}, false, "Chrome"},
		{"parse bot", []string{"parse", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"}, false, "Googlebot"},
		{"parse apple", []string{"parse", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X)"}, false, "iPhone"},
		{"parse unknown device type", []string{"parse", "SmartTV"}, false, ""},
		{"explain bot", []string{"explain", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"}, false, "Googlebot"},
		{"explain apple", []string{"explain", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X)"}, false, "iPhone"},
		{"explain full device", []string{"explain", "Mozilla/5.0 (Linux; Android 10; SM-G981B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.162 Mobile Safari/537.36"}, false, "Galaxy S20 5G"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(tc.args, &out)
			if tc.wantErr && err == nil {
				t.Errorf("want error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("want nil, got %v", err)
			}
			if !tc.wantErr && tc.wantOut != "" && !strings.Contains(out.String(), tc.wantOut) {
				t.Errorf("output missing %q: %s", tc.wantOut, out.String())
			}
		})
	}
}

func TestRunAnalyze(t *testing.T) {
	tmp, _ := os.CreateTemp("", "uax-analyze-*")
	defer os.Remove(tmp.Name())
	tmp.WriteString("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36\n")
	tmp.Close()

	var out bytes.Buffer
	err := run([]string{"analyze", tmp.Name()}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Chrome") {
		t.Errorf("analyze output missing Chrome: %s", out.String())
	}
}

func TestRunAnalyzeNoArg(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"analyze"}, &out)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPct(t *testing.T) {
	if pct(5, 10) != 50.0 {
		t.Error("pct error")
	}
	if pct(0, 0) != 0.0 {
		t.Error("pct zero error")
	}
}

func TestIsCommand(t *testing.T) {
	if !isCommand("parse") || isCommand("not_a_cmd") {
		t.Error("isCommand error")
	}
}

func TestPrintTop(t *testing.T) {
	var out bytes.Buffer
	printTop(&out, "Label", map[string]int{"A": 5, "B": 5, "C": 3, "D": 1, "E": 4, "F": 10, "G": 11, "H": 12, "I": 13, "J": 14, "K": 15}, 100)
	if !strings.Contains(out.String(), "Label") {
		t.Error("printTop missing label")
	}
	var outEmpty bytes.Buffer
	printTop(&outEmpty, "Empty", map[string]int{}, 10)
	if outEmpty.Len() > 0 {
		t.Error("should be empty")
	}
}

func TestMainCmd(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	oldExit := osExit
	defer func() { osExit = oldExit }()

	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}

	// Capture stderr
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Test case where run() fails
	os.Args = []string{"uax", "unknown-command"}
	main()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if exitCode != 1 {
		t.Errorf("expected exitCode 1, got %d", exitCode)
	}
	if !bytes.Contains(buf.Bytes(), []byte("unknown command")) {
		t.Errorf("expected unknown command output, got %s", buf.String())
	}
}
