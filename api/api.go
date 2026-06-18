package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseUrl = "https://developer.nps.gov/api/v1/"
	userAgent      = "go-nps"

	headerApiKey             = "X-Api-Key"
	headerRateLimit          = "X-RateLimit-Limit"
	headerRateLimitRemaining = "X-RateLimit-Remaining"
)

type Client struct {
	httpClient *http.Client
	url        *url.URL
	token      string // api key

	rateLimit *RateLimit
}

type RateLimit struct {
	limit          int
	limitRemaining int
	lastUpdated    time.Time

	*sync.RWMutex
}

// RateLimitSnapshot is a point-in-time view of the API rate limit
// reported by the most recent response.
type RateLimitSnapshot struct {
	Limit       int
	Remaining   int
	LastUpdated time.Time
}

func (c *Client) String() string {
	return fmt.Sprintf("url: %s", c.url.String())
}

// RateLimit returns the rate limit information reported by the most recent
// response. The zero value is returned if no request has been made yet.
func (c *Client) RateLimit() RateLimitSnapshot {
	c.rateLimit.RLock()
	defer c.rateLimit.RUnlock()

	return RateLimitSnapshot{
		Limit:       c.rateLimit.limit,
		Remaining:   c.rateLimit.limitRemaining,
		LastUpdated: c.rateLimit.lastUpdated,
	}
}

// ClientOption configures a Client during construction.
type ClientOption func(*Client) error

// WithHTTPClient sets a custom *http.Client, useful for configuring timeouts,
// transports, or supplying a mock in tests.
func WithHTTPClient(h *http.Client) ClientOption {
	return func(c *Client) error {
		if h == nil {
			return errors.New("http client was nil")
		}
		c.httpClient = h
		return nil
	}
}

// WithBaseURL overrides the default NPS API base URL.
func WithBaseURL(raw string) ClientOption {
	return func(c *Client) error {
		u, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("invalid base url: %w", err)
		}
		c.url = u
		return nil
	}
}

// New creates a new api client to perform http requests.
func New(token string, opts ...ClientOption) (*Client, error) {
	if token == "" {
		return nil, errors.New("api key token was empty")
	}

	baseUrl, err := url.Parse(defaultBaseUrl)
	if err != nil {
		return nil, err
	}

	c := &Client{
		httpClient: &http.Client{},
		url:        baseUrl,
		token:      token,
		rateLimit: &RateLimit{
			RWMutex: &sync.RWMutex{},
		},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// NewRequest wraps http.NewRequestWithContext but adds the ability to supply options as needed.
func (c *Client) NewRequest(ctx context.Context, method string, path string, options ...Option) (*http.Request, error) {
	// escape path and add to base
	u, err := url.JoinPath(c.url.String(), path)
	if err != nil {
		return nil, err
	}

	// craft base request
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set(headerApiKey, c.token)

	// perform options on request
	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Option Pattern
type Option func(*http.Request) error

// WithQuery adds the key value pairs as a raw query to the request.
func WithQuery(queryMap map[string]string) Option {
	return func(r *http.Request) error {
		// start from any query params already on the request
		params := r.URL.Query()

		// add each map key / value as a query param
		for key, value := range queryMap {
			params.Add(key, value)
		}

		// encode the values and add to url
		r.URL.RawQuery = params.Encode()
		return nil
	}
}

// WithOptions adds the parameters in opts as URL query parameters to the request.
// must be a struct whose fields may contain "url" tags.
func WithOptions(opts any) Option {
	return func(r *http.Request) error {
		v := reflect.ValueOf(opts)
		if v.Kind() == reflect.Pointer && v.IsNil() {
			return nil
		}

		params, err := query.Values(opts)
		if err != nil {
			return err
		}

		// merge into any query params already on the request
		existing := r.URL.Query()
		for key, values := range params {
			for _, value := range values {
				existing.Add(key, value)
			}
		}

		r.URL.RawQuery = existing.Encode()
		return nil
	}
}

// DoParse performs standard Do on the request and unmarshals the response body into ans.
// The response is also validated against known NPS error codes and will provide more insight into why a request may have failed.
func (c *Client) DoParse(req *http.Request, v any) error {
	// send it
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close() //nolint:errcheck // normal pattern to ignore this error

	// unmarshal response
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return fmt.Errorf("unable to decode json response: %w", err)
	}

	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// send it
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	// validate response first
	if err := c.validateResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) updateRateLimitAmounts(resp *http.Response) {
	c.rateLimit.Lock()
	defer c.rateLimit.Unlock()

	if limit, err := strconv.Atoi(resp.Header.Get(headerRateLimit)); err == nil {
		c.rateLimit.limit = limit
	}
	if remaining, err := strconv.Atoi(resp.Header.Get(headerRateLimitRemaining)); err == nil {
		c.rateLimit.limitRemaining = remaining
	}
	c.rateLimit.lastUpdated = time.Now()
}

// validateResponse determines if the NPS API returned an error. Any response
// with a status code of 400 or greater is converted into an *APIError, which
// includes the API's structured error body when available.
func (c *Client) validateResponse(resp *http.Response) error {
	// update the rate limit info
	c.updateRateLimitAmounts(resp)

	if resp.StatusCode < 400 {
		return nil
	}

	return newAPIError(resp)
}
