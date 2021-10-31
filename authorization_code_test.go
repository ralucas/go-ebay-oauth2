// +build unit

package oauth2_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oauth2 "github.com/ralucas/go-ebay-oauth2"
	"github.com/ralucas/go-ebay-oauth2/mocks"
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
	client.SetHTTPClient(&mocks.MockHTTPClient{})

	t.Run("RetrievesToken", func(t *testing.T) {
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
}
