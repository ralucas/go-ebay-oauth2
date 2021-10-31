// Package go-ebay-oauth2 provides methods for handling OAuth2 for making eBay REST requests.

// Example Usage:

// 	import oauth2 "github.com/ralucas/go-ebay-oauth2"

// Create the client with your appropriate credentials:

// 	client := oauth2.NewClient(
// 		baseUrl, // oauth2.ProductionBaseUrl or oauth2.SandboxBaseUrl constants are available
// 		clientId, // your eBay application clientId
// 		clientSecret, // your eBay application clientSecret
// 		redirectUri, // your application redirectUri
// 	)

// Then, create the flow you wish to pursue:

// Authorization Code Flow:

// 	myscopes := []string{"offer"}
// 	ac := client.AuthorizationCode(myscopes)

// Authorization Code Flow with options:

// 	var opts []AuthorizationCodeOption
// 	opts = opts.append(WithPrompt("login"))
// 	opts = opts.append(WithState("my-state-token"))

// 	ac := client.AuthorizationCode(config, myscopes, opts...)

// Get the request url for application access:

// 	url, err := ac.GrantApplicationAccessUrl()

// Next, in your redirect route handler, get the access token:

// 	token, err := ac.ExchangeAuthorizationForToken(req.redirectURI)
// 	// do something with the token

// Client Credentials Flow:

// 	myscopes := []string{"offer"}
// 	cc := client.ClientCredentials(myscopes)

// 	token, err := cc.AccessToken()
// 	// do something with the token
package oauth2
