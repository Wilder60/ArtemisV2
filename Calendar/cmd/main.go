package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/logger"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/middleware"

	"github.com/Wilder60/ArtemisV2/Calendar/config"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/adapter"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/db"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/security"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/web"
	"github.com/soheilhy/cmux"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func start(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, router *gin.Engine, server *grpc.Server,
	config *config.Config, logger *logger.Zap) {

	httpSrv := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info(fmt.Sprintf("starting Service on port %s", config.Server.Port))
				l, err := net.Listen("tcp", config.Server.Port)
				if err != nil {
					log.Fatal(err)
				}

				m := cmux.New(l)
				grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
				httpL := m.Match(cmux.HTTP1Fast())

				go server.Serve(grpcL)
				go httpSrv.Serve(httpL)
				go m.Serve()

				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt)

				// Block until a signal is received.
				go func() {
					s := <-c
					logger.Info(fmt.Sprintf("Shutdown signal received %s", s.String()))
					if err := shutdowner.Shutdown(); err != nil {
						os.Exit(1)
					}
				}()

				return nil
			},

			OnStop: func(ctx context.Context) error {
				logger.Info("shutting down service")
				server.GracefulStop()
				return httpSrv.Shutdown(ctx)
			},
		},
	)
}

func main() {

	fx.New(
		config.ConfigModule,
		security.SecurityModule,
		db.FirebaseModule,
		middleware.GRPCMiddlewareModule,
		middleware.HTTPMiddlewareModule,
		adapter.CalendarHandlerModule,
		logger.ZapLoggerModule,
		adapter.GRPCModule,
		web.EngineModule,
		fx.Invoke(start),
	).Run()
}
