package provider

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrContent = errors.New("could not retrieve data from provided URL")
)

type HttpProvider struct {
	baseURL string
	client  *http.Client
}

func NewHTTP(baseURL string) *HttpProvider {
	return &HttpProvider{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

func (r *HttpProvider) Fetch(path string) (io.ReadCloser, error) {
	res, err := r.client.Get(fmt.Sprintf(r.baseURL, path))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrContent
	}

	return res.Body, nil
}
