package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bakhod1r/devicex"
	"github.com/bakhod1r/uax"
)

// The tool surface.
//
// Four tools rather than one general-purpose one. A model picks a tool from its
// description, so a tool that does several things is a tool it will reach for
// at the wrong times — and each of these answers a genuinely different
// question, from different evidence.

func tools() []tool {
	return []tool{
		{
			Name: "parse_user_agent",
			Description: "Identify the client behind a User-Agent string: browser, engine, operating " +
				"system, device, bot, in-app WebView, automation framework or HTTP library. Every " +
				"resolved field carries a confidence and the evidence it came from. Use this instead " +
				"of reading a User-Agent yourself — the strings are deliberately misleading (every " +
				"browser claims 'Mozilla/5.0'), and unresolved fields come back empty rather than " +
				"guessed. Values frozen by Chrome's User-Agent Reduction are flagged as frozen rather " +
				"than reported as observed facts.",
			InputSchema: object(
				props{"user_agent": str("The raw User-Agent header value.")},
				"user_agent",
			),
		},
		{
			Name: "check_consistency",
			Description: "Check whether a client's claims contradict each other, and get a trust/risk " +
				"verdict with the reasoning behind it. Detects impossibilities such as Safari on " +
				"Windows or a non-WebKit engine on iOS. Use when deciding whether to believe a " +
				"client. Findings are evidence of an edited string, not proof of intent, and the " +
				"recommended action is advice — never an authorisation decision.",
			InputSchema: object(
				props{"user_agent": str("The raw User-Agent header value.")},
				"user_agent",
			),
		},
		{
			Name: "lookup_device",
			Description: "Resolve an Android model code such as SM-S928B or CPH2451 to its brand and " +
				"marketing name. Use whenever a model code appears and the handset matters. A code " +
				"the catalogue does not hold resolves its manufacturer from the code's shape and " +
				"reports no name — it is never approximated, because a wrong device name is worse " +
				"than no device name.",
			InputSchema: object(
				props{"code": str("The model code, exactly as it appeared (for example \"SM-S928B\").")},
				"code",
			),
		},
		{
			Name: "database_info",
			Description: "Report what this server can recognise: rule counts per category, the rule " +
				"database fingerprint, and the device catalogue's size and snapshot date. Use to " +
				"tell whether a negative answer means 'not a real client' or 'newer than the data', " +
				"and to record which database produced a classification.",
			InputSchema: object(nil),
		},
	}
}

func call(p toolCallParams) (*toolResult, error) {
	switch p.Name {
	case "parse_user_agent":
		ua, err := stringArg(p.Arguments, "user_agent")
		if err != nil {
			return nil, err
		}
		return parseUserAgent(ua), nil

	case "check_consistency":
		ua, err := stringArg(p.Arguments, "user_agent")
		if err != nil {
			return nil, err
		}
		return checkConsistency(ua), nil

	case "lookup_device":
		code, err := stringArg(p.Arguments, "code")
		if err != nil {
			return nil, err
		}
		return lookupDevice(code), nil

	case "database_info":
		return databaseInfo(), nil

	default:
		return nil, errors.New("unknown tool: " + p.Name)
	}
}

// ------------------------------------------------------------------ handlers

func parseUserAgent(ua string) *toolResult {
	c := parser.Parse(ua)

	var b strings.Builder
	fmt.Fprintf(&b, "kind: %s\n", c.Kind)

	writeField(&b, "browser", c.Browser.Name, c.Browser.Version.String(), float64(c.Browser.Confidence), string(c.Browser.Source))
	if c.Browser.Family != "" && c.Browser.Family != c.Browser.Name {
		fmt.Fprintf(&b, "  family: %s\n", c.Browser.Family)
	}
	if c.Browser.Mode != "" {
		fmt.Fprintf(&b, "  mode: %s\n", c.Browser.Mode)
	}

	writeField(&b, "engine", c.Engine.Name, c.Engine.Version.String(), float64(c.Engine.Confidence), "")
	writeField(&b, "os", c.OS.Name, c.OS.Version.String(), float64(c.OS.Confidence), string(c.OS.Source))
	if c.OS.Arch != "" {
		fmt.Fprintf(&b, "  arch: %s\n", c.OS.Arch)
	}
	if c.OS.Reduced {
		b.WriteString("  note: the platform version is frozen by the browser, not observed\n")
	}

	if c.Device.Type != "" || c.Device.Model != "" || c.Device.Name != "" {
		fmt.Fprintf(&b, "device: %s\n", firstNonEmpty(string(c.Device.Type), "unknown type"))
		writeIf(&b, "  brand", c.Device.Brand)
		writeIf(&b, "  model", c.Device.Model)
		writeIf(&b, "  name", c.Device.Name)
		writeIf(&b, "  family", c.Device.Family)
		if c.Device.Reduced {
			b.WriteString("  note: the model was replaced by a placeholder before this request was sent\n")
		}
	}

	if c.Bot != nil {
		fmt.Fprintf(&b, "bot: %s (%s)\n", c.Bot.Name, c.Bot.Category)
		writeIf(&b, "  provider", c.Bot.Provider)
		b.WriteString("  note: this is what the client called itself, not proof of who sent the request\n")
	}
	if c.AI != nil {
		fmt.Fprintf(&b, "ai: %s (%s)\n", c.AI.Name, c.AI.Category)
	}
	if c.App != nil {
		fmt.Fprintf(&b, "in-app webview: %s\n", c.App.Name)
	}
	if c.Automation != nil {
		fmt.Fprintf(&b, "automation: %s (%s)\n", c.Automation.Name, c.Automation.Category)
	}
	if c.HTTPClient != nil {
		fmt.Fprintf(&b, "http client: %s (%s)\n", c.HTTPClient.Name, c.HTTPClient.Category)
	}
	if c.Truncated {
		b.WriteString("note: the input exceeded the length bound and was truncated before parsing\n")
	}

	// A string that resolved nothing is a real answer, and saying so plainly
	// beats returning a page of empty fields the model has to interpret.
	if c.Browser.Name == "" && c.OS.Name == "" && c.Bot == nil && c.HTTPClient == nil &&
		c.Automation == nil && c.App == nil {
		return errorResult("Nothing in that string identified a client. It is not a User-Agent this " +
			"database recognises — which may mean it is malformed, or newer than the rules. " +
			"Call database_info to see what the database covers.")
	}

	return textResult(strings.TrimRight(b.String(), "\n"))
}

func checkConsistency(ua string) *toolResult {
	c := parser.Parse(ua)
	cons := c.Consistency()

	// Parse leaves Verdict nil: a verdict is a pointer, and ParseTo allocates
	// nothing, so only ParseRequest fills one in. This tool promises a verdict,
	// and a missing one would read as "no risk" rather than "not computed".
	//
	// The verdict is held locally rather than written back into c, which the
	// parser may have handed out of a cache and may hand out again.
	verdict := uax.VerdictFrom(c, cons)

	var b strings.Builder
	if cons.OK() {
		b.WriteString("No contradictions found. The client's claims hold together.\n")
		b.WriteString("This is not proof of honesty — a carefully edited string is also consistent.\n")
	} else {
		fmt.Fprintf(&b, "%d contradiction(s) found:\n", len(cons.Findings))
		for _, f := range cons.Findings {
			fmt.Fprintf(&b, "  [%s] %s\n", f.Severity, f.Rule)
			writeIf(&b, "    detail", f.Detail)
			if f.Left != "" || f.Right != "" {
				fmt.Fprintf(&b, "    conflict: %q vs %q\n", f.Left, f.Right)
			}
		}
	}

	if v := verdict; v != nil {
		fmt.Fprintf(&b, "\ntrust: %.2f  risk: %.2f  recommend: %s\n", v.Trust, v.Risk, v.Recommend)
		for _, r := range v.Reasons {
			fmt.Fprintf(&b, "  %s: %s\n", r.Code, r.Detail)
		}
		b.WriteString("\nThe recommendation is advice. Being a bot is not being a liar: a well-behaved " +
			"crawler that names itself is highly trusted and low risk.\n")
	}

	return textResult(strings.TrimRight(b.String(), "\n"))
}

func lookupDevice(code string) *toolResult {
	// The code may arrive with the build suffix Android appends.
	code = devicex.Code(code)

	if d, ok := devicex.Lookup(code); ok {
		return textResult(fmt.Sprintf("model: %s\nbrand: %s\nname: %s\nsource: %s",
			d.Code, d.Brand, d.Name, devicex.Source))
	}

	if brand, ok := devicex.BrandOf(code); ok {
		return textResult(fmt.Sprintf("model: %s\nbrand: %s\nname: (not in the catalogue)\n\n"+
			"The brand comes from the shape of the code, which is why it resolves for handsets the "+
			"catalogue has never seen. The marketing name needs a catalogue entry and there is none, "+
			"so it is reported as absent rather than guessed. The catalogue's snapshot is %s.",
			code, brand, devicex.Generated))
	}

	return errorResult(fmt.Sprintf("%q matches no catalogue entry and no known code shape. It may not "+
		"be an Android model code. The catalogue holds %d devices as of %s.",
		code, devicex.Len(), devicex.Generated))
}

func databaseInfo() *toolResult {
	var b strings.Builder
	b.WriteString("Rules recognised by category:\n")
	for _, g := range uax.DatabaseSummary() {
		fmt.Fprintf(&b, "  %-12s %d\n", g.Category, g.Rules)
	}
	fmt.Fprintf(&b, "  %-12s %d\n", "total", uax.DatabaseRules())
	fmt.Fprintf(&b, "\nuax release:       %s\n", uax.Release)
	fmt.Fprintf(&b, "rule fingerprint:  %s\n", uax.DatabaseVersion())
	fmt.Fprintf(&b, "device catalogue:  %d devices, snapshot %s\n", devicex.Len(), devicex.Generated)
	fmt.Fprintf(&b, "catalogue source:  %s\n", devicex.Source)
	b.WriteString("\nA client newer than these dates resolves what its shape reveals and no more. " +
		"Record the fingerprint alongside any stored classification: rules change between releases, " +
		"and a result is only reproducible against the database that produced it.\n")
	return textResult(strings.TrimRight(b.String(), "\n"))
}

// -------------------------------------------------------------------- helpers

func stringArg(raw json.RawMessage, name string) (string, error) {
	var args map[string]any
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &args); err != nil {
			return "", errors.New("arguments must be a JSON object")
		}
	}
	v, ok := args[name]
	if !ok {
		return "", errors.New("missing required argument: " + name)
	}
	s, ok := v.(string)
	if !ok {
		return "", errors.New(name + " must be a string")
	}
	if s == "" {
		return "", errors.New(name + " must not be empty")
	}
	return s, nil
}

type props map[string]map[string]any

func str(description string) map[string]any {
	return map[string]any{"type": "string", "description": description}
}

func object(p props, required ...string) map[string]any {
	properties := map[string]any{}
	for k, v := range p {
		properties[k] = v
	}
	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func writeField(b *strings.Builder, label, name, version string, confidence float64, source string) {
	if name == "" {
		return
	}
	fmt.Fprintf(b, "%s: %s", label, name)
	if version != "" {
		fmt.Fprintf(b, " %s", version)
	}
	if confidence > 0 {
		fmt.Fprintf(b, "  (confidence %.2f", confidence)
		if source != "" {
			fmt.Fprintf(b, ", from %s", source)
		}
		b.WriteString(")")
	}
	b.WriteString("\n")
}

func writeIf(b *strings.Builder, label, value string) {
	if value != "" {
		fmt.Fprintf(b, "%s: %s\n", label, value)
	}
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
