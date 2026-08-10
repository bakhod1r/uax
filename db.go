package uax

// CategorySummary is how many rules one category holds.
type CategorySummary struct {
	Category string `json:"category"`
	Rules    int    `json:"rules"`
}

// DatabaseSummary reports the size of the compiled rule database.
//
// It exists so a deployment can record which database produced a given result:
// rules change between releases, and a classification from six months ago is
// only reproducible if you know what it was classified against.
func DatabaseSummary() []CategorySummary {
	groups := builtinGroups()
	return []CategorySummary{
		{"browser", len(groups[catBrowserIdx])},
		{"engine", len(groups[catEngineIdx])},
		{"os", len(groups[catOSIdx])},
		{"bot", len(groups[catBotIdx])},
		{"app", len(groups[catAppIdx])},
		{"automation", len(groups[catAutomationIdx])},
		{"http client", len(groups[catClientIdx])},
	}
}

// DatabaseRules is the total number of compiled rules.
func DatabaseRules() int {
	n := 0
	for _, g := range DatabaseSummary() {
		n += g.Rules
	}
	return n
}
