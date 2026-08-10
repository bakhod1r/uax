package uax

// Config configures a Parser. The zero value is valid and is what the
// package-level functions use: no cache, no recorder, no shared state.
//
// Caching is opt-in on purpose. A cache is a global mutable side effect on what
// is otherwise a pure function, it makes benchmarks report the cache rather
// than the parser, and it retains attacker-supplied strings in memory. Callers
// that want it should say so, and size it themselves.
type Config struct {
	// Cache memoises results by raw User-Agent. Production traffic repeats a
	// few thousand distinct strings across millions of requests, so this is
	// usually worth enabling on a request path — but only when the caller has
	// decided the memory ceiling is acceptable.
	Cache CacheConfig

	// MaxInputLength bounds how much of an input is examined. Longer input is
	// truncated and flagged via Client.Truncated rather than rejected, because
	// a parser that errors on a hostile string is a parser that hands the
	// caller a decision it cannot make. Zero selects DefaultMaxInputLength.
	MaxInputLength int

	// Recorder collects unrecognised product tokens so a deployment can improve
	// its own rules. It is nil by default: recording is a deliberate choice,
	// never something that happens because a caller did not opt out.
	Recorder *Recorder

	// Metrics receives counters. Nil means no metrics are emitted.
	Metrics MetricsReporter

	// Registry layers custom rules over the built-in database. Nil uses the
	// built-in rules alone.
	Registry *Registry

	// Device names the hardware behind a model code. It is nil by default, and
	// then Device.Model reports the code the request carried while the name,
	// brand and form factor stay empty.
	//
	// This package holds no device catalogue. "SM-S928B" does not contain
	// "Galaxy S24 Ultra" — a manufacturer assigned that mapping, so answering
	// it needs published data, and published data moves on its own schedule.
	// github.com/bakhod1r/devicex is such a catalogue and satisfies this signature
	// directly:
	//
	//	p := uax.New(uax.Config{Device: devicex.Describe})
	//
	// The signature carries no types from either package, so neither has to
	// import the other. It is called with the extracted model code and the raw
	// User-Agent: the code answers Android, and the User-Agent answers the
	// hardware that publishes no code and names itself in prose instead, like
	// "Macintosh" or "PlayStation 5". Either argument may be empty.
	//
	// model is the identifier the request actually carried, which is not
	// always the code passed in: hardware that publishes no code names itself
	// in prose, and "iPhone" is both the name and the whole identifier. A hook
	// that cannot point at something in the request returns model empty, and
	// Device.Model then reports the extracted code, or nothing.
	//
	// deviceType is "Mobile", "Tablet", "Desktop", "Console" or empty. It must
	// return ok == false rather than guess: a wrong device name is worse than
	// no device name, because it is a fact-shaped value a caller will store.
	Device func(code, ua string) (name, brand, model, family, deviceType string, confidence float64, ok bool)

	// SkipDevice omits device resolution entirely, including the model code.
	// Useful on hot paths that only need browser and OS.
	SkipDevice bool

	// SkipBot omits bot, AI, automation and HTTP-client resolution.
	SkipBot bool
}

// CacheConfig controls result memoisation.
type CacheConfig struct {
	Enabled bool
	// Size is the maximum number of entries. Zero selects DefaultCacheSize.
	Size int
}

const (
	// DefaultMaxInputLength is a generous bound on a real User-Agent, which
	// rarely exceeds 512 bytes. It exists to cap the work a hostile client can
	// cause, not to reject unusual browsers.
	DefaultMaxInputLength = 4096

	// DefaultCacheSize applies when caching is enabled without a size.
	DefaultCacheSize = 4096
)

// Parser parses User-Agents under a fixed configuration.
//
// A Parser is safe for concurrent use. The rule tables it reads are immutable
// after init, so parsing takes no locks; only an enabled cache synchronises,
// and it shards to keep contention off the hot path.
type Parser struct {
	cfg      Config
	cache    *Cache
	maxInput int
}

// New returns a Parser using cfg.
func New(cfg Config) *Parser {
	p := &Parser{cfg: cfg, maxInput: cfg.MaxInputLength}
	if p.maxInput <= 0 {
		p.maxInput = DefaultMaxInputLength
	}
	if cfg.Cache.Enabled {
		size := cfg.Cache.Size
		if size <= 0 {
			size = DefaultCacheSize
		}
		p.cache = NewCache(size)
	}
	return p
}

// defaultParser backs the package-level functions. Its cache is disabled, so
// Parse is a pure function of its input with no shared state and no retention.
var defaultParser = New(Config{})

// Default returns the parser backing the package-level functions: no cache, no
// recorder, no metrics. It is safe for concurrent use and holds no state.
func Default() *Parser { return defaultParser }
