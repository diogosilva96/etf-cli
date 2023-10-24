package data

import "net/http"

// headerTransport represents an http transport with default headers.
type headerTransport struct {
	headers map[string]string
	base    http.RoundTripper
}

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// attach headers
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}
