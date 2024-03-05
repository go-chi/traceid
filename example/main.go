package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Get("/proxy", func(w http.ResponseWriter, r *http.Request) {
		reverseProxyTo, _ := url.Parse("http://localhost:3333")
		r.URL.Path = "/"
		proxy := httputil.NewSingleHostReverseProxy(reverseProxyTo)
		proxy.ServeHTTP(w, r)
	})

	r.Get("/proxy/proxy", func(w http.ResponseWriter, r *http.Request) {
		reverseProxyTo, _ := url.Parse("http://localhost:3333")
		r.URL.Path = "/proxy"
		proxy := httputil.NewSingleHostReverseProxy(reverseProxyTo)
		proxy.ServeHTTP(w, r)
	})

	fmt.Println(`Try sending a sample request in another terminal and see traceId propagation:`)
	fmt.Println(`$ curl http://localhost:3333/proxy/proxy`)
	fmt.Println(`$ curl -H "TraceId: 018e0f8e-c6e2-7aae-bbf1-000000000000" http://localhost:3333/proxy/proxy`)
	fmt.Println()
	fmt.Println(`server listening on :3333`)

	http.ListenAndServe(":3333", r)
}

func logger() *httplog.Logger {
	return httplog.NewLogger("httplog-example", httplog.Options{
		LogLevel:         slog.LevelDebug,
		RequestHeaders:   false,
		ResponseHeaders:  false,
		JSON:             false,
		Concise:          true,
		MessageFieldName: "message",
		LevelFieldName:   "severity",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"version": "v0.1.0",
			"env":     "dev",
		},
	})
}
