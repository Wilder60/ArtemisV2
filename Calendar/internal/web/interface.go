package web

import (
	"github.com/gin-gonic/gin"
)

type Calendar interface {
	GetEventsInRange(*gin.Context)
	GetPaginatedEvents(*gin.Context)
	AddEvent(*gin.Context)
	UpdateEvent(*gin.Context)
	DeleteEvent(*gin.Context)
}
