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
// the contents of a testdata JSON file.
func newTestServer(t *testing.T, path, fixture string) *hevy.Client {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", fixture))
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixture, err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			http.NotFound(resp, req)
			return
		}
		if req.Header.Get("User-Agent") == "" {
			t.Errorf("request to %s missing User-Agent header", req.URL.Path)
		}
		if req.Header.Get("api-key") == "" {
			t.Errorf("request to %s missing api-key header", req.URL.Path)
		}
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write(body)
	}))
	t.Cleanup(srv.Close)
	return hevy.New("test-key", hevy.WithBaseURL(srv.URL))
}

// newErrorServer returns a server that always responds with the given status code.
func newErrorServer(t *testing.T, status int) *hevy.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"error":"test error"}`))
	}))
	t.Cleanup(srv.Close)
	return hevy.New("test-key", hevy.WithBaseURL(srv.URL))
}
