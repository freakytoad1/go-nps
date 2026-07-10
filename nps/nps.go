package nps

import (
	"github.com/freakytoad1/go-nps/api"
)

type Client struct {
	apiClient *api.Client

	// Services of the NPS API
	Parks *ParksService
}

func New(token string, opts ...api.ClientOption) (*Client, error) {
	apiClient, err := api.New(token, opts...)
	if err != nil {
		return nil, err
	}

	c := &Client{apiClient: apiClient}

	// Add all the services
	c.Parks = &ParksService{client: c}

	return c, nil

}

// RateLimit returns the rate limit information reported by the most recent
// response. The zero value is returned if no request has been made yet.
func (c *Client) RateLimit() api.RateLimitSnapshot {
	return c.apiClient.RateLimit()
}
