package traceid

import (
	"context"
	"log/slog"
)

var LogKey string = "traceId"

/*
You can pass slog.Handler, and all the logs would automatically log
the traceId if logging with context is used.

	slog.Log()
	slog.DebugContext()
	slog.InfoContext()
	slog.WarnContext()
	slog.ErrorContext()

Example:

	ctx := traceid.NewContext(context.Background())
	handler := traceid.LogHandler(slog.NewJSONHandler(os.Stdout, nil))
	logger := slog.New(handler)
	logger.InfoContext(ctx, "message")

This would log

	{"time":"2025-04-01T13:17:43.097789397Z","level":"INFO","msg":"message","traceId":"0195f180-2939-7bc4-bffe-838eb3c62526"}
*/
func LogHandler(handler slog.Handler) slog.Handler {
	return &logHandler{
		handler: handler,
	}
}

type logHandler struct {
	handler slog.Handler
}

func (l *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return l.handler.Enabled(ctx, level)
}

func (l *logHandler) Handle(ctx context.Context, record slog.Record) error {
	if traceID := FromContext(ctx); traceID != "" {
		record.AddAttrs(slog.String(LogKey, traceID))
	}

	return l.handler.Handle(ctx, record)
}

func (l *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{
		handler: l.handler.WithAttrs(attrs),
	}
}

func (l *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{
		handler: l.handler.WithGroup(name),
	}
}
