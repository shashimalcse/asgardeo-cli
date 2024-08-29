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
	"time"

	"github.com/shashimalcse/asgardeo-cli/internal/config"
	"go.uber.org/zap"
)

type httpClient struct {
	client   *http.Client
	baseUrl  *url.URL
	basepath string
	token    string
	logger   *zap.Logger
}

type HTTPClient interface {
	Request(ctx context.Context, method, uri string, opts ...RequestOption) error
	Do(req *http.Request) (*http.Response, error)
	URI(path ...string) string
}

func NewHTTPClientAPI(cfg *config.Config, tenantDomain string, logger *zap.Logger) (HTTPClient, error) {
	tenant, err := cfg.GetTenant(tenantDomain)
	if err != nil {
		logger.Error("failed to get tenant while creating http client", zap.Error(err))
		return nil, err
	}
	basepath := "t/" + tenant.Name + "/api/server/v1"
	u, err := url.Parse("https://api.asgardeo.io/")
	if err != nil {
		logger.Error("failed to parse base URL while creating http client", zap.Error(err))
		return nil, err
	}
	return &httpClient{client: &http.Client{Timeout: 30 * time.Second}, basepath: basepath, baseUrl: u, token: tenant.GetAccessToken(), logger: logger}, nil
}

func (c *httpClient) Request(ctx context.Context, method, uri string, opts ...RequestOption) error {
	options := &requestOptions{}
	for _, opt := range opts {
		opt(options)
	}
	request, err := c.newRequest(ctx, method, uri, options.params, options.payload)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %w", err)
	}
	response, err := c.Do(request)
	if err != nil {
		c.logger.Error("failed to send the request with http client", zap.String("method", method), zap.String("uri", uri), zap.Error(err))
		return fmt.Errorf("failed to send the request: %w", err)
	}
	defer func() {
		if cErr := response.Body.Close(); cErr != nil {
			err = fmt.Errorf("failed to close response body: %w", cErr)
		}
	}()
	if response.StatusCode >= http.StatusBadRequest {
		c.logger.Error("received an error response from the server", zap.String("method", method), zap.String("uri", uri), zap.Int("status_code", response.StatusCode))
		return newError(response)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %w", err)
	}
	if len(responseBody) > 0 && string(responseBody) != "{}" {
		if err = json.Unmarshal(responseBody, &options.payload); err != nil {
			return fmt.Errorf("failed to unmarshal response payload: %w", err)
		}
	}
	return nil
}

func (c *httpClient) newRequest(ctx context.Context, method, uri string, params url.Values, payload interface{}) (*http.Request, error) {
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
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	// Merge existing query parameters with new ones
	query := parsedURL.Query()
	for key, values := range params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	parsedURL.RawQuery = query.Encode()
	request, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), &body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	return request, nil
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	req.Header.Set("Authorization", "Bearer "+c.token)
	response, err := c.client.Do(req)
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

func (c *httpClient) URI(path ...string) string {
	baseURL := &url.URL{
		Scheme: c.baseUrl.Scheme,
		Host:   c.baseUrl.Host,
		Path:   c.basepath + "/",
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

type RequestOption func(*requestOptions)

type requestOptions struct {
	params  url.Values
	payload interface{}
}

func WithParams(params url.Values) RequestOption {
	return func(ro *requestOptions) {
		ro.params = params
	}
}

func WithPayload(payload interface{}) RequestOption {
	return func(ro *requestOptions) {
		ro.payload = payload
	}
}
