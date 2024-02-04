package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errPerformingRequest          = errors.New("error performing request")
	errReadingResponseBody        = errors.New("error reading response body")
	errUnmarshallingResponse      = errors.New("error unmarshalling notifications summary")
	errCreatingRequest            = errors.New("error creating request")
	errNilResponse                = errors.New("error nil response")
	errUnmarshallingErrorResponse = errors.New("error unmarshalling error response")
)

type serviceErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func errPolicyFunc(response *http.Response) error {
	if response.StatusCode < http.StatusBadRequest {
		return nil
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("%w: cannot read error response body: %v", errReadingResponseBody, err)
	}

	var errResponse serviceErrorResponse
	err = json.Unmarshal(responseBody, &errResponse)
	if err != nil {
		return fmt.Errorf("%w: %w", errUnmarshallingErrorResponse, err)
	}

	return fmt.Errorf("%s", errResponse.Message)
}
