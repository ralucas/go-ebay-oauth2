// +build unit

package oauth2_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oauth2 "github.com/ralucas/go-ebay-oauth2"
	"github.com/ralucas/go-ebay-oauth2/double"
)

func TestClientCredentials_AccessToken(t *testing.T) {
	client := oauth2.NewClient(testBaseUrl, testClientId, testClientSecret, testRedirectUri)

	t.Run("SuccessfullyReturnsToken", func(t *testing.T) {
		client.SetHTTPClient(&double.MockHTTPClient{})

		cc := client.ClientCredentials(
			testScopes,
		)

		token, err := cc.AccessToken()
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		assert.NotNil(t, token)
	})

	t.Run("MakesExpectedRequest", func(t *testing.T) {
		spyHttpClient := &double.SpyHTTPClient{}
		spyHttpClient.Reset()

		client.SetHTTPClient(spyHttpClient)

		cc := client.ClientCredentials(
			testScopes,
		)

		_, err := cc.AccessToken()
		require.Nil(t, err)

		require.Equal(t, 1, spyHttpClient.CallCount("Do"))

		calls := spyHttpClient.Calls("Do")
		call := calls[0]
		callArgs := call.Arguments()
		callReq := callArgs[0].(*http.Request)

		reqUrl := fmt.Sprintf("%s%s", testBaseUrl, oauth2.TokenPath)

		assert.Equal(t, reqUrl, callReq.URL.String())

		callClientID, callClientSecret, ok := callReq.BasicAuth()
		require.True(t, ok)

		assert.Equal(t, testClientId, callClientID)
		assert.Equal(t, testClientSecret, callClientSecret)

		assert.Equal(t, "application/x-www-form-urlencoded", callReq.Header.Get("Content-Type"))

		body, err := io.ReadAll(callReq.Body)
		require.Nil(t, err)

		reqBody := fmt.Sprintf(
			"%s=%s\n%s=%s",
			oauth2.FieldGrantType,
			"client_credentials",
			oauth2.FieldRedirectURI,
			testRedirectUri,
		)

		assert.Equal(t, reqBody, string(body))
	})
}
