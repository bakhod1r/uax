package uax

// DeviceType is the physical form factor of the client.
type DeviceType string

const (
	DeviceUnknown  DeviceType = ""
	DeviceDesktop  DeviceType = "Desktop"
	DeviceMobile   DeviceType = "Mobile"
	DeviceTablet   DeviceType = "Tablet"
	DeviceTV       DeviceType = "TV"
	DeviceConsole  DeviceType = "Console"
	DeviceWatch    DeviceType = "Watch"
	DeviceCar      DeviceType = "Car"
	DeviceXR       DeviceType = "XR"   // headsets and glasses that replace or augment the view
	DeviceEInk     DeviceType = "EInk" // slow screen refresh, little or no colour
	DeviceEmbedded DeviceType = "Embedded"
	DeviceBot      DeviceType = "Bot"
)

// BotCategory says what a bot is for. It carries far more decision value than
// a bare "is a bot" boolean: a search crawler, an AI trainer and a scraper all
// warrant different handling.
type BotCategory string

const (
	BotSearchEngine  BotCategory = "SearchEngine"
	BotAICrawler     BotCategory = "AICrawler"
	BotAIAgent       BotCategory = "AIAgent"
	BotSocialCrawler BotCategory = "SocialCrawler"
	BotSEO           BotCategory = "SEO"
	BotMonitoring    BotCategory = "Monitoring"
	BotSecurity      BotCategory = "Security"
	BotGeneric       BotCategory = "Generic"
)

// Kind is the top-level answer to "what is talking to me".
type Kind string

const (
	KindUnknown    Kind = "Unknown"
	KindHuman      Kind = "Human"
	KindBot        Kind = "Bot"
	KindAICrawler  Kind = "AICrawler"
	KindAIAgent    Kind = "AIAgent"
	KindTool       Kind = "Tool"
	KindAutomation Kind = "Automation"
)

// Source records where a resolved field came from. Client Hints are
// authoritative; a value inferred from a sibling field is not.
type Source string

const (
	SourceNone     Source = ""
	SourceUA       Source = "ua"
	SourceHints    Source = "hints"
	SourceBoth     Source = "both"
	SourceInferred Source = "inferred"
)

// Browser is the resolved browser identity.
//
// Family groups derivatives under their base (Edge, Brave and Vivaldi all have
// the family Chrome); Name is the specific product. Mode distinguishes a
// headless or WebView instance of an otherwise ordinary browser.
type Browser struct {
	Name       string  `json:"name,omitempty"`
	Family     string  `json:"family,omitempty"`
	Version    Version `json:"version,omitempty"`
	Mode       string  `json:"mode,omitempty"` // Headless, WebView, "" for a normal browser
	Confidence float32 `json:"confidence,omitempty"`
	Source     Source  `json:"source,omitempty"`
}

// Engine is the layout engine.
type Engine struct {
	Name       string  `json:"name,omitempty"`
	Version    Version `json:"version,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

// OS is the resolved operating system.
type OS struct {
	Name    string  `json:"name,omitempty"`
	Family  string  `json:"family,omitempty"`
	Version Version `json:"version,omitempty"`
	Arch    string  `json:"arch,omitempty"` // amd64, arm64, x86, arm — only when stated

	// Bitness is the platform's CPU width, "32" or "64", exactly as
	// Sec-CH-UA-Bitness stated it. It describes the machine, not the process:
	// a 32-bit browser on a 64-bit Windows reports bitness 64 and sets
	// Sec-CH-UA-WoW64 instead. Empty unless the client sent the hint.
	Bitness string `json:"bitness,omitempty"`

	Confidence float32 `json:"confidence,omitempty"`
	Source     Source  `json:"source,omitempty"`

	// Reduced reports that the platform version in the User-Agent is frozen by
	// the browser rather than observed from the system. Chrome's User-Agent
	// Reduction did this deliberately: every reduced Android build reports 10,
	// macOS caps at 10.15.7, and Windows has reported NT 10.0 since Windows 10.
	//
	// When Reduced is set, Version is either absent (Android, where the value
	// is a placeholder that contradicts reality) or a floor (macOS and Windows,
	// where the real version is the stated one or higher). Sec-CH-UA-Platform-
	// Version carries the true value and clears this flag.
	Reduced bool `json:"reduced,omitempty"`
}

// Device is the resolved hardware.
//
// Model and Name are deliberately separate and must not be conflated:
//
//	Model is the verbatim code from the User-Agent ("SM-S928B"). It is evidence
//	and never changes.
//
//	Name is the marketing name ("Galaxy S24 Ultra"). It comes from a lookup
//	table, because no rule can derive one from the other. When the table has no
//	entry, Name stays empty — it is never approximated, since a wrong device
//	name is worse than none at all.
type Device struct {
	Type       DeviceType `json:"type,omitempty"`
	Brand      string     `json:"brand,omitempty"`
	Model      string     `json:"model,omitempty"`
	Name       string     `json:"name,omitempty"`
	Family     string     `json:"family,omitempty"`
	Series     string     `json:"series,omitempty"`
	Year       int        `json:"year,omitempty"`
	Confidence float32    `json:"confidence,omitempty"`
	Source     Source     `json:"source,omitempty"`

	// Reduced reports that the browser replaced the model code with a
	// placeholder ("K"). The handset exists and has a code; this client simply
	// declined to state it, which is different from a code nobody recognised.
	// Sec-CH-UA-Model carries the real one and clears this flag.
	Reduced bool `json:"reduced,omitempty"`
}

// Bot is a self-identified crawler, agent or scraper.
//
// Every field here is a claim made by the client. A User-Agent saying
// "Googlebot" is not proof; the uax/verify module exists to turn that claim
// into a verified fact using DNS and published address ranges.
type Bot struct {
	Name       string      `json:"name,omitempty"`
	Category   BotCategory `json:"category,omitempty"`
	Provider   string      `json:"provider,omitempty"`
	Version    Version     `json:"version,omitempty"`
	Confidence float32     `json:"confidence,omitempty"`
}

// IsSearchEngine reports whether the bot indexes for a search product.
func (b *Bot) IsSearchEngine() bool { return b != nil && b.Category == BotSearchEngine }

// IsAICrawler reports whether the bot collects data for a model, either for
// training or for live retrieval.
func (b *Bot) IsAICrawler() bool {
	return b != nil && (b.Category == BotAICrawler || b.Category == BotAIAgent)
}

// IsSocialCrawler reports whether the bot is building a link preview.
func (b *Bot) IsSocialCrawler() bool { return b != nil && b.Category == BotSocialCrawler }

// IsSEO reports whether the bot belongs to commercial SEO tooling.
func (b *Bot) IsSEO() bool { return b != nil && b.Category == BotSEO }

// AI describes an AI system fetching the resource. It is derived from the bot
// rule that matched rather than stored separately, so the two can never drift.
type AI struct {
	Name     string      `json:"name,omitempty"`
	Provider string      `json:"provider,omitempty"`
	Category BotCategory `json:"category,omitempty"` // AICrawler or AIAgent
}

// App is the native application hosting an embedded WebView.
type App struct {
	Name       string  `json:"name,omitempty"`
	Category   string  `json:"category,omitempty"`
	Version    Version `json:"version,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

// Automation is a headless browser or a browser-automation framework.
//
// Only self-declared markers are reported. A driver configured to hide is
// invisible here by design, and the absence of an Automation value is
// therefore never evidence that a client is human.
type Automation struct {
	Name       string  `json:"name,omitempty"`
	Category   string  `json:"category,omitempty"` // Headless, Framework, Embedded
	Version    Version `json:"version,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

// HTTPClient is a non-browser client: a CLI tool, or an HTTP library embedded
// in an application.
type HTTPClient struct {
	Name       string  `json:"name,omitempty"`
	Category   string  `json:"category,omitempty"` // CLI, Library, Package
	Version    Version `json:"version,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}
