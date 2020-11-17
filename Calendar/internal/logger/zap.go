package logger

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Wilder60/ArtemisV2/Calendar/config"
)

var _ Logger = Zap{}

type Zap struct {
	Log *zap.Logger
}

func ProvideZapModule(cfg *config.Config) (*Zap, error) {
	// if the file in OutputPath does not exist yet, we want to create it
	if _, err := os.OpenFile(cfg.Logger.OutputPath, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		return nil, err
	}

	logCfg := zap.NewProductionConfig()
	logCfg.OutputPaths = []string{
		"stderr",
		"stdout",
		cfg.Logger.OutputPath,
	}

	logger, err := logCfg.Build()
	if err != nil {
		return nil, err
	}

	defer logger.Sync()

	return &Zap{
		Log: logger,
	}, nil
}

func (z Zap) Debug(msg string) {
	z.Log.Debug(msg)
}

func (z Zap) Info(msg string) {
	z.Log.Info(msg)
}

func (z Zap) Warn(msg string) {
	z.Log.Warn(msg)
}

func (z Zap) Error(msg string) {
	z.Log.Error(msg)
}

func (z Zap) Panic(msg string) {
	z.Log.Panic(msg)
}

var ZapLoggerModule = fx.Option(
	fx.Provide(ProvideZapModule),
)
