package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// apiErrorBody mirrors the JSON error envelope returned by the NPS API.
type apiErrorBody struct {
	Response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// APIError describes a non-2xx response from the NPS API. The parsed Code and
// Message come from the API's error body when available, while Hint provides a
// human-friendly explanation derived from the HTTP status code. Callers can use
// errors.As to inspect these fields.
type APIError struct {
	StatusCode int
	Status     string
	Code       string
	Message    string
	Hint       string
	Body       string
}

func (e *APIError) Error() string {
	msg := e.Status
	if e.Hint != "" {
		msg = fmt.Sprintf("%s: %s", msg, e.Hint)
	}

	switch {
	case e.Message != "" && e.Code != "":
		return fmt.Sprintf("%s (code %s: %s)", msg, e.Code, e.Message)
	case e.Message != "":
		return fmt.Sprintf("%s: %s", msg, e.Message)
	case e.Code != "":
		return fmt.Sprintf("%s (code %s)", msg, e.Code)
	case e.Body != "":
		return fmt.Sprintf("%s: %s", msg, e.Body)
	default:
		return msg
	}
}

// statusHint returns a human-friendly explanation for known error statuses.
func statusHint(code int) string {
	switch code {
	case http.StatusTooManyRequests:
		return "your API key is being temporarily blocked from making further requests; the block will automatically be lifted by waiting an hour"
	case http.StatusUnauthorized:
		return "not authorized for api endpoint; check that your API key is valid"
	case http.StatusBadRequest:
		return "request to api was not understood"
	case http.StatusNotFound:
		return "api endpoint was not found"
	default:
		return ""
	}
}

func newAPIError(resp *http.Response) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Hint:       statusHint(resp.StatusCode),
	}

	// read the body once so both JSON parsing and the raw fallback can use it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiErr
	}
	apiErr.Body = string(body)

	// attempt to extract the API's structured error code and message
	var parsed apiErrorBody
	if json.Unmarshal(body, &parsed) == nil {
		apiErr.Code = parsed.Response.Code
		apiErr.Message = parsed.Response.Message
	}

	return apiErr
}
