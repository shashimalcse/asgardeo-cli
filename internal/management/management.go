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

	tenantDomain = "https://localhost:9443/t/" + tenantDomain

	u, err := url.Parse(tenantDomain)
	if err != nil {
		return nil, err
	}

	m := &Management{
		url:      u,
		basePath: "api/server/v1",
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
