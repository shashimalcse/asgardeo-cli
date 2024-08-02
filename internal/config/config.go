package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var ErrConfigFileMissing = errors.New("config.json file is missing")
var ErrNoAuthenticatedTenants = errors.New("not logged in. Try `is login`")

type Config struct {
	mu            sync.RWMutex
	path          string
	Server        string            `json:"server"`
	DefaultTenant string            `json:"default_tenant"`
	Tenants       map[string]Tenant `json:"tenants"`
	initialized   bool
}

// NewConfig creates a new Config instance
func NewConfig() *Config {
	return &Config{
		path:    defaultPath(),
		Tenants: make(map[string]Tenant),
	}
}

// Initialize loads the configuration from disk
func (c *Config) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return nil
	}

	if err := c.loadFromDisk(); err != nil && !errors.Is(err, ErrConfigFileMissing) {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	c.initialized = true
	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if err := c.Initialize(); err != nil {
		return err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Tenants) == 0 {
		return ErrNoAuthenticatedTenants
	}
	if c.DefaultTenant == "" {
		return errors.New("no default tenant set")
	}
	return nil
}

// IsLoggedInWithTenant checks if the user is logged in with a specific tenant
func (c *Config) IsLoggedInWithTenant(tenantName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if tenantName == "" {
		tenantName = c.DefaultTenant
	}
	_, ok := c.Tenants[tenantName]
	// TODO: validate token
	return ok
}

// GetTenant retrieves a tenant by name
func (c *Config) GetTenant(tenantName string) (Tenant, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tenant, ok := c.Tenants[tenantName]
	if !ok {
		return Tenant{}, fmt.Errorf("tenant not found: %s", tenantName)
	}
	return tenant, nil
}

// AddTenant adds a new tenant to the configuration
func (c *Config) AddTenant(tenant Tenant) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.DefaultTenant == "" {
		c.DefaultTenant = tenant.Name
	}
	c.Tenants[tenant.Name] = tenant
	return c.saveToDisk()
}

// RemoveTenant removes a tenant from the configuration
func (c *Config) RemoveTenant(tenant string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Tenants, tenant)
	if c.DefaultTenant == tenant {
		return c.setDefaultTenant()
	}
	return c.saveToDisk()
}

// SetDefaultTenant sets the default tenant
func (c *Config) SetDefaultTenant(tenantName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Tenants[tenantName]; !ok {
		return fmt.Errorf("tenant not found: %s", tenantName)
	}
	c.DefaultTenant = tenantName
	return c.saveToDisk()
}

func (c *Config) setDefaultTenant() error {
	for tenantName := range c.Tenants {
		c.DefaultTenant = tenantName
		return c.saveToDisk()
	}
	return nil
}

func (c *Config) saveToDisk() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	buffer, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.path, buffer, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func (c *Config) loadFromDisk() error {
	buffer, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrConfigFileMissing
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(buffer, c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

func defaultPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return filepath.Join(cwd, ".config", "config.json")
}
