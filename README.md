# TraceId

Go pkg to propagate `TraceId` header across multiple services.

- Enables simple tracing capabilities and log grouping.
- The value can be exposed to end-users in case of an error.
- The value is [UUIDv7](https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#name-uuid-version-7), which lets you [infer the timestamp](https://github.com/go-chi/traceid?tab=readme-ov-file#get-time-from-uuidv7-value).
- Unlike [OTEL](https://pkg.go.dev/go.opentelemetry.io/otel), this package doesn't trace spans or metrics, doesn't require any backend and doesn't have large dependencies (GRPC).

## Example - HTTP middleware

```go
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/traceid"
)

func main() {
	r := chi.NewRouter()

	r.Use(traceid.Middleware)
	r.Use(httplog.RequestLogger(logger()))
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Log traceId to request logger.
			traceID := traceid.FromContext(r.Context())
			httplog.LogEntrySetField(ctx, "traceId", slog.StringValue(traceID))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
```

See [example/main.go](./example/main.go)

## Example: Send HTTP request with TraceId header

```go
import (
	"github.com/go-chi/traceid"
)

func main() {
	// Set TraceId in context, if not set from parent ctx yet.
	ctx := traceid.NewContext(context.Background())

	// Make a request with TraceId header.
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:3333/proxy", nil)
	traceid.SetHeader(ctx, req)

	resp, err := http.DefaultClient.Do(req)
	//...
}
```

## Example: Set TraceId header in all outgoing HTTP requests globally

```go
import (
	"github.com/go-chi/traceid"
	"github.com/go-chi/transport"
)

func main() {
	// Send TraceId in all outgoing HTTP requests.
	http.DefaultTransport = transport.Chain(
		http.DefaultTransport,
		transport.SetHeader("User-Agent", "my-app/v1.0.0"),
		traceid.Transport,
	)
	
	// This will automatically send TraceId header.
	req, _ := http.NewRequest("GET", "http://localhost:3333/proxy", nil)
	_, _ = http.DefaultClient.Do(req)
}
```

## Get time from UUIDv7 value

```
$ go run github.com/go-chi/traceid/cmd/traceid 018e0ee7-3605-7d75-b344-01062c6fd8bc
2024-03-05 14:56:57.477 +0100 CET
```

You can also create a new UUIDv7:
```
$ go run github.com/go-chi/traceid/cmd/traceid
018e0ee7-3605-7d75-b344-01062c6fd8bc
```

## License
[MIT](./LICENSE)