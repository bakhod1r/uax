package main

import "encoding/json"

// The Model Context Protocol wire format: JSON-RPC 2.0 over newline-delimited
// JSON on stdin and stdout.
//
// The protocol is implemented here rather than pulled in as a dependency. It is
// a few hundred lines of message shuffling, and this repository's rule is that
// the library links nothing external — a rule worth keeping in the commands
// too, since a server that speaks a wire format is exactly the place where an
// unaudited dependency would sit closest to untrusted input.

// protocolVersion is the MCP revision this server implements. It is echoed in
// the initialize result so a client that speaks a different revision can say so
// rather than discovering the mismatch one malformed message at a time.
const protocolVersion = "2024-11-05"

// request is an incoming JSON-RPC message.
//
// ID is json.RawMessage because JSON-RPC allows a string, a number or null, and
// the only thing this server does with it is echo it back. Parsing it into a
// concrete type would mean choosing one, and choosing wrongly for some client.
type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// isNotification reports whether the message expects no reply. A JSON-RPC
// notification is a request without an ID, and replying to one is a protocol
// error rather than a harmless extra.
func (r *request) isNotification() bool { return len(r.ID) == 0 }

type response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSON-RPC 2.0 error codes. Only the ones this server can actually produce are
// defined; an unused constant is a claim about behaviour that does not exist.
const (
	codeParseError     = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
)

// ---------------------------------------------------------------- initialize

type initializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    capabilities `json:"capabilities"`
	ServerInfo      serverInfo   `json:"serverInfo"`
}

// capabilities advertises what this server offers. Tools only: it exposes no
// resources and no prompts, and says so by omitting them rather than declaring
// them empty.
type capabilities struct {
	Tools *toolsCapability `json:"tools,omitempty"`
}

// toolsCapability is empty because the tool list is fixed at build time. A
// server whose tools can change advertises listChanged here; claiming it
// without sending the notification would be a lie a client acts on.
type toolsCapability struct{}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ------------------------------------------------------------------ tools

type tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type toolsListResult struct {
	Tools []tool `json:"tools"`
}

type toolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

// toolResult is what a tool call returns.
//
// IsError marks a failure the model should see and reason about — a
// User-Agent that resolves to nothing, an unknown model code — as distinct
// from a JSON-RPC error, which reports that the call itself was malformed.
// Collapsing the two would tell a client its request was broken when the
// answer was simply "no".
type toolResult struct {
	Content []content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func textResult(s string) *toolResult {
	return &toolResult{Content: []content{{Type: "text", Text: s}}}
}

func errorResult(s string) *toolResult {
	return &toolResult{Content: []content{{Type: "text", Text: s}}, IsError: true}
}
