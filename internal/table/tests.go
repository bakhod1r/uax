package table

// RuleTest is an expectation that travels with the rule it proves. The uax
// package runs these as real tests, so a rule and its proof ship together and
// lint_test.go can refuse a rule that has none.
type RuleTest struct {
	RuleID string
	UA     string
	Expect map[string]string
}

var Tests = []RuleTest{
	{RuleID: "headless_chrome", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/120.0.0.0 Safari/537.36", Expect: map[string]string{
		"browser": "Chrome",
	}},
	{RuleID: "electron", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) MyApp/1.2.3 Chrome/120.0.0.0 Electron/28.1.0 Safari/537.36", Expect: map[string]string{
		"browser": "Electron",
	}},
	{RuleID: "chrome_ios", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/120.0.0.0 Mobile/15E148 Safari/604.1", Expect: map[string]string{
		"browser": "Chrome",
	}},
	{RuleID: "duckduckgo_ios", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 DuckDuckGo/7 Mobile/15E148 Safari/605.1.15", Expect: map[string]string{
		"browser": "DuckDuckGo",
	}},
	{RuleID: "edge_ios", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) EdgiOS/120.0.0.0 Mobile/15E148 Safari/605.1.15", Expect: map[string]string{
		"browser": "Edge",
	}},
	{RuleID: "firefox_ios", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/119.0 Mobile/15E148 Safari/605.1.15", Expect: map[string]string{
		"browser": "Firefox",
	}},
	{RuleID: "opera_ios", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) OPiOS/16.0.0 Mobile/15E148 Safari/9537.53", Expect: map[string]string{
		"browser": "Opera",
	}},
	{RuleID: "opera_gx", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPX/106.0.0.0", Expect: map[string]string{
		"browser": "Opera GX",
	}},
	{RuleID: "opera_mini", UA: "Opera/9.80 (Android; Opera Mini/7.5.33361/34.1990; U; en) Presto/2.8.119 Version/11.10", Expect: map[string]string{
		"browser": "Opera Mini",
	}},
	{RuleID: "brave", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Brave/120", Expect: map[string]string{
		"browser": "Brave",
	}},
	{RuleID: "edge", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0", Expect: map[string]string{
		"browser": "Edge",
	}},
	{RuleID: "edge_android", UA: "Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 EdgA/120.0.0.0", Expect: map[string]string{
		"browser": "Edge",
	}},
	{RuleID: "maxthon", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Maxthon/7.1.0", Expect: map[string]string{
		"browser": "Maxthon",
	}},
	{RuleID: "opera", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0", Expect: map[string]string{
		"browser": "Opera",
	}},
	{RuleID: "qq_browser", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 QQBrowser/11.5", Expect: map[string]string{
		"browser": "QQ Browser",
	}},
	{RuleID: "samsung_internet", UA: "Mozilla/5.0 (Linux; Android 14; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/25.0 Chrome/121.0.0.0 Mobile Safari/537.36", Expect: map[string]string{
		"browser": "Samsung Internet",
	}},
	{RuleID: "silk", UA: "Mozilla/5.0 (Linux; Android 9; KFTRWI) AppleWebKit/537.36 (KHTML, like Gecko) Silk/110.4.4 like Chrome/110.0.5481.154 Safari/537.36", Expect: map[string]string{
		"browser": "Amazon Silk",
	}},
	{RuleID: "uc_browser", UA: "Mozilla/5.0 (Linux; U; Android 13; en-US; SM-A536B) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.58 UCBrowser/15.5.0.1204 Mobile Safari/537.36", Expect: map[string]string{
		"browser": "UC Browser",
	}},
	{RuleID: "vivaldi", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Vivaldi/6.5.3206.63", Expect: map[string]string{
		"browser": "Vivaldi",
	}},
	{RuleID: "yandex", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 YaBrowser/24.1.0.0 Safari/537.36", Expect: map[string]string{
		"browser": "Yandex Browser",
	}},
	{RuleID: "edge_legacy", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0 Safari/537.36 Edge/12.246", Expect: map[string]string{
		"browser": "Edge (Legacy)",
	}},
	{RuleID: "chromium", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chromium/120.0.0.0 Chrome/120.0.0.0 Safari/537.36", Expect: map[string]string{
		"browser": "Chromium",
	}},
	{RuleID: "firefox_focus", UA: "Mozilla/5.0 (Android 13; Mobile; rv:109.0) Gecko/117.0 Firefox/117.0 Focus/122.0", Expect: map[string]string{
		"browser": "Firefox Focus",
	}},
	{RuleID: "seamonkey", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0 SeaMonkey/2.53.17", Expect: map[string]string{
		"browser": "SeaMonkey",
	}},
	{RuleID: "waterfox", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0 Waterfox/6.0.6", Expect: map[string]string{
		"browser": "Waterfox",
	}},
	{RuleID: "ie_msie", UA: "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)", Expect: map[string]string{
		"browser": "Internet Explorer",
	}},
	{RuleID: "ie_trident", UA: "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko", Expect: map[string]string{
		"browser": "Internet Explorer",
	}},
	{RuleID: "firefox", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{
		"browser": "Firefox",
	}},
	{RuleID: "chrome", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", Expect: map[string]string{
		"browser": "Chrome",
	}},
	{RuleID: "android_browser", UA: "Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Android Browser/4.0 Mobile Safari/534.30", Expect: map[string]string{
		"browser": "Android Browser",
	}},
	{RuleID: "safari", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15", Expect: map[string]string{
		"browser": "Safari",
	}},
	{RuleID: "edgehtml", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0 Safari/537.36 Edge/12.246", Expect: map[string]string{
		"engine": "EdgeHTML",
	}},
	{RuleID: "presto", UA: "Opera/9.80 (Windows NT 6.0) Presto/2.12.388 Version/12.14", Expect: map[string]string{
		"engine": "Presto",
	}},
	{RuleID: "trident", UA: "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko", Expect: map[string]string{
		"engine": "Trident",
	}},
	{RuleID: "gecko", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{
		"engine": "Gecko",
	}},
	{RuleID: "webkit", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15", Expect: map[string]string{
		"engine": "WebKit",
	}},
	{RuleID: "ipados", UA: "Mozilla/5.0 (iPad; CPU OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1", Expect: map[string]string{
		"os": "iPadOS",
	}},
	{RuleID: "ios_iphone", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1", Expect: map[string]string{
		"os": "iOS",
	}},
	{RuleID: "ios_ipod", UA: "Mozilla/5.0 (iPod touch; CPU iPhone OS 15_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148", Expect: map[string]string{
		"os": "iOS",
	}},
	{RuleID: "tizen", UA: "Mozilla/5.0 (SMART-TV; LINUX; Tizen 6.0) AppleWebKit/537.36 (KHTML, like Gecko) Version/6.0 TV Safari/537.36", Expect: map[string]string{
		"os": "Tizen",
	}},
	{RuleID: "webos", UA: "Mozilla/5.0 (Web0S; Linux/SmartTV) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 WebAppManager", Expect: map[string]string{
		"os": "webOS",
	}},
	{RuleID: "webos_legacy", UA: "Mozilla/5.0 (webOS/1.4.0; U; en-US) AppleWebKit/532.2 (KHTML, like Gecko) Version/1.0 Safari/532.2 Pre/1.0", Expect: map[string]string{
		"os": "webOS",
	}},
	{RuleID: "android_tv", UA: "Mozilla/5.0 (Linux; Android TV 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36", Expect: map[string]string{
		"os": "Android TV",
	}},
	{RuleID: "nintendo", UA: "Mozilla/5.0 (Nintendo Switch; WifiWebAuthApplet) AppleWebKit/609.4 (KHTML, like Gecko) NF/6.0.2.21.3 NintendoBrowser/5.1.0.22433", Expect: map[string]string{
		"os": "Nintendo",
	}},
	{RuleID: "playstation", UA: "Mozilla/5.0 (PlayStation; PlayStation 5/2.26) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Safari/605.1.15", Expect: map[string]string{
		"os": "PlayStation",
	}},
	{RuleID: "xbox", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; Xbox; Xbox Series X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edge/44.18363.8131", Expect: map[string]string{
		"os": "Xbox",
	}},
	{RuleID: "harmonyos", UA: "Mozilla/5.0 (Linux; HarmonyOS 4.0; ELS-NX9) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36", Expect: map[string]string{
		"os": "HarmonyOS",
	}},
	{RuleID: "kaios", UA: "Mozilla/5.0 (Mobile; LYF/F90M/LYF_F90M; Android; rv:48.0; KAIOS/2.5) Gecko/48.0 Firefox/48.0 KAIOS/2.5", Expect: map[string]string{
		"os": "KaiOS",
	}},
	{RuleID: "windows_phone", UA: "Mozilla/5.0 (Mobile; Windows Phone 8.1; Android 4.0; ARM; Trident/7.0) like Gecko", Expect: map[string]string{
		"os": "Windows Phone",
	}},
	{RuleID: "blackberry", UA: "Mozilla/5.0 (BlackBerry; U; BlackBerry 9900; en) AppleWebKit/534.11+ (KHTML, like Gecko) Version/7.1.0.346 Mobile Safari/534.11+", Expect: map[string]string{
		"os": "BlackBerry",
	}},
	{RuleID: "android", UA: "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36", Expect: map[string]string{
		"os": "Android",
	}},
	{RuleID: "windows", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", Expect: map[string]string{
		"os": "Windows",
	}},
	{RuleID: "chromeos", UA: "Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", Expect: map[string]string{
		"os": "Chrome OS",
	}},
	{RuleID: "arch", UA: "Mozilla/5.0 (X11; Arch Linux; Linux x86_64; rv:126.0) Gecko/20100101 Firefox/126.0", Expect: map[string]string{
		"os": "Arch Linux",
	}},
	{RuleID: "debian", UA: "Mozilla/5.0 (X11; Debian; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0", Expect: map[string]string{
		"os": "Debian",
	}},
	{RuleID: "fedora", UA: "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:126.0) Gecko/20100101 Firefox/126.0", Expect: map[string]string{
		"os": "Fedora",
	}},
	{RuleID: "freebsd", UA: "Mozilla/5.0 (X11; FreeBSD amd64; rv:126.0) Gecko/20100101 Firefox/126.0", Expect: map[string]string{
		"os": "FreeBSD",
	}},
	{RuleID: "openbsd", UA: "Mozilla/5.0 (X11; OpenBSD amd64; rv:126.0) Gecko/20100101 Firefox/126.0", Expect: map[string]string{
		"os": "OpenBSD",
	}},
	{RuleID: "ubuntu", UA: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:127.0) Gecko/20100101 Firefox/127.0", Expect: map[string]string{
		"os": "Ubuntu",
	}},
	{RuleID: "macos", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", Expect: map[string]string{
		"os": "macOS",
	}},
	{RuleID: "linux", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", Expect: map[string]string{
		"os": "Linux",
	}},
	{RuleID: "claude_searchbot", UA: "Mozilla/5.0 (compatible; Claude-SearchBot/1.0; +Claude-SearchBot@anthropic.com)", Expect: map[string]string{
		"bot": "Claude-SearchBot",
	}},
	{RuleID: "applebot_extended", UA: "Mozilla/5.0 (compatible; Applebot-Extended/0.1; +http://www.apple.com/go/applebot)", Expect: map[string]string{
		"bot": "Applebot-Extended",
	}},
	{RuleID: "chatgpt_user", UA: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko); compatible; ChatGPT-User/1.0; +https://openai.com/bot", Expect: map[string]string{
		"bot":      "ChatGPT-User",
		"category": "AIAgent",
	}},
	{RuleID: "claude_user", UA: "Mozilla/5.0 (compatible; Claude-User/1.0; +Claude-User@anthropic.com)", Expect: map[string]string{
		"bot":      "Claude-User",
		"category": "AIAgent",
	}},
	{RuleID: "google_extended", UA: "Mozilla/5.0 (compatible; Google-Extended/1.0)", Expect: map[string]string{
		"bot": "Google-Extended",
	}},
	{RuleID: "oai_searchbot", UA: "Mozilla/5.0 (compatible; OAI-SearchBot/1.0; +https://openai.com/searchbot)", Expect: map[string]string{
		"bot": "OAI-SearchBot",
	}},
	{RuleID: "perplexity_user", UA: "Mozilla/5.0 (compatible; Perplexity-User/1.0; +https://perplexity.ai/perplexity-user)", Expect: map[string]string{
		"bot": "Perplexity-User",
	}},
	{RuleID: "amazonbot", UA: "Mozilla/5.0 (compatible; Amazonbot/0.1; +https://developer.amazon.com/support/amazonbot)", Expect: map[string]string{
		"bot": "Amazonbot",
	}},
	{RuleID: "bytespider", UA: "Mozilla/5.0 (compatible; Bytespider; spider-feedback@bytedance.com)", Expect: map[string]string{
		"bot": "Bytespider",
	}},
	{RuleID: "ccbot", UA: "CCBot/2.0 (https://commoncrawl.org/faq/)", Expect: map[string]string{
		"bot": "CCBot",
	}},
	{RuleID: "claudebot", UA: "Mozilla/5.0 (compatible; ClaudeBot/1.0; +claudebot@anthropic.com)", Expect: map[string]string{
		"bot":      "ClaudeBot",
		"category": "AICrawler",
	}},
	{RuleID: "gptbot", UA: "Mozilla/5.0 (compatible; GPTBot/1.2; +https://openai.com/gptbot)", Expect: map[string]string{
		"bot":      "GPTBot",
		"category": "AICrawler",
	}},
	{RuleID: "meta_external_agent", UA: "meta-externalagent/1.1 (+https://developers.facebook.com/docs/sharing/webmasters/crawler)", Expect: map[string]string{
		"bot": "meta-externalagent",
	}},
	{RuleID: "perplexitybot", UA: "Mozilla/5.0 (compatible; PerplexityBot/1.0; +https://perplexity.ai/perplexitybot)", Expect: map[string]string{
		"bot":      "PerplexityBot",
		"provider": "Perplexity",
	}},
	{RuleID: "googlebot_image", UA: "Googlebot-Image/1.0", Expect: map[string]string{
		"bot": "Googlebot-Image",
	}},
	{RuleID: "ahrefsbot", UA: "Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)", Expect: map[string]string{
		"bot":      "AhrefsBot",
		"category": "SEO",
	}},
	{RuleID: "applebot", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)", Expect: map[string]string{
		"bot": "Applebot",
	}},
	{RuleID: "baiduspider", UA: "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)", Expect: map[string]string{
		"bot": "Baiduspider",
	}},
	{RuleID: "bingbot", UA: "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", Expect: map[string]string{
		"bot":      "Bingbot",
		"provider": "Microsoft",
	}},
	{RuleID: "censys", UA: "Mozilla/5.0 (compatible; CensysInspect/1.1; +https://about.censys.io/)", Expect: map[string]string{
		"bot": "Censys",
	}},
	{RuleID: "discordbot", UA: "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)", Expect: map[string]string{
		"bot": "Discordbot",
	}},
	{RuleID: "dotbot", UA: "Mozilla/5.0 (compatible; DotBot/1.2; +https://opensiteexplorer.org/dotbot)", Expect: map[string]string{
		"bot": "DotBot",
	}},
	{RuleID: "duckduckbot", UA: "DuckDuckBot/1.1; (+http://duckduckgo.com/duckduckbot.html)", Expect: map[string]string{
		"bot": "DuckDuckBot",
	}},
	{RuleID: "facebookexternalhit", UA: "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)", Expect: map[string]string{
		"bot":      "facebookexternalhit",
		"category": "SocialCrawler",
	}},
	{RuleID: "googlebot", UA: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", Expect: map[string]string{
		"bot":      "Googlebot",
		"category": "SearchEngine",
	}},
	{RuleID: "linkedinbot", UA: "LinkedInBot/1.0 (compatible; Mozilla/5.0; Apache-HttpClient +http://www.linkedin.com)", Expect: map[string]string{
		"bot": "LinkedInBot",
	}},
	{RuleID: "mj12bot", UA: "Mozilla/5.0 (compatible; MJ12bot/v1.4.8; http://mj12bot.com/)", Expect: map[string]string{
		"bot": "MJ12bot",
	}},
	{RuleID: "pingdom", UA: "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)", Expect: map[string]string{
		"bot": "Pingdom",
	}},
	{RuleID: "pinterestbot", UA: "Mozilla/5.0 (compatible; Pinterestbot/1.0; +http://www.pinterest.com/bot.html)", Expect: map[string]string{
		"bot": "Pinterestbot",
	}},
	{RuleID: "semrushbot", UA: "Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)", Expect: map[string]string{
		"bot": "SemrushBot",
	}},
	{RuleID: "shodan", UA: "Mozilla/5.0 (compatible; Shodan; +https://www.shodan.io)", Expect: map[string]string{
		"bot": "Shodan",
	}},
	{RuleID: "slackbot", UA: "Slackbot-LinkExpanding 1.0 (+https://api.slack.com/robots)", Expect: map[string]string{
		"bot": "Slackbot",
	}},
	{RuleID: "slurp", UA: "Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)", Expect: map[string]string{
		"bot": "Yahoo! Slurp",
	}},
	{RuleID: "telegrambot", UA: "TelegramBot (like TwitterBot)", Expect: map[string]string{
		"bot": "TelegramBot",
	}},
	{RuleID: "twitterbot", UA: "Twitterbot/1.0", Expect: map[string]string{
		"bot": "Twitterbot",
	}},
	{RuleID: "uptimerobot", UA: "Mozilla/5.0 (compatible; UptimeRobot/2.0; http://www.uptimerobot.com/)", Expect: map[string]string{
		"bot": "UptimeRobot",
	}},
	{RuleID: "yandexbot", UA: "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)", Expect: map[string]string{
		"bot": "YandexBot",
	}},
	{RuleID: "whatsapp", UA: "WhatsApp/2.23.20.0 A", Expect: map[string]string{
		"bot": "WhatsApp",
	}},
	{RuleID: "generic_bot", UA: "SomeRandomBot/1.0 (+http://example.com/bot)", Expect: map[string]string{
		"category": "Generic",
	}},
	{RuleID: "generic_crawler", UA: "custom-crawler/0.1", Expect: map[string]string{
		"category": "Generic",
	}},
	{RuleID: "generic_spider", UA: "somespider/0.1", Expect: map[string]string{
		"category": "Generic",
	}},
	{RuleID: "instagram_app", UA: "Mozilla/5.0 (Linux; Android 13; SM-S911B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 Instagram 300.0.0.0.0 Android", Expect: map[string]string{
		"app": "Instagram",
	}},
	{RuleID: "facebook_app", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [FBAN/FBIOS;FBAV/440.0.0.0]", Expect: map[string]string{
		"app": "Facebook",
	}},
	{RuleID: "facebook_app_android", UA: "Mozilla/5.0 (Linux; Android 13; SM-S911B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/440.0.0.0;]", Expect: map[string]string{
		"app": "Facebook",
	}},
	{RuleID: "line_app", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Safari/604.1 Line/13.19.0", Expect: map[string]string{
		"app": "LINE",
	}},
	{RuleID: "linkedin_app", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 LinkedInApp/9.28.0", Expect: map[string]string{
		"app": "LinkedIn",
	}},
	{RuleID: "pinterest_app", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Pinterest/11.0", Expect: map[string]string{
		"app": "Pinterest",
	}},
	{RuleID: "slack_app", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Slack_SSB/4.36.0 Chrome/120.0.0.0 Electron/28.0.0 Safari/537.36", Expect: map[string]string{
		"app": "Slack",
	}},
	{RuleID: "snapchat_app", UA: "Mozilla/5.0 (Linux; Android 13; SM-S911B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 Snapchat/12.50.0.0", Expect: map[string]string{
		"app": "Snapchat",
	}},
	{RuleID: "telegram_app", UA: "Mozilla/5.0 (Linux; Android 13; SM-S911B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 TelegramWebview/1.0", Expect: map[string]string{
		"app": "Telegram",
	}},
	{RuleID: "tiktok_app", UA: "Mozilla/5.0 (Linux; Android 13; SM-S911B) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/120.0.0.0 Mobile Safari/537.36 BytedanceWebview/d8a21c6", Expect: map[string]string{
		"app": "TikTok",
	}},
	{RuleID: "twitter_app", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Twitter for iPhone/10.30", Expect: map[string]string{
		"app": "X",
	}},
	{RuleID: "wechat_app", UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.42", Expect: map[string]string{
		"app": "WeChat",
	}},
	{RuleID: "cypress", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Cypress/13.6", Expect: map[string]string{
		"automation": "Cypress",
	}},
	{RuleID: "phantomjs", UA: "Mozilla/5.0 (Unknown; Linux x86_64) AppleWebKit/538.1 (KHTML, like Gecko) PhantomJS/2.1.1 Safari/538.1", Expect: map[string]string{
		"automation": "PhantomJS",
	}},
	{RuleID: "playwright", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Playwright/1.40", Expect: map[string]string{
		"automation": "Playwright",
	}},
	{RuleID: "puppeteer", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Puppeteer/21.0", Expect: map[string]string{
		"automation": "Puppeteer",
	}},
	{RuleID: "selenium", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 selenium/4.16", Expect: map[string]string{
		"automation": "Selenium",
	}},
	{RuleID: "webdriver", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 WebDriver", Expect: map[string]string{
		"automation": "WebDriver",
	}},
	{RuleID: "headless_chrome", UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/120.0.0.0 Safari/537.36", Expect: map[string]string{
		"automation": "Headless Chrome",
	}},
	{RuleID: "headless_firefox", UA: "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0 Headless", Expect: map[string]string{
		"automation": "Headless Firefox",
	}},
	{RuleID: "electron_automation", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) MyApp/1.2.3 Chrome/120.0.0.0 Electron/28.1.0 Safari/537.36", Expect: map[string]string{
		"automation": "Electron",
	}},
	{RuleID: "aiohttp", UA: "Python/3.11 aiohttp/3.9.1", Expect: map[string]string{
		"client": "aiohttp",
	}},
	{RuleID: "apache_httpclient", UA: "Apache-HttpClient/4.5.14 (Java/17.0.9)", Expect: map[string]string{
		"client": "Apache-HttpClient",
	}},
	{RuleID: "axios", UA: "axios/1.6.5", Expect: map[string]string{
		"client": "axios",
	}},
	{RuleID: "curl", UA: "curl/8.4.0", Expect: map[string]string{
		"category": "CLI",
		"client":   "curl",
	}},
	{RuleID: "go_http", UA: "Go-http-client/2.0", Expect: map[string]string{
		"category": "Library",
		"client":   "Go-http-client",
	}},
	{RuleID: "guzzle", UA: "GuzzleHttp/7.8.1", Expect: map[string]string{
		"client": "Guzzle",
	}},
	{RuleID: "httpie", UA: "HTTPie/3.2.2", Expect: map[string]string{
		"client": "HTTPie",
	}},
	{RuleID: "insomnia", UA: "insomnia/8.6.1", Expect: map[string]string{
		"client": "Insomnia",
	}},
	{RuleID: "libwww_perl", UA: "libwww-perl/6.72", Expect: map[string]string{
		"client": "libwww-perl",
	}},
	{RuleID: "node_fetch", UA: "node-fetch/1.0 (+https://github.com/bitinn/node-fetch)", Expect: map[string]string{
		"client": "node-fetch",
	}},
	{RuleID: "okhttp", UA: "okhttp/4.12.0", Expect: map[string]string{
		"client": "okhttp",
	}},
	{RuleID: "postman", UA: "PostmanRuntime/7.36.0", Expect: map[string]string{
		"client": "PostmanRuntime",
	}},
	{RuleID: "python_httpx", UA: "python-httpx/0.26.0", Expect: map[string]string{
		"client": "httpx",
	}},
	{RuleID: "python_requests", UA: "python-requests/2.31.0", Expect: map[string]string{
		"client": "python-requests",
	}},
	{RuleID: "python_urllib", UA: "Python-urllib/3.11", Expect: map[string]string{
		"client": "urllib",
	}},
	{RuleID: "wget", UA: "Wget/1.21.4", Expect: map[string]string{
		"client": "Wget",
	}},
	{RuleID: "xh", UA: "xh/0.20.1", Expect: map[string]string{
		"client": "xh",
	}},
	{RuleID: "got", UA: "got (https://github.com/sindresorhus/got)", Expect: map[string]string{
		"client": "got",
	}},
	{RuleID: "java", UA: "Java/17.0.9", Expect: map[string]string{
		"client": "Java",
	}},
	{RuleID: "oculus_browser", UA: "Mozilla/5.0 (X11; Linux x86_64; Quest 3) AppleWebKit/537.36 (KHTML, like Gecko) OculusBrowser/33.0.0.1.1 SamsungBrowser/4.0 Chrome/126.0.0.0 VR Safari/537.36", Expect: map[string]string{"browser": "Meta Quest Browser"}},
	{RuleID: "coc_coc", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) coc_coc_browser/126.0.0 Chrome/120.0.0.0 Safari/537.36", Expect: map[string]string{"browser": "Coc Coc"}},
	{RuleID: "whale", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Whale/3.26.244.21 Safari/537.36", Expect: map[string]string{"browser": "Whale"}},
	{RuleID: "huawei_browser", UA: "Mozilla/5.0 (Linux; Android 12; ELS-NX9) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 HuaweiBrowser/14.0.5.310 Mobile Safari/537.36", Expect: map[string]string{"browser": "Huawei Browser"}},
	{RuleID: "miui_browser", UA: "Mozilla/5.0 (Linux; U; Android 13; 2201116SG) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.0.0 Mobile Safari/537.36 XiaoMi/MiuiBrowser/17.2.71", Expect: map[string]string{"browser": "MIUI Browser"}},
	{RuleID: "heytap_browser", UA: "Mozilla/5.0 (Linux; U; Android 13; CPH2451) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.0.0 Mobile Safari/537.36 HeyTapBrowser/45.11.4.1", Expect: map[string]string{"browser": "OPPO Browser"}},
	{RuleID: "vivo_browser", UA: "Mozilla/5.0 (Linux; Android 13; V2312A) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.0.0 Mobile Safari/537.36 VivoBrowser/22.1.0.0", Expect: map[string]string{"browser": "vivo Browser"}},
	{RuleID: "quark", UA: "Mozilla/5.0 (Linux; U; Android 13; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.0.0 Mobile Safari/537.36 Quark/7.2.5.620", Expect: map[string]string{"browser": "Quark"}},
	{RuleID: "baidu_browser", UA: "Mozilla/5.0 (Linux; Android 12; SM-A536B) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/114.0.0.0 Mobile Safari/537.36 baidubrowser/13.42.0.10", Expect: map[string]string{"browser": "Baidu Browser"}},
	{RuleID: "puffin", UA: "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36 Puffin/10.0.1.51703AP", Expect: map[string]string{"browser": "Puffin"}},
	{RuleID: "opera_touch", UA: "Mozilla/5.0 (Linux; Android 13; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36 OPT/3.4.6", Expect: map[string]string{"browser": "Opera Touch"}},
	{RuleID: "sogou_explorer", UA: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36 SE 2.X MetaSr 1.0", Expect: map[string]string{"browser": "Sogou Explorer"}},
	{RuleID: "qihoo_360ee", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 QIHU 360EE", Expect: map[string]string{"browser": "360 Extreme Explorer"}},
	{RuleID: "qihoo_360se", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 QIHU 360SE", Expect: map[string]string{"browser": "360 Secure Browser"}},
	{RuleID: "srware_iron", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Iron/120.0.6150.0 Chrome/120.0.6150.0 Safari/537.36", Expect: map[string]string{"browser": "SRWare Iron"}},
	{RuleID: "librewolf", UA: "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0 LibreWolf/128.0-1", Expect: map[string]string{"browser": "LibreWolf"}},
	{RuleID: "palemoon", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:4.8) Goanna/20240301 PaleMoon/33.0.2", Expect: map[string]string{"browser": "Pale Moon"}},
	{RuleID: "basilisk", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:52.9) Gecko/20100101 Goanna/6.4 Firefox/52.9 Basilisk/20240211", Expect: map[string]string{"browser": "Basilisk"}},
	{RuleID: "firefox_klar", UA: "Mozilla/5.0 (Android 13; Mobile; rv:128.0) Gecko/128.0 Firefox/128.0 Klar/128.0", Expect: map[string]string{"browser": "Firefox Klar"}},
	{RuleID: "konqueror", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) konqueror/22.12.3 Konqueror/22.12.3 Safari/537.36", Expect: map[string]string{"browser": "Konqueror"}},
	{RuleID: "epiphany", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Epiphany/605.1.15", Expect: map[string]string{"browser": "GNOME Web"}},
	{RuleID: "falkon", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.144 Safari/537.36 Falkon/24.02.2", Expect: map[string]string{"browser": "Falkon"}},
	{RuleID: "midori", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Midori/11.2.1 Safari/537.36", Expect: map[string]string{"browser": "Midori"}},
	{RuleID: "sleipnir", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Sleipnir/6.4.13", Expect: map[string]string{"browser": "Sleipnir"}},
	{RuleID: "goanna", UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:4.8) Goanna/20240301 PaleMoon/33.0.2", Expect: map[string]string{"engine": "Goanna"}},
	{RuleID: "khtml", UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) konqueror/22.12.3 Konqueror/22.12.3 Safari/537.36", Expect: map[string]string{"engine": "KHTML"}},
	{RuleID: "visionos", UA: "Mozilla/5.0 (Vision OS 2.0; like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/2.0 Safari/605.1.15", Expect: map[string]string{"os": "visionOS"}},
	{RuleID: "tvos", UA: "AppleTV6,2/11.1", Expect: map[string]string{"os": "tvOS"}},
	{RuleID: "watchos", UA: "Mozilla/5.0 (Apple Watch; CPU Watch OS 10_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)", Expect: map[string]string{"os": "watchOS"}},
	{RuleID: "roku", UA: "Roku4640X/DVP-7.70 (297.70E04154A)", Expect: map[string]string{"os": "Roku OS"}},
	{RuleID: "chromecast", UA: "Mozilla/5.0 (X11; Linux armv7l) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 CrKey/1.56.500000", Expect: map[string]string{"os": "Chromecast"}},
	{RuleID: "sailfish", UA: "Mozilla/5.0 (Wayland; Linux x86_64; Sailfish 4.5.0.18) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.2 Safari/605.1.15", Expect: map[string]string{"os": "Sailfish OS"}},
	{RuleID: "symbian", UA: "Mozilla/5.0 (SymbianOS/9.4; Series60/5.0 NokiaN97-1/12.0.024; Profile/MIDP-2.1 Configuration/CLDC-1.1) AppleWebKit/525 (KHTML, like Gecko) BrowserNG/7.1.12344", Expect: map[string]string{"os": "Symbian"}},
	{RuleID: "meego", UA: "Mozilla/5.0 (MeeGo; NokiaN9) AppleWebKit/534.13 (KHTML, like Gecko) NokiaBrowser/8.5.0 Mobile Safari/534.13", Expect: map[string]string{"os": "MeeGo"}},
	{RuleID: "windows_ce", UA: "Mozilla/4.0 (compatible; MSIE 6.0; Windows CE; IEMobile 7.11)", Expect: map[string]string{"os": "Windows CE"}},
	{RuleID: "steamos", UA: "Mozilla/5.0 (X11; Linux x86_64; SteamOS 3.5.19) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36", Expect: map[string]string{"os": "SteamOS"}},
	{RuleID: "linux_mint", UA: "Mozilla/5.0 (X11; Linux Mint; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{"os": "Linux Mint"}},
	{RuleID: "centos", UA: "Mozilla/5.0 (X11; CentOS; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0", Expect: map[string]string{"os": "CentOS"}},
	{RuleID: "gentoo", UA: "Mozilla/5.0 (X11; Gentoo; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{"os": "Gentoo"}},
	{RuleID: "manjaro", UA: "Mozilla/5.0 (X11; Manjaro Linux; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{"os": "Manjaro Linux"}},
	{RuleID: "netbsd", UA: "Mozilla/5.0 (X11; NetBSD amd64; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{"os": "NetBSD"}},
	{RuleID: "opensuse", UA: "Mozilla/5.0 (X11; openSUSE Tumbleweed; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{"os": "openSUSE"}},
	{RuleID: "slackware", UA: "Mozilla/5.0 (X11; Slackware; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0", Expect: map[string]string{"os": "Slackware"}},
	{RuleID: "solaris", UA: "Mozilla/5.0 (X11; SunOS i86pc; rv:52.0) Gecko/20100101 Firefox/52.0", Expect: map[string]string{"os": "Solaris"}},
	{RuleID: "haiku", UA: "Mozilla/5.0 (X11; Haiku x86_64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15", Expect: map[string]string{"os": "Haiku"}},
}
