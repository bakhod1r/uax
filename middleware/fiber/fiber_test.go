package fiber_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bakhod1r/uax"
	uaxfiber "github.com/bakhod1r/uax/middleware/fiber"
	"github.com/bakhod1r/uax/uatest"
	"github.com/gofiber/fiber/v2"
)

var chrome = uatest.Chrome(138, uatest.Windows(10))

// run mounts the middleware, sends one request through it and returns what the
// handler saw. Fiber runs on fasthttp, so app.Test is the only way to exercise
// the header copy the adapter performs — building a fiber.Ctx by hand would
// test the constructor instead.
func run(t *testing.T, p *uax.Parser, ua string, headers ...string) *uax.Client {
	t.Helper()

	var got *uax.Client
	app := fiber.New()
	app.Use(uaxfiber.Middleware(p))
	app.Get("/", func(c *fiber.Ctx) error {
		got = uaxfiber.FromContext(c)
		return c.SendStatus(204)
	})

	req := uatest.Request(ua, headers...)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 204 {
		t.Errorf("status: got %d, want 204 — the chain did not reach the handler", resp.StatusCode)
	}
	if got == nil {
		t.Fatal("FromContext returned nil")
	}
	return got
}

func TestMiddleware(t *testing.T) {
	got := run(t, nil, chrome)

	if got.Browser.Name != "Chrome" {
		t.Errorf("browser: got %q, want Chrome", got.Browser.Name)
	}
	if got.OS.Name != "Windows" {
		t.Errorf("os: got %q, want Windows", got.OS.Name)
	}
}

func TestMiddlewareCopiesTheHintHeaders(t *testing.T) {
	// The adapter copies fasthttp's headers into a net/http set by hand. A copy
	// that dropped everything but the User-Agent would still report Chrome, so
	// the assertion is on a source only the hints can produce.
	got := run(t, nil, chrome, uatest.WithHints("Google Chrome", 138, "Windows", false)...)

	if got.Browser.Source != uax.SourceHints && got.Browser.Source != uax.SourceBoth {
		t.Errorf("browser source: got %q — the hint headers were not copied", got.Browser.Source)
	}
}

func TestMiddlewareUsesTheParserItIsGiven(t *testing.T) {
	// A configured parser is why the argument exists; ignoring it would quietly
	// drop the caller's cache, metrics and device hook.
	got := run(t, uax.New(uax.Config{SkipBot: true}), uatest.HeadlessChrome(138, uatest.Linux()))

	// SkipBot covers automation too, and headless Chrome is otherwise resolved.
	if got.Automation != nil {
		t.Errorf("SkipBot was configured but automation still resolved as %q", got.Automation.Name)
	}
}

func TestFromContextWithoutMiddleware(t *testing.T) {
	// A handler reading the client on a route where the middleware was never
	// mounted gets an empty client, not a nil dereference in production.
	var got *uax.Client

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		got = uaxfiber.FromContext(c)
		return c.SendStatus(204)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if got == nil {
		t.Fatal("FromContext returned nil without the middleware mounted")
	}
}
