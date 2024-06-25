package api

import (
	"context"

	"github.com/shashimalcse/is-cli/internal/management"
)

type ApplicationAPI interface {

	// List all applications.
	List(ctx context.Context) (c *management.ApplicationList, err error)

	// Create a new application.
	Create(ctx context.Context, application management.Application) (c *management.Application, err error)
}
