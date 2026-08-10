package uax

// PrivacyRisk is how identifying a client's self-description is.
//
// It is a coarse ordering, not a measurement. Real fingerprinting entropy
// depends on the population a visitor is drawn from — a Samsung phone is
// unremarkable in a consumer audience and conspicuous in a fleet of managed
// laptops — and this package sees one request, never the population. Treat it
// as "how much did this client volunteer about itself", which is a question a
// single request can answer.
type PrivacyRisk string

const (
	PrivacyLow    PrivacyRisk = "low"
	PrivacyMedium PrivacyRisk = "medium"
	PrivacyHigh   PrivacyRisk = "high"
)

// Privacy summarises how identifying a parsed client is.
type Privacy struct {
	Risk PrivacyRisk `json:"risk"`

	// Identifiers lists the fields that narrow this client down, so a caller
	// can see exactly what drove the rating rather than trusting a number.
	Identifiers []string `json:"identifiers,omitempty"`
}

// Privacy reports how identifying this client's description is.
//
// The rating counts the specific, narrowing facts the client disclosed. An
// exact device model is the strongest of them: it can separate one visitor from
// nearly everyone else, which is precisely why Chromium stopped sending it by
// default and why it now requires an explicit Accept-CH request.
func (c *Client) Privacy() Privacy {
	var ids []string

	if c.Device.Model != "" {
		ids = append(ids, "device.model")
	}
	if c.Device.Name != "" {
		ids = append(ids, "device.name")
	}
	// A full four-part browser version is far more identifying than a major
	// version, which is why Client Hints send only the significant version
	// unless the full list is explicitly requested.
	if c.Browser.Version.Patch > 0 {
		ids = append(ids, "browser.version.patch")
	}
	if c.OS.Version.Minor > 0 || c.OS.Version.Patch > 0 {
		ids = append(ids, "os.version.minor")
	}
	if c.OS.Arch != "" {
		ids = append(ids, "os.arch")
	}
	if c.App != nil {
		ids = append(ids, "app")
	}
	if c.Hints != nil && c.Hints.Model != "" {
		ids = append(ids, "hints.model")
	}

	risk := PrivacyLow
	switch {
	case len(ids) >= 4:
		risk = PrivacyHigh
	case len(ids) >= 2:
		risk = PrivacyMedium
	}

	// An exact device model weighs more than the other signals: it is the single
	// most narrowing thing a client can disclose, which is why Chromium stopped
	// sending it by default. Counting it as one signal among equals would
	// understate it.
	if c.Device.Model != "" && len(ids) >= 3 {
		risk = PrivacyHigh
	}

	return Privacy{Risk: risk, Identifiers: ids}
}

// Coarse returns a copy of the client with the narrowing details removed:
// exact device model and name, version components below the major, and the
// architecture.
//
// Use it when the client description is going somewhere it will be retained —
// analytics, logs, a data warehouse. What survives is enough to answer "which
// browsers and platforms do my users have", which is what those systems are
// usually for, without carrying a per-visitor identifier along with it.
func (c *Client) Coarse() *Client {
	if c == nil {
		return nil
	}

	out := *c
	out.hintFindings = nil

	out.Device.Model = ""
	out.Device.Name = ""
	out.Device.Series = ""
	out.Device.Year = 0
	out.OS.Arch = ""
	out.Browser.Version = majorOnly(c.Browser.Version)
	out.OS.Version = majorOnly(c.OS.Version)
	out.Raw = ""

	// The pointers must be re-aimed at the copy's own storage, or the coarse
	// client would still reference the original's fields.
	out.Bot, out.AI, out.App, out.Automation, out.HTTPClient = nil, nil, nil, nil, nil
	if c.Bot != nil {
		out.botStore = *c.Bot
		out.botStore.Version = majorOnly(c.Bot.Version)
		out.Bot = &out.botStore
	}
	if c.AI != nil {
		out.aiStore = *c.AI
		out.AI = &out.aiStore
	}
	if c.App != nil {
		out.appStore = *c.App
		out.appStore.Version = majorOnly(c.App.Version)
		out.App = &out.appStore
	}
	if c.Automation != nil {
		out.automationStore = *c.Automation
		out.automationStore.Version = majorOnly(c.Automation.Version)
		out.Automation = &out.automationStore
	}
	if c.HTTPClient != nil {
		out.httpClientStore = *c.HTTPClient
		out.httpClientStore.Version = majorOnly(c.HTTPClient.Version)
		out.HTTPClient = &out.httpClientStore
	}

	if c.Hints != nil {
		h := *c.Hints
		h.Model = ""
		h.PlatformVersion = ""
		h.Arch = ""
		h.Bitness = ""
		h.FullVersionList = nil
		out.Hints = &h
	}

	return &out
}

func majorOnly(v Version) Version {
	if !v.Known() {
		return Version{}
	}
	return Version{Major: v.Major, KnownVersion: true, Raw: v.Raw}
}
