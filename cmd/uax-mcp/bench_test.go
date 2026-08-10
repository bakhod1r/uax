package main

import (
	"io"
	"strings"
	"testing"
)

// The server is not on a hot path — one call per model turn, against a model
// that took a second to decide to make it. These exist to catch a regression
// that changes the order of magnitude, not to shave nanoseconds: a tool that
// suddenly allocates a megabyte per call is a bug worth seeing.

func benchServe(b *testing.B, line string) {
	b.Helper()

	// All iterations go through a single serve call, which is how a real client
	// uses it: one process, many messages. Restarting serve per iteration would
	// charge every message for the 64KB scanner buffer allocated once at
	// startup, and report a per-call cost dominated by a cost nobody pays.
	//
	// The measurement still covers the whole path — decode, dispatch, the tool,
	// and the encode back out.
	in := strings.NewReader(strings.Repeat(line+"\n", b.N))

	b.ReportAllocs()
	b.ResetTimer()

	if err := serve(in, io.Discard); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkParseUserAgent(b *testing.B) {
	benchServe(b, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"parse_user_agent",`+
		`"arguments":{"user_agent":"`+chromeAndroid+`"}}}`)
}

func BenchmarkCheckConsistency(b *testing.B) {
	benchServe(b, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"check_consistency",`+
		`"arguments":{"user_agent":"`+chromeWindows+`"}}}`)
}

func BenchmarkLookupDevice(b *testing.B) {
	benchServe(b, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"lookup_device",`+
		`"arguments":{"code":"SM-S928B"}}}`)
}

// BenchmarkToolsList is the one call every client makes, once, before it can do
// anything else. It rebuilds the schemas each time, which is fine at handshake
// frequency and worth knowing if that ever changes.
func BenchmarkToolsList(b *testing.B) {
	benchServe(b, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
}
