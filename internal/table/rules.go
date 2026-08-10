// Package table is the classification database.
//
// The rules live here, in Go, and are the source of truth. There is no
// generator and no data file: a rule is a struct literal, reviewed as code,
// and lint_test.go enforces the invariants that keep the set trustworthy —
// unique IDs, exactly one match clause, a priority, a confidence in (0,1],
// and at least one test per rule in tests.go.
//
// Hardware is deliberately absent. Naming a handset is a catalogue problem on
// a different release schedule, so it lives in github.com/bakhod1r/devicex and
// reaches the parser through uax.Config.Device.
package table

// Rule is one classification rule. Rules are sorted by descending priority, so
// the first match wins and evaluation can stop.
type Rule struct {
	ID          string
	Name        string
	Family      string
	Engine      string
	Brand       string
	Model       string
	Series      string
	Category    string
	Provider    string
	Type        string
	Token       string
	Prefix      string
	Contains    string
	VersionFrom string
	Year        int
	Priority    int
	Confidence  float32
}

var Browsers = []Rule{
	{ID: "headless_chrome", Name: "Chrome", Family: "Chrome", Engine: "Blink", Token: "HeadlessChrome", VersionFrom: "token", Priority: 300, Confidence: 0.98},
	{ID: "electron", Name: "Electron", Family: "Chrome", Engine: "Blink", Token: "Electron", VersionFrom: "token", Priority: 260, Confidence: 0.97},
	{ID: "chrome_ios", Name: "Chrome", Family: "Chrome", Engine: "WebKit", Token: "CriOS", VersionFrom: "token", Priority: 250, Confidence: 0.97},
	{ID: "duckduckgo_ios", Name: "DuckDuckGo", Family: "DuckDuckGo", Engine: "WebKit", Token: "DuckDuckGo", VersionFrom: "token", Priority: 250, Confidence: 0.95},
	{ID: "edge_ios", Name: "Edge", Family: "Edge", Engine: "WebKit", Token: "EdgiOS", VersionFrom: "token", Priority: 250, Confidence: 0.97},
	{ID: "firefox_ios", Name: "Firefox", Family: "Firefox", Engine: "WebKit", Token: "FxiOS", VersionFrom: "token", Priority: 250, Confidence: 0.97},
	{ID: "opera_ios", Name: "Opera", Family: "Opera", Engine: "WebKit", Token: "OPiOS", VersionFrom: "token", Priority: 250, Confidence: 0.97},
	{ID: "opera_gx", Name: "Opera GX", Family: "Opera", Engine: "Blink", Token: "OPX", VersionFrom: "token", Priority: 230, Confidence: 0.96},
	{ID: "opera_mini", Name: "Opera Mini", Family: "Opera", Engine: "Presto", Token: "Opera Mini", VersionFrom: "token", Priority: 230, Confidence: 0.96},
	{ID: "oculus_browser", Name: "Meta Quest Browser", Family: "Chrome", Engine: "Blink", Token: "OculusBrowser", VersionFrom: "token", Priority: 225, Confidence: 0.96},
	{ID: "coc_coc", Name: "Coc Coc", Family: "Chrome", Engine: "Blink", Token: "coc_coc_browser", VersionFrom: "token", Priority: 222, Confidence: 0.96},
	{ID: "whale", Name: "Whale", Family: "Chrome", Engine: "Blink", Token: "Whale", VersionFrom: "token", Priority: 221, Confidence: 0.96},
	{ID: "baidu_browser", Name: "Baidu Browser", Family: "Chrome", Engine: "Blink", Token: "baidubrowser", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "brave", Name: "Brave", Family: "Chrome", Engine: "Blink", Token: "Brave", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "edge", Name: "Edge", Family: "Edge", Engine: "Blink", Token: "Edg", VersionFrom: "token", Priority: 220, Confidence: 0.97},
	{ID: "edge_android", Name: "Edge", Family: "Edge", Engine: "Blink", Token: "EdgA", VersionFrom: "token", Priority: 220, Confidence: 0.97},
	{ID: "heytap_browser", Name: "OPPO Browser", Family: "Chrome", Engine: "Blink", Token: "HeyTapBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "huawei_browser", Name: "Huawei Browser", Family: "Chrome", Engine: "Blink", Token: "HuaweiBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.96},
	{ID: "maxthon", Name: "Maxthon", Family: "Chrome", Engine: "Blink", Token: "Maxthon", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "opera", Name: "Opera", Family: "Opera", Engine: "Blink", Token: "OPR", VersionFrom: "token", Priority: 220, Confidence: 0.97},
	{ID: "opera_touch", Name: "Opera Touch", Family: "Opera", Engine: "Blink", Token: "OPT", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "puffin", Name: "Puffin", Family: "Chrome", Engine: "Blink", Token: "Puffin", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "qq_browser", Name: "QQ Browser", Family: "QQ", Engine: "Blink", Token: "QQBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.96},
	{ID: "quark", Name: "Quark", Family: "Chrome", Engine: "Blink", Token: "Quark", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "samsung_internet", Name: "Samsung Internet", Family: "Chrome", Engine: "Blink", Token: "SamsungBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.97},
	{ID: "silk", Name: "Amazon Silk", Family: "Chrome", Engine: "Blink", Token: "Silk", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "uc_browser", Name: "UC Browser", Family: "UC", Engine: "Blink", Token: "UCBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.96},
	{ID: "vivaldi", Name: "Vivaldi", Family: "Chrome", Engine: "Blink", Token: "Vivaldi", VersionFrom: "token", Priority: 220, Confidence: 0.97},
	{ID: "vivo_browser", Name: "vivo Browser", Family: "Chrome", Engine: "Blink", Token: "VivoBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.95},
	{ID: "yandex", Name: "Yandex Browser", Family: "Chrome", Engine: "Blink", Token: "YaBrowser", VersionFrom: "token", Priority: 220, Confidence: 0.97},
	{ID: "miui_browser", Name: "MIUI Browser", Family: "Chrome", Engine: "Blink", Contains: "MiuiBrowser", Priority: 219, Confidence: 0.93},
	{ID: "qihoo_360ee", Name: "360 Extreme Explorer", Family: "Chrome", Engine: "Blink", Contains: "360EE", Priority: 218, Confidence: 0.92},
	{ID: "sogou_explorer", Name: "Sogou Explorer", Family: "Chrome", Engine: "Blink", Contains: "MetaSr", Priority: 218, Confidence: 0.92},
	{ID: "qihoo_360se", Name: "360 Secure Browser", Family: "Chrome", Engine: "Blink", Contains: "360SE", Priority: 217, Confidence: 0.92},
	{ID: "srware_iron", Name: "SRWare Iron", Family: "Chrome", Engine: "Blink", Token: "Iron", VersionFrom: "token", Priority: 216, Confidence: 0.93},
	{ID: "edge_legacy", Name: "Edge (Legacy)", Family: "Edge", Engine: "EdgeHTML", Token: "Edge", VersionFrom: "token", Priority: 215, Confidence: 0.96},
	{ID: "basilisk", Name: "Basilisk", Family: "Firefox", Engine: "Goanna", Token: "Basilisk", VersionFrom: "token", Priority: 210, Confidence: 0.94},
	{ID: "chromium", Name: "Chromium", Family: "Chrome", Engine: "Blink", Token: "Chromium", VersionFrom: "token", Priority: 210, Confidence: 0.96},
	{ID: "firefox_focus", Name: "Firefox Focus", Family: "Firefox", Engine: "Gecko", Token: "Focus", VersionFrom: "token", Priority: 210, Confidence: 0.95},
	{ID: "firefox_klar", Name: "Firefox Klar", Family: "Firefox", Engine: "Gecko", Token: "Klar", VersionFrom: "token", Priority: 210, Confidence: 0.94},
	{ID: "librewolf", Name: "LibreWolf", Family: "Firefox", Engine: "Gecko", Token: "LibreWolf", VersionFrom: "token", Priority: 210, Confidence: 0.95},
	{ID: "palemoon", Name: "Pale Moon", Family: "Firefox", Engine: "Goanna", Token: "PaleMoon", VersionFrom: "token", Priority: 210, Confidence: 0.95},
	{ID: "seamonkey", Name: "SeaMonkey", Family: "Firefox", Engine: "Gecko", Token: "SeaMonkey", VersionFrom: "token", Priority: 210, Confidence: 0.95},
	{ID: "waterfox", Name: "Waterfox", Family: "Firefox", Engine: "Gecko", Token: "Waterfox", VersionFrom: "token", Priority: 210, Confidence: 0.95},
	{ID: "epiphany", Name: "GNOME Web", Family: "Epiphany", Engine: "WebKit", Token: "Epiphany", VersionFrom: "token", Priority: 205, Confidence: 0.94},
	{ID: "falkon", Name: "Falkon", Family: "Falkon", Engine: "Blink", Token: "Falkon", VersionFrom: "token", Priority: 205, Confidence: 0.94},
	{ID: "konqueror", Name: "Konqueror", Family: "Konqueror", Engine: "KHTML", Token: "Konqueror", VersionFrom: "token", Priority: 205, Confidence: 0.95},
	{ID: "midori", Name: "Midori", Family: "Midori", Engine: "WebKit", Token: "Midori", VersionFrom: "token", Priority: 205, Confidence: 0.93},
	{ID: "sleipnir", Name: "Sleipnir", Family: "Chrome", Engine: "Blink", Token: "Sleipnir", VersionFrom: "token", Priority: 205, Confidence: 0.93},
	{ID: "ie_msie", Name: "Internet Explorer", Family: "IE", Engine: "Trident", Token: "MSIE", VersionFrom: "suffix", Priority: 200, Confidence: 0.96},
	{ID: "ie_trident", Name: "Internet Explorer", Family: "IE", Engine: "Trident", Token: "Trident", VersionFrom: "rv", Priority: 200, Confidence: 0.96},
	{ID: "firefox", Name: "Firefox", Family: "Firefox", Engine: "Gecko", Token: "Firefox", VersionFrom: "token", Priority: 150, Confidence: 0.97},
	{ID: "chrome", Name: "Chrome", Family: "Chrome", Engine: "Blink", Token: "Chrome", VersionFrom: "token", Priority: 140, Confidence: 0.95},
	{ID: "android_browser", Name: "Android Browser", Family: "Android", Engine: "WebKit", Contains: "Android Browser", VersionFrom: "token", Priority: 130, Confidence: 0.9},
	{ID: "safari", Name: "Safari", Family: "Safari", Engine: "WebKit", Token: "Safari", VersionFrom: "version_product", Priority: 100, Confidence: 0.92},
}

var Engines = []Rule{
	{ID: "goanna", Name: "Goanna", Token: "Goanna", VersionFrom: "token", Priority: 210, Confidence: 0.96},
	{ID: "edgehtml", Name: "EdgeHTML", Token: "Edge", VersionFrom: "token", Priority: 200, Confidence: 0.95},
	{ID: "presto", Name: "Presto", Token: "Presto", VersionFrom: "token", Priority: 200, Confidence: 0.97},
	{ID: "trident", Name: "Trident", Token: "Trident", VersionFrom: "token", Priority: 200, Confidence: 0.97},
	{ID: "gecko", Name: "Gecko", Token: "Firefox", VersionFrom: "token", Priority: 150, Confidence: 0.95},
	{ID: "khtml", Name: "KHTML", Token: "Konqueror", VersionFrom: "token", Priority: 120, Confidence: 0.9},
	{ID: "webkit", Name: "WebKit", Token: "AppleWebKit", VersionFrom: "token", Priority: 100, Confidence: 0.8},
}

var OSes = []Rule{
	{ID: "visionos", Name: "visionOS", Family: "visionOS", Type: "XR", Contains: "Vision OS", VersionFrom: "apple", Priority: 330, Confidence: 0.94},
	{ID: "tvos", Name: "tvOS", Family: "tvOS", Type: "TV", Contains: "AppleTV", Priority: 325, Confidence: 0.94},
	{ID: "watchos", Name: "watchOS", Family: "watchOS", Type: "Watch", Contains: "Watch OS", VersionFrom: "apple", Priority: 324, Confidence: 0.93},
	{ID: "ipados", Name: "iPadOS", Family: "iOS", Type: "Tablet", Token: "iPad", VersionFrom: "apple", Priority: 320, Confidence: 0.96},
	{ID: "ios_iphone", Name: "iOS", Family: "iOS", Type: "Mobile", Prefix: "CPU iPhone OS", VersionFrom: "apple", Priority: 310, Confidence: 0.97},
	{ID: "ios_ipod", Name: "iOS", Family: "iOS", Type: "Mobile", Prefix: "iPod", VersionFrom: "apple", Priority: 300, Confidence: 0.93},
	{ID: "roku", Name: "Roku OS", Family: "Roku", Type: "TV", Contains: "Roku", Priority: 292, Confidence: 0.94},
	{ID: "chromecast", Name: "Chromecast", Family: "Chrome OS", Type: "TV", Contains: "CrKey", Priority: 291, Confidence: 0.94},
	{ID: "tizen", Name: "Tizen", Family: "Tizen", Type: "TV", Prefix: "Tizen", VersionFrom: "suffix", Priority: 290, Confidence: 0.95},
	{ID: "webos", Name: "webOS", Family: "webOS", Type: "TV", Prefix: "Web0S", Priority: 290, Confidence: 0.94},
	{ID: "webos_legacy", Name: "webOS", Family: "webOS", Type: "TV", Prefix: "webOS", Priority: 285, Confidence: 0.94},
	{ID: "android_tv", Name: "Android TV", Family: "Android", Type: "TV", Prefix: "Android TV", VersionFrom: "suffix", Priority: 280, Confidence: 0.95},
	{ID: "nintendo", Name: "Nintendo", Family: "Nintendo", Type: "Console", Contains: "Nintendo", Priority: 280, Confidence: 0.95},
	{ID: "playstation", Name: "PlayStation", Family: "PlayStation", Type: "Console", Contains: "PlayStation", Priority: 280, Confidence: 0.95},
	{ID: "xbox", Name: "Xbox", Family: "Xbox", Type: "Console", Contains: "Xbox", Priority: 280, Confidence: 0.95},
	{ID: "sailfish", Name: "Sailfish OS", Family: "Sailfish", Type: "Mobile", Contains: "Sailfish", Priority: 272, Confidence: 0.93},
	{ID: "meego", Name: "MeeGo", Family: "MeeGo", Type: "Mobile", Contains: "MeeGo", Priority: 271, Confidence: 0.92},
	{ID: "symbian", Name: "Symbian", Family: "Symbian", Type: "Mobile", Contains: "Symbian", Priority: 271, Confidence: 0.93},
	{ID: "harmonyos", Name: "HarmonyOS", Family: "HarmonyOS", Type: "Mobile", Prefix: "HarmonyOS", VersionFrom: "suffix", Priority: 270, Confidence: 0.94},
	{ID: "kaios", Name: "KaiOS", Family: "KaiOS", Type: "Mobile", Prefix: "KAIOS", VersionFrom: "suffix", Priority: 270, Confidence: 0.93},
	{ID: "windows_phone", Name: "Windows Phone", Family: "Windows", Type: "Mobile", Prefix: "Windows Phone", VersionFrom: "suffix", Priority: 270, Confidence: 0.95},
	{ID: "windows_ce", Name: "Windows CE", Family: "Windows", Type: "Mobile", Contains: "Windows CE", Priority: 269, Confidence: 0.92},
	{ID: "blackberry", Name: "BlackBerry", Family: "BlackBerry", Type: "Mobile", Prefix: "BlackBerry", Priority: 265, Confidence: 0.92},
	{ID: "android", Name: "Android", Family: "Android", Prefix: "Android", VersionFrom: "suffix", Priority: 200, Confidence: 0.96},
	{ID: "windows", Name: "Windows", Family: "Windows", Type: "Desktop", Prefix: "Windows NT", VersionFrom: "nt", Priority: 190, Confidence: 0.97},
	{ID: "steamos", Name: "SteamOS", Family: "Linux", Type: "Console", Contains: "SteamOS", Priority: 186, Confidence: 0.94},
	{ID: "chromeos", Name: "Chrome OS", Family: "Chrome OS", Type: "Desktop", Prefix: "CrOS", VersionFrom: "cros", Priority: 185, Confidence: 0.96},
	{ID: "linux_mint", Name: "Linux Mint", Family: "Linux", Type: "Desktop", Token: "Linux Mint", Priority: 181, Confidence: 0.94},
	{ID: "arch", Name: "Arch Linux", Family: "Linux", Type: "Desktop", Token: "Arch Linux", Priority: 180, Confidence: 0.94},
	{ID: "centos", Name: "CentOS", Family: "Linux", Type: "Desktop", Token: "CentOS", Priority: 180, Confidence: 0.94},
	{ID: "debian", Name: "Debian", Family: "Linux", Type: "Desktop", Token: "Debian", Priority: 180, Confidence: 0.94},
	{ID: "fedora", Name: "Fedora", Family: "Linux", Type: "Desktop", Token: "Fedora", Priority: 180, Confidence: 0.94},
	{ID: "freebsd", Name: "FreeBSD", Family: "BSD", Type: "Desktop", Prefix: "FreeBSD", Priority: 180, Confidence: 0.95},
	{ID: "gentoo", Name: "Gentoo", Family: "Linux", Type: "Desktop", Token: "Gentoo", Priority: 180, Confidence: 0.94},
	{ID: "haiku", Name: "Haiku", Family: "Haiku", Type: "Desktop", Prefix: "Haiku", Priority: 180, Confidence: 0.94},
	{ID: "manjaro", Name: "Manjaro Linux", Family: "Linux", Type: "Desktop", Token: "Manjaro Linux", Priority: 180, Confidence: 0.94},
	{ID: "netbsd", Name: "NetBSD", Family: "BSD", Type: "Desktop", Prefix: "NetBSD", Priority: 180, Confidence: 0.95},
	{ID: "openbsd", Name: "OpenBSD", Family: "BSD", Type: "Desktop", Prefix: "OpenBSD", Priority: 180, Confidence: 0.95},
	{ID: "opensuse", Name: "openSUSE", Family: "Linux", Type: "Desktop", Token: "openSUSE", Priority: 180, Confidence: 0.94},
	{ID: "slackware", Name: "Slackware", Family: "Linux", Type: "Desktop", Token: "Slackware", Priority: 180, Confidence: 0.94},
	{ID: "solaris", Name: "Solaris", Family: "Solaris", Type: "Desktop", Prefix: "SunOS", Priority: 180, Confidence: 0.94},
	{ID: "ubuntu", Name: "Ubuntu", Family: "Linux", Type: "Desktop", Token: "Ubuntu", Priority: 180, Confidence: 0.94},
	{ID: "macos", Name: "macOS", Family: "macOS", Type: "Desktop", Contains: "Mac OS X", VersionFrom: "apple", Priority: 150, Confidence: 0.96},
	{ID: "linux", Name: "Linux", Family: "Linux", Type: "Desktop", Prefix: "Linux", Priority: 100, Confidence: 0.85},
}

var Bots = []Rule{
	{ID: "claude_searchbot", Name: "Claude-SearchBot", Category: "AICrawler", Provider: "Anthropic", Token: "Claude-SearchBot", VersionFrom: "token", Priority: 206, Confidence: 0.98},
	{ID: "applebot_extended", Name: "Applebot-Extended", Category: "AICrawler", Provider: "Apple", Token: "Applebot-Extended", Priority: 205, Confidence: 0.97},
	{ID: "chatgpt_user", Name: "ChatGPT-User", Category: "AIAgent", Provider: "OpenAI", Token: "ChatGPT-User", VersionFrom: "token", Priority: 205, Confidence: 0.98},
	{ID: "claude_user", Name: "Claude-User", Category: "AIAgent", Provider: "Anthropic", Token: "Claude-User", VersionFrom: "token", Priority: 205, Confidence: 0.98},
	{ID: "google_extended", Name: "Google-Extended", Category: "AICrawler", Provider: "Google", Token: "Google-Extended", Priority: 205, Confidence: 0.97},
	{ID: "oai_searchbot", Name: "OAI-SearchBot", Category: "AICrawler", Provider: "OpenAI", Token: "OAI-SearchBot", VersionFrom: "token", Priority: 205, Confidence: 0.98},
	{ID: "perplexity_user", Name: "Perplexity-User", Category: "AIAgent", Provider: "Perplexity", Token: "Perplexity-User", Priority: 205, Confidence: 0.97},
	{ID: "amazonbot", Name: "Amazonbot", Category: "AICrawler", Provider: "Amazon", Token: "Amazonbot", Priority: 200, Confidence: 0.96},
	{ID: "bytespider", Name: "Bytespider", Category: "AICrawler", Provider: "ByteDance", Token: "Bytespider", Priority: 200, Confidence: 0.96},
	{ID: "ccbot", Name: "CCBot", Category: "AICrawler", Provider: "Common Crawl", Token: "CCBot", Priority: 200, Confidence: 0.96},
	{ID: "claudebot", Name: "ClaudeBot", Category: "AICrawler", Provider: "Anthropic", Token: "ClaudeBot", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "gptbot", Name: "GPTBot", Category: "AICrawler", Provider: "OpenAI", Token: "GPTBot", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "meta_external_agent", Name: "meta-externalagent", Category: "AICrawler", Provider: "Meta", Token: "meta-externalagent", Priority: 200, Confidence: 0.97},
	{ID: "perplexitybot", Name: "PerplexityBot", Category: "AICrawler", Provider: "Perplexity", Token: "PerplexityBot", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "googlebot_image", Name: "Googlebot-Image", Category: "SearchEngine", Provider: "Google", Token: "Googlebot-Image", Priority: 195, Confidence: 0.97},
	{ID: "ahrefsbot", Name: "AhrefsBot", Category: "SEO", Provider: "Ahrefs", Token: "AhrefsBot", VersionFrom: "token", Priority: 190, Confidence: 0.97},
	{ID: "applebot", Name: "Applebot", Category: "SearchEngine", Provider: "Apple", Token: "Applebot", Priority: 190, Confidence: 0.96},
	{ID: "baiduspider", Name: "Baiduspider", Category: "SearchEngine", Provider: "Baidu", Token: "Baiduspider", Priority: 190, Confidence: 0.97},
	{ID: "bingbot", Name: "Bingbot", Category: "SearchEngine", Provider: "Microsoft", Token: "bingbot", VersionFrom: "token", Priority: 190, Confidence: 0.98},
	{ID: "censys", Name: "Censys", Category: "Security", Provider: "Censys", Contains: "CensysInspect", Priority: 190, Confidence: 0.95},
	{ID: "discordbot", Name: "Discordbot", Category: "SocialCrawler", Provider: "Discord", Token: "Discordbot", Priority: 190, Confidence: 0.96},
	{ID: "dotbot", Name: "DotBot", Category: "SEO", Provider: "Moz", Token: "DotBot", Priority: 190, Confidence: 0.95},
	{ID: "duckduckbot", Name: "DuckDuckBot", Category: "SearchEngine", Provider: "DuckDuckGo", Token: "DuckDuckBot", Priority: 190, Confidence: 0.97},
	{ID: "facebookexternalhit", Name: "facebookexternalhit", Category: "SocialCrawler", Provider: "Meta", Token: "facebookexternalhit", Priority: 190, Confidence: 0.97},
	{ID: "googlebot", Name: "Googlebot", Category: "SearchEngine", Provider: "Google", Token: "Googlebot", VersionFrom: "token", Priority: 190, Confidence: 0.98},
	{ID: "linkedinbot", Name: "LinkedInBot", Category: "SocialCrawler", Provider: "LinkedIn", Token: "LinkedInBot", Priority: 190, Confidence: 0.97},
	{ID: "mj12bot", Name: "MJ12bot", Category: "SEO", Provider: "Majestic", Token: "MJ12bot", Priority: 190, Confidence: 0.96},
	{ID: "pingdom", Name: "Pingdom", Category: "Monitoring", Provider: "Pingdom", Prefix: "Pingdom.com_bot", Priority: 190, Confidence: 0.96},
	{ID: "pinterestbot", Name: "Pinterestbot", Category: "SocialCrawler", Provider: "Pinterest", Token: "Pinterestbot", Priority: 190, Confidence: 0.96},
	{ID: "semrushbot", Name: "SemrushBot", Category: "SEO", Provider: "Semrush", Token: "SemrushBot", Priority: 190, Confidence: 0.97},
	{ID: "shodan", Name: "Shodan", Category: "Security", Provider: "Shodan", Contains: "Shodan", Priority: 190, Confidence: 0.95},
	{ID: "slackbot", Name: "Slackbot", Category: "SocialCrawler", Provider: "Slack", Prefix: "Slackbot", Priority: 190, Confidence: 0.96},
	{ID: "slurp", Name: "Yahoo! Slurp", Category: "SearchEngine", Provider: "Yahoo", Contains: "Slurp", Priority: 190, Confidence: 0.94},
	{ID: "telegrambot", Name: "TelegramBot", Category: "SocialCrawler", Provider: "Telegram", Token: "TelegramBot", Priority: 190, Confidence: 0.96},
	{ID: "twitterbot", Name: "Twitterbot", Category: "SocialCrawler", Provider: "X", Token: "Twitterbot", VersionFrom: "token", Priority: 190, Confidence: 0.97},
	{ID: "uptimerobot", Name: "UptimeRobot", Category: "Monitoring", Provider: "UptimeRobot", Token: "UptimeRobot", Priority: 190, Confidence: 0.96},
	{ID: "yandexbot", Name: "YandexBot", Category: "SearchEngine", Provider: "Yandex", Token: "YandexBot", VersionFrom: "token", Priority: 190, Confidence: 0.98},
	{ID: "whatsapp", Name: "WhatsApp", Category: "SocialCrawler", Provider: "Meta", Token: "WhatsApp", Priority: 185, Confidence: 0.94},
	{ID: "generic_bot", Name: "Unknown Bot", Category: "Generic", Contains: "bot", Priority: 50, Confidence: 0.55},
	{ID: "generic_crawler", Name: "Unknown Crawler", Category: "Generic", Contains: "crawler", Priority: 50, Confidence: 0.55},
	{ID: "generic_spider", Name: "Unknown Spider", Category: "Generic", Contains: "spider", Priority: 50, Confidence: 0.55},
}

var Apps = []Rule{
	{ID: "instagram_app", Name: "Instagram", Category: "Social", Token: "Instagram", VersionFrom: "suffix", Priority: 210, Confidence: 0.96},
	{ID: "facebook_app", Name: "Facebook", Category: "Social", Contains: "FBAN", Priority: 200, Confidence: 0.96},
	{ID: "facebook_app_android", Name: "Facebook", Category: "Social", Contains: "FB_IAB", Priority: 200, Confidence: 0.96},
	{ID: "line_app", Name: "LINE", Category: "Messaging", Contains: "Line/", Priority: 200, Confidence: 0.95},
	{ID: "linkedin_app", Name: "LinkedIn", Category: "Social", Contains: "LinkedInApp", Priority: 200, Confidence: 0.94},
	{ID: "pinterest_app", Name: "Pinterest", Category: "Social", Contains: "Pinterest/", Priority: 200, Confidence: 0.94},
	{ID: "slack_app", Name: "Slack", Category: "Productivity", Contains: "Slack_SSB", Priority: 200, Confidence: 0.95},
	{ID: "snapchat_app", Name: "Snapchat", Category: "Social", Contains: "Snapchat", Priority: 200, Confidence: 0.95},
	{ID: "telegram_app", Name: "Telegram", Category: "Messaging", Contains: "TelegramWebview", Priority: 200, Confidence: 0.95},
	{ID: "tiktok_app", Name: "TikTok", Category: "Social", Contains: "BytedanceWebview", Priority: 200, Confidence: 0.95},
	{ID: "twitter_app", Name: "X", Category: "Social", Contains: "Twitter for", Priority: 200, Confidence: 0.94},
	{ID: "wechat_app", Name: "WeChat", Category: "Messaging", Contains: "MicroMessenger", Priority: 200, Confidence: 0.96},
}

var Automations = []Rule{
	{ID: "cypress", Name: "Cypress", Category: "Framework", Contains: "Cypress", Priority: 210, Confidence: 0.95},
	{ID: "phantomjs", Name: "PhantomJS", Category: "Headless", Contains: "PhantomJS", Priority: 210, Confidence: 0.97},
	{ID: "playwright", Name: "Playwright", Category: "Framework", Contains: "Playwright", Priority: 210, Confidence: 0.97},
	{ID: "puppeteer", Name: "Puppeteer", Category: "Framework", Contains: "Puppeteer", Priority: 210, Confidence: 0.97},
	{ID: "selenium", Name: "Selenium", Category: "Framework", Contains: "selenium", Priority: 210, Confidence: 0.95},
	{ID: "webdriver", Name: "WebDriver", Category: "Framework", Contains: "WebDriver", Priority: 205, Confidence: 0.94},
	{ID: "headless_chrome", Name: "Headless Chrome", Category: "Headless", Token: "HeadlessChrome", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "headless_firefox", Name: "Headless Firefox", Category: "Headless", Contains: "Headless", Priority: 200, Confidence: 0.85},
	{ID: "electron_automation", Name: "Electron", Category: "Embedded", Token: "Electron", VersionFrom: "token", Priority: 100, Confidence: 0.8},
}

var Clients = []Rule{
	{ID: "aiohttp", Name: "aiohttp", Category: "Library", Token: "aiohttp", VersionFrom: "token", Priority: 200, Confidence: 0.96},
	{ID: "apache_httpclient", Name: "Apache-HttpClient", Category: "Library", Token: "Apache-HttpClient", VersionFrom: "token", Priority: 200, Confidence: 0.97},
	{ID: "axios", Name: "axios", Category: "Library", Token: "axios", VersionFrom: "token", Priority: 200, Confidence: 0.96},
	{ID: "curl", Name: "curl", Category: "CLI", Token: "curl", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "go_http", Name: "Go-http-client", Category: "Library", Token: "Go-http-client", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "guzzle", Name: "Guzzle", Category: "Library", Token: "GuzzleHttp", VersionFrom: "token", Priority: 200, Confidence: 0.96},
	{ID: "httpie", Name: "HTTPie", Category: "CLI", Token: "HTTPie", VersionFrom: "token", Priority: 200, Confidence: 0.97},
	{ID: "insomnia", Name: "Insomnia", Category: "Library", Token: "insomnia", VersionFrom: "token", Priority: 200, Confidence: 0.96},
	{ID: "libwww_perl", Name: "libwww-perl", Category: "Library", Token: "libwww-perl", VersionFrom: "token", Priority: 200, Confidence: 0.96},
	{ID: "node_fetch", Name: "node-fetch", Category: "Library", Token: "node-fetch", VersionFrom: "token", Priority: 200, Confidence: 0.96},
	{ID: "okhttp", Name: "okhttp", Category: "Library", Token: "okhttp", VersionFrom: "token", Priority: 200, Confidence: 0.97},
	{ID: "postman", Name: "PostmanRuntime", Category: "Library", Token: "PostmanRuntime", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "python_httpx", Name: "httpx", Category: "Library", Token: "python-httpx", VersionFrom: "token", Priority: 200, Confidence: 0.97},
	{ID: "python_requests", Name: "python-requests", Category: "Library", Token: "python-requests", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "python_urllib", Name: "urllib", Category: "Library", Token: "Python-urllib", VersionFrom: "token", Priority: 200, Confidence: 0.97},
	{ID: "wget", Name: "Wget", Category: "CLI", Token: "Wget", VersionFrom: "token", Priority: 200, Confidence: 0.98},
	{ID: "xh", Name: "xh", Category: "CLI", Token: "xh", VersionFrom: "token", Priority: 200, Confidence: 0.95},
	{ID: "got", Name: "got", Category: "Library", Token: "got", VersionFrom: "token", Priority: 195, Confidence: 0.92},
	{ID: "java", Name: "Java", Category: "Library", Token: "Java", VersionFrom: "token", Priority: 190, Confidence: 0.94},
}
