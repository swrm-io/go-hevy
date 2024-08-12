package hevy

import (
	"net/http"
	"time"
)

// Client for working with the Hevy API
type Client struct {
	client     http.Client
	ApiURL     string // Base API URL
	ApiVersion string // API endpoint version for requests
}

// NewClient returns a default client instance
func NewClient(apiKey string) *Client {
	httpClient := http.Client{
		Transport: apiTransport{apiKey: apiKey, agent: "go-hevy (https://github.com/swrm-io/go-hevy)"},
		Timeout:   10 * time.Second,
	}

	return &Client{
		client:     httpClient,
		ApiURL:     "https://api.hevyapp.com",
		ApiVersion: "v1",
	}
}