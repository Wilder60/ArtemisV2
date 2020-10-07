package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Wilder60/KeyRing/internal/sql"

	"github.com/Wilder60/KeyRing/internal/adapter"
)

func main() {

	sqlDriver := sql.New()

	srv := &http.Server{
		Handler:      adapter.NewWebAdapter(&sqlDriver),
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server")

	log.Fatal(srv.ListenAndServe())

}
