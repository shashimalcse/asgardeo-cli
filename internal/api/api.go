package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/shashimalcse/is-cli/internal/config"
)

type API struct {
	Application ApplicationAPI
	client      *httpClient
}

func NewAPI(cfg *config.Config, tenantDomain string) (*API, error) {
	client, err := NewHTTPClientAPI(cfg, tenantDomain)
	if err != nil {
		return nil, err
	}
	api := &API{
		client: client,
	}
	api.Application = NewApplicationAPI(api)
	return api, nil
}

func (a *API) Request(ctx context.Context, method, uri string, payload interface{}) error {
	request, err := a.NewRequest(ctx, method, uri, payload)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %w", err)
	}
	response, err := a.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send the request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		return newError(response)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %w", err)
	}
	if len(responseBody) > 0 && string(responseBody) != "{}" {
		if err = json.Unmarshal(responseBody, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal response payload: %w", err)
		}
	}
	return nil
}

func (a *API) NewRequest(ctx context.Context, method, uri string, payload interface{}) (*http.Request, error) {
	const nullBody = "null\n"
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return nil, fmt.Errorf("encoding request payload failed: %w", err)
		}
	}
	if body.String() == nullBody {
		body.Reset()
	}
	request, err := http.NewRequestWithContext(ctx, method, uri, &body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	return request, nil
}

func (a *API) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	response, err := a.client.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	return response, nil
}

func (a *API) URI(path ...string) string {
	baseURL := &url.URL{
		Scheme: a.client.url.Scheme,
		Host:   a.client.url.Host,
		Path:   a.client.basePath + "/",
	}
	const escapedForwardSlash = "%2F"
	var escapedPath []string
	for _, unescapedPath := range path {
		defaultPathEscaped := url.PathEscape(unescapedPath)
		escapedPath = append(
			escapedPath,
			strings.ReplaceAll(defaultPathEscaped, "/", escapedForwardSlash),
		)
	}
	return baseURL.String() + strings.Join(escapedPath, "/")
}

type APIError struct {
	StatusCode  int    `json:"statusCode"`
	Err         string `json:"error"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (m *APIError) Error() string {
	return fmt.Sprintf("%d %s: %s", m.StatusCode, m.Err, m.Message)
}

func (m *APIError) Status() int {
	return m.StatusCode
}

func newError(response *http.Response) error {
	apiError := &APIError{}

	if err := json.NewDecoder(response.Body).Decode(apiError); err != nil {
		return &APIError{
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
