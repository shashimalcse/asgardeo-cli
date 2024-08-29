package core

import (
	"errors"
	"fmt"

	"github.com/shashimalcse/asgardeo-cli/internal/api"
	"github.com/shashimalcse/asgardeo-cli/internal/config"
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
func (c *CLI) SetupWithAuthentication() error {
	if err := c.Config.Validate(); err != nil {
		return err
	}
	if c.Tenant == "" {
		c.Tenant = c.Config.DefaultTenant
	}
	if err := c.checkAndRefreshAuth(); err != nil {
		return fmt.Errorf("authentication check failed: %w", err)
	}
	newApi, err := api.NewAPI(c.Config, c.Tenant, c.Logger)
	if err != nil {
		return fmt.Errorf("failed to initialize API client: %w", err)
	}
	c.API = newApi
	return nil
}

func (c *CLI) checkAndRefreshAuth() error {
	tenant, err := c.Config.GetTenant(c.Tenant)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}
	err = tenant.CheckAuthenticationStatus()
	if err != nil {
		return fmt.Errorf("failed to check authentication status: %w", err)
	}
	switch {
	case errors.Is(err, config.ErrInvalidToken):
		c.Logger.Info("Token is invalid, attempting to refresh")
		if err := c.refreshToken(); err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}
		return nil
	default:
		return nil
	}
}

func (c *CLI) refreshToken() error {
	// Implement token refresh logic here
	// This is a placeholder and should be replaced with actual refresh logic
	c.Logger.Info("Refreshing token")
	// newToken := "refreshed-token"
	// tenant.SetAccessToken(newToken)
	// return c.Config.AddTenant(*tenant)
	return fmt.Errorf("token refresh not implemented")
}
