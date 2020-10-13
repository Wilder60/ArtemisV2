package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Wilder60/KeyRing/internal/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/Wilder60/KeyRing/configs"
	"github.com/Wilder60/KeyRing/internal/web"
)

func start(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, router *gin.Engine, config *configs.Config) {
	srv := &http.Server{
		Handler:      router,
		Addr:         config.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				fmt.Println("Starting server")
				go srv.ListenAndServe()

				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)

				// Block until a signal is received.
				go func() {
					s := <-c
					fmt.Println(s.String())
					if err := shutdowner.Shutdown(); err != nil {
						os.Exit(1)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		},
	)
}

func main() {

	fx.New(
		configs.Module,
		web.ModuleBase,
		sql.ModuleCloudSql,
		fx.Invoke(start),
	).Run()

	// srv := &http.Server{
	// 	Handler:      adapter.NewWebAdapter(),
	// 	Addr:         "127.0.0.1:8001",
	// 	WriteTimeout: 15 * time.Second,
	// 	ReadTimeout:  15 * time.Second,
	// }

	fmt.Println("Starting server")

	// log.Fatal(srv.ListenAndServe())

}
