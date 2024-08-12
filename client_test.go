package hevy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransport(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		userAgent := req.Header.Get("User-Agent")
		assert.Equal(t, "go-hevy (https://github.com/swrm-io/go-hevy)", userAgent)

		apiKey := req.Header.Get("api-key")
		assert.Equal(t, "my-fake-api-key", apiKey)
	}))
	defer testServer.Close()

	client := NewClient("my-fake-api-key")

	client.ApiURL = testServer.URL
	client.client.Get("fake/url")
}