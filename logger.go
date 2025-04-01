package traceid

import (
	"context"
	"log/slog"
)

/*
You can pass slog.Handler and logKey, and all the logs would automatically log
the traceId if logging with context is used.

	slog.Log()
	slog.DebugContext()
	slog.InfoContext()
	slog.WarnContext()
	slog.ErrorContext()

Example:

	ctx := traceid.NewContext(context.Background())
	handler := traceid.NewLoggerWrapper(slog.NewJSONHandler(os.Stdout, nil), "traceId")
	logger := slog.New(handler)
	logger.InfoContext(ctx, "message")

This would log

	{"time":"2025-04-01T13:17:43.097789397Z","level":"INFO","msg":"message","traceId":"0195f180-2939-7bc4-bffe-838eb3c62526"}
*/
func NewLoggerWrapper(handler slog.Handler, logKey string) slog.Handler {
	return &wrapperLogger{
		handler: handler,
		key:     logKey,
	}
}

type wrapperLogger struct {
	handler slog.Handler
	key     string
}

func (l *wrapperLogger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.handler.Enabled(ctx, level)
}

func (l *wrapperLogger) Handle(ctx context.Context, record slog.Record) error {
	if traceID := FromContext(ctx); traceID != "" {
		record.AddAttrs(slog.String(l.key, traceID))
	}

	return l.handler.Handle(ctx, record)
}

func (l *wrapperLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &wrapperLogger{
		handler: l.handler.WithAttrs(attrs),
		key:     l.key,
	}
}

func (l *wrapperLogger) WithGroup(name string) slog.Handler {
	return &wrapperLogger{
		handler: l.handler.WithGroup(name),
		key:     l.key,
	}
}
