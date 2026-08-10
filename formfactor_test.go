package uax

import (
	"reflect"
	"testing"
)

func TestParseFormFactors(t *testing.T) {
	if got := parseFormFactors("  "); got != nil {
		t.Errorf("expected nil for spaces, got %v", got)
	}
	if got := parseFormFactors(`"Mobile", "Tablet"`); !reflect.DeepEqual(got, []string{"Mobile", "Tablet"}) {
		t.Errorf("expected Mobile, Tablet, got %v", got)
	}
	if got := parseFormFactors(`""`); got != nil {
		t.Errorf("expected nil for empty unquoted, got %v", got)
	}
}

func TestDeviceTypeFromFormFactors(t *testing.T) {
	dt, ok := deviceTypeFromFormFactors([]string{"Mobile"})
	if !ok || dt != DeviceMobile {
		t.Errorf("expected Mobile, ok, got %v %v", dt, ok)
	}
	dt, ok = deviceTypeFromFormFactors([]string{"UnknownFactor"})
	if ok || dt != DeviceUnknown {
		t.Errorf("expected Unknown, !ok, got %v %v", dt, ok)
	}
}
