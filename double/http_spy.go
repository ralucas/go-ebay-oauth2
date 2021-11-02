package double

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	oauth2 "github.com/ralucas/go-ebay-oauth2"
	"github.com/ralucas/go-ebay-oauth2/spy"
)

type SpyHTTPClient struct {
	spy.Spy
}

func (s *SpyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	s.Called(req)

	at, _ := json.Marshal(oauth2.AccessToken{
		AccessToken: "test-token",
		ExpiresIn:   60,
		TokenType:   "bearer",
	})

	res := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(at)),
	}

	return res, nil
}
