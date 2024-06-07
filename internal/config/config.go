package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrConfigFileMissing = errors.New("config.json file is missing")
var ErrNoAuthenticatedTenants = errors.New("not logged in. Try `is login`")

type Config struct {
	initError     error
	path          string
	DefaultTenant string            `json:"default_tenant"`
	Tenants       map[string]Tenant `json:"tenants"`
}

func (c *Config) Initialize() error {

	c.initError = c.loadFromDisk()
	return c.initError
}

func (c *Config) Validate() error {
	if err := c.Initialize(); err != nil {
		return err
	}

	if len(c.Tenants) == 0 {
		return ErrNoAuthenticatedTenants
	}

	if c.DefaultTenant != "" {
		return nil
	}

	return c.saveToDisk()
}

func (c *Config) IsLoggedInWithTenant(tenantName string) bool {

	_ = c.Initialize()

	if tenantName == "" {
		tenantName = c.DefaultTenant
	}

	_, ok := c.Tenants[tenantName]
	if !ok {
		return false
	}

	//validate token

	return true
}

func (c *Config) GetTenant(tenantName string) (Tenant, error) {
	if err := c.Initialize(); err != nil {
		return Tenant{}, err
	}

	tenant, ok := c.Tenants[tenantName]
	if !ok {
		return Tenant{}, fmt.Errorf(
			"failed to find tenant: %s.",
			tenantName,
		)
	}

	return tenant, nil
}

func (c *Config) AddTenant(tenant Tenant) error {

	_ = c.Initialize()

	if c.DefaultTenant == "" {
		c.DefaultTenant = tenant.Name
	}

	if c.Tenants == nil {
		c.Tenants = make(map[string]Tenant)
	}

	c.Tenants[tenant.Name] = tenant

	return c.saveToDisk()
}

func (c *Config) SetDefaultTenant(tenantName string) error {
	tenant, err := c.GetTenant(tenantName)
	if err != nil {
		return err
	}

	c.DefaultTenant = tenant.Name

	return c.saveToDisk()
}

func (c *Config) saveToDisk() error {
	dir := filepath.Dir(c.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		const dirPerm os.FileMode = 0700
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return err
		}
	}

	buffer, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	const filePerm os.FileMode = 0600
	return os.WriteFile(c.path, buffer, filePerm)
}

func (c *Config) loadFromDisk() error {

	if c.path == "" {
		c.path = defaultPath()
	}

	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return ErrConfigFileMissing
	}

	buffer, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(buffer, c)
}

func defaultPath() string {

	return "/Users/thilinashashimalsenarath/Documents/my_projects/is-cli/.config/is/config.json"
}
