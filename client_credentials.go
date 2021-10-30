package oauth2

import (
	"fmt"
	"strings"
)

type clientCredentialsFlow struct {
	*oauth2Client
	scopes    []string
	grantType string
}

func (o *oauth2Client) ClientCredentials(scopes []string) *clientCredentialsFlow {
	return &clientCredentialsFlow{
		oauth2Client: o,
		scopes:       scopes,
		grantType:    GrantTypeClientCredentials,
	}
}

func (c *clientCredentialsFlow) Scopes() []string {
	return c.scopes
}

func (c *clientCredentialsFlow) GrantType() string {
	return c.grantType
}

func (c *clientCredentialsFlow) AccessToken() (*AccessToken, error) {
	requestBody := strings.NewReader(
		fmt.Sprintf(
			"grant_type=%s\nredirect_uri=%s",
			c.grantType,
			c.redirectURI,
		),
	)

	return c.accessToken(requestBody)
}
