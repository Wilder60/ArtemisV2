package web

import (
	"net/http"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/middleware"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/adapter"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// New will return a default, engine with no handlers or middleware used for
// inserting mocks to unit test with
func New() *gin.Engine {
	return gin.New()
}

// CreateEngine will take
func CreateEngine(cal *adapter.Calendar, sec *security.Security, logger *zap.Logger, mid *middleware.HTTP) *gin.Engine {
	engine := gin.Default()
	engine.GET("key", func(ctx *gin.Context) {
		token, err := sec.CreateToken("TESTUSER")
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, token)
	})

	calendarGroup := engine.Group("api/v1/calendar")
	calendarGroup.Use(mid.Authorize())
	calendarGroup.Use(mid.AddUser())
	registerCalendarHandlers(calendarGroup, cal)
	return engine
}

func registerCalendarHandlers(group *gin.RouterGroup, api Calendar) {
	// GET api/v1/calendar?limit={val}&offset={val}
	group.GET("/", api.GetPaginatedEvents)
	group.GET("/range", api.GetEventsInRange)
	group.POST("/event", api.AddEvent)
	group.PATCH("/event", api.UpdateEvent)
	group.DELETE("/event", api.DeleteEvent)
}

var EngineModule = fx.Option(
	fx.Provide(CreateEngine),
)
