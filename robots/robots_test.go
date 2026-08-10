package robots

import (
	"strings"
	"testing"

	"github.com/bakhod1r/uax"
	"github.com/bakhod1r/uax/internal/table"
)

func TestRobotsDeniesTheRequestedCategories(t *testing.T) {
	t.Parallel()

	out := Render(Policy{Deny: []uax.BotCategory{uax.BotAICrawler, uax.BotAIAgent}})

	for _, want := range []string{"GPTBot", "ClaudeBot", "CCBot", "Google-Extended", "ChatGPT-User"} {
		if !strings.Contains(out, "User-agent: "+want+"\n") {
			t.Errorf("%s was not denied\n%s", want, out)
		}
	}

	// A search engine is not an AI crawler, and blocking one by accident would
	// cost the operator their traffic.
	if strings.Contains(out, "Googlebot\n") {
		t.Errorf("Googlebot was denied although only AI categories were requested\n%s", out)
	}
}

func TestRobotsIsWellFormed(t *testing.T) {
	t.Parallel()

	out := Render(Policy{
		Deny:    []uax.BotCategory{uax.BotAICrawler},
		Sitemap: "https://example.com/sitemap.xml",
	})

	// Every User-agent line must be followed by a rule; a bare group is
	// ignored by crawlers, which would silently mean "no policy at all".
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "User-agent:") {
			if i+1 >= len(lines) || !strings.HasPrefix(lines[i+1], "Disallow:") {
				t.Errorf("line %d: %q is not followed by a Disallow rule", i+1, l)
			}
		}
	}

	if !strings.Contains(out, "Sitemap: https://example.com/sitemap.xml") {
		t.Error("sitemap was not emitted")
	}
	if !strings.HasSuffix(out, "\n") {
		t.Error("output does not end in a newline")
	}
}

func TestRobotsNeverEmitsTheGenericFallbacks(t *testing.T) {
	t.Parallel()

	// "bot", "crawler" and "spider" are substring heuristics for recognising an
	// unknown crawler. As robots.txt tokens they would match a vast set of
	// clients the operator never intended to name.
	out := Render(Policy{Deny: []uax.BotCategory{uax.BotGeneric, uax.BotAICrawler}})
	for _, bad := range []string{"User-agent: bot", "User-agent: crawler", "User-agent: spider"} {
		if strings.Contains(out, bad+"\n") {
			t.Errorf("generic heuristic %q leaked into robots.txt\n%s", bad, out)
		}
	}
}

func TestRobotsAllowListInvertsThePolicy(t *testing.T) {
	t.Parallel()

	out := Render(Policy{Allow: []uax.BotCategory{uax.BotSearchEngine}})

	if !strings.Contains(out, "User-agent: Googlebot\nAllow: /") {
		t.Errorf("Googlebot was not allowed\n%s", out)
	}
	// An allow list is only meaningful with a default deny after it.
	if !strings.HasSuffix(strings.TrimSpace(out), "User-agent: *\nDisallow: /") {
		t.Errorf("allow list did not end in a default deny\n%s", out)
	}
}

func TestRobotsByProvider(t *testing.T) {
	t.Parallel()

	out := Render(Policy{DenyProviders: []string{"OpenAI"}})
	for _, want := range []string{"GPTBot", "OAI-SearchBot", "ChatGPT-User"} {
		if !strings.Contains(out, "User-agent: "+want+"\n") {
			t.Errorf("%s was not denied for provider OpenAI\n%s", want, out)
		}
	}
	if strings.Contains(out, "ClaudeBot") {
		t.Errorf("a different provider's bot was denied\n%s", out)
	}
}

// Every token the generator emits must be one the parser recognises. If they
// ever diverge, the file would block a name no client actually sends.
func TestRobotsTokensRoundTripThroughTheParser(t *testing.T) {
	t.Parallel()

	for _, b := range Agents(Policy{Deny: allBotCategories()}) {
		c := uax.Parse(b.Token)
		if c.Bot == nil {
			t.Errorf("token %q (%s) is not recognised by the parser", b.Token, b.Name)
			continue
		}
		if c.Bot.Name != b.Name {
			t.Errorf("token %q parsed as %q, want %q", b.Token, c.Bot.Name, b.Name)
		}
	}
}

func allBotCategories() []uax.BotCategory {
	return []uax.BotCategory{
		uax.BotSearchEngine, uax.BotAICrawler, uax.BotAIAgent, uax.BotSocialCrawler,
		uax.BotSEO, uax.BotMonitoring, uax.BotSecurity,
	}
}

func TestRobotsSkipsBotsWithoutTokens(t *testing.T) {
	// Not t.Parallel() because it modifies global state
	oldBots := table.Bots
	defer func() { table.Bots = oldBots }()

	table.Bots = append([]table.Rule(nil), oldBots...) // clone to avoid race if others run
	table.Bots = append(table.Bots, table.Rule{
		Name:     "NoTokenBot",
		Category: string(uax.BotSearchEngine),
	})

	out := Agents(Policy{Allow: []uax.BotCategory{uax.BotSearchEngine}})
	for _, a := range out {
		if a.Name == "NoTokenBot" {
			t.Errorf("expected bot without token to be skipped")
		}
	}
}
