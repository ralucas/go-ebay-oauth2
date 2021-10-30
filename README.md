Ebay OAuth2 for Go
---

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


### Client Credentials Flow
Example:
```go
myscopes := []string{"offer"}
cc := client.ClientCredentials(myscopes)
```