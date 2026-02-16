package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

var (
	log zerolog.Logger
)

func Init(devMode bool) {
	var w io.Writer = os.Stdout

	if devMode {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		w = consoleWriter
	}

	zerolog.TimeFieldFormat = time.RFC3339

	log = zerolog.New(w).
		With().
		Timestamp().
		Str("service", "subscriptions-api").
		Logger()
}

func FromContext(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return log
	}

	if l, ok := ctx.Value(ctxKey{}).(zerolog.Logger); ok {
		return l
	}

	return log
}

func WithContext(ctx context.Context, l zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func With(fields map[string]interface{}) zerolog.Logger {
	l := log
	for k, v := range fields {
		l = l.With().Interface(k, v).Logger()
	}
	return l
}

func Info(ctx context.Context, msg string, fields map[string]interface{}) {
	l := FromContext(ctx)
	e := l.Info()
	for k, v := range fields {
		e = e.Interface(k, v)
	}
	e.Msg(msg)
}

func Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	l := FromContext(ctx)
	e := l.Error().Err(err)
	for k, v := range fields {
		e = e.Interface(k, v)
	}
	e.Msg(msg)
}

func Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	l := FromContext(ctx)
	e := l.Debug()
	for k, v := range fields {
		e = e.Interface(k, v)
	}
	e.Msg(msg)
}

func Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	l := FromContext(ctx)
	e := l.Warn()
	for k, v := range fields {
		e = e.Interface(k, v)
	}
	e.Msg(msg)
}
