// Package middleware attaches a parsed uax.Client to the request context.
//
// Parsing is lazy: the middleware stores the request, and the work happens only
// if a handler actually asks for the client. A route that never calls
// FromContext pays nothing but a context value, so the middleware can be
// mounted globally without auditing which handlers care.
package middleware

import (
	"context"
	"net/http"
	"sync"

	"github.com/bakhod1r/uax"
)

type contextKey struct{}

var clientKey contextKey

// lazyClient defers parsing until the first request for it, then remembers the
// result. sync.Once makes it safe for a handler that fans out to goroutines.
type lazyClient struct {
	once   sync.Once
	parser *uax.Parser
	req    *http.Request
	client *uax.Client
}

func (l *lazyClient) get() *uax.Client {
	l.once.Do(func() { l.client = l.parser.ParseRequest(l.req) })
	return l.client
}

// New returns middleware using the given parser. A nil parser uses the
// package-level default, which has no cache and no shared state.
//
// A server under real load usually wants its own parser with caching enabled:
//
//	p := uax.New(uax.Config{Cache: uax.CacheConfig{Enabled: true}})
//	mux.Use(middleware.New(p))
func New(p *uax.Parser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lazy := &lazyClient{parser: p, req: r}
			if lazy.parser == nil {
				lazy.parser = uax.Default()
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), clientKey, lazy)))
		})
	}
}

// Handler wraps a single handler with the default parser.
func Handler(next http.Handler) http.Handler { return New(nil)(next) }

// FromContext returns the client for this request.
//
// It never returns nil: a handler that reaches for the client always gets one,
// even if the middleware was not mounted, so callers do not have to guard every
// access. A client parsed from an absent middleware simply has empty fields.
func FromContext(ctx context.Context) *uax.Client {
	switch v := ctx.Value(clientKey).(type) {
	case *lazyClient:
		return v.get()
	case *uax.Client:
		return v
	default:
		return &uax.Client{}
	}
}

// WithClient injects a client directly, for tests and for handlers that already
// parsed one.
func WithClient(ctx context.Context, c *uax.Client) context.Context {
	return context.WithValue(ctx, clientKey, c)
}
