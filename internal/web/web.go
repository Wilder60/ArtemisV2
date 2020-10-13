package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

//Note more functions can be added here if metadata
//will be added

func ProvideBase() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return router
}

var ModuleBase = fx.Options(
	fx.Provide(ProvideBase),
)
