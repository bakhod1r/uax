package uax

// MetricsReporter receives operational counters. Supply one via Config.Metrics
// to feed Prometheus, OpenTelemetry, or anything else.
//
// The counters describe work the parser actually performed: a cache hit
// increments IncCacheHit and nothing else, because no parsing happened. That
// makes the ratio between IncParse and IncCacheHit meaningful, which is the
// number worth graphing.
type MetricsReporter interface {
	// IncParse counts inputs actually parsed, excluding cache hits.
	IncParse()

	// IncCacheHit counts results served from the cache.
	IncCacheHit()

	// IncBotDetected counts clients that identified themselves as a bot,
	// crawler or AI agent.
	IncBotDetected()

	// IncSpoofDetected counts clients whose claims contradict each other.
	IncSpoofDetected()
}

// NoopReporter discards every counter. It is useful as an embedded field when
// implementing only some of the interface.
type NoopReporter struct{}

// The bodies are empty because discarding is the whole behaviour.
func (NoopReporter) IncParse()         {}
func (NoopReporter) IncCacheHit()      {}
func (NoopReporter) IncBotDetected()   {}
func (NoopReporter) IncSpoofDetected() {}
