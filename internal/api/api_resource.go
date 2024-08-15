package api

import (
	"context"

	"github.com/shashimalcse/asgardeo-cli/internal/models"
)

type apiResourceAPI struct {
	httpClient *httpClient
}

func NewApiResourceAPI(httpClient *httpClient) *apiResourceAPI {
	return &apiResourceAPI{httpClient: httpClient}
}

func (a *apiResourceAPI) List(ctx context.Context, apiType string) (list *models.APIResourceList, err error) {
	err = a.httpClient.Request(ctx, "GET", a.httpClient.URI("api-resources"), &list)
	return
}

func (m *apiResourceAPI) Create(ctx context.Context, apiResource map[string]interface{}) (err error) {
	err = m.httpClient.Request(ctx, "POST", m.httpClient.URI("api-resources"), apiResource)
	return
}

func (m *apiResourceAPI) Delete(ctx context.Context, id string) (err error) {
	err = m.httpClient.Request(ctx, "DELETE", m.httpClient.URI("api-resources", id), nil)
	return
}
