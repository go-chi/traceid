# TraceId

Go package that helps create and pass `TraceId` header among microservices for simple tracing capabilities and log grouping.

The generated `TraceId` value is [UUIDv7](https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#name-uuid-version-7), which lets you infer the time of the trace creation from its value.

## traceid.Middleware example

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

## traceid.SetHeader()

```go
func main() {
    // Set TraceId in context, if not set yet.
    ctx := traceid.NewContext(context.Background())

    // Make a request with TraceId header.
    req, _ := http.NewRequest("GET", "http://localhost:3333/proxy", nil)
    req.WithContext(ctx)
    traceid.SetHeader(ctx, req)
    
    resp, err := resp.Do(req)
    //...
}
}
```

## Get time from UUIDv7 value

```
$ go run github.com/go-chi/traceid/cmd/traceid 018e0ee7-3605-7d75-b344-01062c6fd8bc
2024-03-05 14:56:57.477 +0100 CET
```

### Create new UUIDv7:
```
$ go run github.com/go-chi/traceid/cmd/traceid
018e0ee7-3605-7d75-b344-01062c6fd8bc
```

## License
[MIT](./LICENSE)