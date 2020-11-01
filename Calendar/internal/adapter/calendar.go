package adapter

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Wilder60/ShadowKeep/internal/domain/requests"

	"github.com/Wilder60/ShadowKeep/internal/domain"

	"github.com/Wilder60/ShadowKeep/internal/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

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

// GET api/v1/calendar?time=string&limit=int&offset=int
func (c *Calendar) GetPaginatedEvents(requestCtx *gin.Context) {
	backgroundCtx := context.Background()
	limitParam := requestCtx.Query("limit")
	offsetParam := requestCtx.Query("offset")
	parsedParams, err := parseIntParameters(limitParam, offsetParam)
	if err != nil {
		requestCtx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	limit, offset := parsedParams[0], parsedParams[1]
	// NOTE: think about replacing this with a struct instead, could make code cleaner
	events, err := c.db.GetEventsPaginated(backgroundCtx,
		requestCtx.GetString("UserID"),
		requestCtx.Param(""),
		limit,
		offset,
	)
	if err != nil {
		requestCtx.Status(http.StatusInternalServerError)
		return
	}
	requestCtx.JSON(http.StatusOK, events)
}

// GET api/v1/calendar/range?sdate=string&edate=string
func (c *Calendar) GetEventsInRange(requestCtx *gin.Context) {
	ctx := context.Background()
	sdate := requestCtx.Query("sdate")
	edate := requestCtx.Query("edate")
	userID := requestCtx.GetString("UserID")
	events, err := c.db.GetEventsInRange(ctx, userID, sdate, edate)
	if err != nil {
		requestCtx.Status(http.StatusInternalServerError)
		return
	}
	requestCtx.JSON(http.StatusOK, events)
}

// POST api/v1/calendar
func (c *Calendar) AddEvent(requestCtx *gin.Context) {
	ctx := context.Background()
	userID := requestCtx.GetString("UserID")
	event := domain.Event{}
	requestCtx.Bind(&event)
	event.UserID = userID
	err := c.db.CreateEvents(ctx, event)
	if err != nil {
		requestCtx.Status(http.StatusInternalServerError)
		return
	}
	requestCtx.Status(http.StatusNoContent)
}

// PATCH api/v1/calendar
func (c *Calendar) UpdateEvent(requestCtx *gin.Context) {
	ctx := context.Background()
	event := domain.Event{}
	requestCtx.Bind(&event)
	userID := requestCtx.GetString("UserID")
	if userID != event.UserID {
		requestCtx.Status(http.StatusBadRequest)
		return
	}
	err := c.db.UpdateEvent(ctx, event)
	if err != nil {
		requestCtx.Status(http.StatusInternalServerError)
		return
	}
	requestCtx.Status(http.StatusNoContent)
}

// DeleteEvent is the handler for the api/v1/calendar DELETE endpoint
func (c *Calendar) DeleteEvent(requestCtx *gin.Context) {
	ctx := context.Background()
	deleteRequest := requests.Delete{}
	requestCtx.Bind(&deleteRequest)
	userID := requestCtx.GetString("UserID")
	err := c.db.DeleteEvents(ctx, userID, deleteRequest.IDs[0])
	if err != nil {
		requestCtx.Status(http.StatusInternalServerError)
		return
	}
	requestCtx.Status(http.StatusNoContent)
}

func parseIntParameters(parameters ...string) ([]int, error) {
	output := make([]int, len(parameters))
	for idx, param := range parameters {
		val, err := strconv.Atoi(param)
		if err != nil {
			return nil, err
		}
		output[idx] = val
	}
	return output, nil
}

// CalendarHandlerModule is the exported
var CalendarHandlerModule = fx.Option(
	fx.Provide(ProvideCalendar),
)
