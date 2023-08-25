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

type Remote struct {
	baseURL string
	client  *http.Client
}

func NewRemote(baseURL string) *Remote {
	return &Remote{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
}

func (r Remote) Audio(path string) (io.ReadCloser, error) {
	res, err := r.client.Get(fmt.Sprintf(r.baseURL, path))

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrContent
	}

	return res.Body, nil
}
