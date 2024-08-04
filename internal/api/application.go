package api

import (
	"context"

	"github.com/shashimalcse/is-cli/internal/models"
)

type ApplicationAPI interface {
	List(ctx context.Context) (c *models.ApplicationList, err error)
	Create(ctx context.Context, application map[string]interface{}) (err error)
	Delete(ctx context.Context, id string) (err error)
}

type applicationAPI struct {
	api *API
}

func NewApplicationAPI(api *API) ApplicationAPI {
	return &applicationAPI{api: api}
}

func (a *applicationAPI) List(ctx context.Context) (list *models.ApplicationList, err error) {
	err = a.api.Request(ctx, "GET", a.api.URI("applications"), &list)
	return
}

func (m *applicationAPI) Create(ctx context.Context, application map[string]interface{}) (err error) {
	err = m.api.Request(ctx, "POST", m.api.URI("applications"), application)
	return
}

func (m *applicationAPI) Delete(ctx context.Context, id string) (err error) {
	err = m.api.Request(ctx, "DELETE", m.api.URI("applications", id), nil)
	return
}
