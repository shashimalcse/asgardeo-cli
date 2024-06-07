package core

import (
	"context"
	"fmt"

	"github.com/shashimalcse/is-cli/internal/api"
	"github.com/shashimalcse/is-cli/internal/config"
	"github.com/shashimalcse/is-cli/internal/management"
	"go.uber.org/zap"
)

type CLI struct {
	Config config.Config
	Logger zap.Logger
	Tenant string
	API    *api.API
}

func (c *CLI) SetupWithAuthentication(ctx context.Context) error {

	if err := c.Config.Validate(); err != nil {
		return err
	}

	if c.Tenant == "" {
		c.Tenant = c.Config.DefaultTenant
	}

	tenant, err := c.Config.GetTenant(c.Tenant)
	if err != nil {
		return err
	}

	// Check authentication status.
	err = tenant.CheckAuthenticationStatus()
	switch err {
	case config.ErrInvalidToken:
		return fmt.Errorf("Invalid token. please login")
	}
	client, err := initializeManagementClient(tenant.Name, tenant.GetAccessToken())
	if err != nil {
		return err
	}
	c.API = api.NewAPI(client)
	return nil
}

func initializeManagementClient(tenantDomain string, accessToken string) (*management.Management, error) {
	client, err := management.New(tenantDomain, accessToken)

	return client, err
}
