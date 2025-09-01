package openapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type HTTPURLLoader struct {
	client http.Client
}

func (l HTTPURLLoader) Load(url string) (any, error) {
	resp, err := l.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	return jsonschema.UnmarshalJSON(resp.Body)
}

func newHTTPURLLoader(insecure bool) *HTTPURLLoader {
	return &HTTPURLLoader{
		client: http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
		},
	}
}
