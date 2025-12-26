package api

import (
	"context"
	"net/url"

	"github.com/shashimalcse/asgardeo-cli/internal/models"
)

type organizationAPI struct {
	httpClient HTTPClient
}

type OrganizationAPI interface {
	List(ctx context.Context, filter string) (list *models.OrganizationList, err error)
	Create(ctx context.Context, organization map[string]interface{}) (err error)
	Update(ctx context.Context, id string, updates map[string]interface{}) (err error)
	UpdateMetadata(ctx context.Context, id string, metadata []models.OrganizationPatch) (err error)
	Delete(ctx context.Context, id string) (err error)
}

func NewOrganizationAPI(httpClient HTTPClient) OrganizationAPI {
	return &organizationAPI{httpClient: httpClient}
}

func (api *organizationAPI) List(ctx context.Context, filter string) (list *models.OrganizationList, err error) {
	params := url.Values{}
	if filter != "" {
		params.Add("filter", filter)
	}
	err = api.httpClient.Request(ctx, "GET", api.httpClient.URI("organizations"), WithParams(params), WithPayload(&list))
	return
}

func (api *organizationAPI) Create(ctx context.Context, organization map[string]interface{}) (err error) {
	err = api.httpClient.Request(ctx, "POST", api.httpClient.URI("organizations"), WithPayload(&organization))
	return
}

func (api *organizationAPI) Update(ctx context.Context, id string, updates map[string]interface{}) (err error) {
	err = api.httpClient.Request(ctx, "PUT", api.httpClient.URI("organizations", id), WithPayload(&updates))
	return
}

func (api *organizationAPI) UpdateMetadata(ctx context.Context, id string, metadata []models.OrganizationPatch) (err error) {
	err = api.httpClient.Request(ctx, "PATCH", api.httpClient.URI("organizations", id), WithPayload(&metadata))
	return
}

func (api *organizationAPI) Delete(ctx context.Context, id string) (err error) {
	err = api.httpClient.Request(ctx, "DELETE", api.httpClient.URI("organizations", id))
	return
}
