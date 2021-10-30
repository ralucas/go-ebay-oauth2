package mocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/stretchr/testify/mock"

	oauth2 "github.com/ralucas/go-ebay-oauth2"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	body := string(reqBody)

	at, err := json.Marshal(oauth2.AccessToken{
		AccessToken: "test-token",
		ExpiresIn:   60,
		TokenType:   "bearer",
	})
	if err != nil {
		return nil, err
	}

	if strings.Contains(body, oauth2.GrantTypeAuthorizationCode) {
		return &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(at)),
		}, nil
	}

	if strings.Contains(body, oauth2.GrantTypeClientCredentials) {
		return &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(at)),
		}, nil
	}

	return nil, fmt.Errorf("bad request")
}
