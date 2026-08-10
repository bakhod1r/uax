// Package echo adapts uax to the Echo framework.
package echo

import (
	"github.com/bakhod1r/uax"
	"github.com/labstack/echo/v4"
)

const contextKey = "uax_client"

// Middleware parses each request with p and stores the client in the Echo
// context. A nil parser uses the package default, which has no cache.
func Middleware(p *uax.Parser) echo.MiddlewareFunc {
	if p == nil {
		p = uax.Default()
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(contextKey, p.ParseRequest(c.Request()))
			return next(c)
		}
	}
}

// FromContext returns the parsed client. It never returns nil.
func FromContext(c echo.Context) *uax.Client {
	if client, ok := c.Get(contextKey).(*uax.Client); ok {
		return client
	}
	return &uax.Client{}
}
