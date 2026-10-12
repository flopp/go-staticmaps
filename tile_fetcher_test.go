package sm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// closeTrackingBody wraps an io.ReadCloser and records whether Close was called.
type closeTrackingBody struct {
	io.ReadCloser
	closed bool
}

func (b *closeTrackingBody) Close() error {
	b.closed = true
	return b.ReadCloser.Close()
}

// closeTrackingTransport wraps a Transport and captures the closeTrackingBody
// of the last response it produced, so a test can check whether it was closed.
type closeTrackingTransport struct {
	http.RoundTripper
	lastBody *closeTrackingBody
}

func (t *closeTrackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	wrapped := &closeTrackingBody{ReadCloser: resp.Body}
	t.lastBody = wrapped
	resp.Body = wrapped
	return resp, nil
}

// withInstrumentedDefaultClient temporarily swaps http.DefaultClient's
// Transport for one that tracks response body closes, and returns a
// restore func plus the tracker itself.
func withInstrumentedDefaultClient(t *testing.T) (*closeTrackingTransport, func()) {
	t.Helper()
	base := http.DefaultClient.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	tracker := &closeTrackingTransport{RoundTripper: base}
	orig := http.DefaultClient.Transport
	http.DefaultClient.Transport = tracker
	return tracker, func() {
		http.DefaultClient.Transport = orig
	}
}

// TestTileFetcherDownloadClosesBodyOn404 verifies that TileFetcher.download
// closes the HTTP response body even when the server responds with a
// non-200 status. Leaving it unclosed leaks the underlying connection's
// file descriptor since it is never returned to the transport's pool.
func TestTileFetcherDownloadClosesBodyOn404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer srv.Close()

	tracker, restore := withInstrumentedDefaultClient(t)
	defer restore()

	provider := NewTileProviderOpenStreetMaps()
	fetcher := NewTileFetcher(provider, nil, true)

	_, err := fetcher.download(srv.URL)
	if err != errTileNotFound {
		t.Fatalf("expected errTileNotFound, got %v", err)
	}

	if tracker.lastBody == nil {
		t.Fatal("no response body was observed")
	}
	if !tracker.lastBody.closed {
		t.Error("response body was not closed on 404 response, leaking the underlying connection")
	}
}

// TestTileFetcherDownloadClosesBodyOnServerError mirrors the 404 case for
// the generic "unexpected status code" error path.
func TestTileFetcherDownloadClosesBodyOnServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	tracker, restore := withInstrumentedDefaultClient(t)
	defer restore()

	provider := NewTileProviderOpenStreetMaps()
	fetcher := NewTileFetcher(provider, nil, true)

	_, err := fetcher.download(srv.URL)
	if err == nil {
		t.Fatal("expected an error for 500 response")
	}

	if tracker.lastBody == nil {
		t.Fatal("no response body was observed")
	}
	if !tracker.lastBody.closed {
		t.Error("response body was not closed on 500 response, leaking the underlying connection")
	}
}
