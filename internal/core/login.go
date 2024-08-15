package core

import (
	"net/http"
	"time"

	"github.com/shashimalcse/asgardeo-cli/internal/auth"
	"github.com/shashimalcse/asgardeo-cli/internal/config"
	"github.com/shashimalcse/asgardeo-cli/internal/keyring"
)

type LoginInputs struct {
	ClientID     string
	ClientSecret string
	Tenant       string
}

func (i *LoginInputs) IsLoggingInAsAMachine() bool {
	return i.ClientID != "" || i.ClientSecret != "" || i.Tenant != ""
}

func AuthenticateWithClientCredentials(inputs LoginInputs, cli *CLI) error {

	result, err := auth.AuthenticateWithClientCredentials(http.DefaultClient, auth.ClientCredentials{ClientID: inputs.ClientID, ClientSecret: inputs.ClientSecret, Tenant: inputs.Tenant})
	if err != nil {
		return err
	}
	expireIn := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	tenant := config.Tenant{
		Name:      inputs.Tenant,
		ClientID:  inputs.ClientID,
		ExpiresIn: expireIn,
	}
	if err := keyring.StoreAccessToken(inputs.Tenant, result.AccessToken); err != nil {
		tenant.AccessToken = result.AccessToken
	}
	err = cli.Config.AddTenant(tenant)
	if err != nil {
		return err
	}
	cli.Config.DefaultTenant = tenant.Name
	return nil
}

func GetDeviceCode(cli *CLI) (auth.State, error) {

	result, err := auth.GetDeviceCode(http.DefaultClient)
	if err != nil {
		return auth.State{}, err
	}
	return result, nil
}

func GetAccessTokenFromDeviceCode(cli *CLI, state auth.State) error {

	result, err := auth.GetAccessTokenFromDeviceCode(http.DefaultClient, state)
	if err != nil {
		return err
	}
	tenant := config.Tenant{Name: "carbon.super", ClientID: "Wkwv5_jmo2DJVoul3bW7qve46C4a", AccessToken: result.AccessToken}
	if err := keyring.StoreAccessToken("carbon.super", result.AccessToken); err != nil {
	}
	err = cli.Config.AddTenant(tenant)
	cli.Config.DefaultTenant = tenant.Name
	return nil
}
