// +build integration

package oauth2_test

import (
	"strings"
	"testing"

	"github.com/plaid/go-envvar/envvar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oauth2 "github.com/ralucas/go-ebay-oauth2"
)

type integrationEnv struct {
	ClientID     string `envvar:"CLIENT_ID"`
	ClientSecret string `envvar:"CLIENT_SECRET"`
	RedirectURI  string `envvar:"REDIRECT_URI"`
	Scopes       string `envvar:"SCOPES"`
}

func TestClientCredentialsFlow(t *testing.T) {
	env := integrationEnv{}
	err := envvar.Parse(&env)

	client := oauth2.NewClient(
		oauth2.SandboxBaseURL,
		env.ClientID,
		env.ClientSecret,
		env.RedirectURI,
	)

	scopes := strings.Split(env.Scopes, " ")

	cc := client.ClientCredentials(scopes)

	token, err := cc.AccessToken()
	require.Nil(t, err)

	assert.NotNil(t, token.AccessToken)
	assert.NotNil(t, token.ExpiresIn)
	assert.NotNil(t, token.TokenType)
}
