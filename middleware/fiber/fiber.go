// Package fiber adapts uax to the Fiber framework.
package fiber

import (
	"net/http"

	"github.com/bakhod1r/uax"
	"github.com/gofiber/fiber/v2"
)

const contextKey = "uax_client"

// Middleware parses each request with p and stores the client in Fiber's
// locals. A nil parser uses the package default, which has no cache.
//
// Fiber runs on fasthttp, so the headers are copied into a net/http header set
// before parsing. That copy is the price of reading Client Hints and fetch
// metadata rather than the User-Agent alone, which is worth far more than the
// copy costs.
func Middleware(p *uax.Parser) fiber.Handler {
	if p == nil {
		p = uax.Default()
	}
	return func(c *fiber.Ctx) error {
		req := &http.Request{Header: make(http.Header, 8)}
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		c.Locals(contextKey, p.ParseRequest(req))
		return c.Next()
	}
}

// FromContext returns the parsed client. It never returns nil.
func FromContext(c *fiber.Ctx) *uax.Client {
	if client, ok := c.Locals(contextKey).(*uax.Client); ok {
		return client
	}
	return &uax.Client{}
}
