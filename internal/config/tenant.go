package config

import "time"

type Tenant struct {
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token,omitempty"`
	ExpiresAt   time.Time `json:"expires_at"`
	ClientID    string    `json:"client_id"`
}
