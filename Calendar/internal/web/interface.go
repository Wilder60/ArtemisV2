package web

import (
	"github.com/gin-gonic/gin"
)

type calendar interface {
	GetEventsInRange(*gin.Context)
	GetPaginatedEvents(*gin.Context)
	AddEvent(*gin.Context)
	UpdateEvent(*gin.Context)
	DeleteEvent(*gin.Context)
}
