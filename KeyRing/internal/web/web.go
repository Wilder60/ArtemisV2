package web

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Wilder60/KeyRing/internal/security"

	"github.com/gin-gonic/gin"
)

//Note more functions can be added here if metadata
//will be added

func RegisterRoutes(kr *KeyRing, sec *security.Security) *gin.Engine {
	r := gin.Default()
	base := r.Group("/")

	base.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	ring := r.Group("/api/v1")
	ring.Use(sec.Authorize())
	ring.GET("/keyring", kr.getEvents)
	ring.POST("/keyring", kr.addEvent)
	return r
}

var RouterModule = fx.Option(
	fx.Provide(RegisterRoutes),
)
