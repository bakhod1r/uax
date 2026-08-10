package main

import (
	"testing"

	"github.com/avct/uasurfer"
	"github.com/bakhod1r/uax"
	"github.com/mileusna/useragent"
	"github.com/mssola/user_agent"
)

// TestWhatEachResolves prints what every library extracts from the same input.
// A latency comparison is meaningless without it: a parser that resolves less
// is expected to be faster, and the only fair reading is cost against work.
func TestWhatEachResolves(t *testing.T) {
	c := uax.Parse(ua)
	t.Logf("uax        browser=%s %s | os=%s %s | device=%s %s %q | type=%s",
		c.Browser.Name, c.Browser.Version.Raw, c.OS.Name, c.OS.Version.Raw,
		c.Device.Brand, c.Device.Model, c.Device.Name, c.Device.Type)

	m := useragent.Parse(ua)
	t.Logf("mileusna   browser=%s %s | os=%s %s | device=%q | mobile=%v",
		m.Name, m.Version, m.OS, m.OSVersion, m.Device, m.Mobile)

	s := user_agent.New(ua)
	bn, bv := s.Browser()
	t.Logf("mssola     browser=%s %s | os=%s | platform=%s | mobile=%v",
		bn, bv, s.OS(), s.Platform(), s.Mobile())

	u := uasurfer.Parse(ua)
	t.Logf("uasurfer   browser=%v %d | os=%v %d | device=%v",
		u.Browser.Name, u.Browser.Version.Major, u.OS.Name, u.OS.Version.Major, u.DeviceType)

	p := uapParser.Parse(ua)
	t.Logf("uap-go     browser=%s %s | os=%s %s | device=%s %s %s",
		p.UserAgent.Family, p.UserAgent.Major, p.Os.Family, p.Os.Major,
		p.Device.Brand, p.Device.Model, p.Device.Family)
}
