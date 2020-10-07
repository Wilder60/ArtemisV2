package adapter

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Wilder60/KeyRing/internal/interfaces"
)

func NewWebAdapter(db interfaces.Database) http.Handler {
	r := mux.NewRouter()

	initKeyRing(r, db, "/KeyRing")
	return r
}
