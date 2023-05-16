package testenv

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

func NewLogger(out io.Writer) (logr.Logger, func() error, error) {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	sink := zapcore.AddSync(out)
	level := zap.NewAtomicLevelAt(zap.DebugLevel)

	core := zapcore.NewCore(encoder, sink, level)
	logger := zap.New(core).WithOptions(
		zap.ErrorOutput(sink),
		zap.Development(),
	)
	log := zapr.NewLogger(logger)
	return log, logger.Sync, nil
}
