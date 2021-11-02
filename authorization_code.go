package oauth2

import (
	"fmt"
	"net/url"
	"strings"
)

type authorizationCodeFlow struct {
	*oauth2Client
	scopes    []string
	state     string
	prompt    string
	grantType string
}

type authorizationCodeOption func(*authorizationCodeFlow)

// AuthorizationCode takes scopes as an array of strings and options and
// returns the authorization code flow object.
func (o *oauth2Client) AuthorizationCode(
	scopes []string,
	options ...authorizationCodeOption,
) *authorizationCodeFlow {

	acf := &authorizationCodeFlow{
		oauth2Client: o,
		scopes:       scopes,
		grantType:    GrantTypeAuthorizationCode,
	}

	acf.applyOptions(options...)

	return acf
}

func (a *authorizationCodeFlow) applyOptions(options ...authorizationCodeOption) {
	for _, opt := range options {
		opt(a)
	}
}

// WithPrompt is used for adding the "prompt" option used in the GrantApplicationAccessURL.
func WithPrompt(prompt string) authorizationCodeOption {
	return func(a *authorizationCodeFlow) {
		a.prompt = prompt
	}
}

// WithState is used for adding the "state" option used for maintaining an optional application
// state token throughout the authorization code flow.
func WithState(state string) authorizationCodeOption {
	return func(a *authorizationCodeFlow) {
		a.state = state
	}
}

// Scopes returns the scopes used in this authorization code flow.
func (a *authorizationCodeFlow) Scopes() []string {
	return a.scopes
}

// GrantType returns the grant type of authorization code flow. Should always be "authorization_code".
func (a *authorizationCodeFlow) GrantType() string {
	return a.grantType
}

// Prompt returns the prompt (if any, if not returns empty string) used in this authorization code flow.
func (a *authorizationCodeFlow) Prompt() string {
	return a.prompt
}

// State returns the state token (if any, if not returns empty string) used in this authorization code flow.
func (a *authorizationCodeFlow) State() string {
	return a.state
}

// GrantApplicationAccessURL builds and returns the url used for starting the
// oauth2 authorization code flow, redirecting the user to authorize (or not)
// your application.
func (a *authorizationCodeFlow) GrantApplicationAccessURL() (*url.URL, error) {
	requestUrl, err := url.Parse(a.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %+v\n", err)
	}

	requestUrl.Path = AuthorizePath

	qs := url.Values{}
	qs.Set(FieldClientID, a.clientID)
	qs.Set(FieldRedirectURI, a.redirectURI)
	qs.Set(FieldResponseType, FieldCode)
	qs.Set(FieldScope, strings.Join(a.scopes, " ")) // a URL-encoded string of space-separated scopes

	if a.prompt != "" {
		qs.Set(FieldPrompt, a.prompt)
	}

	if a.state != "" {
		qs.Set(FieldState, a.state)
	}

	requestUrl.RawQuery = qs.Encode()

	return requestUrl, nil
}

// ExchangeAuthorizationForToken takes the url invoked in the GET request to the redirect URI,
// parses the query out of that url for the "code" and potentially the "state" fields, and exchanges
// those in a request to the token endpoint for the access token.
func (a *authorizationCodeFlow) ExchangeAuthorizationForToken(reqURL *url.URL) (*AccessToken, error) {
	code := reqURL.Query().Get(FieldCode)
	if code == "" {
		return nil, fmt.Errorf("missing code in query %s\n", reqURL.String())
	}

	state := reqURL.Query().Get(FieldState)
	if (state != "" || a.state != "") && state != a.state {
		return nil, fmt.Errorf("mismatching states, got %s, expected %s", state, a.state)
	}

	rb := RequestBody{}
	rb.Set(FieldGrantType, a.grantType)
	rb.Set(FieldCode, code)
	rb.Set(FieldRedirectURI, a.redirectURI)

	requestBody := strings.NewReader(rb.Encode())

	return a.accessToken(requestBody)
}
