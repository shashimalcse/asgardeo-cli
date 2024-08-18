package api

import (
	"github.com/shashimalcse/asgardeo-cli/internal/config"
	"go.uber.org/zap"
)

type API struct {
	Application *applicationAPI
	APIResource *apiResourceAPI
	httpClient  *httpClient
}

func NewAPI(cfg *config.Config, tenantDomain string, logger *zap.Logger) (*API, error) {
	httpClient, err := NewHTTPClientAPI(cfg, tenantDomain, logger)
	if err != nil {
		return nil, err
	}
	api := &API{
		httpClient:  httpClient,
		Application: NewApplicationAPI(httpClient),
		APIResource: NewApiResourceAPI(httpClient),
	}
	return api, nil
}
