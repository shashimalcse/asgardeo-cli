package management

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (m *Management) NewRequest(ctx context.Context, method, uri string, payload interface{}) (*http.Request, error) {
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

func (m *Management) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	response, err := m.http.Do(req)
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

func (m *Management) Request(ctx context.Context, method, uri string, payload interface{}) error {
	request, err := m.NewRequest(ctx, method, uri, payload)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %w", err)
	}
	response, err := m.Do(request)
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

func (m *Management) URI(path ...string) string {
	baseURL := &url.URL{
		Scheme: m.url.Scheme,
		Host:   m.url.Host,
		Path:   m.basePath + "/",
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
