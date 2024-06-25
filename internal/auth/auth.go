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
	ClientID: "Wkwv5_jmo2DJVoul3bW7qve46C4a",
	Tenant:   "carbon.super",
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
	defer resp.Body.Close()
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
	defer resp.Body.Close()
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func GetAccessTokenFromClientCreds(httpClient *http.Client, args ClientCredentials) (Result, error) {

	httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	data := url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"SYSTEM"},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.asgardeo.io/t/%s/oauth2/token", args.Tenant), strings.NewReader(data.Encode()))
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
