package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(level, format, outputPath string) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var logWriter zerolog.LevelWriter
	switch outputPath {
	case "stdout":
		logWriter = zerolog.MultiLevelWriter(os.Stdout)
	default:
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logWriter = zerolog.MultiLevelWriter(os.Stdout)
		} else {
			logWriter = zerolog.MultiLevelWriter(os.Stdout, file)
		}
	}

	zerolog.SetGlobalLevel(parseLevel(level))

	logger := zerolog.New(logWriter).With().Timestamp().Logger()

	return &Logger{
		logger: &logger,
	}
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	e := l.logger.Debug()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			e = e.Interface(fields[i].(string), fields[i+1])
		}
	}
	e.Msg(msg)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	e := l.logger.Info()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			e = e.Interface(fields[i].(string), fields[i+1])
		}
	}
	e.Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	e := l.logger.Warn()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			e = e.Interface(fields[i].(string), fields[i+1])
		}
	}
	e.Msg(msg)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	e := l.logger.Error()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			e = e.Interface(fields[i].(string), fields[i+1])
		}
	}
	e.Msg(msg)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	e := l.logger.Fatal()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			e = e.Interface(fields[i].(string), fields[i+1])
		}
	}
	e.Msg(msg)
}

func (l *Logger) WithRequestID(ctx *gin.Context) *zerolog.Event {
	requestID := ctx.GetString("request_id")
	if requestID == "" {
		requestID = ctx.GetHeader("X-Request-ID")
	}
	logger := l.logger.With().Str("request_id", requestID).Logger()
	return logger.Info()
}

func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

func (l *Logger) InfoEvent() *zerolog.Event {
	return l.logger.Info()
}

func (l *Logger) ErrorEvent() *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) DebugEvent() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) WarnEvent() *zerolog.Event {
	return l.logger.Warn()
}

func GetLoggerFromContext(ctx *gin.Context) *zerolog.Logger {
	if logger, exists := ctx.Get("logger"); exists {
		return logger.(*zerolog.Logger)
	}
	return &log.Logger
}

func Middleware(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Request.Context().Value("start")
		if start == nil {
			c.Next()
			return
		}

		c.Next()

		latency := c.Request.Context().Value("latency")
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		entry := logger.logger.Info()
		if statusCode >= 400 {
			entry = logger.logger.Error()
		}

		entry.
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Interface("latency", latency).
			Msg("request completed")
	}
}
