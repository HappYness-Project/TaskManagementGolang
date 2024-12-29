package loggers

import (
	"io"
	"net/http"
	"os"
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

func (l *AppLogger) WithReqID(r *http.Request) (zerolog.Logger, string) {
	if rID := r.Context().Value(utils.ContextKey(utils.RequestIdentifier)); rID != nil {
		if reqID, ok := rID.(string); ok {
			return l.zLogger.With().Str(utils.RequestIdentifier, reqID).Logger(), reqID
		}
		return l.zLogger, ""
	}
	return l.zLogger, ""
}

func (l *AppLogger) Fatal() *zerolog.Event {
	return l.zLogger.Fatal()
}

func (l *AppLogger) Error() *zerolog.Event {
	return l.zLogger.Error()
}

func (l *AppLogger) Info() *zerolog.Event {
	return l.zLogger.Info()
}

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
