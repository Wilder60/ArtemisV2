package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cloud.google.com/go/logging"
	"github.com/Wilder60/KeyRing/configs"
	"github.com/Wilder60/KeyRing/internal/logger"
	"github.com/Wilder60/KeyRing/internal/security"
	"github.com/Wilder60/KeyRing/internal/sql"
	"github.com/Wilder60/KeyRing/internal/web"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func start(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, router *gin.Engine, config *configs.Config, logger *logging.Logger) {
	srv := &http.Server{
		Handler:      router,
		Addr:         config.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Log(
					logging.Entry{
						Payload: "Starting Server on port " + config.Server.Port,
					},
				)
				err := logger.Flush()
				if err != nil {
					panic(err)
				}
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
		sql.ModuleCloudSql,
		sql.KeyRingSQLModule,
		web.KeyRingModule,
		web.RouterModule,
		logger.CloudLoggerModule,
		security.SecurityModule,
		fx.Invoke(start),
	).Run()

	fmt.Println("Starting server")
}
