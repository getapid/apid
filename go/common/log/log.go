package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// L is global logger
	L *zap.SugaredLogger
)

// Init setus up the logger, takes one parameter - the log level above which it will
// redirect logs to stderr, below which it will redirect logs to stdout
func Init(lvl int) {
	globalLevel := zapcore.Level(lvl)
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= globalLevel && lvl < zapcore.ErrorLevel
	})
	consoleInfos := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	ecfg := zap.NewProductionEncoderConfig()
	ecfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	ecfg.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("Mon Jan _2 15:04:05 2006"))
	}
	consoleEncoder := zapcore.NewConsoleEncoder(ecfg)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority),
	)
	logger := zap.New(core)
	zap.RedirectStdLog(logger)
	L = logger.Sugar()
}
