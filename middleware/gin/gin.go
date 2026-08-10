// Package gin adapts uax to the Gin framework.
package gin

import (
	"github.com/bakhod1r/uax"
	"github.com/gin-gonic/gin"
)

const contextKey = "uax_client"

// Middleware parses each request with p and stores the client in the Gin
// context. A nil parser uses the package default, which has no cache; a server
// under load should pass its own:
//
//	r.Use(uaxgin.Middleware(uax.New(uax.Config{Cache: uax.CacheConfig{Enabled: true}})))
func Middleware(p *uax.Parser) gin.HandlerFunc {
	if p == nil {
		p = uax.Default()
	}
	return func(c *gin.Context) {
		c.Set(contextKey, p.ParseRequest(c.Request))
		c.Next()
	}
}

// FromContext returns the parsed client. It never returns nil, so handlers do
// not have to guard every access.
func FromContext(c *gin.Context) *uax.Client {
	if v, ok := c.Get(contextKey); ok {
		if client, ok := v.(*uax.Client); ok {
			return client
		}
	}
	return &uax.Client{}
}
