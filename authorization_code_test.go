// +build unit

package oauth2_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oauth2 "github.com/ralucas/go-ebay-oauth2"
	"github.com/ralucas/go-ebay-oauth2/double"
)

var (
	testBaseUrl      = "https://test.com"
	testClientId     = "test-client-id"
	testClientSecret = "test-client-secret"
	testRedirectUri  = "https://my.host.com/oauth"
	testScopes       = []string{"a", "b", "c"}
	testPrompt       = "test-prompt"
	testState        = "test-state"
	testCode         = "test-code"
)

func TestAuthorizationCode_GrantApplicationAccessURL(t *testing.T) {
	client := oauth2.NewClient(testBaseUrl, testClientId, testClientSecret, testRedirectUri)

	t.Run("NoOptions", func(t *testing.T) {
		ac := client.AuthorizationCode(
			testScopes,
		)

		assert.Empty(t, ac.Prompt())
		assert.Empty(t, ac.State())

		u, err := ac.GrantApplicationAccessURL()

		assert.Nil(t, err, fmt.Sprintf("%v", err))

		us := u.String()
		assert.Contains(t, us, testBaseUrl)
		assert.Contains(t, us, testClientId)
		assert.Contains(t, us, url.QueryEscape(testRedirectUri))
		assert.Contains(t, us, url.QueryEscape(strings.Join(testScopes, " ")))
	})

	t.Run("WithOptions", func(t *testing.T) {
		ac := client.AuthorizationCode(
			testScopes,
			oauth2.WithPrompt(testPrompt),
			oauth2.WithState(testState),
		)

		assert.Equal(t, ac.Prompt(), testPrompt)
		assert.Equal(t, ac.State(), testState)

		u, err := ac.GrantApplicationAccessURL()

		assert.Nil(t, err, fmt.Sprintf("%v", err))

		us := u.String()
		assert.Contains(t, us, testBaseUrl)
		assert.Contains(t, us, testClientId)
		assert.Contains(t, us, url.QueryEscape(testRedirectUri))
		assert.Contains(t, us, url.QueryEscape(strings.Join(testScopes, " ")))
		assert.Contains(t, us, testPrompt)
		assert.Contains(t, us, testState)
	})
}

func TestAuthorizationCode_ExchangeAuthorizationCodeForToken(t *testing.T) {
	client := oauth2.NewClient(testBaseUrl, testClientId, testClientSecret, testRedirectUri)

	t.Run("RetrievesToken", func(t *testing.T) {
		mockHttpClient := double.MockHTTPClient{}
		client.SetHTTPClient(&mockHttpClient)

		ac := client.AuthorizationCode(
			testScopes,
		)

		testURL, err := url.Parse(fmt.Sprintf("https://test.com/redirect?code=%s", testCode))
		require.Nil(t, err)

		token, err := ac.ExchangeAuthorizationForToken(testURL)
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		assert.NotNil(t, token)
	})

	t.Run("RetrievesTokenWithState", func(t *testing.T) {
		mockHttpClient := double.MockHTTPClient{}
		client.SetHTTPClient(&mockHttpClient)

		ac := client.AuthorizationCode(
			testScopes,
			oauth2.WithState(testState),
		)

		testURL, err := url.Parse(fmt.Sprintf("https://test.com/redirect?code=%s&state=%s", testCode, testState))
		require.Nil(t, err)

		token, err := ac.ExchangeAuthorizationForToken(testURL)
		assert.Nil(t, err, fmt.Sprintf("%v", err))

		assert.NotNil(t, token)
	})

	t.Run("ErrorsOnMissingCodeFromURL", func(t *testing.T) {
		mockHttpClient := double.MockHTTPClient{}
		client.SetHTTPClient(&mockHttpClient)

		ac := client.AuthorizationCode(
			testScopes,
		)

		testURL, err := url.Parse(fmt.Sprintf("https://test.com/redirect?missing=%s", testCode))
		require.Nil(t, err)

		token, err := ac.ExchangeAuthorizationForToken(testURL)

		assert.Nil(t, token, fmt.Sprintf("recieved token %v", token))
		assert.NotNil(t, err)
	})

	t.Run("ErrorsOnUnexpectedStateFromURL", func(t *testing.T) {
		mockHttpClient := double.MockHTTPClient{}
		client.SetHTTPClient(&mockHttpClient)

		ac := client.AuthorizationCode(
			testScopes,
		)

		testURL, err := url.Parse(fmt.Sprintf("https://test.com/redirect?code=%v&state=%v", testCode, testState))
		require.Nil(t, err)

		token, err := ac.ExchangeAuthorizationForToken(testURL)

		assert.Nil(t, token, fmt.Sprintf("recieved token %v", token))
		assert.NotNil(t, err)
	})

	t.Run("ErrorsOnMissingStateFromURL", func(t *testing.T) {
		mockHttpClient := double.MockHTTPClient{}
		client.SetHTTPClient(&mockHttpClient)

		ac := client.AuthorizationCode(
			testScopes,
			oauth2.WithState(testState),
		)

		testURL, err := url.Parse(fmt.Sprintf("https://test.com/redirect?code=%v", testCode))
		require.Nil(t, err)

		token, err := ac.ExchangeAuthorizationForToken(testURL)

		assert.Nil(t, token, fmt.Sprintf("recieved token %v", token))
		assert.NotNil(t, err)
	})

	t.Run("MakesExpectedRequest", func(t *testing.T) {
		spyHttpClient := &double.SpyHTTPClient{}
		spyHttpClient.Reset()

		client.SetHTTPClient(spyHttpClient)

		ac := client.AuthorizationCode(
			testScopes,
		)

		testURL, err := url.Parse(fmt.Sprintf("https://test.com/redirect?code=%s", testCode))
		require.Nil(t, err)

		_, err = ac.ExchangeAuthorizationForToken(testURL)
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
			"%s=%s\n%s=%s\n%s=%s",
			oauth2.FieldGrantType,
			"authorization_code",
			oauth2.FieldCode,
			testCode,
			oauth2.FieldRedirectURI,
			testRedirectUri,
		)

		assert.Equal(t, reqBody, string(body))
	})
}
