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

func WithPrompt(prompt string) authorizationCodeOption {
	return func(a *authorizationCodeFlow) {
		a.prompt = prompt
	}
}

func WithState(state string) authorizationCodeOption {
	return func(a *authorizationCodeFlow) {
		a.state = state
	}
}

func (a *authorizationCodeFlow) Scopes() []string {
	return a.scopes
}

func (a *authorizationCodeFlow) GrantType() string {
	return a.grantType
}

func (a *authorizationCodeFlow) Prompt() string {
	return a.prompt
}

func (a *authorizationCodeFlow) State() string {
	return a.state
}

func (a *authorizationCodeFlow) GrantApplicationAccessUrl() (*url.URL, error) {
	requestUrl, err := url.Parse(a.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %+v\n", err)
	}

	requestUrl.Path = authorizePath

	qs := url.Values{}
	qs.Add("client_id", a.clientID)
	qs.Add("redirect_uri", a.redirectURI)
	qs.Add("response_type", authorizationCodeResponseType)
	qs.Add("scope", strings.Join(a.scopes, " ")) // a URL-encoded string of space-separated scopes

	if a.prompt != "" {
		qs.Add("prompt", a.prompt)
	}

	if a.state != "" {
		qs.Add("state", a.state)
	}

	requestUrl.RawQuery = qs.Encode()

	return requestUrl, nil
}

func (a *authorizationCodeFlow) ExchangeAuthorizationForToken(reqURL *url.URL) (*AccessToken, error) {
	code := reqURL.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("missing code in query %s\n", reqURL.String())
	}

	state := reqURL.Query().Get("state")
	if (state != "" || a.state != "") && state != a.state {
		return nil, fmt.Errorf("mismatching states, got %s, expected %s", state, a.state)
	}

	requestBody := strings.NewReader(
		fmt.Sprintf(
			"grant_type=%s\ncode=%s\nredirect_uri=%s",
			a.grantType,
			code,
			a.redirectURI,
		),
	)

	return a.accessToken(requestBody)
}
