package api

import (
	"context"

	"github.com/shashimalcse/asgardeo-cli/internal/models"
)

type applicationAPI struct {
	httpClient *httpClient
}

func NewApplicationAPI(httpClient *httpClient) *applicationAPI {
	return &applicationAPI{httpClient: httpClient}
}

func (a *applicationAPI) List(ctx context.Context) (list *models.ApplicationList, err error) {
	err = a.httpClient.Request(ctx, "GET", a.httpClient.URI("applications"), &list)
	return
}

func (m *applicationAPI) Create(ctx context.Context, application map[string]interface{}) (err error) {
	err = m.httpClient.Request(ctx, "POST", m.httpClient.URI("applications"), application)
	return
}

func (m *applicationAPI) Delete(ctx context.Context, id string) (err error) {
	err = m.httpClient.Request(ctx, "DELETE", m.httpClient.URI("applications", id), nil)
	return
}
