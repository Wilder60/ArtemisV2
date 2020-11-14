package adapter

import (
	"fmt"
	"net/http"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/db"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain/requests"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var _ Storage = db.Firestore{}

type Calendar struct {
	db Storage
}

// New is for unit testing only, and should not be used otherwise
func New() *Calendar {
	return &Calendar{}
}

// ProvideCalendar is the provider function for
func ProvideCalendar(s *db.Firestore) *Calendar {
	return &Calendar{db: s}
}

// GET api/v1/calendar?time=string&limit=int&offset=int&desc=bool
func (c *Calendar) GetPaginatedEvents(ctx *gin.Context) {
	userID := ctx.GetString("UserID")
	paginationRequest := requests.NewGetPagination(userID)
	bindErr := ctx.BindQuery(paginationRequest)
	if bindErr != nil {
		ctx.Status(http.StatusBadRequest)
	}

	// NOTE: think about replacing this with a struct instead, could make code cleaner
	events, err := c.db.GetEventsPaginated(paginationRequest)
	if err != nil {
		fmt.Println(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, events)
}

// GET api/v1/calendar/range?sdate=string&edate=string
func (c *Calendar) GetEventsInRange(ctx *gin.Context) {
	userID := ctx.GetString("UserID")
	rangeRequest := requests.NewGetRange(userID)
	bindErr := ctx.BindQuery(rangeRequest)
	if bindErr != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	events, err := c.db.GetEventsInRange(rangeRequest)
	if err != nil {
		fmt.Println(err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, events)
}

// POST api/v1/calendar
func (c *Calendar) AddEvent(ctx *gin.Context) {
	userID := ctx.GetString("UserID")
	addRequest := requests.NewAdd(userID)
	bindErr := ctx.Bind(&addRequest)
	if bindErr != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := c.db.CreateEvents(addRequest)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// PATCH api/v1/calendar
func (c *Calendar) UpdateEvent(ctx *gin.Context) {
	userID := ctx.GetString("UserID")
	updateRequest := requests.NewUpdate(userID)
	bindErr := ctx.Bind(&updateRequest)
	if bindErr != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := c.db.UpdateEvent(updateRequest)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// DeleteEvent is the handler for the api/v1/calendar DELETE endpoint
func (c *Calendar) DeleteEvent(requestCtx *gin.Context) {
	userID := requestCtx.GetString("UserID")
	deleteRequest := requests.NewDelete(userID)
	requestCtx.Bind(&deleteRequest)
	err := c.db.DeleteEvents(deleteRequest)
	if err != nil {
		requestCtx.Status(http.StatusInternalServerError)
		return
	}
	requestCtx.Status(http.StatusNoContent)
}

// CalendarHandlerModule is the exported
var CalendarHandlerModule = fx.Option(
	fx.Provide(ProvideCalendar),
)
