package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// L is global logger
	L        *zap.SugaredLogger
	onceInit sync.Once
)

// Init setus up the logger, takes one parameter - the log level above which it will
// redirect logs to stderr, below which it will redirect logs to stdout
func Init(lvl int) {
	onceInit.Do(func() {
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
		consoleEncoder := zapcore.NewJSONEncoder(ecfg)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority),
		)
		logger := zap.New(core)
		zap.RedirectStdLog(logger)
		L = logger.Sugar()
	})
}
