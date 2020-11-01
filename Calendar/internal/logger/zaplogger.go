package logger

import "go.uber.org/zap"

type ZapLogger struct {
	log *zap.Logger
}
