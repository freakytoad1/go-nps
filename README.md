# go-nps

A National Park Service (NPS) API client for Go.

`go-nps` provides a small, typed client for the [NPS Data API](https://www.nps.gov/subjects/developer/index.htm).
It handles authentication, request building, rate-limit tracking, and structured
error responses so you can focus on consuming park data.

## Installation

```sh
go get github.com/freakytoad1/go-nps
```

Requires Go 1.26+.

## Getting an API key

All requests require an API key. Register for a free key on the
[NPS Get Started page](https://www.nps.gov/subjects/developer/get-started.htm).
Keep your key private — do not commit it to source control.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/freakytoad1/go-nps/nps"
)

func main() {
	client, err := nps.New(os.Getenv("NPS_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	parks, err := client.Parks.List(context.Background(), &nps.ParkListOptions{
		StateCode: []string{"WY"},
		Limit:     5,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, park := range parks.Data {
		fmt.Printf("%s (%s)\n", park.FullName, park.ParkCode)
	}
}
```

## Project structure

The module is split into two packages:

```
go-nps/
├── api/          # Low-level HTTP client
│   ├── api.go    # Client, request building, options, rate limiting
│   └── error.go  # Structured APIError type
└── nps/          # High-level, service-oriented client
    ├── nps.go    # Client and service registration
    └── parks.go  # Parks service, options, and response models
```

- **`api`** is the transport layer. It builds authenticated requests, applies
  request options, sends them, tracks rate-limit headers, and converts non-2xx
  responses into a structured `*api.APIError`. Most users do not interact with
  this package directly.
- **`nps`** is the package you use. `nps.Client` wraps an `api.Client` and
  exposes the API through services (e.g. `client.Parks`). Each service has typed
  methods and option structs that map to the API's query parameters.

### Construction options

`nps.New` (and the underlying `api.New`) accept functional options:

```go
client, err := nps.New(token,
	api.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	api.WithBaseURL("https://developer.nps.gov/api/v1/"),
)
```

## API coverage

The NPS API exposes roughly 25 endpoints. This client currently implements:

| Endpoint  | Service / Method        | Status |
| --------- | ----------------------- | ------ |
| `/parks`  | `client.Parks.List`     | ✅      |

All other endpoints (`/alerts`, `/campgrounds`, `/events`, `/activities`,
`/thingstodo`, `/visitorcenters`, `/webcams`, etc.) are not yet implemented.
The service pattern in `nps/nps.go` and `nps/parks.go` is designed to make
adding them straightforward.

## Working with parks

`Parks.List` accepts a `*ParkListOptions` to filter and paginate results. A
`nil` value returns the API's default result set.

```go
opts := &nps.ParkListOptions{
	ParkCode:  []string{"yell", "grca"}, // comma-delimited park codes
	StateCode: []string{"WY", "AZ"},     // comma-delimited state codes
	Q:         "canyon",                 // free-text search
	Limit:     20,                       // results per page (default 50)
	Start:     0,                        // pagination offset
	Sort:      []string{"fullName"},     // prefix with "-" for descending
}

parks, err := client.Parks.List(context.Background(), opts)
```

The response includes pagination metadata and the matched parks:

```go
fmt.Printf("total=%s limit=%s start=%s\n", parks.Total, parks.Limit, parks.Start)
for _, park := range parks.Data {
	fmt.Println(park.FullName, park.States)
}
```

## Rate limiting

The NPS API allows a default of **1,000 requests per hour per key** and returns
`X-RateLimit-Limit` / `X-RateLimit-Remaining` headers on every response. The
client tracks these automatically. Read the most recent values with
`RateLimit()`:

```go
snap := client.RateLimit()
fmt.Printf("%d of %d requests remaining (as of %s)\n",
	snap.Remaining, snap.Limit, snap.LastUpdated)
```

Exceeding the limit results in an HTTP 429; the block lifts automatically after
an hour.

## Error handling

Any response with a status of 400 or greater is returned as an `*api.APIError`,
which carries the HTTP status, a human-friendly hint, and the API's structured
error code/message when available. Use `errors.As` to inspect it:

```go
parks, err := client.Parks.List(context.Background(), opts)
if err != nil {
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusTooManyRequests:
			// back off and retry later
		case http.StatusUnauthorized:
			// check that your API key is valid
		default:
			log.Printf("API error: %v", apiErr)
		}
	}
	return
}
```

## Reference

- [NPS API documentation](https://www.nps.gov/subjects/developer/api-documentation.htm)
- [NPS API guides (auth, rate limits)](https://www.nps.gov/subjects/developer/guides.htm)
- [NPS Data API overview](https://www.nps.gov/subjects/digital/nps-data-api.htm)
- [NPS on GitHub](https://github.com/nationalparkservice)
- [NPS API samples](https://github.com/nationalparkservice/nps-api-samples)