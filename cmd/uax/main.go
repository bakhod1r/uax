// Command uax inspects User-Agent strings from the command line.
//
//	uax parse    "<user-agent>"   one-line summary
//	uax explain  "<user-agent>"   full breakdown, including why
//	uax json     "<user-agent>"   the parsed Client as JSON
//	uax analyze  <file>           aggregate a log file of User-Agents
//	uax db                        rule database summary
//
// With no subcommand, a single argument is explained:
//
//	uax "Mozilla/5.0 (…)"
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/bakhod1r/devicex"
	"github.com/bakhod1r/uax"
)

var osExit = os.Exit

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "uax:", err)
		osExit(1)
	}
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		usage(out)
		return nil
	}

	cmd, rest := args[0], args[1:]

	// A bare User-Agent is the common case; do the obvious thing with it.
	if strings.HasPrefix(cmd, "Mozilla/") || strings.Contains(cmd, "/") && len(rest) == 0 && !isCommand(cmd) {
		return explain(cmd, out)
	}

	switch cmd {
	case "parse":
		if len(rest) == 0 {
			return fmt.Errorf("parse needs a User-Agent")
		}
		return summarize(rest[0], out)

	case "explain":
		if len(rest) == 0 {
			return fmt.Errorf("explain needs a User-Agent")
		}
		return explain(rest[0], out)

	case "json":
		if len(rest) == 0 {
			return fmt.Errorf("json needs a User-Agent")
		}
		b, _ := json.MarshalIndent(parser.Parse(rest[0]), "", "  ")
		_, err := fmt.Fprintln(out, string(b))
		return err

	case "analyze":
		if len(rest) == 0 {
			return fmt.Errorf("analyze needs a file, or - for stdin")
		}
		return analyze(rest[0], out)

	case "device":
		if len(rest) == 0 {
			return fmt.Errorf("device needs a model code, e.g. uax device SM-S928B")
		}
		return describeModel(rest[0], out)

	case "db":
		return describeDB(out)

	case "robots":
		return robotsCmd(rest, out)

	case "version", "--version", "-v":
		// The database has its own identity: a verdict is only reproducible
		// against the rules that produced it.
		_, _ = fmt.Fprintf(out, "uax %s\nrules %d (database %s)\n",
			uax.Release, uax.DatabaseRules(), uax.DatabaseVersion())
		return nil

	case "help", "-h", "--help":
		usage(out)
		return nil

	default:
		return fmt.Errorf("unknown command %q (try: uax help)", cmd)
	}
}

func isCommand(s string) bool {
	switch s {
	case "parse", "explain", "json", "analyze", "device", "db", "robots", "version", "help":
		return true
	}
	return false
}

func usage(out io.Writer) {
	fmt.Fprint(out, `uax — client intelligence for User-Agent strings

  uax parse   "<user-agent>"   one-line summary
  uax explain "<user-agent>"   full breakdown
  uax json    "<user-agent>"   parsed Client as JSON
  uax analyze <file|->         aggregate a log of User-Agents
  uax device  <model-code>     resolve a bare model code
  uax db                       rule database summary
  uax robots [flags]           robots.txt from the bot database
  uax version                  library version and rule database fingerprint

robots flags:
  --deny <categories>          ai, ai-crawler, ai-agent, search, social,
                               seo, monitoring, security (comma-separated).
                               Defaults to ai when nothing is given.
  --allow <categories>         permit these and deny everything else
  --deny-provider <names>      block by operator, e.g. OpenAI,Anthropic
  --path <prefix>              restrict the rules to a path (default /)
  --sitemap <url>              append a Sitemap line
`)
}

func summarize(ua string, out io.Writer) error {
	c := parser.Parse(ua)

	parts := []string{}
	if c.Browser.Name != "" {
		parts = append(parts, strings.TrimSpace(c.Browser.Name+" "+c.Browser.Version.Raw))
	}
	if c.OS.Name != "" {
		parts = append(parts, strings.TrimSpace(c.OS.Name+" "+c.OS.Version.Raw))
	}
	if d := deviceLabel(c); d != "" {
		parts = append(parts, d)
	}
	if c.Bot != nil {
		parts = append(parts, c.Bot.Name+" ("+string(c.Bot.Category)+")")
	}
	if c.HTTPClient != nil {
		parts = append(parts, strings.TrimSpace(c.HTTPClient.Name+" "+c.HTTPClient.Version.Raw))
	}
	if len(parts) == 0 {
		parts = append(parts, "unrecognised")
	}

	_, _ = fmt.Fprintf(out, "%s [%s]\n", strings.Join(parts, " · "), c.Kind)
	return nil
}

func deviceLabel(c *uax.Client) string {
	switch {
	case c.Device.Name != "":
		return c.Device.Brand + " " + c.Device.Name
	case c.Device.Model != "":
		return strings.TrimSpace(c.Device.Brand + " " + c.Device.Model)
	case c.Device.Brand != "":
		return c.Device.Brand
	default:
		return string(c.Device.Type)
	}
}

func explain(ua string, out io.Writer) error {
	c := parser.Parse(ua)
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)

	row := func(k, v string, conf float32) {
		if v == "" {
			return
		}
		if conf > 0 {
			fmt.Fprintf(w, "%s\t%s\t%.2f\n", k, v, conf)
		} else {
			fmt.Fprintf(w, "%s\t%s\t\n", k, v)
		}
	}

	fmt.Fprintf(w, "FIELD\tVALUE\tCONFIDENCE\n")
	row("kind", string(c.Kind), 0)
	row("browser", strings.TrimSpace(c.Browser.Name+" "+c.Browser.Version.Raw), c.Browser.Confidence)
	row("  family", c.Browser.Family, 0)
	row("  mode", c.Browser.Mode, 0)
	row("  source", string(c.Browser.Source), 0)
	row("engine", strings.TrimSpace(c.Engine.Name+" "+c.Engine.Version.Raw), c.Engine.Confidence)
	row("os", strings.TrimSpace(c.OS.Name+" "+c.OS.Version.Raw), c.OS.Confidence)
	row("  arch", c.OS.Arch, 0)
	row("device", string(c.Device.Type), c.Device.Confidence)
	row("  brand", c.Device.Brand, 0)
	row("  model", c.Device.Model, 0)
	row("  name", c.Device.Name, 0)
	row("  family", c.Device.Family, 0)

	if c.Bot != nil {
		row("bot", c.Bot.Name, c.Bot.Confidence)
		row("  category", string(c.Bot.Category), 0)
		row("  provider", c.Bot.Provider, 0)
	}
	if c.AI != nil {
		row("ai", c.AI.Provider+" "+c.AI.Name, 0)
	}
	if c.App != nil {
		row("app", c.App.Name, c.App.Confidence)
	}
	if c.Automation != nil {
		row("automation", c.Automation.Name, c.Automation.Confidence)
	}
	if c.HTTPClient != nil {
		row("http client", strings.TrimSpace(c.HTTPClient.Name+" "+c.HTTPClient.Version.Raw), c.HTTPClient.Confidence)
	}
	_ = w.Flush()

	// The verdict and its reasons are the part a developer actually debugs.
	// One coherence check serves both it and the listing below.
	cons := c.Consistency()
	v := uax.VerdictFrom(c, cons)
	fmt.Fprintf(out, "\nverdict: %s  trust=%.2f  risk=%.2f\n", v.Recommend, v.Trust, v.Risk)
	for _, r := range v.Reasons {
		fmt.Fprintf(out, "  - %s", r.Code)
		if r.Detail != "" {
			fmt.Fprintf(out, ": %s", r.Detail)
		}
		fmt.Fprintln(out)
	}

	if !cons.OK() {
		fmt.Fprintf(out, "\ninconsistencies:\n")
		for _, f := range cons.Findings {
			fmt.Fprintf(out, "  - [%s] %s: %s vs %s\n", f.Severity, f.Rule, f.Left, f.Right)
		}
	}
	return nil
}

// analyze aggregates a file of User-Agents, one per line. It is the offline
// counterpart to the middleware: point it at an access log to see what actually
// reaches a service, and how much of it the database does not recognise.
func analyze(path string, out io.Writer) error {
	f := io.ReadCloser(os.Stdin)
	if path != "-" {
		var err error
		if f, err = os.Open(path); err != nil {
			return err
		}
		// Read-only: nothing was buffered, so a close error reports nothing the
		// caller could act on.
		defer func() { _ = f.Close() }()
	}

	var (
		total    int
		unknown  int
		browsers = map[string]int{}
		oses     = map[string]int{}
		kinds    = map[string]int{}
		bots     = map[string]int{}
	)

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)

	var c uax.Client
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		total++

		parser.ParseTo(&c, line)
		kinds[string(c.Kind)]++

		switch {
		case c.Bot != nil:
			bots[c.Bot.Name]++
		case c.Browser.Name != "":
			browsers[c.Browser.Name]++
		case c.HTTPClient != nil:
			browsers[c.HTTPClient.Name]++
		default:
			unknown++
		}
		if c.OS.Name != "" {
			oses[c.OS.Name]++
		}
	}
	_ = sc.Err()

	fmt.Fprintf(out, "%d User-Agents, %d unrecognised (%.1f%%)\n", total, unknown, pct(unknown, total))
	printTop(out, "clients", browsers, total)
	printTop(out, "operating systems", oses, total)
	printTop(out, "bots", bots, total)
	printTop(out, "kinds", kinds, total)
	return nil
}

func printTop(out io.Writer, title string, counts map[string]int, total int) {
	if len(counts) == 0 {
		return
	}

	type kv struct {
		k string
		v int
	}
	rows := make([]kv, 0, len(counts))
	for k, v := range counts {
		rows = append(rows, kv{k, v})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].v != rows[j].v {
			return rows[i].v > rows[j].v
		}
		return rows[i].k < rows[j].k
	})
	if len(rows) > 10 {
		rows = rows[:10]
	}

	fmt.Fprintf(out, "\n%s:\n", title)
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, r := range rows {
		fmt.Fprintf(w, "  %s\t%d\t%.1f%%\n", r.k, r.v, pct(r.v, total))
	}
	_ = w.Flush()
}

// parser is the command's parser. It is wired to devicex, which is what turns
// a model code into a handset: uax extracts the code, devicex names it.
var parser = uax.New(uax.Config{Device: devicex.Describe})

func pct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) * 100 / float64(total)
}

func describeModel(code string, out io.Writer) error {
	d, ok := devicex.Lookup(code)
	if !ok {
		brand, bok := devicex.BrandOf(code)
		if !bok {
			fmt.Fprintf(out, "%s: unknown\n", code)
			return nil
		}

		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "model\t%s\n", code)
		fmt.Fprintf(w, "brand\t%s\n", brand)
		fmt.Fprintf(w, "name\t(not in the catalogue)\n")
		_ = w.Flush()
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "model\t%s\n", d.Code)
	fmt.Fprintf(w, "brand\t%s\n", d.Brand)

	if d.Name != "" {
		fmt.Fprintf(w, "name\t%s\n", d.Name)
	}
	_ = w.Flush()
	return nil
}

func describeDB(out io.Writer) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CATEGORY\tRULES")
	for _, g := range uax.DatabaseSummary() {
		fmt.Fprintf(w, "%s\t%d\n", g.Category, g.Rules)
	}
	_ = w.Flush()

	named := devicex.Len()
	fmt.Fprintf(out, "\n%d model codes resolve to a marketing name.\n", named)
	return nil
}
