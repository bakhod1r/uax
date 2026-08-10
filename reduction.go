package uax

// User-Agent Reduction.
//
// Chrome froze the platform section of its User-Agent in 2022 and Chromium
// forks followed. The string still looks like it carries a version and a model,
// but the values are constants:
//
//	Mozilla/5.0 (Linux; Android 10; K) ...        ← any Android, any handset
//	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ...   ← any macOS ≥ 10.15.7
//	Mozilla/5.0 (Windows NT 10.0; Win64; x64) ... ← Windows 10 or 11
//
// Reporting those as observed facts is the single most common way a parser
// lies. The two shapes are not equally wrong, and the difference is preserved:
//
//	Android 10 with the "K" placeholder contradicts reality — the device may be
//	on Android 15 — so the version is dropped rather than repeated.
//
//	macOS and Windows state a floor. The real system is that version or newer,
//	which is still usable information, so the value is kept and flagged.
//
// The flag describes the User-Agent only. Client Hints carry the true values,
// and merging them clears it; see merge.go.

// reductionPlaceholder is what Chromium sends in place of the model.
const reductionPlaceholder = "K"

// applyReduction flags frozen platform values after the fields have been
// resolved. It reads the collected tokens, not the raw string.
func applyReduction(t *tokens, c *Client) {
	switch {
	case c.OS.Family == "Android" && c.OS.Version.Major == 10 && hasField(t, reductionPlaceholder):
		// The version and the model are both placeholders. Neither is evidence
		// of anything, so neither is reported.
		c.OS.Reduced = true
		c.OS.Version = Version{}
		c.Device.Reduced = true
		c.Device.Model = ""
		c.Device.Name = ""

	case c.OS.Name == "macOS" && c.OS.Version.Major == 10 && c.OS.Version.Minor == 15 && c.OS.Version.Patch == 7:
		// Indistinguishable from a genuine Catalina, which is the point: the
		// value is untrustworthy either way, and saying so is the honest
		// reading rather than a guess about which case this is.
		c.OS.Reduced = true

	case c.OS.Name == "Windows" && c.OS.Version.Major == 10 && c.OS.Version.Minor == 0:
		// Windows 11 reports NT 10.0. Only Sec-CH-UA-Platform-Version tells
		// them apart, so the major version is a floor rather than a release.
		c.OS.Reduced = true
	}
}

func hasField(t *tokens, want string) bool {
	for i := 0; i < t.nfield; i++ {
		if t.fieldAt(i) == want {
			return true
		}
	}
	return false
}
