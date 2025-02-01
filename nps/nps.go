package nps

import (
	"github.com/freakytoad1/go-nps/api"
)

type Client struct {
	apiClient *api.Client

	// Services of the NPS API
	Parks *ParksService
}

func New(token string) (*Client, error) {
	apiClient, err := api.New(token)
	if err != nil {
		return nil, err
	}

	c := &Client{apiClient: apiClient}

	// Add all the services
	c.Parks = &ParksService{client: c}

	return c, nil

}
