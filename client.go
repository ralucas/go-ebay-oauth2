package oauth2

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	// Content-Type used in OAuth2 requests
	ContentType = "application/x-www-form-urlencoded"

	// Paths
	AuthorizePath = "/oauth2/authorize"
	TokenPath     = "/identity/v1/oauth2/token"

	// Grant Types
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeClientCredentials = "client_credentials"

	// Fields
	FieldCode         = "code"
	FieldState        = "state"
	FieldRedirectURI  = "redirect_uri"
	FieldGrantType    = "grant_type"
	FieldResponseType = "response_type"
	FieldClientID     = "client_id"
	FieldScope        = "scope"
	FieldPrompt       = "prompt"

	// Base urls that can be used in configuration
	ProductionBaseURL = "https://api.ebay.com"
	SandboxBaseURL    = "https://api.sandbox.ebay.com"
)

type Flow int

const (
	AuthorizationCode Flow = iota
	ClientCredentials
)

// AccessToken is the response from each OAuth2 flow.
type AccessToken struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
}

// RequestBody maps a string key to a value string. It is used for creating a "URL Encoded"
// request body used in oauth2 post requests.
type RequestBody map[string]string

// Encode encodes the values into "URL encoded" form ("bar=baz&foo=quux").
func (rb RequestBody) Encode() string {
	sb := strings.Builder{}

	for key, val := range rb {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(val)
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// Set sets the key to value. It replaces any existing values.
func (rb RequestBody) Set(key, val string) {
	rb[key] = val
}

type Oauth2Client interface {
	BaseURL() string
	ClientID() string
	ClientSecret() string
	RedirectURI() string
	SetHTTPClient(HTTPClient)
	AuthorizationCode([]string, ...authorizationCodeOption) *authorizationCodeFlow
	ClientCredentials([]string) *clientCredentialsFlow
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type oauth2Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	redirectURI  string
	httpClient   HTTPClient
}

// NewClient creates a new Oauth2Client.
func NewClient(baseURL, clientID, clientSecret, redirectURI string) Oauth2Client {
	return &oauth2Client{
		baseURL, clientID, clientSecret, redirectURI, http.DefaultClient,
	}
}

// SetHTTPClient sets an http client that satisfies the HTTPClient interface.
func (o *oauth2Client) SetHTTPClient(c HTTPClient) {
	o.httpClient = c
}

// BaseURL returns the base url used by the Oauth2Client
func (o *oauth2Client) BaseURL() string {
	return o.baseURL
}

// ClientID returns the client ID used by the Oauth2Client
func (o *oauth2Client) ClientID() string {
	return o.clientID
}

// ClientSecret returns the client secret used by the Oauth2Client
func (o *oauth2Client) ClientSecret() string {
	return o.clientSecret
}

// RedirectURI returns the redirect uri used by the Oauth2Client
func (o *oauth2Client) RedirectURI() string {
	return o.redirectURI
}

// accessToken takes a request body in io.Reader type, sends the request to the
// oauth2 token path returning either the access token or an error.
func (o *oauth2Client) accessToken(requestBody io.Reader) (*AccessToken, error) {
	requestUrl, err := url.Parse(o.baseURL)
	if err != nil {
		return nil, err
	}

	requestUrl.Path = TokenPath

	newReq, err := http.NewRequest(
		http.MethodPost,
		requestUrl.String(),
		requestBody,
	)
	if err != nil {
		return nil, err
	}

	newReq.SetBasicAuth(o.clientID, o.clientSecret)

	newReq.Header.Add("Content-Type", ContentType)

	res, err := o.httpClient.Do(newReq)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	acr := AccessToken{}
	if json.Unmarshal(resBody, &acr) != nil {
		return nil, err
	}

	return &acr, nil
}
