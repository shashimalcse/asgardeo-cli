package api

import (
	"context"

	"github.com/shashimalcse/asgardeo-cli/internal/models"
)

type applicationAPI struct {
	httpClient HTTPClient
}

type ApplicationAPI interface {
	List(ctx context.Context) (list *models.ApplicationList, err error)
	Create(ctx context.Context, application map[string]interface{}) (err error)
	Delete(ctx context.Context, id string) (err error)
}

func NewApplicationAPI(httpClient HTTPClient) ApplicationAPI {
	return &applicationAPI{httpClient: httpClient}
}

func (api *applicationAPI) List(ctx context.Context) (list *models.ApplicationList, err error) {
	err = api.httpClient.Request(ctx, "GET", api.httpClient.URI("applications"), WithPayload(&list))
	return
}

func (api *applicationAPI) Create(ctx context.Context, application map[string]interface{}) (err error) {
	err = api.httpClient.Request(ctx, "POST", api.httpClient.URI("applications"),
		WithPayload(application))
	return
}

func (api *applicationAPI) Delete(ctx context.Context, id string) (err error) {
	err = api.httpClient.Request(ctx, "DELETE", api.httpClient.URI("applications", id))
	return
}
