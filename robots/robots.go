// Package robots renders a robots.txt from what uax knows about crawlers.
//
// It is separate from the parser because it is the opposite operation. Parsing
// reads a request and reports who sent it; this writes a file that tells
// crawlers what they may do before any request happens. A server that only
// classifies traffic never needs this code, and a build that does not import
// it does not carry the rendering.
package robots

import (
	"sort"
	"strings"

	"github.com/bakhod1r/uax"
	"github.com/bakhod1r/uax/internal/table"
)

// robots.txt generation.
//
// The bot database already knows every crawler's product token, its category
// and its operator. A robots.txt file is that same knowledge written in a
// different syntax, so generating one needs no new data — only the observation
// that the token a parser matches on is exactly the token a crawler matches
// itself against in robots.txt.
//
// This exists because the alternative is what people actually do today: paste a
// list of AI crawler names from a blog post, and never update it.

// Policy describes which crawlers a site wants to admit.
//
// Deny and Allow are opposite styles and should not be combined. Deny names the
// unwanted crawlers and leaves everything else alone, which is the safe default:
// a crawler nobody thought about keeps working. Allow names the wanted ones and
// refuses the rest, which is stricter and will block crawlers that do not exist
// yet — including ones the site would have wanted.
type Policy struct {
	// Deny blocks every crawler in these categories from the whole site.
	Deny []uax.BotCategory

	// DenyProviders blocks every crawler operated by these companies, matched
	// against Bot.Provider ("OpenAI", "Anthropic", "Google") without regard to
	// case. Useful when the objection is to an operator rather than a purpose.
	DenyProviders []string

	// Allow permits these categories and denies everything else through a
	// trailing "User-agent: *" group.
	Allow []uax.BotCategory

	// Path restricts the rules to a prefix instead of the whole site.
	// Defaults to "/".
	Path string

	// Sitemap is emitted once at the end when set.
	Sitemap string
}

// Agent is one crawler as robots.txt sees it.
type Agent struct {
	// Token is the string to put after "User-agent:". It is the same token the
	// parser matches, which is what keeps the two from drifting apart.
	Token string

	Name     string
	Category uax.BotCategory
	Provider string
}

// Agents returns the crawlers a policy selects, sorted by category and
// then name so the generated file is stable across runs.
//
// It is exported separately from Robots so callers can build their own output —
// a Cloudflare rule set, an nginx map, an allowlist in a WAF — from the same
// selection.
func Agents(p Policy) []Agent {
	cats := p.Deny
	if len(cats) == 0 && len(p.DenyProviders) == 0 {
		cats = p.Allow
	}

	wantCat := make(map[uax.BotCategory]bool, len(cats))
	for _, c := range cats {
		wantCat[c] = true
	}
	wantProvider := make(map[string]bool, len(p.DenyProviders))
	for _, s := range p.DenyProviders {
		wantProvider[strings.ToLower(s)] = true
	}

	var out []Agent
	for i := range table.Bots {
		r := &table.Bots[i]

		// The generic rules match on the substrings "bot", "crawler" and
		// "spider". They are a useful last resort when classifying a request,
		// but as robots.txt tokens they would name an enormous set of clients
		// the operator never chose, so they are never emitted.
		if uax.BotCategory(r.Category) == uax.BotGeneric {
			continue
		}

		if !wantCat[uax.BotCategory(r.Category)] && !wantProvider[strings.ToLower(r.Provider)] {
			continue
		}

		token := token(r)
		if token == "" {
			continue
		}
		out = append(out, Agent{
			Token:    token,
			Name:     r.Name,
			Category: uax.BotCategory(r.Category),
			Provider: r.Provider,
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		return out[i].Token < out[j].Token
	})
	return out
}

// robotsToken is the literal a crawler will recognise as naming itself.
// A rule matches by exact token, by prefix or by substring, and in all three
// cases the stored literal is what appears in the crawler's User-Agent.
func token(r *table.Rule) string {
	switch {
	case r.Token != "":
		return r.Token
	case r.Prefix != "":
		return r.Prefix
	default:
		return r.Contains
	}
}

// Render renders a robots.txt body for the policy.
//
//	robots.Render(robots.Policy{Deny: []uax.BotCategory{uax.BotAICrawler, uax.BotAIAgent}})
//
// The result is advisory, as robots.txt always is: a crawler that ignores it is
// not stopped by it. Blocking by request requires the parser at the edge, which
// is what the uax package is for.
func Render(p Policy) string {
	path := p.Path
	if path == "" {
		path = "/"
	}

	agents := Agents(p)
	allowMode := len(p.Deny) == 0 && len(p.DenyProviders) == 0 && len(p.Allow) > 0

	rule := "Disallow: "
	if allowMode {
		rule = "Allow: "
	}

	var b strings.Builder
	b.WriteString("# Generated by uax from its bot database. Advisory only: robots.txt is\n")
	b.WriteString("# obeyed voluntarily, and a crawler that ignores it is unaffected.\n")

	var lastCat uax.BotCategory
	for _, a := range agents {
		if a.Category != lastCat {
			b.WriteString("\n# ")
			b.WriteString(string(a.Category))
			b.WriteByte('\n')
			lastCat = a.Category
		}
		b.WriteString("User-agent: ")
		b.WriteString(a.Token)
		b.WriteByte('\n')
		b.WriteString(rule)
		b.WriteString(path)
		b.WriteByte('\n')
	}

	if allowMode {
		// Without this the allow list means nothing: naming who may crawl does
		// not exclude anyone unless the default is a refusal.
		b.WriteString("\n# everyone else\nUser-agent: *\nDisallow: ")
		b.WriteString(path)
		b.WriteByte('\n')
	}

	if p.Sitemap != "" {
		b.WriteString("\nSitemap: ")
		b.WriteString(p.Sitemap)
		b.WriteByte('\n')
	}
	return b.String()
}
