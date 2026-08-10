package uax

import "strings"

// Device resolution.
//
// This package does not name hardware. "SM-S928B" does not contain "Galaxy S24
// Ultra" and no rule derives it — the mapping was assigned by a manufacturer,
// so answering it takes a catalogue, and a catalogue is on a release schedule
// this parser is not: phones ship weekly.
//
// So the work here stops at the evidence. The model code a request carries is
// extracted from the User-Agent, and everything beyond it — the marketing name,
// the brand, the form factor — is asked of Config.Device, which the caller
// supplies. github.com/bakhod1r/devicex provides one:
//
//	p := uax.New(uax.Config{Device: devicex.Describe})
//
// With no hook configured, Device.Model still reports the code. The rest stays
// empty, which is the honest answer when nothing can supply it.

// resolveDevice extracts the device evidence and asks the configured hook what
// it means.
func (p *Parser) resolveDevice(t *tokens) Device {
	code := modelCode(t)

	d := Device{}
	if code != "" {
		d.Model = code
		d.Source = SourceUA
		// The code is copied out of the request, so it is certain as evidence
		// even when nothing can name it. A hook that recognises it replaces
		// this with its own, better-founded confidence.
		d.Confidence = 1
	}

	if p.cfg.Device == nil {
		return d
	}

	name, brand, model, family, deviceType, confidence, ok := p.cfg.Device(code, t.raw)
	if !ok {
		return d
	}

	d.Name = name
	d.Brand = brand
	d.Family = family
	d.Type = DeviceType(deviceType)
	d.Confidence = float32(confidence)
	d.Source = SourceUA

	// The hook reports the identifier the request carried. For Android that is
	// the code already extracted; for hardware that publishes none it is the
	// token the User-Agent named itself with — "iPhone" is the whole
	// identifier an iPhone sends. A hook that has nothing to point at leaves it
	// empty, and "Macintosh" is exactly that case: it names a class of machine
	// and no model, so Device.Model stays empty rather than repeating a
	// marketing name as though it were evidence.
	if model != "" {
		d.Model = model
	} else if code == "" {
		d.Model = ""
	}
	return d
}

// modelCode is the model code the User-Agent carries, or "" when it carries
// none.
//
// Android publishes it positionally rather than as a labelled token: the
// platform comment is "(Linux; Android 15; SM-S928B)", and the code is whatever
// follows the Android version. There is no rule to match here — the field is
// the code, whatever it happens to be — which is exactly why a table of twelve
// thousand model codes was the wrong way to read it.
//
// Apple publishes none at all. An iPhone says "iPhone", and that is the whole
// identifier; it is handled by the hook reading the User-Agent, not here.
func modelCode(t *tokens) string {
	for i := 0; i < t.nfield; i++ {
		f := t.fieldAt(i)
		if !strings.HasPrefix(f, "Android") {
			continue
		}
		// "Android 15" is the version field; the code is the next one. A
		// trailing field is absent on a UA that names no device, and on the
		// reduced form Chrome now sends, which says "Android 10; K".
		if i+1 >= t.nfield {
			return ""
		}
		return trimBuild(t.fieldAt(i + 1))
	}
	return ""
}

// trimBuild strips the build suffix Android appends inside the model field, as
// in "SM-A546E Build/UP1A.231005.007". The build is firmware provenance, not
// hardware identity, and no catalogue is keyed by it.
func trimBuild(field string) string {
	if i := strings.Index(field, " Build/"); i > 0 {
		return field[:i]
	}
	return field
}
