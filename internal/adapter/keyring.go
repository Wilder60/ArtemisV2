package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Wilder60/KeyRing/internal/domain"
	"github.com/Wilder60/KeyRing/internal/security"
	"github.com/Wilder60/KeyRing/internal/sql"
	"github.com/gin-gonic/gin"
)

type (
	keyRing struct {
		*sql.SQL
	}
)

// InitRoutes will
func initKeyRing(r *gin.Engine, db *sql.SQL, endpoint string) {
	keyRing := &keyRing{db}

	r.GET("/KeyRing", keyRing.getEvents)
}

func (kr keyRing) getEvents(ctx *gin.Context) {
	queryValues := ctx.Request.URL.Query()
	limitStr := queryValues.Get("limit")
	offsetStr := queryValues.Get("offset")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	var offset int64
	if offsetStr != "" {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}

	entries, err := kr.GetKeyRing("", limit, offset)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.JSON(http.StatusOK, entries)
}

func (kr keyRing) addEvent(w http.ResponseWriter, r *http.Request) {
	var event domain.KeyEntry
	if err := decodeRequest(r.Body, &event); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := security.GetUserFromToken(r.Header.Get("Authorization"))

	ret, err := kr.AddKeyRing(event, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
