package logger

import (
	"context"
	"time"

	"go.uber.org/fx"

	"cloud.google.com/go/logging"

	"github.com/Wilder60/ArtemisV2/Calendar/config"
)

var _ Logger = Cloud{}

type Cloud struct {
	logger *logging.Logger
}

func ProvideCloudModule(cfg config.Config) (*Cloud, error) {
	ctx := context.Background()
	client, err := logging.NewClient(ctx, cfg.Logger.Project)
	if err != nil {
		return nil, err
	}

	logName := cfg.Logger.Name

	return &Cloud{
		logger: client.Logger(logName),
	}, nil
}

func (c *Cloud) Close() {
	c.Close()
}

func (c Cloud) log(msg string, lvl logging.Severity) {
	entry := logging.Entry{
		Timestamp: time.Now(),
		Severity:  lvl,
		Payload:   msg,
	}
	c.logger.Log(entry)
}

func (c Cloud) Debug(msg string) {
	c.log(msg, logging.Debug)

}

func (c Cloud) Info(msg string) {
	c.log(msg, logging.Info)
}

func (c Cloud) Warn(msg string) {
	c.log(msg, logging.Warning)
}

func (c Cloud) Error(msg string) {
	c.log(msg, logging.Error)
	c.logger.Flush()
}

func (c Cloud) Panic(msg string) {
	c.log(msg, logging.Error)
	c.logger.Flush()
	panic(msg)
}

var CloudLoggerModule = fx.Option(
	fx.Provide(ProvideCloudModule),
)
