package oauth2

import (
	"strings"
)

type clientCredentialsFlow struct {
	*oauth2Client
	scopes    []string
	grantType string
}

// ClientCredentials takes a scopes string array and returns the client credentials flow.
func (o *oauth2Client) ClientCredentials(scopes []string) *clientCredentialsFlow {
	return &clientCredentialsFlow{
		oauth2Client: o,
		scopes:       scopes,
		grantType:    GrantTypeClientCredentials,
	}
}

// Scopes returns the scopes used in this client credentials flow.
func (c *clientCredentialsFlow) Scopes() []string {
	return c.scopes
}

// GrantType returns the grant type of client credentials flow. Should always be "client_credentials".
func (c *clientCredentialsFlow) GrantType() string {
	return c.grantType
}

// AccessToken retrieves the access token for the client credentials oauth2 flow
func (c *clientCredentialsFlow) AccessToken() (*AccessToken, error) {
	rb := RequestBody{}
	rb.Set(FieldGrantType, c.grantType)
	rb.Set(FieldRedirectURI, c.redirectURI)

	requestBody := strings.NewReader(rb.Encode())

	return c.accessToken(requestBody)
}
