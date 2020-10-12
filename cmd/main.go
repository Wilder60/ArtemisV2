package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/fx"

	"github.com/Wilder60/KeyRing/internal/adapter"
)

func main() {

	fx.New()

	srv := &http.Server{
		Handler:      adapter.NewWebAdapter(),
		Addr:         "127.0.0.1:8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server")

	log.Fatal(srv.ListenAndServe())

}
