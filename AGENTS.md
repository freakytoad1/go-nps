# AGENTS.md

Guidance for AI coding agents working in the **go-nps** repository — a typed Go
client for the National Park Service (NPS) Data API. See [README.md](README.md)
for user-facing usage and examples.

## Required workflow for every code change

Run all of these after modifying any Go code, and fix anything they report
before considering the work done:

```sh
gofmt -l .            # must print nothing; run `gofmt -w .` to fix
go vet ./...
go build ./...
go test ./...
golangci-lint run ./...
```

- Keep the working tree `gofmt`-clean.
- `golangci-lint` must pass. (Note: the repo currently has pre-existing lint
  findings in `api/api.go` — do not introduce new ones, and prefer fixing any
  you touch.)
- There are no tests yet. When adding behavior, add tests — a local
  `httptest.Server` is preferred, and exercising the real NPS API to capture
  authentic response shapes is even better.

## Always verify against the official API

Do not assume request parameters, response fields, or behavior. Confirm against
the official documentation before adding or changing models and endpoints:

- [API documentation](https://www.nps.gov/subjects/developer/api-documentation.htm) — endpoints, parameters, models (the swagger JSON is linked there)
- [API guides](https://www.nps.gov/subjects/developer/guides.htm) — auth and rate limits
- [Get an API key](https://www.nps.gov/subjects/developer/get-started.htm)

Testing against the live API to see actual responses is encouraged over guessing.

## Architecture

Two packages with a deliberate split:

- **`api/`** — low-level transport. `api.Client` builds authenticated requests
  ([api/api.go](api/api.go)), applies request `Option`s, tracks rate-limit
  headers, and converts non-2xx responses into a structured `*api.APIError`
  ([api/error.go](api/error.go)). The `*http.Client` is an unexported field on
  purpose — do not re-embed it or expose raw HTTP methods.
- **`nps/`** — high-level, service-oriented client. `nps.Client` wraps an
  `api.Client` and exposes the API through services (e.g. `client.Parks` in
  [nps/parks.go](nps/parks.go)). Services are registered in
  [nps/nps.go](nps/nps.go).

## Conventions (project-specific)

- **New endpoints follow the existing service pattern.** Add a `XxxService` with
  a `client *Client` field, register it in `nps.New`, and give each method a
  typed `XxxListOptions` struct using `go-querystring` `url:"...,omitempty"`
  tags. Use `api.WithOptions(opts)` + `apiClient.NewRequest` + `apiClient.DoParse`.
  `Parks.List` is the reference implementation.
- **Response model types are singular** for slice elements (`Activity`, `Topic`,
  `PhoneNumber`), and nested slices use pointers (`[]*Image`). Stay consistent.
- **Most NPS fields are strings** even when numeric — match the API's actual
  JSON, do not "improve" types without verifying the real response.
- **Construction uses functional options** (`api.ClientOption`:
  `WithHTTPClient`, `WithBaseURL`), forwarded through `nps.New`.
- **Errors:** any status >= 400 becomes `*api.APIError`; callers inspect via
  `errors.As`. Keep error parsing centralized in `newAPIError`.

## Keep structure faithful & docs in sync

- Preserve the `api` (transport) vs `nps` (services) boundary. Do not merge the
  packages or move responsibilities across the boundary without reason.
- When you add/remove/change endpoints, options, or public types, update
  [README.md](README.md) (including the API coverage table) and this file in the
  same change.
</content>
