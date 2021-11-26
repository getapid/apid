package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// L is global logger
	L    *zap.SugaredLogger
	Atom zap.AtomicLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
)

// Init setus up the logger, takes one parameter - the log level above which it will
// redirect logs to stderr, below which it will redirect logs to stdout
func init() {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= Atom.Level() && lvl < zapcore.ErrorLevel
	})
	consoleInfos := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	ecfg := zap.NewProductionEncoderConfig()
	ecfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	ecfg.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {}
	// ecfg.EncodeLevel = func(lvl zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {}

	consoleEncoder := zapcore.NewConsoleEncoder(ecfg)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority),
	)
	logger := zap.New(core)
	zap.RedirectStdLog(logger)
	L = logger.Sugar()
}
