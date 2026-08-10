package echo_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bakhod1r/uax"
	uaxecho "github.com/bakhod1r/uax/middleware/echo"
	"github.com/bakhod1r/uax/uatest"
	"github.com/labstack/echo/v4"
)

func TestMiddleware(t *testing.T) {
	var got *uax.Client

	e := echo.New()
	e.Use(uaxecho.Middleware(nil))
	e.GET("/", func(c echo.Context) error {
		got = uaxecho.FromContext(c)
		return c.NoContent(200)
	})

	e.ServeHTTP(httptest.NewRecorder(), uatest.Request(uatest.Chrome(138, uatest.Windows(10))))

	if got == nil {
		t.Fatal("FromContext returned nil")
	}
	if got.Browser.Name != "Chrome" {
		t.Errorf("browser: got %q, want Chrome", got.Browser.Name)
	}
	if got.OS.Name != "Windows" {
		t.Errorf("os: got %q, want Windows", got.OS.Name)
	}
}

func TestMiddlewareReadsClientHints(t *testing.T) {
	// The adapter exists to hand the whole request to the parser, not just the
	// User-Agent. If it passed the string alone this would still say Chrome —
	// so the assertion is on a field only the hints can set.
	var got *uax.Client

	e := echo.New()
	e.Use(uaxecho.Middleware(nil))
	e.GET("/", func(c echo.Context) error {
		got = uaxecho.FromContext(c)
		return c.NoContent(200)
	})

	req := uatest.Request(
		uatest.Chrome(138, uatest.Windows(10)),
		uatest.WithHints("Google Chrome", 138, "Windows", false)...,
	)
	e.ServeHTTP(httptest.NewRecorder(), req)

	// The hints agree with the string here, which the parser records as "both".
	// The point is only that a source naming the hints appeared at all: parsing
	// the User-Agent alone can never produce one.
	if got.Browser.Source != uax.SourceHints && got.Browser.Source != uax.SourceBoth {
		t.Errorf("browser source: got %q — the hints were not read", got.Browser.Source)
	}
}

func TestFromContextWithoutMiddleware(t *testing.T) {
	// A handler that reads the client on a route where the middleware was never
	// mounted gets an empty client, not a nil dereference in production.
	e := echo.New()
	c := e.NewContext(httptest.NewRequest("GET", "/", nil), httptest.NewRecorder())

	if uaxecho.FromContext(c) == nil {
		t.Fatal("FromContext returned nil without the middleware mounted")
	}
}

func TestMiddlewareUsesTheParserItIsGiven(t *testing.T) {
	// A configured parser is the whole reason the argument exists; ignoring it
	// would silently drop the caller's cache, metrics and device hook.
	var got *uax.Client
	p := uax.New(uax.Config{SkipBot: true})

	e := echo.New()
	e.Use(uaxecho.Middleware(p))
	e.GET("/", func(c echo.Context) error {
		got = uaxecho.FromContext(c)
		return c.NoContent(200)
	})

	e.ServeHTTP(httptest.NewRecorder(), uatest.Request(uatest.HeadlessChrome(138, uatest.Linux())))

	// SkipBot covers automation too, and headless Chrome is normally resolved.
	if got.Automation != nil {
		t.Errorf("SkipBot was configured but automation still resolved as %q", got.Automation.Name)
	}
}

func TestHandlersDownstreamStillRun(t *testing.T) {
	// Middleware that swallows the chain breaks every route it is mounted on.
	e := echo.New()
	e.Use(uaxecho.Middleware(nil))

	reached := false
	e.GET("/", func(c echo.Context) error {
		reached = true
		return c.NoContent(204)
	})

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, uatest.Request(uatest.Chrome(138, uatest.Windows(10))))

	if !reached {
		t.Error("the handler never ran")
	}
	if rec.Code != 204 {
		t.Errorf("status: got %d, want 204", rec.Code)
	}
}
