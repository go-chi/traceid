package traceid

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

var Header = http.CanonicalHeaderKey("X-TraceId")

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceID := r.Header.Get(Header)
		if uuid.FromStringOrNil(traceID) == uuid.Nil {
			traceID = uuidV7()
		}

		ctx = context.WithValue(ctx, ctxKey{}, traceID)
		w.Header().Set(Header, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ctxKey struct{}

func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

func NewContext(ctx context.Context) context.Context {
	if v, ok := ctx.Value(ctxKey{}).(string); ok && v != "" {
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
