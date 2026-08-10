package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// The tests drive the server the way a client does: whole lines in, whole
// lines out. Calling handle directly would skip the framing, and the framing is
// where a protocol server actually goes wrong — a reply to a notification, a
// response without an ID, a batch of lines that desynchronises.

// run feeds lines to the server and returns the decoded replies.
func run(t *testing.T, lines ...string) []response {
	t.Helper()

	var out bytes.Buffer
	if err := serve(strings.NewReader(strings.Join(lines, "\n")+"\n"), &out); err != nil {
		t.Fatalf("serve: %v", err)
	}

	var got []response
	dec := json.NewDecoder(&out)
	for dec.More() {
		var r response
		if err := dec.Decode(&r); err != nil {
			t.Fatalf("decode reply: %v", err)
		}
		got = append(got, r)
	}
	return got
}

// only asserts a single reply came back and returns it.
func only(t *testing.T, lines ...string) response {
	t.Helper()
	got := run(t, lines...)
	if len(got) != 1 {
		t.Fatalf("got %d replies, want 1: %+v", len(got), got)
	}
	return got[0]
}

// callTool runs one tools/call and returns the tool result. A JSON-RPC error here
// is a test failure, not a result: the two are different outcomes and the tests
// that care about errors ask for them by name.
func callTool(t *testing.T, name, args string) toolResult {
	t.Helper()

	resp := only(t, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"`+name+`","arguments":`+args+`}}`)
	if resp.Error != nil {
		t.Fatalf("%s: rpc error %d: %s", name, resp.Error.Code, resp.Error.Message)
	}

	var res toolResult
	remarshal(t, resp.Result, &res)
	if len(res.Content) == 0 {
		t.Fatalf("%s: empty content", name)
	}
	return res
}

// remarshal moves a response's `any` result into a typed value. serve encodes
// through JSON, so this is the same trip a real client's decoder makes.
func remarshal(t *testing.T, from any, to any) {
	t.Helper()
	b, err := json.Marshal(from)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := json.Unmarshal(b, to); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
}

func text(res toolResult) string { return res.Content[0].Text }

// ------------------------------------------------------------------ handshake

func TestInitializeAnswersWithTheProtocolItImplements(t *testing.T) {
	resp := only(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)

	var res initializeResult
	remarshal(t, resp.Result, &res)

	if res.ProtocolVersion != protocolVersion {
		t.Errorf("protocolVersion = %q, want %q", res.ProtocolVersion, protocolVersion)
	}
	// A client that cannot see the tools capability will never call tools/list.
	if res.Capabilities.Tools == nil {
		t.Error("tools capability not advertised")
	}
	if res.ServerInfo.Name == "" || res.ServerInfo.Version == "" {
		t.Errorf("serverInfo incomplete: %+v", res.ServerInfo)
	}
}

func TestNotificationsAreNotAnswered(t *testing.T) {
	// Replying to a notification is a protocol violation, and the failure it
	// causes is desynchronisation: the client matches this reply to whatever
	// request it sends next.
	for _, method := range []string{"notifications/initialized", "notifications/cancelled", "notifications/unheard_of"} {
		if got := run(t, `{"jsonrpc":"2.0","method":"`+method+`"}`); len(got) != 0 {
			t.Errorf("%s: got %d replies, want none", method, len(got))
		}
	}
}

func TestRepliesStayInStepWhenNotificationsAreInterleaved(t *testing.T) {
	got := run(t,
		`{"jsonrpc":"2.0","id":1,"method":"ping"}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"ping"}`,
		``, // a blank line is framing noise, not a message
		`{"jsonrpc":"2.0","id":3,"method":"ping"}`,
	)
	if len(got) != 3 {
		t.Fatalf("got %d replies, want 3", len(got))
	}
	for i, want := range []string{"1", "2", "3"} {
		if string(got[i].ID) != want {
			t.Errorf("reply %d has id %s, want %s", i, got[i].ID, want)
		}
	}
}

func TestEveryReplyEchoesItsID(t *testing.T) {
	// JSON-RPC allows a string, a number or null, and a client with several
	// calls in flight matches them by this field alone.
	for _, id := range []string{`7`, `"abc"`, `null`} {
		resp := only(t, `{"jsonrpc":"2.0","id":`+id+`,"method":"ping"}`)
		if id == `null` {
			continue // null is indistinguishable from absent after encoding
		}
		if string(resp.ID) != id {
			t.Errorf("id %s echoed as %s", id, resp.ID)
		}
	}
}

// ---------------------------------------------------------------- error paths

func TestMalformedInputIsReportedWithoutAnID(t *testing.T) {
	resp := only(t, `not json at all`)
	if resp.Error == nil || resp.Error.Code != codeParseError {
		t.Fatalf("got %+v, want parse error", resp.Error)
	}
	// No ID could be read, so none may be echoed — the one case where a reply
	// legitimately carries no ID.
	if len(resp.ID) != 0 {
		t.Errorf("parse error carried id %s", resp.ID)
	}
}

func TestParseErrorDoesNotStopTheServer(t *testing.T) {
	// A client that sends one bad line keeps talking; a server that quits on it
	// takes the whole session down.
	got := run(t, `{oops`, `{"jsonrpc":"2.0","id":2,"method":"ping"}`)
	if len(got) != 2 {
		t.Fatalf("got %d replies, want 2", len(got))
	}
	if got[1].Error != nil {
		t.Errorf("request after bad line failed: %+v", got[1].Error)
	}
}

func TestProtocolErrorsUseTheirOwnCodes(t *testing.T) {
	// The codes are not decoration: a client distinguishes "you are speaking
	// the wrong protocol" from "I do not have that method" by these alone.
	for _, tc := range []struct {
		name string
		line string
		code int
	}{
		{"wrong version", `{"jsonrpc":"1.0","id":1,"method":"ping"}`, codeInvalidRequest},
		{"unknown method", `{"jsonrpc":"2.0","id":1,"method":"nope"}`, codeMethodNotFound},
		{"unknown tool", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"nope"}}`, codeInvalidParams},
		{"params not an object", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":[]}`, codeInvalidParams},
		{"missing argument", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"parse_user_agent","arguments":{}}}`, codeInvalidParams},
		{"argument wrong type", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"parse_user_agent","arguments":{"user_agent":42}}}`, codeInvalidParams},
		{"argument empty", `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"parse_user_agent","arguments":{"user_agent":""}}}`, codeInvalidParams},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := only(t, tc.line)
			if resp.Error == nil {
				t.Fatalf("no error; result was %+v", resp.Result)
			}
			if resp.Error.Code != tc.code {
				t.Errorf("code = %d, want %d (%s)", resp.Error.Code, tc.code, resp.Error.Message)
			}
			if resp.Error.Message == "" {
				t.Error("error carries no message")
			}
		})
	}
}

func TestAWrongVersionNotificationIsStillNotAnswered(t *testing.T) {
	// The version check runs before the method switch, so it is the one place
	// where a notification could pick up a reply by accident.
	if got := run(t, `{"jsonrpc":"1.0","method":"notifications/initialized"}`); len(got) != 0 {
		t.Errorf("got %d replies, want none", len(got))
	}
}

func TestEndOfInputIsNotAnError(t *testing.T) {
	// A client says it is finished by closing the pipe, not by sending a
	// shutdown message. Reporting that as a failure would make every clean exit
	// look like a crash.
	if err := serve(strings.NewReader(""), new(bytes.Buffer)); err != nil {
		t.Fatalf("empty input: %v", err)
	}
}

func TestAnOversizedLineDoesNotWedgeTheServer(t *testing.T) {
	// Bigger than the scanner's default 64KB limit, smaller than the 1MB bound.
	huge := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"parse_user_agent","arguments":{"user_agent":"` +
		strings.Repeat("A", 200*1024) + `"}}}`

	got := run(t, huge, `{"jsonrpc":"2.0","id":2,"method":"ping"}`)
	if len(got) != 2 {
		t.Fatalf("got %d replies, want 2", len(got))
	}
	if got[1].Error != nil {
		t.Errorf("request after huge line failed: %+v", got[1].Error)
	}
}

// -------------------------------------------------------------------- tools

func TestToolsListDescribesEveryToolWellEnoughToPick(t *testing.T) {
	resp := only(t, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)

	var res toolsListResult
	remarshal(t, resp.Result, &res)
	if len(res.Tools) == 0 {
		t.Fatal("no tools listed")
	}

	seen := map[string]bool{}
	for _, tl := range res.Tools {
		if seen[tl.Name] {
			t.Errorf("duplicate tool %q", tl.Name)
		}
		seen[tl.Name] = true

		// A model picks a tool from its description alone. A thin one is a tool
		// called at the wrong time.
		if len(tl.Description) < 80 {
			t.Errorf("%s: description is %d chars, too thin to choose from", tl.Name, len(tl.Description))
		}
		if tl.InputSchema["type"] != "object" {
			t.Errorf("%s: input schema is not an object", tl.Name)
		}
		if _, ok := tl.InputSchema["properties"]; !ok {
			t.Errorf("%s: input schema has no properties", tl.Name)
		}
	}

	for _, want := range []string{"parse_user_agent", "check_consistency", "lookup_device", "database_info"} {
		if !seen[want] {
			t.Errorf("tool %q not listed", want)
		}
	}
}

func TestEveryListedToolIsCallable(t *testing.T) {
	// A tool advertised but not wired up is worse than one that is absent: the
	// model spends a turn discovering it.
	args := map[string]string{
		"parse_user_agent":  `{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"}`,
		"check_consistency": `{"user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"}`,
		"lookup_device":     `{"code":"SM-S928B"}`,
		"database_info":     `{}`,
	}

	resp := only(t, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	var res toolsListResult
	remarshal(t, resp.Result, &res)

	for _, tl := range res.Tools {
		a, ok := args[tl.Name]
		if !ok {
			t.Errorf("tool %q is listed but this test has no arguments for it", tl.Name)
			continue
		}
		if got := callTool(t, tl.Name, a); got.IsError {
			t.Errorf("%s: unexpected tool error: %s", tl.Name, text(got))
		}
	}
}
