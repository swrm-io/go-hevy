// Package hevy provides a Go client for the Hevy workout tracking API.
// See https://api.hevyapp.com/docs for the full API reference.
// A Hevy Pro subscription is required to use the API.
package hevy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://api.hevyapp.com"
	apiKeyHeader   = "api-key"
)

// Sentinel errors returned by the client. Use errors.Is to check them.
var (
	// ErrNotFound is returned when the API responds with 404.
	ErrNotFound = errors.New("hevy: not found")
	// ErrUnauthorized is returned when the API responds with 401 or 403.
	ErrUnauthorized = errors.New("hevy: unauthorized")
	// ErrConflict is returned when the API responds with 409 (e.g. duplicate body measurement date).
	ErrConflict = errors.New("hevy: conflict")
	// ErrBadRequest is returned when the API responds with 400.
	ErrBadRequest = errors.New("hevy: bad request")
	// ErrRoutineLimitExceeded is returned when the routine limit is hit (403 on POST /v1/routines).
	ErrRoutineLimitExceeded = errors.New("hevy: routine limit exceeded")
	// ErrExerciseLimitExceeded is returned when the custom exercise limit is hit (403 on POST /v1/exercise_templates).
	ErrExerciseLimitExceeded = errors.New("hevy: exercise template limit exceeded")
	// ErrNoMorePages is returned by List when the requested page exceeds the total page count.
	ErrNoMorePages = errors.New("hevy: no more pages")
	// ErrInvalidPageSize is returned when pageSize exceeds the maximum allowed for an endpoint.
	ErrInvalidPageSize = errors.New("hevy: invalid page size")
)

// APIError is returned when the API responds with a non-2xx status code.
// It wraps a sentinel error for use with errors.Is.
type APIError struct {
	StatusCode int
	Body       string
	sentinel   error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("hevy: API error %d: %s", e.StatusCode, e.Body)
}

func (e *APIError) Is(target error) bool {
	return e.sentinel != nil && errors.Is(e.sentinel, target)
}

func (e *APIError) Unwrap() error {
	return e.sentinel
}

// core is the shared HTTP plumbing used by all service types.
type core struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Client is the Hevy API client. Access resources via its service fields.
type Client struct {
	// Workouts provides access to workout endpoints.
	Workouts *WorkoutsService
	// Routines provides access to routine endpoints.
	Routines *RoutinesService
	// RoutineFolders provides access to routine folder endpoints.
	RoutineFolders *RoutineFoldersService
	// ExerciseTemplates provides access to exercise template endpoints.
	ExerciseTemplates *ExerciseTemplatesService
	// ExerciseHistory provides access to exercise history endpoints.
	ExerciseHistory *ExerciseHistoryService
	// BodyMeasurements provides access to body measurement endpoints.
	BodyMeasurements *BodyMeasurementsService
	// User provides access to user info endpoints.
	User *UserService
}

// Option is a functional option for configuring a Client.
type Option func(*core)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(u string) Option {
	return func(c *core) {
		c.baseURL = u
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *core) {
		c.httpClient = hc
	}
}

// New creates a new Hevy API client authenticated with the given API key.
func New(apiKey string, opts ...Option) *Client {
	co := &core{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, o := range opts {
		o(co)
	}
	return &Client{
		Workouts:          &WorkoutsService{co},
		Routines:          &RoutinesService{co},
		RoutineFolders:    &RoutineFoldersService{co},
		ExerciseTemplates: &ExerciseTemplatesService{co},
		ExerciseHistory:   &ExerciseHistoryService{co},
		BodyMeasurements:  &BodyMeasurementsService{co},
		User:              &UserService{co},
	}
}

func (c *core) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	u := c.baseURL + path
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("hevy: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, c.apiKey)
	req.Header.Set("User-Agent", "github.com/swrm-io/go-hevy")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *core) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("hevy: do request: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("hevy: read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp.StatusCode, string(b), req)
	}

	if out != nil && len(b) > 0 {
		if err := json.Unmarshal(b, out); err != nil {
			return fmt.Errorf("hevy: decode response: %w", err)
		}
	}
	return nil
}

// newAPIError maps HTTP status codes to sentinel errors.
// The path is used to disambiguate 403 responses (routine vs exercise limit).
func newAPIError(statusCode int, body string, req *http.Request) *APIError {
	var sentinel error
	switch statusCode {
	case http.StatusBadRequest:
		sentinel = ErrBadRequest
	case http.StatusUnauthorized:
		sentinel = ErrUnauthorized
	case http.StatusForbidden:
		// 403 means "limit exceeded" on specific POST endpoints, otherwise unauthorized.
		path := req.URL.Path
		switch path {
		case "/v1/routines":
			sentinel = ErrRoutineLimitExceeded
		case "/v1/exercise_templates":
			sentinel = ErrExerciseLimitExceeded
		default:
			sentinel = ErrUnauthorized
		}
	case http.StatusNotFound:
		sentinel = ErrNotFound
	case http.StatusConflict:
		sentinel = ErrConflict
	}
	return &APIError{StatusCode: statusCode, Body: body, sentinel: sentinel}
}

func (c *core) get(ctx context.Context, path string, query url.Values, out any) error {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *core) post(ctx context.Context, path string, body, out any) error {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *core) put(ctx context.Context, path string, body, out any) error {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

// validatePageSize returns ErrInvalidPageSize if pageSize exceeds max.
func validatePageSize(pageSize, max int) error {
	if pageSize > max {
		return fmt.Errorf("%w: maximum is %d, got %d", ErrInvalidPageSize, max, pageSize)
	}
	return nil
}

func pageQuery(page, pageSize int) url.Values {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	return q
}
