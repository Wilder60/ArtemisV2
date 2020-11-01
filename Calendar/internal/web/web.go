package web

import (
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
func CreateEngine(cal *adapter.Calendar, sec *security.Security, logger *zap.Logger) *gin.Engine {
	engine := gin.Default()
	calendarGroup := engine.Group("api/v1/calendar")
	calendarGroup.Use(Authorize(sec))
	calendarGroup.Use(AddUser(sec))
	registerCalendarHandlers(calendarGroup, cal)
	return engine
}

func registerCalendarHandlers(group *gin.RouterGroup, api calendar) {
	// GET api/v1/calendar?limit={val}&offset={val}
	group.GET("/", api.GetEventsInRange)
	group.GET("/range", api.GetEventsInRange)
	group.POST("/event", api.AddEvent)
	group.PATCH("/event", api.UpdateEvent)
	group.DELETE("/event", api.DeleteEvent)
}

var EngineModule = fx.Option(
	fx.Provide(CreateEngine),
)
