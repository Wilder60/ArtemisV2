package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/Wilder60/Calendar/internal/adapter"
	"github.com/Wilder60/Calendar/internal/security"
	"github.com/Wilder60/Calendar/internal/storage"
	"github.com/Wilder60/Calendar/internal/web"

	"github.com/Wilder60/Calendar/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func start(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, router *gin.Engine, config *config.Config) {
	srv := &http.Server{
		Handler:      router,
		Addr:         config.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go srv.ListenAndServe()

				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)

				// Block until a signal is received.
				go func() {
					s := <-c
					fmt.Println(s.String)
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
		fx.Provide(zap.NewProduction),
		config.ConfigModule,
		security.SecurityModule,
		storage.FirebaseModule,
		adapter.CalendarHandlerModule,
		web.EngineModule,
		fx.Invoke(start),
	).Run()
}