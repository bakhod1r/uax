package uax

// Version represents a parsed version string.
type Version struct {
	Major        int    `json:"major,omitempty"`
	Minor        int    `json:"minor,omitempty"`
	Patch        int    `json:"patch,omitempty"`
	KnownVersion bool   `json:"known_version,omitempty"`
	Raw          string `json:"raw,omitempty"`
}

// Known returns true if the version was successfully parsed.
func (v Version) Known() bool {
	return v.KnownVersion
}

// String returns the raw version string.
func (v Version) String() string {
	return v.Raw
}

// Compare returns -1 if v < o, 0 if v == o, and 1 if v > o.
func (v Version) Compare(o Version) int {
	if v.Major != o.Major {
		if v.Major < o.Major {
			return -1
		}
		return 1
	}
	if v.Minor != o.Minor {
		if v.Minor < o.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != o.Patch {
		if v.Patch < o.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// AtLeast checks if the version is greater than or equal to the specified parts.
func (v Version) AtLeast(major int, rest ...int) bool {
	if v.Major > major {
		return true
	}
	if v.Major < major {
		return false
	}

	minor := 0
	if len(rest) > 0 {
		minor = rest[0]
	}

	if v.Minor > minor {
		return true
	}
	if v.Minor < minor {
		return false
	}

	patch := 0
	if len(rest) > 1 {
		patch = rest[1]
	}

	return v.Patch >= patch
}

// ParseVersion parses a version string (e.g., "1.2.3") without allocating memory for parts.
func ParseVersion(s string) Version {
	if s == "" {
		return Version{KnownVersion: false}
	}

	v := Version{Raw: s}

	part := 0
	current := 0
	hasDigits := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			current = current*10 + int(c-'0')
			hasDigits = true
		} else if c == '.' || c == '_' {
			// Apple writes versions with underscores ("17_5", "10_15_7"). The
			// separator is a spelling difference, not a semantic one, so both
			// forms parse to the same value and Raw keeps the original.
			if !hasDigits && part == 0 {
				return Version{KnownVersion: false, Raw: s}
			}
			switch part {
			case 0:
				v.Major = current
			case 1:
				v.Minor = current
			case 2:
				v.Patch = current
			}
			part++
			current = 0
			hasDigits = false
			if part > 2 {
				break
			}
		} else {
			// e.g., '-' or '_' or alphabet
			break
		}
	}

	if hasDigits {
		switch part {
		case 0:
			v.Major = current
		case 1:
			v.Minor = current
		case 2:
			v.Patch = current
		}
	}

	v.KnownVersion = part > 0 || hasDigits
	return v
}
