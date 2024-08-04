package api

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/shashimalcse/is-cli/internal/config"
)

type httpClient struct {
	client   *http.Client
	url      *url.URL
	basePath string
}

func NewHTTPClientAPI(cfg *config.Config, tenantDomain string) (*httpClient, error) {
	tenant, err := cfg.GetTenant(tenantDomain)
	if err != nil {
		return nil, err
	}
	basePath := "t/" + tenant.Name + "/api/server/v1"
	u, err := url.Parse("https://api.asgardeo.io/")
	if err != nil {
		return nil, err
	}
	return &httpClient{client: newHTTPClient(tenant.AccessToken, true), basePath: basePath, url: u}, nil
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
