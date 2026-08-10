package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/bakhod1r/uax"
	"github.com/bakhod1r/uax/robots"
)

// robotsCategories maps the names a person would type to the categories the
// database uses. The short forms exist because "ai" is what the request
// actually is; nobody blocking model training thinks of it as "AICrawler".
var robotsCategories = map[string][]uax.BotCategory{
	"ai":         {uax.BotAICrawler, uax.BotAIAgent},
	"ai-crawler": {uax.BotAICrawler},
	"ai-agent":   {uax.BotAIAgent},
	"search":     {uax.BotSearchEngine},
	"social":     {uax.BotSocialCrawler},
	"seo":        {uax.BotSEO},
	"monitoring": {uax.BotMonitoring},
	"security":   {uax.BotSecurity},
}

func robotsCmd(args []string, out io.Writer) error {
	var p robots.Policy

	for i := 0; i < len(args); i++ {
		flag := args[i]
		value := func() (string, error) {
			if i+1 >= len(args) {
				return "", fmt.Errorf("%s needs a value", flag)
			}
			i++
			return args[i], nil
		}

		switch flag {
		case "--deny":
			v, err := value()
			if err != nil {
				return err
			}
			cats, err := robotsCategoryList(v)
			if err != nil {
				return err
			}
			p.Deny = append(p.Deny, cats...)

		case "--allow":
			v, err := value()
			if err != nil {
				return err
			}
			cats, err := robotsCategoryList(v)
			if err != nil {
				return err
			}
			p.Allow = append(p.Allow, cats...)

		case "--deny-provider":
			v, err := value()
			if err != nil {
				return err
			}
			p.DenyProviders = append(p.DenyProviders, strings.Split(v, ",")...)

		case "--path":
			v, err := value()
			if err != nil {
				return err
			}
			p.Path = v

		case "--sitemap":
			v, err := value()
			if err != nil {
				return err
			}
			p.Sitemap = v

		default:
			return fmt.Errorf("unknown flag %q for robots (try: uax help)", flag)
		}
	}

	// Blocking AI training crawlers is what almost everyone running this
	// command wants, so it is what an argument-free invocation does.
	if len(p.Deny) == 0 && len(p.Allow) == 0 && len(p.DenyProviders) == 0 {
		p.Deny = robotsCategories["ai"]
	}

	_, err := io.WriteString(out, robots.Render(p))
	return err
}

func robotsCategoryList(s string) ([]uax.BotCategory, error) {
	var out []uax.BotCategory
	for _, name := range strings.Split(s, ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		cats, ok := robotsCategories[name]
		if !ok {
			return nil, fmt.Errorf("unknown bot category %q (known: %s)", name, robotsCategoryNames())
		}
		out = append(out, cats...)
	}
	return out, nil
}

func robotsCategoryNames() string {
	names := make([]string, 0, len(robotsCategories))
	for k := range robotsCategories {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
