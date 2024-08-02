package api

import (
	"context"

	"github.com/shashimalcse/is-cli/internal/management"
)

type ApplicationAPI interface {

	// List all applications.
	List(ctx context.Context) (c *management.ApplicationList, err error)

	// Create a new application.
	Create(ctx context.Context, application map[string]interface{}) (err error)

	// Delete an application.
	Delete(ctx context.Context, id string) (err error)
}
