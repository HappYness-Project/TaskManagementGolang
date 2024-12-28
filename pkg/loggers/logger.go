package loggers

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	setupOnce sync.Once
	appLogger *AppLogger
)

type AppLogger struct {
	zLogger zerolog.Logger
}

func Setup(env configs.Env) *AppLogger {
	setupOnce.Do(func() {
		appLogger = &AppLogger{}
		lvl := ZerologLevel(env.LogLevel)
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano
		var logDest io.Writer
		logDest = os.Stdout
		if utils.IsDevMode(env.AppEnv) {
			logDest = zerolog.ConsoleWriter{Out: logDest}
		}
		appLogger.zLogger = zerolog.New(logDest).With().Caller().Timestamp().Logger().Level(lvl)
	})
	return appLogger
}

// func (l *AppLogger) WithReqID(ctx *gin.Context) (zerolog.Logger, string) {
// 	if rID := ctx.Request.Context().Value(utils.ContextKey(utils.RequestIdentifier)); rID != nil {
// 		if reqID, ok := rID.(string); ok {
// 			return l.zLogger.With().Str(utils.RequestIdentifier, reqID).Logger(), reqID
// 		}
// 		return l.zLogger, ""
// 	}
// 	return l.zLogger, ""
// }

// Fatal logs a message with fatal level and exits the program.
func (l *AppLogger) Fatal() *zerolog.Event {
	return l.zLogger.Fatal()
}

// Error logs a message with error level.
func (l *AppLogger) Error() *zerolog.Event {
	return l.zLogger.Error()
}

// Info logs a message with info level.
func (l *AppLogger) Info() *zerolog.Event {
	return l.zLogger.Info()
}

// Debug logs a message with debug level.
func (l *AppLogger) Debug() *zerolog.Event {
	return l.zLogger.Debug()
}

func ZerologLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
func isDevMode(s string) bool {
	return strings.Contains(s, "local") || strings.Contains(s, "dev")
}
