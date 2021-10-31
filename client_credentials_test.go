// +build unit

package oauth2_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	oauth2 "github.com/ralucas/go-ebay-oauth2"
	"github.com/ralucas/go-ebay-oauth2/mocks"
)

func TestClientCredentials_AccessToken(t *testing.T) {
	client := oauth2.NewClient(testBaseUrl, testClientId, testClientSecret, testRedirectUri)
	client.SetHTTPClient(&mocks.MockHTTPClient{})

	cc := client.ClientCredentials(
		testScopes,
	)

	token, err := cc.AccessToken()
	assert.Nil(t, err, fmt.Sprintf("%v", err))

	assert.NotNil(t, token)
}
