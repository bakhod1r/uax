package uax

import (
	"strings"
	"unsafe"

	"github.com/bakhod1r/uax/internal/table"
)

// unsafeBytes and unsafeString convert between the two representations without
// copying. Both are safe here because neither the scanner nor the resolver ever
// writes through the result: every value the parser returns is either a
// constant from the rule tables or a read-only view into the caller's
// input string, which is itself immutable.
func unsafeBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func unsafeString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// ---------------------------------------------------------------- versions

// versionFor applies the rule's declared extraction strategy. Strategies exist
// because there is no single convention: Safari's real version lives in a
// separate product, Windows encodes a kernel number, and Apple uses
// underscores.
func versionFor(t *tokens, r *table.Rule, matched string) Version {
	switch r.VersionFrom {
	case "token":
		return versionAfterToken(t, matched, r)
	case "version_product":
		return productVersion(t, "Version")
	case "rv":
		return rvVersion(t)
	case "apple":
		return appleVersion(t)
	case "nt":
		return ntVersion(t)
	case "cros":
		return crosVersion(t)
	case "suffix":
		return suffixVersion(t, matched, r)
	default:
		return Version{}
	}
}

// versionAfterToken takes the "Name/1.2.3" version attached to the matched
// product. When the match came from a comment field instead, the version trails
// the name inside that field.
func versionAfterToken(t *tokens, matched string, r *table.Rule) Version {
	for i := 0; i < t.nprod; i++ {
		if t.prodName(i) == matched && t.prodVersion(i) != "" {
			return ParseVersion(t.prodVersion(i))
		}
	}
	if v := trailingVersion(matched, r.Token); v.Known() {
		return v
	}
	// Some clients separate the name and version with a space rather than a
	// slash ("Instagram 300.0.0.0"), which the scanner emits as two products.
	for i := 0; i < t.nprod-1; i++ {
		if t.prodName(i) == matched && startsWithDigit(t.prodName(i+1)) {
			return ParseVersion(t.prodName(i + 1))
		}
	}
	return Version{}
}

// productVersion finds a standalone product by name and returns its version.
func productVersion(t *tokens, name string) Version {
	for i := 0; i < t.nprod; i++ {
		if t.prodName(i) == name {
			return ParseVersion(t.prodVersion(i))
		}
	}
	return Version{}
}

// rvVersion reads the "rv:11.0" form used by Trident and Gecko.
func rvVersion(t *tokens) Version {
	for i := 0; i < t.nfield; i++ {
		if f := t.fieldAt(i); strings.HasPrefix(f, "rv:") {
			return ParseVersion(f[len("rv:"):])
		}
	}
	return Version{}
}

// appleVersion reads Apple's underscore form. It looks for the "OS " marker
// because the digits follow it in every shape Apple ships: "CPU iPhone OS
// 17_5", the iPad's "CPU OS 17_5", and the desktop "Intel Mac OS X 10_15_7".
func appleVersion(t *tokens) Version {
	for i := 0; i < t.nfield; i++ {
		f := t.fieldAt(i)
		for {
			j := strings.Index(f, "OS ")
			if j < 0 {
				break
			}
			rest := f[j+3:]
			// "Mac OS X 10_15_7" puts an X between the marker and the digits.
			rest = strings.TrimPrefix(rest, "X ")
			if startsWithDigit(rest) {
				return ParseVersion(numericPrefix(rest))
			}
			f = f[j+3:]
		}
	}
	return Version{}
}

// ntVersion maps the Windows NT kernel version to its marketing release.
//
// Windows 11 also reports NT 10.0, so it is indistinguishable from Windows 10
// by User-Agent alone. Client Hints carry the real platform version, and
// merging them is what resolves the ambiguity.
func ntVersion(t *tokens) Version {
	f, ok := hasPrefixAny(t, "Windows NT")
	if !ok {
		return Version{}
	}

	kernel := ParseVersion(numericPrefix(strings.TrimSpace(f[len("Windows NT"):])))
	if !kernel.Known() {
		return Version{}
	}

	major, minor := 0, 0
	switch {
	case kernel.Major == 10:
		major = 10
	case kernel.Major == 6 && kernel.Minor == 3:
		major, minor = 8, 1
	case kernel.Major == 6 && kernel.Minor == 2:
		major = 8
	case kernel.Major == 6 && kernel.Minor == 1:
		major = 7
	case kernel.Major == 6 && kernel.Minor == 0:
		major = 6 // Vista
	case kernel.Major == 5:
		major, minor = 5, 1 // XP
	default:
		return Version{}
	}

	return Version{Major: major, Minor: minor, KnownVersion: true, Raw: kernel.Raw}
}

// crosVersion reads the ChromeOS platform build, which trails the architecture
// in "CrOS x86_64 14541.0.0".
func crosVersion(t *tokens) Version {
	f, ok := hasPrefixAny(t, "CrOS")
	if !ok {
		return Version{}
	}
	if i := strings.LastIndexByte(f, ' '); i > 0 {
		return ParseVersion(f[i+1:])
	}
	return Version{}
}

// suffixVersion reads the version trailing the match inside the same field,
// as in "Android 15" or "Tizen 6.0".
func suffixVersion(t *tokens, matched string, r *table.Rule) Version {
	needle := r.Prefix
	if needle == "" {
		needle = r.Token
	}
	if v := trailingVersion(matched, needle); v.Known() {
		return v
	}
	for i := 0; i < t.nprod-1; i++ {
		if t.prodName(i) == matched && startsWithDigit(t.prodName(i+1)) {
			return ParseVersion(t.prodName(i + 1))
		}
	}
	return Version{}
}

// trailingVersion extracts the numeric run that follows needle within s.
func trailingVersion(s, needle string) Version {
	if needle == "" || !strings.HasPrefix(s, needle) {
		return Version{}
	}
	rest := strings.TrimLeft(s[len(needle):], " /")
	if !startsWithDigit(rest) {
		return Version{}
	}
	return ParseVersion(numericPrefix(rest))
}

// numericPrefix returns the leading run of digits and version separators.
// Separators include '_' so Apple's form is handled without rewriting it.
func numericPrefix(s string) string {
	i := 0
	for i < len(s) {
		c := s[i]
		if (c >= '0' && c <= '9') || c == '.' || c == '_' {
			i++
			continue
		}
		break
	}
	return s[:i]
}

func startsWithDigit(s string) bool {
	return len(s) > 0 && s[0] >= '0' && s[0] <= '9'
}

// ------------------------------------------------------------------ fields

func resolveBrowser(t *tokens, m match) Browser {
	r, matched := m.rule, m.str
	if r == nil {
		return Browser{}
	}
	return Browser{
		Name:       r.Name,
		Family:     r.Family,
		Version:    versionFor(t, r, matched),
		Confidence: r.Confidence,
		Source:     SourceUA,
	}
}

// resolveEngine prefers the engine declared by the browser rule that matched.
// "AppleWebKit" alone cannot separate Blink from WebKit, so the browser is the
// better authority whenever there is one.
func resolveEngine(t *tokens, m, browser match) Engine {
	r, matched := m.rule, m.str

	var e Engine
	if r != nil {
		e = Engine{Name: r.Name, Version: versionFor(t, r, matched), Confidence: r.Confidence}
	}

	if br := browser.rule; br != nil && br.Engine != "" {
		e.Name = br.Engine
		if e.Confidence < br.Confidence {
			e.Confidence = br.Confidence
		}
	}
	return e
}

func resolveOS(t *tokens, m match) OS {
	r, matched := m.rule, m.str
	if r == nil {
		return OS{}
	}
	return OS{
		Name:       r.Name,
		Family:     r.Family,
		Version:    versionFor(t, r, matched),
		Arch:       resolveArch(t),
		Confidence: r.Confidence,
		Source:     SourceUA,
	}
}

// resolveArch reports the CPU architecture only when the client states it.
// Absence is reported as absence, never guessed from the platform.
func resolveArch(t *tokens) string {
	for i := 0; i < t.nfield; i++ {
		switch f := t.fieldAt(i); {
		case f == "x64", f == "x86_64", f == "amd64", f == "Win64":
			return "amd64"
		case f == "arm64", f == "aarch64":
			return "arm64"
		case f == "WOW64", f == "i686", f == "i386", f == "x86":
			return "x86"
		case f == "ARM", strings.HasPrefix(f, "armv"):
			return "arm"
		case strings.HasPrefix(f, "Linux x86_64"), strings.HasPrefix(f, "CrOS x86_64"):
			return "amd64"
		case strings.HasPrefix(f, "Linux aarch64"), strings.HasPrefix(f, "CrOS aarch64"):
			return "arm64"
		case strings.HasPrefix(f, "Linux armv"):
			return "arm"
		}
	}
	return ""
}

func resolveBot(t *tokens, m match, dst *Bot) *Bot {
	r, matched := m.rule, m.str
	if r == nil {
		return nil
	}
	// Written into storage the caller already owns, so a bot costs no
	// allocation while the API keeps a pointer — nil means "no bot", which a
	// zero struct could never express.
	*dst = Bot{
		Name:       r.Name,
		Category:   BotCategory(r.Category),
		Provider:   r.Provider,
		Version:    versionFor(t, r, matched),
		Confidence: r.Confidence,
	}
	return dst
}

func resolveApp(t *tokens, m match, dst *App) *App {
	r, matched := m.rule, m.str
	if r == nil {
		return nil
	}
	*dst = App{
		Name:       r.Name,
		Category:   r.Category,
		Version:    versionFor(t, r, matched),
		Confidence: r.Confidence,
	}
	return dst
}

func resolveAutomation(t *tokens, m match, dst *Automation) *Automation {
	r, matched := m.rule, m.str
	if r == nil {
		return nil
	}
	*dst = Automation{
		Name:       r.Name,
		Category:   r.Category,
		Version:    versionFor(t, r, matched),
		Confidence: r.Confidence,
	}
	return dst
}

func resolveHTTPClient(t *tokens, m match, dst *HTTPClient) *HTTPClient {
	r, matched := m.rule, m.str
	if r == nil {
		return nil
	}
	*dst = HTTPClient{
		Name:       r.Name,
		Category:   r.Category,
		Version:    versionFor(t, r, matched),
		Confidence: r.Confidence,
	}
	return dst
}

// ------------------------------------------------------------ derived facts

// deviceType decides the form factor. Order matters: explicit evidence from a
// device or OS rule outranks anything inferred from the presence of a token.
func deviceType(t *tokens, c *Client, osRule *table.Rule) DeviceType {
	if c.Bot != nil {
		return DeviceBot
	}
	if c.Device.Type != "" {
		return c.Device.Type
	}

	if osRule != nil && osRule.Type != "" {
		return DeviceType(osRule.Type)
	}

	switch {
	case t.hasTVToken:
		return DeviceTV
	case t.hasTabletToken:
		return DeviceTablet
	case t.hasMobileToken:
		return DeviceMobile
	}

	// An Android client that does not say "Mobile" is a tablet by convention:
	// Chrome omits the token specifically to signal the larger form factor.
	if c.OS.Family == "Android" {
		return DeviceTablet
	}

	switch c.OS.Name {
	case "Windows", "macOS", "Linux", "Ubuntu", "Debian", "Fedora", "Arch Linux",
		"Chrome OS", "FreeBSD", "OpenBSD":
		return DeviceDesktop
	}
	return DeviceUnknown
}

// kind is the top-level classification a caller usually branches on.
func kind(c *Client) Kind {
	switch {
	case c.Bot != nil:
		switch c.Bot.Category {
		case BotAICrawler:
			return KindAICrawler
		case BotAIAgent:
			return KindAIAgent
		default:
			return KindBot
		}
	case c.Automation != nil && c.Automation.Category != "Embedded":
		return KindAutomation
	case c.HTTPClient != nil:
		return KindTool
	case c.Browser.Name != "" || c.OS.Name != "":
		return KindHuman
	default:
		return KindUnknown
	}
}

// aiFor derives the AI view from the bot rule that matched, so the two can
// never disagree.
func aiFor(b *Bot, dst *AI) *AI {
	if b == nil || (b.Category != BotAICrawler && b.Category != BotAIAgent) {
		return nil
	}
	*dst = AI{Name: b.Name, Provider: b.Provider, Category: b.Category}
	return dst
}
