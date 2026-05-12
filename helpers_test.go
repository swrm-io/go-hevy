package hevy_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/swrm-io/go-hevy"
)

func ptr[T any](v T) *T { return &v }

// newTestServer creates an httptest.Server that responds to a single path with
// the contents of a testdata JSON file, using the given status code.
func newTestServer(t *testing.T, path, fixture string, status int) (*httptest.Server, *hevy.Client) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", fixture))
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixture, err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("User-Agent") == "" {
			t.Errorf("request to %s missing User-Agent header", r.URL.Path)
		}
		if r.Header.Get("api-key") == "" {
			t.Errorf("request to %s missing api-key header", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(body)
	}))
	t.Cleanup(srv.Close)
	client := hevy.New("test-key", hevy.WithBaseURL(srv.URL))
	return srv, client
}

// newErrorServer returns a server that always responds with the given status code.
func newErrorServer(t *testing.T, status int) (*httptest.Server, *hevy.Client) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(`{"error":"test error"}`))
	}))
	t.Cleanup(srv.Close)
	client := hevy.New("test-key", hevy.WithBaseURL(srv.URL))
	return srv, client
}
