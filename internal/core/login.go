package core

import (
	"fmt"
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

	cli.Logger.Info("Running login as machine - " + inputs.ClientID + " - " + inputs.ClientSecret + " - " + inputs.Tenant)
	result, err := auth.GetAccessTokenFromClientCreds(http.DefaultClient, auth.ClientCredentials{ClientID: inputs.ClientID, ClientSecret: inputs.ClientSecret, Tenant: inputs.Tenant})
	if err != nil {
		return err
	}
	fmt.Println("Access Token:", result.AccessToken)
	return nil
}
