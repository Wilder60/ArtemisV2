package adpater

import (
	"net/http"

	"github.com/Wilder60/KeyRing/middleware"
	"github.com/gorilla/mux"
)

// InitRoutes will
func InitRoutes(router *mux.Router) {
	router.HandleFunc("/KeyRing", addEvent).Methods("POST")
	router.Use(middleware.Authorize)
}

func getEvents(w http.ResponseWriter, r *http.Request) {

}

func addEvent(w http.ResponseWriter, r *http.Request) {

}

func updateEvent(w http.ResponseWriter, r *http.Request) {

}

func deleteEvents(w http.ResponseWriter, r *http.Request) {

}
