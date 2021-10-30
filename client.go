package oauth2

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// Content-Type used in OAuth2 requests
	contentType = "application/x-www-form-urlencoded"

	authorizationCodeResponseType = "code"

	// Paths
	authorizePath = "/oauth2/authorize"
	tokenPath     = "/identity/v1/oauth2/token"

	// Grant Types
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeClientCredentials = "client_credentials"

	// Base urls that can be used in configuration
	ProductionBaseUrl = "https://auth.ebay.com"
	SandboxBaseUrl    = "https://auth.sandbox.ebay.com"
)

type Flow int

const (
	AuthorizationCode Flow = iota
	ClientCredentials
)

var URLQueryPart = struct {
	Code        string
	State       string
	RedirectURI string
	GrantType   string
}{
	"code", "state", "redirect_uri", "grant_type",
}

// AccessToken is the response from each OAuth2 flow
type AccessToken struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
}

type Oauth2Client interface {
	BaseURL() string
	ClientID() string
	ClientSecret() string
	RedirectURI() string
	SetHTTPClient(HTTPClient)
	GetHTTPClient() HTTPClient
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

func NewClient(baseURL, clientID, clientSecret, redirectURI string) Oauth2Client {
	return &oauth2Client{
		baseURL, clientID, clientSecret, redirectURI, http.DefaultClient,
	}
}

func (o *oauth2Client) SetHTTPClient(c HTTPClient) {
	o.httpClient = c
}

func (o *oauth2Client) GetHTTPClient() HTTPClient {
	return o.httpClient
}

func (o *oauth2Client) BaseURL() string {
	return o.baseURL
}

func (o *oauth2Client) ClientID() string {
	return o.clientID
}

func (o *oauth2Client) ClientSecret() string {
	return o.clientSecret
}

func (o *oauth2Client) RedirectURI() string {
	return o.redirectURI
}

func (o *oauth2Client) basicAuthorization() string {
	return fmt.Sprintf(
		"Basic %s", base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", o.clientID, o.clientSecret)),
		),
	)
}

func (o *oauth2Client) accessToken(requestBody io.Reader) (*AccessToken, error) {
	requestUrl, err := url.Parse(o.baseURL)
	if err != nil {
		return nil, err
	}

	requestUrl.Path = tokenPath

	newReq, err := http.NewRequest(
		http.MethodPost,
		requestUrl.String(),
		requestBody,
	)
	if err != nil {
		return nil, err
	}

	newReq.Header.Add("Content-Type", contentType)
	newReq.Header.Add("Authorization", o.basicAuthorization())

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
