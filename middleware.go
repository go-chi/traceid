package traceid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceID := r.Header.Get(Header)
		if _, err := uuid.Parse(traceID); err != nil {
			traceID = uuidV7()
		}

		ctx = context.WithValue(ctx, ctxKey{}, traceID)
		w.Header().Set(Header, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
