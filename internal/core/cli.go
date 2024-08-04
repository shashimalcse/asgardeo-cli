package core

import (
	"context"
	"fmt"

	"github.com/shashimalcse/is-cli/internal/api"
	"github.com/shashimalcse/is-cli/internal/config"
	"go.uber.org/zap"
)

// CLI represents the main CLI application structure
type CLI struct {
	Config *config.Config
	Logger *zap.Logger
	Tenant string
	API    *api.API
}

// NewCLI creates a new CLI instance
func NewCLI(cfg *config.Config, logger *zap.Logger) *CLI {
	return &CLI{
		Config: cfg,
		Logger: logger,
	}
}

// SetupWithAuthentication sets up the CLI with authentication
func (c *CLI) SetupWithAuthentication(ctx context.Context) error {
	if err := c.Config.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	if c.Tenant == "" {
		c.Tenant = c.Config.DefaultTenant
	}
	if err := c.checkAndRefreshAuth(ctx); err != nil {
		return fmt.Errorf("authentication check failed: %w", err)
	}
	api, err := api.NewAPI(c.Config, c.Tenant)
	if err != nil {
		return fmt.Errorf("failed to initialize API client: %w", err)
	}
	c.API = api
	return nil
}

func (c *CLI) checkAndRefreshAuth(ctx context.Context) error {
	tenant, err := c.Config.GetTenant(c.Tenant)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}
	err = tenant.CheckAuthenticationStatus()
	switch err {
	case nil:
		return nil
	case config.ErrInvalidToken:
		c.Logger.Info("Token is invalid, attempting to refresh")
		if err := c.refreshToken(ctx); err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unexpected error checking auth status: %w", err)
	}
}

func (c *CLI) refreshToken(ctx context.Context) error {
	// Implement token refresh logic here
	// This is a placeholder and should be replaced with actual refresh logic
	c.Logger.Info("Refreshing token")
	// newToken := "refreshed-token"
	// tenant.SetAccessToken(newToken)
	// return c.Config.AddTenant(*tenant)
	return fmt.Errorf("token refresh not implemented")
}
