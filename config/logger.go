package config

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoggerBuild() (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
	})
	zapConfig.EncoderConfig.TimeKey = "time"

	return zapConfig.Build()
}
