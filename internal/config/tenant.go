package config

import (
	"errors"
	"time"

	"github.com/shashimalcse/is-cli/internal/keyring"
)

const accessTokenExpThreshold = 5 * time.Minute

var ErrInvalidToken = errors.New("token is invalid")

type Tenant struct {
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token,omitempty"`
	ExpiresAt   time.Time `json:"expires_at"`
	ClientID    string    `json:"client_id"`
}

func (t *Tenant) HasExpiredToken() bool {
	return time.Now().Add(accessTokenExpThreshold).After(t.ExpiresAt)
}

func (t *Tenant) GetAccessToken() string {
	accessToken, err := keyring.GetAccessToken(t.Name)
	if err == nil && accessToken != "" {
		return accessToken
	}

	return t.AccessToken
}

func (t *Tenant) CheckAuthenticationStatus() error {

	accessToken := t.GetAccessToken()
	if accessToken != "" {
		return nil
	}

	return ErrInvalidToken
}
