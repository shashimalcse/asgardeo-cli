package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode  int    `json:"statusCode"`
	Err         string `json:"error"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (m *Error) Error() string {
	return fmt.Sprintf("%d %s: %s", m.StatusCode, m.Err, m.Message)
}

func (m *Error) Status() int {
	return m.StatusCode
}

func newError(response *http.Response) error {
	apiError := &Error{}

	if err := json.NewDecoder(response.Body).Decode(apiError); err != nil {
		return &Error{
			StatusCode: response.StatusCode,
			Err:        http.StatusText(response.StatusCode),
			Message:    fmt.Errorf("failed to decode json error response payload: %w", err).Error(),
		}
	}

	if apiError.Status() == 0 {
		apiError.StatusCode = response.StatusCode
		apiError.Err = http.StatusText(response.StatusCode)
	}

	return apiError
}
