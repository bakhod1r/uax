package uax

import "github.com/bakhod1r/devicex"

// The device catalogue is a test-only dependency of this package.
//
// uax names no hardware: it extracts the model code and asks Config.Device
// what it means, so the library itself imports no catalogue and a caller who
// does not want twelve thousand device names does not link them. The tests are
// the one place that needs both halves, because the thing worth proving is
// that they fit together — that a real User-Agent, run through the real
// parser, wired to a real catalogue, reports the right handset.
//
// This import therefore never reaches a user of the library. It exists so the
// golden corpus can assert device facts instead of asserting emptiness.
var testParser = New(Config{Device: devicex.Describe})

// parseWithDevices is Parse plus the catalogue, for tests whose subject is the
// device rather than the parse.
func parseWithDevices(ua string) *Client { return testParser.Parse(ua) }
