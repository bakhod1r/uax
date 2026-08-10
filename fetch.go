package uax

// FetchContext represents the parsed Sec-Fetch-* headers.
type FetchContext struct {
	Site string `json:"site,omitempty"`
	Mode string `json:"mode,omitempty"`
	Dest string `json:"dest,omitempty"`
	User string `json:"user,omitempty"`
}

// ParseFetchContext parses the Sec-Fetch headers.
func ParseFetchContext(site, mode, dest, user string) *FetchContext {
	if site == "" && mode == "" && dest == "" && user == "" {
		return nil
	}

	// Some normalization could be done here (e.g. trimming quotes),
	// but for MVP we assume headers are passed mostly raw or already cleaned.
	// Sec-Fetch headers are typically tokens (e.g. `navigate`, `document`, `?1`).
	return &FetchContext{
		Site: site,
		Mode: mode,
		Dest: dest,
		User: user,
	}
}

// IsNavigation returns true if this fetch represents a top-level document navigation.
func (f *FetchContext) IsNavigation() bool {
	return f.Mode == "navigate" && f.Dest == "document"
}
