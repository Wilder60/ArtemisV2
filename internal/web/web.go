package adapter

import (
	"net/http"

	"github.com/Wilder60/KeyRing/internal/sql"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

//Note more functions can be added here if metadata
//will be added

func ProvideGinServer() *gin.Engine {
	gin.DisableConsoleColor()
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middle)
	router.GET("health", healthCheck)
}

func healthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

var Module = fx.Options(
	fx.Provide(ProvideGinServer)
)
