package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Wilder60/KeyRing/internal/domain"
	"github.com/Wilder60/KeyRing/internal/interfaces"

	"github.com/gorilla/mux"
)

type (
	keyRing struct {
		interfaces.Database
	}
)

// InitRoutes will
func initKeyRing(router *mux.Router, db interfaces.Database, endpoint string) {
	keyRing := &keyRing{db}

	router.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	// router.HandleFunc("/KeyRing", getEvents).Methods(http.MethodGet)
	router.HandleFunc("/KeyRing", keyRing.addEvent).Methods(http.MethodPost)
	// router.HandleFunc("/KeyRing", updateEvent).Methods(http.MethodPatch)
	// router.HandleFunc("/KeyRing", deleteEvents).Methods(http.MethodDelete)
}

func (kr keyRing) getEvents(w http.ResponseWriter, r *http.Request) {

}

func (kr keyRing) addEvent(w http.ResponseWriter, r *http.Request) {
	var event domain.KeyEntry
	if err := decodeRequest(r.Body, &event); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ret, _ := kr.AddKeyRing(event)
	w.WriteHeader(200)
	w.Write([]byte(strconv.FormatInt(ret, 10)))
}

func (kr keyRing) updateEvent(w http.ResponseWriter, r *http.Request) {
	var event domain.KeyEntry
	defer r.Body.Close()
	if err := decodeRequest(r.Body, &event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
}

func (kr keyRing) deleteEvents(w http.ResponseWriter, r *http.Request) {

}

func decodeRequest(body io.ReadCloser, out interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(&out)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Okay!"))
}