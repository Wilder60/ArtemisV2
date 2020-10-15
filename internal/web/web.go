package web

import (
	"fmt"
	"net/http"

	"go.uber.org/fx"

	"github.com/Wilder60/KeyRing/internal/security"

	"github.com/gin-gonic/gin"
)

//Note more functions can be added here if metadata
//will be added

func RegisterRoutes(kr *KeyRing) *gin.Engine {
	r := gin.Default()
	base := r.Group("/")
	base.GET("/health", health)
	base.GET("/key", generateKey)

	ring := r.Group("/api/v1")
	ring.Use(security.Authorize())
	ring.GET("/keyring", kr.getEvents)
	ring.POST("/keyring", kr.addEvent)
	return r
}

// These should be moved to a seperate class
func health(c *gin.Context) {
	c.Status(http.StatusOK)
}

func generateKey(c *gin.Context) {
	str, ok := c.GetQuery("user")
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := security.CreateToken(str)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
	return
}

var RouterModule = fx.Option(
	fx.Provide(RegisterRoutes),
)
