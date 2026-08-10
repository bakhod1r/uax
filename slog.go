package uax

import "log/slog"

// Attrs returns the client as structured log attributes.
//
// The keys are stable and namespaced, so dashboards and alerts built on them
// keep working across releases. Empty fields are omitted rather than logged as
// "": absence should not occupy a column.
func (c *Client) Attrs() []slog.Attr {
	if c == nil {
		return nil
	}

	attrs := make([]slog.Attr, 0, 12)
	add := func(k, v string) {
		if v != "" {
			attrs = append(attrs, slog.String(k, v))
		}
	}

	add("client.kind", string(c.Kind))
	add("browser.name", c.Browser.Name)
	add("browser.version", c.Browser.Version.Raw)
	add("browser.mode", c.Browser.Mode)
	add("os.name", c.OS.Name)
	add("os.version", c.OS.Version.Raw)
	add("device.type", string(c.Device.Type))
	add("device.brand", c.Device.Brand)
	add("device.model", c.Device.Model)
	add("device.name", c.Device.Name)

	if c.Bot != nil {
		add("bot.name", c.Bot.Name)
		add("bot.category", string(c.Bot.Category))
		add("bot.provider", c.Bot.Provider)
	}
	if c.HTTPClient != nil {
		add("http_client.name", c.HTTPClient.Name)
	}
	if c.Automation != nil {
		add("automation.name", c.Automation.Name)
	}
	if c.Verdict != nil {
		attrs = append(attrs,
			slog.Float64("verdict.trust", float64(c.Verdict.Trust)),
			slog.Float64("verdict.risk", float64(c.Verdict.Risk)),
		)
		add("verdict.recommend", string(c.Verdict.Recommend))
	}
	return attrs
}

// LogValue makes Client usable directly as a slog value:
//
//	logger.Info("request", "client", client)
func (c *Client) LogValue() slog.Value {
	return slog.GroupValue(c.Attrs()...)
}
