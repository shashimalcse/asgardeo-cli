package core

import (
	"net/http"

	"github.com/shashimalcse/is-cli/internal/auth"
	"github.com/shashimalcse/is-cli/internal/config"
	"github.com/shashimalcse/is-cli/internal/keyring"
)

type LoginInputs struct {
	ClientID     string
	ClientSecret string
	Tenant       string
}

func (i *LoginInputs) IsLoggingInAsAMachine() bool {
	return i.ClientID != "" || i.ClientSecret != "" || i.Tenant != ""
}

func RunLoginAsMachine(inputs LoginInputs, cli *CLI) error {

	result, err := auth.GetAccessTokenFromClientCreds(http.DefaultClient, auth.ClientCredentials{ClientID: inputs.ClientID, ClientSecret: inputs.ClientSecret, Tenant: inputs.Tenant})
	if err != nil {
		return err
	}
	tenant := config.Tenant{Name: inputs.Tenant, ClientID: inputs.ClientID, AccessToken: result.AccessToken}
	if err := keyring.StoreAccessToken(inputs.Tenant, result.AccessToken); err != nil {
		// In case we don't have a keyring, we want the
		// access token to be saved in the config file.
	}
	err = cli.Config.AddTenant(tenant)
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
		// In case we don't have a keyring, we want the
		// access token to be saved in the config file.
	}
	err = cli.Config.AddTenant(tenant)
	cli.Config.DefaultTenant = tenant.Name
	return nil
}
