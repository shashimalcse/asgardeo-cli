package api

import "github.com/shashimalcse/is-cli/internal/management"

type API struct {
	Application ApplicationAPI

	HTTPClient HTTPClientAPI
}

func NewAPI(m *management.Management) *API {
	return &API{
		Application: m.Application,
		HTTPClient:  m,
	}
}
