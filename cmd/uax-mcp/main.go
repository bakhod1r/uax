// Command uax-mcp exposes uax over the Model Context Protocol.
//
// It speaks JSON-RPC 2.0 on stdin and stdout, so an MCP client runs it as a
// subprocess:
//
//	{
//	  "mcpServers": {
//	    "uax": { "command": "uax-mcp" }
//	  }
//	}
//
// The tools answer questions a model cannot answer from a User-Agent string by
// reading it: which browser and operating system sent a request, whether the
// client's claims contradict each other, what a bare model code names, and what
// the rule database contains.
//
// # Why a model needs this
//
// A User-Agent is a string a model will cheerfully guess at, and guessing is
// the failure mode this package exists to prevent. "Mozilla/5.0" appears in
// every browser's string for historical reasons; "SM-S928B" contains no hint
// that it means Galaxy S24 Ultra. Answers here come from a rule database and a
// device catalogue, each field carrying the confidence and source behind it, so
// a model can tell what was established from what was merely plausible.
//
// # What it will not do
//
// Nothing here reaches the network, reads a file, or retains a request. The
// server is a pure function of its input, and a User-Agent that resolves to
// nothing is reported as nothing.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bakhod1r/devicex"
	"github.com/bakhod1r/uax"
)

// parser is wired to the device catalogue: uax extracts the model code a
// request carried, devicex names the handset.
var parser = uax.New(uax.Config{Device: devicex.Describe})

func main() {
	if err := serve(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "uax-mcp: %v\n", err)
		os.Exit(1)
	}
}

// serve reads requests until the input ends, which is how an MCP client says it
// is done: it closes the pipe rather than sending a shutdown message.
func serve(in io.Reader, out io.Writer) error {
	// A User-Agent is bounded, but a hostile or broken client is not. The
	// scanner's default 64KB line limit would turn one oversized line into a
	// permanent error, so the buffer is raised to a megabyte — generous for any
	// real message, and still a bound.
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	enc := json.NewEncoder(out)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req request
		if err := json.Unmarshal(line, &req); err != nil {
			// No ID could be parsed, so the reply carries none. This is the one
			// case where a response without an ID is correct.
			if err := enc.Encode(response{
				JSONRPC: "2.0",
				Error:   &rpcError{Code: codeParseError, Message: "invalid JSON"},
			}); err != nil {
				return err
			}
			continue
		}

		resp, reply := handle(&req)
		if !reply {
			continue
		}
		if err := enc.Encode(resp); err != nil {
			return err
		}
	}

	// A closed pipe is how the conversation ends, not a failure to report.
	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

// handle dispatches one request. The second return reports whether a reply is
// owed: notifications are acted on silently.
func handle(req *request) (response, bool) {
	if req.JSONRPC != "2.0" {
		return fail(req, codeInvalidRequest, "jsonrpc must be \"2.0\""), !req.isNotification()
	}

	switch req.Method {
	case "initialize":
		return ok(req, initializeResult{
			ProtocolVersion: protocolVersion,
			Capabilities:    capabilities{Tools: &toolsCapability{}},
			ServerInfo:      serverInfo{Name: "uax", Version: uax.Release},
		}), true

	case "notifications/initialized", "notifications/cancelled":
		// Acknowledged by doing nothing. There is no per-request state to
		// cancel: every call is answered before the next line is read.
		return response{}, false

	case "ping":
		return ok(req, struct{}{}), true

	case "tools/list":
		return ok(req, toolsListResult{Tools: tools()}), true

	case "tools/call":
		var p toolCallParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return fail(req, codeInvalidParams, "params must be an object with name and arguments"), true
		}
		result, err := call(p)
		if err != nil {
			return fail(req, codeInvalidParams, err.Error()), true
		}
		return ok(req, result), true

	default:
		// Notifications are one-way by definition, so an unknown one is
		// ignored rather than answered — replying would itself be a protocol
		// violation.
		if req.isNotification() {
			return response{}, false
		}
		return fail(req, codeMethodNotFound, "unknown method: "+req.Method), true
	}
}

func ok(req *request, result any) response {
	return response{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func fail(req *request, code int, msg string) response {
	return response{JSONRPC: "2.0", ID: req.ID, Error: &rpcError{Code: code, Message: msg}}
}
