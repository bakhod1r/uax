package uax

import "strings"

// mergeHints folds Client Hints into a Client and reports every place the two
// sources disagree.
//
// Hints win. They are structured, they are what the browser vendor intends
// servers to read, and they survive the User-Agent reduction that is steadily
// emptying the legacy string. But the disagreement is never discarded: a
// User-Agent claiming Chrome 120 while the hints say Chrome 138 is the single
// clearest spoofing signal available, and silently overwriting one with the
// other would destroy it.
func mergeHints(c *Client) []Finding {
	h := c.Hints
	if h == nil {
		return nil
	}

	var findings []Finding

	// ---- browser -----------------------------------------------------------
	if b, ok := h.Brand(); ok {
		hintName := normalizeBrand(b.Brand)

		switch c.Browser.Name {
		case "":
			c.Browser.Name = hintName
			c.Browser.Version = b.Version
			c.Browser.Confidence = 0.99
			c.Browser.Source = SourceHints

		case hintName:
			c.Browser.Source = SourceBoth
			c.Browser.Confidence = 0.99
			// A full version list is more precise than the User-Agent's.
			if len(h.FullVersionList) > 0 && b.Version.Known() {
				c.Browser.Version = b.Version
			}
			// Same browser, different major version: the string was edited.
			if uaV := c.Browser.Version; uaV.Known() && b.Version.Known() &&
				len(h.FullVersionList) == 0 && uaV.Major != b.Version.Major {
				findings = append(findings, Finding{
					Rule:     "hints.browser_version_mismatch",
					Severity: SeverityHigh,
					Detail:   "the User-Agent and Sec-CH-UA report different browser versions",
					Left:     c.Browser.Name + " " + uaV.Raw,
					Right:    b.Brand + " " + b.Version.Raw,
				})
			}

		default:
			findings = append(findings, Finding{
				Rule:     "hints.browser_mismatch",
				Severity: SeverityHigh,
				Detail:   "the User-Agent and Sec-CH-UA name different browsers",
				Left:     c.Browser.Name,
				Right:    hintName,
			})
			c.Browser.Name = hintName
			c.Browser.Version = b.Version
			c.Browser.Confidence = 0.99
			c.Browser.Source = SourceHints
		}
	}

	// ---- platform ----------------------------------------------------------
	if h.Platform != "" {
		switch {
		case c.OS.Name == "":
			c.OS.Name = h.Platform
			c.OS.Confidence = 0.99
			c.OS.Source = SourceHints

		case sameOS(c.OS.Name, h.Platform):
			c.OS.Source = SourceBoth
			c.OS.Confidence = 0.99

		default:
			findings = append(findings, Finding{
				Rule:     "hints.platform_mismatch",
				Severity: SeverityCritical,
				Detail:   "the User-Agent and Sec-CH-UA-Platform name different operating systems",
				Left:     c.OS.Name,
				Right:    h.Platform,
			})
			c.OS.Name = h.Platform
			c.OS.Confidence = 0.99
			c.OS.Source = SourceHints
		}
	}

	// The platform version is the hint worth asking for: it is the only way to
	// tell Windows 11 from Windows 10, and the only source of a real Android
	// version now that the User-Agent reports a frozen "10".
	if h.PlatformVersion != "" {
		if v := ParseVersion(h.PlatformVersion); v.Known() {
			if c.OS.Name == "Windows" && v.Major >= 13 {
				// Chromium reports Windows 11 as platform version 13 or above.
				v = Version{Major: 11, KnownVersion: true, Raw: h.PlatformVersion}
			}
			c.OS.Version = v
			c.OS.Source = SourceHints
			// The frozen value has been replaced by the observed one, so the
			// warning that came with it no longer applies.
			c.OS.Reduced = false
		}
	}

	// ---- architecture ------------------------------------------------------
	//
	// The hint wins, because it is measured rather than parsed out of a string
	// the browser froze. Where the two disagree, the disagreement is recorded
	// before the hint overwrites it: a client whose User-Agent claims one CPU
	// and whose hints claim another has said something impossible about itself.
	uaArch := c.OS.Arch

	if h.Arch != "" {
		if uaArch != "" && uaArch != h.Arch {
			findings = append(findings, Finding{
				Rule:     "hints.arch_mismatch",
				Severity: SeverityMedium,
				Detail:   "the User-Agent and Sec-CH-UA-Arch report different CPU architectures",
				Left:     uaArch,
				Right:    h.Arch,
			})
		}
		c.OS.Arch = h.Arch
	}

	if h.Bitness != "" {
		// Sec-CH-UA-WoW64 does not excuse this. It says the browser process is
		// 32-bit; Sec-CH-UA-Bitness says how wide the machine is. A User-Agent
		// that named a 64-bit CPU cannot be running on a 32-bit platform.
		if want := bitnessOf(uaArch); want != "" && want != h.Bitness {
			findings = append(findings, Finding{
				Rule:     "hints.bitness_mismatch",
				Severity: SeverityMedium,
				Detail:   "the User-Agent's architecture and Sec-CH-UA-Bitness imply different platform widths",
				Left:     want,
				Right:    h.Bitness,
			})
		}
		c.OS.Bitness = h.Bitness
	}

	// ---- device ------------------------------------------------------------
	if h.Model != "" {
		if c.Device.Model != "" && c.Device.Model != h.Model {
			findings = append(findings, Finding{
				Rule:     "hints.model_mismatch",
				Severity: SeverityMedium,
				Detail:   "the User-Agent and Sec-CH-UA-Model report different device models",
				Left:     c.Device.Model,
				Right:    h.Model,
			})
		}
		c.Device.Model = h.Model
		c.Device.Source = SourceHints
		c.Device.Confidence = 0.99
		c.Device.Reduced = false
	}

	// ---- form factor -------------------------------------------------------
	//
	// Sec-CH-UA-Form-Factors states the shape of the device. The User-Agent can
	// only imply it — a tablet is recognised by the absence of a token — so the
	// hint replaces the inference whenever it carries a value this package
	// understands.
	if ff, ok := deviceTypeFromFormFactors(h.FormFactors); ok {
		// Both headers come from the same browser and describe the same
		// machine. A browser that contradicts itself was edited by something
		// in between.
		if h.Mobile && ff != DeviceMobile && ff != DeviceWatch {
			findings = append(findings, Finding{
				Rule:     "hints.form_factor_mismatch",
				Severity: SeverityMedium,
				Detail:   "Sec-CH-UA-Mobile and Sec-CH-UA-Form-Factors describe different devices",
				Left:     "mobile",
				Right:    strings.Join(h.FormFactors, ", "),
			})
		}
		c.Device.Type = ff
		c.Device.Source = SourceHints
		if c.Device.Confidence < 0.99 {
			c.Device.Confidence = 0.99
		}
	}

	if h.Mobile && c.Device.Type == DeviceDesktop {
		findings = append(findings, Finding{
			Rule:     "hints.mobile_mismatch",
			Severity: SeverityMedium,
			Detail:   "Sec-CH-UA-Mobile is ?1 but the User-Agent describes a desktop",
			Left:     string(c.Device.Type),
			Right:    "mobile",
		})
	}

	return findings
}

// normalizeBrand maps hint brand names onto the names this package uses.
// Chromium reports "Google Chrome" and "Microsoft Edge" where the User-Agent
// says "Chrome" and "Edge".
func normalizeBrand(b string) string {
	switch b {
	case "Google Chrome":
		return "Chrome"
	case "Microsoft Edge":
		return "Edge"
	case "Opera GX":
		return "Opera GX"
	case "Yandex Browser":
		return "Yandex Browser"
	default:
		return b
	}
}

// sameOS reports whether a User-Agent OS name and a hint platform describe the
// same system. They differ in granularity: the hints say "Linux" where the
// User-Agent may say "Ubuntu", and "Android" covers Android TV.
func sameOS(uaName, platform string) bool {
	if uaName == platform {
		return true
	}
	switch platform {
	case "Linux":
		switch uaName {
		case "Linux", "Ubuntu", "Debian", "Fedora", "Arch Linux":
			return true
		}
	case "Android":
		return uaName == "Android" || uaName == "Android TV"
	case "macOS":
		return uaName == "macOS"
	case "Windows":
		return uaName == "Windows"
	}
	return false
}

// bitnessOf is the platform width an architecture implies, or "" when the
// architecture says nothing about it.
func bitnessOf(arch string) string {
	switch arch {
	case "amd64", "arm64":
		return "64"
	case "x86", "arm":
		return "32"
	default:
		return ""
	}
}
