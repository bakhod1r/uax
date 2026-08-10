package uax

import (
	"sync"
	"sync/atomic"

	"github.com/bakhod1r/uax/internal/table"
)

// Rule is a caller-supplied classification rule.
//
// It exists so a deployment can teach uax about clients no public database
// will ever carry — an internal mobile app, a partner's crawler, a company
// HTTP client. Custom rules are merged with the built-in database and compete
// on Priority alone, so a rule can also override a built-in one by outranking
// it.
//
// Exactly one match form must be set:
//
//	Token    exact match against a product token or a comment field
//	Prefix   the token or field starts with this string
//	Contains substring anywhere in the input — correct, but the slowest form
type Rule struct {
	ID   string
	Name string

	Family   string
	Engine   string
	Brand    string
	Model    string
	Series   string
	Category string
	Provider string
	Type     string
	Year     int

	Token    string
	Prefix   string
	Contains string

	// VersionFrom selects a version-extraction strategy: "token",
	// "version_product", "rv", "apple", "nt", "cros" or "suffix". Empty means
	// the rule carries no version.
	VersionFrom string

	// Priority decides which rule wins when several match. Built-in rules use
	// 50–320; a custom rule at 1000 always wins.
	Priority int

	// Confidence must be in (0,1]. It is reported verbatim on the resolved
	// field, so it should reflect how much the match really proves.
	Confidence float32
}

func (r Rule) toTable() table.Rule {
	return table.Rule{
		ID: r.ID, Name: r.Name, Family: r.Family, Engine: r.Engine,
		Brand: r.Brand, Model: r.Model, Series: r.Series,
		Category: r.Category, Provider: r.Provider, Type: r.Type,
		Token: r.Token, Prefix: r.Prefix, Contains: r.Contains,
		VersionFrom: r.VersionFrom, Year: r.Year,
		Priority: r.Priority, Confidence: r.Confidence,
	}
}

// builtinGroups is the built-in database, in category order.
func builtinGroups() [numCategories][]table.Rule {
	var g [numCategories][]table.Rule
	g[catBrowserIdx] = table.Browsers
	g[catEngineIdx] = table.Engines
	g[catOSIdx] = table.OSes
	g[catBotIdx] = table.Bots
	g[catAppIdx] = table.Apps
	g[catAutomationIdx] = table.Automations
	g[catClientIdx] = table.Clients
	return g
}

// builtinIndexes is the compiled built-in database. It is created once at init
// and never mutated, so parsing against it takes no locks.
var builtinIndexes = buildIndexes(builtinGroups())

// Registry holds custom rules layered over the built-in database.
//
// Registering rebuilds the compiled indexes and swaps them in atomically, so a
// registry can be updated while requests are in flight: readers never block and
// never observe a half-built database.
type Registry struct {
	mu      sync.Mutex
	rules   [numCategories][]table.Rule
	current atomic.Pointer[indexSet]
}

// NewRegistry returns an empty registry that resolves exactly like the
// built-in database until rules are added.
func NewRegistry() *Registry {
	r := &Registry{}
	r.current.Store(builtinIndexes)
	return r
}

// RegisterBrowser adds browser rules.
func (r *Registry) RegisterBrowser(rules ...Rule) { r.register(catBrowserIdx, rules) }

// RegisterEngine adds layout engine rules.
func (r *Registry) RegisterEngine(rules ...Rule) { r.register(catEngineIdx, rules) }

// RegisterOS adds operating-system rules.
func (r *Registry) RegisterOS(rules ...Rule) { r.register(catOSIdx, rules) }

// RegisterBot adds bot rules.
func (r *Registry) RegisterBot(rules ...Rule) { r.register(catBotIdx, rules) }

// RegisterApp adds in-app WebView rules.
func (r *Registry) RegisterApp(rules ...Rule) { r.register(catAppIdx, rules) }

// RegisterAutomation adds automation and headless rules.
func (r *Registry) RegisterAutomation(rules ...Rule) { r.register(catAutomationIdx, rules) }

// RegisterHTTPClient adds HTTP client and CLI tool rules.
func (r *Registry) RegisterHTTPClient(rules ...Rule) { r.register(catClientIdx, rules) }

func (r *Registry) register(cat category, rules []Rule) {
	if len(rules) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, rule := range rules {
		r.rules[cat] = append(r.rules[cat], rule.toTable())
	}
	r.rebuild()
}

// Reset removes every custom rule, returning the registry to the built-in
// database.
func (r *Registry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rules = [numCategories][]table.Rule{}
	r.current.Store(builtinIndexes)
}

// rebuild compiles a fresh generation and publishes it. The caller holds mu.
func (r *Registry) rebuild() {
	groups := builtinGroups()
	for cat := category(0); cat < numCategories; cat++ {
		groups[cat] = mergeRules(groups[cat], r.rules[cat])
	}
	r.current.Store(buildIndexes(groups))
}

// mergeRules combines built-in and custom rules into one priority-sorted slice.
func mergeRules(builtin, custom []table.Rule) []table.Rule {
	if len(custom) == 0 {
		return builtin
	}

	all := make([]table.Rule, 0, len(builtin)+len(custom))
	all = append(all, custom...)
	all = append(all, builtin...)

	// buildIndexes relies on descending priority order for its early exits,
	// and custom rules arrive unsorted.
	sortByPriority(all)
	return all
}

// sortByPriority orders rules by descending priority. Insertion sort suits a
// nearly-sorted slice, and this runs at registration time, never on a parse.
func sortByPriority(rules []table.Rule) {
	for i := 1; i < len(rules); i++ {
		for j := i; j > 0 && rules[j].Priority > rules[j-1].Priority; j-- {
			rules[j], rules[j-1] = rules[j-1], rules[j]
		}
	}
}

// indexes returns the generation a parse should use.
func (p *Parser) indexes() *indexSet {
	if p.cfg.Registry != nil {
		// A Registry that has never had a rule registered holds nothing
		// compiled — &Registry{} is a natural thing for a caller to write, and
		// NewRegistry is a convenience rather than a requirement. Fall back to
		// the built-in database, which is what an empty registry means: no
		// custom rules, not no rules at all.
		if ix := p.cfg.Registry.current.Load(); ix != nil {
			return ix
		}
	}
	return builtinIndexes
}
