package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bakhod1r/uax"
	"github.com/bakhod1r/uax/middleware"
	"github.com/bakhod1r/uax/uatest"
)

func TestWithClient(t *testing.T) {
	c := &uax.Client{}
	ctx := middleware.WithClient(context.Background(), c)
	if middleware.FromContext(ctx) != c {
		t.Error("middleware mismatch")
	}
}

func TestClientReachesTheHandler(t *testing.T) {
	t.Parallel()

	var got *uax.Client
	h := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = middleware.FromContext(r.Context())
	}))

	h.ServeHTTP(httptest.NewRecorder(), uatest.Request(uatest.Chrome(138, uatest.Windows(10))))

	if got == nil {
		t.Fatal("FromContext returned nil")
	}
	if got.Browser.Name != "Chrome" {
		t.Errorf("browser: got %q, want Chrome", got.Browser.Name)
	}
}

// Parsing must not happen for handlers that never ask for the client. This is
// what makes the middleware safe to mount globally.
func TestParsingIsLazy(t *testing.T) {
	t.Parallel()

	m := &countingMetrics{}
	p := uax.New(uax.Config{Metrics: m})

	silent := middleware.New(p)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	silent.ServeHTTP(httptest.NewRecorder(), uatest.Request(uatest.Chrome(138, uatest.Windows(10))))

	if m.parses != 0 {
		t.Errorf("a handler that ignored the client still caused %d parses", m.parses)
	}

	curious := middleware.New(p)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = middleware.FromContext(r.Context()).Browser.Name
	}))
	curious.ServeHTTP(httptest.NewRecorder(), uatest.Request(uatest.Chrome(138, uatest.Windows(10))))

	if m.parses != 1 {
		t.Errorf("parses: got %d, want 1", m.parses)
	}
}

// Repeated access within one request must parse once.
func TestClientIsParsedOncePerRequest(t *testing.T) {
	t.Parallel()

	m := &countingMetrics{}
	p := uax.New(uax.Config{Metrics: m})

	h := middleware.New(p)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 5; i++ {
			_ = middleware.FromContext(r.Context())
		}
	}))
	h.ServeHTTP(httptest.NewRecorder(), uatest.Request(uatest.Chrome(138, uatest.Windows(10))))

	if m.parses != 1 {
		t.Errorf("parses: got %d, want 1", m.parses)
	}
}

// A handler must never have to nil-check the client.
func TestFromContextWithoutMiddleware(t *testing.T) {
	t.Parallel()

	c := middleware.FromContext(httptest.NewRequest(http.MethodGet, "/", nil).Context())
	if c == nil {
		t.Fatal("FromContext returned nil without the middleware mounted")
	}
	if c.Browser.Name != "" {
		t.Errorf("expected an empty client, got %+v", c.Browser)
	}
}

type countingMetrics struct {
	uax.NoopReporter
	parses int
}

func (m *countingMetrics) IncParse() { m.parses++ }
