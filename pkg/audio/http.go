package audio

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrContent = errors.New("could not retreive data from provided URL")
)

type httpProvider struct {
	baseURL string
	client  *http.Client
}

func NewHTTP(baseURL string) *httpProvider {
	return &httpProvider{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

func (r httpProvider) Audio(path string) (io.ReadCloser, error) {
	res, err := r.client.Get(fmt.Sprintf(r.baseURL, path))

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrContent
	}

	return res.Body, nil
}
