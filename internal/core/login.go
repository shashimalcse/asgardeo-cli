package core

import (
	"net/http"

	"github.com/shashimalcse/is-cli/internal/auth"
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

	_, err := auth.GetAccessTokenFromClientCreds(http.DefaultClient, auth.ClientCredentials{ClientID: inputs.ClientID, ClientSecret: inputs.ClientSecret, Tenant: inputs.Tenant})
	if err != nil {
		return err
	}
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

	_, err := auth.GetAccessTokenFromDeviceCode(http.DefaultClient, state)
	if err != nil {
		return err
	}
	return nil
}
