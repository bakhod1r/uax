package uax

import "testing"

// The contract this package keeps about hardware: it reports the code a request
// carried, and nothing else unless a catalogue is wired in. These tests hold
// both halves — what works without a hook, and what the hook adds.

func TestModelCodeIsReadWithoutACatalogue(t *testing.T) {
	// No Config.Device, so nothing can name the handset. The code is still
	// evidence the request carried, and dropping it would lose a fact that was
	// there in the input.
	tests := []struct {
		name string
		ua   string
		want string
	}{
		{
			name: "android publishes the code positionally",
			ua:   "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36",
			want: "SM-S928B",
		},
		{
			name: "the build suffix is firmware, not hardware",
			ua:   "Mozilla/5.0 (Linux; Android 13; SM-A546E Build/UP1A.231005.007) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
			want: "SM-A546E",
		},
		{
			name: "a code absent from any catalogue is still reported",
			ua:   "Mozilla/5.0 (Linux; Android 15; SM-Z999X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36",
			want: "SM-Z999X",
		},
		{
			name: "apple publishes no code",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
			want: "",
		},
		{
			name: "android as the very last token field yields no model",
			ua:   "Mozilla/5.0 (Linux; Android)",
			want: "",
		},
		{
			name: "a desktop carries no model at all",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Parse(tc.ua).Device.Model; got != tc.want {
				t.Errorf("Device.Model = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestReducedUserAgentYieldsNoModel(t *testing.T) {
	// "K" is Chromium's placeholder, not a handset. Reporting it would be the
	// exact lie User-Agent Reduction was introduced to make possible.
	c := Parse("Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36")
	if c.Device.Model != "" {
		t.Errorf("Device.Model = %q, want empty for the reduction placeholder", c.Device.Model)
	}
	if !c.Device.Reduced {
		t.Error("Device.Reduced = false, want true")
	}
}

func TestNothingIsNamedWithoutAHook(t *testing.T) {
	c := Parse("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36")
	if c.Device.Name != "" || c.Device.Brand != "" {
		t.Errorf("named hardware with no catalogue configured: brand=%q name=%q", c.Device.Brand, c.Device.Name)
	}
	// The form factor is not a catalogue fact — the User-Agent says "Mobile" —
	// so it survives having no hook.
	if c.Device.Type != DeviceMobile {
		t.Errorf("Device.Type = %q, want Mobile", c.Device.Type)
	}
}

func TestHookAnswersWhatTheRulesCannot(t *testing.T) {
	tests := []struct {
		name                          string
		ua                            string
		brand, model, devName, family string
		devType                       DeviceType
	}{
		{
			name:  "catalogued handset",
			ua:    "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36",
			brand: "Samsung", model: "SM-S928B", devName: "Galaxy S24 Ultra", devType: DeviceMobile,
		},
		{
			name:  "uncatalogued code still resolves its maker from the shape",
			ua:    "Mozilla/5.0 (Linux; Android 15; SM-Z999X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36",
			brand: "Samsung", model: "SM-Z999X", devName: "", devType: DeviceMobile,
		},
		{
			name:  "a tablet code is a tablet even when unnamed",
			ua:    "Mozilla/5.0 (Linux; Android 14; SM-X999T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			brand: "Samsung", model: "SM-X999T", family: "Galaxy Tab", devType: DeviceTablet,
		},
		{
			name:  "apple names itself, and the name is the whole identifier",
			ua:    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
			brand: "Apple", model: "iPhone", devName: "iPhone", devType: DeviceMobile,
		},
		{
			name:  "a Mac names a class of machine and no model",
			ua:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			brand: "Apple", model: "", devName: "Mac", devType: DeviceDesktop,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := parseWithDevices(tc.ua).Device
			if d.Brand != tc.brand {
				t.Errorf("Brand = %q, want %q", d.Brand, tc.brand)
			}
			if d.Model != tc.model {
				t.Errorf("Model = %q, want %q", d.Model, tc.model)
			}
			if d.Name != tc.devName {
				t.Errorf("Name = %q, want %q", d.Name, tc.devName)
			}
			if tc.family != "" && d.Family != tc.family {
				t.Errorf("Family = %q, want %q", d.Family, tc.family)
			}
			if d.Type != tc.devType {
				t.Errorf("Type = %q, want %q", d.Type, tc.devType)
			}
		})
	}
}

func TestHookIsNeverAskedToGuess(t *testing.T) {
	// A hook that declines must leave the extracted evidence untouched rather
	// than blanking it: "I cannot name this" is not "there was no code".
	p := New(Config{
		Device: func(code, ua string) (name, brand, model, family, deviceType string, confidence float64, ok bool) {
			return "", "", "", "", "", 0, false
		},
	})
	d := p.Parse("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36").Device
	if d.Model != "SM-S928B" {
		t.Errorf("Model = %q, want the extracted code to survive a declining hook", d.Model)
	}
	if d.Name != "" || d.Brand != "" {
		t.Errorf("a declining hook produced brand=%q name=%q, want both empty", d.Brand, d.Name)
	}
}

func TestSkipDeviceCostsTheModelToo(t *testing.T) {
	p := New(Config{SkipDevice: true, Device: nil})
	if got := p.Parse("Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36").Device.Model; got != "" {
		t.Errorf("Device.Model = %q with SkipDevice, want empty", got)
	}
}
