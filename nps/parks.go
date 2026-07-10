package nps

import (
	"context"
	"net/http"

	"github.com/freakytoad1/go-nps/api"
)

// ParksService handles communication with
// the /parks related methods of the NPS API.
type ParksService struct {
	client *Client
}

// ParkListOptions specifies the optional query parameters
// for the /parks endpoint.
type ParkListOptions struct {
	// ParkCode is a list of park codes (each 4-10 characters in length).
	ParkCode []string `url:"parkCode,omitempty,comma"`
	// StateCode is a list of 2 character state codes.
	StateCode []string `url:"stateCode,omitempty,comma"`
	// Q is a term to search on.
	Q string `url:"q,omitempty"`
	// Limit is the number of results to return per request. Default is 50.
	Limit int `url:"limit,omitempty"`
	// Start is the offset to start returning results from. Default is 0.
	Start int `url:"start,omitempty"`
	// Sort is a list of fields to sort the results by. Prefix a field with
	// a unary negative ("-") for descending order.
	Sort []string `url:"sort,omitempty,comma"`
}

// List retrieves data about national parks. A nil opts is allowed and
// returns the API default result set.
func (s *ParksService) List(ctx context.Context, opts *ParkListOptions) (*Parks, error) {
	req, err := s.client.apiClient.NewRequest(ctx, http.MethodGet, "parks", api.WithOptions(opts))
	if err != nil {
		return nil, err
	}

	parks := &Parks{}
	if err := s.client.apiClient.DoParse(req, parks); err != nil {
		return nil, err
	}

	return parks, nil
}

type Parks struct {
	Total string  `json:"total"`
	Limit string  `json:"limit"`
	Start string  `json:"start"`
	Data  []*Park `json:"data"`
}
type Activity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Topic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type PhoneNumber struct {
	PhoneNumber string `json:"phoneNumber"`
	Description string `json:"description"`
	Extension   string `json:"extension"`
	Type        string `json:"type"`
}
type EmailAddress struct {
	Description  string `json:"description"`
	EmailAddress string `json:"emailAddress"`
}
type Contacts struct {
	PhoneNumbers   []*PhoneNumber  `json:"phoneNumbers"`
	EmailAddresses []*EmailAddress `json:"emailAddresses"`
}
type Exception struct {
	ExceptionHours *Hours `json:"exceptionHours"`
	StartDate      string `json:"startDate"`
	Name           string `json:"name"`
	EndDate        string `json:"endDate"`
}
type Hours struct {
	Wednesday string `json:"wednesday"`
	Monday    string `json:"monday"`
	Thursday  string `json:"thursday"`
	Sunday    string `json:"sunday"`
	Tuesday   string `json:"tuesday"`
	Friday    string `json:"friday"`
	Saturday  string `json:"saturday"`
}
type OperatingHours struct {
	Exceptions    []*Exception `json:"exceptions"`
	Description   string       `json:"description"`
	StandardHours *Hours       `json:"standardHours"`
	Name          string       `json:"name"`
}
type Address struct {
	PostalCode            string `json:"postalCode"`
	City                  string `json:"city"`
	StateCode             string `json:"stateCode"`
	CountryCode           string `json:"countryCode"`
	ProvinceTerritoryCode string `json:"provinceTerritoryCode"`
	Line1                 string `json:"line1"`
	Line2                 string `json:"line2"`
	Line3                 string `json:"line3"`
	Type                  string `json:"type"`
}
type Image struct {
	Credit  string `json:"credit"`
	Title   string `json:"title"`
	AltText string `json:"altText"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
}

type EntranceFee struct {
	Cost        string `json:"cost"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type EntrancePass struct {
	Cost        string `json:"cost"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type Multimedia struct {
	Title string `json:"title"`
	ID    string `json:"id"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type Park struct {
	Activities     []*Activity       `json:"activities"`
	Addresses      []*Address        `json:"addresses"`
	Contacts       *Contacts         `json:"contacts"`
	Description    string            `json:"description"`
	Designation    string            `json:"designation"`
	DirectionsInfo string            `json:"directionsInfo"`
	DirectionsURL  string            `json:"directionsUrl"`
	EntranceFees   []*EntranceFee    `json:"entranceFees"`
	EntrancePasses []*EntrancePass   `json:"entrancePasses"`
	FullName       string            `json:"fullName"`
	ID             string            `json:"id"`
	Images         []*Image          `json:"images"`
	LatLong        string            `json:"latLong"`
	Latitude       string            `json:"latitude"`
	Longitude      string            `json:"longitude"`
	Multimedia     []*Multimedia     `json:"multimedia"`
	Name           string            `json:"name"`
	OperatingHours []*OperatingHours `json:"operatingHours"`
	ParkCode       string            `json:"parkCode"`
	RelevanceScore int               `json:"relevanceScore"`
	States         string            `json:"states"`
	Topics         []*Topic          `json:"topics"`
	URL            string            `json:"url"`
	WeatherInfo    string            `json:"weatherInfo"`
}
