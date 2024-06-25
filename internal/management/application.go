package management

import "context"

type Application struct {
	ID                           string                       `json:"id,omitempty"`
	Name                         string                       `json:"name,omitempty"`
	LogoutReturnURL              string                       `json:"logoutReturnUrl,omitempty"`
	ClientID                     string                       `json:"clientId,omitempty"`
	Issuer                       string                       `json:"issuer,omitempty"`
	Realm                        string                       `json:"realm,omitempty"`
	TemplateID                   string                       `json:"templateId,omitempty"`
	IsManagementApp              bool                         `json:"isManagementApp,omitempty"`
	AssociatedRoles              AssociatedRoles              `json:"associatedRoles,omitempty"`
	ClaimConfiguration           ClaimConfiguration           `json:"claimConfiguration,omitempty"`
	InboundProtocols             []InboundProtocol            `json:"inboundProtocols,omitempty"`
	InboundProtocolConfiguration InboundProtocolConfiguration `json:"inboundProtocolConfiguration,omitempty"`
	AuthenticationSeq            AuthenticationSequence       `json:"authenticationSequence,omitempty"`
	AdvancedConfig               AdvancedConfigurations       `json:"advancedConfigurations,omitempty"`
	ProvisioningConfig           ProvisioningConfigurations   `json:"provisioningConfigurations,omitempty"`
	Access                       string                       `json:"access,omitempty"`
}

type ApplicationList struct {
	TotalResults int           `json:"totalResults"`
	StartIndex   int           `json:"startIndex"`
	Count        int           `json:"count"`
	Applications []Application `json:"applications"`
	Links        []Link        `json:"links"`
}

type ApplicationManager manager

type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type AssociatedRoles struct {
	AllowedAudience string           `json:"allowedAudience"`
	Roles           []AssociatedRole `json:"roles,omitempty"`
}

type AssociatedRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClaimConfiguration struct {
	Dialect         string        `json:"dialect"`
	ClaimMappings   []interface{} `json:"claimMappings"`
	RequestedClaims []interface{} `json:"requestedClaims"`
	Subject         Subject       `json:"subject"`
	Role            Role          `json:"role"`
}

type Subject struct {
	Claim                       Claim `json:"claim"`
	IncludeUserDomain           bool  `json:"includeUserDomain"`
	IncludeTenantDomain         bool  `json:"includeTenantDomain"`
	UseMappedLocalSubject       bool  `json:"useMappedLocalSubject"`
	MappedLocalSubjectMandatory bool  `json:"mappedLocalSubjectMandatory"`
}

type Claim struct {
	URI string `json:"uri"`
}

type Role struct {
	IncludeUserDomain bool  `json:"includeUserDomain"`
	Claim             Claim `json:"claim"`
}

type InboundProtocol struct {
	Type string `json:"type"`
	Self string `json:"self"`
}

type InboundProtocolConfiguration struct {
	OIDC OIDC `json:"oidc"`
}

type OIDC struct {
	AccessToken    AccessToken  `json:"accessToken"`
	GrantTypes     []string     `json:"grantTypes"`
	AllowedOrigins []string     `json:"allowedOrigins"`
	CallbackURLs   []string     `json:"callbackURLs"`
	PKCE           PKCE         `json:"pkce"`
	PublicClient   bool         `json:"publicClient"`
	RefreshToken   RefreshToken `json:"refreshToken"`
}

type PKCE struct {
	Mandatory                      bool `json:"mandatory"`
	SupportPlainTransformAlgorithm bool `json:"supportPlainTransformAlgorithm"`
}

type RefreshToken struct {
	ExpiryInSeconds   int  `json:"expiryInSeconds"`
	RenewRefreshToken bool `json:"renewRefreshToken"`
}

type AccessToken struct {
	ApplicationAccessTokenExpiryInSeconds int    `json:"applicationAccessTokenExpiryInSeconds"`
	BindingType                           string `json:"bindingType"`
	RevokeTokensWhenIDPSessionTerminated  bool   `json:"revokeTokensWhenIDPSessionTerminated"`
	Type                                  string `json:"type"`
	UserAccessTokenExpiryInSeconds        int    `json:"userAccessTokenExpiryInSeconds"`
	ValidateTokenBinding                  bool   `json:"validateTokenBinding"`
}

type AuthenticationSequence struct {
	Type                      string        `json:"type"`
	Steps                     []Step        `json:"steps"`
	RequestPathAuthenticators []interface{} `json:"requestPathAuthenticators"`
	SubjectStepID             int           `json:"subjectStepId"`
	AttributeStepID           int           `json:"attributeStepId"`
}

type Step struct {
	ID      int       `json:"id"`
	Options []Options `json:"options"`
}

type Options struct {
	IDP           string `json:"idp"`
	Authenticator string `json:"authenticator"`
}

type AdvancedConfigurations struct {
	Saas                         bool                   `json:"saas"`
	DiscoverableByEndUsers       bool                   `json:"discoverableByEndUsers"`
	SkipLoginConsent             bool                   `json:"skipLoginConsent"`
	SkipLogoutConsent            bool                   `json:"skipLogoutConsent"`
	ReturnAuthenticatedIdpList   bool                   `json:"returnAuthenticatedIdpList"`
	EnableAuthorization          bool                   `json:"enableAuthorization"`
	Fragment                     bool                   `json:"fragment"`
	EnableAPIBasedAuthentication bool                   `json:"enableAPIBasedAuthentication"`
	AttestationMetaData          AttestationMetaData    `json:"attestationMetaData"`
	AdditionalSpProperties       []AdditionalSpProperty `json:"additionalSpProperties"`
	UseExternalConsentPage       bool                   `json:"useExternalConsentPage"`
}

type AttestationMetaData struct {
	EnableClientAttestation bool   `json:"enableClientAttestation"`
	AndroidPackageName      string `json:"androidPackageName"`
	AppleAppID              string `json:"appleAppId"`
}

type AdditionalSpProperty struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	DisplayName string `json:"displayName"`
}

type ProvisioningConfigurations struct {
	OutboundProvisioningIDPs []interface{} `json:"outboundProvisioningIdps"`
}

func (m *ApplicationManager) List(ctx context.Context) (a *ApplicationList, err error) {
	err = m.management.Request(ctx, "GET", m.management.URI("applications"), &a)
	return
}

func (m *ApplicationManager) Create(ctx context.Context, application Application) (a *Application, err error) {
	err = m.management.Request(ctx, "POST", m.management.URI("applications"), application)
	return
}
