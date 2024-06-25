package management

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

type Management struct {
	Application *ApplicationManager
	http        *http.Client
	url         *url.URL
	basePath    string
	common      manager
}

type manager struct {
	management *Management
}

func New(tenantDomain string, accessToken string) (*Management, error) {

	server := "https://api.asgardeo.io/"
	basePath := "t/" + tenantDomain + "/api/server/v1"
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	m := &Management{
		url:      u,
		basePath: basePath,
		http: &http.Client{
			Transport: &transport{
				underlyingTransport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				token: accessToken,
			},
		},
	}

	m.common.management = m

	m.Application = (*ApplicationManager)(&m.common)

	return m, nil
}

type transport struct {
	underlyingTransport http.RoundTripper
	token               string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.token)
	return t.underlyingTransport.RoundTrip(req)
}
