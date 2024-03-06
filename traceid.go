package traceid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

var Header = http.CanonicalHeaderKey("TraceId")

type ctxKey struct{}

func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

func NewContext(ctx context.Context) context.Context {
	if traceID, ok := ctx.Value(ctxKey{}).(string); ok && traceID != "" {
		return ctx
	}

	traceID := uuidV7()
	return context.WithValue(ctx, ctxKey{}, traceID)
}

func SetHeader(ctx context.Context, req *http.Request) {
	id, ok := ctx.Value(ctxKey{}).(string)
	if !ok || id == "" {
		id = uuidV7()
	}

	req.Header.Set(Header, id)
}

func uuidV7() string {
	// uuid.NewV7() only requires the current time from OS. Let the World panic if we can't get it.
	return uuid.Must(uuid.NewV7()).String()
}
