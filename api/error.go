package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type apiError struct {
	Response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("Code: %s Message: %s", e.Response.Code, e.Response.Message)
}

func newAPIError(resp *http.Response) error {
	// decode response
	apiErrResp := &apiError{}
	if err := json.NewDecoder(resp.Body).Decode(apiErrResp); err != nil {
		// attempt to just turn the body into a string
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			// might be nothing to decode
			return fmt.Errorf("%s: unable to decode error response: %w", resp.Status, err)
		}

		return fmt.Errorf("%s: %s", resp.Status, string(b))
	}

	return apiErrResp
}
