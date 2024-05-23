package auth

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Result struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type ClientCredentials struct {
	ClientID     string
	ClientSecret string
	Tenant       string
}

func GetAccessTokenFromClientCreds(httpClient *http.Client, args ClientCredentials) (Result, error) {

	httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	data := url.Values{
		"grant_type": {"client_credentials"},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://localhost:9443/t/%s/oauth2/token", args.Tenant), strings.NewReader(data.Encode()))
	if err != nil {
		return Result{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Add the Authorization header with the base64 encoded client ID and client secret
	auth := args.ClientID + ":" + args.ClientSecret
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

	resp, err := httpClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Result{}, err
	}
	return result, nil
}
