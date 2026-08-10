package gin_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bakhod1r/uax"
	uaxgin "github.com/bakhod1r/uax/middleware/gin"
	"github.com/bakhod1r/uax/uatest"
	"github.com/gin-gonic/gin"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var got *uax.Client
	r := gin.New()
	r.Use(uaxgin.Middleware(nil))
	r.GET("/", func(c *gin.Context) { got = uaxgin.FromContext(c) })

	r.ServeHTTP(httptest.NewRecorder(), uatest.Request(uatest.Chrome(138, uatest.Windows(10))))

	if got == nil {
		t.Fatal("FromContext returned nil")
	}
	if got.Browser.Name != "Chrome" {
		t.Errorf("browser: got %q, want Chrome", got.Browser.Name)
	}
}

func TestFromContextWithoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if uaxgin.FromContext(c) == nil {
		t.Error("FromContext returned nil without the middleware mounted")
	}
}
