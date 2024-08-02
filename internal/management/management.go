package management

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/shashimalcse/is-cli/internal/config"
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

func New(cfg *config.Config, tenantDomain string) (*Management, error) {
	tenant, err := cfg.GetTenant(tenantDomain)
	if err != nil {
		return nil, err
	}
	basePath := "t/" + tenant.Name + "/api/server/v1"
	u, err := url.Parse("https://api.asgardeo.io/")
	if err != nil {
		return nil, err
	}
	m := &Management{
		url:      u,
		basePath: basePath,
		http:     newHTTPClient(tenant.AccessToken, true),
	}
	m.common.management = m
	m.Application = (*ApplicationManager)(&m.common)
	return m, nil
}

func newHTTPClient(token string, insecureSkipTLS bool) *http.Client {
	return &http.Client{
		Transport: &transport{
			underlyingTransport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipTLS},
			},
			token: token,
		},
	}
}

type transport struct {
	underlyingTransport http.RoundTripper
	token               string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.token)
	return t.underlyingTransport.RoundTrip(req)
}
