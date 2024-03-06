package traceid

import (
	"net/http"
)

func Transport(next http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		r := cloneRequest(req)

		traceID, _ := r.Context().Value(ctxKey{}).(string)
		if traceID == "" {
			traceID = uuidV7()
		}
		r.Header.Set(http.CanonicalHeaderKey(Header), traceID)

		return next.RoundTrip(r)
	})
}

// RoundTripFunc, similar to http.HandlerFunc, is an adapter
// to allow the use of ordinary functions as http.RoundTrippers.
type RoundTripFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// cloneRequest creates a shallow copy of a given request
// to comply with stdlib's http.RoundTripper contract:
//
// RoundTrip should not modify the request, except for
// consuming and closing the Request's Body. RoundTrip may
// read fields of the request in a separate goroutine. Callers
// should not mutate or reuse the request until the Response's
// Body has been closed.
func cloneRequest(orig *http.Request) *http.Request {
	clone := &http.Request{}
	*clone = *orig

	clone.Header = make(http.Header, len(orig.Header))
	for key, value := range orig.Header {
		clone.Header[key] = append([]string{}, value...)
	}

	return clone
}
