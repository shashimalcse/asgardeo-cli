package api

import (
	"context"
	"net/url"

	"github.com/shashimalcse/asgardeo-cli/internal/models"
)

type apiResourceAPI struct {
	httpClient HTTPClient
}

type ResourceAPI interface {
	List(ctx context.Context, apiType string) (list *models.APIResourceList, err error)
	Get(ctx context.Context, id string) (apiResource *models.APIResource, err error)
	Create(ctx context.Context, apiResource map[string]interface{}) (err error)
	Delete(ctx context.Context, id string) (err error)
}

func NewApiResourceAPI(httpClient HTTPClient) ResourceAPI {
	return &apiResourceAPI{httpClient: httpClient}
}

func (api *apiResourceAPI) List(ctx context.Context, apiType string) (list *models.APIResourceList, err error) {
	params := url.Values{}
	params.Add("attributes", "properties")
	params.Add("filter", "type eq "+apiType)
	err = api.httpClient.Request(ctx, "GET", api.httpClient.URI("api-resources"), WithParams(params), WithPayload(&list))
	return
}

func (api *apiResourceAPI) Get(ctx context.Context, id string) (apiResource *models.APIResource, err error) {
	err = api.httpClient.Request(ctx, "GET", api.httpClient.URI("api-resources", id), WithPayload(&apiResource))
	return
}

func (api *apiResourceAPI) Create(ctx context.Context, apiResource map[string]interface{}) (err error) {
	err = api.httpClient.Request(ctx, "POST", api.httpClient.URI("api-resources"), WithPayload(&apiResource))
	return
}

func (api *apiResourceAPI) Delete(ctx context.Context, id string) (err error) {
	err = api.httpClient.Request(ctx, "DELETE", api.httpClient.URI("api-resources", id))
	return
}
