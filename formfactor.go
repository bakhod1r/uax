package uax

import "strings"

// formFactorType maps a specification value onto a DeviceType.
//
// A value that is not in this table is stored and never interpreted. New form
// factors will be added to the specification, and guessing what one means is
// how a parser starts reporting things that are not true.
var formFactorType = map[string]DeviceType{
	"Desktop":    DeviceDesktop,
	"Automotive": DeviceCar,
	"Mobile":     DeviceMobile,
	"Tablet":     DeviceTablet,
	"XR":         DeviceXR,
	"EInk":       DeviceEInk,
	"Watch":      DeviceWatch,
}

// formFactorRank orders the values by how much each one says about the shape of
// the device, most specific first.
//
// The list can hold several values, and they are not equally informative.
// "EInk" describes a screen technology rather than a shape, so an e-ink tablet
// is a Tablet; "EInk" only decides the type when nothing else does. "Desktop"
// ranks last because it is the value a browser falls back to.
var formFactorRank = []string{"XR", "Automotive", "Watch", "Tablet", "Mobile", "EInk", "Desktop"}

// parseFormFactors reads the structured-header list. Unknown values are kept
// verbatim so a caller can see what arrived.
func parseFormFactors(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := unquote(p); s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// deviceTypeFromFormFactors returns the type the hint establishes, and whether
// any value was recognised at all.
func deviceTypeFromFormFactors(factors []string) (DeviceType, bool) {
	for _, want := range formFactorRank {
		for _, got := range factors {
			if got == want {
				return formFactorType[want], true
			}
		}
	}
	return DeviceUnknown, false
}
