package double

import (
	"io"
	"net/http"
	"strings"

	"github.com/ralucas/go-ebay-oauth2/spy"
)

type SpyHTTPClient struct {
	spy.Spy
}

func (s *SpyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	s.Called(req)

	return &http.Response{
		Body: io.NopCloser(strings.NewReader("test")),
	}, nil
}
