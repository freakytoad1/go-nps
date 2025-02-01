package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseUrl = "https://developer.nps.gov/api/v1/"

	headerApiKey             = "X-Api-Key"
	headerRateLimit          = "X-RateLimit-Limit"
	headerRateLimitRemaining = "X-RateLimit-Remaining"
)

type Client struct {
	*http.Client
	url   *url.URL
	token string // api key

	rateLimit *RateLimit
}

type RateLimit struct {
	limit          string
	limitRemaining string
	lastUpdated    time.Time

	*sync.RWMutex
}

func (c *Client) String() string {
	return fmt.Sprintf("url: %s", c.url.String())
}

// New creates a new api client to perform http requests.
func New(token string) (*Client, error) {
	baseUrl, _ := url.Parse(defaultBaseUrl)

	if token == "" {
		return nil, errors.New("api key token was empty")
	}

	return &Client{
		Client: &http.Client{},
		url:    baseUrl,
		token:  token,
	}, nil
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

// WithJSONBody allows a new request to include a JSON body.
func WithJSONBody(b any) Option {
	return func(r *http.Request) error {
		// convert struct to bytes
		body, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("unable to marshal request body: %w", err)
		}

		// add to request
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		return nil
	}
}

// WithQuery adds the key value pairs as a raw query to the request.
func WithQuery(queryMap map[string]string) Option {
	return func(r *http.Request) error {
		params := url.Values{}

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
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return nil
		}

		params, err := query.Values(opts)
		if err != nil {
			return err
		}

		r.URL.RawQuery = params.Encode()
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

	defer resp.Body.Close()

	// unmarshal response
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return fmt.Errorf("unable to decode json response: %w", err)
	}

	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// send it
	resp, err := c.Client.Do(req)
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

	c.rateLimit.limit = resp.Header.Get(headerRateLimit)
	c.rateLimit.limitRemaining = resp.Header.Get(headerRateLimitRemaining)
	c.rateLimit.lastUpdated = time.Now()
}

// validateResponse determines if nps api returned an error.
// 200 = good
// 429 = going above the rate limit
// 400 = bad request
// 404 = api endpoint not found
func (c *Client) validateResponse(resp *http.Response) error {
	// update the rate limit info
	c.updateRateLimitAmounts(resp)

	switch {
	case resp.StatusCode < 400:
		return nil
	case resp.StatusCode == http.StatusTooManyRequests:
		return fmt.Errorf("your API key is being temporarily blocked from making further requests. The block will automatically be lifted by waiting an hour: %s", resp.Status)
	case resp.StatusCode == http.StatusUnauthorized:
		return fmt.Errorf("not authorized for api endpoint: %s", resp.Status)
	case resp.StatusCode == http.StatusBadRequest:
		return fmt.Errorf("request to api was not understood: %s", resp.Status)
	case resp.StatusCode == http.StatusNotFound:
		return fmt.Errorf("api endpoint was not found: %s", resp.Status)
	}

	return newAPIError(resp)
}
