package api

import (
	"context"
	"net/url"

	"github.com/shashimalcse/asgardeo-cli/internal/models"
)

type apiResourceAPI struct {
	httpClient *httpClient
}

func NewApiResourceAPI(httpClient *httpClient) *apiResourceAPI {
	return &apiResourceAPI{httpClient: httpClient}
}

func (a *apiResourceAPI) List(ctx context.Context, apiType string) (list *models.APIResourceList, err error) {
	params := url.Values{}
	params.Add("attributes", "properties")
	params.Add("filter", "type eq "+apiType)
	err = a.httpClient.Request(ctx, "GET", a.httpClient.URI("api-resources"), WithParams(params), WithPayload(&list))
	return
}

func (a *apiResourceAPI) Get(ctx context.Context, id string) (apiResource *models.APIResource, err error) {
	err = a.httpClient.Request(ctx, "GET", a.httpClient.URI("api-resources/"+id), WithPayload(&apiResource))
	return
}

func (a *apiResourceAPI) Create(ctx context.Context, apiResource map[string]interface{}) (err error) {
	err = a.httpClient.Request(ctx, "POST", a.httpClient.URI("api-resources"), WithPayload(&apiResource))
	return
}

func (a *apiResourceAPI) Delete(ctx context.Context, id string) (err error) {
	err = a.httpClient.Request(ctx, "DELETE", a.httpClient.URI("api-resources", id))
	return
}
