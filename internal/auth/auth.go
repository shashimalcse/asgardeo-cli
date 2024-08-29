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

var (
	SystemScope = "SYSTEM"
)

type Result struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
}

type ClientCredentials struct {
	ClientID     string
	ClientSecret string
	Tenant       string
}

type State struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	VerificationURI         string `json:"verification_uri"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type Credentials struct {
	ClientID string
	Tenant   string
}

var credentials = &Credentials{
	ClientID: "",
	Tenant:   "",
}

func GetDeviceCode(httpClient *http.Client) (State, error) {

	httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	a := credentials

	data := url.Values{
		"client_id": {a.ClientID},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.asgardeo.io/t/%s/oauth2/device_authorize", a.Tenant), strings.NewReader(data.Encode()))
	if err != nil {
		return State{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return State{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return State{}, fmt.Errorf("failed to get device code: %s", resp.Status)

	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			err = fmt.Errorf("failed to close response body: %w", cErr)
		}
	}()
	var result State
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return State{}, err
	}
	return result, nil
}

func GetAccessTokenFromDeviceCode(httpClient *http.Client, state State) (Result, error) {

	httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	data := url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"client_id":   {credentials.ClientID},
		"device_code": {state.DeviceCode},
		"scope":       {"SYSTEM"},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.asgardeo.io/t/%s/oauth2/token", credentials.Tenant), strings.NewReader(data.Encode()))
	if err != nil {
		return Result{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("failed to get access token: %s", resp.Status)
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			err = fmt.Errorf("failed to close response body: %w", cErr)
		}
	}()
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func AuthenticateWithClientCredentials(httpClient *http.Client, args ClientCredentials) (Result, error) {

	data := url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {SystemScope},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.asgardeo.io/t/%s/oauth2/token", args.Tenant), strings.NewReader(data.Encode()))
	if err != nil {
		return Result{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	encodedAuth := getBasicAuth(args.ClientID, args.ClientSecret)
	req.Header.Add("Authorization", "Basic "+encodedAuth)
	resp, err := httpClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return Result{}, fmt.Errorf("failed to authenticate. please check your credentials")
		} else if resp.StatusCode == http.StatusNotFound {
			return Result{}, fmt.Errorf("failed to authenticate. tenant not found")
		}
		return Result{}, fmt.Errorf("failed to authenticate")
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			err = fmt.Errorf("failed to close response body: %w", cErr)
		}
	}()
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func getBasicAuth(clientID, clientSecret string) string {
	auth := clientID + ":" + clientSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
