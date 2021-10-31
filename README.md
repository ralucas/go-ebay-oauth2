Ebay OAuth2 for Go
---
[![Go](https://github.com/ralucas/go-ebay-oauth2/actions/workflows/go.yml/badge.svg)](https://github.com/ralucas/go-ebay-oauth2/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ralucas/go-ebay-oauth2.svg)](https://pkg.go.dev/github.com/ralucas/go-ebay-oauth2)

## Introduction
A library for handling OAuth2 for making eBay REST requests.

## Getting Started
Import the module:  
```go
import oauth2 "github.com/ralucas/go-ebay-oauth2"
```

Create the client with your appropriate credentials:
```go
client := oauth2.NewClient(
  baseUrl, // oauth2.ProductionBaseUrl or oauth2.SandboxBaseUrl constants are available
  clientId, // your eBay application clientId
  clientSecret, // your eBay application clientSecret
  redirectUri, // your application redirectUri
)
```

Then, create the flow you wish to pursue:

### Authorization Code Flow

Example:
```go
myscopes := []string{"offer"}
ac := client.AuthorizationCode(myscopes)
```

#### Options
Additional options are available:
- `prompt` - to indicate whether or not to prompt the user to login
  - `WithPrompt(string)`
- `state` - a token created by your application to manage the state
  - `WithState(string)`

Example usage:
```go
myscopes := []string{"offer"}

var opts []AuthorizationCodeOption
opts = opts.append(WithPrompt("login"))
opts = opts.append(WithState("my-state-token"))

ac := client.AuthorizationCode(config, myscopes, opts...)
```

Get the request url for application access:
```go
url, err := ac.GrantApplicationAccessUrl()
```

Next, in your redirect route handler, get the access token:
```go
token, err := ac.ExchangeAuthorizationForToken(req.redirectURI)
// do something with the token
```

### Client Credentials Flow
Example:
```go
myscopes := []string{"offer"}
cc := client.ClientCredentials(myscopes)
```

Get the access token:

```go
token, err := cc.AccessToken()
// do something with the token
```

## AccessToken
```go
type AccessToken struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
}
```