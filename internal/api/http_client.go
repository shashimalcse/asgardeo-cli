package api

import (
	"context"
	"net/http"
)

type HTTPClientAPI interface {
	NewRequest(ctx context.Context, method, uri string, payload interface{}) (*http.Request, error)

	Do(req *http.Request) (*http.Response, error)
}
