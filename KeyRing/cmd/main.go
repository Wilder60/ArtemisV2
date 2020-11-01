package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Wilder60/KeyRing/configs"
	"github.com/Wilder60/KeyRing/internal/security"
	"github.com/Wilder60/KeyRing/internal/sql"
	"github.com/Wilder60/KeyRing/internal/web"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func start(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, router *gin.Engine, config *configs.Config, logger *zap.Logger) {
	srv := &http.Server{
		Handler:      router,
		Addr:         config.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Starting server on port " + config.Server.Port)
				go srv.ListenAndServe()

				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)

				// Block until a signal is received.
				go func() {
					s := <-c
					logger.Info(fmt.Sprintf("Received Signal %s", s.String()))
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
		configs.Module,
		sql.ModuleCloudSql,
		sql.KeyRingSQLModule,
		web.KeyRingModule,
		web.RouterModule,
		security.SecurityModule,
		fx.Invoke(start),
	).Run()
}
